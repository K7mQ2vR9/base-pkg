package interceptors

import (
	"context"

	"github.com/K7mQ2vR9/base-pkg/consts"
	"github.com/K7mQ2vR9/base-pkg/headInfo"
	"github.com/K7mQ2vR9/base-pkg/result/xcode"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ClientInterceptor 复制调用 metadata 并传递链路，按统一规则恢复 RPC 业务错误。
func ClientInterceptor(rpcName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, err := headInfo.EnsureTrace(ctx)
		if err != nil {
			return err
		}
		md, _ := metadata.FromOutgoingContext(ctx)
		md = md.Copy()
		// 自定义传输头统一使用中划线，rpc-name 在当前调用副本内覆盖，不影响上游。
		md.Set(consts.RpcName, rpcName)
		// 链路仅通过标准 W3C 上下文传递，不附加重复 ID 字段。
		trace.Inject(ctx, otel.GetTextMapPropagator(), &md)
		trace.Inject(ctx, propagation.TraceContext{}, &md)
		ctx = metadata.NewOutgoingContext(ctx, md)
		err = invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			return nil
		}
		return xcode.FromGRPCError(err)
	}
}
