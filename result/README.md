# API 响应输出

`result` 为 API 输出统一响应，RPC 数据模型服务继续返回自身协议，由 API 组织业务结果。

| 入口 | 行为 |
|---|---|
| HttpSuccessResult(ctx, w, data) | 自动翻译匹配值，按进程配置决定是否加密成功 data |
| HttpSuccessResultNotAes(ctx, w, data) | 自动翻译匹配值，始终输出普通成功 JSON |
| HttpErrorResult(ctx, w, err) | 沿用现有业务错误分类和公开文案，输出普通错误 JSON |
| HttpErrorResponse(ctx, err) | 仅构建 HTTP 状态与响应对象，供 go-zero 全局错误处理器使用 |

公开函数的参数和返回类型保持不变。前三个输出入口共用内部编码/写入流程；直接调用其他框架输出函数或自行写正文，不在该流程范围内。

## 响应结构与协议

```json
{"code":200,"message":"success","data":{},"trace":"1234567890abcdef1234567890abcdef"}
```

`trace` 来自当前上下文，沿用入口初始化的统一 ID，与 `X-Request-ID` 一致；本包不在响应阶段另行生成 ID。使用 [标准 API 入口](../headInfo/TRACE_ID.md) 初始化链路，业务始终传当前请求 ctx。

正常成功保持 HTTP200/code200。普通 `data=nil` 保持空对象 `{}`；类型化 nil、字节切片、字符串等继续遵循原有 JSON 行为。`message="success"` 的既有协议不变；data 中匹配的字符串值自动翻译，字段名不翻译，`id`/`code` 保护规则保留。

翻译来源为项目本地 `etc/language/*.json`，启动时必须显式加载；未加载或未命中时保留原值。错误 message 直接按原文或业务 key 查找，data 的字符串值自动匹配，无 `messageKeys` 映射。接入方式见 [本地词典说明](../language/README.md)。

## 错误处理器接入

标准启动调用 `app.InitAppServer()` 已通过 `httpx.SetErrorHandlerCtx` 注册 `logs/xcode.ErrHandlerCtx`。只需单独注册 HTTP 错误处理器时，在启动阶段调用下列配置，不必重复初始化日志：

```go
package api

import (
    errorhandler "github.com/K7mQ2vR9/base-pkg/logs/xcode"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// InitHTTPErrorHandler 在接收请求前注册进程共用的错误处理器。
func InitHTTPErrorHandler() {
    httpx.SetErrorHandlerCtx(errorhandler.ErrHandlerCtx)
}
```

业务错误出口使用 `httpx.ErrorCtx(r.Context(), w, err)`，也可直接调用 `result.HttpErrorResult(r.Context(), w, err)`，始终传递当前请求上下文。`ErrHandlerCtx` 仅构建状态和响应对象，由框架写入；它使用上下文中的统一Trace ID和语言，不生成新ID，也不从原始请求头补取数据。链路和语言仍由标准请求入口提前注入，项目词典仍须显式加载。

原无上下文的 `logs/xcode.ErrHandler` 已删除，不提供兼容包装。原 `httpx.SetErrorHandler(errorhandler.ErrHandler)` 调用须迁移为 `httpx.SetErrorHandlerCtx(errorhandler.ErrHandlerCtx)`；不要使用 `context.Background()` 或 `context.TODO()` 替换请求上下文。标准 `app.InitAppServer()` 的注册方式保持，错误的HTTP状态、业务码、公开文案和响应结构沿用现有规则。

## 编码失败

输出前先完整编码，再设置 JSON Content-Type、提交状态码和写入正文。NaN/Inf、循环引用、不可编码类型、无效 RawMessage、自定义 MarshalJSON 返回错误或无效 JSON，都会在提交成功响应之前被发现。

编码或 AES 准备失败时，固定返回现有内部错误协议：**HTTP200、code500、data={}**。message 通过统一多语言入口取得公开内部错误文案，保留当前 trace；底层错误和阶段只写服务端日志，不返回原始数据、密文或诊断文本。

即使自定义序列化器返回带业务码的错误，例如 401，也按响应准备失败处理为 code500。真实业务错误仍通过 `HttpErrorResult` 传入，继续沿用原有业务码和文案映射。

