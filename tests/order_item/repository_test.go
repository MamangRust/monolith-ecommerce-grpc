package order_item_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderItemRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo *repository.Repositories
}

func (s *OrderItemRepositoryTestSuite) SetupSuite() {
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
	s.repo = repository.NewRepositories(queries)
}

func (s *OrderItemRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *OrderItemRepositoryTestSuite) TestOrderItemLifecycle() {
	ctx := context.Background()

	// Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	pageReq := &requests.FindAllOrderItems{Search: "", Page: 1, PageSize: 10}

	// 1. Create
	req := &requests.CreateOrderItemRecordRequest{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  2,
		Price:     5000,
	}
	created, err := s.repo.OrderItemCommand.Create(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(int32(orderID), created.OrderID)
	itemID := int(created.OrderItemID)

	// 2. FindOrderItemByOrder
	foundItems, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, orderID)
	s.NoError(err)
	s.NotEmpty(foundItems)

	// 3. FindAll
	all, err := s.repo.OrderItemQuery.FindAll(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(all)

	// 4. FindActive (before trash)
	active, err := s.repo.OrderItemQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(active)

	// 5. Update
	updateReq := &requests.UpdateOrderItemRecordRequest{
		OrderItemID: itemID,
		OrderID:     orderID,
		ProductID:   productID,
		Quantity:    5,
		Price:       10000,
	}
	updated, err := s.repo.OrderItemCommand.Update(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal(int32(5), updated.Quantity)

	// 6. FindOrderItemByOrder after update
	updatedItems, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, orderID)
	s.NoError(err)
	s.NotEmpty(updatedItems)

	// 7. Trash
	trashed, err := s.repo.OrderItemCommand.Trash(ctx, itemID)
	s.NoError(err)
	s.NotNil(trashed)

	// 8. FindTrashed
	trashedItems, err := s.repo.OrderItemQuery.FindTrashed(ctx, pageReq)
	s.NoError(err)
	s.NotEmpty(trashedItems)

	// 9. FindActive after trash — should NOT include
	activeAfterTrash, err := s.repo.OrderItemQuery.FindActive(ctx, pageReq)
	s.NoError(err)
	for _, item := range activeAfterTrash {
		s.NotEqual(itemID, int(item.OrderItemID))
	}

	// 10. Restore
	restored, err := s.repo.OrderItemCommand.Restore(ctx, itemID)
	s.NoError(err)
	s.NotNil(restored)

	// 11. Trash again then DeletePermanent
	_, err = s.repo.OrderItemCommand.Trash(ctx, itemID)
	s.Require().NoError(err)

	success, err := s.repo.OrderItemCommand.DeletePermanent(ctx, itemID)
	s.NoError(err)
	s.True(success)

	// 12. Create another item then DeleteByOrderIDPermanent
	second, _ := s.repo.OrderItemCommand.Create(ctx, &requests.CreateOrderItemRecordRequest{
		OrderID: orderID, ProductID: productID, Quantity: 1, Price: 1000,
	})
	s.repo.OrderItemCommand.Trash(ctx, int(second.OrderItemID))

	delByOrderSuccess, err := s.repo.OrderItemCommand.DeleteOrderItemByOrderPermanent(ctx, orderID)
	s.NoError(err)
	s.True(delByOrderSuccess)

	// 13. RestoreAll
	s.repo.OrderItemCommand.Create(ctx, &requests.CreateOrderItemRecordRequest{
		OrderID: orderID, ProductID: productID, Quantity: 1, Price: 1000,
	})
	third, _ := s.repo.OrderItemCommand.Create(ctx, &requests.CreateOrderItemRecordRequest{
		OrderID: orderID, ProductID: productID, Quantity: 2, Price: 2000,
	})
	s.repo.OrderItemCommand.Trash(ctx, int(third.OrderItemID))

	resRestore, err := s.repo.OrderItemCommand.RestoreAll(ctx)
	s.NoError(err)
	s.True(resRestore)

	// 14. DeleteAll
	s.repo.OrderItemCommand.Trash(ctx, int(third.OrderItemID))
	resDelete, err := s.repo.OrderItemCommand.DeleteAll(ctx)
	s.NoError(err)
	s.True(resDelete)
}

func TestOrderItemRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemRepositoryTestSuite))
}
