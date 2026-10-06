package resilience

import (
	"context"
	"sync"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"google.golang.org/grpc"
)

// DependencyGuardInterceptor applies a DependencyGuard (per-call deadline +
// circuit breaker + bulkhead) to every unary gRPC call made through the client
// it is attached to. Guards are keyed by the full method name
// ("/package.Service/Method") so each downstream dependency has its own breaker
// and bulkhead: one slow or failing dependency cannot exhaust the caller's
// goroutines or open the circuit for unrelated services.
//
// Usage:
//
//	conn, err := grpc.NewClient(addr,
//	    grpc.WithTransportCredentials(insecure.NewCredentials()),
//	    grpc.WithChainUnaryInterceptor(
//	        resilience.NewDependencyGuardInterceptor(srv.Logger),
//	        pkgmiddleware.TraceUnaryClientInterceptor(),
//	    ),
//	)
type DependencyGuardInterceptor struct {
	logger logger.LoggerInterface

	mu     sync.Mutex
	guards map[string]*DependencyGuard
}

// NewDependencyGuardInterceptor creates the interceptor. A nil logger is
// tolerated (replaced by the guard's no-op logger). The interceptor is safe for
// concurrent use by many goroutines / streams.
func NewDependencyGuardInterceptor(log logger.LoggerInterface) *DependencyGuardInterceptor {
	return &DependencyGuardInterceptor{
		logger: log,
		guards: make(map[string]*DependencyGuard),
	}
}

// guardFor returns (creating on first use) the guard for one downstream
// dependency. Tuning: 5 failures open the breaker for 30s, max 100 concurrent
// in-flight calls, per-call timeout 3s.
func (i *DependencyGuardInterceptor) guardFor(method string) *DependencyGuard {
	i.mu.Lock()
	defer i.mu.Unlock()

	if g, ok := i.guards[method]; ok {
		return g
	}

	g := NewDependencyGuard(method, 5, 30, 100, 3*time.Second, i.logger)
	i.guards[method] = g
	return g
}

// UnaryInterceptor returns a grpc.UnaryClientInterceptor wrapping each call in
// the dependency guard. Only transport-level failures (Unavailable,
// DeadlineExceeded, Aborted, ResourceExhausted) count against the circuit
// breaker; business errors pass through unchanged.
func (i *DependencyGuardInterceptor) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		guard := i.guardFor(method)
		return guard.Call(ctx, func(callCtx context.Context) error {
			return invoker(callCtx, method, req, reply, cc, opts...)
		})
	}
}
