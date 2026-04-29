package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-grpc-review/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-review/service"
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

	userAddr := viper.GetString("GRPC_USER_ADDR")
	
	userConn, err := grpc.NewClient(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}
	userQueryClient := pb.NewUserQueryServiceClient(userConn)

	productAddr := viper.GetString("GRPC_PRODUCT_ADDR")
	
	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	productQueryClient := pb.NewProductQueryServiceClient(productConn)

	repos := repository.NewRepositories(srv.DB, userQueryClient, productQueryClient)
	
	obs, _ := observability.NewObservability("review-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Observability: obs,
		Cache:         cache,
		Repositories:  repos,
		Logger:        srv.Logger,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterReviewQueryServiceServer(gs, h.ReviewQuery)
		pb.RegisterReviewCommandServiceServer(gs, h.ReviewCommand)
	}

	return srv, nil
}
