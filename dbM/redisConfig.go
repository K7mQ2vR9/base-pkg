package dbM

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisTimeoutConfig 定义各阶段超时；nil 使用默认值，显式值须大于 0，YAML 使用带单位时长。
// 请求 ctx 的更早截止时间优先；读写超时是单次 I/O 预算，不是整条命令的总预算。
type RedisTimeoutConfig struct {
	// Startup 覆盖启动 PING 及其连接/认证过程，默认 5s。
	Startup *time.Duration `json:",optional"`
	// Dial 为单次 TCP 建连及 TLS 握手预算，默认 5s。
	Dial *time.Duration `json:",optional"`
	// Read 为单次网络读取超时，默认 3s；阻塞命令仍遵循驱动自身的命令时限规则。
	Read *time.Duration `json:",optional"`
	// Write 为单次网络写入超时，默认 3s，独立于 Read。
	Write *time.Duration `json:",optional"`
	// Pool 为连接池繁忙时等待连接的上限，默认 4s，独立于 Read。
	Pool *time.Duration `json:",optional"`
}

// RedisPoolConfig 管理每个节点的连接池；nil 使用默认值，不表示全体集群的总连接数。
type RedisPoolConfig struct {
	// Size 为池的基础容量，node 默认 10*GOMAXPROCS，cluster 默认 5*GOMAXPROCS，显式值须大于 0。
	// 它不是所有使用方式的连接总数硬上限，额外上限由 MaxActiveConns 配置。
	Size *int `json:",optional"`
	// MaxActiveConns 为驱动每节点连接上限，默认 0 表示不额外限制；非零时须至少为 Size。
	MaxActiveConns *int `json:",optional"`
	// MinIdleConns 为最低空闲连接目标，默认 0，须不超过 Size 和非零 MaxIdleConns。
	MinIdleConns *int `json:",optional"`
	// MaxIdleConns 为最多保留的空闲连接数，默认 0 表示不额外限制，非零时须不超过 Size。
	MaxIdleConns *int `json:",optional"`
	// ConnMaxIdleTime 为连接最大空闲时间，默认 30m，显式 0 关闭按空闲时间回收，负数拒绝。
	ConnMaxIdleTime *time.Duration `json:",optional"`
	// ConnMaxLifetime 为连接最大复用时间，默认 0 不按年龄回收，负数拒绝；驱动按使用时机惰性回收。
	ConnMaxLifetime *time.Duration `json:",optional"`
}

// RedisRetryConfig 配置驱动重发预算；默认关闭，不能提供业务恰好执行一次保证。
type RedisRetryConfig struct {
	// MaxRetries 为每节点命令额外重试次数，默认 0 关闭，负数拒绝；包括流水线，启用前确认幂等性。
	MaxRetries int `json:",optional"`
	// MaxRedirects 仅用于 cluster，默认 0 关闭；非零时同时允许 MOVED/ASK 和部分网络错误重试。
	// node 模式只接受 0；集群转移槽期间默认直接返回重定向错误，业务可显式启用此预算。
	MaxRedirects int `json:",optional"`
	// MinBackoff 为重试最小退避，默认 8ms，显式 0 关闭退避；负数拒绝。
	MinBackoff *time.Duration `json:",optional"`
	// MaxBackoff 为重试最大退避，默认 512ms，须不小于 MinBackoff，负数拒绝。
	MaxBackoff *time.Duration `json:",optional"`
}

// redisSettings 保存校验后的独立配置副本，构造过程不建立连接或修改调用方配置。
type redisSettings struct {
	mode           string
	startupTimeout time.Duration
	options        *redis.UniversalOptions
}

