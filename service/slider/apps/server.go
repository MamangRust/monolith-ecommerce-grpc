package apps

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-slider/cache"
	"github.com/MamangRust/monolith-ecommerce-slider/handler"
	"github.com/MamangRust/monolith-ecommerce-slider/repository"
	"github.com/MamangRust/monolith-ecommerce-slider/service"
	"google.golang.org/grpc"

	pbslider "github.com/MamangRust/monolith-ecommerce-pb/slider"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	cache := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability("slider-server", srv.Logger)

	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Mencache:      cache,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbslider.RegisterSliderQueryServiceServer(gs, h.SliderQuery)
		pbslider.RegisterSliderCommandServiceServer(gs, h.SliderCommand)
	}

	return srv, nil
}
