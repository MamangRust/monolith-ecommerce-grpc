package resilience

import (
	"context"
	"errors"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

var (
	// ErrDependencyGuardOpen is returned when the circuit breaker is open and
	// refuses to forward the call.
	ErrDependencyGuardOpen = errors.New("dependency guard: circuit breaker is open")
	// ErrDependencyGuardBulkhead is returned when the concurrent request limit
	// for this dependency has been reached.
	ErrDependencyGuardBulkhead = errors.New("dependency guard: too many concurrent requests")
)

// DependencyGuard wraps outbound calls to a remote dependency with a bulkhead
// (concurrency limit), a circuit breaker and a per-call timeout. The zero value
// is not usable; build one with NewDependencyGuard.
//
// A nil *DependencyGuard is a valid passthrough so adapters can be constructed
// without a guard (see Call).
type DependencyGuard struct {
	name    string
	limiter *RequestLimiter
	breaker *CircuitBreaker
	timeout time.Duration
	logger  logger.LoggerInterface
}

// NewDependencyGuard builds a guard for a named dependency.
//
//   - threshold: consecutive failures before the circuit breaker opens.
//   - timeoutSecs: how long the breaker stays open before allowing half-open probes.
//   - maxConcurrent: maximum concurrent in-flight calls (bulkhead).
//   - callTimeout: deadline applied to every call; <= 0 disables it.
func NewDependencyGuard(
	name string,
	threshold uint64,
	timeoutSecs uint64,
	maxConcurrent int64,
	callTimeout time.Duration,
	log logger.LoggerInterface,
) *DependencyGuard {
	return &DependencyGuard{
		name:    name,
		limiter: NewRequestLimiter(maxConcurrent, log),
		breaker: NewCircuitBreaker(threshold, timeoutSecs, log),
		timeout: callTimeout,
		logger:  log,
	}
}

// Call runs fn guarded by the bulkhead, circuit breaker and timeout.
//
// A nil receiver is a passthrough: it simply runs fn with the original context,
// which lets adapters work with no guard attached.
func (g *DependencyGuard) Call(ctx context.Context, fn func(ctx context.Context) error) error {
	return g.CallWithOutcome(ctx, fn, nil)
}

// CallWithOutcome behaves like Call, but only records a failure when isFailure
// reports true. This lets callers ignore application-level errors (for example a
// gRPC NotFound) so the circuit breaker only reacts to genuine dependency
// outages. A nil isFailure counts every error as a failure.
//
// A nil receiver is a passthrough, exactly like Call.
func (g *DependencyGuard) CallWithOutcome(ctx context.Context, fn func(ctx context.Context) error, isFailure func(error) bool) error {
	if g == nil {
		return fn(ctx)
	}

	if !g.breaker.ShouldAllowRequest() {
		g.logger.Warn("Dependency guard: circuit breaker open, rejecting call", zap.String("dependency", g.name))
		return ErrDependencyGuardOpen
	}

	if !g.limiter.TryAcquire() {
		g.logger.Warn("Dependency guard: bulkhead full, rejecting call", zap.String("dependency", g.name))
		return ErrDependencyGuardBulkhead
	}
	defer g.limiter.Release()

	callCtx := ctx
	if g.timeout > 0 {
		var cancel context.CancelFunc
		callCtx, cancel = context.WithTimeout(ctx, g.timeout)
		defer cancel()
	}

	err := fn(callCtx)
	if err != nil && (isFailure == nil || isFailure(err)) {
		g.breaker.RecordFailure()
		return err
	}

	g.breaker.RecordSuccess()
	return err
}

// Name returns the dependency name the guard was built with.
func (g *DependencyGuard) Name() string {
	if g == nil {
		return ""
	}
	return g.name
}
