package resilience

import (
	"sync"
	"sync/atomic"
	"time"
)

type LoadMonitor struct {
	requestCount uint64
	lastReset    time.Time
	mu           sync.RWMutex
	window       time.Duration
	history      []LoadSnapshot
	historyMu    sync.Mutex
}

type LoadSnapshot struct {
	Timestamp time.Time
	RPS       uint64
	Count     uint64
}

func NewLoadMonitor() *LoadMonitor {
	return &LoadMonitor{
		lastReset: time.Now(),
		window:    time.Second,
		history:   make([]LoadSnapshot, 0, 60),
	}
}

func NewLoadMonitorWithWindow(window time.Duration) *LoadMonitor {
	return &LoadMonitor{
		lastReset: time.Now(),
		window:    window,
		history:   make([]LoadSnapshot, 0, 60),
	}
}

func (lm *LoadMonitor) RecordRequest() {
	atomic.AddUint64(&lm.requestCount, 1)
}

func (lm *LoadMonitor) GetCurrentRPS() uint64 {
	count := atomic.SwapUint64(&lm.requestCount, 0)

	lm.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(lm.lastReset).Seconds()
	lm.lastReset = now
	lm.mu.Unlock()

	var rps uint64
	if elapsed > 0 {
		rps = uint64(float64(count) / elapsed)
	}

	lm.historyMu.Lock()
	lm.history = append(lm.history, LoadSnapshot{
		Timestamp: now,
		RPS:       rps,
		Count:     count,
	})

	if len(lm.history) > 60 {
		lm.history = lm.history[1:]
	}
	lm.historyMu.Unlock()

	return rps
}

func (lm *LoadMonitor) GetRequestCount() uint64 {
	return atomic.LoadUint64(&lm.requestCount)
}

func (lm *LoadMonitor) GetAverageRPS(duration time.Duration) float64 {
	lm.historyMu.Lock()
	defer lm.historyMu.Unlock()

	if len(lm.history) == 0 {
		return 0
	}

	cutoff := time.Now().Add(-duration)
	var total uint64
	var count int

	for i := len(lm.history) - 1; i >= 0; i-- {
		if lm.history[i].Timestamp.Before(cutoff) {
			break
		}
		total += lm.history[i].RPS
		count++
	}

	if count == 0 {
		return 0
	}

	return float64(total) / float64(count)
}

func (lm *LoadMonitor) GetPeakRPS(duration time.Duration) uint64 {
	lm.historyMu.Lock()
	defer lm.historyMu.Unlock()

	if len(lm.history) == 0 {
		return 0
	}

	cutoff := time.Now().Add(-duration)
	var peak uint64

	for i := len(lm.history) - 1; i >= 0; i-- {
		if lm.history[i].Timestamp.Before(cutoff) {
			break
		}
		if lm.history[i].RPS > peak {
			peak = lm.history[i].RPS
		}
	}

	return peak
}

func (lm *LoadMonitor) GetHistory() []LoadSnapshot {
	lm.historyMu.Lock()
	defer lm.historyMu.Unlock()

	historyCopy := make([]LoadSnapshot, len(lm.history))
	copy(historyCopy, lm.history)
	return historyCopy
}

func (lm *LoadMonitor) Reset() {
	atomic.StoreUint64(&lm.requestCount, 0)
	lm.mu.Lock()
	lm.lastReset = time.Now()
	lm.mu.Unlock()

	lm.historyMu.Lock()
	lm.history = make([]LoadSnapshot, 0, 60)
	lm.historyMu.Unlock()
}
