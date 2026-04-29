package order_item_test

import (
	"context"
	"testing"

	item_cache "github.com/MamangRust/monolith-ecommerce-grpc-order-item/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderItemServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *OrderItemServiceTestSuite) SetupSuite() {
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

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Order Item dependencies
	mencache := item_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(queries)

	s.svc = service.NewService(&service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *OrderItemServiceTestSuite) TestOrderItemLifecycle() {
	ctx := context.Background()

	// 1. Setup Dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderService()

	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	categoryID := s.SeedCategory(ctx)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	// 2. Create Order Item
	req := &requests.CreateOrderItemRecordRequest{
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  2,
		Price:     5000,
	}
	created, err := s.svc.OrderItemCommand.Create(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	orderItemID := int(created.OrderItemID)

	// 3. FindByOrder
	items, err := s.svc.OrderItemQuery.FindByOrder(ctx, orderID)
	s.Require().NoError(err)
	s.NotEmpty(items)

	// 4. Update
	updateReq := &requests.UpdateOrderItemRecordRequest{
		OrderItemID: orderItemID,
		ProductID:   productID,
		Quantity:    3,
		Price:       6000,
	}
	updated, err := s.svc.OrderItemCommand.Update(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(int32(updateReq.Quantity), updated.Quantity)

	// 5. FindAll
	_, total, err := s.svc.OrderItemQuery.FindAll(ctx, &requests.FindAllOrderItems{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 6. Trash
	_, err = s.svc.OrderItemCommand.Trash(ctx, orderItemID)
	s.Require().NoError(err)

	// 7. FindTrashed
	_, totalTrashed, err := s.svc.OrderItemQuery.FindTrashed(ctx, &requests.FindAllOrderItems{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 8. FindActive
	active, _, err := s.svc.OrderItemQuery.FindActive(ctx, &requests.FindAllOrderItems{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, item := range active {
		s.NotEqual(orderItemID, int(item.OrderItemID))
	}

	// 9. Restore
	_, err = s.svc.OrderItemCommand.Restore(ctx, orderItemID)
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, err = s.svc.OrderItemCommand.Trash(ctx, orderItemID)
	s.Require().NoError(err)
	success, err := s.svc.OrderItemCommand.DeletePermanent(ctx, orderItemID)
	s.Require().NoError(err)
	s.True(success)

	// 11. RestoreAll & DeleteAll
	i1, _ := s.svc.OrderItemCommand.Create(ctx, &requests.CreateOrderItemRecordRequest{OrderID: orderID, ProductID: productID, Quantity: 1, Price: 1000})
	i2, _ := s.svc.OrderItemCommand.Create(ctx, &requests.CreateOrderItemRecordRequest{OrderID: orderID, ProductID: productID, Quantity: 1, Price: 1000})
	
	s.svc.OrderItemCommand.Trash(ctx, int(i1.OrderItemID))
	s.svc.OrderItemCommand.Trash(ctx, int(i2.OrderItemID))

	resRestoreAll, err := s.svc.OrderItemCommand.RestoreAll(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.OrderItemCommand.Trash(ctx, int(i1.OrderItemID))
	s.svc.OrderItemCommand.Trash(ctx, int(i2.OrderItemID))

	resDeleteAll, err := s.svc.OrderItemCommand.DeleteAll(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestOrderItemServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemServiceTestSuite))
}
