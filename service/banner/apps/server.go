package apps

import (
	"github.com/MamangRust/monolith-ecommerce-banner/cache"
	"github.com/MamangRust/monolith-ecommerce-banner/handler"
	"github.com/MamangRust/monolith-ecommerce-banner/repository"
	"github.com/MamangRust/monolith-ecommerce-banner/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"google.golang.org/grpc"

	pbbanner "github.com/MamangRust/monolith-ecommerce-pb/banner"
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
		pbbanner.RegisterBannerQueryServiceServer(gs, h.BannerQuery)
		pbbanner.RegisterBannerCommandServiceServer(gs, h.BannerCommand)
	}

	return srv, nil
}
