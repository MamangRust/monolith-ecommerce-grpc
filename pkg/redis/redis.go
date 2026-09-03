package redisclient

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config represents the configuration for the Redis client.
type Config struct {
	Host         string
	Port         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
}

type redisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a new Redis client using provided configuration.
func NewRedisClient(cfg *Config) *redisClient {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	return &redisClient{Client: client}
}
