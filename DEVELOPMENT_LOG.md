# common-pkg 开发记录

本文件记录平台子项目 `common-pkg` 的开发进度、验证结果和待确认事项。日期采用 `Asia/Ho_Chi_Minh`，按日期追加记录并关联修复序号。

按下表序号逐项修改。每项完成后记录具体改动和实际验证结果，等用户确认后再开始下一项。

## 平台架构

| 板块 | 职责 |
|---|---|
| `common-pkg` | 平台公共包，为 RPC 和 API 等项目提供基础框架能力。 |
| RPC 层 | 对上层提供数据模型服务，通过 RPC 接口提供数据读取和写入等操作。 |
| API 层 | 面向用户提供 API，在这一层组织业务流程，调用 RPC 数据模型服务完成业务操作。 |

用户请求进入 API 层，API 层按业务流程调用 RPC 层；API 和 RPC 各自复用 `common-pkg` 提供的基础能力。

当前请求链路统一采用 Request ID = Trace ID，生成和接入规则见 [统一请求 ID 规则](headInfo/TRACE_ID.md)。历史记录中分别生成或允许任意 Request ID 的约定已由此规则替换。

## 当前基线

- 建档基线为 `pre` 分支上的 `2b03a40`。
- Go 模块声明为 `github.com/K7mQ2vR9/base-pkg`，首版 `v0.0.0` 已发布到同名 GitHub 私有仓库。本地 `common-pkg/` 仍是旧仓库的工作区。

## 待办与进度

| 序号 | 修复项 | 状态 |
|---|---|---|
| 1 | JWT 签发、验签和刷新校验 | 已修改、验证并确认 |
| 2 | 删除源码与仓库中的历史凭证，确认轮换范围 | 工作区清理已验证并确认，部署凭证用途待核对 |
| 3 | 请求与 RPC 日志中的凭证脱敏 | 已修改、验证并确认 |
| 4 | 客户端 IP 的可信代理校验 | 已修改、验证并确认 |
| 5 | Kafka v2 死信处理和位点提交 | 已修改、验证并确认 |
| 6 | Kafka v2 消费和关闭流程 | 已修改、验证并确认 |
| 7 | Go 模块路径与独立发版 | `v0.0.0` 已发布并验证下载，用户要求跳过后续处理 |
| 8 | 认证中间件的回调、白名单和错误处理 | 已修改、验证并确认，包含登录用户信息扩展 |
| 9 | 统一响应和特殊错误匹配 | 已修改、验证并确认，包含跨 RPC Trace ID |
| 10 | 多语言转换和 JSON 序列化兼容 | 已恢复自动匹配、验证并确认 |
| 11 | 请求头校验、CORS 和并发上限 | 已完成确认范围内的 CORS 和并发上限调整 |
| 12 | 数据库连接错误、Redis 客户端类型和日志关闭 | MySQL、Redis 和日志关闭均已完成、验证并确认 |
| 13 | AES-GCM nonce、密钥校验和加密失败处理 | 已完成、验证并确认，使用进程级配置开关 |
| 14 | gRPC 客户端成功分支、服务端认证和状态码 | 已完成、验证并确认 |
| 15 | RabbitMQ 队列绑定、重连和重新订阅 | 已实现重试、TG 通知、统一链路、处理超时、业务消息接口、消费错误分类和终态双份确认，已验证；后续完善逐项确认 |
| 16 | 搜索客户端 TLS 与响应体关闭、XXL-JOB 监听生命周期 | Elasticsearch 7、OpenSearch、XXL-JOB 已完成并确认 |
| 17 | HTTP 请求、随机字符串、TOTP URI 和 Redis 分片工具 | 已完成、验证并确认 |
| 18 | 搜索模型 UUID 旧引用与全量编译恢复 | 已完成、验证并确认 |
| 19 | 时间工具独立管理、UTC 默认值与日历边界 | 已完成、验证并确认 |
| 20 | RSA 密钥初始化、签名与验签 | 已完成、验证并确认 |
| 21 | RPC 客户端启动、调用超时与连接配置 | 已完成、验证并确认 |
| 22 | 通用十进制精度与边界校验 | 已完成、验证并确认 |
| 23 | 哈希摘要工具与直接标准库依赖 | 已完成、验证并确认 |
| 24 | MySQL 多库管理器状态与并发关闭 | 已完成、验证并确认 |
| 25 | Redis 超时、每节点连接池与重试配置 | 已完成、验证并确认 |
| 26 | API 请求体读取、大小限制与统一失败响应 | 已完成、验证并确认 |
| 27 | API 响应编码失败与单次写入处理 | 已完成、验证并确认 |
| 28 | 业务 metadata 单值、UID 与 IP 校验 | 已完成、验证并确认 |
| 29 | HTTP 单值头校验与自定义传输字段中划线命名 | 已完成、验证并确认 |
| 30 | 多语言恢复项目本地 etc/language 词典 | 已完成、验证并确认 |
| 31 | 白名单仅豁免用户必须登录、统一校验与错误返回 | 已完成、验证并确认 |
| 32 | 内置错误文案与项目本地词典对齐 | 已完成、验证并确认 |
| 33 | 认证结果、业务标识与上下文日志字段同步 | 已完成、验证并确认 |
| 34 | 认证前错误响应的语言提前注入与统一读取 | 已完成、验证并确认 |
| 35 | API panic 恢复与响应重复写入处理 | 已按用户要求撤回，不继续开发 |
| 36 | JWT 普通解析与刷新有效期校验统一 | 已完成、验证并确认 |
| 37 | GenerateToken 结构体参数与明确时长类型 | 已完成、验证并确认 |
| 38 | 请求头装配保留已有出站自定义 metadata | 已完成、验证并确认 |
| 39 | GetFullHead 返回独立请求头副本 | 已修改并确认，按用户要求不做测试 |
| 40 | CORS 凭证模式与实例来源配置 | 已完成并确认，配置统一通过 NewServerWithCORS 接入 |
| 41 | 收紧 Server 构造参数，防止路由选项覆盖统一入口 | 已按确认修改，按用户要求不做测试和验证 |
| 42 | 入站 Trace 复用与活动链路冲突统一拒绝 | 已按确认修改，沿用用户不做测试和验证的要求 |
| 43 | API中间件先校验后装配，初始化失败返回error | 已按确认修改，按用户要求不做测试和验证 |
| 44 | 独立服务状态中间件初始化失败返回error | 已按确认修改，沿用用户不做测试和验证的要求 |
| 45 | Kafka回调消息副本与重试源数据隔离 | 用户要求跳过Kafka，未修改 |
| 46 | Elasticsearch 7/OpenSearch搜索结果完整性校验 | 已按确认修改，沿用用户不做测试和验证的要求 |
| 47 | Elasticsearch 7/OpenSearch成功响应读取上限 | 已按确认修改，沿用用户不做测试和验证的要求 |
| 48 | JWT迁移至golang-jwt/jwt/v5并保持平台校验规则 | 已按确认修改，沿用用户不做测试和验证的要求 |
| 49 | RabbitMQ驱动迁移至rabbitmq/amqp091-go | 已按确认修改，沿用用户不做测试和验证的要求 |
| 50 | Elasticsearch 7/OpenSearch搜索结果回填原生文档ID | 已按确认修改，沿用用户不做测试和验证的要求 |
| 51 | 同一ServerMiddleware实例防止重复装配 | 已按确认修改，沿用用户不做测试和验证的要求 |
| 52 | HeadInMetadata业务字段与W3C链路共用独立副本 | 已按确认修改，沿用用户不做测试和验证的要求 |
| 53 | 删除无实际行为的WithTestOption和isTest字段 | 已按确认修改，沿用用户不做测试和验证的要求 |
| 54 | XXL-JOB任务结果回调重试 | 用户明确要求跳过，不处理 |
| 55 | 删除无请求上下文的HTTP错误处理入口ErrHandler | 已按确认修改，沿用用户不做测试和验证的要求 |

## 开发记录

### 2026-11-17 / DEV-20261117-001 / 第 1 项 JWT 签发、验签和刷新校验

- **状态** 已修改并完成定向验证，待用户确认。
- **改动** [`jwtToken/jwt.go`](jwtToken/jwt.go) 统一将调用方配置的密钥转换为 `[]byte`，拒绝空白或少于 32 字节的密钥，验签仅允许 `HS256`。刷新只接受未过期的管理员或用户访问令牌。保留三个公开函数的参数和返回类型。
- **验证** 使用临时 Go overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> ./jwtToken -count=1`，通过。覆盖两个有效令牌类型、错误密钥、空白和过短密钥、非 `HS256` 算法、未知或缺失令牌类型、过期或缺少过期时间。`git diff --check` 通过。临时测试和缓存已清理。
- **下一步** 等用户确认后开始第 2 项凭证清理。

### 2026-11-17 / DEV-20261117-002 / 确认平台分层和子项目职责

- **用户确认** 平台由公共 `common-pkg`、RPC 数据模型服务和面向用户的 API 三个板块协作组成。API 层调用 RPC 层实现业务逻辑。
- **改动** 在本文件记录三层职责和调用关系，明确本开发记录的范围为 `common-pkg` 子项目。
- **修复进度** 第 1 项 JWT 仍待确认，第 2 项凭证清理尚未开始。

### 2026-11-17 / DEV-20261117-003 / 第 2 项历史凭证清理

- **进度** 用户要求进入下一项，第 1 项记为已确认。第 2 项工作区清理和验证已完成，等待确认。
- **改动** 删除 `jwtToken/jwt.go` 中注释里的旧 JWT 签名密钥，以及仓库根目录的 `y`、`y.pub` SSH 密钥文件。OpenSearch 集成测试改从 `COMMON_PKG_OPENSEARCH_TEST_URL`、`COMMON_PKG_OPENSEARCH_TEST_USERNAME`、`COMMON_PKG_OPENSEARCH_TEST_PASSWORD` 读取配置，缺失时跳过。XXL-JOB 示例改从 `XXL_JOB_ADDRESS` 和 `XXL_JOB_ACCESS_TOKEN` 读取连接配置，缺失时退出。
- **验证** 使用本地依赖缓存执行 `go build -mod=readonly -buildvcs=false ./...`，构建通过。`go test -mod=readonly -buildvcs=false -v ./es/test -run '^TestLogInfo$' -count=1` 验证空配置时跳过。构建并运行 XXL-JOB 示例，确认缺少环境变量时拒绝启动。对照清理前提交，静态核查已移除的凭证值未出现在当前工作区文件，`git diff --check` 通过。未连接外部服务。
- **测试清理** 未创建测试用例。登记验证结果后已删除临时构建缓存和 XXL-JOB 示例二进制，并核对清理结果。
- **待核对** 历史 Git 提交和已发布或下载的归档中可能仍有旧凭证。部署或外部系统是否仍信任旧 JWT 密钥、SSH 密钥、OpenSearch 账号和 XXL-JOB Token，需要由对应环境的管理方确认。仍有效的凭证需要撤销或轮换。本次未改写 Git 历史或变更部署配置。
- **下一步** 等用户确认第 2 项后处理第 3 项日志凭证脱敏。

### 2026-11-17 / DEV-20261117-004 / 第 3 项请求与 RPC 日志凭证脱敏

- **进度** 用户要求进入下一项，第 2 项工作区清理记为已确认，部署凭证的用途和轮换范围仍待核对。第 3 项已修改并验证，等待确认。
- **改动** `headInfo` 移除 HTTP 日志上下文中的 JWT 字段，`Head.String()` 从凭证已清空的副本生成日志字符串。RPC 服务端按允许列表记录链路、站点、请求路径等已知字段，排除凭证字段及未列入的自定义 metadata。HTTP 认证和 RPC 处理器继续读取原始 Token，RPC 入站及出站 metadata 保持原有内容。
- **验证** 使用临时 Go overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -v ./appMiddleware -run '^Test(HTTPLogRedaction|RPCMetadataLogRedaction)$' -count=1`，两项通过。覆盖 `Head.String()`、HTTP 认证及日志字段、RPC 日志允许列表及原始 metadata 保留。测试使用内存请求和模拟凭证，启动时提示缺少 `etc/language`，未执行外部服务联调。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。
- **测试清理** 登记结果后已删除临时测试用例、overlay 配置和构建缓存，并核对工作区无本次新增的测试文件。
- **下一步** 等用户确认第 3 项后处理第 4 项可信代理校验。

### 2026-11-17 / DEV-20261117-005 / 第 4 项客户端 IP 可信代理校验

- **进度** 用户要求进入下一项，第 3 项记为已确认。第 4 项已修改并验证，等待确认。
- **改动** HTTP 入口默认使用连接对端 IP。仅当对端命中已配置的可信代理 CIDR 时，从右向左检查 `X-Forwarded-For` 链，取第一个不可信的有效 IP。若整条链均为可信节点，则取最左侧有效 IP。可信段内出现无效内容时回退到连接对端。`ClientIp` 和 `AdminClientIp` 都使用同一解析结果，`GetAdminClientIp` 拒绝包含多个地址或格式无效的 metadata。
- **接入** 可在 `app.ConfigAppServer` 中配置 `TrustedProxyCIDRs`，调用 `app.WithTrustedProxyCIDRs(config.TrustedProxyCIDRs)` 检查错误后，将返回的选项传入 `app.NewServerMiddleware`。单独使用 HTTP 头中间件时可调用 `appMiddleware.WithTrustedProxyCIDRs`。只配置会在 `X-Forwarded-For` 末尾追加连接对端 IP 的代理。IP 解析采用 `X-Forwarded-For`，原有 `CF-Connecting-IP`、`TZ-Connecting-IP` 和 `X-Real-IP` 不参与取值。
- **验证** 使用临时 Go overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -v ./app -run '^Test(TrustedProxyResolution|TrustedProxyMiddlewareConfiguration|AdminClientIPMetadata)$' -count=1`，三组通过。覆盖默认直连、可信与不可信对端、多级与多行转发链、IPv6、IPv4 映射地址、异常链回退、配置校验及 HTTP metadata 传递。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。测试使用内存请求，启动时提示缺少 `etc/language`，未连接实际代理或变更部署配置。
- **测试清理** 登记结果后已删除临时测试用例、overlay 配置和构建缓存，并核对工作区无本次新增的测试文件。
- **下一步** 等用户确认第 4 项后处理第 5 项 Kafka v2 死信和位点提交。

### 2026-11-17 / DEV-20261117-006 / 更新公共包 Go 模块路径

- **用户要求** 将 `headInfo` 的导入路径改为 `github.com/K7mQ2vR9/base-pkg/headInfo`。
- **改动** 将 `go.mod` 的模块路径改为 `github.com/K7mQ2vR9/base-pkg`，同步更新包内 Go 源码和现有测试的导入路径。正式发布仍待处理。
- **验证** `go list -mod=readonly -f '{{.ImportPath}}' ./headInfo` 返回新路径。`go build -mod=readonly -buildvcs=false ./...` 通过。`go test -mod=readonly -buildvcs=false ./... -run '^$' -count=1` 编译了测试包，未执行测试用例。独立临时 Go 模块通过本地 `replace` 导入并调用新路径的 `headInfo.GetHead`，运行通过。`git diff --check` 通过。
- **测试清理** 登记验证结果后删除临时导入验证模块及构建缓存。核对工作区无本次新增的测试文件。
- **限制与下一步** 本次验证使用本地源码，Git `origin` 仍为 `github.com/K7mQ2vR9/common-pkg`。远端仓库和版本发布未调整，尚未验证通过 GitHub 下载新模块。第 4 项仍等待确认。

### 2026-11-17 / DEV-20261117-007 / 第 5 项 Kafka v2 死信与位点提交

- **进度** 用户要求进入下一项，第 4 项记为已确认。第 5 项已修改并验证，等待确认。
- **改动** 分区 worker 只提交按拉取顺序连续完成的消息。业务重试耗尽后，死信投递失败会保留失败消息和后续消息；重试该消息时只重试死信投递。已完成的连续前缀必须等 broker 确认提交后才清除。提交失败后优先重试提交，避免在同一 worker 内重复执行业务或投递死信。失败期间暂停接收新消息到该 worker 的批次，保留批次边界。单条消费入口直接按顺序处理消息。
- **接入** 消费者须设置 GroupID。可通过 `WithDLQTopic` 配置同集群中已存在的独立死信主题，内置生产者同步等待 Kafka 副本确认，并在死信消息附上原主题、分区和位点。`WithDLQHandler` 优先，回调应在持久化成功后返回 nil。没有配置死信处理时，处理失败的消息保持未提交。
- **验证** 使用临时 Go overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -v ./amqp/kafka/v2 -run '^Test(DLQFailureBlocksLaterOffsets|CommitFailureRetainsProgress|MissingDLQDoesNotCommit|KafkaDLQPublisher)$' -count=1`，四组通过。覆盖死信失败后提交已完成前缀、恢复后继续、批量处理和死信各自成功但提交失败、缺少 DLQ、生产者配置及消息封装。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。未连接真实 Kafka broker，未验证重平衡或进程退出后的重投。批量处理器仍需幂等，进程崩溃或 broker 回执丢失时可能重投。
- **测试清理** 登记结果后已删除临时测试用例、overlay 配置和构建缓存，并核对工作区无本次新增的测试文件。
- **下一步** 等用户确认第 5 项后处理第 6 项 Kafka v2 消费和关闭流程。

### 2026-11-17 / DEV-20261117-008 / 第 6 项 Kafka v2 消费和关闭流程

- **进度** 用户要求进入下一项，第 5 项记为已确认。第 6 项已修改并验证，等待确认。
- **改动** 消费者只允许运行一个拉取与分发循环。`Close()` 取消拉取和投递并等待收尾，由消费循环退出后关闭 worker 队列。worker 按分区排空已接收的消息，在 Reader 和 DLQ Writer 仍可用时提交已完成位点。Reader 返回 `io.EOF` 时进入收尾，临时读取错误的重试等待可被取消。重复或并发关闭等待同一次收尾，Reader 和 Writer 的关闭错误纳入消费返回结果。
- **接入** `WithShutdownTimeout` 可设置缓冲处理的时限，默认 10 秒。超时会取消 worker 上下文并由消费方法返回超时错误，未完成的位点保持未提交。业务及死信回调需响应上下文取消。`Close()` 由消费回调之外的调用方调用。
- **验证** 使用临时 Go overlay，执行 `go test -race -mod=readonly -buildvcs=false -overlay <临时配置> ./amqp/kafka/v2 -run '^Test(CloseStopsBlockedDispatchBeforeDraining|ReaderEOFDrainsPartialBatch|ShutdownTimeoutCancelsWorkerWithoutCommit|CancelInterruptsFetchRetryAndDuplicateStart|CallerDeadlineAndCloseBeforeStart)$' -count=3`，五组通过。另以相同参数运行 `-run '^TestShutdownRetriesPendingDLQBeforeDraining$' -count=3`，验证关闭时先恢复死信并提交，再处理后续缓冲消息。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。
- **限制** 使用模拟 Reader 和回调，未连接真实 Kafka broker。多个分区共享拉取与分发循环，故障分区缓冲填满时仍会产生跨分区背压。与 Kafka broker 的重平衡和关闭期间的重投待联调。
- **测试清理** 登记结果后已删除临时测试用例、overlay 配置和构建缓存，并核对工作区无本次新增的测试文件。
- **下一步** 等用户确认第 6 项后处理第 7 项剩余的远端和独立发版工作。

### 2026-11-17 / DEV-20261117-009 / 第 7 项发布 base-pkg v0.0.0