// settings 显式应用默认值并将业务零值转换成驱动哨兵值，避免驱动静默启用默认重试或回收。
func (config RedisConfig) settings() (redisSettings, error) {
	mode := strings.TrimSpace(config.Type)
	if mode == "" {
		mode = "node"
	}
	settings := redisSettings{mode: mode}
	addrs, err := config.addresses(mode)
	if err != nil {
		return settings, err
	}
	if config.InsecureSkipVerify && !config.TLS {
		return settings, errors.New("redis InsecureSkipVerify requires TLS")
	}
	settings.startupTimeout = redisConfigValue(config.Timeout.Startup, 5*time.Second)
	size := 10 * runtime.GOMAXPROCS(0)
	if mode == "cluster" {
		size = 5 * runtime.GOMAXPROCS(0)
	}
	options := &redis.UniversalOptions{
		Addrs:                 addrs,
		Username:              config.Username,
		Password:              config.Password,
		DB:                    config.DB,
		ContextTimeoutEnabled: true,
		DialTimeout:           redisConfigValue(config.Timeout.Dial, 5*time.Second),
		ReadTimeout:           redisConfigValue(config.Timeout.Read, 3*time.Second),
		WriteTimeout:          redisConfigValue(config.Timeout.Write, 3*time.Second),
		PoolTimeout:           redisConfigValue(config.Timeout.Pool, 4*time.Second),
		PoolSize:              redisConfigValue(config.Pool.Size, size),
		MaxActiveConns:        redisConfigValue(config.Pool.MaxActiveConns, 0),
		MinIdleConns:          redisConfigValue(config.Pool.MinIdleConns, 0),
		MaxIdleConns:          redisConfigValue(config.Pool.MaxIdleConns, 0),
		ConnMaxIdleTime:       redisConfigValue(config.Pool.ConnMaxIdleTime, 30*time.Minute),
		ConnMaxLifetime:       redisConfigValue(config.Pool.ConnMaxLifetime, 0),
		MaxRetries:            config.Retry.MaxRetries,
		MaxRedirects:          config.Retry.MaxRedirects,
		MinRetryBackoff:       redisConfigValue(config.Retry.MinBackoff, 8*time.Millisecond),
		MaxRetryBackoff:       redisConfigValue(config.Retry.MaxBackoff, 512*time.Millisecond),
	}
	for _, timeout := range []struct {
		name  string
		value time.Duration
	}{
		{"Timeout.Startup", settings.startupTimeout},
		{"Timeout.Dial", options.DialTimeout},
		{"Timeout.Read", options.ReadTimeout},
		{"Timeout.Write", options.WriteTimeout},
		{"Timeout.Pool", options.PoolTimeout},
	} {
		if timeout.value <= 0 {
			return settings, fmt.Errorf("redis %s must be positive", timeout.name)
		}
	}
	if options.PoolSize <= 0 {
		return settings, errors.New("redis Pool.Size must be positive")
	}
	if options.MinIdleConns < 0 || options.MaxIdleConns < 0 || options.MaxActiveConns < 0 {
		return settings, errors.New("redis pool connection counts must be non-negative")
	}
	if options.MinIdleConns > options.PoolSize || options.MaxIdleConns > options.PoolSize {
		return settings, errors.New("redis idle connection counts must not exceed Pool.Size")
	}
	if options.MaxIdleConns > 0 && options.MinIdleConns > options.MaxIdleConns {
		return settings, errors.New("redis Pool.MinIdleConns must not exceed Pool.MaxIdleConns")
	}
	if options.MaxActiveConns > 0 && options.MaxActiveConns < options.PoolSize {
		return settings, errors.New("redis Pool.MaxActiveConns must be zero or at least Pool.Size")
	}
	if options.ConnMaxIdleTime < 0 || options.ConnMaxLifetime < 0 {
		return settings, errors.New("redis connection lifetimes must be non-negative")
	}
	if options.MaxRetries < 0 || options.MaxRedirects < 0 {
		return settings, errors.New("redis retry counts must be non-negative")
	}
	if mode == "node" && options.MaxRedirects != 0 {
		return settings, errors.New("redis Retry.MaxRedirects requires cluster mode")
	}
	if options.MinRetryBackoff < 0 || options.MaxRetryBackoff < options.MinRetryBackoff {
		return settings, errors.New("redis retry backoff must satisfy 0 <= MinBackoff <= MaxBackoff")
	}
	// 驱动的 0 会开启默认值；使用 -1 表示关闭，对外配置仍使用直观的 0。
	if options.MaxRetries == 0 {
		options.MaxRetries = -1
	}
	if options.MaxRedirects == 0 {
		options.MaxRedirects = -1
	}
	if options.ConnMaxIdleTime == 0 {
		options.ConnMaxIdleTime = -1
	}
	if options.MinRetryBackoff == 0 {
		options.MinRetryBackoff = -1
	}
	if options.MaxRetryBackoff == 0 {
		options.MaxRetryBackoff = -1
	}
	if config.TLS {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: config.InsecureSkipVerify}
		// 使用可取消的 TLS 拨号器，TCP 连接及握手均沿用配置的拨号预算与更早的上下文截止时间。
		dialer := &tls.Dialer{
			NetDialer: &net.Dialer{Timeout: options.DialTimeout, KeepAlive: 5 * time.Minute},
			Config:    options.TLSConfig,
		}
		options.Dialer = dialer.DialContext
	}
	settings.options = options
	return settings, nil
}

// redisConfigValue 仅在指针未配置时采用默认值，保留调用方显式传入的零值。
func redisConfigValue[T any](configured *T, fallback T) T {
	if configured != nil {
		return *configured
	}
	return fallback
}
