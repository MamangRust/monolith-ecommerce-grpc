package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CacheStats struct {
	TotalKeys       int64     `json:"total_keys"`
	HitRate         float64   `json:"hit_rate"`
	MemoryUsed      int64     `json:"memory_used"`
	MemoryUsedHuman string    `json:"memory_used_human"`
	ExpiredKeys     int64     `json:"expired_keys"`
	LastCleanupTime time.Time `json:"last_cleanup_time"`
}

type CacheStore struct {
	Redis           *redis.Client
	Logger          logger.LoggerInterface
	metrics         observability.CacheMetricsInterface
	refCount        int64
	lastCleanupTime time.Time
}

func NewCacheStore(redisClient *redis.Client, logger logger.LoggerInterface, metrics observability.CacheMetricsInterface) *CacheStore {
	return &CacheStore{
		Redis:           redisClient,
		Logger:          logger,
		metrics:         metrics,
		refCount:        0,
		lastCleanupTime: time.Now(),
	}
}

func GetFromCache[T any](ctx context.Context, store *CacheStore, key string) (*T, bool) {
	atomic.AddInt64(&store.refCount, 1)
	defer atomic.AddInt64(&store.refCount, -1)

	start := time.Now()
	defer func() {
		store.metrics.RecordCacheOperationLatency(ctx, "get", time.Since(start))
	}()

	cached, err := store.Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		store.metrics.RecordCacheMiss(ctx, key)
		return nil, false
	}
	if err != nil {
		store.Logger.Error("Redis get error", zap.Error(err), zap.String("cacheKey", key))
		store.metrics.RecordCacheError(ctx, "get", key, err)
		return nil, false
	}

	var result T
	if err := json.Unmarshal([]byte(cached), &result); err != nil {
		store.Logger.Error("Failed to unmarshal cache", zap.Error(err), zap.String("cacheKey", key))
		store.metrics.RecordCacheError(ctx, "unmarshal", key, err)
		return nil, false
	}

	store.metrics.RecordCacheHit(ctx, key)
	return &result, true
}

func SetToCache[T any](ctx context.Context, store *CacheStore, key string, data *T, expiration time.Duration) {
	atomic.AddInt64(&store.refCount, 1)
	defer atomic.AddInt64(&store.refCount, -1)

	start := time.Now()
	defer func() {
		store.metrics.RecordCacheOperationLatency(ctx, "set", time.Since(start))
	}()

	jsonData, err := json.Marshal(data)
	if err != nil {
		store.Logger.Error("Failed to marshal cache", zap.Error(err), zap.String("cacheKey", key))
		store.metrics.RecordCacheError(ctx, "marshal", key, err)
		store.metrics.RecordCacheSet(ctx, key, false)
		return
	}

	if err := store.Redis.Set(ctx, key, jsonData, expiration).Err(); err != nil {
		store.Logger.Error("Failed to set cache", zap.Error(err), zap.String("cacheKey", key))
		store.metrics.RecordCacheError(ctx, "set", key, err)
		store.metrics.RecordCacheSet(ctx, key, false)
	} else {
		store.Logger.Debug("Successfully cached data",
			zap.String("cacheKey", key),
			zap.Duration("expiration", expiration))
		store.metrics.RecordCacheSet(ctx, key, true)
	}
}

func DeleteFromCache(ctx context.Context, store *CacheStore, key string) {
	atomic.AddInt64(&store.refCount, 1)
	defer atomic.AddInt64(&store.refCount, -1)

	start := time.Now()
	defer func() {
		store.metrics.RecordCacheOperationLatency(ctx, "delete", time.Since(start))
	}()

	if err := store.Redis.Del(ctx, key).Err(); err != nil {
		store.Logger.Error("Failed to delete cache", zap.Error(err), zap.String("cacheKey", key))
		store.metrics.RecordCacheError(ctx, "delete", key, err)
		store.metrics.RecordCacheDelete(ctx, key, false)
	} else {
		store.metrics.RecordCacheDelete(ctx, key, true)
	}
}

