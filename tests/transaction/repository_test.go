package transaction_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *TransactionRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupShippingAddressService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(&repository.Deps{
		DB:              queries,
		UserQuery:       pb.NewUserQueryServiceClient(s.Conns["user"]),
		MerchantQuery:   pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		OrderQuery:      pb.NewOrderQueryServiceClient(s.Conns["order"]),
		OrderItemQuery:  pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		ShippingQuery:   pb.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
	})
}

func (s *TransactionRepositoryTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, merchID, catID)
	orderID := s.SeedOrder(ctx, userID, merchID, prodID)

	// 2. Create Transaction
	status := "paid"
	req := &requests.CreateTransactionRequest{
		UserID:        userID,
		OrderID:       orderID,
		MerchantID:    merchID,
		PaymentMethod: "Credit Card",
		PaymentStatus: &status,
		Amount:        100000,
	}

	created, err := s.repo.TransactionCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(orderID), created.OrderID)

	// 3. Find by ID
	found, err := s.repo.TransactionQuery.FindByID(ctx, int(created.TransactionID))
	s.NoError(err)
	s.Equal(created.Amount, found.Amount)
}

func TestTransactionRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
