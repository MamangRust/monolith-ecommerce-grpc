package order_item_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderItemRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *OrderItemRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies (optional for repo, but good for real IDs)
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupShippingAddressService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	queries := db.New(s.DBPool())
	s.repo = repository.NewRepositories(queries)
}

func (s *OrderItemRepositoryTestSuite) TestOrderItemLifecycle() {
	ctx := context.Background()

	// Seed all dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	// Create Order Item
	req := &requests.CreateOrderItemRecordRequest{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  2,
		Price:     1000,
	}

	created, err := s.repo.OrderItemCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(orderID), created.OrderID)

	// 3. Find by Order ID
	found, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, orderID)
	s.NoError(err)
	s.NotEmpty(found)
	var foundItem *db.GetOrderItemsByOrderRow
	for _, item := range found {
		if item.Quantity == int32(req.Quantity) {
			foundItem = item
			break
		}
	}
	s.NotNil(foundItem)
	s.Equal(int32(req.Quantity), foundItem.Quantity)

	// 4. Find by Order ID
	items, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, orderID)
	s.NoError(err)
	s.NotEmpty(items)
}

func TestOrderItemRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemRepositoryTestSuite))
}