- **进度** 用户要求进入下一项，第 6 项记为已确认。第 7 项已发布并验证，等待确认。
- **发布** 新建 [K7mQ2vR9/base-pkg](https://github.com/K7mQ2vR9/base-pkg) 私有仓库。由本地源码提交 `f753d52` 生成快照，排除 `.DS_Store`，以独立的首个提交 `4726dbab6bcfcc481d769322fd5252033a005805` 创建远端 `main` 和 `v0.0.0`。旧提交历史和旧 tag 留在原仓库。
- **验证** 发布前检查快照中的原 JWT、OpenSearch、XXL-JOB 和 SSH 凭证特征，通过。快照的 `go build -mod=readonly -buildvcs=false ./...`、`go test -mod=readonly -buildvcs=false ./... -run '^$' -count=1` 和 Git 差异检查通过。测试命令只编译测试包。远端 `git ls-remote` 确认 `main` 和 `v0.0.0` 指向发布提交。首次下载触发环境中已有的公开校验库配置，校验失败。以只针对本模块的 `GOPRIVATE`、`GONOPROXY` 和 `GONOSUMDB` 例外及 SSH 权限重跑，`go mod download -json github.com/K7mQ2vR9/base-pkg@v0.0.0` 从空缓存下载成功，`go list -m -f '{{.Version}}' github.com/K7mQ2vR9/base-pkg@latest` 返回 `v0.0.0`。
- **用户决定** 收到跳过第 7 项的指令时，远端仓库、`main` 和 `v0.0.0` 已经创建，且已验证下载。按最新要求停止第 7 项的后续操作，保留已发布的版本。
- **测试清理** 登记验证后清理了测试缓存、临时发布快照和验收下载目录。Go 下载缓存中的只读文件通过 `GOMODCACHE=<临时目录> go clean -modcache` 清理。核对临时目录不存在，工作区无本次新增测试文件。
- **当前状态** 本地旧工作区的 `origin` 仍指向 `K7mQ2vR9/common-pkg`。第 8 项认证中间件调整待开始。

### 2026-11-17 / DEV-20261117-010 / 第 8 项认证中间件统一

- **进度** 用户要求开始第 8 项。代码已修改并通过定向验证，等待确认。
- **改动** 旧版和带错误回调的认证共用同一个流程。白名单支持原有的精确匹配和尾部 `*` 前缀匹配。未提供 Token 或认证失败时以匿名身份继续，并清除上下文中原有的用户 ID。认证成功时写入已验证的用户 ID。受保护接口会拒绝认证失败，响应使用固定提示，特殊错误沿用原业务码，权限拒绝不包含内部路径。
- **接入** `app.WithCheckTokenHandleWithErrSMOption` 启用带错误的回调，旧的 `WithCheckTokenHandleSMOption` 继续可用。两种选项都传入时以后配置的为准。`app` 装配和各认证、服务状态中间件在初始化时检查必需回调。服务状态检查沿用白名单，独立调用时缺少回调返回服务不可用。
- **验证** 使用临时 Go overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -v ./app -run '^Test(AuthWhiteListAndAppOptions|AuthErrorsUseSafeMessages|MissingChecksFailDuringAssembly|ServiceStatusGuard)$' -count=1`，四组通过。覆盖新旧装配、匿名与已认证请求、白名单匹配、旧身份清除、特殊码、错误文本隔离、启动时缺少回调及状态检查。测试启动时提示缺少 `etc/language`，验证在本地内存请求中完成。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。
- **测试清理** 登记验证结果后已删除临时测试用例、overlay 配置和构建缓存，并核对工作区无本次新增测试文件。
- **下一步** 等用户确认后开始第 9 项统一响应和特殊错误匹配。

### 2026-11-18 / DEV-20261118-001 / 第 8 项补充登录用户信息 JSON

- **用户要求** 认证中间件注入当前登录用户的扩展 JSON，供业务使用。
- **改动** 新增 `app.WithCheckTokenHandleWithUserInfoSMOption` 和 `appMiddleware.WithCheckOptionWithUserInfo`，回调可返回用户 ID、可选的 `json.RawMessage` 和错误。只有认证成功且扩展 JSON 格式有效时才注入。业务通过 `headInfo.GetUserInfo(ctx)` 获取独立副本，范围为当前 HTTP 请求上下文，避免自动进入 RPC metadata、请求头及日志字段。匿名访问和使用旧回调时清空旧扩展，现有认证回调签名继续可用。
- **验证** 使用临时 Go overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -v ./app -run '^Test(AuthenticatedUserInfoFromCallback|InvalidAuthUserInfoIsRejected|ExistingAuthCallbacksStillWork)$' -count=1`，三组通过。覆盖认证后注入、忽略外部同名头和 metadata、扩展副本隔离、匿名清除、非法 JSON 拒绝以及旧回调兼容。测试启动时提示缺少 `etc/language`，验证在本地内存中完成。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。
- **测试清理** 登记验证结果后已删除临时测试用例、overlay 配置和构建缓存，并核对工作区无本次新增测试文件。
- **下一步** 等用户确认第 8 项后开始第 9 项统一响应和特殊错误匹配。

### 2026-11-18 / DEV-20261118-002 / 第 9 项统一 HTTP 错误响应和特殊错误匹配

- **进度** 用户要求进入下一项，第 8 项记为已确认。第 9 项已修改并验证，等待确认。
- **改动** `result.HttpErrorResult` 和 `logs/xcode.ErrHandlerCtx` 共用 HTTP 错误响应生成器。`app.InitAppServer` 注册带上下文的 go-zero 全局错误处理器。沿用 HTTP 200 和响应体业务码，统一 `code`、`message`、`data`、`trace`。普通错误和 5xx 错误仅展示固定提示；业务错误读取显式公开的 `Message()`，底层诊断留在服务端日志。错误响应仅翻译 `message`。旧 `ErrHandler` 签名继续可用，调用方需使用带上下文入口才能获得请求 `trace`。
- **特殊错误** `errors.Is` 按特殊错误码匹配，`IsSpecialNotifyError` 仅识别具有通知码 410 的特殊错误，支持包裹和解包底层原因。特殊错误仅公开通用文案。原解密失败和请求缺少加密数据使用明确的 400 响应体业务码及固定提示。
- **验证** 使用临时 overlay 执行 `go test -c -mod=readonly -buildvcs=false -overlay <临时配置> -o <临时二进制> ./app`，在临时语言目录中运行 `./app-response.test -test.run '^Test(HTTPErrorResponsesAgree|SpecialErrorsMatchByCode)$' -test.count=1`，两组通过。覆盖中间件与 go-zero 输出一致、普通和 5xx 错误脱敏、业务错误和特殊错误码、`trace`、文案翻译、空错误、错误链及解密提示。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。验证采用内存 HTTP 请求，未发布或联调外部服务。
- **测试清理** 登记验证结果后已删除临时测试用例、语言样本、overlay 配置、测试二进制和构建缓存，并核对工作区无本次新增测试文件。
- **下一步** 等用户确认第 9 项后开始第 10 项多语言和 JSON 序列化兼容。

### 2026-11-18 / DEV-20261118-003 / 第 9 项补充跨 API 与 RPC Trace ID

- **用户要求** 同一 API 请求及其多个 RPC 调用使用一致的 Trace ID。
- **改动** HTTP 入口沿用有效的 OpenTelemetry 上下文或 W3C `traceparent`，两者均缺失时初始化链路。`X-Request-ID` 作为独立的 `request_id` 保留。RPC 客户端向下游注入 `traceparent` 和同值的兼容 `trace` metadata，服务端以 W3C 链路为准规范入站和出站 `trace`。响应及 go-zero 结构化日志使用同一个 OpenTelemetry Trace ID。同步修复 RPC 成功调用后空错误引用触发的 panic，原错误映射沿用既有规则。
- **验证** 使用临时 overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> ./appMiddleware -run '^Test(TraceAcrossHTTPAndMultipleRPCs|ExistingW3CTraceIsPreserved)$' -count=3`，通过。内存 gRPC 覆盖 API 分别调用 RPC A、RPC B 及 A 再调用 B，并对比 HTTP 成功与错误响应、RPC metadata 和结构化日志。首轮测试捕获器未实现 Error 方法而失败，补齐后通过。`go build -mod=readonly -buildvcs=false ./...` 和 `git diff --check` 通过。未连接部署环境或验证链路导出。
- **测试清理** 登记验证后已删除临时测试用例、overlay 配置及构建缓存，并核对临时目录已移除。
- **后续开发** 用户明确要求以后先讨论方案，经确认后再开发。已记录在平台 `AGENTS.md`。第 10 项仍待方案讨论。

### 2026-11-18 / DEV-20261118-004 / 核对网关请求 ID 的 API 与 RPC 传递

- **用户确认** 网关优先保留前端提供的非空 `X-Request-ID`，缺失或为空时生成一次再传给 API。请求 ID 与 OpenTelemetry Trace ID 分别传递。
- **现有实现** API 优先读取 `X-Request-ID`，以 `request_id` 写入 RPC metadata。RPC 收到后继续向下游传递同一个值。现有传递逻辑已满足该规则的 API 和 RPC 部分。
- **验证** 使用临时模拟网关和内存 RPC 拦截器，执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -v ./appMiddleware -run '^TestGatewayRequestIDAcrossAPIAndRPC$' -count=3`，通过。覆盖保留前端 ID、缺失时透传网关生成的 ID、连续请求生成不同的 ID、一次 API 请求内分别调用两个 RPC 以及嵌套 RPC，验证 `trace` 独立且一致。`git diff --check` 通过。测试启动时提示缺少 `etc/language`，未连接实际网关。
- **测试清理** 登记验证后删除临时测试用例、overlay 配置和构建缓存，并核对临时目录已移除。
- **待接入** 当前平台工作区未找到网关配置。待确认实际网关种类和配置路径后接入生成规则并联调。

### 2026-11-18 / DEV-20261118-005 / 第 10 项多语言和 JSON 序列化

- **用户确认** 新建独立的 [language-pkg](../language-pkg/DEVELOPMENT_LOG.md)。从 `content_language` 读取请求语言，统一翻译错误 `message` 和业务显式选中的 `data` 字段值。
- **改动** `TranslateMessage` 将常用错误文案映射为稳定 Key，也接受词典 Key。原有 `SwitchLanguage` 入口改为调用业务类型的 `Localize(locale string) any`，其他响应直接交给原有编码流程。明文和 AES 响应共用翻译入口，JSON 编码由业务响应类型负责。
- **迁移** 如已有部署使用 `etc/language`，将词条迁移到 `language-pkg/locales` 的 JSON 文件。成功响应过去依赖自动遍历翻译的字段，需要由业务显式选择。
- **验证** 使用临时 overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^Test(LocalizedResponseData|UnmarkedResponseData|LocalizedErrorResponse|DictionaryLookup)$' -count=1 -v ./common-pkg/result ./language-pkg`，四组通过。覆盖中英文和缺少词条时的回退、AES 解密后字段值、原始业务状态、`json:"-"`、`omitempty`、自定义 `MarshalJSON`、空响应及错误码与 Trace ID。工作区 `go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...`、`go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1` 和 `git diff --check` 通过。
- **验证限制** 现有测试只做全量编译，临时用例在内存 HTTP 中运行。默认构建缓存不可写，最初工作区依赖解析也缺少本地版本和一个公开间接依赖。使用工作区 `replace`、临时缓存和 Go 校验和验证的公开依赖完成构建。模块独立发布和服务联调待进行。
- **测试清理** 登记验证后删除本次临时测试、overlay 配置、下载文件和 Go 缓存。核对临时目录已移除，项目原有五个测试文件保留。

### 2026-11-18 / DEV-20261118-006 / 第 10 项恢复自动匹配词典

- **用户确认** 响应数据恢复自动翻译。业务在展示字段中填入词典 Key，成功响应按请求语言匹配，沿用 `id` 和 `code` 的跳过规则。
- **改动** `SwitchLanguage` 先按响应类型的规则编码结构化数据，再只替换匹配词典的 JSON 字符串值。匹配覆盖嵌套对象和数组，按完整字符串匹配，字段名及大小写不敏感的 `id` 和 `code` 字符串值保持原样。原有字段标签、自定义 `MarshalJSON`、原始数字精度以及 JSON 中未匹配的值由原序列化结果保留。顶层字符串可直接翻译，字节切片沿用原先的明文和 AES 编码方式。移除第 10 项早期版本引入的 `Localize` 方法要求。
- **接入** 业务返回 `{"status":"pending","status_label":"order.status.pending"}`，中文响应输出 `"status_label":"待支付"`，英文输出 `"status_label":"Pending payment"`。未知语言或未配置的 Key 保留原值。
- **验证** 使用临时 overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestAutomatic(DataLocalization|ScalarResponses|LocalizationErrorEnvelope)$' -count=1 -v ./common-pkg/result`，三组通过。覆盖中英文、空语言和未知语言、明文与 AES、嵌套匹配、字段名和受保护字段、`omitempty`、`json:"-"`、自定义序列化、大数和指数原值、顶层字符串与字节、空值、错误码及 Trace ID。
- **验证限制** 工作区 `go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...` 和 `go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1` 通过。现有测试只做全量编译，定向测试使用内存 HTTP 和测试密钥。工作区从本机缓存和临时代理读取依赖，缺少的公开依赖从 Go 公共代理下载并由 `go.work.sum` 校验。发布和部署联调待进行。
- **测试清理** 登记验证后已删除临时测试、overlay、下载文件、日志及 Go 缓存。核对临时目录已移除，两模块无本次新增的测试文件，原有五个测试文件保留。

### 2026-11-18 / DEV-20261118-007 / 第 11 项 CORS 请求头和实例并发上限

- **用户确认** CORS 允许浏览器在预检中列出的全部请求头。API 实例共享一个可配置的并发上限。超限立即返回 HTTP `503`，预检在并发限制之前处理。
- **CORS** `OPTIONS` 响应返回 `204`，把请求的 `Access-Control-Request-Headers` 写入 `Access-Control-Allow-Headers`，并加上 `Vary: Access-Control-Request-Headers`。已覆盖 `Authorization` 和业务自定义头。已注册的 `OPTIONS` 路由经过 API 中间件，无 `OPTIONS` 路由的已注册路径可在创建 go-zero Server 时接入 `rest.WithNotAllowedHandler(appMiddleware.NewCorsMiddleware().Handler())`。该回退处理器对其他方法返回 `404`。已有自定义回退的服务可用 `CorsMiddleware.Handle` 包装原回退逻辑。业务请求的必填头校验与认证继续由各自的中间件处理。
- **并发** `ConfigAppServer.MaxConcurrency` 默认 `0`，通过 `app.WithMaxConcurrency` 接入 API 中间件。正数表示当前 Server 实例跨路由共享的名额，处理结束后释放。负数返回配置错误。超限使用统一 JSON 响应体，`code` 为 `503`，`message` 可按请求语言翻译，`trace` 复用当前请求的 Trace ID。
- **接入顺序** 先调用 `app.WithMaxConcurrency(appConf.MaxConcurrency)` 并处理返回的错误。实例级限制启用时，在创建 go-zero Server 前设置 `restConf.Middlewares.MaxConns = false`。创建 Server 时配置需要的 CORS 回退，然后将并发选项传入 `app.NewServerMiddleware(server, ...)` 并调用 `ApiUseMiddleware()`。
- **验证** 使用临时 overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^Test(SharedConcurrencyLimitAcrossRoutes|ConcurrencyAndCORSConfiguration)$' -count=3 -v ./common-pkg/app`，两组各连续通过三次。覆盖两路由共享两个名额、第三个请求的状态码及响应体、中文文案和 Trace ID、释放名额后放行、饱和时的预检、显式 `OPTIONS` 路由与 go-zero 回退、多个请求头字段、默认值和负数配置。
- **验证限制** `go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...`、`go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1` 及 `git diff --check` 通过。定向验证使用实际 go-zero 路由和内存 HTTP 请求，现有测试只做全量编译，真实浏览器、网关和部署并发值待联调。工作区依赖使用临时缓存，缺少的公开间接依赖由 `go.work.sum` 校验。
- **测试清理** 登记验证后已删除临时用例、overlay、下载文件、日志和缓存。核对临时目录不存在，工作区保留原有五个测试文件。

### 2026-11-18 / DEV-20261118-008 / 第 12 项 MySQL 连接错误和上下文

- **用户确认** 新项目直接替换旧 MySQL 接口。本轮先完成 MySQL 连接和上下文。
- **启动** `ConnectMySQL(ctx, dataSource, logLevel)` 通过同一个启动上下文执行 `PingContext` 和 `SELECT VERSION()`。每个连接池最多有 5 秒，调用方的更早截止时间或取消优先。GORM 用检测到的服务端版本配置方言。错误保留可供 `errors.Is` 或 `errors.As` 判断的原因，默认错误文本省略 DSN 和驱动的原始错误。
- **接入** 用 `NewSelectDb(startupCtx, config.SelectDb, logLevel)` 初始化连接池。空或重复的 `DbKey` 提前报错。任一库初始化失败会关闭已打开的池并返回带 `DbKey` 的错误。业务使用 `GetDb(requestCtx, dbKey)` 获取绑定请求上下文的 GORM 会话，查询和事务继承该上下文。服务退出时调用 `SelectDb.Close()`。
- **验证** 使用临时 overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestMySQL(StartupDeadline|SelectionCleanup|RequestContext)$' -count=1 -v ./common-pkg/dbM`，三组通过。覆盖探活和版本探测卡住、默认 5 秒和调用方更短的超时、第二库失败后的池清理、错误脱敏、版本和 DSN 时区保留、启动上下文结束后继续查询、请求取消和查询超时。工作区 `go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...`、`go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1` 及 Git 差异检查通过。
- **验证限制** 本机执行 `mysqld --no-defaults --initialize-insecure --datadir=<临时目录> --log-error=<临时日志> --tmpdir=<临时目录>` 时 MySQL 9.7.1 初始化进程收到 SIGSEGV。定向测试改用真实 Go MySQL 驱动连接内存协议模拟端，真实 MySQL 联调待进行。现有测试只做全量编译。
- **测试清理** 登记验证后已删除临时用例、overlay、MySQL 临时数据目录和日志、下载文件与 Go 缓存。核对临时目录已移除，原有五个测试文件保留。

### 2026-11-18 / DEV-20261118-009 / 第 12 项 Redis 单机和集群客户端

- **用户确认** 延续新项目规则，直接替换旧的 `RedisConnect`，按配置创建单机或集群客户端。
- **配置** `ConfigAppServer.RedisConf` 接收 `dbM.RedisConfig`。默认 `Type: node`，通过 `Addr` 配置单机地址和非负 `DB`。`Type: cluster` 使用 `Addrs` 中的种子地址并要求 `DB: 0`。一个种子地址也会创建集群客户端。地址采用 `host:port`。开启 `TLS` 后默认校验证书，`InsecureSkipVerify` 仅在明确启用 TLS 时可设置。
- **接入** 有 Redis 配置的服务在启动时调用 `dbM.NewRedisClient(startupCtx, *config.RedisConf)`，处理返回的 `error`。启动 PING 的超时上限为 5 秒并继承更早的调用方截止时间。失败时关闭客户端，错误文本省略原始错误中的配置和服务端细节，原因可通过错误链检查。业务通过 `redis.UniversalClient` 执行命令并传入请求上下文。客户端启用命令 I/O deadline 支持。服务退出时调用客户端的 `Close()`。
- **验证** 使用 `go test -c -mod=readonly -buildvcs=false -overlay <临时配置> -o <临时二进制> ./common-pkg/dbM` 编译定向测试，在本地回环环境执行 `<临时二进制> -test.run '^TestRedis(NodeAndCluster|StartupFailure|ConfigAndTLS)$' -test.v -test.timeout=35s`，三组通过。使用真实 Redis 单机和单节点集群验证客户端类型、单种子集群、认证和独立 DB、读写、业务请求 deadline 和客户端关闭。模拟端验证 PING 超时与失败后的连接回收、错误脱敏，以及 TLS 默认拒绝未受信任证书和显式跳过校验。首轮 TLS 用例限定了底层证书错误类型，更新临时证书和断言后通过。
- **验证限制** `go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...`、`go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1` 及 Git 差异检查通过。现有测试只做全量编译。集群验证使用本机单节点，生产多节点切换和部署配置待联调。
- **测试清理** 定向测试结束时停止并等待临时 Redis 进程退出。登记验证后删除临时用例、overlay、测试二进制、数据、日志、下载文件与 Go 缓存。核对临时目录已移除，原有五个测试文件保留。

### 2026-11-19 / DEV-20261119-001 / 补充应用配置注释

- **用户要求** 为 ConfigAppServer 添加备注，后续改动同步添加注释。
- **改动** 配置结构及各字段增加中文注释，说明默认值、用途和实际接入方式。平台 AGENTS.md 增加后续开发的注释要求。
- **核对** MySQL 当前固定使用 5 秒启动超时、10 个空闲连接、80 个最大连接、30 分钟连接生命周期、10 分钟最大空闲时间。连接参数配置化方案待确认。
- **验证** 执行 gofmt -w app/appServer.go、gofmt -l app/appServer.go 及平台和公共包的 git diff --check，均通过。核对配置字段类型与标签。本次更新为注释和开发规则，未执行构建或运行测试，未创建测试用例及临时产物。

### 2026-11-19 / DEV-20261119-002 / 独立 MySQL 配置与默认值

- **用户确认** MySQL 的连接、连接池、超时和日志参数封装为独立结构体，未配置时走默认值。每个字段补充中文注释。
- **配置与接口** 新增 `dbM.MySQLConfig`。每个 `SelectDbConfig` 通过 `MySQL` 字段组合连接配置。`ConnectMySQL(ctx, config)` 接收 `MySQLConfig`，`NewSelectDb(ctx, configs)` 按各库配置初始化。调用方使用请求上下文调用 `GetDb(ctx, key)`，服务退出时调用 `Close()`。
- **默认值** 启动超时为 5s，驱动拨号超时在 DSN 也未配置时为 5s，读写超时沿用 DSN，默认 0。最大空闲连接为 10、最大连接为 80、连接生命周期为 30m、最大空闲时间为 10m。SQL 日志级别为 warn、慢 SQL 阈值为 200ms、参数占位符默认开启。
- **参数语义** 可选数字、时长及布尔字段用指针区分缺省和显式零值，YAML 和 Go 直接传参应用相同默认值。显式配置的驱动超时覆盖 DSN。最大空闲连接为 0 表示关闭空闲连接保留，最大连接为 0 表示不限数量。连接生命周期、最大空闲时间和慢 SQL 阈值为 0 时，关闭各自的回收或阈值检查。启动超时须为正值，其余时间及连接数须非负。最大连接数为正值时，空闲连接上限须小于或等于它。全部数据库配置在打开第一个连接池前完成校验。
- **验证** 通过临时 overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestMySQLConfig(DefaultsAndYAML|Validation|Connection)$' -count=1 -v ./common-pkg/dbM`，三组通过。覆盖最小 YAML 和 Go 结构体默认值一致、显式 0 和 false、DSN 覆盖、各项非法配置、连接池最大连接数、自定义启动超时、SQL 参数过滤、方言版本和时区保留、启动上下文取消后执行查询，以及连接池关闭。
- **验证限制** `go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...`、`go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1` 和 Git 差异检查通过。定向连接测试使用实际 MySQL 驱动及内存协议模拟端，真实数据库环境的容量压测和联调待进行。现有测试仅做全量编译。
- **测试清理** 登记验证后已删除临时测试用例、overlay、日志、下载文件及 Go 缓存，核对临时目录已移除，原有五个测试文件保留。

最小配置只需填写数据库标识和 DSN，可选字段自动采用默认值。

```yaml
SelectDb:
  - DbKey: main
    MySQL:
      DataSource: "user:password@tcp(127.0.0.1:3306)/platform?parseTime=true&loc=UTC"
```

按库调整参数时，将可选字段写在同一层。

```yaml
SelectDb:
  - DbKey: main
    MySQL:
      DataSource: "user:password@tcp(127.0.0.1:3306)/platform?parseTime=true&loc=UTC"
      StartupTimeout: 5s
      DialTimeout: 5s
      ReadTimeout: 3s
      WriteTimeout: 3s
      MaxIdleConns: 10
      MaxOpenConns: 80
      ConnMaxLifetime: 30m
      ConnMaxIdleTime: 10m
      LogLevel: warn
      SlowThreshold: 200ms
      ParameterizedQueries: true
```

### 2026-11-19 / DEV-20261119-003 / 第 12 项日志关闭

- **用户确认** 修复日志关闭时直接退出进程的问题，验证重复关闭和后续资源清理。
- **改动** LogrusWriter.Close() 正常返回。Logrus 同步写日志，输出及 Hook 的生命周期由配置它们的调用方管理。补充中文注释，说明重复和并发调用的语义。
- **验证** 使用临时 overlay 执行 `go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^Test(LogrusWriterClose|LoggingShutdownContinuesCleanup)$' -count=1 -v ./common-pkg/logs`，两组通过。覆盖重复关闭、16 个并发关闭、已写日志保留、调用方输出资源的生命周期。子进程使用默认退出函数和实际 go-zero 的 logx.Close()，确认正常关闭后仍执行后续代码和 defer 清理。执行 `go build -mod=readonly -buildvcs=false ./common-pkg/logs` 和 Git 差异检查，通过。
- **验证范围** 本轮针对日志包构建和定向测试，使用子进程模拟服务退出收尾。部署环境的完整服务退出流程待联调。
- **测试清理** 登记验证结果后已删除本次临时测试用例、overlay 配置和构建缓存。核对临时目录已移除，项目原有五个测试文件保留。

### 2026-11-19 / DEV-20261119-004 / 第 12 项复查与补充修复

- **用户要求** 重新检查并处理第 12 项。本轮核查 MySQL 配置和请求上下文、Redis 启动取消与清理、日志记录与关闭流程。
- **日志问题与修复** Close() 正常返回的行为通过复核。另发现 Severe() 使用 Logrus.Fatal()，与 go-zero Severe() 的记录行为不一致。修复为使用 FatalLevel 记录并正常返回，使调用方能继续收尾。补充关闭与严重日志的中文注释。
- **Redis 问题与修复** go-redis 的网络读取 deadline 继承了上下文截止时间，但初始化过程中主动取消上下文时，正在进行的读取仍会等待原 I/O 超时。先用临时用例复现取消后仍阻塞的情况，再增加仅在初始化期间生效的 context.AfterFunc 回调，以关闭客户端中断读取。初始化完成即撤销回调，成功的客户端由调用方管理。补充 Redis 配置字段注释。
- **MySQL 核查** YAML 与 Go 直接传参的默认值一致，显式零值与 false 正确保留，结构化超时覆盖 DSN，业务会话绑定请求上下文。查询上下文独立于基础连接池，重复关闭连接池正常返回。
- **定向验证** 临时 overlay 执行日志包的 TestRecheckSevereReturns、TestRecheckCloseAndExtensions、TestRecheckShutdownSubprocess，以及数据库包的 TestRecheckMySQLConfiguration、TestRecheckMySQLRequestContext。子进程验证严重日志后仍能执行正常关闭和 defer。现有的格式化日志、结构化日志和凭证字段过滤通过回归。日志包另以 go test -race 通过。
- **Redis 验证** 编译带 overlay 的 dbM 测试二进制并在本机回环协议端运行 TestRecheckRedisStartupCancellation 和 TestRecheckRedisSuccessfulStartup。覆盖单机与集群初始化取消、连接关闭、成功后取消原启动上下文仍可继续使用客户端，并通过 -race 检测。修复前 Severe() 退出和 Redis 取消阻塞均被用例捕获，修复后通过。
- **执行命令** 以下命令使用本次临时 overlay 和测试二进制。

```text
go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestRecheck(SevereReturns|CloseAndExtensions|ShutdownSubprocess)$' -count=1 -v ./common-pkg/logs
go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestRecheckMySQL(Configuration|RequestContext)$' -count=1 -v ./common-pkg/dbM
go test -race -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestRecheck(SevereReturns|CloseAndExtensions|ShutdownSubprocess)$' -count=1 ./common-pkg/logs
go test -c -race -mod=readonly -buildvcs=false -overlay <临时配置> -o <临时二进制> ./common-pkg/dbM
<临时二进制> -test.run '^TestRecheckRedis(StartupCancellation|SuccessfulStartup)$' -test.v -test.timeout=15s
```

- **构建与限制** go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...、go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1 和 Git 差异检查通过。现有测试做全量编译。本轮数据库请求上下文检查使用 GORM 会话，Redis 使用本地协议端，生产数据库、多节点集群切换和完整部署退出流程待联调。
- **测试清理** 登记验证后已删除临时用例、overlay、测试二进制、日志、下载文件和 Go 缓存。核对临时目录及三份 overlay 目标文件均不存在，原有五个测试文件保留。

### 2026-11-19 / DEV-20261119-005 / 第 13 项 AES-GCM 实例配置及失败处理

- **用户确认** 统一使用随机 nonce，启动校验密钥，按 Server 实例保存配置。序列化及加密失败返回统一错误，完善解密失败和请求体处理。
- **加密实现** 删除固定 nonce 的 EncryptSame，纯函数 Encrypt 和 Decrypt 共用随机 nonce 及认证逻辑。密文采用无 padding 的 URL-safe Base64，字节格式为 nonce、密文和认证标签。解密同时支持带 padding 的输入，并检查 nonce 和认证标签的最小总长度。原有固定 nonce 测试调整为随机 nonce 测试，项目原有测试文件继续保留。
- **实例配置** 新增 aesGCM.Config 和只读 Cipher。Config.Enabled 默认 false，启用后 Config.Key 要求为 16、24 或 32 字节原始 UTF-8 字符串。ConfigAppServer.AESGCM 保存该配置，服务通过 app.WithAESGCM(config.AESGCM) 获取装配选项并处理错误，然后传入 app.NewServerMiddleware。调试模式和加密分别通过各自配置启用。
- **接口迁移** 原来的 WithAesKeyOption、DecryptKeyOption、DecryptWithDebugFalseOption 及全局变量已替换。独立接入解密中间件时调用 NewApiRequestDecryptMiddleware(aesGCM.Config)，处理其返回的错误。加密器通过本地请求上下文供成功响应使用。
- **错误响应** 启用 AES 时，响应序列化失败或随机源读取失败立即返回统一错误响应，响应体业务码 500，沿用 HTTP 200 和当前 Trace ID，客户端仅展示固定错误文案。请求体读取失败、无效密文或认证失败按现有 400 业务码响应。HttpSuccessResultNotAes 继续提供显式明文响应。
- **请求体** 读取原请求体后关闭其 Reader。解密后同步重建 Body、GetBody、ContentLength 和 Content-Length，清除原 TransferEncoding。普通请求按已有规则透传，非空请求体带有 business_code 时要求提供 aes_data。已提供但为空、null 或非字符串的 aes_data 返回解密错误。
- **定向验证** 使用临时 overlay 执行 go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestAES13' -count=1 -v ./common-pkg/aesGCM ./common-pkg/result ./common-pkg/app，五组通过。覆盖 16/24/32 字节密钥、随机 nonce、空明文、带 padding 密文、截断和篡改、非法密钥及短密文。真实 go-zero 路由验证两个不同密钥 Server 与关闭加密的 Server 并发处理请求，解密后的请求正文、响应及 Trace ID 正常。注入序列化和随机源失败验证错误响应，检查请求读取失败、请求体长度及 GetBody。
- **竞态与构建** 同一组临时测试加 -race 通过。合并解密封装判断后，再执行带 -race 的 TestAES13RequestHandling，通过。go test -mod=readonly -buildvcs=false ./common-pkg/aesGCM 运行原有 AES 用例，通过。go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...、go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1 及 Git 差异检查通过。
- **验证限制** 使用内存 HTTP 请求和实际 go-zero 路由，前端加解密及部署配置待联调。原有外部服务测试仅做编译。新增、修改的配置字段和对外接口已补充中文注释。
- **测试清理** 登记验证后删除本次临时用例、overlay、日志、下载文件及 Go 缓存。核对临时目录与三份 overlay 目标文件已移除，项目原有五个测试文件保留。

API 装配时，将已校验的 AES 选项与现有认证等选项一起传入。

```go
aesOption, err := app.WithAESGCM(config.AESGCM)
if err != nil {
    return err
}

middleware := app.NewServerMiddleware(server, aesOption, authOption)
middleware.ApiUseMiddleware()
```

### 2026-11-19 / DEV-20261119-006 / 第 13 项改为配置开关控制

- **用户要求** 加密器通过配置开关控制，移除请求上下文传递。
- **配置逻辑** aesGCM.Configure 校验并应用当前进程共享的 Enabled 和 Key，默认关闭。校验失败时返回错误并保留原配置。开关与已校验的加密器作为同一份只读状态原子发布。请求解密与成功响应加密直接读取配置，因此业务传入 Background 或 TODO 上下文时仍按开关执行加密。
- **接入** ConfigAppServer.AESGCM 和 app.WithAESGCM 的调用方式继续适用。WithAESGCM 先校验配置，在 ApiUseMiddleware 装配 API 时应用。同一进程的 API 共用这份开关和密钥。独立构造 NewApiRequestDecryptMiddleware 时也会应用传入的进程级配置。
- **接口调整** 删除 ContextWithCipher 和 CipherFromContext，仓库 Go 源码中的引用已清理。纯 Encrypt、Decrypt 以及独立 Cipher 继续供直接加解密使用。序列化及加密失败使用第 13 项的统一错误响应。
- **定向验证** 使用临时 overlay 执行 go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestAES(ConfigSwitch|ConfigConcurrentAccess|ResponseWithoutRequestContext|ConfigFailureResponses)$' -count=1 -v ./common-pkg/aesGCM ./common-pkg/result，通过。覆盖开关启停、非法密钥、原子配置读取、Background 和 TODO 上下文返回加密响应、序列化和随机源失败。
- **装配与并发** 带 -race 运行上述用例及 TestAESConfigServerAssembly，通过实际 go-zero 路由验证加密请求、业务丢失上下文后返回加密响应，以及关闭开关后返回普通 JSON。未使用的 WithAESGCM 选项仅校验配置。配置并发发布与读取通过竞态检测。
- **构建与限制** go test -mod=readonly -buildvcs=false ./common-pkg/aesGCM、go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...、go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1 和 Git 差异检查通过。定向验证使用内存 HTTP 与实际 go-zero 路由，前端和部署配置待联调。
- **测试清理** 登记验证后已删除本次临时用例、overlay、日志、下载文件及 Go 缓存。核对临时目录和三份 overlay 目标文件已移除，原有五个测试文件保留。

### 2026-11-19 / DEV-20261119-007 / 第 14 项 RPC 认证和错误传递

- **用户确认** 统一 RPC 白名单和认证回调，采用标准 gRPC 状态码，恢复 API 业务错误，显式复制调用 metadata。
- **服务端认证** NewRpcAuthMiddleware 接收 CheckRPCRequestFunc 并返回中间件和 error，缺少回调时初始化失败。白名单支持精确方法名和末尾 * 前缀，其他方法均执行回调，包括完全缺少 metadata 的请求。回调读取 metadata 或传输层身份并校验调用方。认证失败返回 Unauthenticated，权限不足返回 PermissionDenied，取消、超时和服务不可用使用对应的标准状态。具体服务凭证策略由 RPC 项目提供。
- **错误传递** result/xcode.ToGRPCStatus 将业务错误映射为有效的 gRPC 状态，通过标准 ErrorInfo 携带业务码、公开文案和特殊错误标记。Code.GRPCStatus、SpecialCodeError.GRPCStatus、NewRpc 及 logs/xcode.FromError 共用该规则。普通内部错误、5xx 错误和 panic 使用固定公开文案。诊断信息由服务端日志记录。
- **客户端** 成功调用继续返回 nil。失败时通过 FromGRPCError 恢复业务码、公开文案和特殊通知错误，保留底层错误链。原生取消和超时支持 errors.Is 检查 context.Canceled 和 context.DeadlineExceeded，API 业务码分别为 499 和 504。未认证、权限不足及服务不可用分别映射为 401、403 和 503。原生 gRPC 其他状态按统一映射处理。
- **metadata** 客户端在注入 rpc_name 和链路字段前显式复制 metadata，服务端为入站和出站上下文分别提供副本。当前 gRPC 库的读取函数也会返回副本，测试核对原始 metadata、并发调用各自的名称及相同 Trace ID，未将复制前行为描述为已复现的数据竞争。
- **定向验证** 临时 overlay 执行 go test -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestRPC14' -count=1 ./common-pkg/result/xcode ./common-pkg/appMiddleware ./common-pkg/rpc/interceptors，通过。覆盖业务码往返、特殊通知多级传递、nil 及 typed nil 错误、包装错误的公开文案、原生状态、缺少凭证、伪造 business_code、白名单、权限拒绝和 panic。首轮捕获业务取消文案被标准文案覆盖，增加错误详情的公开文案后通过。
- **真实协议验证** 通过 bufconn 建立实际 gRPC 服务和客户端，验证 API 错误响应、认证拒绝、多级 RPC Trace ID、业务错误与通知往返、调用过程中取消及超时。该组和全部定向用例以 -race 通过。
- **构建与限制** go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...、go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1 及 Git 差异检查通过。原有测试做全量编译。验证范围为一元 RPC，流式 RPC、部署凭证策略及生产联调待对应服务接入时确定。
- **补充验证** 对直接构造零值中间件的情况增加防护检查，受保护调用返回固定 Internal 文案。带 -race 运行 TestRPC14Authentication，通过。
- **测试清理** 登记验证后删除本次临时用例、overlay、日志、下载文件和 Go 缓存。核对临时目录和四份 overlay 目标文件均已移除，原有五个测试文件保留。

go-zero RPC Server 在启动前接入认证拦截器。checkCaller 是服务实现的认证回调，publicMethods 返回公开方法列表，nil 表示全部方法都经过认证。

```go
auth, err := appMiddleware.NewRpcAuthMiddleware(checkCaller)
if err != nil {
    return err
}
server.AddOptions(auth.ServerOption(publicMethods))
```

业务继续返回 xcode.New(code, message) 或特殊通知错误。使用基础包客户端拦截器后，API 可以从返回错误读取原业务码，公开文案由统一 HTTP 响应处理。

### 2026-11-19 / DEV-20261119-008 / 第 14 项 RPC 必传链路信息

- **用户要求** RPC metadata 为空或缺少 Trace ID 时直接返回错误。
- **入口规则** 在白名单、认证回调和业务处理前检查传入 metadata。不存在或为空时返回 InvalidArgument。有效的 traceparent 或兼容 trace 字段可提供 Trace ID，缺失、格式错误、重复值以及两者冲突都返回 InvalidArgument，基础包客户端映射为 API 业务码 400。白名单方法执行同一入口规则。
- **链路处理** 入站校验以实际传入的 metadata 为准，已有服务端 Span 的情况也执行该校验。保留调用方的 Trace ID，仅 trace 字段存在且需要建立 SpanContext 时生成本次接收的 Span ID。接收的 W3C 上下文和兼容 trace 字段继续用于后续 RPC 传递。
- **定向验证** 使用临时 overlay 执行 go test -race -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestRPCTraceGate' -count=1 ./common-pkg/appMiddleware ./common-pkg/rpc/interceptors，通过。覆盖 metadata 缺失及空值、身份字段存在但链路缺失、只有服务端 Span、空和非法 Trace ID、全零 ID、非法 traceparent、重复及冲突链路字段、白名单入口校验和原始 metadata 副本。
- **协议验证** bufconn 上的实际 gRPC 服务验证公开方法也拒绝缺少链路的调用，单独 trace 或 traceparent 能正常调用，基础包客户端注入的链路与服务端收到的 Trace ID 一致。
- **构建与限制** go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...、go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1 和 Git 差异检查通过。本轮覆盖一元 RPC，生产入口和服务凭证配置待联调。
- **测试清理** 登记验证后已删除本次临时用例、overlay、日志、下载文件和 Go 缓存。核对临时目录和两份 overlay 目标文件已移除，原有五个测试文件保留。

### 2026-11-19 / DEV-20261119-009 / 提取 RPC 日志字段名单

- **用户要求** RPC 日志字段名单定义为包级变量，请求处理时复用。
- **改动** 将 contextMetadataInLog 中的字段名单提取为私有包级变量 rpcLogKeys，并补充中文注释。函数使用它创建本次请求的日志字段。
- **验证** 执行 gofmt -w appMiddleware/rpcAuthMiddleware.go 和平台及公共包的 git diff --check，通过。使用 Python 静态对照 HEAD 中原名单与提取后的名单，字段内容和顺序一致，函数内的局部名单声明已移除。
- **验证范围与清理** 本次为字段名单提取，采用语法格式和静态核对，未执行构建或运行测试，未新增测试用例和临时产物。

### 2026-11-19 / DEV-20261119-010 / 日志字段统一到 headInfo

- **用户要求** 日志字段名单放到 headInfo，相关位置统一引用。
- **改动** 新增 headInfo.LogKeys，在包初始化时创建共用名单。RPC 的 metadata 日志注入迁移为 headInfo.ContextMetadataInLog，RPC 中间件改为调用它。HTTP 的 ContextHeadInLog 按同一名单选择本次请求的字段。
- **字段与格式** 静态对照原字段名单，名称及顺序一致。HTTP 字段使用字符串，RPC metadata 字段使用字符串列表。请求字段值由各请求独立构建，已有业务自定义日志上下文字段继续参与日志输出。
- **验证** 临时 overlay 执行 go test -race -mod=readonly -buildvcs=false -overlay <临时配置> -run '^Test(SharedHeadLogKeys|HeadLogKeysSingleSource|HeadLogValuesPerRequest)$' -count=1 ./common-pkg/headInfo，三组通过。覆盖 HTTP 和 RPC 共用名单、字段格式、凭证及未知字段排除、业务自定义字段、原始 metadata 和并发请求的值隔离。测试中临时替换名单，确认两个入口均读取同一来源，测试结束后恢复。
- **构建与范围** go build -mod=readonly -buildvcs=false ./common-pkg/headInfo ./common-pkg/appMiddleware、gofmt 和 Git 差异检查通过。本轮针对相关包构建和日志行为验证。
- **测试清理** 登记验证后已删除临时用例、overlay、日志、下载文件及 Go 缓存。核对临时目录和 overlay 目标文件已移除，原有五个测试文件保留。

### 2026-11-19 / DEV-20261119-011 / 第 15 项 RabbitMQ 重试及 TG 通知

- **用户确认** 进入第 15 项，增加重试队列，并在重试耗尽后发送 TG 通知。第 14 项记为已确认，本项完成后等待确认。
- **配置与接口** 新增独立 `rabbitmq.Config`、`TelegramConfig` 和 `ConfigAppServer.RabbitMQ`。新接口为 `NewRabbitMQ(ctx, config)`、`Consume(Handler)`、`Publish(ctx, amqp.Publishing)`、`PublishBatch`、`Close`；不保留旧接口。配置字段、公开接口和关键逻辑均补充中文注释。完整默认值、消息流程和调用示例见 [RabbitMQ 接入说明](amqp/rabbitMq/README.md)。
- **连接与发布** 初始化声明业务队列并绑定具名交换机，失败立即返回并回收连接。连接、发布通道或订阅失效后重建拓扑和订阅，关闭后停止重连。所有发布启用 mandatory 和 broker 确认，区分未路由、broker 拒绝、连接未就绪和结果不确定。发布超时主动关闭当前连接以中断网络写入，取消原因保留在错误链中。默认持久化消息，并为缺少 ID 的发布生成 ID。
- **重试与归档** 默认首次处理之外重试 3 次，间隔 30s。每次失败先转入独立重试队列，broker 确认后再单条 ACK 原消息；到期定向重投原业务队列，避免对其他交换机订阅者重复广播。重试耗尽后通过专用 fanout 交换机投递失败归档和通知任务，确认后再 ACK 当前消息。转存失败不 ACK 原消息，手动确认模式下重新入队；内部头损坏直接归档，处理器 panic 按失败处理。所有 ACK/NACK 均为 `multiple=false`，不批量确认其他失败消息。
- **TG 通知** 默认关闭，启用时要求 BotToken 和 ChatID。业务注册消费者后启动独立通知消费者，发送服务名、队列、消息 ID、处理次数和错误摘要，不发送正文。成功后单条 ACK 通知任务；失败按独立间隔重试并遵守服务端 retry_after。TG 失败不占用业务消费者，失败归档继续保留；关闭 TG 时通知任务保存在队列，启用后处理积压。
- **运行语义** 默认手动确认。开启 AutoAck 时原消息已自动确认，转存失败无法恢复。断线、确认丢失和 TG 响应丢失可能导致重复，业务按正文中的稳定业务标识幂等处理，生产重发沿用 MessageId。内部转存强制持久化并保留正文、消息 ID 和常用属性，不沿用原 TTL 及生产者 UserId。处理器须响应上下文取消，超过关闭时限会断开连接并返回超时，无法强制终止用户回调。
- **定向验证** 临时 overlay 下 10 组 TestRabbit15 用例通过，并启用 -race。覆盖 YAML 与 Go 默认值、显式零值、协议数值范围、拓扑绑定和发布确认、24 个并发发布、mandatory 返回及 NACK、重试计数和延迟、转存确认先于源 ACK、转存失败不 ACK、单条确认、AutoAck 不重复确认、panic、消费者取消与连接断开后的恢复、初始化失败回收、并发重复关闭和不响应取消的回调超时。补充 net.Pipe 验证阻塞写入 deadline、初始化取消回收、等待发布锁的取消不影响其他发布，以及阻塞发布取消时保留 DeadlineExceeded。
- **通知验证** 使用模拟 HTTP Transport 验证 TG 500 后重试成功、429 的 retry_after、凭证不进入普通错误文本、UTF-16 文本长度限制、失败任务不 ACK，以及 TG 请求阻塞时业务消息继续处理。未向实际 TG 发送消息。
- **执行命令** 以下使用独立临时 GOCACHE 和 `GOPROXY=off`，依赖来自现有本地缓存，未新增依赖或下载文件。

```text
go test -mod=readonly -buildvcs=false ./common-pkg/amqp/rabbitMq -run '^$'
go test -race -mod=readonly -buildvcs=false -overlay <临时配置> ./common-pkg/amqp/rabbitMq -run '^TestRabbit15' -count=1 -v -timeout=30s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/amqp/rabbitMq common-pkg/app/appServer.go
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 上述验证通过。broker 测试使用驱动接口的内存模拟端，网络写超时使用 net.Pipe；环境未安装 RabbitMQ，未验证实际 broker 的持久化、节点重启和生产积压。TG 使用模拟响应，真实 BotToken、ChatID 和网络环境待服务接入时配置联调。项目原有测试本轮只做全量编译。
- **测试清理** 清理前已登记上述测试命令、结果和限制；随后删除本次临时测试、overlay 和独立构建缓存。核对临时目录及仓库内的 overlay 目标文件均不存在，原有五个测试文件全部保留。

### 2026-11-19 / DEV-20261119-012 / 合并 Request ID 与 Trace ID

- **用户确认** Request ID 与 Trace ID 共用一个值，并把生成规则明确备注。当前任务实现统一 ID；RabbitMQ 处理超时、业务消息接口及其余优化仍按先讨论后开发的流程单独确认。
- **生成与校验** 新增 `headInfo.NewTraceID` 和 `ExtractTrace`。使用 crypto/rand 读取 16 个随机字节，编码为非全零的 32 位小写十六进制字符串，无时间戳、业务前缀或连字符。全零重新生成，随机源失败返回错误。入站拒绝空值、重复值、大写、非法长度、空格及字段冲突。代码中文注释及 [统一请求 ID 规则](headInfo/TRACE_ID.md) 记录生成算法、前端 Web Crypto 示例、网关职责及各层接入方式。
- **入口规则** HTTP 优先沿用合法输入，各 ID 与 traceparent 同时存在时必须相等，全部缺失时兜底生成；非法输入返回真实 HTTP 400，不进入业务。RPC 缺少 ID 仍直接返回 InvalidArgument。trace、request_id、x-request-id 和正文 trace 使用同一个值，Span ID 继续区分操作。HTTP 返回 X-Request-ID，CORS 允许前端读取该响应头。
- **原生链路顺序** 新增 `app.NewServer`、`app.WithRequestTraceRouter`，在 go-zero 路由外层校验并初始化 ID，再创建原生 Trace Span 和日志，避免只改业务上下文造成导出链路不一致。仅有 ID 的新链路按服务的 Telemetry.Sampler 决定采样，有 W3C 上游上下文时保留其采样标记。新增 `RpcAuthMiddleware.ServerOption`，通过 server.AddOptions(auth.ServerOption(...)) 将入站校验放在原生 gRPC 链路拦截器前；同一 Server 只设置一次主 UnaryInterceptor。
- **RPC 与消息传递** HTTP Head、RPC 入站/出站 metadata 和标准 W3C 头同步同一 ID，已有活动 Span 保留，复制后修改 metadata。RabbitMQ 自动复制及注入链路头，消费恢复上下文、创建消费 Span，重试、归档与 TG 告警沿用该 ID；消费者不继承生产请求的超时或取消。独立发布或外部无链路消息可生成新链路，非法消息链路不执行业务并归档。具体消息的 MessageId 保持独立。
- **接口影响** SetTrace 现在返回 (context.Context, error)，会校验 ID 并拒绝不同活动链路。API 使用 app.NewServer 或包装后的 Router；RPC 认证使用 ServerOption 在原生链路前接入。未提供旧接口兼容层。仓库内 SetTrace 无已有调用方，调用示例及 API/RPC 子项目记录已更新。
- **验证** 临时 overlay 下 10 组 TestUnified 用例以 -race 通过。覆盖 64 个并发 ID 生成、全零重生成、随机源失败、严格格式/重复/冲突校验、metadata 副本及活动 Span 保留、真实 go-zero HTTP 路由响应与错误入口、导出的原生 HTTP Span、API 到两级实际 gRPC 服务、无 traceparent 的原始 RPC 调用在原生 Span 前归一化、缺少 RPC ID 拒绝，以及 0/1 采样与上游采样标记保留。首轮临时用例因直接比较包含 TraceState 的结构而编译失败，改为比较 Trace ID/Span ID 后通过。
- **消息验证** 模拟 broker 接口与 TG HTTP Transport 验证发布、消费、延迟重投、最终归档、通知和消费内 RPC 的 ID 一致，MessageId 不变、源头不被修改、Span ID 区分操作、取消生产请求后消费者继续工作，以及无链路新任务生成和非法链路归档。TG 未向真实服务发送。提取文档中的 Web Crypto 示例，用 Node 执行 1000 次格式、非零及重复值检查，通过。
- **执行命令** 以下均使用本次独立临时 GOCACHE、GOPROXY=off 与现有本地依赖缓存；未增加或下载依赖。

```text
go test -race -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestUnified' -count=1 -v ./common-pkg/headInfo ./common-pkg/app ./common-pkg/amqp/rabbitMq -timeout=60s
node <临时目录>/frontend-id.cjs
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l <本次修改的 Go 文件>
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 定向验证、竞态检测、全量构建和项目原有测试编译通过。gRPC 使用 bufconn 上的真实协议，MQ 使用内存接口模拟端；真实网关配置、真实 RabbitMQ/TG 和生产部署联调待接入。前端示例使用 Node Web Crypto 验证，未接入实际前端应用。
- **测试清理** 清理前已登记上述命令、结果和限制；随后删除本次三份临时 Go 用例、Web Crypto 验证脚本、overlay 和独立构建缓存。核对任务临时目录及三个仓库内 overlay 目标文件均不存在，项目原有五个测试文件全部保留。

### 2026-11-19 / DEV-20261119-013 / 删除旧请求 ID 字段及接入入口

- **用户要求** 旧实现直接删除，不保留兼容。清理范围为本次统一请求 ID 留下的重复字段、别名函数及旧接入路径。
- **删除项** 删除 Head.RequestID、consts.RequestId、consts.Trace、GetTrance，以及 RPC 认证中间件的公开 Handle。RPC 认证只通过 ServerOption 接入；应用中间件和 ApiHeaderMiddleware 不再重复初始化 ID，入口统一由 app.NewServer 或包装后的 Router 负责。
- **传输与日志** HTTP 只以 X-Request-ID 对外接收和返回统一 Trace ID，正文仍使用 trace。RPC 和 RabbitMQ 只注入 W3C traceparent/tracestate，不再注入 request_id、trace 或内部 x-request-id 别名。RPC 缺少 traceparent 时直接拒绝，其他 ID 字段和已有 Span 均不能替代。HTTP 业务 Head 只保存 Trace，移除日志名单中的重复 ID 字段，日志 Trace ID 由 go-zero 从当前 Span 输出。
- **实现** 标准链路提取 ExtractTrace 与 HTTP 入口提取 ExtractRequestTrace 分开；SetTrace 直接校验并绑定 ID。EnsureTrace 只读取出站 W3C 链路，保留 metadata 副本隔离，清除不属于当前 Span 的旧 tracestate 后重新注入。新请求的 16 字节安全随机数、32 位小写十六进制、非全零生成规则保持一致。
- **文档** 更新统一 ID 规范和全局 AGENTS.md，移除规范中的旧字段及旧接入示例。公共包原 RPC 示例改为 ServerOption，全局及 API/RPC 子项目登记当前协议，历史验证记录继续作为开发过程留档。
- **定向验证** 临时 overlay 执行 6 组 TestTraceClean 用例，以 -race 通过。覆盖删除字段和公开入口、拒绝旧别名作为 RPC/HTTP 链路替代、HTTP 格式及冲突校验、metadata 仅注入 W3C、业务字段保留、日志名单去重、HTTP 原生 Span 和响应 ID 一致、两级真实 gRPC 调用保持同一 Trace ID 及不同 Span ID、实际 RPC 拒绝旧字段和重复/非法 traceparent，以及 RabbitMQ 发布、消费、重试、归档与消费内 RPC 的标准传播。消息 ID 和单条 ACK 行为保持，取消生产上下文不影响消费。
- **静态核查** Python 检查 Go 源码中无 GetTrance、旧常量引用、Head.RequestID 或 RPC 公开 Handle；链路生成及传输代码不包含旧别名写入，业务中间件的重复初始化已删除。gofmt 和 Git 差异检查通过。
- **执行命令** 以下使用独立临时 GOCACHE、GOPROXY=off 和现有本地依赖缓存。

```text
go test -race -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestTraceClean' -count=1 -v ./common-pkg/app ./common-pkg/amqp/rabbitMq -timeout=45s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/headInfo common-pkg/app common-pkg/appMiddleware common-pkg/rpc/interceptors common-pkg/amqp/rabbitMq common-pkg/consts/system.go
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 定向测试、竞态检测、全量构建及原有测试编译通过。gRPC 使用 bufconn 实际协议，MQ 为接口模拟端；实际网关、RabbitMQ/TG 及生产服务尚未联调。本次不新增依赖，不修改其他业务模块的接口。
- **测试清理** 清理前已登记上述命令、结果和限制，随后删除两份临时用例、overlay 与独立构建缓存。核对任务临时目录和两个仓库内 overlay 目标文件均不存在，原有五个测试文件全部保留。

### 2026-11-19 / DEV-20261119-014 / RabbitMQ 单条业务处理超时

- **用户确认** 为单条业务处理增加独立超时，默认 30s；超时后复用现有重试及耗尽后的 TG 通知。回调退出后再转存，服务关闭取消不计为业务失败。
- **配置与接口** rabbitmq.Config 新增 HandlerTimeout，零值默认 30s，负值拒绝，支持按业务队列配置。计时覆盖每次业务回调，不包含排队、重试等待、发布确认或 TG HTTP。Handler 签名保持不变，传入的 ctx 带独立截止时间，回调结束后取消并回收。配置及关键行为增加中文注释，接入说明补充 YAML、时限范围和取消边界。
- **处理与转存** 在原 worker 中同步调用业务，仅派生限时子上下文；超时即使回调返回 nil 仍按失败处理。超时错误保留 context.DeadlineExceeded，消费 Span 记录错误，日志包含队列、消息 ID、时限、耗时和同一个 Trace ID。回调未退出时也会记录超时，但不提前转存或创建新的业务 goroutine。正常完成时取消定时通知，避免清理调度滞后误报成功消息。
- **重试及关闭** 转存使用原消费生命周期上下文，并由 ConfirmTimeout 独立限时，避免回调的过期上下文阻断重试或归档。等待 broker 确认后才单条 ACK 原消息，转存失败不 ACK。处理前及回调返回后均检查生命周期取消；退出或连接失效优先于超时及 nil 结果，不确认取消后的任务、不增加业务重试次数，关闭连接后由 broker 恢复投递。默认手动确认提供该语义，AutoAck 仍为尽力转存，无法恢复已自动确认的源消息。
- **定向验证** 临时 overlay 下 11 组 TestHandlerTimeout 用例以 -race 通过。覆盖 Go/YAML 缺省与显式 0 的 30s 默认、自定义时限、负值拒绝，正常处理单条确认和回调上下文回收，panic 和下游短超时保留原行为，处理超时后的独立转存、确认前不 ACK、转存 NACK/确认超时不 ACK，超时后返回 nil 仍失败，下一条消息获得新截止时间，DeadlineExceeded 错误链、Trace ID、Span 状态和超时日志不含正文。
- **生命周期及通知验证** 单 worker 的不响应取消回调超时后仍阻止后续业务启动，退出前不转存，退出后恢复消费；日志在等待退出期间已产生且每次处理只记录一次。验证提前父级 deadline、已取消上下文、超时后再关闭、正常关闭及关闭超时均不误 ACK 或推进业务重试。四次处理耗尽默认三次重试后归档，模拟 TG 发送一次并带超时摘要，AutoAck 不重复确认。bufconn 上实际 gRPC 服务验证处理器传入的截止时间取消下游请求，之后仍能成功转存重试消息。
- **执行命令** 以下使用本次独立临时 GOCACHE、GOPROXY=off 和已有本地依赖缓存，无新增依赖或下载文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestHandlerTimeout' -count=1 -v ./common-pkg/amqp/rabbitMq -timeout=45s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/amqp/rabbitMq/config.go common-pkg/amqp/rabbitMq/consume.go common-pkg/amqp/rabbitMq/conn.go
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 定向验证、竞态检测、全量构建及原有测试编译通过。RabbitMQ 使用驱动接口的内存模拟端，TG 使用模拟 HTTP Transport，未发送实际通知；下游 RPC 使用真实 gRPC 协议的内存连接。环境未安装 RabbitMQ，真实 broker 的恢复、确认超时及部署参数待联调。Go 无法强制终止忽略 ctx 的回调，超时后返回仍可能造成业务重复，调用方须按稳定业务标识幂等处理。逐条处理错误通过日志/Span/队列观察，Consume 的返回值只表示订阅初始化。
- **测试清理** 清理前已登记上述命令、结果和限制，随后删除临时 Go 用例、overlay 与独立构建缓存。核对任务临时目录和仓库内 overlay 目标文件均不存在，项目原有五个测试文件全部保留。

### 2026-11-19 / DEV-20261119-015 / RabbitMQ 业务消息接口及使用说明

- **用户确认** 简化业务接入，新增 Message、JSON 编解码、发布 ID 返回值和批量部分成功结果；旧接口直接替换。另要求独立建立带使用范例的说明文件，并确认业务自增 ID 可以用作消息 ID。
- **消息接口** 发布和消费统一使用 Message，提供 ID、正文、普通嵌套消息头、内容属性、RetryCount、Redelivered 及投递来源。Handler 改为 func(context.Context, Message) error，不暴露协议 ACK/NACK。消费前复制正文、嵌套表、数组和二进制头，业务修改 ID、正文、持久化标记或消息头不会污染重试原文。每次回调继续受 HandlerTimeout 限制，关闭及转存确认语义保留。
- **JSON 与发布** 新增 Message.DecodeJSON、NewJSONMessage、PublishJSON、WithMessageID、WithHeaders 和 WithCorrelationID。Publish(ctx, Message) 返回实际 ID 和 error，分配 ID 后即使发布失败也返回该 ID；JSON 编码失败不会发送，编解码错误保留标准错误链并使用不含正文的公开错误文本。NewJSONMessage 可在连接 broker 前生成稳定 ID 及 JSON 正文，供保存和批量发送。保留必要的高级内容属性，业务不必导入驱动。
- **ID 规则** 默认生成 36 位 UUIDv7，不是自增数字，同一毫秒不保证递增。允许用 strconv.FormatInt 将业务自增 ID 转为字符串后指定，也支持业务类型/ID/事件等组合。已指定 ID 不再生成 UUID，重发及包内消费重试沿用原 ID。消息 ID 与 Trace ID 独立，当前封装不提供自动消息去重，业务负责幂等。ID 和内容短字符串按 255 字节校验，避免驱动静默截断；非法元数据在发送前拒绝。
- **批量结果** PublishBatch(ctx, []Message) 返回 BatchResult 和 error。MessageIDs 对应输入顺序，已尝试项返回实际 ID，未尝试项保留输入 ID；ConfirmedCount 为已确认且已路由的连续前缀；FailedIndex 从 0 开始，首次失败即停止，成功或批次初始化失败为 -1。错误链支持 ErrPublishUnknown，失败那条可能已接收，其后未尝试。传入切片不修改，由调用方保留 ID 并只重发未确认部分。
- **使用说明** 新建 [USAGE.md](amqp/rabbitMq/USAGE.md)，含初始化/消费/退出的完整 main、业务自增 ID、字节正文、离线 JSON 消息准备、稳定 ID 重发、批量部分成功恢复、错误分类和幂等职责。README 作为配置与运行入口链接该文件，删除旧 Publishing 发布及 []byte 回调范例，统一链路规范改用 Message.ID 字段名称。
- **定向验证** 临时 overlay 执行 13 组 TestMessage 用例，-race 通过。覆盖结构体 JSON 往返和大 int64 精度、中文、空目标/语法/类型错误链、编码失败不发送且返回已分配 ID，显式业务数字/组合/大整数 ID 原样发送、自动 ID 与 24 个并发发布、常用属性与默认值、嵌套头和正文隔离、循环头和不支持类型拒绝、过长字段/非法 TTL 校验，关闭/未就绪/未路由/NACK/结果不确定时的 ID 保留。
- **流程验证** 覆盖批量部分成功、失败停止、已分配 ID 的尾部恢复及输入切片不变；消费元数据、JSON 解码、业务修改不影响原文及 Trace ID，单条确认时序、解码失败归档、处理超时后独立转存、关闭取消不确认。检查旧的驱动 Publishing 发布签名和 []byte Handler 已删除。4 段完整 Go 使用范例均抽取到临时文件编译通过，未执行会连接真实 broker 的示例。
- **执行命令** 以下使用本次独立临时 GOCACHE、GOPROXY=off 和现有本地依赖缓存，没有新增依赖或下载文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay <临时配置> -run '^TestMessage' -count=1 -v ./common-pkg/amqp/rabbitMq -timeout=45s
go build -mod=readonly -buildvcs=false -o <临时产物1> <临时目录>/usage_1.go
go build -mod=readonly -buildvcs=false -o <临时产物2> <临时目录>/usage_2.go
go build -mod=readonly -buildvcs=false -o <临时产物3> <临时目录>/usage_3.go
go build -mod=readonly -buildvcs=false -o <临时产物4> <临时目录>/usage_4.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/amqp/rabbitMq/message.go common-pkg/amqp/rabbitMq/provider.go common-pkg/amqp/rabbitMq/conn.go common-pkg/amqp/rabbitMq/consume.go
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 定向测试、竞态检测、使用范例编译、全量构建和原有测试编译通过。broker 使用驱动接口模拟端，无实际 RabbitMQ/TG 联调或消息发送。JSON 解码错误仍按现有错误策略重试，不提前实现永久错误分类；消息去重、业务事务与失败任务持久化由业务接入。使用范例的实际运行需提供真实连接及服务配置。
- **测试清理** 清理前已登记上述命令、结果和限制；随后删除临时测试、overlay、4 份抽取范例源码、全部范例编译产物和独立构建缓存。核对当前文档与已编译范例一致、任务临时目录和仓库内 overlay 目标文件均不存在，原有五个测试文件全部保留。


### 2026-11-19 / DEV-20261119-016 / RabbitMQ 消费错误分类和使用范例

- **用户确认** 开发消费错误分类，并单独建立带使用范例的说明文件。本次替换此前“所有 JSON 解码错误均重试”的行为，普通业务错误仍默认重试。
- **公开接口** 新增 Permanent(err) 和 Retryable(err)，nil 输入均返回 nil，保留错误文本和错误链，支持 errors.Is/As 及 %w。先查找显式标记，外层覆盖内层；多分支按标准 errors.As 深度优先、从左到右找到第一个标记。无标记时，标准 JSON SyntaxError、UnmarshalTypeError、InvalidUnmarshalError 直接归档，其他自定义解码器错误与普通业务错误默认重试。Retryable 可覆盖默认 JSON 分类或内层 Permanent，仍受 MaxRetries 限制，不根据业务文案或 HTTP/gRPC 状态码猜测策略。
- **消费与确认** 永久错误不进入业务重试，直接转存失败交换机；确认转存后才单条 ACK，NACK 或确认超时不 ACK 源消息。服务生命周期取消优先，不重试或确认；HandlerTimeout 仍覆盖超时后返回的 nil 或永久标记，按超时重试。归档补齐缺失消息 ID，原有业务 ID、正文及合法 Trace ID 保留。
- **终态通知** 新增内部 x-base-failure-reason，分别保存 permanent、retries_exhausted、invalid_metadata；生产者的同名头不控制策略，业务重试移除旧终态原因。TG 按原因显示不可重试错误、重试耗尽或元数据错误，永久错误归档后即可通知，普通错误仅耗尽后通知。TG 仍默认关闭、独立重试并在发送成功后 ACK 通知任务。处理次数为本包 RetryCount+1，不统计 broker 重投的额外执行；元数据错误注明本次未执行，未知原因使用通用文案。
- **文档与注释** 新建 [消费错误处理说明](amqp/rabbitMq/ERROR_HANDLING.md)，包含完整领域处理器与注册、显式重试、包装和 nil 的范例，说明优先级、超时、归档和 TG 行为。同步 README、USAGE 的入口与字段校验例子，以及 Handler、DecodeJSON、TelegramConfig 和关键分支的中文注释。平台与 API/RPC 开发记录同步接入约定。
- **定向验证** 临时 overlay 的 10 组 TestProcessingPolicy 用例及竞态检测通过。覆盖 nil、errors.Is/As、%w、嵌套与多分支策略，标准 JSON 和自定义解码错误；永久错误直接归档、等待发布确认才 ACK、原文/ID/链路不变；显式重试、连续失败耗尽和零次重试；NACK/确认超时不 ACK；元数据归档原因和缺失 ID；超时与关闭优先级；三类 TG 文案、HTTP 失败后独立重试及成功后确认；普通成功与 AutoAck 不重复协议确认。
- **执行命令** 在平台根目录执行，使用临时 GOCACHE、GOPROXY=off 和现有依赖缓存，无新增依赖。overlay、测试输出、抽取的范例和编译产物均位于 /private/tmp/new-platform-rabbit-policy.Agu6KF0T。文档通过临时 Python 抽取脚本分别编译 USAGE 的 4 段和 ERROR_HANDLING 的 3 段 Go 范例，未运行需要真实 broker 的程序。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-rabbit-policy.Agu6KF0T/overlay.json -run '^TestProcessingPolicy' -count=1 -v ./common-pkg/amqp/rabbitMq -timeout=45s
go build -mod=readonly -buildvcs=false -o <临时目录>/example_N <临时目录>/example_N.go  （N=1..7，逐一执行）
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/amqp/rabbitMq/{errors,message,consume,telegram,config,conn}.go
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 10 组定向测试及 -race、7 段范例编译、全量构建、原有测试编译和格式检查均通过，RabbitMQ 文档相对链接检查通过。broker 使用驱动接口模拟端、TG 使用模拟 HTTP Transport，未连接真实 RabbitMQ 或发送实际 TG 消息；真实服务、凭证、断线恢复及部署环境仍需联调。未执行原有集成测试的业务逻辑。错误分类不提供业务去重，业务须幂等。
- **测试清理** 清理前已登记实际命令、结果和限制；随后删除本次临时测试、overlay、测试输出、7 份范例源码和编译产物及独立构建缓存。已核对任务临时目录、早先合并缓存的临时目录及仓库内 overlay 目标文件均不存在，原有五个测试文件全部保留。API/RPC 本次仅更新文档，无新增测试或临时产物。


### 2026-11-19 / DEV-20261119-017 / RabbitMQ 失败归档与通知双份确认

- **用户确认** 删除失败广播投递代码，改为先定向保存失败队列，再定向创建通知任务，两份都确认后才单条 ACK 源消息；说明部分成功后重投可能重复。
- **实现** archive 使用默认交换机和 FailedQueue/NotifyQueue 的队列名分别发布，每次均沿用 mandatory 和发布确认，各受独立 ConfirmTimeout 限制。第一份未确认不发送第二份，任一步失败或结果不确定不 ACK；两份确认后再检查生命周期取消。错误包装区分保存失败消息和创建通知任务，保留 ErrUnroutable、ErrPublishRejected、ErrPublishUnknown 与 context 错误链。
- **拓扑清理** 删除内部 failedExchange 字段、生成及冲突校验、失败 fanout 声明及两条绑定；仍声明业务、重试、失败和通知四个队列，具名业务交换机和绑定按业务配置保留。无旧投递兼容路径，不自动删除 broker 历史资源，无公开接口或新配置字段。
- **语义与说明** 两份副本使用相同消息 ID、正文、失败原因和原有链路；通知关闭时仍可靠保存通知任务。发布确认丢失、第一份成功但第二份失败、源消息 ACK 失败后重投均可能重复归档或通知，业务与归档处理须幂等。外部未指定 ID 的源消息重投可能重新生成 ID，生产者应提供稳定 ID。同步 [运行说明](amqp/rabbitMq/README.md#失败归档与通知的投递确认)、[使用说明](amqp/rabbitMq/USAGE.md) 和 [错误处理说明](amqp/rabbitMq/ERROR_HANDLING.md)，补充中文关键注释。
- **定向验证** 临时 overlay 的 10 组 TestDualArchive 测试及 -race 通过。覆盖默认/具名业务拓扑与辅助队列声明失败、删除旧保留名称约束；先归档再通知及双确认前不 ACK；两个目标分别 basic.return、NACK、确认超时、写入错误；归档前/第一份/第二份取消、取消与第二份确认同时到达、通知投递时 Close；永久错误/耗尽/JSON/业务及重试元数据错误统一归档；原文、ID、Trace 和原因一致；成功及普通重试路径；部分成功后重复处理保持已有 ID；AutoAck 不重复协议确认及源消息 ACK 失败失效连接。
- **执行命令** 在平台根目录使用任务独立 GOCACHE=/private/tmp/new-platform-rabbit-dual.B6wkAXD9/cache 和 GOPROXY=off，复用本地依赖，无下载或新增依赖。测试源码、overlay、测试输出和缓存位于同一临时目录，没有在仓库创建测试文件。

```text
gofmt -w common-pkg/amqp/rabbitMq/config.go common-pkg/amqp/rabbitMq/conn.go common-pkg/amqp/rabbitMq/consume.go
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-rabbit-dual.B6wkAXD9/overlay.json -run '^TestDualArchive' -count=1 -v ./common-pkg/amqp/rabbitMq -timeout=45s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/amqp/rabbitMq/{config,conn,consume}.go
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 10 组定向测试及竞态检测、全量构建、原有测试编译通过；格式、diff 和说明相对链接/新增锚点检查通过，原有五个测试保留。驱动接口模拟发布、确认与路由失败，无真实 RabbitMQ/TG 联调或外部消息发送；真实队列删除、重连和部署参数仍需联调。未运行原有集成测试业务逻辑。
- **测试清理** 清理前已登记全部实际命令、结果和限制，随后删除本次临时测试、overlay、测试输出和独立构建缓存。核对任务临时目录及仓库内 overlay 目标均不存在，原有五个测试全部保留。API/RPC 本次只更新文档，无新增测试或临时产物。


### 2026-11-19 / DEV-20261119-018 / Elasticsearch 7 独立配置及请求生命周期

- **用户确认** Elasticsearch 7 与 OpenSearch 分为两套独立实现，配置、初始化、错误、连接池和文档各自维护。本轮仅开发 es/esv7，完成后等用户确认再处理 OpenSearch；XXL-JOB 尚未开始。
- **独立配置** 新增 core.Config、TLSConfig 和 TransportConfig，全部字段补齐中文用途、默认值和单位。Addresses 必填，HTTP/HTTPS 可带路径前缀，认证使用独立字段；默认系统证书验证及 TLS 1.2，自定义 CAFile 追加信任，无法读取或无有效证书初始化失败，InsecureSkipVerify 默认 false。StartupTimeout 默认 5s、RequestTimeout 默认 10s，连接/TLS/响应头阶段默认 5s，连接池三个上限默认 100，空闲保留 90s；零值取默认、负值拒绝，校验每节点空闲数不超过连接总数。配置不修改调用方对象。
- **重试与认证** MaxRetries 改为 *int，nil 默认 3，显式 0 设置 SDK DisableRetry；针对当前 v7.17.10 SDK 直接 io.EOF 仍重试的特例，在关闭重试时包装该错误，保留 errors.Is(io.EOF)，避免绕过配置。SDK 自带认证替代重复手工 Basic Auth 头，单次只发送一条认证头；网络/502/503/504 重试沿用 SDK 规则，操作总预算不随重试重置。
- **初始化与关闭** NewEs(ctx, conf)/MustNewEs(ctx, conf) 接收启动上下文并执行 Info 探活，使用 StartupTimeout 与上游较早的截止时间，读取完成再取消，失败释放空闲连接。Close 可并发和重复调用，拒绝后续请求并释放连接池空闲连接，不强制取消在途 I/O。OperationContext 派生模型请求预算，调用方完成响应处理后再 cancel。移除旧签名、手工认证函数和未启用的旧 tracing transport/孤立指标定义，无兼容包装。
- **响应与错误** ReadResponse、CheckResponse 和 Get/Search 解析均接管响应体，在成功、404、错误状态、读取/关闭/解码失败及取消时释放；错误响应最多读取 4KiB，不无限 drain。独立 ErrNotFound/ErrInvalidResponse/ErrClosed 和 ResponseError 支持 errors.Is/As，HTTP 状态与最多 128 字节 error.type 用于诊断，不输出原始正文或 error.reason；JSON/网络/上下文错误保留底层错误链。Get 的 404/found=false 返回 ErrNotFound，成功直接解码到业务指针并设置 Schema ID，修复原 &interface{} 目标传递。
- **模型操作** 创建索引、create 写入、单查、删除、局部更新和搜索均继承 model.Ctx 并应用 RequestTimeout，不再用 Background 替换。增加 nil 模型/上下文/客户端、typed nil、空索引或 ID 校验。局部更新补齐 Elasticsearch 必需的 doc 包装。Search 删除已关闭的原始响应返回值，只返回 (int64, error) 并写入目标切片；使用原始 JSON 片段保留大 int64 精度，支持对象/数字 total 和空命中，track_total_hits=true 获取准确总数，统计成本在文档注明。
- **使用说明** 新建 [Elasticsearch 7 使用说明](es/esv7/README.md)，包含独立 YAML 默认值、完整初始化/FindOne/Search 示例、低层 SDK 请求预算与响应释放、错误语义和签名变更。原有 esv7/core/es_test.go 仅更新注释中的旧签名和指针配置，原测试继续保留；未修改 OpenSearch 或 XXL-JOB 源码。
- **定向验证** 最终 20 组 TestES7 用例及 -race 通过，其中 core 13 组、model 7 组。真实本机 HTTP/TLS 验证默认不信任证书拒绝、自定义 CA 信任及显式测试跳过、无效 CA、Basic Auth 单头、连接池覆盖、初始化上下文取消/超时和阶段超时；验证 SDK 重试 0/2 次以及关闭重试时直接 EOF 仅一次、开启时重试。
- **模型与清理验证** 覆盖所有六类操作的成功/上下文超时/取消、404 映射及 doc 更新封装、父级在途取消、正文读取时限、多个重试共享预算、Close 并发与后续请求拒绝；响应体正常/失败/解码/关闭错误链、4KiB 读取及 UTF-8 摘要边界、nil 响应/目标及必要字段校验、文档 ID 写回、大 int64、空搜索结果及不同 total 格式。README 的两段完整 Go 范例抽取编译通过，未运行需要真实服务的文档程序。
- **实际执行与排障** 初次定向命令被沙箱禁止本地监听端口拦住，获得仅本地 HTTP/TLS 测试所需权限后继续。首次本地运行 core 通过，model 临时服务未读取请求体导致测试服务清理等待并触发 60s 测试超时；修正临时服务读取请求体后通过。随后补充 EOF 边界与定向验证，最终完整 20 组 -race 回归通过。未把测试服务等待问题当作产品通过结果。
- **执行命令** 平台根目录执行，使用独立 GOCACHE=/private/tmp/new-platform-es7.hk7ZJfwS/cache、GOPROXY=off，定向命令额外设置 ES7_TEST_ARTIFACTS=/private/tmp/new-platform-es7.hk7ZJfwS/artifacts。overlay 测试、本地证书、测试输出、抽取范例与缓存均在此临时目录，未下载依赖。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/es/esv7/...
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-es7.hk7ZJfwS/overlay.json -run '^TestES7' -count=1 -v ./common-pkg/es/esv7/... -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-es7.hk7ZJfwS/overlay.json -run '^TestES7NoRetryEOF$' -count=1 -v ./common-pkg/es/esv7/core -timeout=30s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-es7.hk7ZJfwS/overlay.json -run '^TestES7' -count=1 -v ./common-pkg/es/esv7/... -timeout=45s
go build -mod=readonly -buildvcs=false -o <临时目录>/example_1 <临时目录>/example_1.go
go build -mod=readonly -buildvcs=false -o <临时目录>/example_2 <临时目录>/example_2.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/es/esv7/core/*.go common-pkg/es/esv7/model/*.go
git -C common-pkg diff --check
git diff --check
git -C common-pkg diff --exit-code -- es/esaws xxlJob
```

- **结果与限制** 最终 20 组定向测试、竞态检测、两段范例编译、全量构建和原有测试编译通过，格式和 diff 检查通过，OpenSearch/XXL-JOB 无改动。真实本机 HTTP/TLS 服务验证协议和取消，响应体关闭计数、读取故障及 EOF 使用测试替身；未连接真实 Elasticsearch 集群，也未执行原有外部服务集成测试。直接用底层 SDK 须自行提供上下文/预算并释放响应，JSON 同步计算不可强制中断，准确总数查询有统计成本。写入结果不确定和业务幂等由接入方处理。
- **测试清理** 清理前已登记实际命令、结果、初次失败原因和限制；随后删除本次两份 overlay 测试、overlay 配置、全部测试输出、CA 证书、两份范例及编译产物和独立构建缓存。核对任务临时目录与仓库内两处 overlay 目标均不存在，原有五个测试全部保留；原 ES7 测试只更新注释范例。API/RPC 本轮仅更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-001 / OpenSearch 独立配置、原生响应与聚合

- **用户确认** 接续 Elasticsearch 7，本轮独立开发 es/esaws，自行维护配置、TLS、连接池、错误和文档，不建立共用类型或兼容入口；保留 BusinessCode 后缀、聚合及 IndexSchema 全量覆盖。XXL-JOB 尚未处理。
- **独立类型与配置** 根包声明改为 esaws，自行定义 IndexTable/Schema；新增独立 core.Config/TLSConfig/TransportConfig，不引用 es/esv7。缺省系统证书验证及 TLS 1.2，自定义 CAFile 追加信任；StartupTimeout 5s、模型 RequestTimeout 10s，TCP/TLS/响应头阶段 5s，连接池三项上限 100、空闲保留 90s，省略或零值采用默认，负值及非法地址/连接上限拒绝。MaxRetries 改为 *int，nil 默认 3、显式 0 关闭，配置不修改调用方对象。
- **初始化与 SDK 边界** NewEs(ctx, conf)/MustNewEs(ctx, conf) 执行受控 Info 探活，失败释放空闲连接；删除 InitESClient 和旧 MustNewEs 签名。当前 opensearch-go v1.1.0 内部产品探活使用 Background，关闭该隐式探活，改由本包检查连接、认证、HTTP 状态和 JSON 正文。删除原不可取消的随机指数退避，不设置使用 time.Sleep 的 SDK RetryBackoff；重试仍按默认网络错误及 502/503/504 触发并共用操作额度，不额外开启超时重试。关闭重试时包装直接 io.EOF，防止 SDK 特例绕过 DisableRetry，并保留 errors.Is 链。
- **路径和生命周期** SDK v1 在重试时重复追加节点路径前缀；通过本套 Perform 记录原始路径和 RawPath，网络请求使用副本，恢复 SDK 原始路径以供下一次重试/切换节点，保留文档 ID 转义及调用方 URL。生成 API 绑定到当前实例，底层原生 Transport 保留 SDK Metrics/DiscoverNodes 能力。Close 可并发或重复调用，释放空闲连接并拒绝新请求，在途请求由自身 ctx 取消；OperationContext 提供单次模型总预算，父级更早截止时间优先。
- **原生响应与错误** 删除 Elasticsearch v7/v8 响应、聚合类型及强制转换，统一使用 opensearchapi.Response。ReadResponse/CheckResponse/Get/Search 在正常、404、错误状态、取消、读取/解码/关闭失败时释放正文，不先 drain 再取错误，也不先关闭再读取；错误正文最多读取 4KiB，错误类型摘要最多 128 字节，不直接输出原始正文或 error.reason。独立 ErrNotFound/ErrInvalidResponse/ErrClosed 和 ResponseError 支持 errors.Is/As，删除旧 EsNotFound/Es500Err；Get 的 404/found=false 返回 ErrNotFound，直接解码到业务指针并设置 Schema ID。
- **模型与聚合** 七类操作继承 model.Ctx，使用单次 RequestTimeout，删除业务调用中的 Background。BusinessCode 非空时模型文档与搜索操作使用基础名加下划线及业务码；IndexAuto 继续接受完整索引名，不自动加后缀。保留 create 自动 ID/稳定 ID、doc 局部更新和 IndexSchema 非空 ID 全量覆盖及 refresh=true。Search 返回 (int64, core.Aggregations, error)，不返回已关闭响应；Aggregations 为本套定义的 map[string]json.RawMessage，未知/嵌套聚合和服务端原始数值保留，空聚合返回空 map，空命中为空切片，固定 track_total_hits=true 并注明统计成本。
- **注释、文档与现有调用** 新增 [OpenSearch 使用说明](es/esaws/README.md)，包含独立 YAML 默认值、初始化/业务后缀/查询与聚合解码、全量覆盖及底层 SDK 超时释放范例。删除未启用旧 transport、手工重复认证头函数及孤立指标；认证由 SDK 自行处理。原核心测试仅更新注释范例，原 es/test 集成测试更新 ctx、指针重试配置、Close 和 Search 范例，用 errors.As/StatusCode 替换旧错误文案匹配，保留原测试。API/RPC 仅登记接入约定。
- **定向验证** 最终 23 组 TestOS 测试及 -race 通过，core 14 组、model 9 组。真实本机 HTTP/TLS 覆盖单次受控探活、Basic Auth 单头、连接池覆盖、默认不信任证书拒绝/自定义 CA/测试跳过/无效 CA、启动响应头和正文超时、父级取消、阶段超时；验证重试 0/2 次、直接 EOF 的关闭/开启重试、节点前缀不重复和文档 ID 转义、调用方 URL 不变、Close 并发与后续拒绝。
- **模型与解析验证** 覆盖七类操作成功及取消/超时、BusinessCode 有无后缀和显式索引名、404 错误映射、doc 更新与 refresh 覆盖写入、稳定 ID 重试时正文复用及自动 ID 行为、在途父级取消/正文读取/多次重试共用预算；正常/错误/读取及关闭失败的正文释放、摘要长度与 UTF-8 边界、nil 及必要字段校验、大 int64、total 两种格式、空命中/空聚合、嵌套桶和未知聚合原样解码。go list 生产依赖树无 Elasticsearch 或 es/esv7 引用，按任务开始时的 SHA-256 快照验证 ES7 和 XXL-JOB 文件未被本轮修改。
- **执行与结果** 使用 GOCACHE=/private/tmp/new-platform-opensearch.djrNAnWE/cache、GOPROXY=off，定向测试设置 OS_TEST_ARTIFACTS=/private/tmp/new-platform-opensearch.djrNAnWE/artifacts；本地监听测试经权限审核运行，不连接外部服务。首次 core 14 组通过，model 两个用例错误假设 SDK create 路径为 /index/_create/id；按实际 v1 原生 /index/_doc/id/_create 修正临时测试预期后通过，未为迎合测试改变业务请求。最后调整 API 绑定保留原生 Transport 能力后再次完整回归，最终 23 组通过。三段 README Go 范例抽取编译通过；全量构建、原有测试编译、gofmt 和 diff 检查通过，无下载或新依赖。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/es/esaws/...
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-opensearch.djrNAnWE/overlay.json -run '^TestOS' -count=1 -v ./common-pkg/es/esaws/... -timeout=45s
go list -deps -mod=readonly -buildvcs=false ./common-pkg/es/esaws/...
go build -mod=readonly -buildvcs=false -o <临时目录>/example_N <临时目录>/example_N.go  （N=1..3，逐一执行）
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/es/esaws/core/*.go common-pkg/es/esaws/model/*.go common-pkg/es/esaws/esInterface.go common-pkg/es/test/es_test.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 使用真实本地 HTTP/TLS 和当前 SDK 验证协议行为，故障正文和 EOF 使用测试替身；未连接真实 OpenSearch 集群或执行原有集成测试业务逻辑。初始化探活不替代实际集群/索引权限及版本联调；SDK 高级调用须自行提供上下文预算并释放正文，JSON 同步计算不可强制中断，准确总数和 refresh=true 有额外成本，写入结果不确定与幂等由接入方处理。当前 SDK 节点发现等高级功能保留，未新增其真实集群联调。
- **测试清理** 清理前已登记实际命令、结果、首次测试预期错误和限制；随后删除本次两份临时 overlay 测试、配置、全部测试输出、CA 文件、依赖清单、快照、三份抽取范例及编译产物和独立构建缓存。已核对任务临时目录与两处仓库 overlay 目标均不存在，原有五个测试全部保留，原 OpenSearch 测试仅按新接口同步；按删除前快照对照确认本轮未改 ES7/XXL-JOB 文件。API/RPC 仅更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-002 / XXL-JOB 监听、心跳及任务关闭生命周期

- **用户确认** 接管当前 SDK 无法完整关闭的 HTTP 服务、注册心跳和任务上下文，增加 Start(ctx)/Stop(ctx) 错误返回、明确启动结果、关闭预算及任务注册校验，更新范例和说明。用户随后明确回复“沿用现有协议”，确认现有 /run 与数组回调协议，无需修改当前实现。协议结构以官方 2.4.1 为验证基准，未将该基准视为用户指定的部署版本，实际版本待联调。
- **配置与 API** NewClient(Config) 返回 (*Client,error)，不监听或注册；配置 AdminAddress/AccessToken/RegistryKey/ExecutorAddress 必填，ListenAddress 缺省 :9999。RequestTimeout 默认 5s、RegistryInterval 20s、ShutdownTimeout 10s，HTTP 请求体上限 MaxBodyBytes 默认 1MiB，已接受未完成任务总数 MaxTasks 默认 1024；零值采用默认、负值及非法地址/端口拒绝，注册地址不得为未指定 IP。新增和修改的配置、公开接口和关键路径补齐中文注释。
- **启动与关闭** Start 使用服务生命周期 ctx，先监听再完成首次注册，监听或注册失败返回阶段错误并完成有界清理；异步 Serve 错误通过 Wait 返回，不退出进程。Stop 立即拒绝新任务、关闭监听并取消任务，等首次注册和心跳本地请求退出后再摘除，防止后台再次周期注册；ShutdownTimeout 覆盖 HTTP 收尾、摘除、任务和回调等待，超限强制关闭 HTTP/取消回调后返回超时。未启动、重复或并发 Stop 均安全。调用方 Stop/Wait 的 ctx 只中止等待，不放弃后台有界清理；父 Start ctx 取消自动触发关闭。一实例只能启动一次，关闭后须重新创建。
- **任务注册与调度** Register 改为 error 返回且批量原子提交，拒绝 nil/typed nil、空名、重复名、nil 函数和启动后追加；可选 SetLogHandler 仅启动前设置。保留现有 HandlerFun 的 string 结果签名，普通返回按成功，超时、取消、panic 按失败。每次调度独立快照避免 SDK 复用任务对象的数据竞争；同 JobID 串行 FIFO、DISCARD_LATER 拒绝忙碌任务、COVER_EARLY 取消旧任务及排队任务并等待旧处理器退出后执行最新任务，不强制并发启动不能终止的 Go 代码。kill 幂等，未退出任务仍忙碌。受控排队与回调计入统一容量，活跃 LogID 仅对相同请求去重，不提供持久化去重。
- **上下文与回调** 调度入口沿用合法 HTTP Request/Trace ID，缺省按平台规则生成独立 Trace；不继承启动请求的 ID 或 /run 响应后的取消。ExecutorTimeout 按秒从实际处理开始计时，为 0 不加执行时限，所有任务仍响应服务关闭。业务修改参数不会影响原 LogID/LogDateTime 回调；业务退出后回调沿用原 Trace，通过独立网络和关闭预算上报取消失败。回调不阻塞业务串行队列，结果摘要最多 1024 字节并规范无效 UTF-8；本轮保持单次回调，失败记录日志，不增加自动重试或持久化队列。
- **协议与日志** 五个执行器端点均要求单一正确令牌、POST 和有界 JSON 正文，拒绝重复令牌、未知 handler/策略、非法参数及非 BEAN 脚本任务。管理请求使用独立连接池、有界正文读取和 RequestTimeout，拒绝重定向，不打印令牌、完整参数或调度中心原始正文。/log 支持带 ctx 的日志提供器，默认明确提示使用应用日志系统，不伪造历史日志。删除 SDK Init/Run/Stop、MustStart、旧构造/链式注册及无参 Stop、旧 Logger 适配器和 demo.Old，依赖仅复用协议 DTO。
- **说明与示例** 新建 [XXL-JOB 使用说明](xxlJob/README.md)，说明配置、API 签名、Start ctx 的整个生命周期含义、Stop 等待取消、任务策略、链路及回调限制，并提供完整程序、主动关闭及日志提供器三段 Go 范例；demo/main.go 改为环境变量、应用退出信号、Start/Wait 和 ctx 业务示例，不运行外部执行器或无限等待。
- **定向验证** 16 组本机 HTTP/协议/生命周期测试及 -race 通过，另 1 组 UTF-8/JSON 边界 -race 通过，共 17 组。覆盖配置缺省/非法值和注册原子性；Stop-before-Start、端口冲突不注册、注册失败释放端口、首次注册时并发 Stop 无死锁；多次/并发 Stop 仅摘除一次、停止后无心跳、监听端口可重新绑定；父 ctx 自动停止、Wait 取消不停止服务；非法/重复令牌、HTTP 方法/路径、超大/null 请求、未知任务/脚本和非法 Trace；任务快照、并行独立 JobID 与 Trace、串行/丢弃/覆盖/重复活跃调度；kill 及 Stop 对排队任务失败回调；实际调度超时及 panic 脱敏；不响应 ctx 的任务关闭超时和晚退出不再发送回调；调用方 Stop 截止后后台继续清理；摘除超时及异步监听错误可观察；日志提供器、总容量、回调失败不阻塞后续业务、1024 字节 UTF-8 和额外 JSON 拒绝。
- **执行记录与排障** 初次包构建缺少 time 导入、临时测试混用有名/无名参数导致编译失败，均修正后验证通过。第一次本地端口测试的自动权限审核未在期限内完成，工具允许重试；先以沙箱内编译检查验证脚本，再按提示重试获准，16 组本地测试通过。最后规范无效 UTF-8 摘要，并以额外无监听边界用例验证。三段 README 范例抽取编译、全量构建和原有测试编译通过，格式、diff 和本地文档链接检查通过；按开始时 SHA-256 快照确认搜索和 AMQP 源码未被本轮修改。
- **实际命令** 在平台根目录使用独立 GOCACHE=/private/tmp/new-platform-xxljob.tQv5N685/cache 和 GOPROXY=off，无下载或新增依赖。临时测试、overlay、输出、源码快照、抽取范例与产物均在同一临时目录。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/xxlJob/...
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-xxljob.tQv5N685/overlay.json -run '^$' ./common-pkg/xxlJob
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-xxljob.tQv5N685/overlay.json -run '^TestXXL' -count=1 -v ./common-pkg/xxlJob -timeout=45s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-xxljob.tQv5N685/overlay.json -run '^TestXXLHelperBoundaries$' -count=1 -v ./common-pkg/xxlJob
go build -mod=readonly -buildvcs=false -o <临时目录>/example_N <临时目录>/example_N.go  （N=1..3，逐一执行）
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/xxlJob/*.go common-pkg/xxlJob/demo/main.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 使用真实本机 HTTP 监听和本地调度中心替身，未连接实际调度中心或发送外部回调；版本和部署地址仍须联调，未运行原有外部服务集成测试。Go 无法强制终止忽略取消的任务/日志提供器，关闭预算耗尽后业务须自行退出；队列和回调是内存状态，进程崩溃或关闭预算耗尽不保证回调送达；正常 string 返回没有业务 error 结果通道，业务须明确现有签名语义并自行幂等。注册响应取消后的远端瞬间状态不保证，摘除失败通过错误返回并依靠调度中心过期判定或管理处理。
- **测试清理** 清理前已登记实际命令、结果、初次编译/审核问题及限制；随后删除本次临时测试、overlay、全部输出、快照、三份范例源码及编译产物和独立构建缓存。核对任务临时目录与仓库 overlay 目标均不存在，原有五个测试全部保留；删除前摘要对照确认搜索和 AMQP 未被本轮修改。API/RPC 只更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-003 / HTTP 请求客户端与有界响应

- **用户确认与范围** 接续第 17 项，按已讨论方案先处理 commonTool HTTP 请求；随机字符串、TOTP URI 和 Redis 分片本轮不修改。原 GET 10000 秒超时与注释不符，POST 无超时，两者逐次创建客户端、缺少请求上下文及响应体上限。删除旧 GetRequest/PostRequest，不保留兼容包装；仓库 Go 调用检索无残留。
- **独立配置** 新增 HTTPConfig、HTTPTransportConfig 和 NewHTTPClient，配置字段、公开接口及关键逻辑补齐中文备注。总预算默认 10s，响应正文默认 4MiB；TCP/TLS/响应头默认 5s，连接池三项上限默认 100，空闲保留 90s。零值采用默认、负值拒绝，响应上限 MaxInt64 拒绝以避免加一溢出，每节点空闲数不得超过连接总数；不修改调用方配置或进程级默认客户端。HTTPS 使用系统证书验证和最低 TLS 1.2，代理沿用标准环境变量。
- **请求与生命周期** Get/PostJSON 必须传入 ctx，保留原值和取消，父级更早截止时间优先；总预算覆盖 JSON 编码耗时及网络读取，编码失败或编码后超时不发送。共享连接池支持并发，Close 可并发重复调用，拒绝后续请求并释放空闲连接，在途请求受自身上下文控制。请求头规范大小写并复制多值切片，支持单一非空 Host 覆盖；POST 默认 application/json。
- **响应与错误** 所有 2xx 成功，非 2xx 返回可 errors.As 识别的 HTTPStatusError，同时保留 HTTPResponse 的 StatusCode/Header/Body/BodyComplete。正文最多读取上限加一字节，超限保留上限内前缀并返回 ErrHTTPResponseTooLarge；自动 gzip 按解压后大小计。状态、读取、关闭、超限和取消错误可共存，保留 errors.Is/As 链，公开错误文本不输出 URL、凭证或正文。正常、异常和取消均释放响应体，调用方无需再 Close 响应。
- **重定向与链路** 不自动跟随重定向，3xx 原样返回；不增加业务重试，POST 关闭已发送正文的自动重建/重放，标准库尚未发送正文前的连接恢复及 GET 原有恢复机制保留，共用调用预算。业务写入须自行幂等，网络错误不证明远端未处理。工具不自行生成或注入 Trace ID；跨 HTTP 传递已有 ID 的显式范例见新建 [HTTP 使用说明](commonTool/HTTP_CLIENT.md)，其中包含配置、GET、POST JSON、错误和服务关闭说明。
- **定向验证** 20 组 TestHTTPTool 及 -race 全部通过：配置缺省/覆盖/非法值及不修改输入，200/201/202/204/206、查询和多值头，大 int64 JSON/Content-Type/Host/正文不可重建，301/304/400/401/429/500/503 与错误脱敏；空正文/恰好上限/超限/容量/状态和超限并存，gzip 解压后限额；读取及关闭故障和响应快照；JSON/nil ctx/未初始化/非法 URL/Host/头校验，编码消耗预算；GET/POST 取消/总超时、父级截止、正文超时及响应头阶段超时；GET/POST 不跟随重定向、带 Idempotency-Key 的 POST 不重放、连接复用与并发 Close、在途请求正常收尾、24 个并行 POST 隔离及不信任 TLS 证书拒绝。
- **执行与排障** 定向测试以 overlay 注入，未在仓库创建测试文件；实际本机端口测试经权限审核运行。首次 19 组通过，非法头用例错误假设自定义 RoundTripper 会执行标准库头校验，调整临时测试为真实本机服务和原生 Transport 后完整 20 组通过，生产代码未为测试替身增加行为。范例复核脚本曾因源码缩进格式不同触发相等性断言，重新直接抽取当前文档源码后逐一编译通过，不涉及产品或范例编译失败。全量构建、原有测试编译、三段当前文档 Go 范例编译、gofmt 和 diff 检查均通过；按开始时 SHA-256 快照核对 AMQP、搜索、XXL-JOB 和其余三类工具未改动。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-http.L9u6TtTp/cache、GOPROXY=off，无新增依赖或下载。临时测试、overlay、输出、快照、范例源码/产物和缓存均在同一任务目录。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/commonTool
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-http.L9u6TtTp/overlay.json -run '^$' ./common-pkg/commonTool
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-http.L9u6TtTp/overlay.json -run '^TestHTTPTool' -count=1 -v ./common-pkg/commonTool -timeout=45s
go build -mod=readonly -buildvcs=false -o <临时目录>/example_N <临时目录>/example_N.go  （N=1..3，逐一执行）
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/commonTool/curl.go common-pkg/commonTool/http_config.go
git -C common-pkg diff --check
git diff --check
rg -n 'GetRequest\(|PostRequest\(' common-pkg --glob '*.go'  （无匹配，退出码 1）
```

- **结果与限制** 20 组定向测试及竞态检测、三段范例编译、全量构建及原有测试编译通过。仅使用本机 HTTP/TLS 和故障替身，未调用外部业务接口，未执行原有外部集成测试逻辑。同步 JSON 编码无法强制中断，Close 不强制结束在途请求；实际目标接口认证、响应契约和部署网络由接入业务联调。API/RPC 本轮仅更新接入约定，无业务代码或测试产物。
- **测试清理** 清理前已登记实际命令、结果、排障和限制；随后删除本次临时测试、overlay、全部测试输出、源码快照、三份范例源码及编译产物和独立构建缓存。核对任务临时目录及仓库 overlay 目标均不存在，原有五个测试全部保留；AMQP、搜索、XXL-JOB 及其余三类工具的摘要对照通过。API/RPC 仅更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-004 / 加密安全随机字符串

- **用户确认** 第 17 项依次处理，HTTP 已确认；本轮只调整 GenerateRandomString，不改 TOTP、Redis 分片或统一 Trace ID 生成规则。长度 1～4096、固定 A-Z/a-z/0-9、使用加密安全随机源、拒绝采样及错误返回均已在开发前确认。
- **实现** 删除 tool.go 原 math/rand、逐次 rand.Seed 和无错误返回实现；新增 random_string.go，公开签名为 GenerateRandomString(length int) (string, error)，无旧签名包装。使用 crypto/rand.Reader，局部 128 字节缓冲按剩余长度批量 io.ReadFull；只接受 0～247，每个字符对应四个字节值，拒绝 248～255，避免取模偏差。非法长度在读取或分配结果前拒绝，包含极端 int 值；读取失败返回空字符串并用 %w 保留错误链，已生成前缀不返回。无全局播种、缓存或可变源设置，支持并发调用，字段范围和关键逻辑补齐中文备注。
- **接入与说明** 当前仓库 Go 源码无业务调用方，无需业务迁移；新建 [随机字符串使用说明](commonTool/RANDOM_STRING.md)，包含完整错误处理范例、长度/字符集、均匀映射、签名变化与用途边界。生成值不保证唯一，业务按需提供唯一约束、有效期或其他消费控制；本函数不替代专用 Trace ID/TOTP。API/RPC 本轮仅更新接入记录。
- **验证结果** 6 组 TestRandomString 及 -race 全部通过：非法长度（负数、0、4097、最大/最小 int）拒绝且不读源；公开接口 1/2/32/61/62/127/128/129/255/256/4095/4096 长度和字符集；确定性输入全部 256 个字节，先拒绝八个高值，再验证四轮完整字符集、每字符恰好四次且源耗尽；连续两批拒绝值、逐字节短读和跨批补足；读源即时/部分/已生成前缀后的故障、EOF/UnexpectedEOF 的 errors.Is 及不返回前缀；64 个 goroutine 各 16 次共 1024 次调用的长度/字符集与竞态检查。不以概率性的“所有结果不重复”判断正确性。
- **构建与范围核对** 一段文档 Go 范例抽取编译通过；全量构建、原有五个测试的编译（不执行用例）、gofmt、diff 检查通过。按开始时 SHA-256 快照核对除本轮 tool.go 和开发记录外的原有公共包文件未变，HTTP、AMQP、搜索、XXL-JOB、TOTP、Redis 分片和原有测试全部保留。定向测试使用临时 overlay，不在仓库新增测试文件，无下载或新依赖。
- **实际执行命令** 在平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-random.gmftllno/cache、GOPROXY=off。临时测试/overlay、结果输出、源码快照、范例源码/二进制及构建缓存均在同一任务目录。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-random.gmftllno/overlay.json -run '^TestRandomString' -count=1 -v ./common-pkg/commonTool -timeout=45s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-random.gmftllno/example /private/tmp/new-platform-random.gmftllno/example.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/commonTool/tool.go common-pkg/commonTool/random_string.go
git -C common-pkg diff --check
git diff --check
rg -n 'GenerateRandomString|rand.Seed' common-pkg --glob '*.go'
```

- **已知限制** 故障源经内部读取函数注入，未替换全局 crypto/rand.Reader 或人为破坏系统熵源；真实公开接口使用本机系统随机源验证。部分 Go 版本的系统熵源发生底层致命故障会终止进程，本工具只传播以 error 返回的失败，不捕获运行时致命退出；当前工具链为 go1.26.5，模块基线 go1.22，未另装旧工具链运行。定向测试验证字符映射和拒绝逻辑，不构成对操作系统随机性的统计认证；原有外部服务集成测试只编译，未连接外部服务。
- **测试清理** 已在清理前登记实际命令、结果和限制，随后删除本次临时测试、overlay、测试输出、源码快照、文档范例源码/二进制及独立构建缓存。核对任务临时目录和仓库 overlay 目标均不存在，原有五个测试全部保留；无关原有文件 SHA-256 摘要一致，公共包新增文件仅为随机工具源码及使用说明。API/RPC 仅更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-005 / TOTP 密钥、URI 与校验

- **用户确认与范围** 随机字符串已确认，继续第 17 项 TOTP；固定 HMAC-SHA1、六位、30 秒，接受当前及上一个窗口，不计算未来窗口。删除旧接口，统一密钥生成、URI 和字符串验证码校验，不修改 Redis 分片或其他工具。
- **接口与密钥** 删除 google_code.go 的 GenGoogleAuthSecret、BuildGoogleQR、CheckGoogleAuthCode 及三个旧私有转换/计算函数，新增 totp.go 的 GenerateTOTPSecret、BuildTOTPURI、VerifyTOTP，无兼容包装。新密钥由 crypto/rand.Reader 读取 20 字节，直接 Base32 编码为 32 个大写无填充字符；读源失败不返回部分值并保留错误链。输入密钥统一严格校验字母表、无填充和规范尾部位，拒绝空值、小写、空白、CR/LF 和非法编码，按 RFC 4226 最小要求解码后至少 16 字节；生成仍采用推荐的 20 字节。配置常量、公开接口、错误及关键流程补齐中文注释。
- **URI** account/issuer 必填、非纯空白、有效 UTF-8，不含冒号或控制字符，保留原值。net/url 分别编码标签与查询，RawPath 保护嵌入斜杠，空格按 %20 编码，发行方前缀与 issuer 参数一致；固定 algorithm=SHA1/digits=6/period=30。只返回 otpauth URI，不生成图片或访问二维码服务，错误文本不输出账号、密钥或 URI。
- **校验与错误** VerifyTOTP 接收字符串，严格六位 ASCII 数字，保留前导零，不自动补零或清理空白。一次读取服务端时间，8 字节大端计数器和 RFC 动态截断，当前及前一候选均做常量时间比较，不提前命中返回；起始窗口避免减一下溢，负 Unix 时间拒绝。ErrInvalidTOTPSecret、ErrInvalidTOTPCode、ErrTOTPMismatch 支持 errors.Is。工具无状态且可并发，重复消费、尝试次数限制及绑定确认由业务原子处理。
- **说明与接入** 新建 [TOTP 使用说明](commonTool/TOTP.md)，包含密钥规则、URI 编码、窗口边界、两个完整 Go 范例、错误分类及 API/RPC 状态管理职责；URI 按密钥保护并交本地二维码组件，绑定需首次校验确认，RPC 负责受控状态存取，API 编排。当前仓库无旧接口业务调用方，源码检索无旧函数残留；API/RPC 只登记接入约定。
- **定向验证** 12 组 TestTOTP 及 -race 全部通过：确定性随机字节编码和真实系统生成；短读、EOF、UnexpectedEOF 和注入故障不返回前缀；16/20/32/64 字节密钥及全部非法格式/非规范填充位；中文/空格/斜杠/加号/百分号/查询字符 URI 往返和五个参数唯一性；非法标签；RFC 4226 十组六位向量；RFC 6238 六组 SHA1 向量按固定六位取末六位，覆盖前导零及 2038 年后；Python 标准库 hmac/hashlib/struct/base64 独立生成 40 组参考结果（四种密钥长度、十个时间点，含计数器超过 32 位）；当前/上一/未来/过旧窗口、Unix 起点、秒/时区边界及负时间；输入格式及错误脱敏；公开接口和重复调用保持无状态；64 goroutine 各 8 次共 512 组生成/URI/校验及竞态检查。两段当前文档 Go 范例编译通过。
- **全量检查与既有失败** commonTool 包单独构建通过，未受影响的四个原有测试包及 language-pkg 测试编译通过。全量构建与原有测试编译失败：es/esv7/model/es_model.go:75 和 es/esaws/model/es_model.go:83 仍调用已不存在的 commonTool.GenUUID，连带 es/test 无法编译。以 HEAD 原 google_code.go 字节建立基线 overlay、隐藏本轮 totp.go 后再次全量构建，同样两处失败；原文件与本轮开始快照哈希一致，确认问题在任务前已存在。本轮未改动搜索或恢复旧 UUID 接口，不把全量失败记为通过；该问题后续单独讨论处理。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-totp.izj61k8k/cache、GOPROXY=off；定向设置 TOTP_REFERENCE=/private/tmp/new-platform-totp.izj61k8k/reference.json。临时测试和两个 overlay、Python 参考数据、测试/构建输出、原源码快照、抽取范例与产物、缓存都在此任务目录，无依赖下载或新增依赖。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-totp.izj61k8k/overlay.json -run '^TestTOTP' -count=1 -v ./common-pkg/commonTool -timeout=45s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-totp.izj61k8k/example_1 /private/tmp/new-platform-totp.izj61k8k/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-totp.izj61k8k/example_2 /private/tmp/new-platform-totp.izj61k8k/example_2.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...  （失败：既有 GenUUID 缺失）
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1  （同一既有编译失败）
go build -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-totp.izj61k8k/baseline-overlay.json ./common-pkg/... ./language-pkg/...  （还原原 TOTP 后重现同一失败）
go build -mod=readonly -buildvcs=false ./common-pkg/commonTool
go test -mod=readonly -buildvcs=false ./common-pkg/commonTool ./common-pkg/aesGCM ./common-pkg/es/esv7/core ./common-pkg/es/esaws/core ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/commonTool/totp.go
git -C common-pkg diff --check
git diff --check
rg -n 'GenGoogleAuthSecret|BuildGoogleQR|CheckGoogleAuthCode|oneTimePassword|func toBytes|func toUint32' common-pkg --glob '*.go'  （无匹配，退出码 1）
```

- **限制** 本机系统随机源、标准向量及故障注入验证，无真实认证器扫码、部署时钟或外部业务联调。读源故障经内部函数注入，未替换全局随机源；Go 系统熵源的致命退出不转换为普通错误，未另装 Go 1.22 工具链验证。格式正确后的候选比较为常量时间，输入校验和整个函数不宣称恒定耗时；不保存使用状态或实现业务限流。原有集成测试未执行，仅尝试编译，全量失败如上登记。
- **测试清理** 清理前已登记实际命令、结果、基线失败及限制；随后删除本轮临时测试、两个 overlay、Python 参考数据、全部测试/构建输出、原源码快照、两份范例源码和编译产物及独立构建缓存。核对任务临时目录及仓库 overlay 目标均不存在，原有五个测试全部保留；除删除旧 TOTP 源码及更新开发记录外，全部原有公共包文件摘要一致，新增文件仅为 totp.go 和使用说明。API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-006 / Redis 逻辑分片错误处理

- **用户确认与范围** 第 17 项最后处理 Redis 分片，TOTP 已确认。只调整 ShardIndexXX/ShardTagXX 的错误返回和说明；分片数必须显式大于 0，不设默认，合法参数路由保持不变。上轮既有 GenUUID 编译问题另列处理，不混入本轮。
- **接口与错误** 两个函数增加 error 返回，参数 n 改名 shardCount；新增 ErrInvalidShardCount，零分片时索引返回 0、标签返回空串且支持 errors.Is，不再发生除零 panic，也不静默选用默认桶。ShardTagXX 先传播索引错误，成功后保持 {s%03d} 格式。删除旧签名，无兼容包装；当前仓库无业务调用方需要迁移。
- **路由语义与注释** 保留有符号十进制用户 ID 文本的 xxhash64 取模，支持全部 int64 ID 和非零 uint64 桶数，返回索引在 [0, shardCount)。标签至少三位，超过三位自然扩展。补齐接口和关键逻辑中文注释，明确业务桶数、Redis 槽和节点的区别，同标签同槽、不同标签可同槽，以及变更桶数的迁移影响。
- **使用说明** 新建 [Redis 分片使用说明](commonTool/REDIS_SHARDS.md)，包含新签名、错误处理、固定算法、完整关联 key 范例及配置/迁移约定。范例同时保留逻辑桶标签和用户 ID，防止同桶不同用户覆盖 key；固定前缀避免花括号改变 Redis 标签解析。工具只计算索引/标签，不选择节点、访问 Redis 或迁移数据；API 继续通过 RPC 编排，实际数据存取方统一管理分片配置。
- **定向验证** 6 组 TestRedisShards 及 -race 通过。先抽取本轮修改前原函数执行，记录 11 个用户 ID（含最小/最大 int64、0、负数及超过浮点精确范围值）与 11 个桶数（含 1、999/1000/1001、16384/16385、2^63、最大 uint64）的 121 组索引和标签，再逐项验证新实现完全一致。覆盖所有代表 ID 的零桶错误和 errors.Is、单桶固定 0/{s000}、10000 个输入的索引范围/标签形状及数值一致性、大编号无截断/溢出，以及 64 goroutine 各 128 次共 8192 组并发稳定路由。使用修改前结果作回归基准，没有靠重复实现同一哈希公式生成预期。
- **构建与既有问题** 使用说明的一段 Go 范例编译、commonTool 单独构建、四个未受影响原有测试包及 language-pkg 测试编译通过；gofmt、diff 和文档本地链接检查通过。全量构建仍失败于 es/esv7/model/es_model.go:75 与 es/esaws/model/es_model.go:83 对已删除 commonTool.GenUUID 的调用，与上轮已还原基线复核的问题一致；搜索源码按任务开始摘要核对未改，不恢复旧接口、不将全量失败记为通过。本轮未重复运行已知无法编译的 es/test 集成测试。
- **实际命令** 在平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-shards.ks0fb_be/cache、GOPROXY=off；定向设置 SHARDS_GOLDEN=/private/tmp/new-platform-shards.ks0fb_be/golden.json。原函数抽取程序、121 组结果、overlay 测试/配置、测试和构建输出、源码摘要、文档范例源码/产物及缓存均在此任务目录，无依赖下载或新依赖。

```text
go run -mod=readonly -buildvcs=false /private/tmp/new-platform-shards.ks0fb_be/baseline.go  （原函数产生 golden.json）
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-shards.ks0fb_be/overlay.json -run '^TestRedisShards' -count=1 -v ./common-pkg/commonTool -timeout=45s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-shards.ks0fb_be/example /private/tmp/new-platform-shards.ks0fb_be/example.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...  （失败：既有 GenUUID 缺失）
go test -mod=readonly -buildvcs=false ./common-pkg/commonTool ./common-pkg/aesGCM ./common-pkg/es/esv7/core ./common-pkg/es/esaws/core ./language-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false ./common-pkg/commonTool
gofmt -l common-pkg/commonTool/redis_shards.go
git -C common-pkg diff --check
git diff --check
rg -n 'ShardIndexXX|ShardTagXX' common-pkg --glob '*.go'
```

- **限制** 纯函数、原路由样本与并发验证，不连接真实 Redis 集群，不验证实际负载热点或执行数据迁移。原有外部集成测试未执行，相关编译检查仅使用 -run '^$'；全量构建限制如上登记。API/RPC 本轮只更新约定，无业务代码或测试产物。
- **测试清理** 清理前已登记实际命令、结果和既有失败；随后删除本次原函数抽取程序、121 组结果、临时测试及 overlay、全部测试/构建输出、源码摘要、范例源码/编译产物及独立构建缓存。核对任务临时目录及仓库 overlay 目标均不存在，原有五个测试全部保留；除 redis_shards.go 与开发记录外，原有公共包文件摘要一致，新增文件仅为使用说明。API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-007 / 搜索模型 UUID 旧引用修复

- **用户确认与范围** 第 17 项 Redis 分片已确认，新增第 18 项处理阻断全量构建的 ES7/OpenSearch GenUUID 旧引用。两套模型空 ID 创建均复用现有 NewUUidV7，生成失败直接返回并保留错误链；不恢复旧接口，不改 UUID 生成器、RPC/API 业务或其他组件。
- **实现与备注** 两处 documentID 在空 ID 且 generate=true 时调用 commonTool.NewUUidV7()，失败以 %w 包装包含各自组件和阶段的错误，返回空 ID；InsertSchema 在正文编码和 SDK 调用前退出。补齐私有 ID 解析与公开创建接口的中文注释，说明 UUIDv7、失败无写入、业务 ID 保留及不写回对象。已有 ID 继续原样使用，UpdateSchema 和 OpenSearch IndexSchema 缺少 ID 仍返回错误，不隐式生成。
- **使用说明** 更新两套独立 README，明确空 create ID 为 36 字符 UUIDv7、同一 SDK 操作重试复用 ID、再次业务调用会重新生成、需要稳定重发标识时预先指定 ID；生成值不是业务自增序号，同一毫秒不保证递增。无新增文档代码范例或接口签名变化。
- **定向验证** 8 组测试及 -race 全部通过，两套模型各 3 组正常行为和 1 组生成失败。正常验证指定业务 ID 在生成开关两种模式下保留、真实生成的 UUIDv7 格式/版本/variant、nil/typed nil/非法 Schema、空 ID 更新拒绝；SDK 503 后成功重试保持同一文档 ID 和正文，自动 ID 不回写且正文只编码一次，响应关闭；两次指定 ID 创建使用同一 ID，更新保留 ID。OpenSearch 另覆盖缺少 ID 的全量覆盖拒绝。
- **失败路径验证** 第二份临时 overlay 仅将 NewUUidV7 替换为返回 os.ErrPermission，分别验证 documentID/InsertSchema 的 errors.Is 错误链、空返回值、零正文编码、零 SDK/传输调用及调用方未修改；指定 ID 创建和更新仍成功，证明这两类操作不依赖生成器。通过临时 SDK 客户端工厂设置正常请求预算并使用内存 HTTP RoundTripper，无本地监听或外部请求；未在产品代码引入可替换函数变量、测试开关或兼容入口。
- **全量恢复与范围核对** 无 overlay 的全量构建和全部原有测试编译恢复通过，原 es/test 的编译阻断解除；未执行原有外部服务集成测试逻辑。Go 源码检索无 GenUUID 残留，gofmt、diff 和说明本地链接检查通过。按任务开始 SHA-256 快照核对，除两套模型源码、README 和开发记录外的原有公共包文件全部未变，UUID 生成器与原有五个测试保留。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-uuid-ref.51cm8eoh/cache、GOPROXY=off；故障命令额外 UUID_REF_FAILURE=1。两份模型测试、两份临时 core 工厂、失败生成器、两个 overlay、测试/构建输出、源码摘要及缓存均在同一任务目录，仓库内仅通过 overlay 虚拟注入，无依赖下载或新增依赖。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-uuid-ref.51cm8eoh/overlay.json -run '^TestUUIDRef(DocumentID|CreateAndRetry|UpdateRequiresID)$' -count=1 -v ./common-pkg/es/esv7/model ./common-pkg/es/esaws/model -timeout=45s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-uuid-ref.51cm8eoh/failure-overlay.json -run '^TestUUIDRefFailure$' -count=1 -v ./common-pkg/es/esv7/model ./common-pkg/es/esaws/model -timeout=45s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/es/esv7/model/es_model.go common-pkg/es/esaws/model/es_model.go
git -C common-pkg diff --check
git diff --check
rg -n 'GenUUID' common-pkg --glob '*.go'  （无匹配，退出码 1）
```

- **限制** SDK 传输替身验证请求和重试行为，无真实 Elasticsearch/OpenSearch 集群联调。生成失败用临时源码 overlay 注入，未破坏系统熵源；现有生成器和其运行时行为不变，不能将 Go 的致命随机源失败转换为普通错误。自动 ID 不回写及跨业务重发可能生成新 ID 的现有行为保留，业务幂等仍由稳定 ID 控制。
- **测试清理** 清理前已登记实际命令、结果和限制；随后删除本轮两份模型测试、两份临时 core 工厂、失败生成器、两个 overlay、全部测试/构建输出、源码摘要及独立构建缓存。核对任务临时目录和四个仓库虚拟新增目标均不存在，原有五个测试全部保留；除两套模型源码、README 及开发记录外的原有公共包文件摘要一致，UUID 生成器未改动，无新增持久测试或产物。API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-008 / 独立时间工具与 UTC 默认值

- **用户确认与补充** 第 18 项 UUID 修复已确认，开始第 19 项时间工具；用户要求新建独立文件管理，并再次明确未指定时区默认按零时区计算。先前已确认统一 TimeTool、显式基准时间、日历边界、严格解析与秒级差值、删除旧接口的范围。
- **文件与接口** 新增 commonTool/time_tool.go，将八个旧时间函数从 tool.go 移除，无兼容包装；当前仓库无业务调用方。NewTimeTool(timezone) 构造只读实例，提供 DayStart/YesterdayStart/MonthStart/LastMonthStart、Date、FormatUnix、Parse、DiffSeconds，以及 DateTimeLayout/DateLayout 和 ErrAmbiguousLocalTime；公开接口、字段、单位、默认值和关键逻辑均补齐中文注释。tool.go 仅保留原有非时间工具，不改其行为。
- **时区默认** 空配置、UTC、零值和 nil 工具均使用 UTC（零时区），不读取或修改 time.Local。显式 IANA 时区在初始化加载，非法名称和 Local 返回错误不回退。引入标准库 time/tzdata，在系统缺少时区数据时提供内置后备，不新增第三方依赖。调用方处理构造错误，实例可跨请求并发复用，基准时间显式传入。
- **日历边界** 先将基准时刻转换到工具时区，再按年月日取日初、前一日和月初；上月从本月第一天回退，避免 3 月 31 日减月归一化回本月。相邻时区切换的偏移与边界参与候选判断，零点不存在取首个有效时刻、重复零点取较早者，整日被跳过时边界为跳变后的首个有效时刻，对应空日区间；不假设每个日历日都等于 24 小时。
- **格式与解析** FormatUnix 输入单位为秒，0 为合法 Unix 起点，负数为起点之前；Date 按指定时区直接格式化。Parse 严格 YYYY-MM-DD HH:MM:SS，拒绝非法日期、额外空白、小数秒及不规范格式，拒绝跳变中不存在或有歧义的本地时间；底层 time.ParseError 保留 errors.As 链，重复时间以 ErrAmbiguousLocalTime 供 errors.Is 判断。DiffSeconds 返回 end-start 实际秒差，允许负值，解析失败不继续计算；使用 Unix 秒差避免 time.Duration 跨约 292 年的饱和。
- **使用说明与迁移** 新建 [时间工具使用说明](commonTool/TIME_TOOL.md)，突出默认 UTC，包含接口/单位、两个完整 Go 范例、日历和夏令时行为、八个旧接口迁移表，以及 API/RPC 按业务契约传入时区的约定。当前时间直接使用 time.Now()/Unix，不新增只读当前时间的包装；工具不自动读取请求头或修改服务时区。
- **定向验证** 15 组 TestTimeTool 及 -race 全部通过：默认/零值/nil UTC、非法时区，UTC/上海/胡志明/印度/尼泊尔/洛杉矶跨日，闰年/跨年/月底和时区跨月，纽约 23/25 小时日，圣保罗零点跳过、哈瓦那重复零点、阿皮亚整日跳过；0/负秒格式化及往返，严格格式/非法日期/不存在时刻与 ParseError 链，歧义时刻及前后合法时刻；秒差正负零、夏令时实际秒差和 0001～9999 年超长跨度，错误阶段和歧义链；Python datetime/zoneinfo 在十个时区和十个时间点独立生成的 100 组日/月边界与格式化参考全部一致；64 goroutine 各 64 次共 4096 组跨实例只读调用及原始基准时间不变。
- **非零服务器时区与全量检查** 单独以 TZ=America/Los_Angeles 再执行默认 UTC 用例及 -race，验证不随服务器本地时区改变；两段说明 Go 范例均抽取编译，独立 main 实际运行输出与文档一致。无 overlay 全量构建、全部原有测试编译、gofmt、diff 和本地文档链接检查通过，Go 源码无旧时间函数残留。按开始时摘要核对，除 tool.go 与开发记录外的原有公共包文件全部未变，上一项未提交的 UUID 修复完整保留。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-time.nz7q40l0/cache、GOPROXY=off，定向设置 TIME_TOOL_REFERENCE=/private/tmp/new-platform-time.nz7q40l0/reference.json；附加默认值命令使用 TZ=America/Los_Angeles。临时测试/overlay、100 组参考、测试/构建输出、源码摘要、抽取范例/产物和缓存均在此任务目录，无依赖下载。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-time.nz7q40l0/overlay.json -run '^TestTimeTool' -count=1 -v ./common-pkg/commonTool -timeout=45s
TZ=America/Los_Angeles go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-time.nz7q40l0/overlay.json -run '^TestTimeToolDefaultUTC$' -count=1 -v ./common-pkg/commonTool
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-time.nz7q40l0/example_1 /private/tmp/new-platform-time.nz7q40l0/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-time.nz7q40l0/example_2 /private/tmp/new-platform-time.nz7q40l0/example_2.go
/private/tmp/new-platform-time.nz7q40l0/example_1
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/commonTool/tool.go common-pkg/commonTool/time_tool.go
git -C common-pkg diff --check
git diff --check
rg -n 'TimeToString|DiffTimeUnix|GetNowTime|GetTodayZeroTimeInt|GetYesterdayZeroTimeInt|GetMonthZeroTimeInt|GetLastMonthZeroTimeInt|GetXZeroTodayTimeInt' common-pkg --glob '*.go'  （无匹配，退出码 1）
```

- **限制** 不执行原有外部服务集成测试，仅编译；未在实际 API/RPC 服务部署或另外安装 Go 1.22 工具链。内置 tzdata 后备通过标准库接入，未删除系统数据以模拟缺失，实际规则版本仍由部署数据或 Go 内置数据决定。无偏移本地字符串不能表达回拨歧义时刻，工具返回错误，需要精确表达时使用带偏移的跨服务协议。工具不提供闰秒、业务工作日或节假日计算，不自动修改业务时区配置。
- **测试清理** 清理前已登记实际命令、结果和限制；随后删除本轮临时测试、overlay、100 组参考数据、全部测试/构建输出、源码摘要、两份范例源码/编译产物及独立构建缓存。核对任务临时目录和仓库 overlay 目标均不存在，原有五个测试全部保留；除 tool.go 和开发记录外的原有公共包文件摘要一致，包括上一项未提交的 UUID 修复；新增持久文件仅为 time_tool.go 与使用说明。API/RPC 只更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-009 / RSA 初始化与可复用签名验签

- **用户确认与范围** 第 19 项时间工具已确认，开始第 20 项 RSA。密钥初始化校验一次，签名/验签实例可复用，支持完整 PEM 的四种格式，项目最小 2048 位；保留 SHA256 + RSA PKCS#1 v1.5 和标准 Base64，删除旧接口，补齐使用说明与中文备注。
- **接口与生命周期** 重写 commonTool/rsaTool.go，新增 NewRSASigner/RSASigner.Sign 与 NewRSAVerifier/RSAVerifier.Verify，删除 GenerateSignature/VerifySignature，当前仓库无业务调用方需迁移。密钥字段私有且初始化后只读，私钥校验并预计算后交给并发调用，不引入可变全局密钥、每次请求解析或兼容包装。nil/零值实例返回 ErrRSAUninitialized，签名失败不返回部分值。
- **密钥规则** 私钥支持 PKCS#1 RSA PRIVATE KEY、未加密 PKCS#8 PRIVATE KEY，公钥支持 PKIX PUBLIC KEY 和 PKCS#1 RSA PUBLIC KEY。只接受完整唯一 PEM，允许外层空白和 LF/CRLF，拒绝非 RSA、加密类型或 PEM 头属性、多个块、垃圾前缀/尾部和损坏 DER；额外进行 ASN.1 单值检查，补足部分标准库版本 PKCS#8 会忽略外层 DER 尾部的问题。两侧模数至少 2048 位，检查正奇数模数及 3～2^31-1 范围的奇数指数，私钥额外 Validate。
- **签名与错误** Sign 对原始 data 计算一次 SHA256，固定 PKCS#1 v1.5，输出标准 Base64；空正文合法，不做 JSON/文本归一化。Verify 先按模数限制编码长度，再严格解码、检查解码长度和规范重编码，拒绝换行、空白、缺失填充、非零填充位及 Base64URL 专用字符，阻止超大输入分配。ErrInvalidRSAKey/ErrInvalidRSASignature/ErrRSAUninitialized 支持 errors.Is，内容不匹配保留 rsa.ErrVerification，Base64 等底层错误保留 errors.As 链；公开错误文本仅输出阶段，不包含密钥、签名或正文。
- **说明与接入** 新建 [RSA 使用说明](commonTool/RSA.md)，包含密钥格式、实例复用、两个完整范例、错误分类、字节协议、轮换和业务状态边界。主范例从环境变量指定的文件读取配对密钥并本地签名验签，不打印密钥/签名。API/RPC 各按职责注入实例，业务约定序列化和重放控制；本轮不新增业务接口或密钥管理服务。
- **定向验证** 最终 14 组 TestRSATool 及 -race 全部通过：四种格式交叉签名验签、空/中文/二进制/大整数 JSON 原始字节不变；外层空白/CRLF 和全部非法 PEM 包装；加密/属性/错误类型/证书拒绝；ECDSA/Ed25519 非 RSA 拒绝；1024 位、非法指数、偶数模数、损坏私钥数学参数及坏 DER 拒绝；正文/签名篡改与同位数错钥的 rsa.ErrVerification；标准 Base64 编码边界、尾部位、超大输入及底层 CorruptInputError；nil/零值；与标准库互通；本机 openssl（LibreSSL 3.3.6）双向签名/验签及确定性签名字节一致；错误脱敏；32 goroutine 各 16 次共 512 组同实例并发复用；四种 DER 的尾部数据拒绝；3072 位密钥完整成功路径。
- **执行与排障** 首轮 12 组中 11 组通过，故障私钥样本在辅助函数尝试先编码 PKCS#8 时被标准库以 invalid CRT exponent 提前拒绝，未进入被测构造函数；临时样本改为直接编码故障 PKCS#1 后通过，未为测试绕过生产校验。随后根据标准库源码确认 PKCS#8 尾部检查差异，新增统一单值检查及对应测试，补充 3072 位用例，最终完整 14 组 -race 通过。两个文档 Go 范例抽取编译，main 使用本轮测试密钥实际运行并完成验签；全量构建、全部原有测试编译、gofmt、diff、旧接口检索及本地链接检查通过。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-rsa.dk6es3qg/cache、GOPROXY=off，定向设置 RSA_TEST_ARTIFACTS=/private/tmp/new-platform-rsa.dk6es3qg/artifacts。临时测试/overlay、生成的测试密钥和签名、输出、源码摘要、范例源码/产物和缓存均在此任务目录；密钥文件按 0600 写入 0700 临时目录，无外部请求、依赖下载或新依赖。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-rsa.dk6es3qg/overlay.json -run '^TestRSATool' -count=1 -v ./common-pkg/commonTool -timeout=60s
openssl dgst -sha256 -sign <临时 artifacts>/test-private.pem -out <临时 artifacts>/openssl.sig <临时 artifacts>/message.bin
openssl dgst -sha256 -verify <临时 artifacts>/test-public.pem -signature <临时 artifacts>/our.sig <临时 artifacts>/message.bin
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-rsa.dk6es3qg/example_1 /private/tmp/new-platform-rsa.dk6es3qg/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-rsa.dk6es3qg/example_2 /private/tmp/new-platform-rsa.dk6es3qg/example_2.go
RSA_PRIVATE_KEY_FILE=<临时 artifacts>/test-private.pem RSA_PUBLIC_KEY_FILE=<临时 artifacts>/test-public.pem /private/tmp/new-platform-rsa.dk6es3qg/example_1
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/commonTool/rsaTool.go
git -C common-pkg diff --check
git diff --check
rg -n 'GenerateSignature|VerifySignature' common-pkg --glob '*.go'  （无匹配，退出码 1）
```

- **限制** 使用本机临时密钥、标准库及 openssl 命令验证，无外部签名协议、密钥托管或部署模式联调，未另外安装 Go 1.22 工具链。原有外部集成测试只编译不执行。工具不处理密钥加密口令、轮换、正文规范化或业务防重放，进程中的密钥生命周期和访问控制由应用管理。
- **测试清理** 清理前已登记实际命令、结果、首轮临时样本问题和限制；随后删除本轮临时测试、overlay、全部测试密钥/签名/正文、测试和构建输出、源码摘要、两份范例源码/编译产物及独立构建缓存。核对任务临时目录与仓库 overlay 目标均不存在，原有五个测试全部保留；除 rsaTool.go 与开发记录外的原有公共包文件摘要一致，新增持久文件仅为 RSA 使用说明。API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-010 / RPC 客户端启动与连接配置

- **用户确认与范围** 第 20 项 RSA 已确认，开始第 21 项 RPC 客户端。独立配置默认调用 10s/启动等待 5s、可选非阻塞、Keepalive 默认关闭，保留现有 Trace/metadata/错误协议，移除硬编码重试并增加幂等 Close。现有 Timeout:10000 未因 go-zero 中间件零值而生效，KeepaliveTime:30000 实际为 30 微秒，本轮显式修正。
- **配置与默认** 新增 rpc/config.go，Config 改用必填 Target/RpcName，增加 CallTimeout、ConnectTimeout、NonBlock 及独立 KeepaliveConfig。零时长取默认、负值拒绝，非阻塞默认 false；Keepalive.Enabled 默认 false，Time/Timeout 默认 5m/20s，PermitWithoutStream 默认 false，启用时 Time 至少 10s 防止静默最小值修正。Target 拒绝空白/控制字符，RpcName 必须为非空白可打印 ASCII，符合 metadata 值要求；配置值复制不修改原对象，配置及公开接口补齐中文用途/默认值/单位备注。
- **初始化与生命周期** NewClient(ctx, Config) 返回自有 *Client，仍实现 zrpc.Client，并提供 Conn()/Close()；删除 Config.Host、GenRpcTarget 和旧无 ctx 构造签名，不保留包装。内部通过 go-zero 非阻塞创建以保留解析/负载均衡，再用启动 ctx 和 ConnectTimeout 等待 Ready，失败关闭连接且保留 context 错误链；NonBlock=true 不等待 Ready。成功后的初始化 ctx 取消不关闭实例，MustNewClient 复用 NewClient 并在失败时 panic，不使用进程退出函数。
- **调用与关闭** 自有一元拦截器为每次调用派生 CallTimeout，保留父级值/取消及更早截止时间，结束释放，不依赖 go-zero 的可选 Timeout 中间件。Close 用 sync.Once/原子关闭标记支持并发重复、nil/零值安全，释放连接并取消在途 RPC，后续一元调用返回 ErrClientClosed；外部已直接关闭 Conn 后仍可安全 Close。关闭不做业务排空，业务先停止新请求并按需等待收尾。
- **链路与重试** 保留原 ClientInterceptor 函数体，继续复制 metadata、设置 rpc_name 和 W3C 链路、恢复业务错误。删除 RetryDialOption，通过 grpc.WithDisableRetry 关闭服务配置驱动的策略重试，保留 gRPC 未写入/远端未处理场景的透明恢复；连接重连不等同于业务重发。公开启动错误只输出阶段，底层诊断保留错误链，不直接输出 Target。
- **定向验证** 16 组 TestRPCClient 协议/生命周期测试及 -race 通过，另 1 组 YAML 加载 -race 通过，共 17 组。覆盖缺省/覆盖/非法值及不修改配置、默认 10s 截止/父级优先/值保留/派生释放、nil 与预取消/非法配置/Must panic/零值关闭；Ready 返回、共享单连接及启动 ctx 成功后取消不影响业务；握手不响应的启动超时、较早父级截止/主动取消均检测到 TCP 连接关闭；非阻塞创建与延迟启动服务；客户端预算/父级超时/主动取消及服务端 ctx 取消和 504/499 映射；两级真实 RPC 保持同一 Trace、业务 metadata 和独立 rpc_name、不修改原上下文；409 业务文案及 401 现有公开文案映射。
- **补充场景** 自定义解析器提供有效三次 UNAVAILABLE 策略，实际业务仅处理一次，证明策略重试已关闭；32 并发 Close、在途取消和 Shutdown/关闭后拒绝；24 并发调用用确定性不同 Trace ID 验证隔离和单连接复用，并启用合法 Keepalive 正常通信；dns/direct 两种目标；NonBlock 的 WaitForReady 等待仍受调用预算并可关闭；原生 Conn 先关后自有 Close。go-zero YAML 显式加载 125ms/350ms/5m/20s，最小配置应用正确默认，避免时长单位退化。
- **执行与排障** 先在沙箱内编译测试，获准后仅本机回环运行。首次 13 组中 12 组通过，临时测试误以为原生 Unauthenticated 的文案会固定替换成 Unauthorized；核对原 xcode 协议保留该类公开文案后，调整测试服务为显式公开的 authorization required 并按原协议断言，未改错误转换代码。补齐目标解析、连接等待和外部关闭场景后完整 16 组通过，再单独验证 YAML 一组；无产品行为为测试绕过。两段 README Go 范例抽取编译通过，未运行外部服务范例；全量构建、原有测试编译、gofmt、diff、旧入口检索及文档链接检查通过。
- **范围核对** 按开始时 SHA-256 快照核对，除 rpc/rpcClient.go、删除 RetryDialOption 的文件及开发记录外的原有公共包文件未变；原拦截器函数体按 HEAD 对比完全一致，上一项未提交 RSA 改动和原有五个测试全部保留。API/RPC 子项目仅登记约定，无业务代码或测试产物。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-rpc-client.4gwddwq4/cache、GOPROXY=off，无下载或新增依赖。临时 overlay 测试、初次/最终输出、源码摘要、两份抽取范例与编译产物、缓存均在同一任务目录，未创建仓库持久测试。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/rpc/...
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-rpc-client.4gwddwq4/overlay.json -run '^$' ./common-pkg/rpc
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-rpc-client.4gwddwq4/overlay.json -run '^TestRPCClient' -count=1 -v ./common-pkg/rpc -timeout=45s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-rpc-client.4gwddwq4/overlay.json -run '^TestRPCClientYAMLConfig$' -count=1 -v ./common-pkg/rpc
go build -mod=readonly -buildvcs=false -o <临时目录>/example_1 <临时目录>/example_1.go
go build -mod=readonly -buildvcs=false -o <临时目录>/example_2 <临时目录>/example_2.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/rpc/config.go common-pkg/rpc/rpcClient.go common-pkg/rpc/interceptors/clientInterceptor.go
git -C common-pkg diff --check
git diff --check
rg -n 'GenRpcTarget|RetryDialOption|KeepaliveTime: *30000' common-pkg --glob '*.go'  （无匹配，退出码 1）
```

- **限制与说明** 新建 [RPC 客户端说明](rpc/README.md)，包含默认值、两个完整范例、构造迁移、上下文/重试/关闭语义。只覆盖一元 RPC 的自动预算与现有链路/错误封装，流式 RPC 的上下文由调用方管理；传输沿用既有明文设置，未新增 TLS/凭证配置。Keepalive 检查参数与启用后通信，未等待 5 分钟验证 PING 线上时序；无外部生产发现、认证和代理联调。第三方自定义 resolver 的同步阻塞回调不能被启动 ctx 强制终止，解析器需自行实现时限。原有集成测试只编译未执行，未另装 Go 1.22 工具链。
- **测试清理** 清理前已登记实际命令、结果、初次测试预期问题和限制；随后删除本轮临时测试、overlay、全部编译/初次/最终测试输出、源码摘要、两份范例源码与编译产物及独立构建缓存。核对任务临时目录与仓库 overlay 目标均不存在，原有五个测试全部保留；除两份本轮 RPC 源码与开发记录外的原有公共包文件摘要一致，原 ClientInterceptor 函数体和上一项 RSA 改动完整保留，新增持久文件仅为 rpc/config.go 与使用说明。API/RPC 子项目只更新文档，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-011 / 通用十进制精度工具

- **用户确认与范围** 第 21 项 RPC 客户端已确认，开始第 22 项通用十进制工具；新建 commonTool/decimal_tool.go，字符串输入/输出、小数位 0～18、输入最多 256 字符、拒绝科学计数法，保持向负无穷取整。删除旧 FloorAmountNum，无兼容包装，不修改其他基础工具或业务结构。
- **实现** 新增 FloorDecimal(value string, places int) (string, error)，先校验精度与长度，再检查可选正负号、至少一位 ASCII 整数及可选小数点和至少一位小数。支持前导零/正号，拒绝空白、Unicode 数字、NaN/Inf、指数、千位分隔、十六进制及不完整小数；符号和点计入输入上限。校验通过后才将精度转为 int32 并交现有 shopspring/decimal.NewFromString，直接 RoundFloor + StringFixed，不使用 float64、类型转换 cast 或先放大再 Div，不依赖全局 DivisionPrecision。
- **语义与错误** 输出固定 places 位，0 位返回整数，去掉多余前导零/正号，数值零无负号；负数向负无穷，-1.239 的两位结果为 -1.24，极小负数也按该方向取整。ErrInvalidDecimal 与 ErrInvalidDecimalPlaces 支持 errors.Is，非法时返回空串，不默默修正负位数，不回显原输入或返回可用的部分数值。配置常量、公开接口和关键处理逻辑补齐中文注释；工具无共享可变状态。
- **迁移与说明** 从 tool.go 移除原 FloorAmountNum 及仅用于该函数的 decimal/cast 导入，Md5 原逻辑保留；当前仓库无业务调用方。新建 [十进制使用说明](commonTool/DECIMAL_TOOL.md)，包含完整格式、符号/精度/长度规则、向下取整样例、两个可编译范例和 JSON 字符串字段接入。API/RPC 业务契约从入口使用精确字符串；本轮不改数据库或业务接口，不恢复上游已经因 float64 丢失的精度。
- **定向验证** 9 组 TestFloorDecimal 及 -race 全部通过：正负数取整方向、跨整数与极小值；所有 0～18 位固定显示、18 位精确截取及大于 2^53 的整数/小数；符号/前导零/负零规范化；全部非法文本类型（含超大指数）；负数/19/最大和最小 int 精度；恰好 256 字符、257/1MiB 拒绝、符号和小数点计长、254 位极小小数与最长 275 字符输出；错误分类与不回显；Python Decimal（600 位上下文、ROUND_FLOOR）独立生成 1000 组正负/大整数/长小数和 0～18 位参考均一致，再用 math/big.Rat 验证结果不大于输入且差值小于最小精度单位；32 goroutine 各 128 次共 4096 组并发结果保持一致。
- **构建与核对** 两段当前文档 Go 范例抽取编译通过，主程序实际运行输出与文档四行示例一致。全量构建、全部原有测试编译、gofmt、diff、旧函数检索及本地文档链接检查通过。按开始时 SHA-256 快照核对，除 tool.go 和开发记录外的原有公共包文件均未改动；上一轮尚未提交的 RSA/RPC 修改与原有五个测试全部保留，无依赖新增、下载或 go.mod 调整。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-decimal.4pwpkyo0/cache、GOPROXY=off，定向设置 DECIMAL_REFERENCE=/private/tmp/new-platform-decimal.4pwpkyo0/reference.json。临时 overlay 测试、1000 组参考、测试/构建输出、源码摘要、抽取范例及产物、缓存均在同一任务目录，未创建仓库持久测试。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-decimal.4pwpkyo0/overlay.json -run '^TestFloorDecimal' -count=1 -v ./common-pkg/commonTool -timeout=45s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-decimal.4pwpkyo0/example_1 /private/tmp/new-platform-decimal.4pwpkyo0/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-decimal.4pwpkyo0/example_2 /private/tmp/new-platform-decimal.4pwpkyo0/example_2.go
/private/tmp/new-platform-decimal.4pwpkyo0/example_1
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/commonTool/tool.go common-pkg/commonTool/decimal_tool.go
git -C common-pkg diff --check
git diff --check
rg -n 'FloorAmountNum' common-pkg --glob '*.go'  （无匹配，退出码 1）
```

- **限制** 公共工具只处理有界普通十进制文本的向负无穷取整，不实现其他舍入模式、单位换算或业务数据迁移；科学计数法不自动展开。原有外部集成测试只编译未执行，未接入实际 API/RPC DTO/存储协议或另装 Go 1.22 工具链；上游已丢失精度无法恢复，需要业务从源头统一表示。
- **测试清理** 清理前已登记实际命令、结果和限制；随后删除本轮临时测试、overlay、1000 组参考、全部测试/构建输出、源码摘要、两份范例源码及编译产物、运行输出和独立构建缓存。核对任务临时目录与仓库 overlay 目标均不存在，原有五个测试全部保留；除 tool.go 和开发记录外的原有公共包文件摘要一致，包括此前未提交的 RSA/RPC 工作；新增持久文件仅为 decimal_tool.go 和使用说明。API/RPC 子项目只更新约定，无新增测试或临时产物。


### 2026-11-20 / DEV-20261120-012 / 独立哈希摘要工具

- **用户确认与范围** 第 22 项十进制工具已确认，开始第 23 项哈希摘要整理；独立 hash_tool.go 提供 MD5Hex/SHA256Hex，字节切片输入、小写十六进制输出，空输入合法，删除旧 Md5 及只包含该函数的 tool.go，无兼容包装。当前仓库无旧函数调用方，不修改其他基础工具或业务协议。
- **实现与备注** 分别直接调用标准库 md5.Sum/sha256.Sum256，再 hex.EncodeToString，输出 32/64 个字符。避免依赖 crypto.Hash.New 间接注册，不修改输入、不去空白/换行、不进行 Unicode 或 JSON 归一化，nil 与空切片返回同一空内容摘要。函数无共享可变状态，可并发调用，调用方不并发改写同一输入；公开接口补齐用途、编码、输入所有权及 MD5 运行限制的中文注释。旧实现当前可由 TLS 间接引入 MD5 注册，未将其描述为现有调用必然 panic。
- **使用边界与说明** 新建 [哈希摘要说明](commonTool/HASH_TOOL.md)，包含两个接口、完整 main 范例、精确输出、原始字节协议及旧调用迁移。MD5 仅用于协议指定的非安全校验，不用于密码存储或认证签名；SHA256 为内容摘要，不替代来源认证。Go FIPS 只允许模式禁用 MD5 时继承标准库 panic，不 recover、不静默更换算法；本轮不调整部署模式或现有认证流程。
- **定向验证** 6 组 TestHashTool 及 -race 全部通过：MD5 七组标准向量；SHA256 空串、abc、多块及百万字符四组向量；nil/空切片与固定长度；Python hashlib 独立生成的 128 组参考，包括中文、Unicode 组合/分解形式、二进制全部 256 个字节、NUL/非法 UTF-8、空白/换行、55/56/63/64/65/119/120/127/128/129 字节分块边界及 1024/4096 字节和确定性随机内容，两个算法均一致且为小写 hex；输入切片不变及计算完成后改写原切片不影响已返回结果，换行保留；32 goroutine 各 128 次共 4096 组并发双算法结果稳定。
- **构建与核对** 一段当前文档 Go 范例抽取编译并实际运行，两行输出与文档一致。全量构建、全部原有测试编译、gofmt、diff、旧函数检索及本地链接检查通过。按开始时 SHA-256 快照核对，除删除 tool.go 和更新开发记录外的原有公共包文件未变，之前尚未提交的 RSA/RPC/十进制工作及原有五个测试全部保留，无依赖下载或新增依赖。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-hash.4ri5dx3d/cache、GOPROXY=off，定向设置 HASH_REFERENCE=/private/tmp/new-platform-hash.4ri5dx3d/reference.json。临时 overlay 测试、128 组参考、测试/构建输出、源码摘要、抽取范例/编译产物及缓存均在同一任务目录，仓库未写入新测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-hash.4ri5dx3d/overlay.json -run '^TestHashTool' -count=1 -v ./common-pkg/commonTool -timeout=45s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-hash.4ri5dx3d/example /private/tmp/new-platform-hash.4ri5dx3d/example.go
/private/tmp/new-platform-hash.4ri5dx3d/example
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/commonTool/hash_tool.go
git -C common-pkg diff --check
git diff --check
rg -n '\bMd5\s*\(' . --glob '*.go'  （无匹配，退出码 1）
```

- **限制** 仅正常 Go 运行模式验证，未执行 FIPS 模式认证/部署验证，未另装 Go 1.22 工具链；MD5 在禁止该算法的环境中遵循标准库拒绝行为。摘要函数只处理传入字节，不执行 I/O、密码存储、来源认证或业务协议自动迁移。原有外部服务集成测试只编译未执行，API/RPC 本轮只更新接入约定。
- **测试清理** 清理前已登记实际命令、结果和限制；随后删除本轮临时测试、overlay、128 组参考、全部测试/构建输出、源码摘要、范例源码与编译产物、运行输出及独立构建缓存。核对任务临时目录和仓库 overlay 目标均不存在，原有五个测试全部保留；除删除 tool.go 与更新开发记录外的原有公共包文件摘要一致，包括此前未提交的 RSA/RPC/十进制工作；新增持久文件仅为 hash_tool.go 和使用说明。API/RPC 子项目只更新约定，无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-001 / MySQL 多库管理器生命周期

- **用户确认与范围** 第 24 项按已确认方案完善 SelectDb，公开函数签名不变；本轮跨日收尾，用户“下一项”确认本项结果。只修改 dbM/selectDb.go，新增 [使用说明](dbM/SELECT_DB.md)，不修改 MySQL 连接配置和实际连接函数。
- **状态与上下文** 新增 ErrSelectDbUninitialized、ErrDbNotConfigured、ErrSelectDbClosed，支持 errors.Is；nil/零值取库返回未初始化错误，构造的空管理器正常区分缺少配置和已关闭。保留上下文优先校验与 WithContext 独立请求会话，查询/事务继承值、取消及截止时间。启动 ctx 不控制成功后的实例生命周期。
- **关闭行为** 初始化后映射只读，atomic.Bool 标记开始关闭，sync.Once 确保并发/重复 Close 返回同一首次汇总结果；按 DbKey 排序尝试全部池，失败不跳过剩余池。nil/零值 Close 安全；不强制取消在途查询，也不撤销已返回会话，与关闭重叠且已通过检查的 GetDb 仍可能返回会话。应用先停止新业务/按需收尾，再关闭；直接关闭底层池不会同步管理器状态。
- **定向验证** 内存 SQL 驱动配合真实 GORM，共 13 组及 -race 通过：nil/零值、空管理器/未知 key、配置与 ctx 预检查、独立查询/事务上下文、查询超时、全部池关闭及聚合错误/32 并发重复关闭、关闭阻塞期间立即拒绝取库、32 并发取库与关闭、损坏条目不阻断其他池回收、旧会话服从底层池关闭；另三组初始化替身验证成功后启动 ctx 取消不影响请求、部分初始化失败回收并保留全部错误、初始化超时及主动取消回收。
- **执行与排障** 初次临时测试引用了标准库未导出的 sql.ErrDBClosed 名称，编译失败；改为取得已关闭池 PingContext 实际返回的错误，再通过 errors.Is 验证 GORM 查询错误链，无产品代码绕过。普通 overlay 保留实际 ConnectMySQL，只运行无建连的预检查/空构造；初始化 overlay 仅在测试构建中替换连接入口，不接触实际 MySQL。两份文档 Go 范例抽取编译通过，未运行外部数据库示例。
- **实际命令** 平台根目录执行；独立 GOCACHE=/private/tmp/new-platform-selectdb.rj1c3ixl/cache、GOPROXY=off。临时测试、fixture、overlay、初次/最终输出、范例源码和编译产物、快照及缓存均在此任务目录，仓库没有写入新增测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-selectdb.rj1c3ixl/overlay.json -run '^TestSelectDb' -count=1 -v ./common-pkg/dbM -timeout=45s
SELECTDB_INIT_FIXTURE=1 go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-selectdb.rj1c3ixl/init-overlay.json -run '^TestSelectDbFixture' -count=1 -v ./common-pkg/dbM -timeout=45s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-selectdb.rj1c3ixl/example_1 /private/tmp/new-platform-selectdb.rj1c3ixl/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-selectdb.rj1c3ixl/example_2 /private/tmp/new-platform-selectdb.rj1c3ixl/example_2.go
gofmt -l common-pkg/dbM/selectDb.go
git -C common-pkg diff --check
git diff --check
```

- **结果与限制** 10 组正常生命周期测试通过，3 组初始化测试在普通 overlay 下明确跳过、在初始化 overlay 下全部通过；两次均启用竞态检测。全量构建、全部原有测试编译、范例编译、格式、diff 和说明链接检查通过。无新依赖或下载；未连接真实 MySQL，未验证认证/网络/服务部署退出，原有外部集成测试只编译未执行。开始时快照核对除 selectDb.go 外的原有公共包文件一致，已有五个测试保留。
- **测试清理** 已在清理前登记实际命令、结果及限制；随后删除本轮临时测试、fixture、overlay、全部输出、两份范例源码及编译产物、快照与独立构建缓存。核对任务目录和仓库 overlay 目标均不存在，原有五个测试全部保留；除 selectDb.go 与开发记录外的原有公共包文件摘要一致，新增持久文件仅为使用说明。API/RPC 本轮只更新约定，无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-002 / Redis 超时、连接池和重试配置

- **用户确认与范围** 按已确认第 25 项实施：保持 NewRedisClient(ctx, config) 和 redis.UniversalClient 不变，新增独立 Timeout/Pool/Retry 结构体，未配置取明确默认值，默认关闭命令自动重试。修改 dbM/redis.go、ConfigAppServer.RedisConf 中文说明，新增 dbM/redisConfig.go 与 [Redis 使用说明](dbM/REDIS.md)，不修改 MySQL 或其他业务代码。
- **超时与 TLS** RedisTimeoutConfig 使用可选时长指针，Startup/Dial/Read/Write/Pool 默认 5s/5s/3s/3s/4s，显式必须大于 0；读写和等待池超时独立配置，保留 ContextTimeoutEnabled 及父级更早截止时间。TLS 使用配置的 Dial 预算覆盖 TCP/握手，至少 TLS1.2，默认校验证书，沿用已有显式跳过选项。不对某个业务请求取消关闭共享客户端；启动成功后的 ctx 与实例生命周期分离。
- **连接池与校验** RedisPoolConfig 每节点应用，Size 保留 node=10*GOMAXPROCS/cluster=5*GOMAXPROCS 默认，显式须正数；MaxActiveConns/MinIdleConns/MaxIdleConns 默认 0，校验非负、空闲不超过 Size、非零 maxIdle 不低于 minIdle、非零 maxActive 不低于 Size。空闲回收默认 30m、复用年龄默认 0，显式 0 关闭对应回收，负时长拒绝。配置副本独立，不修改原地址切片或指针值。
- **重试与行为变化** RedisRetryConfig.MaxRetries 默认 0，内部映射驱动 -1，防止 node 沿用驱动默认三次重试；MaxRedirects 默认 0 同样关闭，node 非零拒绝。当前 go-redis v9.7.0 集群外层同时处理 MOVED/ASK 和部分网络错误重发，因此明确提供此预算，不宣称只重定向不重发。默认集群可能直接返回 MOVED/ASK，启用两层预算前需业务确认幂等性，内外层可能叠加。退避默认 8ms～512ms，非负且 max>=min，显式零不被驱动替换。
- **驱动细节修正** 集群先把退避 -1 转成 0，节点初始化又会把 0 恢复成默认退避；通过 ClusterOptions.NewClient 在节点构造前重传已校验的原始退避配置，验证节点最终 Options 中显式关闭仍为 0。补充测试也实测出 socket 超时可能早于 ctx 定时器回调，导致启动返回普通网络超时；启动错误处理增加实际截止时间检查，预算已耗尽时统一返回 context.DeadlineExceeded，保留 errors.Is 语义。
- **定向验证** 最终 19 组 TestRedisConfig 及 -race 通过：默认值与两种模式驱动初始化、配置覆盖/副本不变、显式零与哨兵转换、29 类非法组合和 nil/预取消 ctx、YAML 缺省/覆盖/时长单位、TLS 选项；真实 TCP+go-redis 的 node/单节点拓扑 cluster 通信及并发共享；INCR 执行后丢弃响应时默认执行一次、启用一次重试后执行两次；MOVED 预算关闭/启用；启动预算/父级更早截止/主动取消、成功后启动 ctx 取消不影响使用；读超时/父级 deadline、池占满时等待超时/父级取消；TLS 握手无响应时 Dial 超时与连接回收、自签证书默认拒绝且显式跳过可通信；net.Pipe 实际写阻塞超时；空闲/年龄过期连接在下次使用时回收；集群节点最终超时/池/重试配置；node/cluster 流水线响应丢失默认不重发。启动超时/取消组三个场景另重复 20 次并启用竞态检测，全部通过。
- **执行与排障** 首次沙箱内非网络组通过，但协议组因禁止监听回环端口而失败，httptest 同样受限；获执行权限后仅监听本机回环，不连接外部服务，初次 17 组通过。补齐集群节点配置和流水线两组后，其中父级启动超时暴露定时器时序问题；修复生产错误分类后最终 19 组全部通过并重复验证。没有为了通过测试放宽断言或关闭竞态检测；测试使用临时 RESP 替身、内存证书及真实 go-redis 解析/连接池，未改生产连接入口。
- **文档与全量检查** 两段使用说明 Go 范例抽取编译，未运行外部服务范例；抽取文档 YAML 通过 go-zero conf.LoadFromYamlBytes 加载到实际 app.ConfigAppServer，核对全部配置/单位和最小配置的 nil 默认。最终全量构建、全部原有测试编译、gofmt、diff 和本地说明链接检查通过。开始时 SHA-256 快照核对，除两份已说明源码外的原有公共包文件一致，此前未提交开发记录保留；新增持久文件只有配置源码和使用说明。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-redis-config.iqk341il/cache、GOPROXY=off，无新增依赖或下载。临时 overlay 测试、全部初次/中间/最终输出、源码快照、两份抽取范例/编译产物、文档 YAML/加载校验脚本及缓存均在此任务目录，仓库未新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-redis-config.iqk341il/overlay.json -run '^$' ./common-pkg/dbM
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-redis-config.iqk341il/overlay.json -run '^TestRedisConfig' -count=1 -v ./common-pkg/dbM -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-redis-config.iqk341il/overlay.json -run '^TestRedisConfigStartupTimeoutAndCancellation$' -count=20 ./common-pkg/dbM -timeout=30s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-redis-config.iqk341il/example_1 /private/tmp/new-platform-redis-config.iqk341il/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-redis-config.iqk341il/example_2 /private/tmp/new-platform-redis-config.iqk341il/example_2.go
go run -mod=readonly -buildvcs=false /private/tmp/new-platform-redis-config.iqk341il/check_document.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/dbM/redis.go common-pkg/dbM/redisConfig.go common-pkg/app/appServer.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 未启动真实 Redis 服务或多节点集群，未做 Redis ACL/故障转移/跨节点迁槽/生产证书与部署网络联调；协议替身为单节点槽拓扑，验证配置传递和指定故障行为，不替代实际集群验收。业务主动取消不能保证立即打断已阻塞的 socket I/O，业务需设置请求 deadline；阻塞命令有驱动自己的超时规则。原有外部集成测试只编译未执行，未另装 Go1.22 工具链。API/RPC 本轮只更新接入约定，无业务源码变更。
- **测试清理** 已在清理前登记实际命令、结果、排障和限制；随后删除本轮临时 overlay 测试、全部初次/中间/最终输出、源码快照、两份范例源码与编译产物、文档 YAML/校验脚本及独立构建缓存。核对任务目录和仓库虚拟测试目标均不存在，原有五个测试及其摘要保持不变；除 redis.go、appServer.go 和开发记录外的原有公共包文件摘要一致，新增持久文件仅为 redisConfig.go 与使用说明。API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-003 / API 请求体限量读取与统一错误

- **用户确认与范围** 第 25 项已通过用户“下一项”确认，第 26 项按讨论方案实施：删除请求头中间件的无用正文读取，独立配置默认 1 MiB，实际读取超限返回统一 HTTP 413，读取失败不再静默返回，保留 AES 进程配置及原有解密业务错误协议。本轮只修改 API 装配/请求体处理与接入说明，不修改 RPC/数据库/消息处理源码。
- **配置与顺序** 新增 appMiddleware.RequestBodyConfig/RequestBodyMiddleware，MaxBytes 单位字节，未填/0 默认 1048576，负数拒绝，不提供无限制模式。ConfigAppServer.RequestBody 与 app.WithRequestBody 接入；ApiUseMiddleware 在 CORS 和可选并发控制之后、解密之前自动应用，未提供选项也有默认保护，正文读取占用已有并发名额。配置及修改的公开接口和关键流程均补齐中文注释。
- **实际读取与复用** 使用 http.MaxBytesReader 限量读取，不依赖 Content-Length；已知超限提前拒绝，未知/分块/虚报较小长度也无法超过实际读取限制。失败时不进入后续认证/业务，不交付部分正文，并关闭当前读取流。成功后用私有 bufferedRequestBody 提供 Body/GetBody，普通正文只完整读取一次；解密复用未消费的缓存字节，不再复制一遍封装。ApiHeaderMiddleware 只装配请求头/上下文，移除其 io.ReadAll 和 Body 替换。
- **AES 与错误** 独立 NewApiRequestDecryptMiddleware 可选传一份相同 RequestBodyConfig，省略默认 1 MiB；配置先校验再应用 AES，读取错误保留底层链，超限 HTTP 413、其他读取失败 HTTP 400，公开响应不输出原正文或底层错误。AES 密文/封装业务校验失败继续沿用 HTTP 200/code400；空正文、业务头判断与普通请求透传规则保留。解密成功重建 Body/GetBody、更新 ContentLength/Content-Length 并清除旧传输编码；加密器仍从进程配置读取，不经上下文传递。
- **原生长度检查与链路** go-zero v1.7.4 原生 MaxBytesHandler 只判断声明长度且返回空 413。NewServer 在 Trace 初始化后对启用的 RestConf.MaxBytes 进行同值预检查，输出统一结构及 CORS，不创建新的请求 ID。WithRequestTraceRouter 增加可选第三参数对齐原生声明上限，省略默认 1 MiB、0 关闭预检查；原有调用方式继续适用。实际读取上限仍由 ApiUseMiddleware 装配。文档推荐关闭原生 MaxBytes 并统一用平台配置，调大上限同时核对网关及原生约束；自定义原生路由限制的提前响应仍由对应框架层负责。
- **定向验证** 最终共 21 组及 -race 通过：12 组中间件测试覆盖默认/负值/零值实例/重复配置拒绝、精确边界和重复 GetBody、已知/未知/虚报长度超限及最多一个判定字节、立即/部分读取失败/脱敏及底层错误链、nil/NoBody/空正文、请求头不触碰正文、AES 缓存复用/明文长度与精度/链路/重读、独立与组合 AES 超限/读取错误、原有 AES 封装与透传协议、关闭 AES 与 32 并发隔离、读取层取消错误；8 组 app 测试覆盖实际 YAML/选项、原生已知长度统一 413、未知/虚报长度、读失败、完整认证/AES/Trace 链、gzip 解压后超限、预检不读正文及调大至 2 MiB、自定义 Trace Router。
- **真实 HTTP 验证** 另 1 组含两个场景，获执行权限后只监听本机回环，以原始 TCP 发送真正的 Transfer-Encoding: chunked 超限正文和声明长度未完成就半关闭的正文；分别得到统一 HTTP413/400、正文 trace 和响应 X-Request-ID 一致，未进入业务。普通沙箱测试明确跳过这组；单独启用 BODY_HTTP_LOOPBACK=1 后通过，无外部网络请求。
- **范例与全量核对** 新建 [请求体说明](app/REQUEST_BODY.md)，更新 [Trace 接入说明](headInfo/TRACE_ID.md)。两份 Go 范例抽取编译通过；文档 YAML 加载到包含真实 RestConf/ConfigAppServer 的配置，核对默认中间件与 2 MiB 配置，并实际 ServeHTTP 验证超限得到统一 413。全量构建、全部原有测试编译、gofmt、diff 和说明链接检查通过。按开始时 SHA-256 快照核对，除本轮六份源码/接入说明外的原有公共包文件完全一致，包括上轮 Redis 工作；原有五个测试摘要不变。验证中无测试失败或为测试放宽实现的情况。
- **实际命令** 平台根目录执行；独立 GOCACHE=/private/tmp/new-platform-request-body.s6jrfnqg/cache、GOPROXY=off，无新依赖或下载。两个 overlay 测试源、全部输出、源码摘要、两份文档范例源码与编译产物、文档 YAML/校验脚本及缓存位于此任务目录，仓库未创建新增持久测试文件。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/app/... ./common-pkg/appMiddleware/...
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-request-body.s6jrfnqg/overlay.json -run '^TestRequestBody' -count=1 -v ./common-pkg/appMiddleware -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-request-body.s6jrfnqg/overlay.json -run '^TestRequestBody' -count=1 -v ./common-pkg/app ./common-pkg/appMiddleware -timeout=60s
BODY_HTTP_LOOPBACK=1 go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-request-body.s6jrfnqg/overlay.json -run '^TestRequestBodyHTTPChunkedAndTruncated$' -count=1 -v ./common-pkg/appMiddleware -timeout=30s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-request-body.s6jrfnqg/example_1 /private/tmp/new-platform-request-body.s6jrfnqg/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-request-body.s6jrfnqg/example_2 /private/tmp/new-platform-request-body.s6jrfnqg/example_2.go
go run -mod=readonly -buildvcs=false /private/tmp/new-platform-request-body.s6jrfnqg/document_check.go
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/app/appServer.go common-pkg/app/start.go common-pkg/app/traceServer.go common-pkg/appMiddleware/apiHeader.go common-pkg/appMiddleware/apiRequestDecrypt.go common-pkg/appMiddleware/requestBodyMiddleware.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 未对接部署网关/代理、真实业务接口或 HTTP/2；原有外部集成测试只编译未执行，未另装 Go1.22。平台限量读取位于原生 Gunzip 之后，限制解压后的正文，网关传输体积限制独立设置；原生 gzip 头解析/路由额外限制等前置拒绝不由本中间件重写。大小限制不是读取超时，不保证主动取消会打断任意阻塞 Reader，需服务器读取时限配合。先缓存有界完整正文后执行业务，上传/流式接口使用独立处理链；本轮未修改业务解析规则或响应加密协议。API/RPC 子项目只更新接入记录。
- **测试清理** 已在清理前登记实际命令、结果和限制；随后删除本轮两个临时 overlay 测试源、全部测试/构建输出、源码摘要、两份范例源码与编译产物、文档 YAML/校验脚本及独立构建缓存。核对任务目录和仓库两个虚拟测试目标均不存在，原有五个测试与其摘要保持不变；除本轮明确修改的装配/请求体源码、Trace 接入说明和开发记录外的原有公共包文件摘要一致，包括上轮 Redis 工作。新增持久文件仅为请求体中间件和使用说明；API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-004 / API 响应序列化与写入失败处理

- **用户确认与范围** 第 26 项已由用户“下一项”确认，第 27 项按讨论方案开发：响应完整编码成功后再提交，编码失败统一 HTTP200/code500，保留 Trace/多语言；写入失败只记录、不重发或追加。原有公开函数签名不变。修改 result/apiResponse.go，新增 result/responseWriter.go 和 [响应说明](result/README.md)，不修改语言包、AES 实现、请求中间件或业务源码。
- **统一输出** HttpSuccessResult、HttpSuccessResultNotAes、HttpErrorResult 共用 writeJSONResponse/writeResponseBytes。先 json.Marshal 完整响应，再设置既有 application/json; charset=utf-8、提交一次状态码和调用一次 Write。移除依赖框架纯文本编码错误出口及旧的响应打印注释；正常数据不做“预编码再重编码”，自定义 MarshalJSON 每次只执行一次，失败后不再处理同一 data。
- **准备失败协议** 编码失败以及 AES 明文编码/加密失败固定构造安全内部错误（HTTP200/code500/data={}），不沿自定义序列化错误的 XCode 采用 401 等业务码；使用当前上下文语言翻译公开内部错误，trace 沿用当前请求值，不另生成 ID。原始错误只记录服务端阶段诊断；备用错误结构由本包创建且只含数值、字符串和空对象，无可失败自定义序列化器。业务主动调用 HttpErrorResult 的原有错误分类及 HttpErrorResponse 回调行为保持不变。
- **AES 与数据兼容** 仍使用进程配置，明文 string/[]byte 原始字节协议不变，其他类型按 JSON；结构化 data 先翻译再加密，成功封装 aes_data，错误普通 JSON，NotAes 强制普通成功响应。普通 nil data 保持 {}，类型化 nil、字节切片、RawMessage 和数字沿用既有协议；不经 float64 中转，不改变字段名及 id/code 保护规则。
- **传输边界** Write 返回错误或短写时记录当前 Trace 的日志，包含错误或预期/实际长度，不追加备用 JSON、不重试 Write，也不修改已经提交的状态码。一个 handler 选择一个响应出口；本包不接管任意调用方此前已提交的正文。自定义 MarshalJSON 主动 panic 仍由上层恢复机制管理；不新增流式输出或响应大小限制。
- **定向验证** 最终 11 组 TestResponse 及 -race 通过：提交前完整编码、自定义 MarshalJSON 单次调用（普通/AES、显式 NotAes、缺省/中文语言组合）；11 类序列化失败共 88 种入口/语言组合（NaN、正负 Inf、循环、函数、channel、非法 map key、损坏 RawMessage、自定义错误/携带401/返回非法 JSON），全部返回 code500、空对象、正确语言与同一 Trace，公开无内部错误；nil/类型化nil/数组/bytes/string/RawMessage 精确协议；递归翻译/id-code 保护/未知字典、大于2^53整数及小数尾零文本；AES 原始 string/二进制/空正文/null/结构体协议和新 nonce；临时替换随机源模拟加密失败（包含401诊断）仍正确返回500且 NotAes 不受影响；5 类写入失败×4种入口均只提交/写入一次，日志含Trace且不打印原正文；原有包装业务错误、5xx、typed-nil、特殊错误与中文协议；32 goroutine 各16次、普通/AES 两轮共1024份响应并发隔离；真实 HTTP 中间件内存链的 Trace/语言/AES/失败返回；go-zero 全局错误回调协议。
- **执行与排障** 初轮 9 组中，3 组因临时断言漏写现有中文字典“服务暂不可用。”的句号而失败，其中并发组重复报告同一文案差异，未出现竞态报告。核对现有字典原文后仅修正临时预期，未改产品翻译逻辑或字典；加上两组入口/回调测试后最终 11 组全部通过。测试使用内存 HTTP 记录器、受控部分写入器和错误随机源，不监听端口、不请求外部服务；随机源和日志 writer 均按测试恢复，生产没有增加注入开关。
- **范例与全量检查** 说明包含两个完整 Go 范例，两份抽取编译通过，main 在内存 HTTP 请求中实际运行，成功数据保留精确整数/中文翻译，失败返回 HTTP200/code500 与同一请求 ID。全量构建、全部原有测试编译、gofmt、diff 和说明链接检查通过。按开始时 SHA-256 快照核对，除 apiResponse.go 外的原有公共包文件未变，包括上轮请求体处理；原有五个测试摘要不变。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-response.l5own5bq/cache、GOPROXY=off，无新依赖或下载。两个 overlay 测试源、初次/最终测试输出、构建输出、源码摘要、两份文档范例源码与编译产物、运行输出及缓存都在同一任务目录；仓库没有新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-response.l5own5bq/overlay.json -run '^TestResponse' -count=1 -v ./common-pkg/result -timeout=60s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-response.l5own5bq/example_1 /private/tmp/new-platform-response.l5own5bq/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-response.l5own5bq/example_2 /private/tmp/new-platform-response.l5own5bq/example_2.go
/private/tmp/new-platform-response.l5own5bq/example_1
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/result/apiResponse.go common-pkg/result/responseWriter.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 使用内存 HTTP 链与受控 ResponseWriter 验证，未进行真实 TCP/TLS 客户端断连或部署网关联调；传输失败时客户端可能仅收到部分 JSON，本实现不能修复已经发生的网络中断。直接使用其他框架写入函数和上游已提交响应不由本包接管。原有外部服务集成测试只编译未执行，未另装 Go1.22 工具链。API/RPC 子项目只更新接入记录，无业务代码、测试或临时产物。
- **测试清理** 已在清理前登记实际命令、结果、初轮预期问题和限制；随后删除本轮两个临时 overlay 测试源、初次/最终测试输出、构建输出、源码摘要、两份范例源码与编译产物、运行输出及独立构建缓存。核对任务目录和仓库两个虚拟测试目标均不存在，原有五个测试及其摘要保持不变；除 apiResponse.go 和开发记录外的原有公共包文件摘要一致，包括上轮请求体工作。新增持久文件仅为 responseWriter.go 与使用说明；API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-005 / 请求业务 metadata 单值校验

- **用户确认与范围** 第 27 项已由用户“下一项”确认，第 28 项按讨论方案实施：现有 Getter 签名不变，单值字段不再拼接，RPC 入口拒绝重复/非法业务字段，统一 UID/IP 解析，保留 Trace 规则及可选字段。修改 headInfo/headInfo.go、appMiddleware/rpcAuthMiddleware.go，新增 headInfo/metadata.go 和 [使用说明](headInfo/METADATA.md)，同步 [Trace 接入说明](headInfo/TRACE_ID.md)，不修改 HTTP 请求头选择、RPC 客户端拦截器、身份认证判断或其他业务源码。
- **单值定义** 包内只读集合集中定义 16 个已知字段：认证令牌、UID、角色、两种 IP、路径、版本、来源、业务/业务标识、语言、来源地址、时区、User-Agent、设备及 rpc_name。可选字段缺失或单个空字符串允许，存在空列表/nil 列表或两个以上值拒绝，即使内容相同。未知自定义 metadata 及 tracestate 多段传播不新增约束，字段集合不依赖可由业务修改的日志名单。
- **Getter 与格式** 公共 Get... 统一读取出站 metadata 的唯一值，不再 strings.Join，不回退入站；缺失/重复字符串返回空，nil ctx 按未设置。非空 UID 只接受 ASCII 十进制 0..MaxInt64，允许前导零，拒绝符号/空白/其他进制/小数/溢出，非法 Getter 返回 0；空或0表示无身份，并不表示通过认证。ClientIp/AdminClientIp 同用 netip 纯地址校验，拒绝端口/zone/CIDR/列表/域名/空白，Getter 规范化 IPv4-mapped IPv6。其他业务文本不新增内容格式约束，Terminal 和 OriginHostUrl 原有转换保留。
- **校验与入口** 新增 ValidateMetadata 与 ErrInvalidMetadata，支持 errors.Is，校验不改写原 map，错误只描述字段/规则而不回显输入值。RPC 先恢复有效 W3C 链路，再校验业务字段，失败返回 gRPC InvalidArgument，白名单回调/认证/业务均不执行，公开方法也不豁免；公共客户端继续按既有协议映射为业务 code400。仅有业务校验不要求 Trace 必须存在，完整入口原有 metadata/Trace 必填与取消检查仍保持。
- **上下文隔离** SetTokenUid/SetBusinessCode/SetBothBusinessCode 继续复制后覆盖单值，保留其他字段、值/截止时间/取消及 Trace，不对写入值进行身份或格式认证；移除旧注释实现，补齐所有修改公开接口和关键处理的中文说明。入站/出站仍分别持有副本，不通过更改 Getter 来源来隐式混用两个方向。
- **定向验证** 最终 12 组 TestMetadata 及 -race 全部通过。7 组 headInfo 覆盖所有 Getter 的唯一值/重复/缺失/nil/入站不回退、终端及来源地址既有转换，UID 空/0/前导零/MaxInt64/溢出/符号/空白/小数/科学计数/其他进制/Unicode，双 IP 格式与规范化，全部16字段的重复同值/空列表/重复空值、可选空值/未知多值/tracestate，不修改原 map、直接原始 map 的大小写别名检查及日志名单隔离，Setter 双向副本/父级值与截止时间/取消/Trace、32 goroutine各64次共2048份分支隔离。
- **RPC 验证** 4 组入口测试覆盖16字段各类非法结构均在白名单/认证/业务之前拒绝，非法 UID/IP 不回显输入，缺失/空值/自定义多值可通过且上下文不污染，缺失/非法/重复 Trace 与已取消请求保持原规则。另 1 组通过 bufconn 内存连接启动两个真实 gRPC 服务，沿用现有 ClientInterceptor；验证两级 UID/语言/两种 IP/自定义多值/统一 Trace，分别设置 rpc_name，下游修改业务标识不污染上游；真实重复 UID 调用未进入认证/业务并映射 code400；16并发双级调用与初次调用合计17组双跳请求保持副本隔离。
- **执行与排障** 首轮7组中，集合不变性测试用 metadata.MD.Copy 制作快照，该库会把 nil 切片复制为空切片，导致一次测试组错误报告输入被修改；临时快照改为精确保留 nil/空差别后全部7组通过。新增 RPC 临时服务描述时漏写一个复合字面量闭合括号，测试编译失败，补齐并 gofmt 后恢复。上述问题均仅修正临时测试，未修改生产逻辑迁就断言；最终完整12组通过，未出现竞态报告。所有通信通过进程内 bufconn，不监听端口或连接外部服务，无生产注入开关。
- **范例与全量检查** 两份说明 Go 范例抽取编译通过，main 实际输出 12 192.0.2.1 / 0 true / 12，验证规范化、重复拒绝和原分支不变。最终全量构建、全部原有测试编译、gofmt、diff、说明链接检查通过。开始时 SHA-256 快照核对，除本轮 headInfo Getter、RPC 入口及 Trace 接入说明外的原有公共包文件未变，包括上一项响应处理；原有五个测试摘要不变。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-metadata.w3dgb3y1/cache、GOPROXY=off，无新增依赖或下载。三个 overlay 测试源、初次/中间/最终输出、源码快照、两份范例源码与编译产物、运行输出及缓存都在此任务目录；仓库未新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-metadata.w3dgb3y1/overlay.json -run '^TestMetadata' -count=1 -v ./common-pkg/headInfo -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-metadata.w3dgb3y1/overlay.json -run '^TestMetadata' -count=1 -v ./common-pkg/headInfo ./common-pkg/appMiddleware -timeout=60s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-metadata.w3dgb3y1/example_1 /private/tmp/new-platform-metadata.w3dgb3y1/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-metadata.w3dgb3y1/example_2 /private/tmp/new-platform-metadata.w3dgb3y1/example_2.go
/private/tmp/new-platform-metadata.w3dgb3y1/example_1
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/headInfo/headInfo.go common-pkg/headInfo/metadata.go common-pkg/appMiddleware/rpcAuthMiddleware.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 未进行真实 TCP/TLS、部署认证/发现或网关联调；原有外部集成测试只编译未执行，未另装 Go1.22。标准 gRPC 传输键名小写；手工构造大小写重复 map 后再放入上下文时，原生读取 API 可能已覆盖别名，须使用 metadata.Pairs/md.Set 创建规范字段或在放入上下文前调用校验。结构校验不证明身份或 IP 真实性，HTTP 重复头策略、用户扩展 JSON 跨 RPC 传递不在本轮范围。API/RPC 子项目只更新接入记录。
- **测试清理** 已在清理前登记实际命令、结果、临时测试排障和限制；随后删除本轮三个 overlay 测试源、全部初次/中间/最终测试及构建输出、源码快照、两份范例源码与编译产物、运行输出及独立构建缓存。核对任务目录和仓库三个虚拟测试目标均不存在，原有五个测试及其摘要保持不变；除 Getter、RPC入口、Trace接入说明和开发记录外的原有公共包文件摘要一致，包括上一项响应处理。新增持久文件仅为 metadata.go 与 METADATA.md；API/RPC 本轮只更新文档，无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-006 / HTTP 单值头校验与传输字段统一命名

- **用户确认与范围** 第 28 项经用户“下一项”确认，第 29 项按已确认方案开发：HTTP 已知单值头重复时拒绝，放在正文/解密/认证前，保留 Trace、CORS 和可信代理链。验证期间用户要求自定义头改为中划线，并明确选择“HTTP 和 RPC 统一改为中划线”；本轮整合为最终范围，调用方须同步更新，不增加旧名兼容映射，业务 JSON/JWT 标签保持原样。
- **HTTP 规则与副本** 新增 headInfo.NormalizeHTTPHeaders/ErrInvalidHTTPHeaders 和 appMiddleware.RequestHeaderMiddleware，10个实际读取字段（项目令牌、版本、来源、业务/业务标识、语言、时区、设备、Origin、User-Agent）大小写不敏感，重复同值、无值列表和大小写别名均拒绝；可选缺失/单空值允许。成功返回深复制请求头，已知字段转为标准 MIME 大小写，值及其他多值头原样保留，不读取或关闭正文、不修改原请求/上下文。单个逗号文本不拆分，其他组件的原有值格式规则保留。
- **装配与错误** 在 NewServer/WithRequestTraceRouter 的路由外层，Trace 初始化后、原生声明长度检查与 go-zero 中间件前执行，避免原生 gzip 先读取正文。非法头 HTTP400/code400，统一 JSON、同一 trace/X-Request-ID 和 CORS，消息经统一翻译入口，不输出令牌或原始头值；GET、POST、OPTIONS、白名单和旧 NotVerify 均不豁免。有效预检头名列表正常返回204。ApiUseMiddleware 的正文/并发/AES顺序不变，独立 net/http 必须显式装配新头校验器。
- **代理与协议边界** X-Forwarded-For、多行/逗号代理链、Accept/Cookie/自定义多值头保留，由原协议处理；可信代理源码未改。标准HTTP User-Agent/Origin、X-Request-ID 和 W3C 协议沿用；UID、角色、IP 仍由认证/可信代理产生，不采用客户端任意同名头作为身份。
- **统一中划线命名** consts 的传输值改为 authorization-jwt、client-ip、admin-client-ip、client-user-agent、token-uid、token-uid-role、req-path、business-code、content-language、origin-url、device-id，新增共享 RpcName=rpc-name 并替换客户端/字段校验/日志名单中的硬编码。version/source/business/timezone 单词名称不变。配置备注、解密头读取备注、当前使用说明及全局 AGENTS 命名规则同步更新，公开 Go 调用方式不变。
- **gRPC 保留头与兼容** 原终端标准 HTTP User-Agent 经 client-user-agent 传 RPC，避开 gRPC 保留并会忽略应用传入值的 user-agent。CORS 暴露 Authorization-Jwt，日志业务字段键随共享常量同步且新令牌字段继续过滤。Head JSON 标签、JWT claims 和 aes_data 等业务 JSON 不改；旧名称不读取、不映射，仅带旧令牌的受保护请求按未提供令牌处理，未知旧字段可按原多值规则保留但不再成为对应业务值。
- **最终验证** 共20组及 -race 通过。6组 headInfo 覆盖10头的空列表/重复同值/不同值/重复空值、大小写别名（含空别名不能掩盖重复）、可选头/值不变/未知多值、独立副本与32 goroutine各64次共2048分支、16个最终传输名称及新 RPC 字段单值识别、新 Getter 读取/旧名不映射、JSON标签不变和 Head 字符串令牌脱敏。6组中间件覆盖校验在读体/AES/认证前、原上下文/正文保持、可信/不可信/非法/缺失 XFF 回归、Trace与预检规则、nil/空头、新 CORS 名称及旧令牌不登录。6组 app 覆盖 GET/POST/OPTIONS×白名单×10头、先于原生 gzip/声明长度/AES 判断、合法规范化/原请求不变、正常AES及预检、自定义路由，以及 HTTP→真实 gRPC(bufconn) 的身份/业务/语言/Trace/rpc-name 传递，并确认原浏览器UA位于client-user-agent、框架原生user-agent仍为grpc-go。1组日志验证新字段键及令牌过滤。
- **真实 HTTP 回环** 另1组包含3个场景，获执行权限后仅监听本机回环，用原始 TCP 发送新头名的同值重复、大小写别名及多行XFF；前两者统一400且不进业务，后者正确取得可信链客户端IP，响应ID一致。普通沙箱执行明确跳过该组，HTTP_HEADER_LOOPBACK=1 单独启用后通过，不连接外部服务。
- **执行过程** 原单值校验范围先完成16组验证；新增 app 临时RPC测试初次漏导入 grpc/metadata 导致编译失败，仅补齐临时测试导入后全部通过。用户确认扩大到HTTP+RPC命名后，同步调整测试输入到新协议，增加4组命名/Getter/JSON/CORS/日志验证，最终20组通过，无竞态报告。核对本地 gRPC transport 源码确认 user-agent 保留行为，再用真实 gRPC 内存连接验证原终端信息未丢失；未为测试放宽校验或增加生产注入开关。
- **说明与全量核对** 新建 [HTTP 请求头说明](headInfo/HTTP_HEADERS.md)，包含完整更名表；同步请求体、Trace、metadata、RPC说明及全局命名规范。4份最终 Go 范例抽取编译通过，main 在新字段名下实际输出 business 2 / true / business，验证读取、代理值数量、拒绝分类与原副本不变。最终全量构建、全部原有测试编译、gofmt、diff、说明链接及旧传输字面量扫描通过；Go 生产旧字面量仅剩明确保留的JSON标签，无旧头名读写分支。SHA-256核对 JSON/JWT、可信代理及本轮范围外源码不变，原有五个测试摘要不变。
- **实际命令** 平台根目录执行，独立 GOCACHE=/private/tmp/new-platform-http-headers.kfexj2pu/cache、GOPROXY=off，无新增依赖或下载。六个 overlay 测试源、全部初次/中间/最终输出、原始/改名前快照与阶段验证记录、四份范例源码/编译产物/运行输出及缓存均位于同一任务目录，仓库没有新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-http-headers.kfexj2pu/overlay.json -run '^TestHTTPHeaders' -count=1 -v ./common-pkg/headInfo ./common-pkg/appMiddleware -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-http-headers.kfexj2pu/overlay.json -run '^TestHTTPHeaders' -count=1 -v ./common-pkg/headInfo ./common-pkg/appMiddleware ./common-pkg/app -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-http-headers.kfexj2pu/overlay.json -run '^TestHTTPHeaders' -count=1 -v ./common-pkg/headInfo ./common-pkg/appMiddleware ./common-pkg/app ./common-pkg/logs -timeout=60s
HTTP_HEADER_LOOPBACK=1 go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-http-headers.kfexj2pu/overlay.json -run '^TestHTTPHeadersWireDuplicatesAndProxy$' -count=1 -v ./common-pkg/appMiddleware -timeout=30s
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-http-headers.kfexj2pu/header_example_1 /private/tmp/new-platform-http-headers.kfexj2pu/header_example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-http-headers.kfexj2pu/header_example_2 /private/tmp/new-platform-http-headers.kfexj2pu/header_example_2.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-http-headers.kfexj2pu/body_example_1 /private/tmp/new-platform-http-headers.kfexj2pu/body_example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-http-headers.kfexj2pu/body_example_2 /private/tmp/new-platform-http-headers.kfexj2pu/body_example_2.go
/private/tmp/new-platform-http-headers.kfexj2pu/header_example_2
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/consts/system.go common-pkg/headInfo/httpHeaders.go common-pkg/headInfo/metadata.go common-pkg/headInfo/logKeys.go common-pkg/appMiddleware/requestHeaderMiddleware.go common-pkg/appMiddleware/corsMiddleware.go common-pkg/appMiddleware/apiRequestDecrypt.go common-pkg/app/traceServer.go common-pkg/app/start.go common-pkg/rpc/config.go common-pkg/rpc/interceptors/clientInterceptor.go
git -C common-pkg diff --check
git diff --check
```

- **限制与调用方迁移** 前端、网关及手工构造RPC metadata的调用方须同步迁移新头名；本仓库API/RPC子项目暂无业务源码，本轮仅登记约定。上游若把重复行合并成单个逗号文本，本层无法还原原行数；其他字段文本格式与令牌有效性仍由业务负责。HTTP解析器/部署代理提前拒绝的响应不由本包重写；真实网络仅验证本机HTTP/1.1，RPC使用bufconn，未验证部署网关、HTTP/2或TLS。原有外部集成测试只编译未执行，未另装Go1.22工具链。RPC业务结构校验不代替身份认证，标准入口之外需手动组合Trace与新头校验器。
- **测试清理** 已在清理前登记实际命令、结果、编译排障、改名过程与限制；随后删除本轮六个 overlay 测试源、全部初次/中间/最终测试及构建输出、原始/改名前快照和阶段记录、四份范例源码与编译产物、运行输出及独立构建缓存。核对任务目录和仓库六个虚拟测试目标均不存在，原有五个测试及摘要不变；范围外源码、Head JSON/JWT及可信代理实现摘要保持一致。新增持久文件仅为httpHeaders.go、requestHeaderMiddleware.go和HTTP_HEADERS.md；全局AGENTS及各层开发文档已同步，API/RPC子项目无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-007 / 多语言恢复项目本地词典

- **用户确认与范围** 用户要求回退为读取各项目 etc 配置翻译，并确认实施方案：翻译逻辑归回 base-pkg，默认 etc/language，可配置目录，启动一次加载，请求只查内存；错误 message 和 data 字符串值自动精确匹配，缺失返回原文。移除 messageKeys 及内嵌 language-pkg 依赖，不恢复旧反射实现，保留已修复 JSON 编码与响应错误处理。上一项 HTTP/RPC 中划线命名保留。
- **实现** 新增 language.Config/Load，ConfigAppServer.Language 作为独立配置字段。Dir 为空默认工作目录下 etc/language，由启动代码显式 Load(c.Language)，仅填写配置或调用 InitAppServer 不隐式加载。读取直接子级小写 .json 文件，忽略子目录/其他文件；要求扁平字符串对象，拒绝重复 key、重复规范化语言、null/非字符串、尾随数据。缺失目录、无词典或文件错误返回错误，调用方应停止启动。全部成功后原子发布只读词典，失败保留旧词典，不监视文件变更，不在请求中读取文件。
- **翻译语义** 直接查实际原文或业务 key，移除常见错误到稳定 key 的中间映射；完整字符串精确匹配，不做子串替换或裁剪文案。语言大小写/下划线规范化及 en-US 无专用词典时回退 en 保留；缺失语言/key、未加载原样返回，空字符串译文合法。SwitchLanguage 捕获 data 只读词典快照，保留 JSON 标签、omitempty、自定义 MarshalJSON、数字文本和 id/code 直接字符串保护；不改业务对象或字段名，[]byte 不翻译，AES 原始字节约定不变。固定成功 message=success 协议保留。
- **依赖与说明** 删除 base-pkg/go.mod 的 language-pkg 依赖及 go.work 的对应 replace；保留独立 language-pkg 项目和 workspace use，不删除用户已有模块。新增 [本地词典使用说明](language/README.md) 和中英文 [配置示例](language/examples/etc/language)，同步 [响应范例](result/README.md)。示例是可复制的长期配置资料，不是临时测试夹具；历史开发记录保留，当前约定以本记录为准。
- **定向验证** 最终10组 TestLocal 及 -race 全部通过，含16项非法词典子场景。覆盖未加载/空目录配置与实际默认路径、删除文件后请求仍用内存、缺失目录失败保持旧词典、忽略子目录及非 JSON/大写扩展名文件、空对象、损坏/嵌套/null/数字/布尔/数组词典、重复转义 key、重复规范化语言、空语言文件、尾随 JSON、断裂文件链接与加载不部分发布、整体替换、en-US 专用优先及不跨词典补 key。并发16 goroutine各64次共1024次编码，同时64次完整词典切换，无竞态且每份 data 不混用词典。
- **响应与配置验证** 使用长期示例词典，验证错误原文和业务 key、未匹配回退、内部诊断不外泄、编码失败固定500且翻译正确，嵌套数组/map、自定义编码一次、标签/omitempty/隐藏字段、nil、转义字符串、大整数及小数科学计数原始文本，字段名和 id/code 保留。AES 解密验证普通字符串翻译/字节不翻译/结构体先翻译、NaN失败及强制NotAes。内存 HTTP 接入 content-language，覆盖成功、重复业务头400、正文413、并发503，保留统一Trace；go-zero YAML验证未配置/空 Language 与自定义Dir。
- **执行过程** 初轮及第二轮临时夹具使用仅大小写不同的文件名，分别导致 en.json 被 en.JSON 覆盖、ignored.json 目录与 ignored.JSON 文件冲突；这是本机大小写不敏感文件系统的夹具问题。调整为互不冲突的文件名后最终10组通过，未修改生产逻辑迁就测试。全部验证只使用临时文件及内存 HTTP，无网络监听、无外部服务。
- **范例与完整检查** 三份Go范例抽取后编译通过，两个main实际运行：输出中文错误、已启用、保留9007199254740993及id，普通API成功及编码失败均保留原Trace。全量构建及原有测试编译通过；单模块 GOWORK=off list -m all 确认不再依赖 language-pkg。gofmt、diff、文档链接、原有文件SHA-256核对通过，范围外公共包文件及五个原有测试摘要不变。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-local-language.lAa005xI/cache、GOPROXY=off，无下载或新增依赖。两个overlay测试源、快照、各次测试输出、三份范例源码/编译产物/运行输出、测试词典副本与缓存均在该任务目录；测试内部t.TempDir自动删除。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-local-language.lAa005xI/overlay.json -run '^TestLocal' -count=1 -v ./common-pkg/language ./common-pkg/result -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-local-language.lAa005xI/language-example /private/tmp/new-platform-local-language.lAa005xI/language_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-local-language.lAa005xI/response-example /private/tmp/new-platform-local-language.lAa005xI/response_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-local-language.lAa005xI/response_2.a /private/tmp/new-platform-local-language.lAa005xI/response_2.go
# 以下两个范例在任务目录执行，该目录已复制示例etc/language
./language-example
./response-example
GOWORK=off go -C common-pkg list -mod=readonly -m all
gofmt -l common-pkg/language/config.go common-pkg/language/language.go common-pkg/app/appServer.go
git -C common-pkg diff --check
git diff --check
```

- **限制** API/RPC子项目暂无业务入口，本轮仅更新接入记录；启动代码需显式加载并处理错误，部署工作目录/词典挂载与实际网关联调待服务接入。未运行原有依赖外部设施的集成测试，只做全量编译；使用本机Go1.26.5，未另装Go1.22。不提供文件监听、多租户独立词典或自动业务枚举保护；自动匹配词条应明确命名，避免其他业务值意外命中。未改动独立language-pkg实现，不再阻塞base-pkg独立使用或发布。
- **测试清理** 已在清理前登记实际命令、结果、夹具修正与限制；随后删除本轮两个overlay测试源、配置、源码快照、所有测试/编译/依赖输出、三份范例源码及编译产物、运行输出、临时词典副本和独立缓存。核对任务目录及两个仓库虚拟测试目标均不存在，原有五个测试文件及摘要不变。清理核对脚本初次把子场景计入10组总数，改为只统计顶层测试后核对通过。范围外公共包原有文件摘要一致；持久新增文件仅为language/config.go、使用说明及两份长期示例词典。API/RPC/language-pkg仅更新记录，无新增测试或临时产物。


### 2026-11-21 / DEV-20261121-008 / 白名单仅豁免用户必须登录

- **用户确认与范围** 用户明确白名单配置不改，白名单接口也执行统一校验，只对userid=0允许匿名；进一步选择无效/过期Token返回认证错误，随后确认开发。本轮保留map[string]int、公开白名单配置接口及精确/尾部星号字符串前缀匹配，不做原建议的独立配置结构、目录边界变更或map副本改造。第30项本地多语言已由用户“下一项”确认，上一轮源码及未提交变更保留。
- **HTTP认证** authenticateToken将没有Token（含空值/空白）返回正常匿名结果；所有非空Token照常调用业务认证回调，错误不再被白名单吞掉。仅无错误且uid=0时查白名单，非白名单返回原未登录401，白名单显式写入token-uid=0并清空扩展用户JSON；uid>0使用回调身份和JSON，uid=-1无权限、其他负数无效、非空非法JSON服务不可用均拒绝。三类回调入口及错误映射保留，原对象/metadata/截止时间/Trace不被破坏。
- **回调约定** 有error的回调必须把无效/过期Token作为错误返回，不可用0,nil冒充匿名；无error整数回调仍保留，0表示真实匿名，-1无权限，其他负数无效。旧业务若用0表达认证失败，须调整返回值或使用带error接口，框架无法仅凭0判断业务是否实际完成Token验证。
- **统一校验** ApiHeaderMiddleware移除GET、白名单与NotVerify跳过Head.Verify的分支，保留CloseVerifyOption、NotVerify字段及SetNoVerify公开调用形状并说明不再豁免。Head.Verify当前为空实现，本轮未擅加business/source等必填规则，实际单值校验继续由标准入口前置RequestHeaderMiddleware负责。CheckServiceStatusRequest保留白名单参数，但不据此跳过回调，缺少回调/返回false仍拒绝；状态中间件由业务显式接入，不增加自动回调或新配置。
- **RPC规则** Trace、metadata结构与取消校验保留；之后对所有方法执行CheckRPCRequestFunc，白名单也不跳过调用方认证。调用方认证成功后再调用原白名单回调，UID缺失/空/0视为匿名，仅白名单放行，非白名单返回标准Unauthenticated/客户端业务401；正数UID不受白名单限制。业务回调仍负责核实服务调用方与Token/用户真实性，不直接信任任意正数metadata；后台无用户调用需明确匿名方法范围。原RPC认证错误安全映射及panic恢复保持。
- **使用说明** 新增 [认证与匿名白名单](appMiddleware/AUTH.md)，含行为矩阵、HTTP/RPC装配范例、三类回调返回约定、状态检查和Head.Verify实际边界。同步app白名单选项注释、[HTTP头](headInfo/HTTP_HEADERS.md)、[metadata](headInfo/METADATA.md)和[RPC](rpc/README.md)说明，新增/修改公开接口与关键逻辑有中文注释。历史记录中的白名单吞认证错误/跳过RPC调用方认证规则由本条替换。
- **定向验证** 首轮8组中间件测试及-race通过；增加API/真实RPC内存链路后，最终10组及-race全部通过，无竞态报告。HTTP认证矩阵覆盖5路径×13结果，包括无Token/空白、正常匿名及扩展信息丢弃、登录/无扩展、无效/过期/正UID伴错误、权限/负UID、特殊503、合法/非法匿名扩展JSON，验证精确/旧前缀规则、原身份/metadata不变、Trace及截止时间保留。三类适配器共15个子场景覆盖白名单成功与错误不入业务，另验证缺少回调的装配拒绝。状态检查验证精确/前缀白名单和受保护路径均执行、匿名上下文保留、不可用/缺少回调拒绝。
- **并发与RPC验证** 16goroutine各64次共1024个HTTP认证分支，匿名与登录身份及扩展JSON不串用、父上下文不变。RPC定向矩阵覆盖公开/受保护×缺失/空/0/前导零/正UID，调用方认证总执行、白名单只豁免匿名；错误/权限/取消/超时/不可用、不可信调用方、缺少认证器与panic均不能绕过，公开诊断不泄露。Trace/metadata缺失、重复UID、负UID及已取消在认证和白名单之前拒绝。
- **API与真实RPC内存链路** 使用实际app.NewServer/ApiUseMiddleware、服务状态组件和grpc.ServerOption，通过bufconn执行真实gRPC及既有ClientInterceptor。13个业务场景验证匿名/登录、受保护未登录、原匹配边界、过期/无效/权限/后端错误以及白名单调用方认证失败；另3个场景验证白名单重复头400、正文413、状态不可用503均不进入RPC业务。验证API用户JSON保留或清空、RPC不自动接收用户扩展JSON、语言及统一Trace传递、CORS保持。直接真实RPC另验证5个访问场景，16goroutine各8次共128个匿名/登录请求并发隔离。
- **范例与全量检查** AUTH说明的两份Go装配范例抽取编译通过。全量构建、原有测试全量编译、gofmt、diff及说明链接检查通过；SHA-256对照本轮开始快照，五个原有测试和上一轮本地多语言实现/配置/响应范例摘要不变。原有字段常量、可信IP、JWT/AES、错误分类实现和白名单匹配函数未修改。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-anonymous-auth.mDSG7Vy8/cache、GOPROXY=off，无下载或新增依赖。两个overlay测试源、配置、源码快照、首轮/最终测试及构建/编译输出、两份抽取范例源码/编译产物与独立缓存均放在本轮任务目录，仓库未增加持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-anonymous-auth.mDSG7Vy8/overlay.json -run '^TestAnon' -count=1 -v ./common-pkg/appMiddleware -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-anonymous-auth.mDSG7Vy8/overlay.json -run '^TestAnon' -count=1 -v ./common-pkg/appMiddleware ./common-pkg/app -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-anonymous-auth.mDSG7Vy8/auth_example_1.a /private/tmp/new-platform-anonymous-auth.mDSG7Vy8/auth_example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-anonymous-auth.mDSG7Vy8/auth_example_2.a /private/tmp/new-platform-anonymous-auth.mDSG7Vy8/auth_example_2.go
gofmt -l common-pkg/appMiddleware/apiAuthToken.go common-pkg/appMiddleware/apiHeader.go common-pkg/appMiddleware/checkServiceStatusMiddleware.go common-pkg/appMiddleware/rpcAuthMiddleware.go common-pkg/app/start.go
git -C common-pkg diff --check
git diff --check
```

- **限制与接入影响** HTTP使用实际路由加内存请求，RPC为真实grpc但通过进程内bufconn，无TCP/TLS/部署网关联调；认证回调由临时业务桩提供，不验证生产Token服务/凭证配置。原有外部设施集成测试只编译未执行，本机Go1.26.5，未另装Go1.22。Head.Verify没有业务规则，不能声称已验证新增业务必填限制；预检继续沿用CORS流程。白名单map共享方式保持，调用方不得并发写入。非白名单RPC现在要求正数UID，纯服务任务需规划匿名方法；RPC服务调用方认证和用户身份核实由业务回调落实。API/RPC子项目暂无业务源码，本轮只更新接入记录。
- **测试清理** 已在清理前登记实际命令、结果与限制；随后删除两个overlay测试源及配置、源码快照、首轮/最终测试日志、构建/全量编译输出、两份抽取范例源码与编译产物、全部本轮临时脚本产物和独立缓存。核对任务目录与两个仓库虚拟测试目标均不存在，原有五个测试及摘要不变。白名单匹配函数逐字一致，上一轮本地多语言及范围外公共包原有文件摘要保持不变；新增持久文件仅为AUTH使用说明。API/RPC子项目仅更新开发记录，无新增测试或临时产物。


### 2026-11-22 / DEV-20261122-001 / 内置错误文案与项目本地词典对齐

- **用户确认与范围** 用户确认第32项：清理内置文案多余尾部空格，核对本地示例词典，保留精确匹配、错误码、认证及白名单接口。第31项经用户继续下一项确认；本轮不实施第33项日志身份同步。
- **改动** result/xcode/common_err.go的UserNotFound、TokenInvalid，以及appMiddleware/apiRequestDecrypt.go的RequestDecryptError、RequestBadError公开文案去除尾部空格，分别与User not login.、Token invalid.、Request decryption failed.、Request forbid.词典key一致，补齐中文注释。其余常用固定公开文案与示例词典核对一致，原有内部诊断及业务动态文案不改；语言加载、翻译器、两份词典JSON、错误分类和认证逻辑均保持不变。
- **说明与影响** [本地词典说明](language/README.md)增加四条内置原文/中文对照。项目若使用这些内置错误带空格的旧key，应同步去除key末尾空格；不增加兼容别名或对所有文案自动TrimSpace，业务自定义message及data保持精确匹配。
- **验证结果** 5组TestErrorLocale全部通过：12条常用内置文案的中文/英文与en-US回退；8类错误×4种语言×原始/包装共64个统一HTTP响应场景，验证原业务码、HTTP200、Trace、空data、内部诊断不外泄；5类错误经ToGRPCStatus/FromGRPCError转换再输出中文响应；6类带空格/未命中业务文案原样保留，data字段精确匹配、id保护及9007199254740993精度保留；实际Trace/请求头/认证中间件覆盖受保护未登录、白名单无效/过期Token、权限不足及正常匿名5个场景，错误不进入业务，白名单/回调次数规则保持。
- **完整核对** 全量构建及原有测试全量编译通过，gofmt、diff与说明链接检查通过。按本轮开始SHA-256快照，除两份错误定义及语言说明外其他原有公共包文件不变；五个原有测试、词典JSON、语言实现、白名单/认证/RPC及日志源码摘要不变。没有新增依赖或发布操作。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-error-language.L6cxJ3Bz/cache、GOPROXY=off，无下载。单个overlay测试源、overlay配置、源码摘要快照、测试/构建/编译输出与独立缓存均在本轮任务目录，仓库未新增持久测试文件。

```text
go test -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-error-language.L6cxJ3Bz/overlay.json -run '^TestErrorLocale' -count=1 -v ./common-pkg/result -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
gofmt -l common-pkg/result/xcode/common_err.go common-pkg/appMiddleware/apiRequestDecrypt.go
git -C common-pkg diff --check
git diff --check
```

- **限制** HTTP认证使用内存请求和真实中间件，RPC仅验证错误状态的本地转换，未启动实际RPC或做部署网关联调；解密两类错误通过统一响应入口验证，未测试完整AES请求链路或改变早期中间件的语言注入时机。原有外部设施集成测试只编译未运行；使用本机Go1.26.5，未另装Go1.22。本次只改公开文案，无并发逻辑改动，未重复进行竞态检测。API/RPC子项目暂无业务源码，仅更新接入记录。
- **测试清理** 已在清理前登记实际命令、结果及限制；随后删除本轮overlay测试源及配置、源码摘要快照、测试/构建/编译输出、临时脚本产物和独立构建缓存。核对任务目录及仓库虚拟测试目标均不存在，原有五个测试及摘要不变；未新增任何持久测试或其他文件。范围外原有源码、词典JSON、翻译/认证/白名单/RPC/日志实现摘要一致，API/RPC子项目仅更新开发记录，无新增测试或临时产物。


### 2026-11-22 / DEV-20261122-002 / 认证结果与上下文日志同步

- **用户确认与范围** 第32项经用户继续下一项确认，本轮第33项按已确认方案实施：认证后的业务日志同步用户ID，匿名为0；SetBusinessCode同步当前分支业务标识，按字段覆盖，保留Trace/其他日志字段，Token和用户JSON不自动入日志，公共函数签名和认证规则不变。
- **实现** headInfo.SetTokenUid在复制出站metadata后覆盖LogKeys中的用户ID日志；现有HTTP认证已调用该Setter，无需修改认证代码。SetBusinessCode继续分别复制入站/出站metadata，并同步日志字段。HTTP初始化后的字段继续为字符串，RPC初始化后继续为字符串列表；未初始化日志模式的Setter默认字符串，空值显式覆盖旧值，不改变Setter不验证身份/格式的约定。
- **字段覆盖** ContextHeadInLog与ContextMetadataInLog共用内部覆盖逻辑；同名上下文字段最后写入值生效、保留其他自定义字段，不修改父上下文。核对当前go-zero实现，ContextWithFields只追加且未公开字段key；通过公开ContextWithFields构造仅含日志字段的标记上下文，在构造新快照时合并字段并交回logx保存，不使用反射/unsafe或私有key。其余Value、取消、截止时间委托原上下文。RPC字段列表复制，避免调用方后续修改原metadata影响日志；Token即使出现在LogKeys中也不由metadata入口输出。
- **说明** 新增 [日志上下文使用说明](headInfo/LOG_CONTEXT.md)及可运行示例，更新 [metadata说明](headInfo/METADATA.md)和[认证说明](appMiddleware/AUTH.md)。必须使用返回的新ctx记录后续日志；旧ctx、绑定旧ctx的logger及外层访问日志不自动回写。去重限于合并的上下文字段，业务另加单条/全局同名字段按其自身接口处理。文档同时纠正已不存在的SetBothBusinessCode旧名称引用，未新增兼容接口。
- **验证** 首轮5组、加入真实RPC内存链路后的最终6组TestLogSync及-race全部通过。原生logx JSON和平台Logrus两种writer实际输出验证UID/业务标识更新、无重复上下文字段、Trace及其他业务字段保留；RPC数组格式、原metadata修改后日志快照不变、空值清除；64次连续更新及自定义字段交错不持续累积同名字段；直接Setter默认格式及原始前导零保留；LogKeys排除字段不被Setter扩大，Token/扩展JSON不自动输出；外部context值（含同类型日志字段切片）、截止时间、取消和metadata父/子副本保持。16goroutine各32次共512个日志分支的UID/业务标识与Trace准确，无竞态及父上下文污染。
- **链路验证** 实际HTTP Trace→单值头→头上下文→认证链，随后调用真实grpc.ServerOption/ClientInterceptor，通过bufconn验证4个成功请求（匿名/匿名Token/登录公开/登录受保护）及3个拒绝场景（受保护未登录、无效Token、无权限）。API日志为字符串，RPC为字符串列表；RPC子分支修改业务标识不影响RPC父级和API调用方，统一Trace及rpc-name保留，用户JSON/令牌不进入自动字段。认证和白名单行为未改变。
- **完整检查与范例** 全量构建、原有测试全量编译、gofmt、diff及说明链接核对通过。说明中的main抽取编译并实际运行，解析三条JSON日志，验证UID均为12、业务标识依次catalog/child/catalog且Trace一致。对照本轮开始SHA-256快照，五个原有测试、上一轮文案/词典调整及其他范围外源码保持不变。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-context-logs.CYIYrgAK/cache、GOPROXY=off，无下载或新增依赖。两个overlay测试源、配置、源码快照、首轮/最终测试及构建/编译输出、抽取范例源码/二进制/运行输出与独立缓存均位于本轮任务目录；无持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-context-logs.CYIYrgAK/overlay.json -run '^TestLogSync' -count=1 -v ./common-pkg/headInfo -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-context-logs.CYIYrgAK/example /private/tmp/new-platform-context-logs.CYIYrgAK/example.go
/private/tmp/new-platform-context-logs.CYIYrgAK/example
gofmt -l common-pkg/headInfo/logKeys.go common-pkg/headInfo/head.go common-pkg/headInfo/headInfo.go
git -C common-pkg diff --check
git diff --check
```

- **限制** HTTP使用内存请求，真实gRPC走进程内bufconn，无TCP/TLS/部署日志平台联调；原有外部设施集成测试只编译未运行。本机Go1.26.5，未另装Go1.22；日志适配按当前go-zero v1.7.4公开ContextWithFields行为验证，框架升级需复核覆盖行为。Logger/旧ctx不自动刷新，不解析日志正文，手工改metadata仍需显式同步。LogKeys仅初始化时配置，请求期间不得并发修改。API/RPC子项目暂无业务源码，仅更新接入约定。
- **测试清理** 已在清理前登记实际命令、结果与限制；随后删除两个overlay测试源及配置、源码摘要快照、首轮/最终测试和构建/编译输出、抽取范例源码及二进制/运行输出、全部本轮临时脚本产物与独立缓存。核对任务目录与两个仓库虚拟测试目标均不存在，原有五个测试及摘要不变；上一轮错误文案/词典变更及其他范围外原有文件摘要一致。新增持久文件仅为LOG_CONTEXT使用说明，API/RPC子项目仅更新记录，无新增测试或临时产物。


### 2026-11-22 / DEV-20261122-003 / 认证前错误响应语言统一

- **用户确认与范围** 用户确认第34项：入口单值头校验通过后提前注入语言，请求头、正文、并发、AES错误统一从上下文读取；重复语言头拒绝且不任选，未知语言/词条回退原文。中间件顺序、AES配置开关、错误码和Trace规则保持。第33项经用户继续下一项确认，本轮保留其日志同步实现。
- **实现** headInfo/httpHeaders.go新增ContextWithRequestLanguage，将HTTP唯一Content-Language写入出站metadata副本，仅裁剪与原Head一致的首尾空格，其他字段/入站metadata/Trace/取消/截止时间保持。缺失、空值、重复、大小写别名和空值列表选用空语言，清除继承的旧语言，不任选首个值或猜测语言。该函数只选择语言，不代替完整NormalizeHTTPHeaders。
- **入口与统一读取** RequestHeaderMiddleware成功校验后对下游新ctx提前注入语言，发生于原生声明长度检查/并发/正文/AES之前；错误分支也通过同一函数选择唯一语言。其他头非法但语言唯一时400可翻译，语言头自身重复/别名/无值列表时400使用原文。正文错误移除原始头回退，限流错误改读GetContentLanguage(ctx)，AES和统一响应原有ctx读取即可取得语言，无需修改解密器或响应实现。原请求/请求头/正文不被语言函数修改，不注入身份或加密器。
- **说明与接入影响** 更新[HTTP入口](headInfo/HTTP_HEADERS.md)、[正文限制](app/REQUEST_BODY.md)、[多语言说明及认证前解密错误范例](language/README.md)。标准NewServer自动接入；独立net/http需显式接入RequestHeaderMiddleware，或在完成单值头校验后注入语言。仅调用NormalizeHTTPHeaders不会写ctx，单独正文/并发/AES组件若无上下文语言则返回原文，不再从未验证原始头自行选择；后续头变化也不覆盖已注入ctx。RPC语言Getter/协议保持，Trace自身和原生组件自行生成的错误不在本轮改写范围。
- **验证结果** 首轮4组头部/上下文测试通过，完整API加入后8组及-race通过，另外定向增加1组缺少词条回退验证通过，共9组，无失败或竞态报告。15种语言头结构覆盖缺失/空/空白、大小写/下划线、未知/逗号原文、重复同值/不同值、空值列表、大小写别名含空别名、旧下划线名称不读取，父上下文旧语言被覆盖但本身不变；其他metadata/自定义多值/入站值/Trace/取消/截止时间保持、不读或关闭正文。头校验拒绝路径覆盖有效单语言翻译与重复语言不沿用旧值；512个并发请求分支语言不串用。
- **实际API验证** 使用app.NewServer/ApiUseMiddleware及真实go-zero路由（内存请求，不监听网络）：5种语言×6种AES请求共30场景，涵盖非法密文、错误密钥、空密文、缺少aes_data、非法封装及正常加解密；确认解密失败在认证前按语言翻译，HTTP200/code400保持。成功响应验证翻译、实际语言、9007199254740993精度及AES密文解码。正文/普通响应4语言×4路径共16场景，覆盖原生声明长度413、未知长度实际超限413、读取失败400与正常明文返回，错误不进入认证/业务且无内部诊断。另验证3种实际路由重复头拒绝与1组空词典缺词条原文回退。
- **并发上限与错误源** 在同一实际服务的业务handler内保持名额占用并嵌套请求，4种语言均得到503对应译文/原文，拒绝时未读取正文或执行认证；CORS预检204、重复语言先400、释放后正常请求可用。单独组件验证已注入zh-CN后即使原始头改成en，正文/限流仍读ctx；没有注入时不再回退未校验头。此组通过嵌套请求确定占用时机，不声称是多连接并发压测。
- **范例与完整核对** 两份language说明Go范例抽取编译并实际运行：原示例中文/精确数字保持，新示例认证前AES失败输出中文、HTTP200/code400且Trace一致。全量构建、原有测试全量编译、gofmt、diff及说明链接核对通过。SHA-256对照本轮开始快照，五个原有测试、上一轮日志字段实现、AES/响应/认证/Trace/中间件装配顺序/错误码和词典JSON摘要不变。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-early-language.G6C4mqYa/cache、GOPROXY=off，无下载或新增依赖。两个overlay测试源、配置、源码快照、首轮/完整/缺词条测试输出、构建/编译输出、两份抽取范例源码/二进制/运行输出、词典临时副本与独立缓存都在本轮任务目录；缺词条测试夹具由t.TempDir自动清理。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-early-language.G6C4mqYa/overlay.json -run '^TestEarlyLanguage' -count=1 -v ./common-pkg/headInfo -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-early-language.G6C4mqYa/overlay.json -run '^TestEarlyLanguage' -count=1 -v ./common-pkg/headInfo ./common-pkg/app -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-early-language.G6C4mqYa/overlay.json -run '^TestEarlyLanguageMissingDictionaryKey$' -count=1 -v ./common-pkg/app -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-early-language.G6C4mqYa/example_1 /private/tmp/new-platform-early-language.G6C4mqYa/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-early-language.G6C4mqYa/example_2 /private/tmp/new-platform-early-language.G6C4mqYa/example_2.go
# 两个范例均在任务目录运行，该目录已复制示例etc/language
./example_1
./example_2
gofmt -l common-pkg/headInfo/httpHeaders.go common-pkg/appMiddleware/requestHeaderMiddleware.go common-pkg/appMiddleware/requestBodyMiddleware.go common-pkg/appMiddleware/concurrencyLimiterMiddleware.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 实际API使用内存HTTP，无TCP/TLS/网关或生产词典挂载联调；原有外部设施集成测试只编译未运行，未另装Go1.22，本机Go1.26.5。为确定各拒绝路径，临时API配置关闭可选日志/统计/熔断/负载丢弃/超时/原生MaxConns，Trace/原生MaxBytes等保持启用；未验证原生组件自行生成错误的译文。RPC未改，仅验证metadata保留，不另启RPC服务。API/RPC子项目暂无业务源码，本轮只更新接入约定。
- **测试清理** 已在清理前登记实际命令、结果及限制；随后删除两个overlay测试源及配置、源码摘要快照、首轮/完整/缺词条测试输出、构建/编译输出、两份抽取范例源码及二进制/运行输出、词典临时副本、全部本轮脚本产物和独立构建缓存。核对任务目录与两个仓库虚拟测试目标均不存在，原有五个测试及摘要不变；无新增持久文件，范围外原有源码、上轮日志同步、AES/认证/响应/Trace/顺序/错误码及词典JSON摘要一致。API/RPC子项目仅更新记录，无新增测试或临时产物。


### 2026-11-22 / DEV-20261122-004 / 按用户要求撤回第35项

- **决定** 用户明确要求撤回第35项API panic恢复改动，停止该项开发，不作为已完成功能交付。恢复本轮修改的app/start.go、app/traceServer.go、appMiddleware/apiHeader.go及result/responseWriter.go，删除本轮新增recovery.go、recovery_writer.go、recovery_interfaces.go；第34项及之前内容保留，不改认证/白名单/语言/Trace等现有规则。
- **已执行验证** 撤回前初次构建通过，第一轮8组及-race通过；扩展到实际API后，英文正常成功响应夹具误期望原业务key而实际为Enabled，导致该测试组的两个子场景失败，修正临时预期后最终13组及-race通过。覆盖提交前/提交后panic、序列化/AES/认证/正文异常、ErrAbortHandler、原生超时包装、64种writer接口组合、FlushError/ResponseController、net.Pipe接管连接、超时后异常和256并发请求。上述测试仅说明已尝试实现的验证状态，相关实现与测试均按用户要求撤回，不代表当前基线新增这些行为。
- **实际命令** 平台根目录使用GOCACHE=/private/tmp/new-platform-api-recovery.UD8lBECW/cache、GOPROXY=off执行，无新增依赖或下载；首轮、扩展轮及修正后完整验证均保存在任务目录。初次全量构建后有恢复逻辑调整，未执行最终全量构建/原有测试编译，不将初次构建结果冒充最终检查。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/... ./language-pkg/...
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-api-recovery.UD8lBECW/overlay.json -run '^TestRecovery' -count=1 -v ./common-pkg/appMiddleware -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-api-recovery.UD8lBECW/overlay.json -run '^TestRecovery' -count=1 -v ./common-pkg/appMiddleware ./common-pkg/app -timeout=60s
git -C common-pkg diff --check
```

- **限制与清理前登记** 验证使用内存HTTP、go-zero原生包装、net.Pipe和模拟writer，未进行TCP/TLS/部署联调；本机Go1.26.5，原有外部集成测试未运行。四个overlay测试源、配置、源码摘要快照、测试与构建输出、独立缓存及本轮生成产物均须删除。回退前已确认四个原有文件的HEAD内容摘要与本轮开始快照一致，接下来逐文件恢复并核对全量快照，清理后追加核对结果。

- **撤回与清理结果** 已逐文件恢复四个原有源码文件，并删除三个本轮新增实现文件；对照本轮开始SHA-256快照，除本撤回记录外，所有原有公共包文件内容完全一致，第34项及之前实现保持。清理前登记后，删除四个overlay测试源及配置、全部构建/测试输出、源码快照、生成脚本产物与独立缓存；核对任务目录、四个仓库虚拟测试目标及三个新增实现均不存在，原有五个测试及摘要保留。公共包仅开发记录存在差异，无第35项代码残留。


### 2026-11-22 / DEV-20261122-005 / JWT有效期校验统一

- **用户确认与范围** 第36项按已确认方案开发：普通ParseToken及ParseRefreshToken都要求exp存在且严格晚于当前校验秒，缺失/无效/已到期/到期当秒均拒绝。默认7天签发、HS256、密钥要求、Claims字段、现有接口和匿名白名单规则保持。第34项经后续开发确认，第35项继续保持撤回，不恢复任何相关代码。
- **实现** 仅修改jwtToken/jwt.go：签名及原有Claims校验通过后，ParseToken要求ExpiresAt>0并且ExpiresAt>jwt.TimeFunc().Unix()，失败返回nil Claims及错误。沿用底层解析器同一时间来源（默认time.Now），不增加宽限。ParseRefreshToken直接复用ParseToken，再执行已有AdminToken/UserToken类型检查，删除原独立time.Now有效期判断。验证签名过程中跨过到期秒也拒绝；不扩展普通TokenType/UID/issuer/audience校验或新增配置。
- **使用说明** 新增 [JWT说明与范例](jwtToken/README.md)，说明必填exp、到期边界、刷新只校验不自动签发、Claims与既有规则保留。旧无exp或到期当秒仍被普通解析接受的Token现在会返回错误，不增加兼容分支。API认证回调应返回该error，不得把失败作为0,nil匿名结果。
- **定向验证** 6组TestJWTExpiry及-race全部通过。普通/刷新两种解析覆盖缺失、null、0、负数、过期、等于当前秒、未来1秒、MaxInt64及字符串/小数/布尔/对象错误类型；到期前1ns、到期时、到期后1ns/1s及验签过程中跨越到期边界，无sleep的固定时钟验证。错误密钥、篡改签名、HS384、none、未来nbf/iat、短密钥/全空白密钥继续拒绝，失败不返回Claims。GenerateToken的两种类型保持默认7天、HS256、原业务字段/JSON标签及Issuer；普通未知TokenType策略及刷新类型限制保持原样。
- **并发与API联动** 稳定只读测试时钟下16goroutine各32轮、每轮两种解析共1024次调用，正常/到期结果正确，无竞态。实际Trace→单值头/语言→请求头上下文→认证中间件覆盖两路径×六种Token共12场景：无Token仅白名单匿名，缺失/0/到期/已过期exp即使白名单也返回原业务402，受保护无Token401，正常Token注入原UID，Trace和HTTP200业务协议保持。使用已签名测试Token，不输出Token或密钥。
- **完整核对与范例** 当前go.work仅包含common-pkg，因此全量构建/原有测试编译按./common-pkg/...执行；未恢复用户已移除的旧language-pkg工作区引用。README范例抽取编译并通过JWT_SECRET注入临时测试值实际运行，输出12 true，不打印密钥/Token。gofmt、diff及说明链接通过；SHA-256对照本轮开始快照，原有五个测试、用户已有go.mod/go.sum改动、go.work/go.work.sum及范围外公共包文件保持不变。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-jwt-expiry.vtFJbArJ/cache、GOPROXY=off，无下载或新增依赖。两个overlay测试源、配置、原始文件/工作区摘要快照、测试/构建/编译输出、抽取范例源码/二进制/运行输出及独立缓存均在本轮任务目录；仓库无新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-jwt-expiry.vtFJbArJ/overlay.json -run '^TestJWTExpiry' -count=1 -v ./common-pkg/jwtToken -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-jwt-expiry.vtFJbArJ/example /private/tmp/new-platform-jwt-expiry.vtFJbArJ/example.go
# 通过JWT_SECRET注入本轮临时32字节测试值，不记录或输出值
/private/tmp/new-platform-jwt-expiry.vtFJbArJ/example
gofmt -l common-pkg/jwtToken/jwt.go
git -C common-pkg diff --check
git diff --check
```

- **限制** JWT测试通过真实HS256签名及内存HTTP，未进行部署时钟/NTP、生产密钥或真实API/RPC联调；其他业务子项目未修改或测试。原有外部设施集成测试仅编译未执行，本机Go1.26.5，未另装Go1.22。测试临时替换上游jwt.TimeFunc后恢复，替换期间不并发写时钟；生产未新增时钟hook或可变全局配置。本次不调整令牌撤销、续签策略、issuer/audience或最大可接受时长。
- **测试清理** 已在清理前登记实际命令、结果与限制；随后删除两个overlay测试源及配置、原始/工作区摘要快照、测试/构建/编译输出、抽取范例源码及二进制/运行输出、全部本轮临时脚本产物和独立构建缓存。核对任务目录及两个仓库虚拟测试目标均不存在，原有五个测试及摘要不变；范围外原有源码、用户已有go.mod/go.sum整理、go.work/go.work.sum和第35项撤回状态保持。新增持久文件仅为JWT使用说明；API/RPC目录仅更新记录，其他业务子项目未新增或删除测试/产物。


### 2026-11-22 / DEV-20261122-006 / GenerateToken结构体参数

- **用户确认与基线** 用户确认将GenerateToken参数改为GenerateTokenParams结构体，ExpireDuration使用time.Duration，直接替换旧签名，保留当前自定义时长/Claims及解析规则。开始时用户已自行改为六个位置参数（密钥/UID/Source/Sign/BusinessCode/纳秒int64），移除原固定7天/角色参数/刷新类型限制；本轮按这一最新基线调整，不恢复旧设计。第36项完成确认，第35项继续撤回。
- **实现** 新增按值传入的GenerateTokenParams，包含JWTSecret、Uid、Source、Sign、BusinessCode、ExpireDuration，字段和公开接口均有中文用途/单位/零值注释。签发读取结构体字段，time.Now().Add直接接受time.Duration，返回类型、HS256、密钥验证、Claims标签、issuer及当前TokenType默认内容保持不变；不增加UID/业务字段或时长正数校验，无默认有效期，零/负时长沿用签发即过期行为，exp精度仍为Unix秒。旧多参数函数不保留兼容入口。
- **并行工作区变化** 开发过程中检测到ParseRefreshToken已从工作区另行移除（非本项删除），本轮保留此修改，不恢复函数。初轮临时测试因此报undefined: platform.ParseRefreshToken；核对当前源码后，仅调整测试和说明为现有GenerateToken/ParseToken接口，最终通过。当前README也同步清理之前已过时的固定7天、角色/类型参数和刷新专用接口示例，历史记录保留。
- **使用说明与调用方** 更新 [JWT参数及范例](jwtToken/README.md)，使用ExpireDuration:24*time.Hour，原int64秒配置应明确乘time.Second。当前工作区Go源码中检索到GenerateToken的引用只有函数定义，没有其他业务Go调用点需要迁移；外部调用方需迁移为结构体形式。
- **验证** 最终5组TestJWTParams及-race全部通过。3秒/15分钟/24小时/7天真实签发并解析，验证字段值/Unicode及空白保留、HS256签名、exp单位和JSON标签，不将参数密钥/时长写入Claims，不恢复role_id参数。覆盖UID=0及可选业务字段为空、短/空/全空白密钥拒绝、密钥首尾空白不自动裁剪；零/负/极小值/极大值time.Duration与原处理一致，小于一秒按实际Unix秒截断核对。ParseToken继续拒绝缺失/0/过去/相等exp，未来exp通过且当前类型策略保持。16goroutine各16次共256个结构体参数签发/解析，UID/业务标识不串用，无竞态。
- **范例与完整检查** README新范例抽取编译，通过JWT_SECRET注入临时测试值运行输出12 catalog，不输出密钥/Token。按当前单模块go.work执行全量构建、原有测试全量编译，gofmt、diff和说明链接检查通过。SHA-256核对五个原有测试、用户已有go.mod/go.sum整理、go.work/go.work.sum及范围外原有公共包文件保持不变；其他业务子项目源码和测试均未改动。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-jwt-params.TuFCiygj/cache、GOPROXY=off，无新增依赖或下载。单个overlay测试源/配置、原始文件及工作区摘要快照、初次失败/最终通过输出、构建/原有测试编译输出、抽取范例源码/二进制/运行输出和独立缓存均在本轮任务目录，仓库无新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-jwt-params.TuFCiygj/overlay.json -run '^TestJWTParams' -count=1 -v ./common-pkg/jwtToken -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-jwt-params.TuFCiygj/example /private/tmp/new-platform-jwt-params.TuFCiygj/example.go
# 运行时通过JWT_SECRET注入本轮临时测试值，不记录或输出值
/private/tmp/new-platform-jwt-params.TuFCiygj/example
gofmt -l common-pkg/jwtToken/jwt.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 使用真实JWT签名与本地解析，未进行生产密钥、API/RPC或部署联调；其他业务模块未编译或测试。原有外部设施集成测试只编译未执行，本机Go1.26.5，未另装Go1.22。未加入时长默认值、正数校验或新的业务认证规则；调用方应显式提供有效时长。ParseRefreshToken的移除来自并行工作区变化，本轮只按最新接口验证，不把该移除归入结构体重构。
- **测试清理** 已在清理前登记实际命令、结果、接口变动导致的测试编译失败及限制；随后删除overlay测试源及配置、原始/工作区摘要快照、初次失败与最终通过输出、两次构建/原有测试编译输出、抽取范例源码及二进制/运行输出、全部本轮临时脚本产物和独立构建缓存。核对任务目录及仓库虚拟测试目标均不存在，原有五个测试及摘要不变；未新增持久文件。范围外原有源码、用户已有go.mod/go.sum整理及go.work/go.work.sum保持，另行移除的ParseRefreshToken未恢复，其他业务项目测试和文件未改动。


### 2026-11-22 / DEV-20261122-007 / HeadInMetadata保留自定义字段

- **用户确认与范围** 第38项按确认实施：HeadInMetadata复制已有出站metadata，在副本覆盖当前Head字段（空值也覆盖），保留自定义字段及多值，W3C按活动Trace重新注入，父上下文不修改。第37项经用户继续下一项确认，第35项继续撤回。仅修改headInfo/head.go及当前使用说明，不调整认证、白名单、JWT、语言、日志实现或RPC拦截器。
- **实现** 先构造原有15个Head平台字段，再复制metadata.FromOutgoingContext的结果，逐字段Set覆盖为单值，杜绝旧重复值、旧UID/角色/令牌通过合并继承。其他出站字段和值顺序保留，包含AppendToOutgoingContext的追加内容；不从入站metadata补字段。先清除旧出站traceparent/tracestate，再调用既有traceMetadata按活动Trace注入；无活动Trace时不生成ID或采用Head.Trace，保持原重建metadata时不会继承旧W3C字段的语义。rpc-name不在Head字段中，保留至客户端拦截器按目标覆盖。
- **说明** 更新 [metadata合并规则及示例](headInfo/METADATA.md)和 [HTTP入口](headInfo/HTTP_HEADERS.md)，明确仅保留显式写入出站metadata的自定义字段，不自动转发HTTP头；用户扩展JSON不自动进入metadata，敏感进程内数据不应放在出站metadata。函数签名及现有头名不变，不新增别名或未知字段校验，不扩大日志名单。
- **验证** 首轮5组及-race通过，补充真实RPC内存链路后最终6组TestMetadataMerge及-race全部通过。15字段各用非空/全空Head覆盖旧重复值，检查自定义多值/空字符串/二进制/rpc-name保留，ValidateMetadata可通过；已有AppendToOutgoingContext多段自定义值顺序保持；无出站/nil/空metadata正常初始化。入站独立、父级metadata/切片不污染、额外context值/取消/截止时间保留。W3C测试覆盖有效活动Trace/Span/TraceState覆盖旧链路与伪Head.Trace、没有活动Trace时清除旧W3C字段不生成ID。16goroutine各32次共512个分支的UID/业务标识、自定义多值和Trace隔离，无竞态。
- **实际HTTP→RPC** 使用Trace→单值头/语言→ApiHeader→认证组件，经真实grpc.ServerOption及ClientInterceptor通过bufconn调用服务，验证匿名与登录请求。API上游metadata含自定义多值和追加值、旧UID/角色/令牌及业务标识；HTTP伪UID/角色不采用，旧身份清空后仅认证回调可写入12，匿名为0；自定义字段完整到达调用方认证与RPC业务，语言/业务/IP/统一Trace/rpc-name正确，父上下文不变。未把自定义metadata加入响应外的自动日志字段。
- **范例及完整检查** METADATA.md三份Go范例抽取编译通过，两份main实际运行：新示例输出[one two] 0 catalog / 999及统一Trace，既有示例12 192.0.2.1 / 0 true / 12保持。全量构建、原有测试编译、gofmt、diff及说明链接检查通过。SHA-256对照本轮开始快照，五个原有测试、JWT参数重构、用户已有依赖/工作区、Trace/认证/白名单/日志/语言及其他范围外公共包源码不变。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-metadata-merge.BIcDXLg3/cache、GOPROXY=off，无下载或新增依赖。两个overlay测试源及配置、原始文件/工作区摘要快照、首轮/最终测试及构建/编译输出、三份范例源码/编译产物、两份运行输出和独立缓存均在任务目录；仓库无新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-metadata-merge.BIcDXLg3/overlay.json -run '^TestMetadataMerge' -count=1 -v ./common-pkg/headInfo -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-metadata-merge.BIcDXLg3/example_1 /private/tmp/new-platform-metadata-merge.BIcDXLg3/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-metadata-merge.BIcDXLg3/example_2 /private/tmp/new-platform-metadata-merge.BIcDXLg3/example_2.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-metadata-merge.BIcDXLg3/example_3.a /private/tmp/new-platform-metadata-merge.BIcDXLg3/example_3.go
/private/tmp/new-platform-metadata-merge.BIcDXLg3/example_1
/private/tmp/new-platform-metadata-merge.BIcDXLg3/example_2
gofmt -l common-pkg/headInfo/head.go
git -C common-pkg diff --check
git diff --check
```

- **限制** HTTP为内存请求，真实grpc通过进程内bufconn，未做TCP/TLS/网关或部署认证联调；自定义metadata仍受gRPC既有键值/大小限制，合并不证明字段真实可信，不应承载仅进程内使用的敏感值。原有外部设施集成测试只编译未运行，本机Go1.26.5，未另装Go1.22；其他业务模块未修改或测试。API/RPC目录仅更新接入记录。
- **测试清理** 已在清理前登记实际命令、结果与限制；随后删除两个overlay测试源及配置、原始/工作区摘要快照、首轮/最终测试与构建/编译输出、三份抽取范例源码及编译产物、两份运行输出、全部本轮临时脚本产物和独立构建缓存。核对任务目录及两个仓库虚拟测试目标均不存在，原有五个测试及摘要不变；未新增持久文件。范围外原有源码、JWT参数、已有依赖及工作区配置保持；API/RPC目录仅更新记录，其他业务项目文件和测试未改动。


### 2026-11-22 / DEV-20261122-008 / GetFullHead深复制请求头

- **用户确认与范围** 用户确认第39项并明确“此项不做测试”。第38项经用户继续下一项确认。仅调整headInfo/head.go中的GetFullHead，不改签名、请求头字段/值、认证/Trace/metadata装配逻辑。
- **改动** 使用标准库http.Header.Clone深复制map与值切片，保留字段名、值及多值顺序；原Header为nil时仍返回可写的非nil空map，维持原调用行为。补充中文注释说明副本隔离、非nil请求条件及凭证不自动脱敏，未扩大日志或自动转发范围。
- **检查与限制** 按用户要求未新增或运行测试、测试编译、构建或运行时验证；只阅读本机标准库Header.Clone实现，并执行gofmt -w common-pkg/headInfo/head.go、git -C common-pkg diff --check及git diff --check进行格式/差异核对。未声称运行验证双向修改行为。API/RPC目录仅更新接入记录，其他业务子项目不改动。
- **测试清理** 本轮未创建测试文件、测试脚本、缓存目录或临时产物，无新增测试内容需要清理；原有测试和用户文件保留。


### 2026-11-22 / DEV-20261122-009 / CORS凭证模式与实例配置

- **用户确认与范围** 第40项按确认实施：默认任意来源、关闭跨域凭证；需要凭证时显式开启并配置精确来源，返回匹配Origin及Vary:Origin；保留开放所有预检请求头，正常和早期错误采用同一策略。第39项确认完成且仍按该项要求不增加专项测试，本轮保留GetFullHead改动，第35项继续撤回。
- **配置** 新增appMiddleware.CORSConfig（AllowCredentials/AllowedOrigins）及ConfigAppServer.CORS备注，NewCorsMiddlewareWithConfig校验并复制配置。凭证开启必须有名单，禁止*、null、空/空白/控制字符、路径/查询/片段/用户信息、通配子域及非法端口；使用HTTP/HTTPS精确序列化来源（大小写/显式端口精确匹配）。非凭证模式也可配置来源名单，否则默认*。名单外/缺少/重复/大小写别名Origin不授予跨域头，普通请求仍按原业务认证处理，不新增403或改变白名单。
- **服务装配** NewServer保持原调用方式并使用默认策略；新增NewServerWithCORS(restConf,config,options...)显式应用实例配置。自定义router新增WithRequestTraceCORSRouter(...)(rest.RunOption,error)，原WithRequestTraceRouter保持默认入口。外层Bind绑定实例只读策略并写头，覆盖Trace/单值头错误、路由404/405和原生组件错误；标准入口在Trace/单值头/声明长度检查后处理有效OPTIONS并返回204（含未注册路径），不进入正文/认证/业务。ApiUseMiddleware原默认CORS继续沿用已绑定策略，无进程全局设置，实例互不影响。
- **错误出口与头** 原Trace/头/正文错误调用的setCORSHeaders从ctx取统一策略，不回退成开放来源，清理旧授权头后按策略重建。凭证匹配时输出具体Origin与true，默认不再输出Allow-Credentials。Vary保留已有项并大小写不敏感去重，Vary:*不追加；预检继续回显全部Access-Control-Request-Headers（包括Authorization和项目自定义头），方法列表及请求ID/项目令牌暴露头保持。CORS不代替服务认证/CSRF或Cookie自身规则，不叠加其他CORS/裸router覆盖配置。
- **说明与范例** 新增 [CORS配置与接入](app/CORS.md)，含YAML、API/独立net-http例子、默认和凭证模式、迁移/上下文边界，引用Fetch官方规范的凭证规则。更新HTTP_HEADERS.md及装配注释。仅填写配置字段不会自动接入，业务须传给NewServerWithCORS；独立HTTP使用Bind在Trace外层写入策略，Handle在校验后处理预检。旧凭证客户端需显式配置来源，默认Token客户端不必开启credentials。
- **验证** 首轮4组、完整7组TestCORS及-race全部通过，另1组原生超时验证通过，合计8组无竞态/失败。配置校验覆盖18类非法来源和4类有效来源、空凭证名单拒绝及输入切片副本；默认/凭证/非凭证名单3模式×GET/OPTIONS×4来源共24场景，验证授权字段清理、请求头全放行、暴露头、Origin精确/重复/别名不任选、Vary去重和原值保持。16goroutine各32次共512次请求跨两策略实例隔离，父ctx值/截止时间保持。
- **实际API验证** 使用真实app.NewServerWithCORS及go-zero内存路由，2模式×3来源×12路径共72场景：正常、Trace错误、重复业务头、声明/实际正文413、未登录401、404、405、预检/未知路径预检、非法/超长预检，CORS策略一致且Trace/业务码不变。另验证AES失败、嵌套占用并发名额时503及来源拒绝时不错误开放。YAML独立配置、初始化错误早于服务创建、自定义router及默认/凭证两个Server互不影响均通过。原生Timeout开启的独立场景验证503仍保留Origin/true及Vary，不改原超时正文。
- **完整检查与范例** 两份Go范例抽取编译，独立预检main实际输出204/https://app.example.com/true/Authorization-Jwt, X-Custom。全量构建、原有测试全量编译、gofmt、diff及说明链接通过。对照本轮开始SHA-256快照，五个原有测试、第39项GetFullHead、JWT/Trace实现、认证/白名单/语言/日志、已撤回第35项及其他范围外源码保持。
- **实际命令** 平台根目录执行，GOCACHE=/private/tmp/new-platform-cors-config.PcAw2jTg/cache、GOPROXY=off，无新增依赖或下载。两个overlay测试源/配置、源码摘要快照、各轮测试/构建/编译输出、两份抽取范例源码及编译产物、范例运行输出和独立缓存都在本轮任务目录，仓库无新增持久测试文件。

```text
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-cors-config.PcAw2jTg/overlay.json -run '^TestCORS' -count=1 -v ./common-pkg/appMiddleware -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-cors-config.PcAw2jTg/overlay.json -run '^TestCORS' -count=1 -v ./common-pkg/appMiddleware ./common-pkg/app -timeout=60s
go test -race -mod=readonly -buildvcs=false -overlay=/private/tmp/new-platform-cors-config.PcAw2jTg/overlay.json -run '^TestCORSNativeTimeoutPolicy$' -count=1 -v ./common-pkg/app -timeout=60s
go build -mod=readonly -buildvcs=false ./common-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-cors-config.PcAw2jTg/example_1.a /private/tmp/new-platform-cors-config.PcAw2jTg/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-cors-config.PcAw2jTg/example_2 /private/tmp/new-platform-cors-config.PcAw2jTg/example_2.go
/private/tmp/new-platform-cors-config.PcAw2jTg/example_2
gofmt -l common-pkg/appMiddleware/corsConfig.go common-pkg/appMiddleware/corsMiddleware.go common-pkg/app/traceServer.go common-pkg/app/appServer.go common-pkg/app/start.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 请求通过内存HTTP/真实go-zero路由验证，不监听网络，未做真实浏览器credentials/Cookie、TLS/网关/代理缓存或生产域名联调；浏览器标准依据Fetch规范，不能将响应头测试声称为浏览器联调。完整API为确定路径关闭可选访问日志/统计/熔断/丢弃/原生MaxConns/超时等，另独立验证开启原生超时场景。原有外部设施集成测试只编译未运行，本机Go1.26.5，未另装Go1.22。默认模式行为调整：不再输出凭证头；来源匹配是精确字符串，部署须使用浏览器实际Origin。API/RPC目录仅更新接入记录，其他业务项目未改动。
- **测试清理** 已在清理前登记实际命令、结果及限制；随后删除两个overlay测试源及配置、源码摘要快照、各轮测试/构建/编译输出、两份抽取范例源码及编译产物/运行输出、全部本轮临时脚本产物和独立构建缓存。核对任务目录及两个仓库虚拟测试目标均不存在，原有五个测试及摘要不变；第39项GetFullHead及其他范围外源码保持。新增持久文件仅为corsConfig.go与CORS使用说明；API/RPC目录仅更新记录，无新增测试或临时产物。


### 2026-11-22 / DEV-20261122-010 / CORS公开入口简化

- **确认与范围** 按用户“保留建议的，删除额外的”指示，保留推荐的 `NewServerWithCORS`，删除公开的 `WithRequestTraceCORSRouter`，不保留别名；原装配函数改为包内私有复用。已有 `NewServer`、默认策略的 `WithRequestTraceRouter` 及独立 HTTP 中间件保持。调用被删除入口的代码需改用 `NewServerWithCORS(restConf, config.CORS)`。
- **实现与说明** 同步配置字段中文注释、CORS接入说明和HTTP头说明；CORS来源校验、凭证、预检和统一请求ID处理不变。仅修改公共包两份Go文件和相关文档，API/RPC目录仅更新开发记录。
- **验证** 平台根目录使用本机Go、`GOCACHE=/private/tmp/new-platform-cors-entry.uYV9L56q/cache`、`GOPROXY=off`，以下命令均通过；两段Go范例从当前CORS.md抽取，第二段运行输出为204、https://app.example.com、true、Authorization-Jwt, X-Custom。当前源码和使用文档已无被删除的公开符号，历史开发记录保留原始事实。

```text
go build -mod=readonly -buildvcs=false ./common-pkg/...
go test -mod=readonly -buildvcs=false ./common-pkg/... -run '^$' -count=1
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-cors-entry.uYV9L56q/example_1.a /private/tmp/new-platform-cors-entry.uYV9L56q/example_1.go
go build -mod=readonly -buildvcs=false -o /private/tmp/new-platform-cors-entry.uYV9L56q/example_2 /private/tmp/new-platform-cors-entry.uYV9L56q/example_2.go
/private/tmp/new-platform-cors-entry.uYV9L56q/example_2
gofmt -l common-pkg/app/traceServer.go common-pkg/app/appServer.go
git -C common-pkg diff --check
git diff --check
```

- **限制** 本轮仅收敛公开接口，未新增测试用例；原有测试仅编译未运行，未重新执行上一轮完整行为/竞态测试或浏览器联调。范例通过内存HTTP请求验证响应头，无新增依赖、业务项目改动或网络监听。
- **清理完成** 已在清理前登记实际命令、结果及限制，随后删除本轮抽取范例源码、编译产物、源码摘要和独立缓存，任务目录已不存在。核对原有五个测试及摘要保持，无新增测试；公共包仅本条记录、两份Go文件及两份使用说明发生本轮改动，第39项GetFullHead等范围外文件保持。


### 2026-11-22 / DEV-20261122-011 / 删除重复Server构造入口

- **用户确认与范围** 按用户要求保留推荐的NewServerWithCORS，删除NewServer，不增加兼容包装或别名。CORS配置、路由装配和中间件行为保持；原有WithRequestTraceRouter不属于本次两个构造函数的精简范围。
- **修改** 将入口校验、默认策略和正文限制的中文注释集中到NewServerWithCORS；同步中间件注释及CORS、请求体、Trace、HTTP头、认证和多语言使用说明。仓库源码中没有需要迁移的app.NewServer业务调用，现有文档范例改为传入CORS配置。对外原调用需改为NewServerWithCORS(restConf, corsConfig)，零值配置允许任意来源且关闭凭证。
- **验证与限制** 用户明确“修改不验证测试”，本轮未执行构建、编译、测试、范例运行、格式或差异检查；仅为实施修改读取和定位相关源码、注释及文档，不将前轮结果作为本轮验证结论。
- **测试与产物** 本轮未新增测试用例、测试脚本、缓存或临时产物，无需清理；原有测试和用户文件保留。API/RPC目录仅更新开发记录。


### 2026-11-22 / DEV-20261122-012 / 删除公开链路路由包装入口

- **用户确认与范围** 按用户确认删除公开的WithRequestTraceRouter，不保留兼容入口或别名；保留包内私有的withRequestTraceCORSRouter，由NewServerWithCORS统一调用，业务不直接装配链路路由。此前记录中保留该公开入口的约定由本条替代。
- **实现与说明** 删除默认CORS路由选项包装函数，将采样率、声明长度上限的单位及默认值注释移至私有装配函数；原有CORS、请求ID、单值头和正文预检查处理保持。同步Trace、HTTP头、正文限制与CORS使用说明及EnsureRequestTrace注释；当前源码没有需迁移的业务调用。已有外部调用应迁移到NewServerWithCORS(restConf,corsConfig)。
- **验证与限制** 用户明确本轮不用测试验证，未运行构建、编译、测试、范例、格式或差异检查；仅读取和定位修改所需源码及文档，不将前轮检查结果作为本轮验证结论。
- **测试与产物** 未新增测试用例、测试脚本、独立缓存或临时产物，无需清理；原有测试和用户文件保留。API/RPC目录仅更新开发记录。

### 2026-11-22 / DEV-20261122-013 / 第41项受控Server构造配置

- **用户确认与范围** 按确认方案删除NewServerWithCORS的任意rest.RunOption透传，统一请求ID、单值头校验和CORS路由仅由包内装配，不保留旧运行选项兼容分支。用户明确本轮不需要测试和验证。
- **接口** 构造函数改为NewServerWithCORS(restConf,corsConfig,handlerConfig ...ServerHandlerConfig)。可选配置最多一项，多项返回初始化错误；结构只开放NotFoundHandler和NotAllowedHandler两个http.Handler字段，不接收Router或CORS替换选项。两参数调用不变，省略、零值或字段nil使用默认404/405。处理器仅通过代码传入，不加入YAML配置。
- **实现** 包内先安装统一路由，再在同一路由上设置受控404/405处理器；默认404保留go-zero处理，默认405保留原生Allow头。自定义处理器自行输出状态、正文及需要的Allow等响应头，不应覆盖入口CORS和X-Request-ID。有效OPTIONS仍在入口结束，不执行回退处理器；回退路径仍不执行已注册业务路由的认证等应用中间件。
- **迁移与文档** 原rest.WithNotFoundHandler/rest.WithNotAllowedHandler参数迁移到明确字段，其余rest.RunOption不再接受。同步中文接口、字段和关键顺序注释及CORS、Trace接入说明，补充自定义回退处理器示例。定位到的现有构造调用为两参数使用方式，无需修改其他项目业务代码。此前文档中通过任意构造选项设置回退处理器的约定由本条替代。
- **验证与限制** 仅为实施修改读取和定位相关源码、文档及本地go-zero依赖实现；按用户要求未执行构建、编译、测试、范例运行、格式或差异检查，未进行运行时验证，不沿用前轮结果声称本轮通过。
- **测试与产物** 本轮未新增测试文件、测试脚本、临时产物或独立缓存，无需清理；原有测试和用户已有文件保留。API/RPC目录仅更新开发记录。

### 2026-11-22 / DEV-20261122-014 / 第42项先恢复请求ID并拒绝覆盖活动链路

- **用户确认与范围** 用户明确先读取请求值，合法入站ID直接复用，缺失时再生成，并在了解traceparent含义及来源后确认开发。保持HTTP先解析再初始化、RPC缺少traceparent直接拒绝及消息内部W3C传递规则；仅补齐已有活动上下文冲突时的统一保护。沿用此前不执行测试和验证的要求。
- **实现** ExtractTrace解析合法W3C上下文后，同Trace保留当前Span，不同Trace返回原上下文及ErrTraceConflict，不再覆盖。contextWithTraceID对只有HTTP请求ID和SetTrace调用采用相同保护：同链路不补随机Span，不同链路直接报错，仅无活动Span时随机补齐非零8字节Span ID，不重新生成Trace ID。EnsureTrace、SetTrace及RabbitMQ恢复移除已由公共入口覆盖的重复冲突判断，统一依赖核心规则；缺失、格式及错误返回路径保持。
- **入口与错误** HTTP继续先ExtractRequestTrace再EnsureTrace；入站合法值不变，活动上下文和出站链路均不存在时才生成新的Trace ID。HTTP冲突使用现有400，RPC使用现有InvalidArgument；RabbitMQ沿用原冲突拒绝及终态处理，无新增字段、别名或兼容入口。错误响应诊断ID约定保持，不将非法输入作为下游业务链路。
- **注释与说明** 更新公开函数、错误和关键绑定逻辑的中文注释，明确Trace ID与补齐Span ID的区别；同步统一请求ID规则，记录只有traceparent时的恢复、同链路Span保留、冲突拒绝及HTTP/RPC/MQ边界。对外函数签名和ID生成格式保持，已有依赖覆盖不同活动链路的调用需改为传入正确上下文。
- **验证与限制** 仅为实施修改读取和定位源码、规则及开发记录；未执行构建、编译、测试、范例运行、格式或差异检查，无运行时和链路联调结果，不声称已验证通过。
- **测试与产物** 本轮未新增测试用例、测试脚本、缓存或临时产物，无需清理；原有测试和用户已有文件保留。其他业务项目源码未修改，API/RPC目录仅补充开发记录。

### 2026-11-23 / DEV-20261123-001 / 第43项API中间件初始化错误返回

- **用户确认与范围** 按确认方案将API启动阶段的配置失败改为返回error：ApiUseMiddleware返回error，两个HTTP认证中间件构造函数返回组件指针与error，不保留旧接口包装。用户明确不验证，本轮不执行测试或验证。
- **接口与校验** ApiUseMiddleware拒绝nil接收者/Server、nil装配选项、缺少认证回调和非法并发/正文/AES配置；NewServerMiddleware继续记录选项，nil选项错误保存并在装配时返回，不执行后续选项。NewUserAgentMiddleware、NewUserAuthMiddlewareWithErr拒绝nil选项或缺少回调，返回nil组件和错误，不再因这些配置错误panic。回调签名及运行时认证规则保持。
- **装配顺序** 先创建认证、正文、请求头、CORS及可选并发组件，最后构造解密组件并应用校验后的进程级AES配置；全部成功后才逐个调用Server.Use。配置失败不追加本次任何中间件，AES配置无效时保留原配置。成功链路仍为CORS、可选并发、正文限制、解密、请求头和认证；此接口仅在启动阶段调用一次，不新增运行时重装配或异常恢复逻辑。
- **实现** authMiddleware返回(rest.Middleware,error)，请求头辅助方法改为仅创建组件，删除原先直接装配解密及panic的辅助路径。补充中文接口、选项错误和初始化顺序注释。调用方必须处理返回错误、记录原因、清理已有资源并停止启动。
- **说明与调用范围** 同步CORS、请求体、认证、HTTP头使用说明；已有可执行范例均增加ApiUseMiddleware错误处理，创建Server的范例在失败时Stop后返回错误。认证配置范例改为返回error，在API装配前明确检查独立服务状态组件的必需回调并创建该组件；独立服务状态构造接口保持原样。定位到user-api同名调用引用github.com/doubunv/common-pkg并替换到另一个旧模块，非当前base-pkg消费者，本轮未修改其源码或依赖。当前base-pkg源码中的认证构造调用统一处理两个返回值。
- **验证与限制** 仅为实施修改读取和定位相关源码、文档及本地go-zero的Use/Stop实现；未执行构建、编译、测试、范例运行、格式或差异检查，未启动服务或访问外部设施，无运行时验证结论。未新增自动资源清理或全局日志退出策略，资源生命周期由调用方管理。
- **测试与产物** 未新增测试文件、测试脚本、独立缓存或临时产物，无需清理；原有测试和用户已有文件保留。API/RPC目录仅更新记录，其他项目源码保持。

### 2026-11-23 / DEV-20261123-002 / 第44项独立服务状态组件初始化错误返回

- **用户确认与范围** 用户确认开发独立服务状态中间件的初始化错误处理，沿用此前不执行测试和验证的要求。此前AES进程级配置一致性提案按用户要求不实施，由运维保证，本轮不修改AES配置逻辑。
- **接口** NewCheckServiceStatusMiddleware改为返回(*CheckServiceStatusMiddleware,error)，成功返回组件与nil，缺少状态回调或包含nil选项时返回nil与配置错误；nil选项错误标明索引，不执行该项及后续选项。WithCheckServiceStatusCheckOption接受nil时由构造函数统一返回错误，不保留旧单返回值包装。
- **行为边界** 仅修改启动构造阶段，不自动追加中间件。CheckServiceStatusFunc签名、Handle、CheckServiceStatusRequest及运行时响应保持：含白名单和匿名在内的所有请求仍执行服务状态检查，false或运行时缺少回调继续返回原服务不可用错误。启动调用方先处理构造错误，再装配并启动服务。
- **注释与说明** 补充公开构造函数和回调选项的中文用途及条件注释。AUTH.md示例先接收组件与error、失败直接返回，再装配API链路，并将状态检查接在认证之后；删除示例中重复的回调预检查和不再使用的导入。更新三类HTTP认证/状态组件的统一构造约定，原第43项中独立状态构造接口保持不变的记录由本条补充替代。当前定位到的调用仅为公共包使用范例，无需修改其他项目业务代码。
- **验证与限制** 仅为实施修改读取和定位相关源码、调用位置及文档；未执行构建、编译、测试、范例运行、格式或差异检查，未启动服务或进行外部联调，不声称已验证通过。
- **测试与产物** 未新增测试文件、测试脚本、缓存或临时产物，无需清理；原有测试和用户已有文件保留。API/RPC目录仅补充记录，其他项目源码保持。

### 2026-11-23 / DEV-20261123-003 / 第46项搜索结果完整性校验

- **用户确认与范围** 用户确认分别修改Elasticsearch 7与OpenSearch，禁止部分搜索结果，并在解析时校验超时、失败分片及总数relation，检查通过后再写入目标。两套实现保持独立，无共享错误或类型别名。沿用此前不执行测试和验证的要求。此前TG通知重试优化和第45项Kafka回调隔离均按用户要求跳过，未修改对应源码。
- **请求** 两套EsModel.Search固定追加WithAllowPartialSearchResults(false)，继续固定track_total_hits=true，保持原参数、返回值、超时预算及SDK重试策略，不修改调用方query或引入允许部分结果的开关。
- **解析与错误** 两套core各自新增ErrIncompleteSearch，供errors.Is判断。timed_out=true、_shards.failed>0、terminated_early=true或对象total.relation=gte在写入目标前返回该错误，不将部分结果或总数下界当成成功。timed_out、_shards及failed为必需字段，缺失/null、负数failed/total、对象total缺失value或relation及未知relation返回ErrInvalidResponse；JSON类型错误仍保留原解码错误链。对象total只接受eq，原数字形式非负total继续支持。
- **结果与边界** 完整性检查失败返回总数0且不修改目标，OpenSearch同时返回nil聚合；正常空命中仍使用非nil空切片，无聚合的成功响应仍返回空map。检查位于文档解码前，不新增对自定义JSON解码器副作用的回滚保证。非2xx仍按原ResponseError等路径处理，不强制转换为ErrIncompleteSearch；响应体清理和错误脱敏保持。
- **说明与迁移** 补充公开错误、解析函数和搜索模型的中文注释，两套README分别更新完整性字段、filter_path必要字段、错误分类和errors.Is使用示例。已有依赖部分结果、下界总数或省略必要字段的调用现在收到错误，应先处理错误再使用结果；函数签名保持，不修改其他业务项目源码。
- **验证与限制** 仅为实施修改读取和定位源码、文档及本地SDK参数接口；方案阶段参考官方Search API说明。未执行构建、编译、测试、范例运行、格式或差异检查，未连接Elasticsearch/OpenSearch集群或进行代理/分片故障联调，不声称已验证通过。
- **测试与产物** 本轮未新增测试文件、测试脚本、缓存或临时产物，无需清理；原有测试和用户已有文件保留。API/RPC目录仅补充开发记录。

### 2026-11-23 / DEV-20261123-004 / 第47项搜索客户端成功响应读取上限

- **用户确认与范围** 按确认给Elasticsearch 7与OpenSearch分别增加成功响应大小限制，默认16MiB，可配置，超限返回错误并关闭响应体，不解析截断数据。两套配置、错误和实现保持独立，模型和公开响应解析函数签名不变。沿用前面的不测试、不验证要求。
- **配置** 两套core.Config增加MaxResponseBytes，单位字节，0默认16777216；负数及MaxInt64初始化失败，不提供无限读取模式。客户端构造时复制已校验值并固定为本实例配置，供其成功响应使用；不通过业务上下文或全局变量传递上限。
- **读取实现** 传输层为每个2xx响应绑定独立的limitedResponseBody，覆盖SDK启动产品检查/Info探活、模型查询/写入及通过本客户端SDK读取的成功正文。按实际读取字节计数，不信任Content-Length；自动解压时在解压后计数。最多读取上限字节及一个额外判定字节，额外字节不交给调用方，正好上限且EOF仍成功。超限立即关闭底层正文并返回各自ErrResponseTooLarge，后续非空读取仍返回该错误。Close通过sync.Once保证可重复或并发调用，超限和解析器清理不会重复关闭底层正文。
- **解析与边界** ReadResponse沿用客户端响应上的实例上限；独立SDK响应默认绑定16MiB，不叠加默认值缩小已配置的较大上限。超限在读取阶段返回错误及nil正文，Get/Search不会进入目标解码，Search返回总数0，OpenSearch同时返回nil聚合。原4KiB错误响应摘要和取消优先规则保持，错误文本不含原始正文或查询。此限制不等同于整个进程内存上限，解码和业务对象仍需额外内存；直接SDK调用方仍须处理Read/解码错误并关闭正文。
- **说明与行为影响** 更新公开配置、错误、客户端及读取逻辑中文注释；两套README分别增加YAML字段、单位与默认值、errors.Is用法、精确边界和大响应调优说明。合法大响应超限需显式调大配置或减少单次查询量；写操作响应超限不代表服务端写入未发生，不新增自动重试。Info超限按现有启动失败流程释放连接资源。其他项目业务源码未修改。
- **验证与限制** 仅为实施修改读取和定位当前源码、调用位置及文档；未执行构建、编译、测试、范例运行、格式或差异检查，未连接搜索集群或模拟大响应，不声称已验证通过。
- **测试与产物** 本轮未新增测试文件、测试脚本、缓存或临时产物，无需清理；原有测试和用户已有文件保留。API/RPC目录仅补充开发记录。

### 2026-11-23 / DEV-20261123-005 / 第48项JWT依赖迁移至v5

- **用户确认与范围** 用户确认将公共包JWT依赖从已归档的dgrijalva/jwt-go迁移至golang-jwt/jwt/v5，保持GenerateToken结构体参数、ParseToken调用方式、HS256、密钥要求、时间校验及业务JSON字段。沿用前面的不测试、不验证要求。
- **依赖维护** 固定github.com/golang-jwt/jwt/v5 v5.3.1，其go.mod要求Go1.21，项目go1.22.0/toolchain声明保持。复用会话前已有的模块缓存，执行GOWORK=off GOPROXY=off GOSUMDB=off go mod download -json github.com/golang-jwt/jwt/v5@v5.3.1取得依赖元数据和两项校验和，不发起网络下载，不编译代码。更新go.mod/go.sum，删除本模块旧库直接依赖及对应校验和；现有框架所用jwt/v4间接依赖保留，不作为本包兼容入口。未运行go mod tidy或全量依赖升级。
- **Claims** 本包直接声明业务和标准字段，实现v5的六个Claims Getter，不再嵌入旧StandardClaims或保留旧库方法包装。exp/iat/nbf保持int64 Unix秒，aud保持单字符串，Id及其他Go字段和JSON标签保持；JSON解析不经float64转换，不扩展为小数秒、字符串时间或数组受众。仅向校验器提供NumericDate转换，不修改库全局TimePrecision或序列化设置。直接构造Claims的旧StandardClaims字段需改为填写本包标准字段。
- **签发与解析** 签发保持原业务字段、TokenType默认0、exp秒精度和issuer=test-user，零/负时长仍按原签发行为处理。解析使用WithValidMethods限定HS256、WithExpirationRequired要求exp、WithIssuedAt保持未来iat拒绝，nbf保持默认校验；显式配置当前解析器的time.Now，不再依赖旧库全局jwt.TimeFunc。签名和声明校验后再检查exp为正数且严格晚于当前Unix秒，到期错误保留v5.ErrTokenExpired错误链；所有失败仍返回nil Claims，不新增issuer/audience/subject/UID规则。
- **调用与说明** 定位到其他子项目引用旧模块或自有JWT实现，不是当前base-pkg消费者，本轮不迁移其源码或依赖。当前公共包未发现需调整的业务Claims调用。补充公开Claims字段、Getter及关键校验的中文注释，README增加版本、Claims迁移、时间源及底层错误链说明，现有签发/解析范例调用方式保持。
- **验证与限制** 仅读取和定位源码、调用、现有缓存及官方发布/迁移资料，并执行上述离线依赖维护命令；未执行构建、编译、测试、范例运行、格式或差异检查，也未验证Token跨版本互验或API/RPC联调，不声称已验证通过。底层错误类型及已移除StandardClaims的类型断言需由使用方迁移。
- **测试与产物** 未新增测试文件、测试脚本、临时产物或本轮独立缓存，无需清理；复用原有模块缓存，不删除原有依赖或测试。API/RPC目录仅更新开发记录，其他业务项目源码保持。

### 2026-11-23 / DEV-20261123-006 / 第49项RabbitMQ驱动维护迁移

- **用户确认与范围** 按确认从streadway/amqp迁移至RabbitMQ团队维护的amqp091-go，保持本包发布、消费、重连、重试、确认和单条ACK规则。沿用不执行测试或验证的要求。用户已跳过的TG通知重试队列优化、Kafka及AES运维一致性检查均未实施。
- **依赖维护** 固定github.com/rabbitmq/amqp091-go v1.17.0，选定版本go.mod要求Go1.20，项目Go1.22和toolchain声明保持。更新go.mod/go.sum，删除旧驱动直接依赖和对应校验和。当前其他业务项目未定位到旧驱动直接引用，本轮不迁移其模块；未运行go mod tidy或升级其他依赖。
- **内部适配** 八份RabbitMQ源文件改为新驱动原生类型，统一覆盖配置URI解析、连接/Channel、Delivery、Publishing、Confirmation、Return、Error、Table及Decimal。内部brokerChannel接口继续使用现有Publish/Consume/通知方法，业务Message和公开发布/消费签名保持。不保留旧包类型兼容转换；直接使用旧Table/Decimal或原生Error断言的调用方须迁移。
- **恢复与超时** amqp.Config.Recovery明确设为nil，驱动自动恢复关闭，本包继续重建完整会话、声明拓扑、恢复绑定/消费者并注册发布与关闭通知。发布串行锁、确认序号、mandatory、归档双份确认、ACK/NACK的multiple=false及原有上下文关闭底层连接策略保持，没有启用另一套恢复或重发机制。
- **心跳适配** 新驱动将Config.Heartbeat=0当作缺省10s，适配器将归一化的Heartbeat按整秒写入独立连接URI，显式0得以按原配置约定参与broker协商，避免静默恢复默认值。Config.Heartbeat优先于原URL的heartbeat参数，其他URI参数保留，调用方原字符串不修改或回显。补充配置字段和适配逻辑中文注释。
- **依赖获取命令** 在common-pkg目录执行 `GOWORK=off GOMODCACHE=/private/tmp/new-platform-amqp-driver.iAnHaWxQ/modules GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org go mod download -json github.com/rabbitmq/amqp091-go@v1.17.0`。首次默认沙箱因DNS解析失败未取得模块；同命令经获准的网络访问下载成功，取得模块和校验和。这是依赖获取，不是项目构建或测试，下载源码仅存放本轮临时目录。
- **注释与说明** 更新连接/驱动边界、消息头原生类型及Heartbeat的中文注释，README和USAGE增加版本、恢复所有权、原生类型迁移及心跳优先级说明，既有普通业务范例调用方式保持。
- **验证与限制** 仅为适配读取和定位项目源码、本轮下载的驱动源码及官方资料，并执行上述依赖获取命令；未执行构建、编译、测试、范例运行、格式、差异检查或broker联调，不声称已验证驱动迁移、确认事件及重连行为通过。
- **清理完成** 已在清理前记录实际依赖命令、结果和限制，随后删除下载到任务目录的依赖源码、zip、模块元数据和sumdb缓存，任务目录已不存在。本轮无新增项目测试、测试脚本或遗留临时产物，原有模块缓存及项目测试未清理；API/RPC目录仅更新记录。

### 2026-11-23 / DEV-20261123-007 / 第50—53项公共组件集中修改

- **用户确认与范围** 用户确认一起开发第50—53项：搜索回填文档ID、同一实例防止重复装配、减少metadata重复复制、删除无效测试开关。仅修改公共包实现、中文注释与使用说明，并同步平台/API/RPC开发记录；不迁移其他业务项目。沿用此前不测试、不验证的要求。
- **第50项搜索ID** 两套独立解析器读取每条命中的原生_id，支持非nil切片指针接收[]T和[]*T。元素或其地址实现本套Schema时调用SetId，嵌入EsSchema的模型同样适用；原生_id优先于正文映射到同一字段的值，普通非Schema元素继续只解码_source。需要回填ID但_id为空/缺失时返回本套ErrInvalidResponse，不回退到正文ID。两套模型、Schema及解析接口的中文注释和README范例同步。
- **第50项提交与边界** JSON先解码到新的同类型切片，数量与命中一致后按位置回填ID，提交前检查上下文，再替换目标；普通解码/ID错误不替换旧目标，失败总数0，OpenSearch同时返回nil聚合。成功结果不复用旧元素，不保留响应未提供的旧字段；空结果保持非nil空切片。自定义切片解码器须保持命中数量和顺序，数量不符返回ErrInvalidResponse；自定义解码器或SetId的外部副作用不回滚。目标不为非nil切片指针时返回解码错误。filter_path需为Schema保留hits.hits._id及_source；公开Search签名、完整性和正文上限规则保持，两套实现不共用类型或错误。
- **第51项重复装配** ServerMiddleware增加实例锁和成功标记，同一对象的装配调用串行化；成功追加全部中间件后标记，后续调用返回错误，不再次追加或应用AES配置。配置初始化失败不置成功标记，先校验再统一追加的规则保持。保护仅针对同一对象，不跨实例或按rest.Server建立全局登记；同一Server应保留一个装配对象，开始使用后不得复制或并发修改配置。仍只在启动阶段装配，不支持运行时重装配。
- **第52项metadata** 提取私有injectTraceMetadata，在调用方独占的可写MD副本上清理旧tracestate并注入当前W3C链路。HeadInMetadata在同一副本完成15个平台字段覆盖、旧traceparent/tracestate清理及当前链路注入，不再建立中间出站上下文后重新提取/复制MD。无活动Trace时仅清理旧链路，不生成ID；SetTrace/EnsureTrace等独立入口仍先复制再调用私有注入函数。保留出站未知字段及多值，不合并入站数据，不修改父上下文或并行分支，不改变字段命名、Trace冲突和生成规则。
- **第53项无效开关** 删除WithTestOption及仅被它写入的私有isTest字段，不提供兼容入口。原开关没有实际执行分支，调用方直接移除选项；WithDebugOption及调试行为保持。CORS接入说明记录迁移方式。
- **文档** 更新两套搜索README和ID读取示例，CORS/正文限制/认证说明补充同一实例防重及使用后不得复制的边界；METADATA说明同步一次副本装配流程。上述规则补充此前开发记录，以本条当前行为为准。
- **验证与限制** 本轮仅为实施修改读取和定位源码、使用说明及开发记录；按用户要求未执行构建、编译、测试、范例运行、格式、差异或其他验证命令，未进行搜索集群/API/RPC联调，无运行时或并发验证结论，不声称已验证通过。
- **测试与产物** 本轮未新增测试用例、测试脚本、独立缓存或临时产物，无需清理；原有测试及用户已有文件保留。API/RPC目录仅补充接入边界记录，其他项目业务源码和依赖未修改。

### 2026-11-23 / DEV-20261123-008 / 第54项按用户要求跳过

- **决策** 用户明确跳过XXL-JOB任务结果回调重试，不实施该方案。现有回调失败记录日志的处理保持，不新增重试配置或队列。
- **范围** 仅记录决策并更新进度，未修改业务代码、接口或依赖；未执行构建、测试或验证，无新增测试及临时产物。

### 2026-11-23 / DEV-20261123-009 / 第55项统一带请求上下文的HTTP错误入口

- **用户确认与范围** 用户确认删除logs/xcode.ErrHandler，仅保留ErrHandlerCtx作为HTTP框架错误处理入口，不保留旧签名包装。同步中文注释、接入说明和开发记录；第54项XXL-JOB回调重试继续跳过。沿用此前不执行测试和验证的要求。
- **实现与协议** 删除使用context.Background调用ErrHandlerCtx的旧函数，避免通过该公开入口丢失请求Trace ID和语言。ErrHandlerCtx继续委托result.HttpErrorResponse构建状态及响应对象，现有签名、HTTP200业务错误协议、错误分类、公开文案及翻译规则保持；不在响应阶段生成ID，不新增上下文缺失回退或请求头读取。InitAppServer原本已使用httpx.SetErrorHandlerCtx注册，执行逻辑保持，仅补充启动和上下文要求注释。
- **接入与兼容性** 原httpx.SetErrorHandler(errorhandler.ErrHandler)需改为httpx.SetErrorHandlerCtx(errorhandler.ErrHandlerCtx)；业务使用httpx.ErrorCtx(r.Context(),w,err)或result.HttpErrorResult(r.Context(),w,err)，不能传Background/TODO替代请求ctx。标准请求入口继续负责Trace及语言注入，词典继续显式加载。本轮定位到当前base-pkg源码只有旧函数定义，标准启动使用正确入口；未修改其他旧模块业务项目或依赖。此前记录中保留旧ErrHandler签名的约定由本条替代，历史记录原文保留。
- **说明** result/README新增独立注册范例和旧入口迁移说明；TRACE_ID规则补充HTTP错误出口须沿用请求ctx并链接接入说明。ErrHandlerCtx和InitAppServer中文注释同步。
- **验证与限制** 仅读取和定位修改所需源码、引用、文档及开发记录；按用户要求未运行构建、编译、测试、范例、格式、差异或其他验证命令，不声称已验证错误响应或框架联调通过。
- **测试与产物** 本轮未新增测试用例、测试脚本、独立缓存或临时产物，无需清理；原有测试和用户已有文件保留。API/RPC目录仅更新接入边界记录，其他业务项目源码和依赖未修改。
