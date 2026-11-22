# CORS 来源与凭证配置

默认允许任意来源，不输出 `Access-Control-Allow-Credentials`。前端以自定义 `authorization-jwt` 请求头携带 Token 时，通常不需要启用 `credentials: "include"` 或 `withCredentials: true`。

需要跨域 Cookie／浏览器凭证请求时，显式开启凭证并配置精确的来源名单。凭证请求不能使用 `Access-Control-Allow-Origin: *`，见 [Fetch 标准](https://fetch.spec.whatwg.org/#cors-protocol-and-credentials)。

## 配置和默认值

`app.ConfigAppServer.CORS` 使用独立的 `appMiddleware.CORSConfig`：

```yaml
CORS:
  AllowCredentials: true
  AllowedOrigins:
    - https://app.example.com
    - http://localhost:3000
```

| 字段 | 默认值与规则 |
|---|---|
| AllowCredentials | false，不输出凭证授权头；true 时 AllowedOrigins 必须非空 |
| AllowedOrigins | 空列表且凭证关闭时允许任意来源；非空时仅允许精确匹配来源 |

来源只允许 `http://host[:port]` 或 `https://host[:port]`，不带尾部斜杠、路径、查询串、片段或用户信息；禁止 `*`、`null` 和通配子域。端口必须有效。匹配区分大小写及显式端口，填写浏览器实际发送的序列化 Origin，例如默认 HTTPS 端口通常写为 `https://app.example.com`。配置非法返回初始化错误，不静默回退为开放来源。构造时复制配置，服务运行中不会共享调用方可变切片。

## API 服务接入

API 服务使用 `app.NewServerWithCORS(restConf, c.CORS)`；仅填写配置字段不会自动启用，传入构造函数后才生效。`c.CORS` 未配置时使用零值，默认允许任意来源且关闭跨域凭证。原有 ApiUseMiddleware 中的默认 CORS 组件会沿用外层已绑定策略，不会覆盖为任意来源。

```go
package api

import (
    "github.com/K7mQ2vR9/base-pkg/app"
    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "github.com/zeromicro/go-zero/rest"
)

// BuildAPI 在启动前校验CORS配置；认证回调仍由项目业务实现。
func BuildAPI(restConf rest.RestConf, config app.ConfigAppServer,
    verify appMiddleware.CheckRequestTokenFuncWithErr) (*rest.Server, error) {
    server, err := app.NewServerWithCORS(restConf, config.CORS)
    if err != nil {
        return nil, err
    }
    app.NewServerMiddleware(server,
        app.WithCheckTokenHandleWithErrSMOption(verify),
    ).ApiUseMiddleware()
    return server, nil
}
```

服务策略在最外层绑定，正常、Trace/请求头错误、正文限制、AES/认证/并发错误，以及404/405都使用同一来源规则。Trace和头校验仍先于预检放行；声明长度检查后，有效 OPTIONS 在路由入口返回204（含未注册路径），不进入正文、认证或业务。普通请求的业务处理顺序不变。

每个 Server 独立保存只读策略，不设置全局变量。API 服务统一通过 `app.NewServerWithCORS` 接入来源和凭证配置，链路路由由包内私有函数装配。构造函数不接受 `rest.RunOption`，不能再传入 `rest.WithRouter`、`rest.WithCors` 或 `rest.WithCustomCors` 覆盖统一入口。

### 自定义 404 / 405

可选的第三个参数为 `app.ServerHandlerConfig`，仅在代码中配置，不放入 YAML。最多传一项，省略、零值或对应字段为 nil 时沿用 go-zero 默认处理；传入多项返回初始化错误。现有两参数调用方式保持。

```go
// notFound 和 notAllowed 为业务提供的 http.Handler；不需要自定义的字段可留空。
server, err := app.NewServerWithCORS(restConf, corsConfig, app.ServerHandlerConfig{
    NotFoundHandler: notFound,
    NotAllowedHandler: notAllowed,
})
if err != nil {
    return err
}
```

`NotFoundHandler` 处理未匹配路径，默认返回原生404；`NotAllowedHandler` 处理路径存在但方法不匹配的请求，默认返回原生405及 `Allow` 头。自定义处理器自行写入状态码和正文，并负责需要的响应头（包括自定义405的 `Allow`），不要覆盖入口已设置的 CORS 和 `X-Request-ID`。处理器可能并发执行，共享状态需由业务保证并发安全。

回退处理器只在统一 Trace、单值头校验、声明长度预检查和 CORS 处理之后执行，不能替换路由或绕过上述入口。有效 OPTIONS 仍返回204，不执行自定义404/405。回退路径不经过已注册业务路由的认证等应用中间件，不应在这些处理器中放置需要登录认证的业务操作。

旧的 `rest.WithNotFoundHandler` / `rest.WithNotAllowedHandler` 参数改为以上结构字段；其余原生运行选项不再透传，无兼容分支。

## 授权头和预检行为

| 情况 | Allow-Origin | Allow-Credentials |
|---|---|---|
| 默认模式 | `*` | 不输出 |
| 名单内来源、凭证开启 | 当前 Origin | `true` |
| 名单内来源、凭证关闭 | 当前 Origin | 不输出 |
| 名单外／缺少来源，且配置了名单 | 不输出 | 不输出 |

名单模式补充 `Vary: Origin`，保留原有 Vary 并大小写不敏感去重。缺失、重复或大小写别名的 Origin 不会任选一个来源授权。

继续开放所有预检请求头：将 `Access-Control-Request-Headers` 中列出的名称返回到 `Access-Control-Allow-Headers`，包括 `Authorization`、`Authorization-Jwt` 和业务自定义头，并添加 `Vary: Access-Control-Request-Headers`。未提供头列表时沿用 `*`。方法保持 GET/POST/PUT/DELETE/OPTIONS/PATCH，统一请求 ID 及项目令牌响应头继续可被浏览器读取。

名单外来源不新增服务端403，正常请求仍按原认证和业务规则处理，只是不授予浏览器跨域读取权限；CORS 不替代调用方认证或 CSRF 防护。Cookie 是否发送还取决于浏览器及 Cookie 自身属性。普通同源或服务端调用不以 CORS 头作为认证依据。

## 独立 net/http 接入与可运行范例

使用 `NewCorsMiddlewareWithConfig` 校验策略，`Bind` 放在 Trace 之前以覆盖早期错误，`Handle` 放在 Trace/单值头校验之后处理预检。Bind 不读正文或提前结束 OPTIONS；后续错误出口从当前请求上下文取同一策略。

```go
package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"

    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "github.com/K7mQ2vR9/base-pkg/headInfo"
)

func main() {
    cors, err := appMiddleware.NewCorsMiddlewareWithConfig(appMiddleware.CORSConfig{
        AllowCredentials: true,
        AllowedOrigins: []string{"https://app.example.com"},
    })
    if err != nil {
        panic(err)
    }
    handler := cors.Bind(appMiddleware.EnsureRequestTrace(
        appMiddleware.NewRequestHeaderMiddleware().Handle(cors.Handle(
            func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            },
        )),
    ))
    request := httptest.NewRequest(http.MethodOptions, "http://api.example.test/items", nil)
    request.Header.Set(headInfo.RequestIDHeader, "1234567890abcdef1234567890abcdef")
    request.Header.Set("Origin", "https://app.example.com")
    request.Header.Set("Access-Control-Request-Headers", "Authorization-Jwt, X-Custom")
    response := httptest.NewRecorder()
    handler(response, request)
    fmt.Println(response.Code)
    fmt.Println(response.Header().Get("Access-Control-Allow-Origin"))
    fmt.Println(response.Header().Get("Access-Control-Allow-Credentials"))
    fmt.Println(response.Header().Get("Access-Control-Allow-Headers"))
}
```

输出204、`https://app.example.com`、`true`、`Authorization-Jwt, X-Custom`。这演示请求和响应头，不是浏览器Cookie联调。

原先依赖跨域凭证的客户端须同步传入来源配置；默认模式将不再输出旧的凭证授权头。仍通过自定义Token头访问的客户端可继续使用默认配置，并关闭跨域凭证选项。
