package headInfo

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/K7mQ2vR9/base-pkg/consts"
)

// ErrInvalidHTTPHeaders 表示已知单值 HTTP 请求头存在重复或无值列表。
// 可用 errors.Is 判断；错误仅包含固定字段名与规则，不回显令牌或其他原始值。
var ErrInvalidHTTPHeaders = errors.New("invalid single-value HTTP headers")

// singleHTTPHeaders 是 API 实际读取的单值业务头，与 RPC 字段及允许多值的标准头分别管理。
// X-Forwarded-For、Trace 头及客户端不能决定的 UID/IP 不在此集合内。
var singleHTTPHeaders = [...]string{
	consts.HeaderToken, consts.Version, consts.Source, consts.Business, consts.BusinessCode,
	consts.ContentLanguage, consts.Timezone, consts.DeviceId, "Origin", "User-Agent",
}

// NormalizeHTTPHeaders 校验已知单值头，成功时返回独立副本并规范化这些头的大小写，供 Header.Get 一致读取。
// 可选头缺失或单个空值允许；空列表及重复值（即使相同）拒绝，大小写别名合并计数。
// 原请求头、值与未知多值头保持不变，不拆分逗号文本，不改变可信代理链或 Trace 校验规则。
func NormalizeHTTPHeaders(headers http.Header) (http.Header, error) {
	normalized := headers.Clone()
	for _, name := range singleHTTPHeaders {
		var originalKey, value string
		present := false
		for key, values := range headers {
			if !strings.EqualFold(key, name) {
				continue
			}
			// 大小写别名也视为重复头，别名中的空列表不能被另一项的值掩盖。
			if present || len(values) != 1 {
				return nil, fmt.Errorf("HTTP header %s must contain exactly one value: %w", name, ErrInvalidHTTPHeaders)
			}
			present, originalKey, value = true, key, values[0]
		}
		if !present {
			continue
		}
		delete(normalized, originalKey)
		normalized.Set(name, value)
	}
	return normalized, nil
}
