package appMiddleware

import (
	"net/http"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"github.com/K7mQ2vR9/base-pkg/headInfo"
)

type CorsMiddleware struct {
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w.Header(), r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func (m *CorsMiddleware) Handler() http.Handler {
	return m.Handle(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}

func setCORSHeaders(header http.Header, request *http.Request) {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Headers", "*")
	header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	// 暴露统一请求 ID 及使用中划线的项目认证头，不保留旧认证头别名。
	header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Access-Control-Allow-Origin, Access-Control-Allow-Headers,"+
		http.CanonicalHeaderKey(consts.HeaderToken)+","+headInfo.RequestIDHeader)
	header.Set("Access-Control-Allow-Credentials", "true")
	if request.Method == http.MethodOptions {
		const requestHeaders = "Access-Control-Request-Headers"
		header.Add("Vary", requestHeaders)
		if requested := strings.TrimSpace(strings.Join(request.Header.Values(requestHeaders), ", ")); requested != "" {
			header.Set("Access-Control-Allow-Headers", requested)
		}
	}
}
