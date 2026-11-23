# JWT 签发与有效期校验

`GenerateToken` 改为接收一个 `GenerateTokenParams` 值，返回类型仍为 `(string, error)`，直接替换旧的多个位置参数。`ParseToken` 调用方式不变。签发使用 HS256，密钥至少32字节且不能全为空白，有效时长由调用方显式提供，Claims 的 JSON 字段及当前签发内容保持不变。

底层依赖为 `github.com/golang-jwt/jwt/v5 v5.3.1`，已移除本模块对 `github.com/dgrijalva/jwt-go` 的直接依赖。选定版本的go.mod要求Go1.21，项目Go1.22声明保持。版本与迁移背景见 [v5.3.1发布记录](https://github.com/golang-jwt/jwt/releases/tag/v5.3.1) 和 [官方迁移说明](https://github.com/golang-jwt/jwt/blob/v5.3.1/MIGRATION_GUIDE.md)。

## Claims 与底层接口

本包的 `Claims` 直接实现v5的 `jwt.Claims` 接口，不再嵌入底层 `StandardClaims`，也不提供旧类型或方法包装。业务字段 `uid/source/sign/Business_code/token_type` 及标准声明JSON标签保持。

| 本包字段 | 类型及JSON形式 |
|---|---|
| ExpiresAt / IssuedAt / NotBefore | int64 Unix秒，对应exp/iat/nbf，零值不输出；不接受小数、字符串或溢出数字 |
| Audience | 单个string，对应aud，保持字符串形式，不改为数组 |
| Id | string，对应jti，保留既有Go字段名 |
| Issuer / Subject | string，对应iss/sub，空值不输出 |

原先构造 `Claims{StandardClaims: ...}` 的代码应改为直接填写本包字段，例如 `Claims{ExpiresAt: expires, Issuer: issuer}`；直接读取 `claims.ExpiresAt` 等字段仍使用原int64类型。`GetExpirationTime/GetIssuedAt/GetNotBefore/GetAudience/GetIssuer/GetSubject` 仅提供v5校验接口，不执行签名验证，业务应通过 `ParseToken` 获取已验证身份。

未增加issuer/audience/subject匹配要求，也未改变业务UID或TokenType策略。旧 `Valid` / `VerifyExpiresAt` 等由StandardClaims提供的方法不再存在，不能沿用对旧库的类型断言。底层解析错误现使用v5错误链；需要分类时可用 `errors.Is(err, jwt.ErrTokenExpired)` 或 `errors.Is(err, jwt.ErrTokenRequiredClaimMissing)`，但任何解析错误都必须拒绝身份使用。

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

时间来源为服务端 `time.Now`，通过 `jwt.WithTimeFunc(time.Now)` 配置在每次创建的解析器上，不再使用旧库的全局 `jwt.TimeFunc`，也不新增全局可变时钟。显式启用 `WithExpirationRequired` 与 `WithIssuedAt`：exp必填，未来iat拒绝，nbf沿用v5默认生效时间校验，iat/nbf为0或缺失时不新增必填要求。不增加时钟宽限。签名验证和声明校验完成后，再次检查exp为正数且严格晚于当前Unix秒，避免验证过程中越过到期边界仍返回Claims；此处到期错误保留v5的ErrTokenExpired错误链。失败时返回nil Claims，不把已解码但未通过验证的身份交给业务。

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
