# API层架构与公共包接入

API层接收用户操作、编排业务流程、调用RPC数据模型服务并构建用户响应。数据模型和数据连接归对应RPC服务管理；API通过服务契约协作。公共能力和配置完整说明见 [base-pkg汇总](../README.md)。

本目录是架构与接入说明，不表示仓库中已有业务API完成了新模块迁移，也不定义首批业务接口。

## 启动与装配顺序

1. 读取项目配置，显式加载本地词典，初始化日志及带ctx的框架错误处理器。
2. 创建服务级RPC客户端和该API实际使用的其他共享资源，保存关闭入口。
3. 用 `app.NewServerWithCORS` 创建Server，使用一份 `ServerMiddleware` 装配配置及认证。
4. 需要服务状态检查时，先成功构造其组件，再在认证之后追加；注册业务路由并启动服务。
5. 退出时停止新请求和任务，等待或取消在途操作，再关闭下游客户端和共享资源。

`app.InitAppServer()` 会初始化日志并注册 `ErrHandlerCtx`，不会加载词典或建立数据库连接。`language.Load` 读取项目 `etc/language`，加载失败须按启动错误处理。初始化失败时释放此前成功创建的资源，停止继续启动。

## 标准Server构造范例

以下函数装配已有配置，不启动监听；调用方负责先初始化词典、日志和共享资源。认证回调须核实Token，无效Token返回error，不能以 `0,nil,nil` 表示匿名成功。

```go
package service

import (
    "github.com/K7mQ2vR9/base-pkg/app"
    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "github.com/zeromicro/go-zero/rest"
)

// BuildAPI 创建并装配当前服务，white为匿名路径，verify为业务认证回调。
// 配置零值沿用公共包默认值；返回前不注册业务路由或启动监听。
func BuildAPI(restConf rest.RestConf, cfg app.ConfigAppServer,
    white map[string]int,
    verify appMiddleware.CheckRequestTokenFuncWithUserInfo) (*rest.Server, error) {
    proxy, err := app.WithTrustedProxyCIDRs(cfg.TrustedProxyCIDRs)
    if err != nil {
        return nil, err
    }
    concurrency, err := app.WithMaxConcurrency(cfg.MaxConcurrency)
    if err != nil {
        return nil, err
    }
    body, err := app.WithRequestBody(cfg.RequestBody)
    if err != nil {
        return nil, err
    }
    encryption, err := app.WithAESGCM(cfg.AESGCM)
    if err != nil {
        return nil, err
    }
    options := []app.SMOption{
        proxy, concurrency, body, encryption,
        app.WithWhiteHeaderPathSMOption(white),
        app.WithCheckTokenHandleWithUserInfoSMOption(verify),
    }
    if cfg.IsDebug {
        options = append(options, app.WithDebugOption())
    }
    server, err := app.NewServerWithCORS(restConf, cfg.CORS)
    if err != nil {
        return nil, err
    }
    assembly := app.NewServerMiddleware(server, options...)
    if err := assembly.ApiUseMiddleware(); err != nil {
        server.Stop()
        return nil, err
    }
    return server, nil
}
```

字段只有通过对应选项或构造函数传入才生效，不能只填写ConfigAppServer。默认正文1MiB、并发0关闭、AES关闭、无可信代理、CORS任意来源且关闭凭证。正文上限与网关/go-zero原生限制协调；开启跨域Cookie时配置精确来源。

同一装配实例成功后重复调用返回错误；不跨实例防重，不按值复制或在运行时重新装配。完整入口、404/405和CORS规则见 [Server接入](../app/CORS.md)，请求大小见 [正文限制](../app/REQUEST_BODY.md)。

## 请求上下文与身份

- 网关/HTTP入口在原生Span创建前校验并初始化统一Trace，已有合法X-Request-ID或traceparent时复用；业务不自行覆盖ID。
- 自定义头使用 `authorization-jwt`、`business-code`、`content-language` 等中划线名称。请求头单值校验不等于身份认证。
- IP由连接对端及已配置的可信代理链解析，不采信用户随意指定的client-ip。代理/CDN需按实际部署整理可信来源和转发链。
- 无Token或回调确认的UID=0只在白名单放行；提供了无效/过期Token时仍返回认证错误。
- 认证后用 `headInfo.GetTokenUid(ctx)`、`GetUserInfo(ctx)` 读取当前身份；扩展JSON仅在当前API请求内使用，不自动发送给RPC。
- 下游调用传当前ctx。修改业务标识使用 `headInfo.SetBusinessCode` 的返回上下文，避免丢弃新分支；不要用Background替换请求ctx。

服务状态组件、不同认证回调及RPC调用方认证详见 [认证说明](../appMiddleware/AUTH.md)。

## RPC调用与用户响应

下游RPC连接在启动时通过 `rpc.NewClient(startupCtx, config)` 创建，使用生成客户端的zrpc.Client接口或 `Conn()`；每个请求沿用已有ctx，可在业务边界进一步收紧预算。默认一元调用10s，父请求更早截止时间优先。不要对结果不确定的写操作盲目重试。

成功响应选择 `result.HttpSuccessResult` 或 `HttpSuccessResultNotAes`，错误选择 `result.HttpErrorResult` 或框架 `httpx.ErrorCtx`，输出后立即返回。框架只注册带上下文的 `ErrHandlerCtx`，无上下文旧入口已删除。

响应沿用 `code/message/data/trace`。业务错误一般HTTP200加业务code，入口400/413/503等保持自身HTTP状态。RPC业务码由公共客户端恢复，API不按驱动错误字符串判断、不输出底层响应正文。

语言从ctx读取，词典从项目etc目录显式加载；data匹配的字符串值自动翻译，JSON字段名不翻译。AES按进程配置控制成功data，错误仍为普通JSON。完整范例见 [响应](../result/README.md) 和 [本地词典](../language/README.md)。

## 异步业务与架构边界

API发布消息时保存实际消息ID，结果不确定时重发沿用该ID，ctx继续携带当前Trace。消息ID可以是业务自增ID的字符串，但不能代替业务幂等或Trace ID。批次失败处理已确认前缀，避免整批重发导致重复效果。

如API服务承担独立任务，RabbitMQ回调用其任务ctx；XXL-JOB的Start用服务生命周期ctx，不能用单次用户请求ctx启动长期执行器。任务处理、数据库事务、消息持久化和恢复流程由业务定义，不由HTTP请求返回成功自动保证。

本指南保留当前接入约定，实际权限、首批接口、服务发现、密钥来源、部署和联调由业务项目确定。未增加全局panic/重复写保护；已跳过的调度回调重试、名称校验及独立中间件构造改造不视为现有能力。

## 维护与验证状态

本指南替代本目录对公共包修改的逐日镜像流水。后续接口变化同步更新本指南和组件说明，具体开发仍遵循 [AGENTS.md](DEVELOPMENT_RULES.md)。本次只整理文档，没有修改其他API项目源码、编译范例或执行测试；公共包近期改动也不能据此视为已完成当前环境验证。
