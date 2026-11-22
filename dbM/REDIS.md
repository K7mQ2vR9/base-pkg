# Redis 客户端使用说明

`NewRedisClient(ctx, RedisConfig)` 创建 `redis.UniversalClient`，支持 node 和 cluster。服务启动时创建并共享实例，退出时先停止新业务并按需收尾，再调用 `Close()`；不要按请求创建或关闭连接池。API 通过 RPC 数据模型服务组织业务，实际 Redis 读写由对应数据服务持有客户端。

## 基础配置

`Type` 未填为 `node`；node 使用必填 `Addr: host:port`，cluster 使用非空 `Addrs` 种子列表。`Username`、`Password` 按服务端认证要求配置；node 的 `DB` 默认为 0 且不能为负，cluster 必须为 0。当前不提供 Sentinel 模式。

`TLS` 默认 false；启用时至少 TLS 1.2，使用系统证书信任并校验证书。`InsecureSkipVerify` 默认 false，只在 TLS=true 时可显式启用，开启后跳过证书校验。TCP 和 TLS 握手共用配置的拨号预算，也受更早的上下文截止时间约束。

## 独立配置与默认值

`Timeout`、`Pool`、`Retry` 是三个独立结构体，可不配置。指针字段未配置取默认值，显式 0 按下面的规则校验或关闭功能；不直接向使用方暴露驱动的 -1/-2 哨兵值。YAML 时长使用带单位的字符串，如 `5s`、`200ms`、`30m`、`0s`。

### Timeout / RedisTimeoutConfig

| 字段 | 默认值 | 作用 |
|---|---|---|
| Startup | 5s | 启动 PING 总预算，覆盖连接和认证 |
| Dial | 5s | 单次 TCP 建连及 TLS 握手预算 |
| Read | 3s | 单次网络读取超时 |
| Write | 3s | 单次网络写入超时，独立于 Read |
| Pool | 4s | 连接池繁忙时等待连接的上限，独立于 Read |

以上时长显式配置必须大于 0；不接受关闭超时的负值或 0。启动上下文与 Startup 取更早截止时间，启动失败/取消回收客户端。启动成功后取消启动上下文，不关闭已返回的客户端。

业务命令须传当前请求 ctx，可通过 `context.WithTimeout` 设置整次业务预算。`ContextTimeoutEnabled` 固定开启，socket I/O 使用 ctx 截止时间与配置时限中更早的值；等待连接和重试退避也感知上下文。**已经阻塞的 socket I/O 不保证在无截止时间 ctx 被提前取消时立即中断**，通常仍受既有读写时限约束；不要为了取消某次请求而关闭共享客户端。驱动也可能返回 `net.Error` 超时，并非所有阶段都返回 `context.DeadlineExceeded`。

阻塞命令（如 BLPOP）使用驱动自己的命令时限规则，Read 不是所有命令的统一总预算；重试、多次 I/O 也不等于一个 Read 时长。请按业务为请求设置总截止时间。

### Pool / RedisPoolConfig

| 字段 | 默认值 | 规则 |
|---|---|---|
| Size | node: 10 × GOMAXPROCS；cluster: 5 × GOMAXPROCS | 基础池容量，显式值必须大于 0 |
| MaxActiveConns | 0 | 驱动每节点连接上限，0 不额外限制；非零须至少为 Size |
| MinIdleConns | 0 | 最低空闲目标，非负且不大于 Size |
| MaxIdleConns | 0 | 空闲保留上限，0 不额外限制；非零不大于 Size，且不小于 MinIdleConns |
| ConnMaxIdleTime | 30m | 最大空闲时间，0s 关闭按空闲时间回收 |
| ConnMaxLifetime | 0s | 最大复用时间，0s 关闭按年龄回收 |

连接回收时长不能为负。驱动会在使用时机惰性检查过期连接，不承诺到期瞬间回收。Size 不是全部连接使用方式的硬上限，需要额外限制时配置 MaxActiveConns，并按业务评估专用连接等用法。

**cluster 中上述容量均按每个节点应用**，不是整个集群共享一个总额；副本/节点数及多实例部署都会影响总连接数。原有两种模式的 Size 默认倍率保留。

### Retry / RedisRetryConfig

