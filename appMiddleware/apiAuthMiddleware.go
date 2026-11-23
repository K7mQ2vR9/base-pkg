package appMiddleware

import (
	"fmt"
	"net/http"

	"github.com/K7mQ2vR9/base-pkg/result"
)

type MiddlewareOption func(m *UserAgentMiddleware)

// WithCheckOption 配置 Token 回调，nil 在 NewUserAgentMiddleware 初始化时返回错误。
func WithCheckOption(check CheckRequestTokenFunc) MiddlewareOption {
	return func(m *UserAgentMiddleware) {
		m.check = check
	}
}

type UserAgentMiddleware struct {
	check        CheckRequestTokenFunc
	noVerifyPath map[string]int
}

// NewUserAgentMiddleware 创建整数结果认证组件，缺少回调或包含 nil 选项时返回 nil 和配置错误。
// 白名单仅豁免必须登录，默认空名单；成功后保持原认证流程，不因配置错误 panic。
func NewUserAgentMiddleware(noVerifyPath map[string]int, ops ...MiddlewareOption) (*UserAgentMiddleware, error) {
	res := &UserAgentMiddleware{
		noVerifyPath: noVerifyPath,
	}
	for index, op := range ops {
		if op == nil {
			return nil, fmt.Errorf("authentication option at index %d is nil", index)
		}
		op(res)
	}
	if res.check == nil {
		return nil, fmt.Errorf("token checker is required: use WithCheckOption")
	}
	return res, nil
}

func (m *UserAgentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newReq, err := MustAuthTokenRequest(r, m.check, m.noVerifyPath)
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		next(w, newReq)
	}
}
