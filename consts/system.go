package consts

// HTTP 自定义业务头与 RPC metadata 统一使用小写中划线，不读取下划线旧名或其他别名。
// 此处只定义传输字段名，JSON 标签与 Go 字段名按各自的数据协议保持不变。
const (
	// HeaderToken 为项目认证令牌的 HTTP 头及 RPC metadata 名称。
	HeaderToken = "authorization-jwt"
	// ClientIp 为可信代理解析后传入 RPC 的客户端地址，不直接信任同名 HTTP 入站值。
	ClientIp = "client-ip"
	// AdminClientIp 为后台客户端地址的 RPC metadata 名称。
	AdminClientIp = "admin-client-ip"
	// UserAgent 保存原终端 User-Agent，使用独立名称避开 gRPC 保留的 user-agent。
	UserAgent = "client-user-agent"
	// TokenUid 为认证回调确认的用户 ID，不直接采用客户端同名头作为身份。
	TokenUid = "token-uid"
	// TokenUidRole 为业务认证确定的角色文本。
	TokenUidRole = "token-uid-role"
	// Version 为客户端版本，单词名称无需分隔符。
	Version = "version"
	// Source 为来源终端，单词名称无需分隔符。
	Source = "source"
	// ReqPath 为 API 请求路径的 RPC metadata 名称。
	ReqPath = "req-path"
	// Business 为业务名称，单词名称无需分隔符。
	Business = "business"
	// BusinessCode 为业务标识，在 HTTP 与 RPC 两层使用同一名称。
	BusinessCode = "business-code"
	// ContentLanguage 为当前请求语言，在 HTTP 与 RPC 两层使用同一名称。
	ContentLanguage = "content-language"
	// OriginUrl 为 HTTP Origin 转入 RPC 后的字段名称，HTTP 端继续使用标准 Origin。
	OriginUrl = "origin-url"
	// Timezone 为时区文本，单词名称无需分隔符。
	Timezone = "timezone"
	// DeviceId 为客户端设备标识，在 HTTP 与 RPC 两层使用同一名称。
	DeviceId = "device-id"
	// RpcName 为客户端拦截器写入的当前目标服务诊断名称。
	RpcName = "rpc-name"
)

const (
	APP     = "APP"
	IOS     = "IOS"
	Android = "Android"
	PC      = "PC"
	H5      = "H5"
)
