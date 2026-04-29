package cart_test

import (
	"context"
	"testing"

	cart_cache "github.com/MamangRust/monolith-ecommerce-grpc-cart/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-cart/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CartServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *CartServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Cart dependencies
	mencache := cart_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(
		queries,
		pb.NewUserQueryServiceClient(s.Conns["user"]),
		pb.NewProductQueryServiceClient(s.Conns["product"]),
	)

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *CartServiceTestSuite) TestCartServiceLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchantID, categoryID)

	// 2. Add to Cart via Service
	req := &requests.CreateCartRequest{
		UserID:    userID,
		ProductID: prodID,
		Quantity:  3,
	}

	created, err := s.svc.CartCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)

	// 3. Get Cart Items via Service
	items, _, err := s.svc.CartQuery.FindAll(ctx, &requests.FindAllCarts{
		UserID:   userID,
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotEmpty(items)
}

func TestCartServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartServiceTestSuite))
}
