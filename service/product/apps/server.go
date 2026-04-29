package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-product/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-product/service"
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

	categoryAddr := viper.GetString("GRPC_CATEGORY_ADDR")

	categoryConn, err := grpc.NewClient(categoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to category service: %w", err)
	}
	categoryQueryClient := pb.NewCategoryQueryServiceClient(categoryConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pb.NewMerchantQueryServiceClient(merchantConn)

	repos := repository.NewRepositories(srv.DB, categoryQueryClient, merchantQueryClient)
	obs, _ := observability.NewObservability("product-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterProductQueryServiceServer(gs, h.ProductQuery)
		pb.RegisterProductCommandServiceServer(gs, h.ProductCommand)
	}

	return srv, nil
}
