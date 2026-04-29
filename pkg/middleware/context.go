package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/randomstring"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ctxKey string

const (
	requestIDKey ctxKey = "requestId"
	startTimeKey ctxKey = "startTime"
	methodKey    ctxKey = "method"
	operationKey ctxKey = "operation"
)

func GenerateRequestID() string {
	rand, _ := randomstring.GenerateRandomString(8)
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), rand)
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, startTimeKey, t)
}

func StartTime(ctx context.Context) time.Time {
	if v, ok := ctx.Value(startTimeKey).(time.Time); ok {
		return v
	}
	return time.Time{}
}

func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodKey, method)
}

func Method(ctx context.Context) string {
	if v, ok := ctx.Value(methodKey).(string); ok {
		return v
	}
	return ""
}

func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey, operation)
}

func Operation(ctx context.Context) string {
	if v, ok := ctx.Value(operationKey).(string); ok {
		return v
	}
	return ""
}

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
