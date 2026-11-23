// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package jwtToken

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const (
	hs256MinKeyBytes = 32
)

// Claims 保存平台业务声明及标准声明，直接实现 jwt/v5.Claims，不嵌入已移除的 StandardClaims。
// 时间字段保持 int64 Unix 秒，JSON 解码不接受小数、字符串或溢出值；不修改底层库的全局配置。
type Claims struct {
	// Uid 为业务用户 ID，签发时不增加身份或范围校验。
	Uid int64 `json:"uid"`
	// Source 为来源渠道，默认空字符串。
	Source string `json:"source"`
	// Sign 为业务会话标识，默认空字符串。
	Sign string `json:"sign"`
	// BusinessCode 为业务标识，保留原 JSON 字段名 Business_code。
	BusinessCode string `json:"Business_code"`
	// TokenType 保留业务字段，签发默认 0，不新增类型校验规则。
	TokenType int `json:"token_type"`
	// Audience 为可选的单个受众文本，默认空且不输出，本包不增加受众匹配规则。
	Audience string `json:"aud,omitempty"`
	// ExpiresAt 为到期 Unix 秒，0 不输出；ParseToken 要求正数且严格晚于当前秒。
	ExpiresAt int64 `json:"exp,omitempty"`
	// Id 为可选令牌标识，保留既有 Go 字段名及 jti 标签，默认空且不输出。
	Id string `json:"jti,omitempty"`
	// IssuedAt 为可选签发 Unix 秒，0 表示未设置；非零且晚于当前秒时拒绝。
	IssuedAt int64 `json:"iat,omitempty"`
	// Issuer 为可选签发者，GenerateToken 沿用 test-user，不增加签发者匹配规则。
	Issuer string `json:"iss,omitempty"`
	// NotBefore 为可选生效 Unix 秒，0 表示未设置；非零且晚于当前秒时拒绝。
	NotBefore int64 `json:"nbf,omitempty"`
	// Subject 为可选主题文本，默认空且不输出，不增加主题匹配规则。
	Subject string `json:"sub,omitempty"`
}

// GetExpirationTime 为 jwt/v5 提供到期时间，ExpiresAt=0 返回 nil，由解析器要求 exp 必填。
func (c Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return claimNumericDate(c.ExpiresAt), nil
}

// GetIssuedAt 为 jwt/v5 提供签发时间，IssuedAt=0 返回 nil，其他值按整数 Unix 秒解释。
func (c Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return claimNumericDate(c.IssuedAt), nil
}

// GetNotBefore 为 jwt/v5 提供生效时间，NotBefore=0 返回 nil，不新增必填要求。
func (c Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return claimNumericDate(c.NotBefore), nil
}

// GetIssuer 为 jwt/v5 提供签发者原文；返回值本身不代表已完成签发者校验。
func (c Claims) GetIssuer() (string, error) { return c.Issuer, nil }

// GetSubject 为 jwt/v5 提供主题原文，不自动增加业务权限或主题校验。
func (c Claims) GetSubject() (string, error) { return c.Subject, nil }

// GetAudience 为 jwt/v5 提供受众列表，空文本返回 nil；Claims 的 JSON 仍使用原单字符串形式。
func (c Claims) GetAudience() (jwt.ClaimStrings, error) {
	if c.Audience == "" {
		return nil, nil
	}
	return jwt.ClaimStrings{c.Audience}, nil
}

// claimNumericDate 只为校验接口转换整数 Unix 秒；0 表示可选字段未设置，不使用全局 TimePrecision。
func claimNumericDate(seconds int64) *jwt.NumericDate {
	if seconds == 0 {
		return nil
	}
	return &jwt.NumericDate{Time: time.Unix(seconds, 0)}
}

// GenerateTokenParams 定义签发 JWT 的参数，按值传入，不修改调用方配置。
type GenerateTokenParams struct {
	// JWTSecret 为 HS256 签名密钥，至少32字节且不能全为空白，使用原始字节签名。
	JWTSecret string
	// Uid 为业务提供的用户 ID，本函数不增加身份或范围校验。
	Uid int64
	// Source 为来源渠道，按原值写入 Claims。
	Source string
	// Sign 为业务会话标识，按原值写入 Claims。
	Sign string
	// BusinessCode 为业务标识，按原值写入 Claims。
	BusinessCode string
	// ExpireDuration 为自当前时刻起的有效时长，使用 time.Duration，例如24*time.Hour。
	// 无默认值，零值按当前时刻到期；沿用现有时长处理，签发时不新增正数校验。
	ExpireDuration time.Duration
}

// GenerateToken 使用结构体参数签发 HS256 Token，返回值及 Claims 格式保持不变。
// 有效期使用调用方提供的时长，exp 精度为 Unix 秒；密钥无效时返回空 Token 和错误。
func GenerateToken(params GenerateTokenParams) (accessToken string, err error) {
	key, err := jwtSigningKey(params.JWTSecret)
	if err != nil {
		return "", err
	}

	nowTime := time.Now()
	expireTime := nowTime.Add(params.ExpireDuration)
	claims := Claims{
		Uid:          params.Uid,
		Source:       params.Source,
		Sign:         params.Sign,
		BusinessCode: params.BusinessCode,
		ExpiresAt:    expireTime.Unix(),
		Issuer:       "test-user",
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to get accessToken")
	}

	return accessToken, err
}

// ParseToken 验证 HS256 签名及时间声明，exp 必须是正数 Unix 秒且严格晚于校验时刻。
// 使用服务端 time.Now，并为当前解析器显式配置时间源；到期当秒即拒绝，不增加宽限时间或全局时钟变量。
// 失败不返回 Claims；密钥要求、Claims 字段和其他时间声明校验保持不变。
func ParseToken(jwtSecret string, token string) (*Claims, error) {
	key, err := jwtSigningKey(jwtSecret)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	// v5 默认不校验 iat，显式开启以保持原有未来签发时间拒绝规则；exp 必填，不配置时钟宽限。
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(time.Now),
	)
	parsedToken, err := parser.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}
	if parsedToken == nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	// 签名和声明校验成功后再检查正数及秒级到期边界，避免校验过程中跨过 exp 仍返回身份。
	if claims.ExpiresAt <= 0 || claims.ExpiresAt <= time.Now().Unix() {
		return nil, errors.Wrap(jwt.ErrTokenExpired, "token has expired or has no expiration")
	}
	return claims, nil
}

// jwtSigningKey 校验 HS256 密钥配置，保留调用方配置的原始字节。
func jwtSigningKey(secret string) ([]byte, error) {
	if len(secret) < hs256MinKeyBytes || strings.TrimSpace(secret) == "" {
		return nil, errors.New("JWT signing key must be at least 32 bytes and not blank")
	}
	return []byte(secret), nil
}
