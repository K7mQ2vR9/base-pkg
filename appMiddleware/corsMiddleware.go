package appMiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"github.com/K7mQ2vR9/base-pkg/headInfo"
)

type CorsMiddleware struct {
	policy *corsPolicy
}

// NewCorsMiddleware 创建默认策略：任意来源、开放请求头，不允许浏览器跨域凭证模式。
func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

// Handle 应在 Trace/单值头校验后执行；有效 OPTIONS 返回204，其他请求沿用原业务流程。
func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, bound := r.Context().Value(corsPolicyKey{}).(*corsPolicy); !bound {
			r = r.WithContext(context.WithValue(r.Context(), corsPolicyKey{}, m.policy))
		}
		setCORSHeaders(w.Header(), r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

// Handler 提供独立服务回退入口，预检204，其他方法404，使用已绑定的服务策略。
func (m *CorsMiddleware) Handler() http.Handler {
	return m.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}

func setCORSHeaders(header http.Header, request *http.Request) {
	// 每个出口从同一策略重建授权头，避免错误分支留下之前的 * 或凭证授权。
	for _, key := range []string{"Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Access-Control-Expose-Headers"} {
		header.Del(key)
	}
	policy, _ := request.Context().Value(corsPolicyKey{}).(*corsPolicy)
	if policy != nil && len(policy.origins) > 0 {
		appendCORSVary(header, "Origin")
		origin := requestOrigin(request.Header)
		if _, allowed := policy.origins[origin]; !allowed {
			return
		}
		header.Set("Access-Control-Allow-Origin", origin)
		if policy.credentials {
			header.Set("Access-Control-Allow-Credentials", "true")
		}
	} else {
		header.Set("Access-Control-Allow-Origin", "*")
	}
	header.Set("Access-Control-Allow-Headers", "*")
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	// 暴露统一请求 ID 及使用中划线的项目认证头，不保留旧认证头别名。
	header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers,"+
		http.CanonicalHeaderKey(consts.HeaderToken)+","+headInfo.RequestIDHeader)
	if request.Method == http.MethodOptions {
		const requestHeaders = "Access-Control-Request-Headers"
		appendCORSVary(header, requestHeaders)
		if requested := strings.TrimSpace(strings.Join(request.Header.Values(requestHeaders), ", ")); requested != "" {
			header.Set("Access-Control-Allow-Headers", requested)
		}
	}
}
