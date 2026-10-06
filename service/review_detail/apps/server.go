package apps

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-review-detail/cache"
	"github.com/MamangRust/monolith-ecommerce-review-detail/handler"
	"github.com/MamangRust/monolith-ecommerce-review-detail/repository"
	"github.com/MamangRust/monolith-ecommerce-review-detail/service"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"google.golang.org/grpc"

	pbreview_detail "github.com/MamangRust/monolith-ecommerce-pb/review_detail"
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
		pbreview_detail.RegisterReviewDetailQueryServiceServer(gs, h.ReviewDetailQuery)
		pbreview_detail.RegisterReviewDetailCommandServiceServer(gs, h.ReviewDetailCommand)
	}

	return srv, nil
}
