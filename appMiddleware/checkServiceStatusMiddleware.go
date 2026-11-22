package appMiddleware

import (
	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/K7mQ2vR9/base-pkg/result/xcode"
	"net/http"
)

// CheckServiceStatusFunc 对所有业务请求检查服务状态，包括白名单及匿名请求。
// token 可能为空；返回 false 表示服务不可用，不得以是否登录代替服务状态判断。
type CheckServiceStatusFunc func(r *http.Request, token string) bool
type CheckServiceStatusMiddlewareOption func(m *CheckServiceStatusMiddleware)

func WithCheckServiceStatusCheckOption(check CheckServiceStatusFunc) CheckServiceStatusMiddlewareOption {
	return func(m *CheckServiceStatusMiddleware) {
		m.check = check
	}
}

type CheckServiceStatusMiddleware struct {
	check        CheckServiceStatusFunc
	noVerifyPath map[string]int
}

func NewCheckServiceStatusMiddleware(noVerifyPath map[string]int, ops ...CheckServiceStatusMiddlewareOption) *CheckServiceStatusMiddleware {
	res := &CheckServiceStatusMiddleware{
		noVerifyPath: noVerifyPath,
	}
	for _, op := range ops {
		op(res)
	}
	if res.check == nil {
		panic("service status checker is required: use WithCheckServiceStatusCheckOption")
	}
	return res
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
