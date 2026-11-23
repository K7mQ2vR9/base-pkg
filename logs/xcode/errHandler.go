package xcode

import (
	"context"

	"github.com/K7mQ2vR9/base-pkg/result"
	"github.com/K7mQ2vR9/base-pkg/result/xcode"
	"google.golang.org/grpc/status"
)

// ErrHandlerCtx 为 httpx.SetErrorHandlerCtx 构建统一错误响应，不直接写入正文。
// ctx 须为当前请求上下文，沿用其中的 Trace ID 和语言；不在响应阶段生成 ID 或补建后台上下文。
func ErrHandlerCtx(ctx context.Context, err error) (int, any) {
	return result.HttpErrorResponse(ctx, err)
}

// FromError 使用统一规则构建 gRPC 状态，支持 nil 和包装后的错误。
func FromError(err error) *status.Status {
	return xcode.ToGRPCStatus(err)
}

// Deprecated : 使用 pkg/xcode/New 替换
func New(code int, msg string) xcode.Code {
	return xcode.New(code, msg)
}
