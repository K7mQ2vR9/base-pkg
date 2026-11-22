# 认证与匿名白名单

白名单配置、接口和匹配方式保持不变：`map[string]int`，值不参与判断，支持精确路径和以 `/` 开头、尾部 `*` 的字符串前缀匹配。例如 `/public` 不匹配 `/public/child`，`/open*` 仍匹配 `/open/item` 和 `/opener`。运行中不要并发修改传入的白名单 map。

白名单只豁免“必须登录”，请求头、正文、解密、认证和已接入的服务状态检查均按各自规则执行；不会因为路径在白名单中而吞掉错误。

## HTTP 认证

| 请求/回调结果 | 白名单接口 | 非白名单接口 |
|---|---|---|
| 未提供 Token（含空值/纯空白） | uid=0，匿名放行 | 返回未登录 |
| Token 回调返回 uid=0、err=nil | uid=0，匿名放行 | 返回未登录 |
| Token 回调返回 uid>0、err=nil | 注入用户身份后放行 | 注入用户身份后放行 |
| Token 无效、过期或回调返回错误 | 返回认证错误 | 返回认证错误 |
| 回调返回 uid=-1 | 返回无权限 | 返回无权限 |
| 其他负数 UID | 返回 Token 无效 | 返回 Token 无效 |
| 非空扩展信息不是合法 JSON | 返回服务不可用 | 返回服务不可用 |

没有 Token 时不调用 Token 校验回调；启动装配仍要求提供认证回调。回调须核实已提供的 Token，只有真实匿名结果才能返回 `0,nil`；**无效或过期 Token 必须返回 error**。中间件不能凭返回的 0 区分业务主动允许匿名与业务错误地忽略了校验失败。

`CheckRequestTokenFuncWithUserInfo` 支持用户扩展 JSON。匿名放行时显式写入 metadata 的 `token-uid="0"` 并清空 `headInfo.GetUserInfo(ctx)`，不会继承父上下文已有的用户 ID/扩展信息。正数 UID 只使用回调给出的身份及 JSON；原请求上下文不修改。非法 JSON 即使配合 uid=0 也拒绝，合法 JSON 在匿名时丢弃。

三种现有回调接口保留。仅返回整数的 `CheckRequestTokenFunc`：0 表示正常匿名、-1 无权限、其他负数无效；原来把过期/无效 Token 返回 0 的业务实现需要改为负数，或改用可返回 error 的接口。不要让真实认证失败表现为正常匿名。

错误响应沿用既有协议 HTTP200 与业务 code：未登录401、Token无效402、无权限405；特殊错误按原安全映射处理，不回显底层诊断。请求头400、正文413等入口错误仍保留各自 HTTP 状态。此处没有重定义业务码或全局错误映射。

## API 装配示例

```go
package service

import (
    "github.com/K7mQ2vR9/base-pkg/app"
    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "github.com/zeromicro/go-zero/rest"
)

// ConfigureAPI 在启动前装配；verify 验证 Token，serviceAvailable 检查服务状态。
// 两个回调均由业务提供，不能把无效 Token 返回为 0,nil,nil。
func ConfigureAPI(server *rest.Server,
    verify appMiddleware.CheckRequestTokenFuncWithUserInfo,
    serviceAvailable appMiddleware.CheckServiceStatusFunc) {
    white := map[string]int{
        "/public": 1,
        "/open/*": 1,
    }
    app.NewServerMiddleware(server,
        app.WithWhiteHeaderPathSMOption(white),
        app.WithCheckTokenHandleWithUserInfoSMOption(verify),
    ).ApiUseMiddleware()
    // 状态检查接在认证之后；匿名请求的 UID 为 0，Token 可能为空。
    server.Use(appMiddleware.NewCheckServiceStatusMiddleware(white,
        appMiddleware.WithCheckServiceStatusCheckOption(serviceAvailable),
    ).Handle)
}
```

server 使用 `app.NewServer` 创建，保留 Trace 和单值业务头校验。服务状态中间件仍需显式接入，`ApiUseMiddleware` 不会自动配置业务状态回调。`CheckServiceStatusRequest` 与构造函数继续接收原白名单参数，但不会据此跳过状态检查，缺少回调或返回 false 均拒绝。

`ApiHeaderMiddleware` 对 GET、白名单及其他请求统一调用 `Head.Verify`。原 `CloseVerifyOption`、`NotVerify` 字段及 `SetNoVerify` 调用入口保留，但不再绕过校验。当前 `Head.Verify` 没有额外业务必填规则，本次没有擅自要求 business/source 等字段必填；实际单值格式检查仍由前置 `RequestHeaderMiddleware` 执行。

有效 CORS 预检仍由 CORS 中间件返回204，不当作需要登录的业务调用；入口 Trace/请求头校验照常执行。自定义 net/http 装配需要显式接入这些入口组件，单独使用认证组件不会自动补齐正文、请求头或状态检查。

## RPC 调用方认证与用户登录

RPC 的处理顺序为：

```text
metadata/Trace → 业务字段校验 → 取消检查 → 调用方认证
→ 白名单/用户登录判断 → 业务方法
```

`CheckRPCRequestFunc` 对所有方法执行，包括白名单。该回调负责核实服务调用方及其携带的 Token/用户身份；不得只凭一个客户端传入的正数 UID 认定可信。回调通过后，`headInfo.GetTokenUid(ctx)>0` 可进入业务；UID缺失、空或0均视为匿名，仅白名单允许。重复/负数/非法 UID 先按已有 metadata 规则拒绝。后台无用户的服务任务应显式配置可匿名访问的方法。

携带错误 Token、调用方未通过认证、权限不足、取消或超时等错误，白名单同样返回错误。现有错误映射保留；非白名单匿名请求返回 gRPC Unauthenticated，由公共客户端恢复为业务401。白名单回调在调用方认证通过后调用，认证失败时不执行白名单及业务。

```go
package service

import (
    "context"

    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "google.golang.org/grpc"
)

// RPCAuthOption 接入所有方法均执行的调用方认证，并声明可匿名的方法。
func RPCAuthOption(check appMiddleware.CheckRPCRequestFunc) (grpc.ServerOption, error) {
    auth, err := appMiddleware.NewRpcAuthMiddleware(check)
    if err != nil {
        return nil, err
    }
    return auth.ServerOption(func(context.Context) map[string]int {
        return map[string]int{
            "/catalog.CatalogService/ListPublic": 1,
        }
    }), nil
}
```

原生 gRPC 使用 `grpc.NewServer(option)`，go-zero 使用 `server.AddOptions(option)`。同一 Server 只设置一次主 UnaryInterceptor。API 的用户扩展 JSON 不自动传入 RPC；用户 ID、语言及统一 Trace 继续沿用现有 metadata 协议。
