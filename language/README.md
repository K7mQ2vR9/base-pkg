# 项目本地多语言

翻译由 `github.com/K7mQ2vR9/base-pkg/language` 提供。各 API 项目维护自己的 `etc/language/*.json`，不再从 `language-pkg` 内嵌词典读取，也没有 `messageKeys` 中间映射。

## 目录和词典

```text
api/
  etc/
    app.yaml
    language/
      zh-CN.json
      en.json
```

`zh-CN.json` 的内容为扁平字符串对象：

```json
{
  "Internal Server Error": "服务器内部错误。",
  "Service Unavailable": "服务暂不可用。",
  "status.enabled": "已启用"
}
```

文件名就是语言标识。可复制 [完整示例词典](examples/etc/language) 到项目中，再增补业务词条。JSON 中的 key 可以是原文，也可以是业务 key；查找按完整字符串精确匹配，不去除文案空白、不替换子串。仅配置 `common.internal_error` 不会翻译原文 `Internal Server Error`，需要把实际返回的原文作为 key。

内置未登录、Token 无效及请求解密相关错误的公开原文已去除多余尾部空格，与示例词典一致：

| 原文 key | zh-CN 译文 |
|---|---|
| `User not login.` | 用户未登录。 |
| `Token invalid.` | 令牌无效。 |
| `Request decryption failed.` | 请求解密失败。 |
| `Request forbid.` | 请求被拒绝。 |

项目已有词典若为这些内置错误配置了带尾部空格的 key，应同步去除 key 的尾部空格。业务自定义文案和 data 字符串仍按原值精确匹配，框架不会自动 TrimSpace 或兼容旧文案。错误码及认证规则保持不变。

## 启动接入

`app.ConfigAppServer` 新增独立的 `Language language.Config` 字段，供项目配置嵌入。未填写 `Dir` 时默认使用**进程工作目录**下的 `etc/language`；部署到其他工作目录时配置绝对路径。

```yaml
Language:
  Dir: etc/language
```

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/K7mQ2vR9/base-pkg/language"
)

