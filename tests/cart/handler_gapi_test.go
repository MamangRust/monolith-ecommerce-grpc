package cart_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	cart_cache "github.com/MamangRust/monolith-ecommerce-cart/cache"
	cart_handler "github.com/MamangRust/monolith-ecommerce-cart/handler"
	cart_repo "github.com/MamangRust/monolith-ecommerce-cart/repository"
	cart_service "github.com/MamangRust/monolith-ecommerce-cart/service"
	pbcart "github.com/MamangRust/monolith-ecommerce-pb/cart"
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
)

type CartGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbcart.CartQueryServiceClient
	commandClient pbcart.CartCommandServiceClient
}

func (s *CartGapiTestSuite) SetupSuite() {
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
	repos := cart_repo.NewRepositories(
		queries,
		pbuser.NewUserQueryServiceClient(s.Conns["user"]),
		pbproduct.NewProductQueryServiceClient(s.Conns["product"]),
	)
	svc := cart_service.NewService(&cart_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := cart_handler.NewHandler(&cart_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbcart.RegisterCartQueryServiceServer(server, handler.CartQuery)
	pbcart.RegisterCartCommandServiceServer(server, handler.CartCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbcart.NewCartQueryServiceClient(conn)
	s.commandClient = pbcart.NewCartCommandServiceClient(conn)
}

func (s *CartGapiTestSuite) TestGapiLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchantID, categoryID)

	// Add to Cart
	createRes, err := s.commandClient.Create(ctx, &pbcart.CreateCartRequest{
		UserId:    int32(userID),
		ProductId: int32(prodID),
		Quantity:  5,
	})
	s.Require().NoError(err)
	s.NotNil(createRes)

	// Get
	listRes, err := s.queryClient.FindAll(ctx, &pbcart.FindAllCartRequest{UserId: int32(userID)})
	s.Require().NoError(err)
	s.NotEmpty(listRes.Data)
}

func TestCartGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CartGapiTestSuite))
}
