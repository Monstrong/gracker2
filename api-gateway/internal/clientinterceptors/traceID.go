package clientinterceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/monstrong/gracker2/pkg/logger"
)

func TraceIDClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		if traceID, ok := logger.TraceIDFromContext(ctx); ok {
			ctx = metadata.AppendToOutgoingContext(
				ctx,
				string(logger.LoggerTraceIDKey),
				traceID,
			)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
