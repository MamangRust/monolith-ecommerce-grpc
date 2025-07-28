package middleware

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ContextMiddleware(timeout time.Duration, logger logger.LoggerInterface) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		ctx = WithRequestID(ctx, GenerateRequestID())
		ctx = WithMethod(ctx, info.FullMethod)
		ctx = WithStartTime(ctx, time.Now())

		resp, err := handler(ctx, req)

		duration := time.Since(StartTime(ctx))

		logger.Info("gRPC Request",
			zap.String("method", Method(ctx)),
			zap.String("request_id", RequestID(ctx)),
			zap.Duration("duration", duration),
			zap.Error(err),
		)

		return resp, err
	}
}
