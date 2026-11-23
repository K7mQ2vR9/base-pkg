# base-pkg 公共组件开发汇总与接入指南

本文面向后续业务架构和服务开发，汇总当前公共包的能力、接入约定、默认配置和使用边界。整理日期：2026-11-23。按当前实现维护，不按开发时间或优化序号排列。

本地目录为 `common-pkg`，Go 模块为 `github.com/K7mQ2vR9/base-pkg`。新项目使用该模块路径，不继续引用 `github.com/doubunv/common-pkg`。代码声明 `go 1.22.0`、`toolchain go1.22.4`；平台根工作区的 `go.work` 声明 Go 1.25.0，独立使用模块与使用整个工作区时需分别考虑工具链要求。

首版版本号为 `v0.1.0`，本文说明随该标签对应的源码快照发布。后续分支上的修改不自动包含在此版本中，业务服务应固定明确的依赖版本。

在业务项目目录中接入：

```sh
go get github.com/K7mQ2vR9/base-pkg@v0.1.0
```

仓库为私有仓库，使用方需配置对应访问权限及匹配模块路径的GOPRIVATE；本命令不替代私有仓库认证，也不自动迁移旧模块导入。包内 `docs` 提供独立可读的架构和接入说明，不要求同时获取上层平台仓库。

## 平台分层与公共包边界

| 板块 | 职责 | 使用公共包的方式 |
|---|---|---|
| 公共包 | 提供传输、认证基础设施、上下文、响应、数据连接、消息和工具能力 | 不定义具体业务表、业务流程或服务契约 |
| RPC 层 | 提供数据模型服务，负责数据访问及对应的数据一致性处理 | 持有服务级数据库/搜索连接，使用请求级模型和上下文 |
| API 层 | 接收用户操作，编排业务流程并调用 RPC | 接入标准 HTTP 入口、认证、统一响应和 RPC 客户端 |

典型链路为：用户 → 网关 → API → 一个或多个 RPC → MySQL/Redis/搜索。异步处理使用 RabbitMQ 或已有 Kafka 组件。各业务项目自行定义 DTO、protobuf、业务错误、权限、事务与幂等策略。

跨层指南：[平台架构](docs/PLATFORM.md)、[API 接入](docs/API.md)、[RPC 接入](docs/RPC.md)。

## 组件与详细说明索引

以下包路径均相对于 `github.com/K7mQ2vR9/base-pkg`。

| 包/目录 | 当前能力 | 详细说明 |
|---|---|---|
| `app`、`appMiddleware` | 标准 API Server、CORS、请求体限制、认证、并发控制、服务状态 | [CORS 与 Server](app/CORS.md)、[请求体](app/REQUEST_BODY.md)、[认证](appMiddleware/AUTH.md) |
| `headInfo`、`consts` | 请求 ID、业务头、可信代理 IP、metadata、用户 JSON、日志字段 | [Trace](headInfo/TRACE_ID.md)、[HTTP 头](headInfo/HTTP_HEADERS.md)、[metadata](headInfo/METADATA.md)、[日志上下文](headInfo/LOG_CONTEXT.md) |
| `jwtToken` | HS256 JWT 签发和解析，结构体签发参数 | [JWT](jwtToken/README.md) |
| `result`、`result/xcode`、`logs/xcode` | 统一 JSON 响应、错误分类、HTTP/gRPC 错误映射 | [响应与错误接入](result/README.md) |
| `language` | 从项目本地 JSON 词典翻译错误文案及 data 字符串值 | [多语言](language/README.md) |
| `aesGCM` | 进程共用 AES-GCM 配置、随机 nonce、请求解密和响应加密 | [请求协议](app/REQUEST_BODY.md)、[响应协议](result/README.md) |
| `dbM` | MySQL/GORM、多库管理、Redis node/cluster | [MySQL](dbM/SELECT_DB.md)、[MySQL 字段](dbM/mysqlConfig.go)、[Redis](dbM/REDIS.md) |
| `rpc` | 共享一元 RPC 客户端、启动等待、调用预算及错误恢复 | [RPC 客户端](rpc/README.md)、[服务端认证](appMiddleware/AUTH.md) |
| `amqp/rabbitMq` | 发布确认、消费、重连、重试、失败归档及 TG 通知 | [机制与配置](amqp/rabbitMq/README.md)、[业务范例](amqp/rabbitMq/USAGE.md)、[错误策略](amqp/rabbitMq/ERROR_HANDLING.md) |
| `amqp/kafka/v2` | 按分区消费、批处理、连续完成位点提交、DLQ 和关闭收尾 | [消费入口](amqp/kafka/v2/consume.go)、[选项](amqp/kafka/v2/option.go) |
| `es/esv7` | 独立 Elasticsearch 7 客户端和模型 | [ES7](es/esv7/README.md) |
| `es/esaws` | 独立 OpenSearch 客户端、模型和原始 JSON 聚合 | [OpenSearch](es/esaws/README.md) |
| `xxlJob` | Go BEAN 任务注册、调度监听、任务取消及结果回调 | [XXL-JOB](xxlJob/README.md) |
| `commonTool` | HTTP、时间、十进制、随机数、TOTP、RSA、摘要、Redis 分片、UUIDv7 | 见本文工具表 |
| `logs` | go-zero 到 logrus 的同步日志适配 | [实现](logs/logrus.go)、[字段约定](headInfo/LOG_CONTEXT.md) |

