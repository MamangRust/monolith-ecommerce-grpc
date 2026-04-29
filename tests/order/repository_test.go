package order_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(&repository.Deps{
		DB:               queries,
		MerchantQuery:    pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		ProductQuery:     pb.NewProductQueryServiceClient(s.Conns["product"]),
		ProductCommand:   pb.NewProductCommandServiceClient(s.Conns["product"]),
		OrderItemQuery:   pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		OrderItemCommand: pb.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
		UserQuery:        pb.NewUserQueryServiceClient(s.Conns["user"]),
		ShippingCommand:  pb.NewShippingCommandServiceClient(s.Conns["shipping-address"]),
		TransactionCommand: pb.NewTransactionCommandServiceClient(s.Conns["transaction"]),
	})
}

func (s *OrderRepositoryTestSuite) TestOrderLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	s.SeedProduct(ctx, merchID, catID)

	req := &requests.CreateOrderRecordRequest{
		UserID:     userID,
		MerchantID: merchID,
		TotalPrice: 5000,
	}

	created, err := s.repo.OrderCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(1), created.UserID)

	// 2. Find by ID
	found, err := s.repo.OrderQuery.FindByID(ctx, int(created.OrderID))
	s.NoError(err)
	s.Equal(int32(5000), found.TotalPrice)
}

func TestOrderRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderRepositoryTestSuite))
}
