package apps

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/handler"
	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-auth/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/database"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/dotenv"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-ecommerce-pkg/otel"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	port = viper.GetInt("GRPC_AUTH_ADDR")
	if port == 0 {
		port = 50053
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

type Server struct {
	Logger       logger.LoggerInterface
	DB           *db.Queries
	TokenManager *auth.Manager
	Services     *service.Service
	Handlers     *handler.Handler
	Ctx          context.Context
}

func NewServer() (*Server, error) {
	flag.Parse()

	logger, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		logger.Fatal("Failed to create token manager", zap.Error(err))
	}

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	DB := db.New(conn)

	ctx := context.Background()
	hash := hash.NewHashingPassword()
	mapperRecord := recordmapper.NewRecordMapper()
	mapperResponse := response_service.NewResponseServiceMapper()

	depsRepo := repository.Deps{
		DB:           DB,
		Ctx:          ctx,
		MapperRecord: mapperRecord,
	}
	repositories := repository.NewRepositories(depsRepo)

	kafka := kafka.NewKafka(logger, []string{viper.GetString("KAFKA_BROKERS")})

	shutdownTracerProvider, err := otel_pkg.InitTracerProvider("Auth-service", ctx)
	if err != nil {
		logger.Fatal("Failed to initialize tracer provider", zap.Error(err))
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			logger.Fatal("Failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	services := service.NewService(service.Deps{
		Context:      ctx,
		Repositories: repositories,
		Hash:         hash,
		Token:        tokenManager,
		Logger:       logger,
		Mapper:       mapperResponse.UserResponseMapper,
		Kafka:        *kafka,
	})

	handlers := handler.NewHandler(handler.Deps{
		Service: *services,
	})

	return &Server{
		Logger:       logger,
		DB:           DB,
		TokenManager: tokenManager,
		Services:     services,
		Handlers:     handlers,
		Ctx:          ctx,
	}, nil
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRICS_AUTH_ADDR"))
	metricsLis, err := net.Listen("tcp", metricsAddr)
	if err != nil {
		s.Logger.Fatal("failed to listen on", zap.Error(err))
	}

	if err != nil {
		log.Fatalf("Failed to listen for metrics: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
				otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
			),
		),
	)

	pb.RegisterAuthServiceServer(grpcServer, s.Handlers.Auth)

	metricsServer := http.NewServeMux()
	metricsServer.Handle("/metrics", promhttp.Handler())

	s.Logger.Info(fmt.Sprintf("Server running on port %d", port))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.Logger.Info("Metrics server listening on :8081")
		if err := http.Serve(metricsLis, metricsServer); err != nil {
			s.Logger.Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		s.Logger.Debug("gRPC server listening on port", zap.Int("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()

	wg.Wait()
}
