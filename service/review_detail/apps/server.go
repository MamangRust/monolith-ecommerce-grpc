package apps

import (
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/service"
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

	obs, _ := observability.NewObservability("review_detail-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Observability: obs,
		Cache:         cache,
		Repositories:  repos,
		Logger:        srv.Logger,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterReviewDetailQueryServiceServer(gs, h.ReviewDetailQuery)
		pb.RegisterReviewDetailCommandServiceServer(gs, h.ReviewDetailCommand)
	}

	return srv, nil
}
