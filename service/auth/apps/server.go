package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-auth/cache"
	"github.com/MamangRust/monolith-ecommerce-auth/handler"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-ecommerce-auth/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}


	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	roleAddr := viper.GetString("GRPC_ROLE_ADDR")


	userAddr := viper.GetString("GRPC_USER_ADDR")

	roleConn, err := grpc.NewClient(
		roleAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to role service: %w", err)
	}

	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	roleQueryClient := pb.NewRoleQueryServiceClient(roleConn)
	roleCommandClient := pb.NewRoleCommandServiceClient(roleConn)
	userQueryClient := pb.NewUserQueryServiceClient(userConn)
	userCommandClient := pb.NewUserCommandServiceClient(userConn)

	hasher := hash.NewHashingPassword()
	repositories := repository.NewRepositories(srv.DB, userQueryClient, userCommandClient, roleQueryClient, roleCommandClient)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	observability, _ := observability.NewObservability("auth-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      cache,
		Repositories:  repositories,
		Token:         tokenManager,
		Hash:          hasher,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Observability: observability,
	})

	handlers := handler.NewHandler(&handler.Deps{Service: services, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	return srv, nil
}