失败时不再次序列化原始 data，只编码固定的安全错误结构。正常路径不做“先编码校验、再让框架重编码”，每次响应只对业务数据执行一次序列化；多语言处理可能内部解析 JSON，但不重复调用业务对象的 MarshalJSON。

自定义 MarshalJSON 主动 panic 仍属于上层恢复中间件处理范围，本包不新增 recover。响应编码使用完整内存缓冲，不提供流式输出或额外响应体大小限制。

## AES 与原始字节

AES 仍由 `aesGCM.Configure` / `app.WithAESGCM` 的进程级配置控制，不通过请求上下文传递加密器。

- 开启 AES 的成功响应使用 `data={"aes_data":"..."}`，每次加密生成新 nonce。
- string 和 []byte 在 AES 路径仍作为原始明文字节，不额外包 JSON 引号、不把 []byte 转 Base64 后再加密；其他类型序列化为 JSON 明文。
- 结构化 data 先做已有多语言处理，再加密。普通 JSON 中的大整数和 JSON 数字文本不会经过 float64 中转。
- 错误响应保持普通 JSON；`HttpSuccessResultNotAes` 在进程开启 AES 时仍明确输出普通成功 JSON。

## 写入失败

完成编码后只提交一次状态码、调用一次 Write。客户端断连、写入超时或短写仅记录当前请求日志，包含错误或实际/预期字节数；不重试写入、不追加第二份 JSON，也不改写已经提交的状态码。传输失败时客户端可能只收到部分正文，这不能通过在后面追加错误响应修复。

一个 handler 应选择一个响应出口并立即返回。调用这些函数之前，不要先写状态码或正文；本包不接管任意上游已提交的响应，也不重新设置入口已经提供的 CORS 或请求 ID 响应头。

## 普通 API 范例

先将 [示例词典](../language/examples/etc/language) 复制到项目的 `etc/language`，从项目根目录启动。这个例子接入 Trace 和请求头上下文，展示数据自动翻译、精确数字以及编码失败的统一响应。实际 API 使用平台服务器装配中间件；此 main 用内存 HTTP 请求演示，不启动监听。

```go
package main

import (
    "encoding/json"
    "fmt"
    "math"
    "net/http"
    "net/http/httptest"

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
    handler := appMiddleware.EnsureRequestTrace(
        appMiddleware.NewApiHeaderMiddleware().Handle(func(w http.ResponseWriter, r *http.Request) {
            if r.URL.Path == "/invalid" {
                result.HttpSuccessResult(r.Context(), w, math.NaN())
                return
            }
            result.HttpSuccessResultNotAes(r.Context(), w, map[string]any{
                "n": json.Number("9007199254740993"),
                "notice": "Service Unavailable",
            })
        }),
    )
    for _, path := range []string{"/items", "/invalid"} {
        request := httptest.NewRequest(http.MethodGet, "http://example.test"+path, nil)
        request.Header.Set(headInfo.RequestIDHeader, "1234567890abcdef1234567890abcdef")
        request.Header.Set(consts.ContentLanguage, "zh-CN")
        response := httptest.NewRecorder()
        handler(response, request)
        fmt.Println(response.Code, response.Body.String())
    }
}
```

第一条返回成功 JSON，保留 `9007199254740993`，notice 为 `服务暂不可用。`；第二条返回 HTTP200/code500，message 为 `服务器内部错误。`。两条响应各自沿用请求传入的 ID；范例复用了固定 ID 便于对照，实际独立请求按统一 ID 规则生成。

## 业务错误范例

业务层返回明确的公开错误时，继续调用 `HttpErrorResult`，不需要自行编码。下面展示普通业务冲突的入口形式，错误响应不受成功 data 的 AES 开关影响。

```go
package api

import (
    "net/http"

    "github.com/K7mQ2vR9/base-pkg/result"
    "github.com/K7mQ2vR9/base-pkg/result/xcode"
)

// Conflict 向调用方返回明确的业务冲突，底层内部诊断由业务另行记录。
func Conflict(w http.ResponseWriter, r *http.Request) {
    result.HttpErrorResult(r.Context(), w, xcode.New(http.StatusConflict, "conflict"))
}
```

此例保持既有 HTTP200/code409 协议。普通内部 error 则仍映射为公开内部错误，包装后的业务错误、特殊错误和 nil 错误沿用现有分类规则。
