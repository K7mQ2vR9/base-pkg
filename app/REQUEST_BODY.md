# API 请求体读取与大小限制

`RequestBodyMiddleware` 在解密和业务读取之前，按实际字节数限量读取请求体，成功后提供 `Body` 和可重复调用的 `GetBody`。请求头中间件只装配头信息和上下文，不读取正文。

## 配置与默认值

`ConfigAppServer.RequestBody` 对应 `appMiddleware.RequestBodyConfig`：

| 字段 | 单位 | 默认值与校验 |
|---|---|---|
| MaxBytes | 字节 | 0 或未配置使用 1 MiB（1048576），正数覆盖上限，负数拒绝；不提供无限制模式 |

通过 `app.WithRequestBody(config)` 校验并接入；不传此选项，`ApiUseMiddleware` 仍自动启用默认限制。

上限作用于中间件实际看到的完整正文，包含 AES 请求中的 JSON、Base64 密文、nonce 和认证标签，不只是业务明文。超过上限的声明长度可提前拒绝；未知长度、chunked 或虚报较小长度仍会按实际读取量检查，最多从当前流读取上限加一个判定字节，不把截断正文交给业务。

使用 go-zero 的 Gunzip 中间件时，本中间件位于它之后，解压后的正文也受上限约束。网关仍可独立限制传输体积。大小上限不等于读取时间上限，应同时设置服务器读取超时；它也不直接限制 JSON 字段数量或业务解析复杂度。

## 中间件顺序

标准接入使用 `app.NewServer` 和 `NewServerMiddleware(...).ApiUseMiddleware()`：

```text
入口 Trace 校验/初始化 → 单值业务头校验并注入语言 → 原生声明长度预检查 → go-zero 原生中间件
→ CORS → 并发控制（若启用） → 请求体限量读取 → AES 解密 → 请求头装配 → 认证 → 业务
```

并发名额涵盖正文读取与解密，不会在等待名额之前预先缓存完整正文。CORS 预检在进入请求体中间件前结束。若上游原生中间件提前拒绝，请求不会进入后续链路。

普通正文只进行一次完整限量读取。AES 解密复用尚未消费的缓存字节，解密后重建 `Body`、`GetBody`、`ContentLength` 和 `Content-Length`，清除旧的 `TransferEncoding`/`Transfer-Encoding`。正文的二进制字节、JSON 数字精度及空白不会由本中间件重新编码。

## 错误响应

| 情况 | HTTP 状态 | JSON code |
|---|---|---|
| 声明长度或实际读取超限 | 413 | 413 |
| 正文读取失败，包括截断或读取层返回取消/超时错误 | 400 | 400 |
| 正文读取成功，但 AES 密文/封装不符合既有协议 | 沿用原有 HTTP 200 业务错误协议 | 400 |

请求体错误使用现有 `{code, message, data, trace}` 结构，`trace` 与 `X-Request-ID` 一致，保留 CORS 响应头。公开消息经过统一多语言入口，未配置翻译时使用标准状态文案，不输出原始正文、密文或底层读取错误。前置 Trace 校验失败仍使用原有链路错误协议。

正文错误只从上下文读取语言，不再回退原始 HTTP 语言头。独立 net/http 必须先接入 RequestHeaderMiddleware（或在完成请求头校验后显式注入语言），再进入正文/并发/解密组件；单独使用正文组件而未注入语言时返回原文。完整接入见 [HTTP 入口说明](../headInfo/HTTP_HEADERS.md)。

没有读到完整正文时不进入认证或业务。空正文、nil Body 和 `http.NoBody` 可继续处理；具体业务是否允许空正文，由接口解析规则决定。

## 与 go-zero / 网关上限的关系

go-zero v1.7.4 的原生 `MaxBytesHandler` 只比较 `Content-Length`，不能覆盖未知长度读取。平台请求体中间件始终检查实际字节数。

