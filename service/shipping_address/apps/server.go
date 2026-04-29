package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"
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
	mencache := cache.NewMencache(srv.CacheStore)
	obs, _ := observability.NewObservability("shipping_address-server", srv.Logger)

	svc := service.NewService(&service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{
		Service: svc,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterShippingQueryServiceServer(gs, h.ShippingQuery)
		pb.RegisterShippingCommandServiceServer(gs, h.ShippingCommand)
	}

	return srv, nil
}
