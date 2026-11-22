# HTTP 单值业务请求头

标准 API 使用 `app.NewServerWithCORS` 时，入口在 Trace 初始化后自动执行单值业务头校验，再进入原生中间件、正文读取、解密和认证。GET、白名单及其他方法使用同一套规则；白名单仅允许正常匿名结果，不跳过已提供 Token 的认证或请求头结构校验，详见 [认证与匿名白名单](../appMiddleware/AUTH.md)。

## 自定义字段统一使用中划线

HTTP 与 RPC 的自定义传输字段统一更名，Go 常量/Getter 名称不变，不读取或映射旧名称。前端、网关及手工构造 RPC metadata 的调用方须一起更新：

| 原字段 | 当前传输字段 |
|---|---|
| authorization_jwt | authorization-jwt |
| client_ip | client-ip |
| admin_client_ip | admin-client-ip |
| user_agent（原终端信息） | client-user-agent（RPC） |
| token_uid | token-uid |
| token_uid_role | token-uid-role |
| req_path | req-path |
| business_code | business-code |
| content_language | content-language |
| origin_url | origin-url |
| device_id | device-id |
| rpc_name | rpc-name |

version/source/business/timezone 本来就是单词，名称不变。标准 HTTP `User-Agent`、`Origin`、`X-Forwarded-For`、`X-Request-ID` 及 W3C 链路头保持各自协议；HTTP 名称仍大小写不敏感，规范化后可显示为 `Authorization-Jwt`、`Business-Code` 等。

RPC 原终端信息使用 `client-user-agent`，因为 `user-agent` 是 gRPC 自身的保留头，直接占用会被框架忽略或替换。当前 CORS 暴露头使用 `Authorization-Jwt`，日志业务字段键随共享传输常量同步。业务 JSON DTO、JWT claims 及响应里的字段名保持原样，例如 `business_code`、`user_agent`、`aes_data` 不因这次头改名而调整。

旧名称不作为新名称的别名。仅提供旧令牌头的受保护请求会按未提供令牌处理；其他旧头不会填入对应业务字段。未知头/metadata 仍可按原规则保留，不能据此认为旧业务字段仍然生效。

## 校验范围

| 分类 | 单值头名 |
|---|---|
| 项目认证令牌 | authorization-jwt |
| 业务标识 | business、business-code |
| 来源与版本 | source、version |
| 显示与设备 | content-language、timezone、device-id |
| 标准来源信息 | Origin、User-Agent |

头名大小写不敏感。同一字段两个以上值一律拒绝，内容相同也不例外；手工构造的大小写别名也合并计数。可选字段缺失或只有一个空字符串继续允许，不新增必填头或业务文本格式约束。存在字段但没有值列表（nil/空切片）拒绝。

成功后，下游使用独立的请求头副本，已知单值字段转成标准 MIME 大小写，解密、认证和请求头装配通过 `Header.Get` 读取同一份值。字段值不去空白、不改内容；后续组件原有的修剪/验证规则仍然生效。原请求的头 map、值切片和上下文不修改，校验器不读取或关闭正文。

校验成功后，RequestHeaderMiddleware 通过 `headInfo.ContextWithRequestLanguage(r)` 将唯一的 Content-Language 提前写入新 ctx 的出站 metadata。正文、并发、解密及统一响应均通过 `GetContentLanguage(ctx)` 读取，不再各自选择原始语言头。仅裁剪原 Head 读取规则中的首尾空格；不读正文、不注入身份、不改变 AES 配置或其他 metadata。

缺少/空语言写入空值，覆盖父上下文的旧语言；重复、大小写别名及空值列表均拒绝且不任选语言，错误响应使用原文。其他头非法但语言头唯一时，该400响应仍可翻译。未知语言、缺少词条或未加载词典按原规则返回原文。语言注入发生在 Trace 成功之后，Trace 错误协议保持原样。

本项目实际读取的令牌头是 `authorization-jwt`，本轮不新增其他令牌别名。UID、角色和客户端 IP 仍从认证及可信代理逻辑获得，不将客户端随意提供的同名头直接作为身份或地址。

后续 ApiHeaderMiddleware 将 Head 写入出站 metadata 时，保留前置中间件已写入的自定义字段及多值，平台字段则按当前请求覆盖（空值也覆盖），避免继承旧用户身份。具体覆盖和W3C边界见 [metadata合并规则](METADATA.md)。仅显式写入ctx出站metadata的自定义字段会保留，不会自动转发所有HTTP头。

## 错误响应与执行顺序

```text
Trace 校验/初始化 → 单值业务头校验并提前注入语言 → 原生声明长度预检查
→ go-zero 原生中间件 → CORS → 可选并发控制 → 请求体限制
→ AES 解密 → 请求头上下文装配 → 认证 → 业务/RPC
```

