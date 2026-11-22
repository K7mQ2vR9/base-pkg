package headInfo

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strconv"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"google.golang.org/grpc/metadata"
)

// ErrInvalidMetadata 表示平台已知的单值业务字段存在重复、空列表或非法 UID/IP 格式。
// 通过 errors.Is 判断；错误文本只描述字段和规则，不包含令牌或其他原始值。
var ErrInvalidMetadata = errors.New("invalid business metadata")

// singleMetadataKeys 是平台定义的单值业务字段集合，与允许多值的自定义 metadata、tracestate 区分。
// 仅包内只读复用，不依赖可由业务修改的日志字段列表。
var singleMetadataKeys = [...]string{
	consts.HeaderToken, consts.TokenUid, consts.TokenUidRole,
	consts.ClientIp, consts.AdminClientIp,
	consts.ReqPath, consts.Version, consts.Source, consts.Business, consts.BusinessCode,
	consts.ContentLanguage, consts.OriginUrl, consts.Timezone, consts.UserAgent, consts.DeviceId,
	consts.RpcName,
}

// ValidateMetadata 校验已知业务字段的单值约束，以及非空 UID/IP 的格式；不修改原始 metadata。
// 可选字段缺失或单个空字符串允许，存在的字段必须恰好一个值；相同内容重复也拒绝。
// 未知字段允许多值；Trace 的必填、格式及传播仍由原有链路函数独立校验。
func ValidateMetadata(md metadata.MD) error {
	for _, key := range singleMetadataKeys {
		values, present := headerValues(md, key)
		if !present {
			continue
		}
		if len(values) != 1 {
			return fmt.Errorf("metadata %s must contain exactly one value: %w", key, ErrInvalidMetadata)
		}
		if values[0] == "" {
			continue
		}
		switch key {
		case consts.TokenUid:
			if _, err := parseMetadataUID(values[0]); err != nil {
				return fmt.Errorf("metadata %s must be decimal digits in range 0..MaxInt64: %w", key, ErrInvalidMetadata)
			}
		case consts.ClientIp, consts.AdminClientIp:
			if normalizeMetadataIP(values[0]) == "" {
				return fmt.Errorf("metadata %s must be an IP address without port or zone: %w", key, ErrInvalidMetadata)
			}
		}
	}
	return nil
}

// outgoingMetadataValue 只读取出站业务 metadata 的唯一值，缺失、空列表或重复返回空字符串。
// Getter 保持现有出站读取语义，不回退到入站 metadata；nil ctx 按未设置处理。
func outgoingMetadataValue(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	md, _ := metadata.FromOutgoingContext(ctx)
	// gRPC 的上下文读取 API 已将字段名规范化为小写，可直接读取唯一键，避免再次遍历与复制。
	values := md.Get(key)
	if len(values) != 1 {
		return ""
	}
	return values[0]
}

// parseMetadataUID 接受空身份或 0..MaxInt64 的 ASCII 十进制数字，不接受符号、空白及其他进制。
// 前导零允许；返回错误不携带原始值，避免被入口错误响应带出。
func parseMetadataUID(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return 0, ErrInvalidMetadata
		}
	}
	uid, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, ErrInvalidMetadata
	}
	return uid, nil
}

// normalizeMetadataIP 统一两种客户端地址读取规则，拒绝空白、端口、列表与 zone，规范化 IPv4-mapped IPv6。
func normalizeMetadataIP(value string) string {
	ip, err := netip.ParseAddr(value)
	if err != nil || ip.Zone() != "" {
		return ""
	}
	return ip.Unmap().String()
}
