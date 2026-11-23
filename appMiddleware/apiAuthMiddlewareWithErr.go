package appMiddleware

import (
	"fmt"
	"net/http"

	"github.com/K7mQ2vR9/base-pkg/result"
)

type AuthMiddlewareOption func(m *UserAuthMiddlewareWithErr)

// WithCheckOptionWithErr 配置可返回错误的 Token 回调，nil 在构造认证组件时返回配置错误。
func WithCheckOptionWithErr(check CheckRequestTokenFuncWithErr) AuthMiddlewareOption {
	return WithCheckOptionWithUserInfo(adaptTokenCheck(check))
}

// WithCheckOptionWithUserInfo 配置可提供登录用户信息的认证回调，nil 在构造组件时返回配置错误。
func WithCheckOptionWithUserInfo(check CheckRequestTokenFuncWithUserInfo) AuthMiddlewareOption {
	return func(m *UserAuthMiddlewareWithErr) {
		m.check = check
	}
}

type UserAuthMiddlewareWithErr struct {
	check        CheckRequestTokenFuncWithUserInfo
	noVerifyPath map[string]int
}

// NewUserAuthMiddlewareWithErr 创建错误/用户信息认证组件，缺少回调或包含 nil 选项时返回 nil 和配置错误。
// 白名单默认空名单，运行时认证、匿名和用户扩展信息规则保持，不因配置错误 panic。
func NewUserAuthMiddlewareWithErr(noVerifyPath map[string]int, ops ...AuthMiddlewareOption) (*UserAuthMiddlewareWithErr, error) {
	res := &UserAuthMiddlewareWithErr{
		noVerifyPath: noVerifyPath,
	}
	for index, op := range ops {
		if op == nil {
			return nil, fmt.Errorf("authentication option at index %d is nil", index)
		}
		op(res)
	}
	if res.check == nil {
		return nil, fmt.Errorf("token checker is required: use WithCheckOptionWithErr or WithCheckOptionWithUserInfo")
	}
	return res, nil
}

func (m *UserAuthMiddlewareWithErr) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newReq, err := mustAuthTokenRequestWithUserInfo(r, m.check, m.noVerifyPath)
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		next(w, newReq)
	}
}