func (store *CacheStore) ClearExpired(ctx context.Context) (int64, error) {
	atomic.AddInt64(&store.refCount, 1)
	defer atomic.AddInt64(&store.refCount, -1)

	start := time.Now()
	defer func() {
		store.metrics.RecordCacheOperationLatency(ctx, "clear_expired", time.Since(start))
	}()

	var scanned int64
	var cursor uint64
	var expiredCount int64

	for {
		keys, nextCursor, err := store.Redis.Scan(ctx, cursor, "*", 1000).Result()
		if err != nil {
			return 0, fmt.Errorf("scan error: %w", err)
		}

		scanned += int64(len(keys))

		for _, key := range keys {
			ttl, err := store.Redis.TTL(ctx, key).Result()
			if err != nil {
				store.Logger.Error("Failed to get TTL", zap.Error(err), zap.String("cacheKey", key))
				continue
			}

			if ttl == -1 {
				continue
			}

			if ttl <= 0 {
				if err := store.Redis.Del(ctx, key).Err(); err != nil {
					store.Logger.Error("Failed to delete expired key", zap.Error(err), zap.String("cacheKey", key))
				} else {
					expiredCount++
					store.Logger.Debug("Deleted expired key", zap.String("cacheKey", key))
				}
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	store.lastCleanupTime = time.Now()

	return scanned, nil
}

func (store *CacheStore) InvalidateCache(ctx context.Context, pattern string) (int64, error) {
	atomic.AddInt64(&store.refCount, 1)
	defer atomic.AddInt64(&store.refCount, -1)

	start := time.Now()
	defer func() {
		store.metrics.RecordCacheOperationLatency(ctx, "invalidate", time.Since(start))
	}()

	var deleted int64
	var cursor uint64

	for {
		keys, nextCursor, err := store.Redis.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return 0, fmt.Errorf("scan error: %w", err)
		}

		if len(keys) > 0 {
			count, err := store.Redis.Del(ctx, keys...).Result()
			if err != nil {
				return 0, fmt.Errorf("delete error: %w", err)
			}
			deleted += count
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	store.Logger.Info("Invalidated cache keys",
		zap.String("pattern", pattern),
		zap.Int64("count", deleted))

	return deleted, nil
}

func (store *CacheStore) GetStats(ctx context.Context) (*CacheStats, error) {
	atomic.AddInt64(&store.refCount, 1)
	defer atomic.AddInt64(&store.refCount, -1)

	info, err := store.Redis.Info(ctx, "memory", "keyspace").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis info: %w", err)
	}

	stats := &CacheStats{
		LastCleanupTime: store.lastCleanupTime,
	}

	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "used_memory:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				if memory, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					stats.MemoryUsed = memory
					stats.MemoryUsedHuman = formatBytes(memory)
				}
			}
			continue
		}

		if strings.HasPrefix(line, "db") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				for item := range strings.SplitSeq(parts[1], ",") {
					if strings.HasPrefix(item, "keys=") {
						if keys, err := strconv.ParseInt(strings.TrimPrefix(item, "keys="), 10, 64); err == nil {
							stats.TotalKeys += keys
						}
					}
				}
			}
		}
	}

	statsInfo, err := store.Redis.Info(ctx, "stats").Result()
	if err == nil {
		var hits, misses float64
		for line := range strings.SplitSeq(statsInfo, "\r\n") {
			if after, ok := strings.CutPrefix(line, "keyspace_hits:"); ok {
				hits, _ = strconv.ParseFloat(after, 64)
			}
			if after, ok := strings.CutPrefix(line, "keyspace_misses:"); ok {
				misses, _ = strconv.ParseFloat(after, 64)
			}
		}

		if hits+misses > 0 {
			stats.HitRate = (hits / (hits + misses)) * 100
		}
	}

	return stats, nil
}

func (store *CacheStore) GetRefCount() int64 {
	return atomic.LoadInt64(&store.refCount)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
