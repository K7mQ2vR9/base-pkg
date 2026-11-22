# JWT 签发与有效期校验

`GenerateToken` 改为接收一个 `GenerateTokenParams` 值，返回类型仍为 `(string, error)`，直接替换旧的多个位置参数。`ParseToken` 调用方式不变。签发使用 HS256，密钥至少32字节且不能全为空白，有效时长由调用方显式提供，Claims 的 JSON 字段及当前签发内容保持不变。

## 签发参数

```text
GenerateToken(params GenerateTokenParams) (string, error)
```

| 字段 | 类型 | 用途和零值行为 |
|---|---|---|
| JWTSecret | string | 原始签名密钥；空、少于32字节或全空白时拒绝 |
| Uid | int64 | 用户ID，按原值写入，不新增身份/范围校验 |
| Source | string | 来源渠道，允许为空，按原值写入 |
| Sign | string | 业务会话标识，允许为空，按原值写入 |
| BusinessCode | string | 业务标识，允许为空，保留现有 JSON 标签 Business_code |
| ExpireDuration | time.Duration | 自签发时刻起的有效时长，例如24*time.Hour，无默认值 |

ExpireDuration 使用 Go 时长值，代替原先裸 int64 的纳秒参数。如果已有 int64 配置，须确认单位后显式转换；秒数应写为 `time.Duration(seconds) * time.Second`。零或负时长的现有处理保持：签发时不新增拒绝规则，但生成的 Token 会被有效期校验拒绝；小于一秒的时长受 exp 的 Unix 秒精度影响，可能签发即到期。签发时长由每次调用指定，旧多参数调用需要迁移为结构体。

## 有效期规则

`ParseToken` 要求 `exp` 是正数 Unix 秒，并且 **`exp > 当前校验时间的 Unix 秒`**。

| exp 情况 | ParseToken |
|---|---|
| 缺失、null、0、负数 | 拒绝 |
| 早于当前秒 | 拒绝 |
| 等于当前秒 | 拒绝 |
| 晚于当前秒，且其他校验通过 | 接受 |

时间来源与底层解析器一致，沿用 `jwt.TimeFunc`，默认为服务端当前时间，不增加时钟宽限。签名验证完成后再次确认未到期，避免验证过程中越过到期边界仍返回 Claims。失败时返回 nil Claims 和错误，不把已解码但未通过验证的身份交给业务。

刷新场景通过 `ParseToken` 校验现有令牌，业务身份及续签许可由调用方判断，不自动签发新令牌。

HS256 限制、密钥要求、已有 nbf/iat 校验及当前 TokenType 策略保持。签发写入现有 Uid/Source/Sign/BusinessCode、exp 和 issuer，保留 Claims 格式及当前 TokenType 默认值，不增加 issuer/audience/UID 校验规则。

## 使用示例

从项目配置或环境注入密钥，示例不输出密钥或完整 Token：

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/K7mQ2vR9/base-pkg/jwtToken"
)

func main() {
    secret := os.Getenv("JWT_SECRET")
    token, err := jwtToken.GenerateToken(jwtToken.GenerateTokenParams{
        JWTSecret: secret,
        Uid: 12,
        Source: "APP",
        Sign: "session-sign",
        BusinessCode: "catalog",
        ExpireDuration: 24 * time.Hour,
    })
    if err != nil {
        panic(err)
    }
    claims, err := jwtToken.ParseToken(secret, token)
    if err != nil {
        panic(err)
    }
    fmt.Println(claims.Uid, claims.BusinessCode)
}
```

配置至少32字节且非全空白的 `JWT_SECRET` 后运行，输出 `12 catalog`。所有 GenerateToken 旧调用需迁移为结构体；过期规则仍要求 exp 存在且严格晚于当前秒。

## 与 API 认证对接

认证回调中调用 `ParseToken`，失败时把 error 返回给认证中间件，不能将过期或缺少有效期伪装成 `0,nil` 匿名结果。白名单规则保持原样：没有 Token 可以匿名访问白名单，有 Token 且校验失败仍返回认证错误；详见 [认证与匿名白名单](../appMiddleware/AUTH.md)。

刷新流程只以解析时的有效期判断为准，解析成功不会延长原 Token 的 exp；业务需要继续按既有身份、会话和权限规则决定是否签发新 Token。
