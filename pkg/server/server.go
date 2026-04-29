package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/database"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/dotenv"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/middleware"
	otel_pkg "github.com/MamangRust/monolith-ecommerce-pkg/otel"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/grafana/pyroscope-go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	Logger           logger.LoggerInterface
	DB               *db.Queries
	Ctx              context.Context
	Cancel           context.CancelFunc
	CacheStore       *cache.CacheStore
	Redis            *redis.Client
	Telemetry        *otel_pkg.Telemetry
	Config           *Config
	RegisterServices func(*grpc.Server)
}

func New(cfg *Config) (*GRPCServer, error) {
	if err := dotenv.Viper(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	if err := initPyroscope(cfg); err != nil {
		log.Printf("Warning: Failed to initialize pyroscope: %v", err)
	}

	telemetry := initTelemetry(cfg)
	if err := telemetry.Init(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	cacheMetrics, err := observability.NewCacheMetrics("cache")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache metrics: %w", err)
	}

	l, err := logger.NewLogger(cfg.ServiceName, telemetry.GetLogger())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	dbConn, err := database.NewClient(l)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	queries := db.New(dbConn)

	ctx, cancel := context.WithCancel(context.Background())

	redisClient, err := initRedisServer(ctx, l, cfg.ServiceName)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	cacheStore := cache.NewCacheStore(redisClient, l, cacheMetrics)

	return &GRPCServer{
		Logger:     l,
		DB:         queries,
		Ctx:        ctx,
		Cancel:     cancel,
		CacheStore: cacheStore,
		Redis:      redisClient,
		Telemetry:  telemetry,
		Config:     cfg,
	}, nil
}

func (s *GRPCServer) Run() error {
	defer s.Cleanup()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", s.Config.Port, err)
	}

	loadMonitor := resilience.NewLoadMonitor()
	circuitBreaker := resilience.NewCircuitBreaker(100, 10, s.Logger)
	requestLimiter := resilience.NewRequestLimiter(800, s.Logger)
	resilienceHandler := middleware.NewResilienceInterceptor(loadMonitor, circuitBreaker, requestLimiter)

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(DefaultMaxConcurrentConn),
		grpc.InitialConnWindowSize(DefaultWindowSize),
		grpc.InitialWindowSize(DefaultWindowSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    DefaultKeepaliveTime,
			Timeout: DefaultKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             DefaultMinKeepaliveTime,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			middleware.ContextMiddleware(30*time.Second, s.Logger),
			middleware.RecoveryMiddleware(s.Logger),
			middleware.PyroscopeUnaryInterceptor(),
			resilienceHandler.UnaryInterceptor(),
		),
	)

	if s.RegisterServices != nil {
		s.RegisterServices(grpcServer)
	}

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	if os.Getenv("ENABLE_REFLECTION") == "true" {
		reflection.Register(grpcServer)
		s.Logger.Info("gRPC reflection enabled")
	}

	monitoringDone := s.spawnMonitoringTask()
	cleanupDone := s.spawnCleanupTask()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	errChan := make(chan error, 1)
	go func() {
		s.Logger.Info("gRPC server starting",
			zap.Int("port", s.Config.Port),
			zap.String("address", lis.Addr().String()),
		)
		if err := grpcServer.Serve(lis); err != nil {
			errChan <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	select {
	case sig := <-sigChan:
		s.Logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	case err := <-errChan:
		s.Logger.Error("Server error", zap.Error(err))
		return err
	}

	return s.gracefulShutdown(grpcServer, healthServer, monitoringDone, cleanupDone)
}

func (s *GRPCServer) gracefulShutdown(
	grpcServer *grpc.Server,
	healthServer *health.Server,
	monitoringDone, cleanupDone <-chan struct{},
) error {
	s.Logger.Info("Starting graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer shutdownCancel()

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	s.Cancel()

	tasksDone := make(chan struct{})
	go func() {
		<-monitoringDone
		<-cleanupDone
		close(tasksDone)
	}()

	select {
	case <-tasksDone:
		s.Logger.Info("Background tasks stopped successfully")
	case <-shutdownCtx.Done():
		s.Logger.Warn("Background tasks shutdown timeout, forcing stop")
	}

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.Logger.Info("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		s.Logger.Warn("Graceful shutdown timeout, forcing stop")
		grpcServer.Stop()
	}

	s.Logger.Info("Graceful shutdown completed")
	return nil
}

func (s *GRPCServer) Cleanup() {
	s.Logger.Info("Cleaning up resources...")

	if s.Redis != nil {
		if err := s.Redis.Close(); err != nil {
			s.Logger.Error("Failed to close Redis connection", zap.Error(err))
		} else {
			s.Logger.Info("Redis connection closed")
		}
	}

	if s.Telemetry != nil {
		if err := s.Telemetry.Shutdown(context.Background()); err != nil {
			s.Logger.Error("Failed to shutdown telemetry", zap.Error(err))
		} else {
			s.Logger.Info("Telemetry shutdown successfully")
		}
	}

	s.Logger.Info("Cleanup completed")
}

func initPyroscope(cfg *Config) error {
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cfg.ServiceName,
		ServerAddress:   os.Getenv("PYROSCOPE_SERVER"),
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
		Tags: map[string]string{
			"service": cfg.ServiceName,
			"env":     cfg.Environment,
			"version": cfg.ServiceVersion,
		},
	})
	return err
}

