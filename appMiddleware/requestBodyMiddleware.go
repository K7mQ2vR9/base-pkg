package appMiddleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/language"
	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const defaultRequestBodyMaxBytes int64 = 1 << 20

// RequestBodyConfig 配置 API 实际读取的请求体上限，适用于普通请求与 AES 封装。
type RequestBodyConfig struct {
	// MaxBytes 单位为字节，0 使用默认 1 MiB（1048576），负数拒绝，不提供无限制模式。
	// AES 请求按包含 Base64 密文的 JSON 封装计费；不含请求头，其他框架或网关可施加更小上限。
	MaxBytes int64 `json:",optional"`
}

// RequestBodyMiddleware 在解密和业务解析前按实际字节数限量读取一次，成功后提供可重读的请求体。
// 应放在 Trace 初始化和并发控制之后；不能替代服务端读取超时配置。
type RequestBodyMiddleware struct {
	maxBytes int64
}

// NewRequestBodyMiddleware 校验请求体配置，MaxBytes=0 默认 1 MiB，负数返回错误。
func NewRequestBodyMiddleware(config RequestBodyConfig) (*RequestBodyMiddleware, error) {
	if config.MaxBytes < 0 {
		return nil, errors.New("request body MaxBytes must be non-negative")
	}
	if config.MaxBytes == 0 {
		config.MaxBytes = defaultRequestBodyMaxBytes
	}
	return &RequestBodyMiddleware{maxBytes: config.MaxBytes}, nil
}

// Handle 拒绝超限或读取失败的请求，不调用后续业务；成功后保留原始字节供解析或解密复用。
func (m *RequestBodyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := readRequestBody(w, r, m.maxBytes); err != nil {
			WriteRequestBodyError(w, r, err)
			return
		}
		next(w, r)
	}
}

// WriteRequestBodyError 输出统一请求体错误：MaxBytesError 使用 HTTP 413，其余读取错误使用 HTTP 400。
// 仅使用公开状态文案，不输出原始请求体或底层错误；调用方须先初始化请求 Trace。
// 语言由前置 RequestHeaderMiddleware 注入上下文，不再独立读取或选择原始语言头。
func WriteRequestBodyError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusBadRequest
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		status = http.StatusRequestEntityTooLarge
	}
	message := language.TranslateMessage(http.StatusText(status), headInfo.GetContentLanguage(r.Context()))
	setCORSHeaders(w.Header(), r)
	httpx.WriteJsonCtx(r.Context(), w, status, result.Error(status, message, headInfo.GetTrace(r.Context())))
}

// requestBodyReadError 保留底层读取错误链，让独立解密中间件区分读取失败与业务解密失败。
type requestBodyReadError struct {
	cause error
}

func (err *requestBodyReadError) Error() string { return "request body read failed" }
func (err *requestBodyReadError) Unwrap() error { return err.cause }

// bufferedRequestBody 保存已验证字节，后续解密可复用未消费的数据，避免再次完整读取和复制。
type bufferedRequestBody struct {
	*bytes.Reader
	data []byte
}

func (*bufferedRequestBody) Close() error { return nil }

// readRequestBody 最多向原始流读取上限加一个判定字节，不信任 Content-Length，失败时不提供部分正文。
// 已缓存且尚未消费的正文直接复用；替换后的 Body/GetBody 由进程内字节提供，当前读取流在本函数关闭。
func readRequestBody(w http.ResponseWriter, r *http.Request, maxBytes int64) ([]byte, error) {
	// 零值中间件同样应用默认上限，避免依赖构造函数赋值才具备基本保护。
	if maxBytes == 0 {
		maxBytes = defaultRequestBodyMaxBytes
	}
	if r.ContentLength > maxBytes {
		if r.Body != nil {
			_ = r.Body.Close()
		}
		return nil, &requestBodyReadError{cause: &http.MaxBytesError{Limit: maxBytes}}
	}
	if r.Body == nil || r.Body == http.NoBody {
		return nil, nil
	}
	if body, ok := r.Body.(*bufferedRequestBody); ok && body.Len() == len(body.data) {
		if int64(len(body.data)) > maxBytes {
			return nil, &requestBodyReadError{cause: &http.MaxBytesError{Limit: maxBytes}}
		}
		return body.data, nil
	}
	limited := http.MaxBytesReader(w, r.Body, maxBytes)
	data, err := io.ReadAll(limited)
	_ = limited.Close()
	if err != nil {
		return nil, &requestBodyReadError{cause: err}
	}
	replaceRequestBody(r, data)
	return data, nil
}
