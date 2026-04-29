package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/service"
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
		pb.RegisterSliderQueryServiceServer(gs, h.SliderQuery)
		pb.RegisterSliderCommandServiceServer(gs, h.SliderCommand)
	}

	return srv, nil
}
