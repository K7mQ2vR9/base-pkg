package headInfo

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
)

const PLATFORM_BUSINESS_CODE = "99999999"

type Head struct {
	AuthorizationJwt string `json:"authorization_jwt"` // 用户token
	Version          string `json:"version"`           // APP版本
	Source           string `json:"source"`            // 来源渠道
	ClientIp         string `json:"client_ip"`         // 客户端IP
	AdminClientIp    string `json:"admin_client_ip"`   // 后台请求的客户端 IP
	Trace            string `json:"trace"`             // 统一请求 ID，32 位小写十六进制，整个链路保持一致。
	TokenUid         string `json:"token_uid"`         // 用户ID
	ReqPath          string `json:"req_path"`          // 请求path
	Business         string `json:"business"`
	BusinessCode     string `json:"business_code"`
	ContentLanguage  string `json:"content_language"`
	TokenUidRole     string `json:"token_uid_role"`
	ReqOrigin        string `json:"req_origin"` // 请求地址
	Timezone         string `json:"timezone"`
	UserAgent        string `json:"user_agent"`
	DeviceId         string `json:"device_id"`
}

// GetHead 读取请求头，直接使用连接对端作为客户端 IP。
func GetHead(r *http.Request) *Head {
	return ClientIPResolver{}.GetHead(r)
}

// GetHead 读取请求头，并按已配置的可信代理解析客户端 IP。
func (resolver ClientIPResolver) GetHead(r *http.Request) *Head {
	header := r.Header
	clientIP := resolver.clientIP(r)
	return &Head{
		AuthorizationJwt: strings.Trim(header.Get(consts.HeaderToken), " "),
		Version:          strings.Trim(header.Get(consts.Version), " "),
		Source:           strings.Trim(header.Get(consts.Source), " "),
		ClientIp:         clientIP,
		AdminClientIp:    clientIP,
		Trace:            GetTrace(r.Context()),
		ReqPath:          r.URL.Path,
		Business:         strings.Trim(header.Get(consts.Business), " "),
		BusinessCode:     strings.Trim(header.Get(consts.BusinessCode), " "),
		ContentLanguage:  strings.Trim(header.Get(consts.ContentLanguage), " "),
		TokenUid:         "",
		TokenUidRole:     "",
		ReqOrigin:        strings.Trim(header.Get("Origin"), " "),
		Timezone:         strings.Trim(header.Get(consts.Timezone), " "),
		UserAgent:        strings.Trim(header.Get("user-agent"), ""),
		DeviceId:         strings.Trim(header.Get(consts.DeviceId), ""),
	}
}

func (h *Head) Verify() error {
	return nil
}

func (h *Head) String() string {
	if h == nil {
		return "null"
	}

	// 字符串表示用于日志，清空副本中的凭证以保留原始认证数据。
	logHead := *h
	logHead.AuthorizationJwt = ""
	data, _ := json.Marshal(logHead)
	return string(data)
}

// ContextHeadInLog 按共用 LogKeys 选择 HTTP 日志字段，字段值使用当前请求的字符串。
// 同名字段覆盖，不重复追加；保留其他业务日志字段，后续 Setter 沿用字符串格式。
func ContextHeadInLog(ctx context.Context, h *Head) context.Context {
	fields := []logx.LogField{
		logx.Field(consts.Version, h.Version),
		logx.Field(consts.Source, h.Source),
		logx.Field(consts.ClientIp, h.ClientIp),
		logx.Field(consts.AdminClientIp, h.AdminClientIp),
		logx.Field(consts.TokenUid, h.TokenUid),
		logx.Field(consts.ReqPath, h.ReqPath),
		logx.Field(consts.Business, h.Business),
		logx.Field(consts.BusinessCode, h.BusinessCode),
		logx.Field(consts.ContentLanguage, h.ContentLanguage),
		logx.Field(consts.TokenUidRole, h.TokenUidRole),
		logx.Field(consts.OriginUrl, h.ReqOrigin),
		logx.Field(consts.Timezone, h.Timezone),
		logx.Field(consts.UserAgent, h.UserAgent),
		logx.Field(consts.DeviceId, h.DeviceId),
	}
	selected := fields[:0]
	for _, field := range fields {
		if slices.Contains(LogKeys, field.Key) {
			selected = append(selected, field)
		}
	}
	ctx = context.WithValue(ctx, metadataLogContextKey{}, false)
	return contextWithLogFields(ctx, selected...)
}

// GetFullHead 返回原始请求头的独立副本，map 和各值切片均不与请求共享。
// 保留字段名、值及多值顺序；空请求头返回可写空 map，r 必须非 nil。
// 结果仍可能包含认证凭证，本函数不做脱敏，由调用方控制使用范围。
func GetFullHead(r *http.Request) map[string][]string {
	headers := r.Header.Clone()
	if headers == nil {
		return make(map[string][]string)
	}
	return headers
}

// HeadInMetadata 写入 RPC 业务头和 W3C 链路，仅使用上下文中的统一 Trace ID。
// 复制已有出站 metadata 后覆盖 Head 对应的平台字段，包括空值，保留其他自定义字段及多值。
// 不合并入站 metadata、不修改父上下文；清除原出站 W3C 字段，再按活动上下文重新注入。
func HeadInMetadata(ctx context.Context, h Head) context.Context {
	platformFields := metadata.Pairs(
		consts.HeaderToken, h.AuthorizationJwt,
		consts.TokenUid, h.TokenUid,
		consts.ClientIp, h.ClientIp,
		consts.AdminClientIp, h.AdminClientIp,
		consts.ReqPath, h.ReqPath,
		consts.Version, h.Version,
		consts.Source, h.Source,
		consts.Business, h.Business,
		consts.BusinessCode, h.BusinessCode,
		consts.ContentLanguage, h.ContentLanguage,
		consts.TokenUidRole, h.TokenUidRole,
		consts.OriginUrl, h.ReqOrigin,
		consts.Timezone, h.Timezone,
		consts.UserAgent, h.UserAgent,
		consts.DeviceId, h.DeviceId,
	)

	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	for key, values := range platformFields {
		// 使用 Set 覆盖为单值，不能追加，否则旧身份/重复平台字段可能留在新请求中。
		md.Set(key, values...)
	}
	// 原实现重建 metadata 时不会沿用旧链路；合并后仍只信任当前活动 Trace。
	md.Delete("traceparent")
	md.Delete("tracestate")
	ctxNew := metadata.NewOutgoingContext(ctx, md)
	if GetTrace(ctx) != "" {
		return traceMetadata(ctxNew)
	}
	return ctxNew
}
