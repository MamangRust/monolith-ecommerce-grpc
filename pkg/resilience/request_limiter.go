package resilience

import (
	"sync/atomic"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type RequestLimiter struct {
	semaphore     *semaphore.Weighted
	maxConcurrent int64
	inFlight      int64
	logger        logger.LoggerInterface
}

func NewRequestLimiter(maxConcurrent int64, logger logger.LoggerInterface) *RequestLimiter {
	if maxConcurrent <= 0 {
		logger.Warn("RequestLimiter maxConcurrent must be positive, defaulting to 100", zap.Int64("provided", maxConcurrent))
		maxConcurrent = 100
	}

	return &RequestLimiter{
		semaphore:     semaphore.NewWeighted(maxConcurrent),
		maxConcurrent: maxConcurrent,
		logger:        logger,
	}
}

func (rl *RequestLimiter) TryAcquire() bool {
	if rl.semaphore.TryAcquire(1) {
		atomic.AddInt64(&rl.inFlight, 1)
		return true
	}
	return false
}

func (rl *RequestLimiter) Release() {
	rl.semaphore.Release(1)

	newInFlight := atomic.AddInt64(&rl.inFlight, -1)

	if newInFlight < 0 {
		rl.logger.Error("RequestLimiter: Released more permits than were acquired",
			zap.Int64("in_flight", newInFlight),
			zap.Int64("max_concurrent", rl.maxConcurrent),
		)

		atomic.StoreInt64(&rl.inFlight, 0)
	}
}

func (rl *RequestLimiter) AvailablePermits() int64 {
	return rl.maxConcurrent - atomic.LoadInt64(&rl.inFlight)
}

func (rl *RequestLimiter) MaxConcurrent() int64 {
	return rl.maxConcurrent
}

func (rl *RequestLimiter) InFlight() int64 {
	return atomic.LoadInt64(&rl.inFlight)
}
