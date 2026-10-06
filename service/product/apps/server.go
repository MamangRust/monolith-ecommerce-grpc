package apps

import (
	"fmt"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-pkg/server"
	"github.com/MamangRust/monolith-ecommerce-product/cache"
	"github.com/MamangRust/monolith-ecommerce-product/handler"
	"github.com/MamangRust/monolith-ecommerce-product/repository"
	"github.com/MamangRust/monolith-ecommerce-product/service"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	categoryAddr := viper.GetString("GRPC_CATEGORY_ADDR")

	categoryConn, err := grpc.NewClient(categoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to category service: %w", err)
	}
	categoryQueryClient := pbcategory.NewCategoryQueryServiceClient(categoryConn)

	merchantAddr := viper.GetString("GRPC_MERCHANT_ADDR")

	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to merchant service: %w", err)
	}
	merchantQueryClient := pbmerchant.NewMerchantQueryServiceClient(merchantConn)

	guardCategory := resilience.NewDependencyGuard("category", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(srv.DB,
		categoryQueryClient,
		merchantQueryClient,
		repository.GuardOptions{
			Category: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardCategory),
			},
			Merchant: []adapter.GuardOption{
				adapter.WithDependencyGuard(guardMerchant),
			},
		},
	)
	obs, _ := observability.NewObservability("product-server", srv.Logger)
	cache := cache.NewMencache(srv.CacheStore)

	svc := service.NewService(&service.Deps{
		Cache:         cache,
		Logger:        srv.Logger,
		Repository:    repos,
		Observability: obs,
	})

	h := handler.NewHandler(&handler.Deps{Service: svc, Logger: srv.Logger})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbproduct.RegisterProductQueryServiceServer(gs, h.ProductQuery)
		pbproduct.RegisterProductCommandServiceServer(gs, h.ProductCommand)
	}

	return srv, nil
}
