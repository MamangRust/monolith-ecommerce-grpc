package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := viper.GetString("GRPC_USER_ADDR")

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")

	userConn, err := grpc.NewClient(
		userAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	productConn, err := grpc.NewClient(
		productAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	userQueryClient := pb.NewUserQueryServiceClient(userConn)
	productQueryClient := pb.NewProductQueryServiceClient(productConn)

	repos := repository.NewRepositories(srv.DB, userQueryClient, productQueryClient)

	obs, _ := observability.NewObservability("cart-service", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repositories:  repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterCartQueryServiceServer(gs, h.CartQuery)
		pb.RegisterCartCommandServiceServer(gs, h.CartCommand)
	}

	return srv, nil
}
