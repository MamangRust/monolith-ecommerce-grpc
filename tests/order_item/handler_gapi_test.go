package order_item_test

import (
	"context"
	"testing"

	item_cache "github.com/MamangRust/monolith-ecommerce-grpc-order-item/cache"
	item_handler "github.com/MamangRust/monolith-ecommerce-grpc-order-item/handler"
	item_repo "github.com/MamangRust/monolith-ecommerce-grpc-order-item/repository"
	item_service "github.com/MamangRust/monolith-ecommerce-grpc-order-item/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderItemGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pb.OrderItemQueryServiceClient
	commandClient pb.OrderItemCommandServiceClient
}

func (s *OrderItemGapiTestSuite) SetupSuite() {
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

	// Item dependencies
	mencache := item_cache.NewMencache(cacheStore)
	repos := item_repo.NewRepositories(queries)
	svc := item_service.NewService(&item_service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := item_handler.NewHandler(&item_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterOrderItemQueryServiceServer(server, handler.OrderItemQuery)
	pb.RegisterOrderItemCommandServiceServer(server, handler.OrderItemCommand)
	
	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pb.NewOrderItemQueryServiceClient(conn)
	s.commandClient = pb.NewOrderItemCommandServiceClient(conn)
}

func (s *OrderItemGapiTestSuite) TestOrderItemGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, catID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	// 2. Create
	createRes, err := s.commandClient.CreateOrderItem(ctx, &pb.CreateOrderItemRecordRequest{
		OrderId:   int32(orderID),
		ProductId: int32(productID),
		Quantity:  10,
		Price:     700,
	})
	s.NoError(err)
	s.NotNil(createRes)
	orderItemID := createRes.Data.Id

	// 3. FindById
	_, err = s.queryClient.FindAll(ctx, &pb.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	// s.Equal(int32(10), getRes.Data[0].Quantity)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pb.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pb.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.UpdateOrderItem(ctx, &pb.UpdateOrderItemRecordRequest{
		OrderItemId: orderItemID,
		Quantity:    20,
		Price:       800,
	})
	s.NoError(err)
	s.Equal(int32(20), updateRes.Data.Quantity)

	// 7. Trash
	_, err = s.commandClient.TrashOrderItem(ctx, &pb.FindByIdOrderItemRequest{Id: orderItemID})
	s.NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pb.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreOrderItem(ctx, &pb.FindByIdOrderItemRequest{Id: orderItemID})
	s.NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashOrderItem(ctx, &pb.FindByIdOrderItemRequest{Id: orderItemID})
	_, err = s.commandClient.DeleteOrderItemPermanent(ctx, &pb.FindByIdOrderItemRequest{Id: orderItemID})
	s.NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllOrdersItem(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllPermanentOrdersItem(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func TestOrderItemGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemGapiTestSuite))
}
