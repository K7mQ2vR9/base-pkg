package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/appMiddleware"
	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/router"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// NewServerWithCORS 创建统一请求 ID 的 API Server，并校验当前实例的独立 CORS 配置。
// 入口校验先于 go-zero 原生链路、日志和业务中间件，最外层统一正常与错误响应的 CORS 头。
// Trace 就绪后先校验单值业务头，失败返回统一 400，正文、解密及认证均不执行。
// 原生 MaxBytes 启用时提前为声明长度超限输出统一 413；实际读取限制由 ApiUseMiddleware 装配。
// CORS 零值允许任意来源且关闭凭证；开启凭证须配置精确来源，配置非法时返回初始化错误。
// 其余配置沿用 rest.NewServer；勿叠加其他 CORS 组件或用裸 Router 覆盖入口。
func NewServerWithCORS(config rest.RestConf, corsConfig appMiddleware.CORSConfig, options ...rest.RunOption) (*rest.Server, error) {
	maxContentLength := int64(0)
	if config.Middlewares.MaxBytes {
		maxContentLength = config.MaxBytes
	}
	option, err := withRequestTraceCORSRouter(nil, config.Telemetry.Sampler, corsConfig, maxContentLength)
	if err != nil {
		return nil, err
	}
	return rest.NewServer(config, append([]rest.RunOption{option}, options...)...)
}

// withRequestTraceCORSRouter 在包内统一装配实例 CORS、请求 ID 和单值头校验。
// 仅供 NewServerWithCORS 内部调用，base 为 nil 时使用默认路由；CORS 配置非法时返回初始化错误。
// samplingRate 对应 RestConf.Telemetry.Sampler：0 不采样、1 全采样；已有 W3C 上下文沿用原决定。
// maxContentLength 单位为字节，最多一项，省略默认 1 MiB；原生 MaxBytes 未启用时传 0。
// 此处仅预检查声明长度并输出带 Trace 的统一 413，实际读取上限由 ApiUseMiddleware 的 RequestBody 配置负责。
func withRequestTraceCORSRouter(base httpx.Router, samplingRate float64, config appMiddleware.CORSConfig, maxContentLength ...int64) (rest.RunOption, error) {
	if len(maxContentLength) > 1 {
		panic("only one native Content-Length limit is allowed")
	}
	cors, err := appMiddleware.NewCorsMiddlewareWithConfig(config)
	if err != nil {
		return nil, err
	}
	declaredLimit := int64(1 << 20)
	if len(maxContentLength) == 1 {
		declaredLimit = maxContentLength[0]
	}
	if base == nil {
		base = router.NewRouter()
	}
	sampler := sdktrace.TraceIDRatioBased(samplingRate)
	headerGuard := appMiddleware.NewRequestHeaderMiddleware()
	wrapped := &requestTraceRouter{Router: base}
	wrapped.serve = cors.Bind(func(w http.ResponseWriter, r *http.Request) {
		inherited := trace.SpanContextFromContext(r.Context()).IsValid()
		for key := range r.Header {
			if strings.EqualFold(key, "traceparent") {
				inherited = true
			}
		}
		appMiddleware.EnsureRequestTrace(headerGuard.Handle(func(w http.ResponseWriter, request *http.Request) {
			// 原生 MaxBytesHandler 只检查声明长度且返回空 413，在它之前保留统一响应和请求 ID。
			if declaredLimit > 0 && request.ContentLength > declaredLimit {
				appMiddleware.WriteRequestBodyError(w, request, &http.MaxBytesError{Limit: declaredLimit})
				return
			}
			cors.Handle(func(w http.ResponseWriter, request *http.Request) {
				if !inherited {
					// 只有请求 ID 时尚无上游采样决定，按服务的原采样比例建立传播上下文。
					// 避免人为补齐的未采样父上下文使 go-zero 的 ParentBased 采样器丢弃所有新请求。
					ctx := request.Context()
					sc := trace.SpanContextFromContext(ctx)
					decision := sampler.ShouldSample(sdktrace.SamplingParameters{ParentContext: context.Background(), TraceID: sc.TraceID()})
					sc = sc.WithTraceFlags(sc.TraceFlags().WithSampled(decision.Decision == sdktrace.RecordAndSample))
					ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
					ctx, _ = headInfo.EnsureTrace(ctx)
					request = request.WithContext(ctx)
					(propagation.TraceContext{}).Inject(ctx, propagation.HeaderCarrier(request.Header))
				}
				base.ServeHTTP(w, request)
			})(w, request)
		}))(w, r)
	})
	return func(server *rest.Server) {
		rest.WithRouter(wrapped)(server)
		rest.WithNotFoundHandler(nil)(server)
	}, nil
}

type requestTraceRouter struct {
	httpx.Router
	serve http.HandlerFunc
}

func (r *requestTraceRouter) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	r.serve(w, request)
}
