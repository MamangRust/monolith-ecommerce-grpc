package product_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbproduct "github.com/MamangRust/monolith-ecommerce-pb/product"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	prod_cache "github.com/MamangRust/monolith-ecommerce-product/cache"
	prod_handler "github.com/MamangRust/monolith-ecommerce-product/handler"
	prod_repo "github.com/MamangRust/monolith-ecommerce-product/repository"
	prod_service "github.com/MamangRust/monolith-ecommerce-product/service"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
)

type ProductGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient   pbproduct.ProductQueryServiceClient
	commandClient pbproduct.ProductCommandServiceClient
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
		pbcategory.NewCategoryQueryServiceClient(s.Conns["category"]),
		pbmerchant.NewMerchantQueryServiceClient(s.Conns["merchant"]),
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
	pbproduct.RegisterProductQueryServiceServer(server, handler.ProductQuery)
	pbproduct.RegisterProductCommandServiceServer(server, handler.ProductCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbproduct.NewProductQueryServiceClient(conn)
	s.commandClient = pbproduct.NewProductCommandServiceClient(conn)
}

func (s *ProductGapiTestSuite) TestProductGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create
	createRes, err := s.commandClient.Create(ctx, &pbproduct.CreateProductRequest{
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
	getRes, err := s.queryClient.FindById(ctx, &pbproduct.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)
	s.Equal("GAPI Item", getRes.Data.Name)

	// 4. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbproduct.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbproduct.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.commandClient.Update(ctx, &pbproduct.UpdateProductRequest{
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
	_, err = s.commandClient.TrashedProduct(ctx, &pbproduct.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbproduct.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.commandClient.RestoreProduct(ctx, &pbproduct.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.commandClient.TrashedProduct(ctx, &pbproduct.FindByIdProductRequest{Id: prodID})
	_, err = s.commandClient.DeleteProductPermanent(ctx, &pbproduct.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.commandClient.RestoreAllProduct(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.commandClient.DeleteAllProductPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func (s *ProductGapiTestSuite) TestProductGapiNotFound() {
	ctx := context.Background()

	// FindById with a non-existent ID must map to codes.NotFound (404), not Internal.
	_, err := s.queryClient.FindById(ctx, &pbproduct.FindByIdProductRequest{Id: 999999})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok, "expected a gRPC status error")
	s.Equal(codes.NotFound, st.Code(), "non-existent product must be NotFound, got %v: %s", st.Code(), st.Message())
}

func (s *ProductGapiTestSuite) TestProductGapiInvalidID() {
	ctx := context.Background()

	// Command mutations with id=0 must be rejected as InvalidArgument before any DB work.
	_, err := s.commandClient.TrashedProduct(ctx, &pbproduct.FindByIdProductRequest{Id: 0})
	s.Require().Error(err)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.InvalidArgument, st.Code(), "trash with id=0 must be InvalidArgument, got %v", st.Code())

	_, err = s.commandClient.Update(ctx, &pbproduct.UpdateProductRequest{
		Name:         "Invalid",
		Description:  "Invalid",
		Price:        1,
		CountInStock: 1,
	})
	s.Require().Error(err)
	st, ok = status.FromError(err)
	s.Require().True(ok)
	s.Equal(codes.InvalidArgument, st.Code(), "update with missing id must be InvalidArgument, got %v", st.Code())
}

func TestProductGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductGapiTestSuite))
}
