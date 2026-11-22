package appMiddleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CORSConfig 配置单个 API Server 的跨域策略，不修改进程全局配置。
type CORSConfig struct {
	// AllowCredentials 默认 false；开启后允许浏览器凭证请求，必须配置非空 AllowedOrigins。
	AllowCredentials bool `json:",optional"`
	// AllowedOrigins 为精确来源列表，例如 https://app.example.com（不带尾部斜杠）。
	// 禁止 *、null、路径、查询串和用户信息；未配置且凭证关闭时允许任意来源。
	// 匹配区分大小写及显式端口，调用方应填写浏览器发送的序列化 Origin。
	AllowedOrigins []string `json:",optional"`
}

type corsPolicy struct {
	credentials bool
	origins     map[string]struct{}
}

type corsPolicyKey struct{}

// NewCorsMiddlewareWithConfig 校验并复制来源配置，供单个服务并发复用。
// 凭证模式不允许通配来源；错误只包含字段和索引，不回显配置原值。
func NewCorsMiddlewareWithConfig(config CORSConfig) (*CorsMiddleware, error) {
	policy := &corsPolicy{credentials: config.AllowCredentials, origins: make(map[string]struct{}, len(config.AllowedOrigins))}
	if config.AllowCredentials && len(config.AllowedOrigins) == 0 {
		return nil, fmt.Errorf("CORS AllowedOrigins is required when AllowCredentials is enabled")
	}
	for index, origin := range config.AllowedOrigins {
		parsed, err := url.Parse(origin)
		if err != nil || strings.TrimSpace(origin) != origin || strings.ContainsAny(origin, "\r\n\t ,\\#?") ||
			parsed == nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Hostname() == "" ||
			parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.Opaque != "" ||
			origin != parsed.Scheme+"://"+parsed.Host || strings.ContainsAny(parsed.Host, "*%") {
			return nil, fmt.Errorf("CORS AllowedOrigins[%d] must be an exact HTTP/HTTPS origin", index)
		}
		if port := parsed.Port(); port != "" {
			number, err := strconv.Atoi(port)
			if err != nil || number < 0 || number > 65535 {
				return nil, fmt.Errorf("CORS AllowedOrigins[%d] has an invalid port", index)
			}
		} else if strings.HasSuffix(parsed.Host, ":") {
			return nil, fmt.Errorf("CORS AllowedOrigins[%d] has an empty port", index)
		}
		policy.origins[origin] = struct{}{}
	}
	return &CorsMiddleware{policy: policy}, nil
}

// Bind 在最外层绑定当前服务策略并写入 CORS 头，覆盖入口 Trace/请求头错误和路由404/405。
// 不拦截 OPTIONS、不读取正文、不校验身份；正常预检由位于 Trace/头校验之后的 Handle 结束。
// 嵌套的默认 CORS 组件和错误出口优先使用此策略，避免覆盖已配置的来源约束。
func (m *CorsMiddleware) Bind(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := r.WithContext(context.WithValue(r.Context(), corsPolicyKey{}, m.policy))
		setCORSHeaders(w.Header(), request)
		next(w, request)
	}
}

// requestOrigin 只采用唯一非空 Origin，不信任重复值或大小写别名中的任意一项。
func requestOrigin(header http.Header) string {
	var origin string
	present := false
	for key, values := range header {
		if !strings.EqualFold(key, "Origin") {
			continue
		}
		if present || len(values) != 1 {
			return ""
		}
		present = true
		origin = values[0]
	}
	return origin
}

// appendCORSVary 保留已有缓存区分字段，大小写不敏感去重，不覆盖 Vary: *。
func appendCORSVary(header http.Header, name string) {
	for _, line := range header.Values("Vary") {
		for _, item := range strings.Split(line, ",") {
			if item = strings.TrimSpace(item); item == "*" || strings.EqualFold(item, name) {
				return
			}
		}
	}
	header.Add("Vary", name)
}
