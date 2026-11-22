# 请求业务 metadata

API 完成请求头装配和认证后，将业务字段写入出站 metadata；RPC 通过 `RpcAuthMiddleware.ServerOption` 校验，再复制为当前请求的入站/出站 metadata。业务继续使用现有 `headInfo.Get...` 函数，参数和返回类型不变。

自定义字段已与 HTTP 一起统一为中划线，不读取下划线旧名。具体更名列表及客户端调整见 [HTTP 传输字段命名](HTTP_HEADERS.md)。原终端 `User-Agent` 通过 `client-user-agent` 传递，与 gRPC 保留的 `user-agent` 分开；业务 JSON 标签不变。

## 单值规则

以下字段每个最多一个值，**内容相同的重复值也拒绝**，不拼接、不任取首个或末个：

| 分类 | 字段 |
|---|---|
| 认证与身份 | authorization-jwt、token-uid、token-uid-role |
| 地址 | client-ip、admin-client-ip |
| 请求与来源 | req-path、version、source、client-user-agent、origin-url、device-id |
| 业务与显示 | business、business-code、content-language、timezone |
| RPC 诊断 | rpc-name |

可选字段缺失或只有一个空字符串时允许，以兼容 API 初始化时写入的空字段。存在字段但值列表为 nil/空列表，或包含两个以上值，均属于非法结构。

未知自定义字段保留多值能力，不修改或丢弃其值。`tracestate` 不属于业务单值字段，继续遵循原有 W3C 规则；`traceparent` 的必填、格式与重复规则仍由 [统一请求 ID 规则](TRACE_ID.md) 管理。

## UID 与 IP 格式

非空 `token-uid` 只接受 ASCII 十进制数字，范围 **0..9223372036854775807**，允许前导零。拒绝正负号、首尾空白、小数、科学计数法、十六进制及溢出；空值或 0 表示未设置身份。角色文本不在本层解析，单个字符串中的角色表示保持不变。

`client-ip` 和 `admin-client-ip` 非空值必须为纯 IPv4/IPv6 地址，拒绝端口、IPv6 zone、CIDR、地址列表、域名和首尾空白。Getter 返回规范化地址，IPv4-mapped IPv6（如 `::ffff:192.0.2.1`）输出 `192.0.2.1`；校验器本身不改写原 metadata。

其他字段只检查单值结构，不在本轮新增版本、语言、时区或业务文本约束，也不修剪文本。

## RPC 入口与错误

```text
非空 metadata → 有效 W3C Trace → ValidateMetadata → 上下文/日志装配
→ 取消检查 → 调用方认证 → 匿名白名单判断 → 业务
```

`headInfo.ValidateMetadata(md)` 只检查上述业务字段，不要求 metadata 或 Trace 必须存在，不能单独替代完整 RPC 入口。错误支持 `errors.Is(err, headInfo.ErrInvalidMetadata)`，文本只包含字段名与规则，不回显令牌、UID 或 IP 原文。

RPC 将非法业务 metadata 转换为 **gRPC InvalidArgument**，公开方法也不豁免。失败不执行白名单回调、认证回调或业务；公共 RPC 客户端按既有转换还原为业务 code400。

结构校验不证明调用方身份，也不证明 UID、角色或 IP 真实可信；原有认证回调、调用方信任与 HTTP 可信代理规则继续负责这些判断。

白名单允许用户 ID 为 0 的匿名请求，但不豁免调用方认证，也不限制已认证用户访问。缺失/空 UID 由 Getter 视为 0，非白名单拒绝；携带无效 Token 或调用方认证失败时，白名单同样返回错误。完整约定见 [认证与匿名白名单](../appMiddleware/AUTH.md)。

## Getter 与 Setter

