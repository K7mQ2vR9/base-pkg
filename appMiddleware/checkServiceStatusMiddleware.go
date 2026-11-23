package appMiddleware

import (
	"fmt"
	"net/http"

	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/K7mQ2vR9/base-pkg/result/xcode"
)

// CheckServiceStatusFunc 对所有业务请求检查服务状态，包括白名单及匿名请求。
// token 可能为空；返回 false 表示服务不可用，不得以是否登录代替服务状态判断。
type CheckServiceStatusFunc func(r *http.Request, token string) bool
type CheckServiceStatusMiddlewareOption func(m *CheckServiceStatusMiddleware)

// WithCheckServiceStatusCheckOption 配置必需的服务状态回调，nil 在构造中间件时返回配置错误。
func WithCheckServiceStatusCheckOption(check CheckServiceStatusFunc) CheckServiceStatusMiddlewareOption {
	return func(m *CheckServiceStatusMiddleware) {
		m.check = check
	}
}

type CheckServiceStatusMiddleware struct {
	check        CheckServiceStatusFunc
	noVerifyPath map[string]int
}

// NewCheckServiceStatusMiddleware 创建独立服务状态组件，不向 Server 追加中间件。
// 必须提供非 nil 状态回调；缺少回调或包含 nil 选项时返回 nil 和配置错误，不直接 panic。
// noVerifyPath 默认空名单，仅保留配置形式，所有请求均执行状态检查，白名单不豁免。
func NewCheckServiceStatusMiddleware(noVerifyPath map[string]int, ops ...CheckServiceStatusMiddlewareOption) (*CheckServiceStatusMiddleware, error) {
	res := &CheckServiceStatusMiddleware{
		noVerifyPath: noVerifyPath,
	}
	for index, op := range ops {
		if op == nil {
			return nil, fmt.Errorf("service status option at index %d is nil", index)
		}
		op(res)
	}
	if res.check == nil {
		return nil, fmt.Errorf("service status checker is required: use WithCheckServiceStatusCheckOption")
	}
	return res, nil
}

func (m *CheckServiceStatusMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newReq, err := CheckServiceStatusRequest(r, m.check, m.noVerifyPath)
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		next(w, newReq)
	}
}

// CheckServiceStatusRequest 始终执行状态检查，白名单不豁免；缺少回调或返回 false 均拒绝。
// 保留 noVerifyPath 配置参数及调用方式，该参数不再用于跳过服务状态检查。
func CheckServiceStatusRequest(r *http.Request, checkServiceStatus CheckServiceStatusFunc, noVerifyPath map[string]int) (*http.Request, error) {
	if checkServiceStatus == nil || !checkServiceStatus(r, headInfo.GetJwtToken(r.Context())) {
		return r, xcode.StatusServiceUnavailable
	}
	return r, nil
}
