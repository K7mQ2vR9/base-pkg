package app

import (
	"fmt"

	"github.com/K7mQ2vR9/base-pkg/aesGCM"
	"github.com/K7mQ2vR9/base-pkg/appMiddleware"
	"github.com/zeromicro/go-zero/rest"
)

type SMOption func(s *ServerMiddleware)

// WithWhiteHeaderPathSMOption 配置允许匿名用户访问的路径，沿用 map 和现有匹配规则。
// 白名单不跳过请求头、Token 校验或服务状态检查，仅允许无认证错误的 uid=0 请求。
func WithWhiteHeaderPathSMOption(whiteHeader map[string]int) SMOption {
	return func(s *ServerMiddleware) {
		s.whiteHeader = whiteHeader
	}
}

func WithDebugOption() SMOption {
	return func(s *ServerMiddleware) {
		s.isDebug = true
	}
}

func WithTestOption() SMOption {
	return func(s *ServerMiddleware) {
		s.isTest = true
	}
}

// WithAESGCM 校验进程级加解密配置，在 ApiUseMiddleware 装配时应用，默认关闭。
func WithAESGCM(config aesGCM.Config) (SMOption, error) {
	if _, err := aesGCM.New(config); err != nil {
		return nil, err
	}
	return func(s *ServerMiddleware) {
		s.aesConfig = config
	}, nil
}

// WithMaxConcurrency sets the shared per-server request limit. Zero disables it.
func WithMaxConcurrency(max int) (SMOption, error) {
	if max < 0 {
		return nil, fmt.Errorf("MaxConcurrency must be non-negative, got %d", max)
	}
	return func(s *ServerMiddleware) {
		s.maxConcurrency = max
	}, nil
}

// WithRequestBody 校验 API 请求体上限，MaxBytes 单位为字节，0 默认 1 MiB，负数拒绝。
// 与 ConfigAppServer.RequestBody 配合使用；网关或 go-zero 的更小限制仍然生效。
func WithRequestBody(config appMiddleware.RequestBodyConfig) (SMOption, error) {
	if _, err := appMiddleware.NewRequestBodyMiddleware(config); err != nil {
		return nil, err
	}
	return func(s *ServerMiddleware) {
		s.requestBodyConfig = config
	}, nil
}

func WithCheckTokenHandleSMOption(fun appMiddleware.CheckRequestTokenFunc) SMOption {
	return func(s *ServerMiddleware) {
		s.checkTokenHandle = fun
		s.authOption = nil
	}
}

// WithCheckTokenHandleWithErrSMOption 启用可返回错误的认证回调。
func WithCheckTokenHandleWithErrSMOption(fun appMiddleware.CheckRequestTokenFuncWithErr) SMOption {
	return withAuthOption(appMiddleware.WithCheckOptionWithErr(fun))
}

// WithCheckTokenHandleWithUserInfoSMOption 启用可返回登录用户 JSON 信息的认证回调。
func WithCheckTokenHandleWithUserInfoSMOption(fun appMiddleware.CheckRequestTokenFuncWithUserInfo) SMOption {
	return withAuthOption(appMiddleware.WithCheckOptionWithUserInfo(fun))
}

func withAuthOption(option appMiddleware.AuthMiddlewareOption) SMOption {
	return func(s *ServerMiddleware) {
		s.authOption = option
		s.checkTokenHandle = nil
	}
}

// WithTrustedProxyCIDRs 设置 HTTP 中间件信任的代理网段。
func WithTrustedProxyCIDRs(cidrs []string) (SMOption, error) {
	option, err := appMiddleware.WithTrustedProxyCIDRs(cidrs)
	if err != nil {
		return nil, err
	}
	return func(s *ServerMiddleware) {
		s.trustedProxyOption = option
	}, nil
}

type ServerMiddleware struct {
	whiteHeader        map[string]int
	checkTokenHandle   appMiddleware.CheckRequestTokenFunc
	authOption         appMiddleware.AuthMiddlewareOption
	trustedProxyOption appMiddleware.ApiHeadOption
	maxConcurrency     int
	// requestBodyConfig 在并发控制之后、解密之前应用，零配置使用默认 1 MiB。
	requestBodyConfig appMiddleware.RequestBodyConfig
	Server            *rest.Server

	isDebug bool
	isTest  bool
	// aesConfig 在 API 装配时应用为进程共用的加解密配置。
	aesConfig aesGCM.Config
}

// NewServerMiddleware 装配应用中间件；Server 用 NewServerWithCORS 创建，以便在原生 Span 前初始化统一 ID。
func NewServerMiddleware(s *rest.Server, opt ...SMOption) *ServerMiddleware {
	res := &ServerMiddleware{
		Server: s,
	}

	for _, item := range opt {
		item(res)
	}

	return res
}

// ApiUseMiddleware 依次装配 CORS、可选并发控制、正文限制、解密、请求头和认证。
// 正文上限未配置时使用 1 MiB；Trace、单值业务头校验及原生长度预检查由 NewServerWithCORS 负责。
// 标准入口的 CORS 策略覆盖正常及早期错误；此处默认 CORS 组件沿用已绑定的实例配置。
func (s *ServerMiddleware) ApiUseMiddleware() {
	authMiddleware := s.authMiddleware()
	s.Server.Use(appMiddleware.NewCorsMiddleware().Handle)
	if s.maxConcurrency > 0 {
		s.Server.Use(appMiddleware.NewConcurrencyLimiter(s.maxConcurrency).Handle)
	}
	body, err := appMiddleware.NewRequestBodyMiddleware(s.requestBodyConfig)
	if err != nil {
		panic(fmt.Errorf("invalid request body configuration: %w", err))
	}
	s.Server.Use(body.Handle)
	s.useApiRequestDecrypt()
	s.useApiHeaderMiddleware()
	s.Server.Use(authMiddleware)
}

func (s *ServerMiddleware) useApiHeaderMiddleware() {
	var apiHeaderOption = []appMiddleware.ApiHeadOption{
		appMiddleware.CloseVerifyOption(s.whiteHeader),
	}
	if s.trustedProxyOption != nil {
		apiHeaderOption = append(apiHeaderOption, s.trustedProxyOption)
	}
	if s.isDebug {
		apiHeaderOption = append(apiHeaderOption, appMiddleware.WithDebugOption())
	}
	s.Server.Use(appMiddleware.NewApiHeaderMiddleware(
		apiHeaderOption...,
	).Handle)
}

func (s *ServerMiddleware) authMiddleware() rest.Middleware {
	if s.authOption != nil {
		return appMiddleware.NewUserAuthMiddlewareWithErr(s.whiteHeader, s.authOption).Handle
	}
	if s.checkTokenHandle == nil {
		panic("token checker is required: use an app token checker option")
	}

	return appMiddleware.NewUserAgentMiddleware(
		s.whiteHeader,
		appMiddleware.WithCheckOption(s.checkTokenHandle),
	).Handle
}

func (s *ServerMiddleware) useApiRequestDecrypt() {
	middleware, err := appMiddleware.NewApiRequestDecryptMiddleware(s.aesConfig, s.requestBodyConfig)
	if err != nil {
		panic(fmt.Errorf("invalid AES-GCM configuration: %w", err))
	}
	s.Server.Use(middleware.Handle)
}