- Getter 保持读取**出站** metadata，不自动回退入站字段；RPC 标准入口已完成复制。
- 缺失、空值、重复或非法 UID 返回 0；两种 IP 的非法或重复值返回空字符串；其他字符串字段重复时返回空字符串。nil ctx 按未设置处理。
- `GetTerminal` 继续将唯一 IOS/Android 映射为 APP；`GetOriginHostUrl` 保持既有协议文本移除规则，不在本轮改为 URL 解析。
- `SetTokenUid` 复制出站 metadata，`SetBusinessCode` 分别复制入站与出站 metadata 后覆盖单个值；不改原上下文、其他字段或 Trace，保留父级取消、截止时间和值。
- 两个 Setter 同步覆盖 LogKeys 中相应日志字段，HTTP 保持字符串、RPC 保持字符串列表；使用返回的新 ctx 记录后续日志，见 [日志同步说明](LOG_CONTEXT.md)。
- Setter 不校验格式或身份；写入非法 UID 后，Getter 返回 0，RPC 入口返回明确错误。

创建 metadata 使用 `metadata.Pairs` / `md.Set` 等标准 API，保持键名规范化。`ValidateMetadata` 对直接传入的原始 map 会合并大小写别名检查；不要手工构造混用大小写的重复键后再存入 gRPC 上下文，因为 gRPC 的上下文读取 API 已经会规范化键名，可能覆盖别名信息。

HTTP 入口的已知单值业务头现在按 [HTTP 请求头规则](HTTP_HEADERS.md) 校验，与 RPC 字段集合分别管理。登录用户扩展 JSON 仍不自动复制为 RPC metadata。

## 本地读写示例

以下 main 可直接运行，不连接外部服务，只演示结构和 Getter，不是认证实现。

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/K7mQ2vR9/base-pkg/consts"
    "github.com/K7mQ2vR9/base-pkg/headInfo"
    "google.golang.org/grpc/metadata"
)

func main() {
    ctx, err := headInfo.SetTrace(context.Background(), "1234567890abcdef1234567890abcdef")
    if err != nil {
        panic(err)
    }
    ctx = headInfo.SetTokenUid(ctx, "0012")
    md, _ := metadata.FromOutgoingContext(ctx)
    md.Set(consts.ClientIp, "::ffff:192.0.2.1")
    md.Set("custom-values", "a", "b")
    if err := headInfo.ValidateMetadata(md); err != nil {
        panic(err)
    }
    ctx = metadata.NewOutgoingContext(ctx, md)
    fmt.Println(headInfo.GetTokenUid(ctx), headInfo.GetClientIp(ctx))

    invalid := md.Copy()
    invalid.Set(consts.TokenUid, "1", "2")
    invalidCtx := metadata.NewOutgoingContext(ctx, invalid)
    fmt.Println(headInfo.GetTokenUid(invalidCtx),
        errors.Is(headInfo.ValidateMetadata(invalid), headInfo.ErrInvalidMetadata))
    fmt.Println(headInfo.GetTokenUid(ctx))
}
```

输出依次为 `12 192.0.2.1`、`0 true`、`12`：重复 UID 不会拼接成 12，原分支也未受影响。

## RPC 装配示例

下面返回包含完整入口校验的选项，认证回调和公开方法列表由业务提供。原生 gRPC 创建服务器时使用，go-zero 则调用 `server.AddOptions(option)`，与现有接入方式一致。

```go
package service

import (
    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "google.golang.org/grpc"
)

// MetadataAuthOption 在认证与业务之前校验 Trace 及已知业务 metadata。
func MetadataAuthOption(check appMiddleware.CheckRPCRequestFunc,
    publicMethods appMiddleware.PathWhiteHandle) (grpc.ServerOption, error) {
    auth, err := appMiddleware.NewRpcAuthMiddleware(check)
    if err != nil {
        return nil, err
    }
    return auth.ServerOption(publicMethods), nil
}
```

同一 Server 只设置一次主 UnaryInterceptor。正常下游调用继续沿用当前 ctx，客户端拦截器复制 metadata 并更新独立 rpc-name，不为同一链路重新生成 Trace ID。