`app.NewServer` 根据 `RestConf.Middlewares.MaxBytes` 和 `RestConf.MaxBytes`，在 Trace 初始化之后预检查标准服务器声明长度，避免原生默认检查提前返回空 413。`WithRequestTraceRouter` 自定义接入时，第三个参数填写对应原生上限；原生检查关闭时传 0，省略则按 1 MiB 预检查。

两套配置及网关的更小上限仍可能先拒绝。**调大 RequestBody.MaxBytes 时，需要同步调整原生/网关限制。** 推荐按下面示例关闭原生 MaxBytes，统一交由平台配置控制。自定义原生路由限制（如 `rest.WithMaxBytes`）或网关提前拒绝时，其响应形式由相应层负责；要求统一外观时使用平台全局配置。

## 标准 API 配置示例

以下将上限设为 2 MiB，关闭原生声明长度检查，避免保留另一个更小上限。其余原生中间件按 go-zero 默认配置加载。

```yaml
Name: api-service
Host: 0.0.0.0
Port: 8888
Middlewares:
  MaxBytes: false
RequestBody:
  MaxBytes: 2097152
```

对应启动代码如下。认证逻辑由业务注入，这里只演示配置和中间件装配，范例本身不启动监听。

```go
package api

import (
    "net/http"

    "github.com/K7mQ2vR9/base-pkg/app"
    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "github.com/zeromicro/go-zero/rest"
)

// Config 将 HTTP 原生配置与平台配置一起加载。
type Config struct {
    rest.RestConf
    app.ConfigAppServer
}

// NewHTTPServer 在启动时校验正文上限，并按标准顺序装配请求链。
func NewHTTPServer(c Config, check appMiddleware.CheckRequestTokenFunc) (*rest.Server, error) {
    bodyOption, err := app.WithRequestBody(c.RequestBody)
    if err != nil {
        return nil, err
    }
    server, err := app.NewServer(c.RestConf)
    if err != nil {
        return nil, err
    }
    app.NewServerMiddleware(server,
        bodyOption,
        app.WithCheckTokenHandleSMOption(check),
    ).ApiUseMiddleware()
    server.AddRoute(rest.Route{
        Method: http.MethodPost,
        Path: "/items",
        Handler: func(w http.ResponseWriter, r *http.Request) {
            // 此时正文已通过上限检查，业务自行解析和处理。
            w.WriteHeader(http.StatusNoContent)
        },
    })
    return server, nil
}
```

如果启用 AES，再将 `app.WithAESGCM(c.AESGCM)` 返回的选项传给同一个 `NewServerMiddleware`；`WithRequestBody` 的配置会自动同步给其内部解密中间件。AES 仍按进程级配置开关工作，不通过上下文传递加密器。

## 独立 net/http 接入

不使用 `ApiUseMiddleware` 时，手动保证 Trace、单值业务头校验、限量读取、解密及业务的顺序。独立解密中间件自身也有读取保护，省略可选 RequestBodyConfig 默认 1 MiB；关闭 AES 时的普通请求限制由前置请求体中间件负责。

```go
package api

import (
    "net/http"

    "github.com/K7mQ2vR9/base-pkg/aesGCM"
    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
)

// BodyHandler 在业务前同时装配正文限制与可选 AES 解密。
func BodyHandler(config appMiddleware.RequestBodyConfig, cipher aesGCM.Config,
    next http.HandlerFunc) (http.HandlerFunc, error) {
    body, err := appMiddleware.NewRequestBodyMiddleware(config)
    if err != nil {
        return nil, err
    }
    decrypt, err := appMiddleware.NewApiRequestDecryptMiddleware(cipher, config)
    if err != nil {
        return nil, err
    }
    headers := appMiddleware.NewRequestHeaderMiddleware()
    return appMiddleware.EnsureRequestTrace(headers.Handle(body.Handle(decrypt.Handle(next)))), nil
}
```

独立构造解密中间件会应用进程共用 AES 配置，应在服务初始化阶段完成。中间件先缓存完整的有界正文，再执行业务；上传或流式接口应使用独立处理链并配置自己的大小与时间限制。
