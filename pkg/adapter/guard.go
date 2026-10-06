package adapter

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
)

// DefaultGuard builds a dependency guard with the project-wide defaults:
// 5 consecutive failures trip the circuit breaker, which stays open for 30s,
// at most 100 concurrent in-flight calls, and a 3s per-call deadline.
func DefaultGuard(name string, log logger.LoggerInterface) *resilience.DependencyGuard {
	return resilience.NewDependencyGuard(name, 5, 30, 100, 3*time.Second, log)
}

// GuardSetter is implemented by every gRPC repository so WithDependencyGuard
// can attach a guard without changing each constructor's client parameter.
type GuardSetter interface {
	SetGuard(*resilience.DependencyGuard)
}

// GuardOption configures a GuardSetter (a gRPC-backed repository). Repositories
// accept zero or more options at construction; the guard wraps every outbound
// gRPC call with a per-call deadline, bulkhead and circuit breaker.
type GuardOption func(GuardSetter)

// WithDependencyGuard returns a GuardOption that attaches a dependency guard
// (per-call timeout + circuit breaker + bulkhead) to a gRPC repository.
// Passing nil disables guarding.
//
// Usage:
//
//	guard := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)
//	repos := repository.NewRepositories(
//		db,
//		merchantadapter.NewQueryAdapter(merchantClient, adapter.WithDependencyGuard(guard)),
//	)
func WithDependencyGuard(guard *resilience.DependencyGuard) GuardOption {
	return func(s GuardSetter) {
		s.SetGuard(guard)
	}
}

// Call runs fn under g (bulkhead + circuit breaker + per-call timeout) and
// returns its value. A nil guard runs fn with the original context, so adapters
// keep working when no guard is attached.
func Call[T any](g *resilience.DependencyGuard, ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	var out T
	err := g.Call(ctx, func(ctx context.Context) error {
		v, err := fn(ctx)
		if err != nil {
			return err
		}
		out = v
		return nil
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return out, nil
}
