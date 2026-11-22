package appMiddleware

import (
	"net/http"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/language"
	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// RequestHeaderMiddleware 在正文读取、解密和认证前校验单值业务头，不读取或关闭请求体。
// 应位于 Trace 初始化之后、框架原生中间件之前，避免 gzip 等前置处理提前读取正文。
type RequestHeaderMiddleware struct{}

// NewRequestHeaderMiddleware 创建固定规则的单值请求头校验器，不设置必填头或白名单豁免。
// X-Forwarded-For 等允许多值的头不受此规则限制，仍由各自协议处理。
func NewRequestHeaderMiddleware() *RequestHeaderMiddleware {
	return &RequestHeaderMiddleware{}
}

// Handle 拒绝非法单值头并输出带当前 Trace 的 HTTP400/code400；成功时下游使用规范化头副本。
// GET、白名单及其他方法使用相同规则，不改变原请求、正文或上游上下文。
func (m *RequestHeaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headers, err := headInfo.NormalizeHTTPHeaders(r.Header)
		if err != nil {
			locale := headInfo.GetContentLanguage(r.Context())
			if locale == "" {
				// 错误响应也只使用唯一语言值，语言头自身重复时不任选一个值。
				var values []string
				for key, entries := range r.Header {
					if strings.EqualFold(key, consts.ContentLanguage) {
						values = append(values, entries...)
					}
				}
				if len(values) == 1 {
					locale = values[0]
				}
			}
			message := language.TranslateMessage(http.StatusText(http.StatusBadRequest), locale)
			setCORSHeaders(w.Header(), r)
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest,
				result.Error(http.StatusBadRequest, message, headInfo.GetTrace(r.Context())))
			return
		}
		request := r.WithContext(r.Context())
		request.Header = headers
		next(w, request)
	}
}
