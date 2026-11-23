package app

import (
	"fmt"
	"sync"

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

// WithAESGCM 校验进程级加解密配置，在 ApiUseMiddleware 装配时应用，默认关闭。
func WithAESGCM(config aesGCM.Config) (SMOption, error) {
	if _, err := aesGCM.New(config); err != nil {
		return nil, err
	}
	return func(s *ServerMiddleware) {
		s.aesConfig = config
	}, nil
}

// WithMaxConcurrency 设置当前 Server 的共享并发请求数上限，单位为请求数，0 关闭，负数返回配置错误。
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

// WithCheckTokenHandleSMOption 配置整数结果的认证回调；nil 在 ApiUseMiddleware 中返回初始化错误。
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

// ServerMiddleware 保存单个装配实例的配置与状态；开始使用后不得复制或并发修改配置。
// 同一个 Server 应复用此实例完成一次启动装配，不跨实例检查重复装配。
type ServerMiddleware struct {
	// applyMu 串行化同一实例的装配调用，applied 仅在全部中间件追加成功后置为 true。
	applyMu sync.Mutex
	applied bool
	// initErr 保存装配选项错误，统一由 ApiUseMiddleware 返回；选项执行前不会追加中间件。
	initErr            error
	whiteHeader        map[string]int
	checkTokenHandle   appMiddleware.CheckRequestTokenFunc
	authOption         appMiddleware.AuthMiddlewareOption
	trustedProxyOption appMiddleware.ApiHeadOption
	maxConcurrency     int
	// requestBodyConfig 在并发控制之后、解密之前应用，零配置使用默认 1 MiB。
	requestBodyConfig appMiddleware.RequestBodyConfig
	Server            *rest.Server

	isDebug bool
	// aesConfig 在 API 装配时应用为进程共用的加解密配置。
	aesConfig aesGCM.Config
}

// NewServerMiddleware 记录应用中间件选项，不立即向 Server 追加中间件。
// Server 用 NewServerWithCORS 创建；启动时须调用 ApiUseMiddleware 并处理返回的错误。
// nil 选项保存为初始化错误，不执行后续选项；不会因 nil 选项直接 panic。
func NewServerMiddleware(s *rest.Server, opt ...SMOption) *ServerMiddleware {
	res := &ServerMiddleware{
		Server: s,
	}

	for index, item := range opt {
		if item == nil {
			res.initErr = fmt.Errorf("server middleware option at index %d is nil", index)
			return res
		}
		item(res)
	}

	return res
}

// ApiUseMiddleware 依次装配 CORS、可选并发控制、正文限制、解密、请求头和认证。
// 正文上限未配置时使用 1 MiB；Trace、单值业务头校验及原生长度预检查由 NewServerWithCORS 负责。
// 标准入口的 CORS 策略覆盖正常及早期错误；此处默认 CORS 组件沿用已绑定的实例配置。
// 仅在启动阶段调用一次；先完成全部配置校验和组件创建，再应用进程级 AES 配置并统一追加。
// nil Server/选项、缺少认证回调及正文/AES 配置错误均返回 error，失败不追加任何本次中间件。
// 同一实例成功后再次调用返回错误，不追加中间件或再次应用 AES 配置；初始化失败不标记成功。
func (s *ServerMiddleware) ApiUseMiddleware() error {
	if s == nil {
		return fmt.Errorf("API server is required")
	}
	s.applyMu.Lock()
	defer s.applyMu.Unlock()
	if s.applied {
		return fmt.Errorf("API middleware has already been applied")
	}
	if s.Server == nil {
		return fmt.Errorf("API server is required")
	}
	if s.initErr != nil {
		return s.initErr
	}
	if s.maxConcurrency < 0 {
		return fmt.Errorf("MaxConcurrency must be non-negative, got %d", s.maxConcurrency)
	}
	authMiddleware, err := s.authMiddleware()
	if err != nil {
		return fmt.Errorf("invalid authentication configuration: %w", err)
	}
	body, err := appMiddleware.NewRequestBodyMiddleware(s.requestBodyConfig)
	if err != nil {
		return fmt.Errorf("invalid request body configuration: %w", err)
	}
	headerMiddleware := s.apiHeaderMiddleware()
	middlewares := []rest.Middleware{appMiddleware.NewCorsMiddleware().Handle}
	if s.maxConcurrency > 0 {
		middlewares = append(middlewares, appMiddleware.NewConcurrencyLimiter(s.maxConcurrency).Handle)
	}
	// 解密构造会发布进程级 AES 配置，因此放在其他可失败的配置与组件创建之后。
	decrypt, err := appMiddleware.NewApiRequestDecryptMiddleware(s.aesConfig, s.requestBodyConfig)
	if err != nil {
		return fmt.Errorf("invalid AES-GCM configuration: %w", err)
	}
	middlewares = append(middlewares, body.Handle, decrypt.Handle, headerMiddleware, authMiddleware)
	for _, middleware := range middlewares {
		s.Server.Use(middleware)
	}
	s.applied = true
	return nil
}

// apiHeaderMiddleware 仅创建请求头组件，实际追加由 ApiUseMiddleware 在全部初始化成功后执行。
func (s *ServerMiddleware) apiHeaderMiddleware() rest.Middleware {
	var apiHeaderOption = []appMiddleware.ApiHeadOption{
		appMiddleware.CloseVerifyOption(s.whiteHeader),
	}
	if s.trustedProxyOption != nil {
		apiHeaderOption = append(apiHeaderOption, s.trustedProxyOption)
	}
	if s.isDebug {
		apiHeaderOption = append(apiHeaderOption, appMiddleware.WithDebugOption())
	}
	return appMiddleware.NewApiHeaderMiddleware(
		apiHeaderOption...,
	).Handle
}

// authMiddleware 创建已校验认证回调的组件，缺少回调或选项无效时返回错误，不追加中间件。
func (s *ServerMiddleware) authMiddleware() (rest.Middleware, error) {
	if s.authOption != nil {
		middleware, err := appMiddleware.NewUserAuthMiddlewareWithErr(s.whiteHeader, s.authOption)
		if err != nil {
			return nil, err
		}
		return middleware.Handle, nil
	}
	middleware, err := appMiddleware.NewUserAgentMiddleware(
		s.whiteHeader,
		appMiddleware.WithCheckOption(s.checkTokenHandle),
	)
	if err != nil {
		return nil, err
	}
	return middleware.Handle, nil
}