## 服务初始化与资源生命周期

服务级资源在启动时构造，通过依赖注入供请求复用；请求仅创建上下文会话或模型，不按请求关闭共享池。

| 资源 | 构造入口 | 上下文与退出要求 |
|---|---|---|
| MySQL | `ConnectMySQL` / `NewSelectDb` | 启动 ctx 只约束初始化；业务 `WithContext` / `GetDb(ctx, key)`；退出关闭底层池或管理器 |
| Redis | `NewRedisClient` | 启动 ctx 只约束初始化；每条命令传业务 ctx；退出 `Close()` |
| RPC 客户端 | `rpc.NewClient` | 启动 ctx 控制连接等待；业务调用继承请求预算；退出 `Close()` 会取消在途 RPC |
| HTTP 工具 | `commonTool.NewHTTPClient` | 每次 `Get` / `PostJSON` 传业务 ctx；收尾后 `Close()` 释放空闲连接 |
| 搜索客户端 | 两套各自的 `core.NewEs` | 启动 ctx 控制探活；按请求创建模型；收尾后 `Close()` |
| RabbitMQ | `NewRabbitMQ` | 启动 ctx 不控制成功后的实例；消费者使用回调 ctx；退出 `Close()` |
| XXL-JOB | `NewClient` → `Register` → `Start` | **Start 的 ctx 控制整个服务生命周期**，不能启动成功后立即 cancel；退出 `Stop(ctx)` 并处理 `Wait` 结果 |

初始化失败必须处理错误，释放已经创建的资源并停止启动。资源关闭顺序由服务管理：停止入口与新任务、等待或取消在途业务、再关闭其依赖。上下文超时不等于写入未发生，跨服务调用、数据库写入和消息重发都需由业务处理幂等。

## API 标准入口

创建服务使用 `app.NewServerWithCORS(restConf, corsConfig, handlerConfig...)`。可选 `ServerHandlerConfig` 最多一项，只开放 404/405 处理器；不透传任意 `rest.RunOption`，不开放路由替换入口。

入口先初始化统一 ID，再校验单值业务头并注入语言，随后进行声明长度检查和 CORS 预检。有效 OPTIONS 返回204；业务路径再进入原生框架中间件，以及平台的 CORS、可选并发限制、正文限制、解密、业务头装配和认证。服务状态检查由业务显式接在认证之后。

`ConfigAppServer` 是配置容器。填写 `RequestBody`、`MaxConcurrency`、`TrustedProxyCIDRs`、`AESGCM` 等字段后，仍须通过对应 `app.With...` 选项装配；数据库、Redis、RabbitMQ和词典也由启动代码显式初始化，不会因字段存在自动启用。完整启动函数见 [API 指南](docs/API.md)。

`ApiUseMiddleware() error` 在配置校验和组件创建完成后才追加中间件。失败停止启动；同一 `ServerMiddleware` 实例成功后再次调用返回错误，不重复装配。此保护不跨实例，同一 Server 保留一个装配对象，使用后不按值复制或在运行期间修改配置。

| 配置 | 默认值与接入约定 |
|---|---|
| `CORS` | 任意来源、关闭浏览器凭证；开放客户端声明的预检请求头；开启凭证须填写精确来源名单 |
| `MaxConcurrency` | 0 关闭；正数为单个 Server 跨路由共享上限，超限立即 HTTP503/code503 |
| `RequestBody.MaxBytes` | 0 使用1MiB；按实际读取字节限制，包含 AES JSON/Base64 封装；超限 HTTP413/code413 |
| `TrustedProxyCIDRs` | 空列表只信任连接对端，不采用用户任意透传的 IP |
| `AESGCM` | 默认关闭；开启时原始密钥须为16、24或32字节；进程共用配置 |
| `Language.Dir` | 空值使用工作目录下 `etc/language`；须显式 `language.Load` |

