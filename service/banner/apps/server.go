package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/service"
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

	observability, _ := observability.NewObservability("banner-server", srv.Logger)

	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: observability,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterBannerQueryServiceServer(gs, h.BannerQuery)
		pb.RegisterBannerCommandServiceServer(gs, h.BannerCommand)
	}

	return srv, nil
}
