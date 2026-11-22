package appMiddleware

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/result/xcode"
)

type (
	// CheckRequestTokenFuncWithErr 校验已提供的 Token，返回用户 ID 或认证错误。
	// uid=0 且 err=nil 表示已确认的匿名结果，仅白名单可放行；负数和任何错误均拒绝。
	CheckRequestTokenFuncWithErr func(r *http.Request, token string) (int64, error)
	// CheckRequestTokenFunc 校验已提供的 Token；0 表示匿名，-1 表示无权限，其他负数表示无效。
	// 需要区分过期、无效 Token 与匿名结果时，应使用可返回 error 的回调。
	CheckRequestTokenFunc func(r *http.Request, token string) int64
	// CheckRequestTokenFuncWithUserInfo 返回已验证的用户 ID、用户信息 JSON 和认证错误。
	// uid=0 且 err=nil 表示匿名，扩展信息清空；无效/过期 Token 必须返回错误，不得用 0,nil 代替。
	CheckRequestTokenFuncWithUserInfo func(r *http.Request, token string) (uid int64, userInfo json.RawMessage, err error)
)

func isWhitelistedPath(urlPath string, noVerifyPath map[string]int) bool {
	if _, ok := noVerifyPath[urlPath]; ok {
		return true
	}

	for path := range noVerifyPath {
		if strings.HasPrefix(path, "/") && strings.HasSuffix(path, "*") {
			prefix := strings.TrimSuffix(path, "*")
			if strings.HasPrefix(urlPath, prefix) {
				return true
			}
		}
	}
	return false
}

// MustAuthTokenRequest 校验 Token，白名单仅允许正常匿名结果，不吞掉认证或权限错误。
func MustAuthTokenRequest(r *http.Request, checkToken CheckRequestTokenFunc, noVerifyPath map[string]int) (*http.Request, error) {
	var check CheckRequestTokenFuncWithErr
	if checkToken != nil {
		check = func(r *http.Request, token string) (int64, error) {
			return checkToken(r, token), nil
		}
	}
	return MustAuthTokenRequestWithErr(r, check, noVerifyPath)
}

// MustAuthTokenRequestWithErr 对所有已提供的 Token 执行认证，错误在白名单上同样返回。
// 缺少 Token 或回调返回 0,nil 时，仅白名单允许以 uid=0、无用户扩展信息继续。
func MustAuthTokenRequestWithErr(r *http.Request, checkToken CheckRequestTokenFuncWithErr, noVerifyPath map[string]int) (*http.Request, error) {
	return mustAuthTokenRequestWithUserInfo(r, adaptTokenCheck(checkToken), noVerifyPath)
}

func adaptTokenCheck(checkToken CheckRequestTokenFuncWithErr) CheckRequestTokenFuncWithUserInfo {
	if checkToken == nil {
		return nil
	}
	return func(r *http.Request, token string) (int64, json.RawMessage, error) {
		uid, err := checkToken(r, token)
		return uid, nil, err
	}
}

func mustAuthTokenRequestWithUserInfo(r *http.Request, checkToken CheckRequestTokenFuncWithUserInfo, noVerifyPath map[string]int) (*http.Request, error) {
	uid, userInfo, err := authenticateToken(r, checkToken)
	if err != nil {
		return r, err
	}
	if uid == 0 {
		if !isWhitelistedPath(r.URL.Path, noVerifyPath) {
			return r, xcode.UserNotFound
		}
		// 匿名请求显式写入 0 并清空扩展信息，不继承父上下文的登录身份。
		userInfo = nil
	}
	ctx := headInfo.SetTokenUid(r.Context(), strconv.FormatInt(uid, 10))
	return r.WithContext(headInfo.SetUserInfo(ctx, userInfo)), nil
}

func authenticateToken(r *http.Request, checkToken CheckRequestTokenFuncWithUserInfo) (int64, json.RawMessage, error) {
	token := strings.TrimSpace(r.Header.Get(consts.HeaderToken))
	if token == "" {
		return 0, nil, nil
	}
	if checkToken == nil {
		return 0, nil, xcode.StatusServiceUnavailable
	}

	uid, userInfo, err := checkToken(r, token)
	switch {
	case err != nil:
		return 0, nil, xcode.TokenSpecial(err)
	case uid == -1:
		return 0, nil, xcode.TokenPermissionDenied
	case uid < 0:
		return 0, nil, xcode.TokenInvalid
	case len(userInfo) != 0 && !json.Valid(userInfo):
		return 0, nil, xcode.StatusServiceUnavailable
	}
	return uid, userInfo, nil
}
