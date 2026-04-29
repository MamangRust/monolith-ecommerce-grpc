package resilience

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

type CircuitBreaker struct {
	failureCount         uint64
	successCount         uint64
	isOpen               uint32
	lastFailureTime      time.Time
	lastFailureTimeMutex sync.RWMutex
	threshold            uint64
	timeout              time.Duration
	halfOpenMaxRequests  uint64
	logger               logger.LoggerInterface
}

func NewCircuitBreaker(threshold uint64, timeoutSecs uint64, logger logger.LoggerInterface) *CircuitBreaker {
	return &CircuitBreaker{
		threshold:           threshold,
		timeout:             time.Duration(timeoutSecs) * time.Second,
		halfOpenMaxRequests: 5,
		logger:              logger,
	}
}

func (cb *CircuitBreaker) IsOpen() bool {
	return atomic.LoadUint32(&cb.isOpen) == 1
}

func (cb *CircuitBreaker) GetFailureCount() uint64 {
	return atomic.LoadUint64(&cb.failureCount)
}

func (cb *CircuitBreaker) GetSuccessCount() uint64 {
	return atomic.LoadUint64(&cb.successCount)
}

func (cb *CircuitBreaker) ShouldAllowRequest() bool {
	if !cb.IsOpen() {
		return true
	}

	cb.lastFailureTimeMutex.RLock()
	lastTime := cb.lastFailureTime
	cb.lastFailureTimeMutex.RUnlock()

	if time.Since(lastTime) > cb.timeout {
		cb.logger.Info("Circuit breaker timeout passed, entering half-open state")

		success := atomic.LoadUint64(&cb.successCount)
		if success < cb.halfOpenMaxRequests {
			return true
		}

		if success >= cb.halfOpenMaxRequests {
			cb.CloseCircuit()
			return true
		}

		return false
	}

	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	success := atomic.AddUint64(&cb.successCount, 1)

	if cb.IsOpen() && success >= cb.halfOpenMaxRequests {
		cb.CloseCircuit()
	}

	atomic.StoreUint64(&cb.failureCount, 0)
}

func (cb *CircuitBreaker) RecordFailure() {
	count := atomic.AddUint64(&cb.failureCount, 1)

	if count >= cb.threshold && !cb.IsOpen() {
		cb.OpenCircuit()
	}
}

func (cb *CircuitBreaker) OpenCircuit() {
	cb.logger.Warn("Circuit breaker opened", zap.Uint64("failures", cb.GetFailureCount()))
	atomic.StoreUint32(&cb.isOpen, 1)
	atomic.StoreUint64(&cb.successCount, 0)

	cb.lastFailureTimeMutex.Lock()
	cb.lastFailureTime = time.Now()
	cb.lastFailureTimeMutex.Unlock()
}

func (cb *CircuitBreaker) CloseCircuit() {
	cb.logger.Info("Circuit breaker closed - service recovered")
	atomic.StoreUint32(&cb.isOpen, 0)
	atomic.StoreUint64(&cb.failureCount, 0)
	atomic.StoreUint64(&cb.successCount, 0)
}

func (cb *CircuitBreaker) Reset() {
	cb.CloseCircuit()
}
