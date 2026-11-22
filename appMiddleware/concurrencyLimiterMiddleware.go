package appMiddleware

import (
	"net/http"

	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/language"
	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type ConcurrencyLimiter struct {
	limiter syncx.Limit
}

func NewConcurrencyLimiter(n int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		limiter: syncx.NewLimit(n),
	}
}

// Handle 超限时按上下文语言返回原有 HTTP503/code503；语言由前置请求头校验器注入。
func (m *ConcurrencyLimiter) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.TryBorrow() {
			const status = http.StatusServiceUnavailable
			ctx := r.Context()
			message := language.TranslateMessage(http.StatusText(status), headInfo.GetContentLanguage(ctx))
			httpx.WriteJsonCtx(ctx, w, status, result.Error(status, message, headInfo.GetTrace(ctx)))
			return
		}
		defer m.limiter.Return()
		next(w, r)
	}
}
