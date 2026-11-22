package appMiddleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/zeromicro/go-zero/core/logc"
)

type ApiHeadOption func(m *ApiHeaderMiddleware)

// CloseVerifyOption 保持现有白名单配置入口；白名单不再跳过请求头校验。
func CloseVerifyOption(path map[string]int) ApiHeadOption {
	return func(m *ApiHeaderMiddleware) {
		m.noVerifyPath = path
	}
}

func WithDebugOption() ApiHeadOption {
	return func(m *ApiHeaderMiddleware) {
		m.debug = true
	}
}

// WithTrustedProxyCIDRs 在装配中间件前校验可信代理配置。
func WithTrustedProxyCIDRs(cidrs []string) (ApiHeadOption, error) {
	resolver, err := headInfo.NewClientIPResolver(cidrs)
	if err != nil {
		return nil, err
	}
	return func(m *ApiHeaderMiddleware) {
		m.clientIPResolver = resolver
	}, nil
}

type ApiHeaderMiddleware struct {
	// NotVerify 保留现有配置字段，但不再豁免请求头校验；白名单仅影响是否必须登录。
	NotVerify        bool
	debug            bool
	noVerifyPath     map[string]int
	clientIPResolver headInfo.ClientIPResolver
}

func NewApiHeaderMiddleware(arg ...ApiHeadOption) *ApiHeaderMiddleware {
	res := &ApiHeaderMiddleware{}
	for _, o := range arg {
		o(res)
	}

	return res
}

// SetNoVerify 保持现有调用方式；该标记不再跳过请求头校验或放宽认证规则。
func (m *ApiHeaderMiddleware) SetNoVerify(b bool) *ApiHeaderMiddleware {
	m.NotVerify = b
	return m
}

// Handle 读取入口已初始化的 Trace 上下文并装配业务头；统一 ID 由路由外层负责。
// 不读取或复制请求体，请求体限制由前置 RequestBodyMiddleware 负责。
// GET、白名单和 NotVerify 请求同样调用 Head.Verify，不因允许匿名而跳过校验。
func (m *ApiHeaderMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if m.debug {
					logc.Errorf(r.Context(), "ApiHeaderMiddleware error:%v, %s", err, string(debug.Stack()))
				}
				logc.Error(r.Context(), err, string(debug.Stack()))
				result.HttpErrorResult(r.Context(), w, errors.New("Server error. "))
				return
			}
		}()

		h := m.clientIPResolver.GetHead(r)
		if err := h.Verify(); err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}

		newCtx := headInfo.ContextHeadInLog(r.Context(), h)
		newCtx = headInfo.HeadInMetadata(newCtx, *h)
		newReq := r.WithContext(newCtx)

		next(w, newReq)
	}
}