校验放在路由外层，避免原生 gzip 在校验前读取正文。非法业务头返回 **HTTP400/code400**，使用现有 `{code,message,data,trace}` JSON，保留同一 `X-Request-ID` 和 CORS；公开消息走统一翻译入口，没有对应翻译时使用标准 Bad Request。不输出原始头值或令牌。

Trace 自身的非法/重复/冲突规则先执行，仍遵循 [统一请求 ID 规则](TRACE_ID.md)。HTTP 单值头失败不进入框架解压、正文读取、解密、认证或业务。有效 CORS 预检仍正常返回204，`Access-Control-Request-Headers` 中的逗号头名列表不受本规则限制；预检若实际附带重复的单值业务头，则同样拒绝。

`headInfo.NormalizeHTTPHeaders(headers)` 可单独使用：成功返回深复制且规范化的头，失败返回可用 `errors.Is(err, headInfo.ErrInvalidHTTPHeaders)` 判断的错误。错误只描述固定字段名与规则。该函数不初始化 Trace，也不执行身份认证。

`ContextWithRequestLanguage` 只提取语言并返回新 ctx，不能替代 `NormalizeHTTPHeaders`；标准入口已自动调用，业务无需重复注入。自定义入口使用它时，须保留完整单值头校验。HTTP 成功头校验会以本次请求的语言覆盖旧语言，不从入站 RPC metadata 猜测语言。

## X-Forwarded-For 和其他多值头

`X-Forwarded-For` 不属于单值集合，多行和逗号分隔的代理链原样保留，由既有可信代理逻辑处理：

1. 连接对端不可信时，忽略转发链，使用连接对端 IP。
2. 连接对端可信时，从右向左跳过可信代理，取第一个不可信 IP。
3. 没有转发链或检查到非法地址时，按既有规则回退到连接对端。

因此 `伪造IP, 用户IP, CDN IP` 在 CDN 与连接对端均可信时，会在用户 IP 处停止，不采用更左侧伪造值。代理应正确追加实际连接对端，可信名单只配置实际代理网段。

Accept、Cookie、自定义多值头及其他不在集合中的字段也不被本规则拒绝；它们继续由各自协议处理。校验只检查 HTTP 解析后可见的值条目，不自动把单行逗号文本拆成多个字段。若上游已将重复行合并成一段文本，本层无法恢复原始行数；业务值格式和令牌合法性仍由相应校验/认证负责。

## 标准与独立接入

平台标准入口自动启用单值头校验，无新增单值校验配置或白名单开关。`ApiUseMiddleware` 继续装配正文限制、AES 和认证，不需要再手动追加一次头校验。默认 CORS 允许任意来源且不启用凭证；需要配置来源或凭证时，统一通过 `NewServerWithCORS` 接入，入口错误沿用同一策略，见 [CORS接入](../app/CORS.md)。链路路由由构造函数内部装配，避免用裸 Router 覆盖标准入口。

独立 net/http 接入时，显式保持下面的顺序：

```go
package api

import (
    "net/http"

    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
)

// HeaderHandler 在正文和业务之前完成 Trace 与单值头校验。
func HeaderHandler(next http.HandlerFunc) (http.HandlerFunc, error) {
    body, err := appMiddleware.NewRequestBodyMiddleware(appMiddleware.RequestBodyConfig{})
    if err != nil {
        return nil, err
    }
    headers := appMiddleware.NewRequestHeaderMiddleware()
    cors := appMiddleware.NewCorsMiddleware()
    return appMiddleware.EnsureRequestTrace(
        headers.Handle(cors.Handle(body.Handle(next)))), nil
}
```

如需 AES，在正文限制之后、业务之前装配解密器。单独调用旧请求头/认证/解密组件，不会自动补齐所有前置中间件，应使用标准入口或上述组合。

## 规范化与错误分类示例

下面 main 不启动服务或连接网络，也不打印令牌：

```go
package main

import (
    "errors"
    "fmt"
    "net/http"

    "github.com/K7mQ2vR9/base-pkg/consts"
    "github.com/K7mQ2vR9/base-pkg/headInfo"
)

func main() {
    original := http.Header{
        "business-code": {"business"},
        "X-Forwarded-For": {"203.0.113.20", "198.51.100.10"},
    }
    normalized, err := headInfo.NormalizeHTTPHeaders(original)
    if err != nil {
        panic(err)
    }
    fmt.Println(normalized.Get(consts.BusinessCode), len(normalized.Values("X-Forwarded-For")))
    invalid := http.Header{"Authorization-Jwt": {"first", "second"}}
    _, err = headInfo.NormalizeHTTPHeaders(invalid)
    fmt.Println(errors.Is(err, headInfo.ErrInvalidHTTPHeaders))
    normalized.Set(consts.BusinessCode, "branch")
    fmt.Println(original["business-code"][0])
}
```

输出为 `business 2`、`true`、`business`。合法 API 请求后续仍转为单值 [RPC metadata](METADATA.md)，保留既有认证身份和同一 Trace ID。