网关、go-zero 原生 MaxBytes 与平台 RequestBody 的较小限制仍可能先拒绝请求。调大正文配置需同步考虑其他层；大小上限不替代读取超时。自定义 404/405 在入口校验内部运行，但不经过已注册业务路由的认证链，不放置需要登录的业务。

## 统一 ID、传输字段与上下文

Request ID 与 Trace ID 共用一个值。HTTP 对外头为 `X-Request-ID`，响应正文为 `trace`；RPC 与 RabbitMQ 内部仅通过 W3C `traceparent` / `tracestate` 传播链路。不同调用使用不同 Span ID，但不另生成 Trace ID。业务统一调用 `headInfo.GetTrace(ctx)`。

新 Trace ID 使用加密安全随机源生成16字节，编码为非全零32位小写十六进制，无时间戳、业务前缀或连字符。HTTP 优先复用合法入站 ID；只有合法 traceparent 时复用其 Trace ID。二者同时存在必须一致，非法、重复和冲突直接拒绝；缺少入站 ID 时复用已有链路，均不存在才生成。

RPC 必须携带有效 traceparent，白名单和已有服务端 Span 不能豁免此要求。MQ外部消息完全缺少链路时可以初始化独立链路，非法或冲突链路不交给业务；合法消息的重试和通知沿用其 Trace。独立后台工作在入口调用一次 `EnsureTrace`，后续操作复用返回的 ctx。

自定义 HTTP/RPC 传输字段使用小写中划线，HTTP大小写不敏感，RPC使用小写；不兼容下划线旧名。常用字段为 `authorization-jwt`、`business-code`、`content-language`、`device-id`、`token-uid`、`req-path` 和 `rpc-name`，完整表见 [HTTP 字段说明](headInfo/HTTP_HEADERS.md)。原终端标准 HTTP User-Agent 在 RPC 中使用 `client-user-agent`，不占用 gRPC 自身的 `user-agent`。JSON/JWT 标签不随传输头改名。

已知业务头采用单值校验；重复值即使内容相同也拒绝。RPC另校验 UID 为0至MaxInt64的十进制文本，IP为无端口、无zone的纯地址。HTTP中的 UID、角色和客户端 IP 不直接作为可信身份或地址：身份取认证结果，IP按可信代理链从右向左寻找第一个不可信节点；无可信代理配置时使用连接对端。

`HeadInMetadata` 在独立出站副本覆盖平台字段，包括空值，保留其他出站自定义字段和多值；不自动透传全部HTTP头或合并入站metadata。旧W3C字段清理后只从活动上下文注入，不采用 `Head.Trace` 建立链路。Getter读取出站业务metadata；`SetBusinessCode` 同时复制并更新入站、出站和当前日志字段，不修改父上下文或其他分支。

## 认证、JWT与用户信息

白名单仅豁免用户必须登录：无Token或已确认的匿名结果允许UID=0。提供了无效/过期Token、无权限、非法扩展JSON或服务不可用时，即使白名单也返回错误。请求头、正文、解密与已接入的状态检查不因白名单跳过。

API注入身份只使用认证回调结果。支持普通UID、带error以及带用户JSON的三种回调；新业务可使用 `WithCheckTokenHandleWithUserInfoSMOption`。扩展JSON经校验后通过 `headInfo.GetUserInfo(ctx)` 取副本；匿名时清空，不自动传给RPC或写入日志。无效Token不能用 `0,nil` 表示正常匿名。

RPC通过 `NewRpcAuthMiddleware(check)` 接入调用方认证。全部方法包括白名单先做链路/metadata检查和调用方认证，再判断用户是否必须登录。正数UID本身不证明可信，回调须核实服务调用方及携带的Token/身份。服务端接入范围为一元RPC。

JWT使用 `golang-jwt/jwt/v5`，签发入口为 `GenerateToken(GenerateTokenParams)`。HS256密钥至少32字节且不能全为空白；`ExpireDuration` 使用 `time.Duration`，无默认有效期。解析要求签名有效、exp为正数且严格晚于当前Unix秒，未来iat/nbf拒绝，失败返回nil Claims。未自动增加issuer/audience/subject、UID或TokenType业务校验，签发issuer沿用 `test-user`；会话撤销、刷新许可和业务身份策略由服务定义。

