package headInfo

import (
	"context"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"google.golang.org/grpc/metadata"
)

// GetTokenUid 读取出站 metadata 中唯一的十进制用户 ID，范围为 0..MaxInt64。
// 缺失、空值、重复、非法格式或溢出返回 0；0 表示未设置身份，不代表通过认证。
func GetTokenUid(ctx context.Context) int64 {
	value, _ := parseMetadataUID(outgoingMetadataValue(ctx, consts.TokenUid))
	return value
}

// GetTokenUidRole 返回唯一的角色文本，缺失、空值或重复返回空字符串，不解析角色内容。
func GetTokenUidRole(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.TokenUidRole)
}

// GetJwtToken 返回唯一的认证令牌，缺失、空值或重复返回空字符串，不拼接或任选一个值。
func GetJwtToken(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.HeaderToken)
}

// GetClientIp 返回唯一、合法的纯 IP 地址，IPv4-mapped IPv6 统一为 IPv4；非法或重复返回空字符串。
func GetClientIp(ctx context.Context) string {
	return normalizeMetadataIP(outgoingMetadataValue(ctx, consts.ClientIp))
}

// GetAdminClientIp 与 GetClientIp 使用相同规则，拒绝端口、zone、重复及其他非法地址。
func GetAdminClientIp(ctx context.Context) string {
	return normalizeMetadataIP(outgoingMetadataValue(ctx, consts.AdminClientIp))
}

// GetUserAgent 返回唯一的 User-Agent，缺失、空值或重复返回空字符串。
func GetUserAgent(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.UserAgent)
}

// GetVersion 返回唯一的版本文本，缺失、空值或重复返回空字符串。
func GetVersion(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.Version)
}

// GetSource 返回唯一的来源文本，缺失、空值或重复返回空字符串。
func GetSource(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.Source)
}

// SetTokenUid 复制出站 metadata 后覆盖为单个值，不修改原上下文；空值清除身份值。
// 不校验写入格式或身份，非法格式由 Getter 返回 0、RPC 入口拒绝，身份仍须认证回调核实。
func SetTokenUid(ctx context.Context, value string) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	md = md.Copy()
	md.Set(consts.TokenUid, value)
	return metadata.NewOutgoingContext(ctx, md)
}

// GetContentLanguage 返回唯一的语言文本，缺失、空值或重复返回空字符串。
func GetContentLanguage(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.ContentLanguage)
}

// GetBusinessCode 返回唯一的业务标识，缺失、空值或重复返回空字符串。
func GetBusinessCode(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.BusinessCode)
}

// SetBusinessCode 分别复制入站与出站 metadata 并覆盖业务标识，不共享两侧的可变切片。
func SetBusinessCode(ctx context.Context, value string) context.Context {
	incomingMD, _ := metadata.FromIncomingContext(ctx)
	incomingMD = incomingMD.Copy()
	incomingMD.Set(consts.BusinessCode, value)
	ctx = metadata.NewIncomingContext(ctx, incomingMD)
	outgoingMD, _ := metadata.FromOutgoingContext(ctx)
	outgoingMD = outgoingMD.Copy()
	outgoingMD.Set(consts.BusinessCode, value)
	return metadata.NewOutgoingContext(ctx, outgoingMD)
}

// GetBusiness 返回唯一的业务文本，缺失、空值或重复返回空字符串。
func GetBusiness(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.Business)
}

// GetOriginHostUrl 从唯一来源地址中移除既有 http://、https:// 协议文本，不新增 URL 解析规则。
// 缺失、空值或重复来源地址返回空字符串。
func GetOriginHostUrl(ctx context.Context) string {
	value := outgoingMetadataValue(ctx, consts.OriginUrl)
	value = strings.ReplaceAll(value, "http://", "")
	return strings.ReplaceAll(value, "https://", "")
}

// GetTerminal 将唯一的 IOS/Android 来源映射为 APP，其余唯一值原样返回；重复值返回空字符串。
func GetTerminal(ctx context.Context) string {
	value := GetSource(ctx)
	if value == consts.IOS || value == consts.Android {
		return consts.APP
	}
	return value
}

// GetTimezone 返回唯一的时区文本，缺失、空值或重复返回空字符串，不在此解析时区。
func GetTimezone(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.Timezone)
}

// GetDeviceId 返回唯一的设备标识，缺失、空值或重复返回空字符串。
func GetDeviceId(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.DeviceId)
}

// GetPath 返回唯一的请求路径，缺失、空值或重复返回空字符串。
func GetPath(ctx context.Context) string {
	return outgoingMetadataValue(ctx, consts.ReqPath)
}