| 字段 | 默认值 | 规则 |
|---|---|---|
| MaxRetries | 0 | 每节点命令额外重试次数，0 关闭，非负 |
| MaxRedirects | 0 | cluster 外层重定向/重试预算，0 关闭，非负；node 只能填 0 |
| MinBackoff | 8ms | 最小退避，非负；0s 关闭退避 |
| MaxBackoff | 512ms | 最大退避，非负且不小于 MinBackoff |

默认将两个重发预算都关闭，避免单机模式悄然采用驱动默认三次重试。**行为变化：原有配置未填写 Retry 时，node 的命令重试及 cluster 的外层自动重试/重定向现在关闭。** 不保留旧行为的兼容开关；需要时直接显式填写正数。

当前固定依赖的 go-redis v9.7.0 使用同一个 MaxRedirects 循环处理 MOVED/ASK 和部分网络错误，因此不能把它解释为“只做安全路由，不重发命令”。默认 0 时，槽迁移或拓扑缓存落后可能直接返回 MOVED/ASK；客户端仍进行必要的拓扑发现，但不自动沿错误继续发送本条命令。业务需要故障切换/重定向恢复时，可以显式启用预算，并同时考虑重复执行。

启用 MaxRetries 或 MaxRedirects 可能在响应丢失后重发写命令，包括 INCR、列表操作及流水线。集群内外层预算可能叠加；它们都不是业务幂等保障。不保证恰好一次执行，也不能根据超时/断连推断写入未发生。连接重新建立和后台拓扑发现不等于重复发送业务写命令。

## YAML 示例

以下沿用 [ConfigAppServer.RedisConf](../app/appServer.go)，地址仅为配置示例，密码通过部署配置注入。

```yaml
RedisConf:
  Type: node
  Addr: 127.0.0.1:6379
  DB: 0
  TLS: false
  Timeout:
    Startup: 5s
    Dial: 3s
    Read: 2s
    Write: 2s
    Pool: 500ms
  Pool:
    Size: 20
    MaxActiveConns: 40
    MinIdleConns: 2
    MaxIdleConns: 10
    ConnMaxIdleTime: 10m
    ConnMaxLifetime: 30m
  Retry:
    MaxRetries: 0
    MaxRedirects: 0
    MinBackoff: 8ms
    MaxBackoff: 512ms
```

cluster 时改为 `Type: cluster`，用 `Addrs` 替代 `Addr`，保持 DB=0，并按节点评估池容量。仅在业务允许重发时显式增加 MaxRetries/MaxRedirects。

## Go 示例

函数签名与返回接口保持不变。最小配置即可使用全部默认值：

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/K7mQ2vR9/base-pkg/dbM"
)

func main() {
    startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    client, err := dbM.NewRedisClient(startupCtx, dbM.RedisConfig{
        Addr: os.Getenv("REDIS_ADDR"),
        Password: os.Getenv("REDIS_PASSWORD"),
    })
    cancel()
    if err != nil {
        log.Print("Redis 初始化失败")
        return
    }
    defer func() {
        if err := client.Close(); err != nil {
            log.Print("Redis 关闭失败")
        }
    }()

    // 实际服务注入并复用 client；请求沿用当前 RPC 上下文，而非重新使用启动 ctx。
    ctx, release := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer release()
    if err := client.Ping(ctx).Err(); err != nil {
        log.Print("Redis 请求失败")
    }
}
```

Go 直接配置时，为可选数值传指针；下面的显式零值关闭空闲回收和重试退避。

```go
package config

import (
    "time"

    "github.com/K7mQ2vR9/base-pkg/dbM"
)

// Redis 返回单个服务的配置，零回收时长与未配置含义不同。
func Redis(addr string) dbM.RedisConfig {
    size, active := 20, 40
    read, zero := 2*time.Second, time.Duration(0)
    return dbM.RedisConfig{
        Addr: addr,
        Timeout: dbM.RedisTimeoutConfig{Read: &read},
        Pool: dbM.RedisPoolConfig{
            Size: &size, MaxActiveConns: &active, ConnMaxIdleTime: &zero,
        },
        Retry: dbM.RedisRetryConfig{MinBackoff: &zero, MaxBackoff: &zero},
    }
}
```

配置错误在建连前返回。启动网络/认证错误使用阶段描述并保留底层错误链，通过 `errors.Is`/`errors.As` 判断；底层诊断可能包含服务端或地址信息，不应直接返回给用户。命令执行仍返回原 go-redis 错误，业务自行按统一错误协议处理。