## 响应、多语言、加密与日志

标准响应为 `{"code":200,"message":"success","data":{},"trace":"..."}`。业务错误沿用HTTP200及业务code；入口请求头400、正文413、并发503等保持各自HTTP状态。公开错误和内部诊断分离；跨RPC使用标准gRPC状态及ErrorInfo恢复业务码，使用 `errors.Is/As` 分类。

框架错误注册使用 `httpx.SetErrorHandlerCtx(errorhandler.ErrHandlerCtx)`，标准 `app.InitAppServer()` 已注册；输出使用 `httpx.ErrorCtx(r.Context(), w, err)` 或 `result.HttpErrorResult`。不以Background/TODO替换请求ctx，响应阶段不重新生成Trace ID。

多语言从各业务项目 `etc/language/*.json` 加载扁平字符串词典。错误message和data中命中的字符串值自动翻译，字段名不翻译；名为id/code的直接字符串值受保护，字节切片不翻译。缺少语言、词条或词典时保留原值；没有messageKeys映射，不依赖独立language-pkg。加载失败由启动代码处理，不在请求中读取词典文件。

AES-GCM使用进程共用开关和密钥，不通过请求上下文传递加密器。每次加密生成新nonce；请求使用 `aes_data` 封装，成功响应启用AES时data包含 `aes_data`。非空请求体携带business-code时，开启加密的解密组件要求有效封装；其他普通请求按组件既有规则透传。错误响应保持普通JSON，`HttpSuccessResultNotAes` 强制明文成功响应。所有实例的进程配置一致性由应用和运维保证。

响应先完成编码/加密，再提交状态和正文；准备失败使用统一内部错误。写入失败只记录日志，不追加第二份JSON；不提供对任意上游已提交响应的全局恢复或重复写保护。每个handler选择一个出口后立即返回。

日志按 `headInfo.LogKeys` 选择字段，认证后的UID和业务标识在当前分支同步；令牌和用户扩展JSON不自动加入业务上下文日志。logrus适配器不接管外部输出/Hook生命周期，严重日志由调用方决定退出与资源清理。

## 数据访问与默认预算

| 组件 | 启动/调用默认值 | 其他关键默认值 |
|---|---|---|
| MySQL | 启动5s；拨号沿用DSN配置，缺失时5s | 空闲10/最大80；连接生命周期30m、空闲10m；warn日志、慢SQL200ms、参数化日志开启 |
| Redis | 启动5s、拨号5s、读3s、写3s、取连接4s | node；自动重试和集群重定向默认关闭 |
| RPC客户端 | 启动5s、一元调用10s | 默认等待Ready，关闭策略重试，Keepalive默认关闭 |
| HTTP工具 | 调用总预算10s | 成功/失败正文统一受4MiB读取上限约束，不自动跟随重定向 |
| ES7/OpenSearch | 启动5s、模型操作10s | 成功正文16MiB、SDK默认额外重试3次，显式0关闭 |

配置字段按各组件注释解释。指针字段的nil表示未配置，显式0可能表示关闭或非法，不能全局套用同一个规则；时长使用带单位值，容量注明字节或连接数。配置容器定义见 [ConfigAppServer](app/appServer.go)。

MySQL使用 `SelectDbConfig{DbKey, MySQL}` 按库创建独立池。启动探活与版本探测共享预算；业务用 `GetDb(ctx, key)` 获取独立GORM会话，查询和事务继承ctx，管理器不额外增加统一查询超时。DSN读写时限未指定时保留驱动零值，显式连接字段覆盖DSN对应项。池关闭可并发/重复调用，尝试回收全部库并汇总错误。

Redis返回 `redis.UniversalClient`，支持node/cluster，cluster只用DB0，无Sentinel封装。连接池配置按每节点应用；默认Size为node的10×GOMAXPROCS、cluster的5×GOMAXPROCS，MaxActiveConns=0不额外限制。命令须传业务ctx；已阻塞socket I/O不保证在无deadline的ctx主动取消时立即中断。重试和MOVED/ASK预算显式开启前须确认写入幂等与集群迁槽要求。

两套搜索实现独立，不共用Config、Schema、错误或聚合类型。Search固定请求准确total并禁止部分结果；超时、失败分片、提前终止和total下界返回 `ErrIncompleteSearch`。成功正文超限关闭响应并返回 `ErrResponseTooLarge`，不解析截断结果。

