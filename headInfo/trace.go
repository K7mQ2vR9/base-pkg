package headInfo

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	ztrace "github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// RequestIDHeader 为 HTTP 对外请求头名，值直接使用统一 Trace ID。
const RequestIDHeader = "X-Request-ID"

var (
	// ErrTraceRequired 表示缺少标准链路头，HTTP 入口可为全新请求生成 ID，RPC 必须拒绝。
	ErrTraceRequired = errors.New("trace ID is required")
	// ErrTraceInvalid 表示 ID 格式错误、字段为空或同一字段重复出现。
	ErrTraceInvalid = errors.New("trace metadata is invalid")
	// ErrTraceConflict 表示请求 ID、W3C 链路或已有活动上下文之间的 Trace ID 不一致。
	ErrTraceConflict = errors.New("request ID and trace ID must match")
)

// NewTraceID 为全新请求生成统一 ID：crypto/rand 读取 16 个随机字节，编码为 32 位小写十六进制。
// 不含时间戳、业务前缀或连字符；全零时重新生成，随机源失败直接返回错误。
// 前端与网关采用相同格式；入口已有合法 ID 时复用，下游不重新生成。
func NewTraceID() (string, error) {
	var id trace.TraceID
	for !id.IsValid() {
		if _, err := io.ReadFull(rand.Reader, id[:]); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(id[:]), nil
}

// GetTrace 返回本次调用唯一的 Trace ID，供业务、日志和响应共用。
func GetTrace(ctx context.Context) string {
	span := trace.SpanContextFromContext(ctx)
	if !span.IsValid() {
		return ""
	}
	return span.TraceID().String()
}

// ExtractTrace 仅从 W3C traceparent/tracestate 恢复 RPC 或消息链路，不读取其他 ID 别名。
// traceparent 缺失返回 ErrTraceRequired，空值、重复和非法格式返回 ErrTraceInvalid。
// 原始 headers 不修改，不生成 Trace ID；已有活动 Span 同链路时保留，不同则返回 ErrTraceConflict。
func ExtractTrace(ctx context.Context, headers map[string][]string) (context.Context, error) {
	parents, present := headerValues(headers, "traceparent")
	if !present {
		return ctx, ErrTraceRequired
	}
	if len(parents) != 1 || parents[0] == "" {
		return ctx, ErrTraceInvalid
	}
	carrier := propagation.MapCarrier{"traceparent": parents[0]}
	if states, _ := headerValues(headers, "tracestate"); len(states) > 0 {
		carrier["tracestate"] = strings.Join(states, ",")
	}
	extracted := (propagation.TraceContext{}).Extract(context.Background(), carrier)
	parent := trace.SpanContextFromContext(extracted)
	if !parent.IsValid() {
		return ctx, ErrTraceInvalid
	}
	if active := trace.SpanContextFromContext(ctx); active.IsValid() {
		if active.TraceID() != parent.TraceID() {
			return ctx, ErrTraceConflict
		}
		return ctx, nil
	}
	return trace.ContextWithRemoteSpanContext(ctx, parent), nil
}

// ExtractRequestTrace 校验 HTTP 的 X-Request-ID 和标准 traceparent，二者同时存在时必须相同。
// 请求 ID 必须是非零的 32 位小写十六进制；不接收任意字符串、空值或重复值。
// 优先恢复请求中的合法 ID，不生成新 Trace ID；与已有活动链路不同则返回 ErrTraceConflict。
// 只有请求 ID 且无活动 Span 时生成 8 字节非零 Span ID 补齐上下文，全部缺失返回 ErrTraceRequired。
func ExtractRequestTrace(ctx context.Context, headers map[string][]string) (context.Context, error) {
	extracted, err := ExtractTrace(ctx, headers)
	if err != nil && !errors.Is(err, ErrTraceRequired) {
		return ctx, err
	}
	ids, present := headerValues(headers, RequestIDHeader)
	if !present {
		return extracted, err
	}
	if len(ids) != 1 {
		return ctx, ErrTraceInvalid
	}
	id, parseErr := trace.TraceIDFromHex(ids[0])
	if parseErr != nil {
		return ctx, ErrTraceInvalid
	}
	if err == nil {
		if GetTrace(extracted) != id.String() {
			return ctx, ErrTraceConflict
		}
		return extracted, nil
	}
	return contextWithTraceID(ctx, id)
}

// EnsureTrace 沿用活动上下文或出站 W3C 链路，缺失时才生成统一 ID，再同步标准传播字段。
// RPC 和 MQ 仅使用 traceparent/tracestate；活动链路与 metadata 冲突时返回错误。
// 仅当活动上下文和出站 traceparent 均缺失时生成新 ID，非法或冲突字段不会触发重新生成。
func EnsureTrace(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromOutgoingContext(ctx)
	extracted, err := ExtractTrace(ctx, md)
	if err == nil {
		ctx = extracted
	} else if !errors.Is(err, ErrTraceRequired) {
		return ctx, err
	} else if GetTrace(ctx) == "" {
		id, err := NewTraceID()
		if err != nil {
			return ctx, err
		}
		parsed, err := trace.TraceIDFromHex(id)
		if err != nil {
			return ctx, err
		}
		ctx, err = contextWithTraceID(ctx, parsed)
		if err != nil {
			return ctx, err
		}
	}
	return traceMetadata(ctx), nil
}

// SetTrace 校验并绑定指定 ID，同步 RPC 出站 W3C 链路；活动 Trace 不同则拒绝覆盖。
func SetTrace(ctx context.Context, value string) (context.Context, error) {
	id, err := trace.TraceIDFromHex(value)
	if err != nil {
		return ctx, ErrTraceInvalid
	}
	ctx, err = contextWithTraceID(ctx, id)
	if err != nil {
		return ctx, err
	}
	return traceMetadata(ctx), nil
}

func traceMetadata(ctx context.Context) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	// tracestate 以当前 Span 为准，避免缺少父级时沿用其他链路的状态。
	md.Delete("tracestate")
	ztrace.Inject(ctx, propagation.TraceContext{}, &md)
	return metadata.NewOutgoingContext(ctx, md)
}

// contextWithTraceID 绑定已校验的 Trace ID，不生成或替换已有 Trace ID。
// 同链路直接保留当前 Span；冲突返回原上下文及错误，仅无活动 Span 时随机补齐非零 8 字节 Span ID。
func contextWithTraceID(ctx context.Context, id trace.TraceID) (context.Context, error) {
	if active := trace.SpanContextFromContext(ctx); active.IsValid() {
		if active.TraceID() != id {
			return ctx, ErrTraceConflict
		}
		return ctx, nil
	}
	var span trace.SpanID
	for !span.IsValid() {
		if _, err := io.ReadFull(rand.Reader, span[:]); err != nil {
			return ctx, err
		}
	}
	return trace.ContextWithRemoteSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{TraceID: id, SpanID: span, Remote: true})), nil
}

func headerValues(headers map[string][]string, name string) ([]string, bool) {
	var values []string
	present := false
	for key, list := range headers {
		if strings.EqualFold(key, name) {
			present = true
			values = append(values, list...)
		}
	}
	return values, present
}
