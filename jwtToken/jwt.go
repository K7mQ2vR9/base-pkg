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

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	hs256MinKeyBytes = 32
)

type Claims struct {
	Uid          int64  `json:"uid"`
	Source       string `json:"source"`
	Sign         string `json:"sign"`
	BusinessCode string `json:"Business_code"`
	TokenType    int    `json:"token_type"`
	jwt.StandardClaims
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
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "test-user",
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to get accessToken")
	}

	return accessToken, err
}

// ParseToken 验证 HS256 签名及时间声明，exp 必须是正数 Unix 秒且严格晚于校验时刻。
// 与底层解析器使用同一 jwt.TimeFunc（默认服务端当前时间），到期当秒即拒绝，不增加宽限时间。
// 失败不返回 Claims；密钥要求、Claims 字段和其他时间声明校验保持不变。
func ParseToken(jwtSecret string, token string) (*Claims, error) {
	key, err := jwtSigningKey(jwtSecret)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}
	parser := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Alg()}}
	parsedToken, err := parser.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}
	if parsedToken == nil || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	// 底层 StandardClaims 将 exp 视为可选且允许相等；平台访问令牌统一要求必填和严格未到期。
	if claims.ExpiresAt <= 0 || claims.ExpiresAt <= jwt.TimeFunc().Unix() {
		return nil, errors.New("token has expired or has no expiration")
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