Search目标为非nil切片指针，支持 `[]T` / `[]*T`；实现本套Schema的元素回填原生_id，缺少非空_id返回错误。解码及ID回填在新切片完成，成功后替换旧切片；不复用旧元素。OpenSearch同时返回本套 `Aggregations` 原始JSON map，并保留BusinessCode索引后缀。使用filter_path必须保留完整性字段、准确total、_source及Schema需要的_id；详情见各自README。

## RabbitMQ与异步业务

当前使用 `rabbitmq/amqp091-go` 驱动，驱动自动恢复关闭，本包管理连接、拓扑和订阅恢复。业务通过 `PublishJSON`、`Publish`、`PublishBatch` 及 `Consume(Handler)` 接入，不直接操作ACK/NACK。连接和消费者应作为服务级资源复用。

| 配置/语义 | 当前默认值或规则 |
|---|---|
| 业务并发/预取 | WorkerNum=1；PrefetchCount默认跟随WorkerNum，显式0不限制未确认数 |
| 确认与持久化 | AutoAck=false；业务队列、具名交换机和发布消息持久化默认开启 |
| 超时 | DialTimeout=5s、ConfirmTimeout=5s、HandlerTimeout=30s、ShutdownTimeout=10s |
| 重试 | MaxRetries=3，为首次处理之外的次数；RetryDelay=30s |
| 辅助队列 | 业务队列名追加.retry、.failed、.notify，可显式分别配置 |
| 心跳 | 默认10s，显式0按broker协商；配置字段优先于URL中的心跳参数 |
| TG | 默认关闭；启用时单次发送5s、失败重试间隔30s并遵守retry_after |

消息ID与Trace ID独立：可传业务自增ID的字符串，缺省生成36位UUIDv7。ID不保证业务去重或同毫秒单调递增；结果不确定的重发须沿用稳定ID，业务用唯一键/事务/幂等记录处理重复。批次按顺序等待确认，失败返回已确认前缀和失败位置，不自动重发整个批次。

手动确认始终针对单条消息，multiple=false，不会确认此前失败消息。普通业务错误默认重试；`Permanent`和标准JSON格式/类型/目标错误默认直接归档，`Retryable`可显式覆盖但不解除次数上限。回调须沿用传入ctx；超时只取消ctx，不强制终止Go代码，旧回调退出前不在同worker另开业务重投。

终态失败先确认FailedQueue归档，再确认NotifyQueue通知任务，两份都成功后才ACK源消息。两个发布不构成事务，部分成功、确认丢失或ACK失败可能造成重复归档和通知。TG关闭时通知任务仍保存；启用后通知发送成功才确认通知消息，持续失败会占用当前通知消费者并推迟后续通知。归档队列供排查和业务控制的人工重放，不自动恢复执行。

Kafka保留已有按分区worker、连续完成前缀提交、死信持久化后提交以及有界关闭流程。死信或提交失败不越过失败位点，批量回调可能重试，业务需幂等；各分区并非完全无背压隔离。未实施后续Kafka消息副本隔离改造，不将现有组件描述为已经完成该能力。

## XXL-JOB任务调度

沿用现有 `/run`、`/kill`、`/beat`、`/idleBeat`、`/log` 和数组结果回调协议，只执行已注册的Go BEAN任务，不执行GLUE脚本。Start先监听并注册；生命周期ctx取消会触发停止，Stop有界收尾，Wait报告最终状态。

默认RequestTimeout=5s、RegistryInterval=20s、ShutdownTimeout=10s、MaxBodyBytes=1MiB、MaxTasks=1024。MaxTasks包括运行、排队和等待回调；同JobID按配置执行串行、丢弃后续或覆盖策略。kill/超时/覆盖只发出取消，不能强制结束忽略ctx的任务。

任务和回调状态在内存中，LogID仅对当前未完成任务去重。进程重启不恢复队列或持久化去重，回调失败只记录日志，无自动重试队列；业务是否幂等和是否补偿由服务自行定义。普通字符串返回表示成功，取消、超时或panic按失败回调；实际调度中心仍需按所用协议联调。

## 通用工具选择

