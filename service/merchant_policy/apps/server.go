package apps

import (
	"fmt"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/service"
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

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(
		merchantAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}

	merchantQueryClient := pb.NewMerchantQueryServiceClient(merchantConn)

	repos := repository.NewRepositories(srv.DB, merchantQueryClient)
	obs, _ := observability.NewObservability("merchant_policy-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterMerchantPolicyQueryServiceServer(gs, h.MerchantPolicyQuery)
		pb.RegisterMerchantPolicyCommandServiceServer(gs, h.MerchantPolicyCommand)
	}

	return srv, nil
}