func main() {
    // 实际服务在配置读取后、启动 HTTP/RPC 之前调用 language.Load(c.Language)。
    // 未配置目录时使用默认值；加载失败必须停止启动，避免词典部署遗漏被忽略。
    if err := language.Load(language.Config{}); err != nil {
        panic(err)
    }
    fmt.Println(language.TranslateMessage("Internal Server Error", "zh-CN"))
    data := language.SwitchLanguage(map[string]any{
        "label": "status.enabled",
        "id": "status.enabled",
        "count": json.Number("9007199254740993"),
    }, "zh-CN")
    raw, err := json.Marshal(data)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(raw))
}
```

输出中文错误和 `{"count":9007199254740993,"id":"status.enabled","label":"已启用"}`。仅填写配置或调用 `app.InitAppServer()` 不会隐式读取词典；不需要翻译的进程可以不调用 `Load`，所有值原样输出。

`Load` 读取目录内扩展名为 `.json` 的文件，忽略子目录及其他文件。缺失/不可读目录、没有 JSON 词典、无效 JSON、非字符串值（含 null）、重复 key、规范化后的重复语言文件和尾随 JSON 都返回错误。空对象 `{}` 合法，表示该语言暂无词条。一次完整加载成功后才原子替换内存词典；失败保留上一次词典。文件修改不会自动生效，正常运维方式为更新文件后重启进程；请求中不读文件。词典为进程级配置，多租户独立词典不在此接口范围内。

## 请求与响应

客户端通过 `Content-Language: zh-CN` 选择语言，API 从请求上下文读取，RPC 沿用 `content-language` metadata。语言标识忽略大小写，`zh_CN` 与 `zh-CN` 等价；`en-US` 没有专用词典时回退到 `en`，其他语言不做猜测。未加载、没有语言、未知语言或缺少 key 时保留原文；词条值为空字符串时按配置输出空字符串。

标准 API 在入口单值业务头校验通过后立即注入语言，早于正文读取、并发控制、AES 解密及认证。各错误组件统一从 ctx 读取，因此 AES 解密失败也能翻译。请求语言缺失时不会沿用父 ctx 的旧语言；重复/大小写别名/空值列表拒绝且错误文案不任选语言。头校验失败但语言值唯一时仍可翻译400，完整规则见 [HTTP 请求头](../headInfo/HTTP_HEADERS.md)。

独立使用正文、并发或 AES 组件时，也应接入 RequestHeaderMiddleware。未注入语言的独立组件返回原文，不再自行回退原始语言头；仅调用 NormalizeHTTPHeaders 不会注入上下文。Trace 校验自身及 go-zero 原生组件自行生成的错误协议不在此次调整范围内。

正常业务继续调用 `result.HttpSuccessResult` / `HttpSuccessResultNotAes` / `HttpErrorResult`，无需指定字段清单或实现 Localize。错误 `message` 直接匹配词典；`data` 中匹配的字符串值递归自动翻译，字段名不翻译。正常成功响应的固定 `message="success"` 沿用现有协议。

`id`/`code`（大小写不敏感）下的直接字符串值受保护；其他业务枚举字符串如果恰好命中词典，也会被替换，因此业务词条应使用明确的命名空间。JSON 标签、omitempty、自定义序列化和数字文本保持编码语义，原业务对象不修改。同一次 data 编码捕获只读词典快照。`[]byte` 不翻译，string/[]byte 在 AES 响应中保持原始明文字节约定，结构化 data 先翻译再加密。

RPC 通常返回数据模型和业务 key，由 API 在最终输出阶段翻译；仅传递语言无需在 RPC 启动时加载词典。RPC 自身确实调用翻译接口时，才在自身项目配置并加载独立词典。

## 认证前解密错误示例

先将示例词典复制到项目 etc/language，并从项目根目录运行。下面仅用内存 HTTP 请求演示独立组件的错误翻译，不启动监听；实际 API 使用 app.NewServerWithCORS 和 ApiUseMiddleware 自动接入。

```go
package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/K7mQ2vR9/base-pkg/aesGCM"
    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "github.com/K7mQ2vR9/base-pkg/consts"
    "github.com/K7mQ2vR9/base-pkg/headInfo"
    "github.com/K7mQ2vR9/base-pkg/language"
    "github.com/K7mQ2vR9/base-pkg/result"
)

func main() {
    if err := language.Load(language.Config{}); err != nil {
        panic(err)
    }
    body, err := appMiddleware.NewRequestBodyMiddleware(appMiddleware.RequestBodyConfig{})
    if err != nil {
        panic(err)
    }
    // 演示密钥；实际服务从项目配置读取，仍由进程级开关控制加解密。
    decrypt, err := appMiddleware.NewApiRequestDecryptMiddleware(aesGCM.Config{
        Enabled: true, Key: "0123456789abcdef",
    })
    if err != nil {
        panic(err)
    }
    defer aesGCM.Configure(aesGCM.Config{})
    handler := appMiddleware.EnsureRequestTrace(
        appMiddleware.NewRequestHeaderMiddleware().Handle(
            appMiddleware.NewCorsMiddleware().Handle(body.Handle(decrypt.Handle(
                func(w http.ResponseWriter, r *http.Request) {
                    result.HttpSuccessResultNotAes(r.Context(), w, nil)
                },
            ))),
        ),
    )
    request := httptest.NewRequest(http.MethodPost, "http://example.test/items",
        strings.NewReader(`{"aes_data":"invalid-ciphertext"}`))
    request.Header.Set(headInfo.RequestIDHeader, "1234567890abcdef1234567890abcdef")
    request.Header.Set(consts.ContentLanguage, "zh-CN")
    response := httptest.NewRecorder()
    handler(response, request)
    fmt.Println(response.Code, response.Body.String())
}
```

输出 HTTP200、code400，message 为 `请求解密失败。`，trace 与传入的请求 ID 一致。缺少语言或未命中词典时保留 `Request decryption failed.`；该错误返回发生在后续请求头业务装配及认证之前。
