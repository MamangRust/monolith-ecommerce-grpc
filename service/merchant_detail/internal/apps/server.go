package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	obs, _ := observability.NewObservability("merchant-detail-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterMerchantDetailQueryServiceServer(gs, h.MerchantDetailQuery)
		pb.RegisterMerchantDetailCommandServiceServer(gs, h.MerchantDetailCommand)
		pb.RegisterMerchantSocialCommandServiceServer(gs, h.MerchantSocialLinkCommand)
	}

	return srv, nil
}
