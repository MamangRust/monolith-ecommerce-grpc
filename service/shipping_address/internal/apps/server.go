package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/handler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/database"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/dotenv"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-ecommerce-pkg/otel"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	port int
)

func init() {
	port = viper.GetInt("GRPC_SHIPPING_ADDRESS_ADDR")
	if port == 0 {
		port = 50063
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

type Server struct {
	Logger   logger.LoggerInterface
	DB       *db.Queries
	Services *service.Service
	Handlers *handler.Handler
	Ctx      context.Context
}

func NewServer() (*Server, func(context.Context) error, error) {
	logger, err := logger.NewLogger()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}
	flag.Parse()

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	DB := db.New(conn)

	ctx := context.Background()

	depsRepo := &repository.Deps{
		DB:  DB,
		Ctx: ctx,
	}

	repositories := repository.NewRepositories(depsRepo)

	shutdownTracerProvider, err := otel_pkg.InitTracerProvider("ShippingAddress-service", ctx)
	if err != nil {
		logger.Fatal("Failed to initialize tracer provider", zap.Error(err))
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			logger.Fatal("Failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	myredis := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST"), viper.GetString("REDIS_PORT")),
		Password:     viper.GetString("REDIS_PASSWORD"),
		DB:           viper.GetInt("REDIS_DB_USER"),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	if err := myredis.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to ping redis", zap.Error(err))
	}

	mencache := mencache.NewMencache(&mencache.Deps{
		Ctx:    ctx,
		Redis:  myredis,
		Logger: logger,
	})

	errorhandler := errorhandler.NewErrorHandler(logger)

	services := service.NewService(&service.Deps{
		Ctx:          ctx,
		ErrorHandler: errorhandler,
		Mencache:     mencache,
		Repositories: repositories,
		Logger:       logger,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
	})

	return &Server{
		Logger:   logger,
		DB:       DB,
		Services: services,
		Handlers: handlers,
		Ctx:      ctx,
	}, shutdownTracerProvider, nil
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}
	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRIC_SHIPPING_ADDRESS_ADDR"))
	metricsLis, err := net.Listen("tcp", metricsAddr)
	if err != nil {
		s.Logger.Fatal("failed to listen on", zap.Error(err))
	}

	if err != nil {
		s.Logger.Fatal("Failed to listen for metrics", zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
				otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
			),
		),
	)

	pb.RegisterShippingServiceServer(grpcServer, s.Handlers.Shipping)

	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	s.Logger.Info(fmt.Sprintf("Server running on port %d", port))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.Logger.Info("Metrics server listening on :8083")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			s.Logger.Fatal("Metrics server error", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		s.Logger.Info("gRPC server listening on :50053")
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("gRPC server error", zap.Error(err))
		}
	}()

	wg.Wait()
}
