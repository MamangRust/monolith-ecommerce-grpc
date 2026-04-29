package order_test

import (
	"context"
	"testing"

	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	order_handler "github.com/MamangRust/monolith-ecommerce-grpc-order/handler"
	order_repo "github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	order_service "github.com/MamangRust/monolith-ecommerce-grpc-order/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pb.OrderQueryServiceClient
	commandClient pb.OrderCommandServiceClient
}

func (s *OrderGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	
	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Order dependencies
	mencache := order_cache.NewMencache(cacheStore)
	repos := order_repo.NewRepositories(&order_repo.Deps{
		DB:               queries,
		MerchantQuery:    pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		ProductQuery:     pb.NewProductQueryServiceClient(s.Conns["product"]),
		ProductCommand:   pb.NewProductCommandServiceClient(s.Conns["product"]),
		OrderItemQuery:   pb.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		OrderItemCommand: pb.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
		UserQuery:        pb.NewUserQueryServiceClient(s.Conns["user"]),
		ShippingCommand:  pb.NewShippingCommandServiceClient(s.Conns["shipping-address"]),
		ShippingQuery:    pb.NewShippingQueryServiceClient(s.Conns["shipping-address"]),
		TransactionCommand: pb.NewTransactionCommandServiceClient(s.Conns["transaction"]),
	})
	svc := order_service.NewService(&order_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := order_handler.NewHandler(&order_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterOrderQueryServiceServer(server, handler.OrderQuery)
	pb.RegisterOrderCommandServiceServer(server, handler.OrderCommand)
	
	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pb.NewOrderQueryServiceClient(conn)
	s.commandClient = pb.NewOrderCommandServiceClient(conn)
}

func (s *OrderGapiTestSuite) TestOrderGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pb.CreateOrderRequest{
		UserId:     int32(userID),
		MerchantId: int32(merchID),
		TotalPrice: 10000,
		Items: []*pb.CreateOrderItemRequest{
			{
				ProductId: int32(prodID),
				Quantity:  1,
				Price:     10000,
			},
		},
		Shipping: &pb.CreateShippingAddressRequest{
			Alamat:          "Test Address",
			Provinsi:        "Test Province",
			Kota:            "Test City",
			Negara:          "Test Country",
			Courier:         "Test Courier",
			ShippingMethod:  "Test Method",
			ShippingCost:    1000,
		},
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	orderID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)
	s.Equal(int32(userID), getRes.Data.UserId)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pb.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pb.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	// Fetch terms first
	itemClient := pb.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	itemsRes, err := itemClient.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{Id: orderID})
	s.Require().NoError(err)
	s.NotEmpty(itemsRes.Data)
	orderItemID := itemsRes.Data[0].Id

	_, err = s.commandClient.Update(ctx, &pb.UpdateOrderRequest{
		OrderId:    orderID,
		UserId:     int32(userID),
		TotalPrice: 15000,
		Items: []*pb.UpdateOrderItemRequest{
			{
				OrderItemId: orderItemID,
				ProductId:   int32(prodID),
				Quantity:    1,
				Price:       15000,
			},
		},
		Shipping: &pb.UpdateShippingAddressRequest{
			Alamat:          "Updated Address",
			Provinsi:        "Updated Province",
			Kota:            "Updated City",
			Negara:          "Updated Country",
			Courier:         "Updated Courier",
			ShippingMethod:  "Updated Method",
			ShippingCost:    1500,
		},
	})
	s.Require().NoError(err)

	// 7. Trash
	_, err = s.commandClient.TrashedOrder(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pb.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreOrder(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedOrder(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	_, err = s.commandClient.DeleteOrderPermanent(ctx, &pb.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllOrder(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllOrderPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestOrderGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderGapiTestSuite))
}
