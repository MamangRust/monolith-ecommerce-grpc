package shipping_address_test

import (
	"context"
	"testing"

	ship_cache "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/cache"
	ship_handler "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/handler"
	ship_repo "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	ship_service "github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ShippingAddressGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pb.ShippingQueryServiceClient
	commandClient pb.ShippingCommandServiceClient
}

func (s *ShippingAddressGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupOrderService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Shipping dependencies
	mencache := ship_cache.NewMencache(cacheStore)
	repos := ship_repo.NewRepositories(queries)
	svc := ship_service.NewService(&ship_service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := ship_handler.NewHandler(&ship_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterShippingQueryServiceServer(server, handler.ShippingQuery)
	pb.RegisterShippingCommandServiceServer(server, handler.ShippingCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pb.NewShippingQueryServiceClient(conn)
	s.commandClient = pb.NewShippingCommandServiceClient(conn)
}

func (s *ShippingAddressGapiTestSuite) TestShippingAddressGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	categoryID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	// 2. Create
	createRes, err := s.commandClient.CreateShipping(ctx, &pb.CreateShippingAddressRequest{
		OrderId:        int32(orderID),
		Alamat:         "GAPI Home",
		Provinsi:       "GAPI Province",
		Negara:         "GAPI Country",
		Kota:           "GAPI City",
		Courier:        "GAPI Courier",
		ShippingMethod: "GAPI Method",
		ShippingCost:   1000,
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	addrID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pb.FindByIdShippingRequest{Id: addrID})
	s.Require().NoError(err)
	s.Equal("GAPI Home", getRes.Data.Alamat)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pb.FindAllShippingRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pb.FindAllShippingRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.UpdateShipping(ctx, &pb.UpdateShippingAddressRequest{
		ShippingId:     addrID,
		OrderId:        int32(orderID),
		Alamat:         "GAPI Home Updated",
		Provinsi:       "GAPI Province Updated",
		Negara:         "GAPI Country Updated",
		Kota:           "GAPI City Updated",
		Courier:        "GAPI Courier Updated",
		ShippingMethod: "GAPI Method Updated",
		ShippingCost:   2000,
	})
	s.Require().NoError(err)
	s.Equal("GAPI Home Updated", updateRes.Data.Alamat)

	// 7. Trash
	_, err = s.commandClient.TrashedShipping(ctx, &pb.FindByIdShippingRequest{Id: addrID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pb.FindAllShippingRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreShipping(ctx, &pb.FindByIdShippingRequest{Id: addrID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedShipping(ctx, &pb.FindByIdShippingRequest{Id: addrID})
	_, err = s.commandClient.DeleteShippingPermanent(ctx, &pb.FindByIdShippingRequest{Id: addrID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllShipping(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllShippingPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestShippingAddressGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ShippingAddressGapiTestSuite))
}