| 需求 | 入口和使用约定 | 说明 |
|---|---|---|
| 外部HTTP | `NewHTTPClient`、`Get`、`PostJSON`；共享连接池、传ctx，响应已读取并关闭，错误时检查BodyComplete | [HTTP](commonTool/HTTP_CLIENT.md) |
| 时间/统计日界 | `NewTimeTool`；未指定时区、零值及nil实例使用UTC；日/月边界按日历计算，不固定减24小时 | [时间](commonTool/TIME_TOOL.md) |
| 十进制截断 | `FloorDecimal`；字符串输入、0～18小数位、向负无穷取整，负数不向零取整 | [十进制](commonTool/DECIMAL_TOOL.md) |
| 随机业务文本 | `GenerateRandomString`；长度1～4096、加密安全等概率字母数字，不保证唯一 | [随机字符串](commonTool/RANDOM_STRING.md) |
| UUIDv7 | `NewUUidV7`；36字符，用于消息或文档ID，处理随机源错误，不替代Trace生成规则 | [实现](commonTool/uuidv7.go) |
| TOTP | 生成密钥、构造otpauth URI、校验六位字符串；SHA1/30秒，当前及上个窗口，无状态 | [TOTP](commonTool/TOTP.md) |
| RSA | 初始化后复用Signer/Verifier；至少2048位，SHA256+PKCS#1 v1.5，原始字节签名 | [RSA](commonTool/RSA.md) |
| 内容摘要 | `MD5Hex` / `SHA256Hex`；不提供身份认证或密码存储 | [摘要](commonTool/HASH_TOOL.md) |
| Redis逻辑分片 | `ShardIndexXX` / `ShardTagXX`；分片数必须大于0，同标签用于同槽，修改桶数需考虑迁移 | [分片](commonTool/REDIS_SHARDS.md) |

TOTP防重放/尝试次数、RSA防重放、业务唯一性、消息与数据事务一致性都不由这些无状态工具自动完成。HTTP工具也不自动注入平台链路头，调用外部系统时按目标协议显式传递已有ID。

## 现有边界与未实施能力

以下内容是架构选型时需要保留的事实，不是待自动执行的优化清单：

- 未增加全局API panic恢复及任意重复响应写入保护；已有局部处理不等于覆盖全部输出链。
- 未增加Kafka回调消息副本隔离、TG通知重试队列重构、XXL-JOB结果回调重试或持久化恢复。
- RabbitMQ队列/路由名称的额外配置预校验未实施，业务按目标broker规则控制最终名称及默认后缀长度。
- 独立 `NewConcurrencyLimiter` 未改为返回error，传负数可能panic、0会拒绝全部请求；`NewApiHeaderMiddleware` 不能传nil选项。标准装配对并发值已有校验并跳过0值组件。
- RPC客户端沿用现有默认明文传输，未增加TLS/服务凭证配置封装；流式调用未统一封装一元拦截器的认证、预算和错误传播能力。
- 进程共用AES配置的一致性由应用/运维保持，没有另加一致性检查或请求上下文加密器。
- 历史凭证的源码清理不代表部署已完成轮换，实际用途及轮换范围仍需部署方核对；生产配置、真实网关、数据库/消息/搜索集群和实际调度中心的接入由部署与业务项目落实。

## 接口选择、验证状态与维护方式

新代码使用当前公开入口：`NewServerWithCORS`、`ErrHandlerCtx`、结构体参数 `GenerateToken`、带ctx的资源构造和业务操作。旧 `NewServer`、公开路由包装入口、`WithTestOption`、无上下文 `ErrHandler` 以及旧Trace别名均已删除，不为它们补兼容分支。`SetBusinessCode` 是当前同时更新双向metadata的入口。

早期部分组件曾进行定向、竞态、构建及范例检查；随后多次接口、依赖和实现修改按用户要求没有再次构建、测试或运行验证。因此当前整体是源码汇总，不能以早期结果代表当前全部改动已验证。JWT v5、RabbitMQ驱动迁移、搜索完整性/上限/ID回填、装配防重和入口清理等近期变动没有当前运行验证结论，生产联调也未由本次整理补做。

首版发布包含公共组件源码、开发汇总、架构指南及各组件使用说明，不另建开发流水。本次发布准备仅更新版本说明和包内文档，不修改Go业务代码或依赖；按用户要求不执行构建、测试或范例，不新增测试和临时产物。原有组件说明、测试及其他业务项目文件保留。

后续开发仍按 [平台开发规则](docs/DEVELOPMENT_RULES.md) 先讨论方案，再修改实现；变更同步更新本指南及组件文档。开发过程的验证和清理事实按规则记录，但不把临时修复编号、被替代方案或过去的验证结果当作业务接入契约。
