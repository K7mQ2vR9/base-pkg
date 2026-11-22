package headInfo

import (
	"context"
	"slices"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

// LogKeys 为 HTTP 和 RPC 日志共用的业务字段名单，请求期间只读复用。
// Trace ID 由 go-zero 从当前 Span 自动输出，不再重复记录请求 ID 别名。
var LogKeys = []string{
	consts.RpcName,
	consts.ReqPath, consts.BusinessCode, consts.Business, consts.Source, consts.Version,
	consts.ClientIp, consts.AdminClientIp, consts.TokenUid, consts.TokenUidRole,
	consts.ContentLanguage, consts.Timezone, consts.UserAgent, consts.DeviceId, consts.OriginUrl,
}

type metadataLogContextKey struct{}

// ContextMetadataInLog 按 LogKeys 选择 metadata 日志字段，保留每个字段的字符串列表格式。
// 同名字段覆盖而不重复追加，字段切片取副本，不受调用方后续修改 metadata 影响。
func ContextMetadataInLog(ctx context.Context, md metadata.MD) context.Context {
	ctx = context.WithValue(ctx, metadataLogContextKey{}, true)
	fields := make([]logx.LogField, 0, len(LogKeys))
	for _, key := range LogKeys {
		if key == consts.HeaderToken {
			continue
		}
		if values := md.Get(key); len(values) > 0 {
			fields = append(fields, logx.Field(key, slices.Clone(values)))
		}
	}
	return contextWithLogFields(ctx, fields...)
}

// contextWithLogValue 同步 Setter 对日志字段的修改；沿用入口的 HTTP 字符串或 RPC 字符串列表格式。
// 只同步 LogKeys 中的非凭证字段，不因调用 Setter 扩大日志名单。
func contextWithLogValue(ctx context.Context, key, value string) context.Context {
	if key == consts.HeaderToken || !slices.Contains(LogKeys, key) {
		return ctx
	}
	var logValue any = value
	if metadataMode, _ := ctx.Value(metadataLogContextKey{}).(bool); metadataMode {
		logValue = []string{value}
	}
	return contextWithLogFields(ctx, logx.Field(key, logValue))
}

// contextWithLogFields 使用公开 logx API 构造字段快照，同名字段以后写入的值为准。
// go-zero 的 ContextWithFields 只追加且未公开字段 key，因此通过一个只含日志字段的上下文定位该 key，
// 在复制时合并字段；其他 Value、取消和截止时间均委托原上下文，不使用反射或框架私有 key。
func contextWithLogFields(ctx context.Context, fields ...logx.LogField) context.Context {
	if len(fields) == 0 {
		return ctx
	}
	overrides := logx.ContextWithFields(context.Background(), fields...)
	// 只在构造时合并一次；最终字段由 logx 自身保存，后续记录不重复遍历上下文分支。
	return logx.ContextWithFields(logFieldsContext{Context: ctx, overrides: overrides})
}

type logFieldsContext struct {
	context.Context
	overrides context.Context
}

func (ctx logFieldsContext) Value(key any) any {
	updates, ok := ctx.overrides.Value(key).([]logx.LogField)
	if !ok {
		return ctx.Context.Value(key)
	}
	previous, _ := ctx.Context.Value(key).([]logx.LogField)
	fields := make([]logx.LogField, 0, len(previous)+len(updates))
	positions := make(map[string]int, len(previous)+len(updates))
	for _, group := range [][]logx.LogField{previous, updates} {
		for _, field := range group {
			if position, exists := positions[field.Key]; exists {
				fields[position] = field
				continue
			}
			positions[field.Key] = len(fields)
			fields = append(fields, field)
		}
	}
	return fields
}
