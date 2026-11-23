# RPC层架构与公共包接入

RPC层向API及其他服务提供数据模型能力，持有对应的数据连接、模型、事务和数据访问规则。API负责编排用户业务，RPC按明确契约提供服务。基础组件总览见 [base-pkg汇总](../README.md)，本目录不表示现有RPC项目已迁移到新模块。

## 资源与模型组织

| 生命周期 | 放置内容 | 使用约定 |
|---|---|---|
| 服务级 | MySQL多库管理器、Redis、搜索、下游RPC、HTTP与消息客户端 | 启动初始化，依赖注入，退出统一关闭 |
| 请求级 | GORM会话、搜索模型、调用预算、业务标识及身份上下文 | 从当前RPC ctx派生，不共享可变请求状态 |
| 业务持久状态 | 数据表、唯一约束、幂等记录、消息任务与补偿信息 | 由业务定义并控制事务边界 |

MySQL用 `NewSelectDb(startupCtx, config)` 创建多库管理器，请求中 `GetDb(ctx, dbKey)` 后执行查询或事务。未知库、未初始化及关闭错误用errors.Is分类；不要关闭请求会话背后的共享sql.DB。

Redis用 `NewRedisClient` 创建UniversalClient，请求命令使用ctx。node/cluster配置、每节点连接池以及默认关闭的重试/重定向策略见 [Redis说明](../dbM/REDIS.md)。超时或断连不证明写入未发生，启用重试前先明确幂等策略。

搜索每个请求创建本套 `model.NewEsModel`，复用服务级客户端。ES7与OpenSearch保持独立，不混用配置、Schema和错误；OpenSearch额外接收businessCode并按现有规则追加索引后缀。

## 服务端认证接入

服务端一元请求的处理顺序为：metadata/W3C检查 → 单值及UID/IP检查 → 取消检查 → 调用方认证 → 匿名白名单判断 → 业务方法。

```go
package service

import (
    "context"

    "github.com/K7mQ2vR9/base-pkg/appMiddleware"
    "google.golang.org/grpc"
)

// RPCAuthOption 接入业务提供的调用方认证，并明确可匿名的方法。
// 白名单仅豁免用户必须登录，不跳过调用方、Trace或metadata校验。
func RPCAuthOption(check appMiddleware.CheckRPCRequestFunc) (grpc.ServerOption, error) {
    auth, err := appMiddleware.NewRpcAuthMiddleware(check)
    if err != nil {
        return nil, err
    }
    return auth.ServerOption(func(context.Context) map[string]int {
        return map[string]int{
            "/catalog.CatalogService/ListPublic": 1,
        }
    }), nil
}
```

原生gRPC将返回选项传给 `grpc.NewServer(option)`，go-zero使用 `server.AddOptions(option)`；同一Server只设置一次主UnaryInterceptor。回调必须核实服务调用方及其携带的Token/身份，不能只凭正数UID认定可信。

缺少有效traceparent直接返回InvalidArgument，即使方法在白名单也不放行；UID缺失/空/0为匿名，仅白名单允许。身份或权限错误不能被白名单吞掉。后台任务无最终用户身份时，需要在服务契约中明确其调用方认证和匿名方法范围。

当前封装以一元RPC为范围，不把其认证、错误映射与调用预算自动视为覆盖流式RPC。完整规则见 [服务端认证](../appMiddleware/AUTH.md)。

## 下游调用、metadata和错误

下游RPC在启动时通过 `rpc.NewClient` 构造共享连接，交给生成客户端；默认连接等待5s，一元调用预算10s。启动ctx只控制初始化，业务方法继续传当前ctx；连接Ready只表示传输就绪，不证明业务健康或认证通过。

Request ID=Trace ID，RPC间只传播 `traceparent` / `tracestate`。自定义metadata用小写中划线，原终端信息用 `client-user-agent`。Getter读取出站metadata，入口已建立双向独立副本；变更业务标识使用 `SetBusinessCode` 的返回ctx，不自行共享可变MD切片。

API的用户扩展JSON不自动进入RPC。额外数据应由明确的protobuf或传输契约定义；未知出站metadata虽可保留，也不应作为隐式权限证明或敏感数据通道。

业务错误使用 `result/xcode` 定义明确公开的业务码和Message，服务端编码为标准gRPC状态及ErrorInfo；API客户端恢复业务协议。内部诊断留在日志，通过errors.Is/As区分未找到、取消、超时或驱动状态，不匹配错误文案。

目前RPC客户端未增加TLS及服务凭证配置封装，传输安全与身份机制需结合服务部署设计。自动策略重试关闭，连接重连和gRPC透明恢复不等于业务重试承诺。

## 搜索结果契约

- ES7 Search返回总数和error，OpenSearch同时返回本套Aggregations；跨RPC输出业务DTO，不泄露SDK响应或驱动类型。
- Search固定准确total及禁止部分结果；失败分片、超时、提前终止或下界总数返回ErrIncompleteSearch，不当作正常空数据。
- 成功正文默认16MiB，超限返回ErrResponseTooLarge；需要大查询时显式配置并考虑分页及业务对象内存。
- 非nil切片指针支持[]T和[]*T，Schema元素自动回填原生_id；成功替换为新切片，不沿用旧元素缺失字段的值。
- filter_path保留完整性、准确total、_source及Schema需要的_id。错误时先处理error，OpenSearch失败不使用本次聚合。

详细配置与范例：[ES7](../es/esv7/README.md)、[OpenSearch](../es/esaws/README.md)。

## 消息、任务与关闭

RabbitMQ业务处理器接收独立消息快照，由公共包执行单条确认、重试和终态归档。回调的ctx传给数据库、Redis、HTTP及下游RPC；成功返回前完成本次应持久化的业务工作，写入副作用用稳定业务标识保证幂等。

业务重试、失败归档和TG通知分别有其语义。终态源消息在归档及通知任务分别确认后才ACK，两次发布非原子，重复归档和通知需要业务识别。消息ID、数据库自增ID和Trace ID各司其职。

XXL-JOB用服务生命周期ctx启动，任务及结果回调在内存中；回调失败无自动重试，进程重启无待执行任务恢复。Kafka沿用现有按分区的连续完成位点提交及DLQ收尾，没有实施后续消息副本隔离增强。

关闭时先停止入口和新任务，按业务要求等待/取消在途工作，再关闭资源。不要在某个请求超时时关闭共享连接池。MySQL、Redis、HTTP、搜索与RPC的关闭语义不同，尤其RPC Close会取消在途调用；XXL-JOB Stop和RabbitMQ Close使用各自的关闭预算。

## 维护与验证状态

本指南保留当前模型服务接入约定，替代此前按公共包变更逐日追加的镜像流水。表结构、protobuf、数据隔离、事务、幂等和补偿仍由业务项目明确。后续开发遵循 [平台规则](DEVELOPMENT_RULES.md)，同步维护公共包和本项目说明。

本次仅整理文档，未修改其他RPC项目源码或依赖，未编译范例或执行测试。当前公共包近期修改未再次做运行验证，不将历史定向测试等同于当前服务联调结论。
