package middleware

import (
	"context"
	"runtime/debug"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryMiddleware(logger logger.LoggerInterface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("gRPC panic recovered",
					zap.Any("panic", r),
					zap.String("method", info.FullMethod),
					zap.ByteString("stack", debug.Stack()),
				)

				err = status.Errorf(codes.Internal, "Internal server error: %v", r)
			}
		}()

		return handler(ctx, req)
	}
}
