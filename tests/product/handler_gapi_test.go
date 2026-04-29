package product_test

import (
	"context"
	"testing"

	prod_cache "github.com/MamangRust/monolith-ecommerce-grpc-product/cache"
	prod_handler "github.com/MamangRust/monolith-ecommerce-grpc-product/handler"
	prod_repo "github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	prod_service "github.com/MamangRust/monolith-ecommerce-grpc-product/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pb.ProductQueryServiceClient
	commandClient pb.ProductCommandServiceClient
}

func (s *ProductGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Product dependencies
	mencache := prod_cache.NewMencache(cacheStore)
	repos := prod_repo.NewRepositories(
		queries,
		pb.NewCategoryQueryServiceClient(s.Conns["category"]),
		pb.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
	svc := prod_service.NewService(&prod_service.Deps{
		Cache:         mencache,
		Repository:    repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := prod_handler.NewHandler(&prod_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pb.RegisterProductQueryServiceServer(server, handler.ProductQuery)
	pb.RegisterProductCommandServiceServer(server, handler.ProductCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pb.NewProductQueryServiceClient(conn)
	s.commandClient = pb.NewProductCommandServiceClient(conn)
}

func (s *ProductGapiTestSuite) TestProductGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pb.CreateProductRequest{
		MerchantId:   int32(merchID),
		CategoryId:   int32(catID),
		Name:         "GAPI Item",
		Description:  "GAPI Description",
		Price:        1000,
		CountInStock: 10,
		Brand:        "GAPI Brand",
		Weight:       100,
		SlugProduct:  "gapi-item",
		ImageProduct: "gapi.jpg",
		Barcode:      "GAPI-123",
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	prodID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.queryClient.FindById(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)
	s.Equal("GAPI Item", getRes.Data.Name)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.Update(ctx, &pb.UpdateProductRequest{
		ProductId:    prodID,
		MerchantId:   int32(merchID),
		CategoryId:   int32(catID),
		Name:         "GAPI Item Updated",
		Description:  "Updated Description",
		Price:        2000,
		CountInStock: 20,
		Brand:        "Updated Brand",
		Weight:       200,
		SlugProduct:  "gapi-item-updated",
		ImageProduct: "gapi-updated.jpg",
		Barcode:      "GAPI-456",
	})
	s.Require().NoError(err)
	s.Equal("GAPI Item Updated", updateRes.Data.Name)

	// 7. Trash
	_, err = s.commandClient.TrashedProduct(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pb.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreProduct(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedProduct(ctx, &pb.FindByIdProductRequest{Id: prodID})
	_, err = s.commandClient.DeleteProductPermanent(ctx, &pb.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllProduct(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllProductPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestProductGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductGapiTestSuite))
}
