package middleware

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResilienceInterceptor struct {
	LoadMonitor    *resilience.LoadMonitor
	CircuitBreaker *resilience.CircuitBreaker
	RequestLimiter *resilience.RequestLimiter
}

func NewResilienceInterceptor(lm *resilience.LoadMonitor, cb *resilience.CircuitBreaker, rl *resilience.RequestLimiter) *ResilienceInterceptor {
	return &ResilienceInterceptor{
		LoadMonitor:    lm,
		CircuitBreaker: cb,
		RequestLimiter: rl,
	}
}

func (im *ResilienceInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		im.LoadMonitor.RecordRequest()

		if im.CircuitBreaker.IsOpen() {
			return nil, status.Error(codes.Unavailable, "Service temporarily unavailable due to high error rate. Please try again later.")
		}

		if !im.RequestLimiter.TryAcquire() {
			return nil, status.Error(codes.ResourceExhausted, "Too many concurrent requests. Please try again later.")
		}

		defer im.RequestLimiter.Release()

		return handler(ctx, req)
	}
}
