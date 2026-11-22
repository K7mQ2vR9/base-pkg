# 请求上下文与日志字段同步

API 认证完成后，业务日志中的 `token-uid` 与已验证的用户 ID 一致；白名单正常匿名请求记录 `"0"`。业务调用 `SetBusinessCode` 后，当前分支日志使用新的 `business-code`。

## 使用规则

| 调用 | metadata 行为 | 日志行为 |
|---|---|---|
| `SetTokenUid(ctx, value)` | 复制并覆盖出站用户 ID | 覆盖当前分支的用户 ID 字段 |
| `SetBusinessCode(ctx, value)` | 分别复制并覆盖入站、出站业务标识 | 覆盖当前分支的业务标识字段 |
| `ContextHeadInLog(ctx, head)` | 不修改 metadata | 按 LogKeys 写入 HTTP 字符串字段 |
| `ContextMetadataInLog(ctx, md)` | 不修改 metadata | 按 LogKeys 写入 RPC 字符串列表字段，列表使用副本 |

所有调用都返回新 ctx，业务必须接住返回值，并用该 ctx 记录后续日志、调用 RPC。父上下文及其他分支不改变，原有业务日志字段、Trace、取消和截止时间保留。同名上下文字段以后写入的值为准，重复调用不会持续追加同名字段。

HTTP 字段格式如 `"token-uid":"12"`；RPC 保持原格式 `"token-uid":["12"]`。Setter 在没有 HTTP/RPC 日志初始化时默认使用字符串。传入空字符串会覆盖旧日志值为空，不保留旧值；匿名身份由认证流程显式写入 `"0"`。Setter 仍不执行身份验证或格式校验，日志同步不等于认证成功。

字段继续受 `headInfo.LogKeys` 控制，请求期间只读。Setter 不会把名单之外的字段加入日志；Token 不由 metadata 日志入口自动输出，用户扩展 JSON 也不会随同步写入日志。

标准 API 的认证流程已调用 `SetTokenUid`，正常业务无需另外补写用户 ID 日志。RPC 入口已调用 `ContextMetadataInLog`，业务仍使用相同 Setter 修改上下文。

## 可运行示例

```go
package main

import (
    "context"
    "os"

    "github.com/K7mQ2vR9/base-pkg/headInfo"
    "github.com/K7mQ2vR9/base-pkg/logs"
    "github.com/zeromicro/go-zero/core/logx"
)

func main() {
    // 独立示例直接输出 JSON 日志；实际服务沿用现有日志初始化。
    logx.SetWriter(logx.NewWriter(os.Stdout))
    ctx, err := headInfo.SetTrace(context.Background(), "1234567890abcdef1234567890abcdef")
    if err != nil {
        panic(err)
    }
    ctx = headInfo.ContextHeadInLog(ctx, &headInfo.Head{
        ReqPath: "/items",
        BusinessCode: "original",
    })
    // 这里演示同步效果；实际用户 ID 必须来自认证后的身份。
    ctx = headInfo.SetTokenUid(ctx, "12")
    ctx = headInfo.SetBusinessCode(ctx, "catalog")
    logs.Info(ctx, "current branch")

    child := headInfo.SetBusinessCode(ctx, "child")
    logx.WithContext(child).Info("child branch")
    logx.WithContext(ctx).Info("parent branch")
}
```

三条日志的 `token-uid` 均为 `"12"`，`business-code` 依次为 `catalog`、`child`、`catalog`；Trace 相同，其他字段保留。`logs.Info/Error`、`logc` 和 `logx.WithContext` 都读取相同的上下文字段，原生 go-zero JSON 与平台 Logrus 输出均可使用。

## 上下文边界

已经绑定旧 ctx 的 logger 不会自动切换到新 ctx；修改后应使用返回的 ctx 创建或获取 logger。外层中间件持有的旧请求上下文同样不自动回写，本功能同步的是使用新 ctx 的后续日志。

直接调用 `metadata.NewOutgoingContext`、手工修改 metadata 或单独构造 `Head` 不会自动同步日志。优先使用上述 Setter；需要批量同步时显式调用 `ContextMetadataInLog`，字段格式随该入口切换为 RPC 列表。不会解析日志正文或用户扩展 JSON 来推断身份。

去重范围是平台合并的上下文字段。业务在单条日志、全局日志或第三方 writer 中另行追加同名字段时，应按其接口规则处理。