func initTelemetry(cfg *Config) *otel_pkg.Telemetry {
	samplingFraction := cfg.OtelSamplingFraction
	if samplingFraction <= 0 {
		samplingFraction = 1.0
	}

	return otel_pkg.NewTelemetry(otel_pkg.Config{
		ServiceName:            cfg.ServiceName,
		ServiceVersion:         cfg.ServiceVersion,
		Environment:            cfg.Environment,
		Endpoint:               cfg.OtelEndpoint,
		Insecure:               true,
		EnableRuntimeMetrics:   true,
		RuntimeMetricsInterval: 15 * time.Second,
		SamplingFraction:       samplingFraction,
	})
}


func initRedisServer(ctx context.Context, logger logger.LoggerInterface, serviceName string) (*redis.Client, error) {
	// Dynamically resolve Redis host/port from env based on service name if needed,
	// but for now we follow the existing pattern in server.go which often uses specific env keys.
	// We might need to standardize this.
	// Let's check how services define their redis env.

	// Defaulting to generic REDIS_HOST/PORT if not specified per service.
	// The original code used viper.GetString("REDIS_HOST_AUTH").

	// hostKey := fmt.Sprintf("REDIS_HOST_%s", os.Getenv("REDIS_SUFFIX"))
	// This might be tricky. Let's look at auth server.go again.
	// It used: fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_AUTH"), viper.GetString("REDIS_PORT_AUTH"))

	// I'll use a more flexible approach or allow passing redis config.
	// For now, let's just use what's in viper.

	// Actually, let's keep it simple for now and rely on viper being initialized by dotenv.Viper().
	// We might need to pass the specific keys.

	return redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST"), viper.GetString("REDIS_PORT")),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB"),
		DialTimeout:  RedisDialTimeout,
		ReadTimeout:  RedisReadTimeout,
		WriteTimeout: RedisWriteTimeout,
		PoolSize:     RedisPoolSize,
		MinIdleConns: RedisMinIdleConns,
	}), nil
}

func (s *GRPCServer) spawnMonitoringTask() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(MonitoringInterval)
		defer ticker.Stop()
		for {
			select {
			case <-s.Ctx.Done():
				return
			case <-ticker.C:
				s.monitorCache()
			}
		}
	}()
	return done
}

func (s *GRPCServer) monitorCache() {
	refCount := s.CacheStore.GetRefCount()
	stats, err := s.CacheStore.GetStats(s.Ctx)
	if err != nil {
		s.Logger.Error("Failed to get cache stats", zap.Error(err))
		return
	}
	logLevel := zap.InfoLevel
	if refCount > CacheRefCountThreshold {
		logLevel = zap.WarnLevel
	}
	if ce := s.Logger.Check(logLevel, "Cache statistics"); ce != nil {
		ce.Write(
			zap.Int64("ref_count", refCount),
			zap.Int64("total_keys", stats.TotalKeys),
			zap.Float64("hit_rate", stats.HitRate),
			zap.String("memory_used", stats.MemoryUsedHuman),
		)
	}
}

func (s *GRPCServer) spawnCleanupTask() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(CleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-s.Ctx.Done():
				return
			case <-ticker.C:
				s.CacheStore.ClearExpired(s.Ctx)
			}
		}
	}()
	return done
}
