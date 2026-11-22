package dbM

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig 配置 Redis 的部署模式、地址、认证和 TLS。
type RedisConfig struct {
	// Type 默认 node，可选 cluster；Go 直接传参时空值也使用 node。
	Type string `json:",default=node,options=node|cluster"`
	// Addr 为 node 模式的必填地址，格式为 host:port。
	Addr string `json:",optional"`
	// Addrs 为 cluster 模式的种子地址列表，至少填写一个 host:port。
	Addrs []string `json:",optional"`
	// Username 为 ACL 用户名，默认空值，由服务端认证方式决定是否填写。
	Username string `json:",optional"`
	// Password 为认证密码，默认空值。
	Password string `json:",optional"`
	// DB 默认 0，node 模式允许非负值，cluster 模式使用 0。
	DB int `json:",optional"`
	// TLS 默认 false，开启后使用 TLS 并默认校验证书。
	TLS bool `json:",optional"`
	// InsecureSkipVerify 默认 false，仅在启用 TLS 时允许显式跳过证书校验。
	InsecureSkipVerify bool `json:",optional"`
	// Timeout 管理启动、拨号、读写和等待连接超时，未配置时使用各字段默认值。
	Timeout RedisTimeoutConfig `json:",optional"`
	// Pool 管理连接池容量与回收；集群模式下每个节点分别应用。
	Pool RedisPoolConfig `json:",optional"`
	// Retry 管理命令重试及集群重定向，默认均关闭，启用前须确认业务幂等性。
	Retry RedisRetryConfig `json:",optional"`
}

// NewRedisClient 按部署模式创建客户端，启动预算由 Timeout.Startup 指定，默认 5s，父级更早截止时间优先。
// 启动失败或取消时关闭客户端，成功后由调用方在服务退出时调用 Close。
// 命令自动重试与集群重定向默认关闭；请求须传当前 ctx，启动 ctx 不控制成功后的实例生命周期。
func NewRedisClient(ctx context.Context, config RedisConfig) (redis.UniversalClient, error) {
	if ctx == nil {
		return nil, errors.New("redis startup context is required")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	settings, err := config.settings()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, settings.startupTimeout)
	defer cancel()
	mode, options := settings.mode, settings.options

	var client redis.UniversalClient
	if mode == "cluster" {
		clusterOptions := options.Cluster()
		// 集群初始化会把退避的 -1 转为 0，节点再初始化时又会将 0 恢复成默认退避。
		// 节点构造前重新传入已校验的原始值，保证显式关闭退避在两层均生效。
		clusterOptions.NewClient = func(nodeOptions *redis.Options) *redis.Client {
			nodeOptions.MinRetryBackoff = options.MinRetryBackoff
			nodeOptions.MaxRetryBackoff = options.MaxRetryBackoff
			return redis.NewClient(nodeOptions)
		}
		client = redis.NewClusterClient(clusterOptions)
	} else {
		client = redis.NewClient(options.Simple())
	}

	// go-redis 的 I/O deadline 不会因 ctx 提前取消而立即更新。
	// 启动取消时关闭连接中断读取，PING 完成后撤销回调，交给调用方管理客户端。
	stopClose := context.AfterFunc(ctx, func() {
		_ = client.Close()
	})
	pingErr := client.Ping(ctx).Err()
	stopClose()
	if err := ctx.Err(); err != nil {
		pingErr = err
	} else if deadline, ok := ctx.Deadline(); pingErr != nil && ok && !time.Now().Before(deadline) {
		// socket deadline 可能先于上下文定时器回调返回；已耗尽启动预算时仍统一保留上下文超时语义。
		pingErr = context.DeadlineExceeded
	}
	if pingErr != nil {
		_ = client.Close()
		if errors.Is(pingErr, context.Canceled) || errors.Is(pingErr, context.DeadlineExceeded) {
			return nil, fmt.Errorf("redis %s startup ping: %w", mode, pingErr)
		}
		return nil, &redisStartupError{mode: mode, cause: pingErr}
	}
	return client, nil
}

func (config RedisConfig) addresses(mode string) ([]string, error) {
	var addresses []string
	switch mode {
	case "node":
		if config.DB < 0 {
			return nil, errors.New("redis DB must be non-negative")
		}
		addresses = []string{config.Addr}
	case "cluster":
		if config.DB != 0 {
			return nil, errors.New("redis cluster requires DB 0")
		}
		if len(config.Addrs) == 0 {
			return nil, errors.New("redis cluster requires at least one seed address in Addrs")
		}
		addresses = config.Addrs
	default:
		return nil, errors.New("redis Type must be node or cluster")
	}

	normalized := make([]string, len(addresses))
	for i, address := range addresses {
		address = strings.TrimSpace(address)
		host, port, err := net.SplitHostPort(address)
		portNumber, portErr := strconv.ParseUint(port, 10, 16)
		if err != nil || host == "" || portErr != nil || portNumber == 0 {
			if mode == "node" {
				return nil, errors.New("redis Addr must be a host:port address")
			}
			return nil, fmt.Errorf("redis Addrs[%d] must be a host:port address", i)
		}
		normalized[i] = address
	}
	return normalized, nil
}

type redisStartupError struct {
	mode  string
	cause error
}

func (err *redisStartupError) Error() string {
	return "redis " + err.mode + " startup ping failed"
}

func (err *redisStartupError) Unwrap() error {
	return err.cause
}
