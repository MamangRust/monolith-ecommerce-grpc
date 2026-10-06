package category_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	cat_cache "github.com/MamangRust/monolith-ecommerce-category/cache"
	cat_handler "github.com/MamangRust/monolith-ecommerce-category/handler"
	cat_repo "github.com/MamangRust/monolith-ecommerce-category/repository"
	cat_service "github.com/MamangRust/monolith-ecommerce-category/service"
	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	tests "github.com/MamangRust/monolith-ecommerce-test"
)

type CategoryGapiTestSuite struct {
	tests.BaseTestSuite
	queryClient           pbcategory.CategoryQueryServiceClient
	commandClient         pbcategory.CategoryCommandServiceClient
	statsClient           pbcategory.CategoryStatsServiceClient
	statsByIdClient       pbcategory.CategoryStatsByIdServiceClient
	statsByMerchantClient pbcategory.CategoryStatsByMerchantServiceClient
}

func (s *CategoryGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	queries := db.New(s.DBPool())

	// Category dependencies
	mencache := cat_cache.NewMencache(cacheStore)
	repos := cat_repo.NewRepositories(queries)
	svc := cat_service.NewService(&cat_service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := cat_handler.NewHandler(&cat_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// GRPC Server
	server := grpc.NewServer()
	pbcategory.RegisterCategoryQueryServiceServer(server, handler.CategoryQuery)
	pbcategory.RegisterCategoryCommandServiceServer(server, handler.CategoryCommand)
	pbcategory.RegisterCategoryStatsServiceServer(server, handler.CategoryStats)
	pbcategory.RegisterCategoryStatsByIdServiceServer(server, handler.CategoryStatsById)
	pbcategory.RegisterCategoryStatsByMerchantServiceServer(server, handler.CategoryStatsByMerchant)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.queryClient = pbcategory.NewCategoryQueryServiceClient(conn)
	s.commandClient = pbcategory.NewCategoryCommandServiceClient(conn)
	s.statsClient = pbcategory.NewCategoryStatsServiceClient(conn)
	s.statsByIdClient = pbcategory.NewCategoryStatsByIdServiceClient(conn)
	s.statsByMerchantClient = pbcategory.NewCategoryStatsByMerchantServiceClient(conn)
}

func (s *CategoryGapiTestSuite) TestCategoryGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.commandClient.Create(ctx, &pbcategory.CreateCategoryRequest{
		Name:          "GAPI Category",
		Description:   "Testing via GRPC",
		SlugCategory:  "gapi-cat",
		ImageCategory: "gapi.jpg",
	})
	s.NoError(err)
	s.NotNil(createRes)
	catID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)
	s.Equal("GAPI Category", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.commandClient.Update(ctx, &pbcategory.UpdateCategoryRequest{
		CategoryId:    catID,
		Name:          "GAPI Category Updated",
		Description:   "Updated via GRPC",
		SlugCategory:  "gapi-cat-updated",
		ImageCategory: "gapi-updated.jpg",
	})
	s.NoError(err)
	s.Equal("GAPI Category Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.commandClient.TrashedCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbcategory.FindAllCategoryRequest{Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.commandClient.RestoreCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 9. DeletePermanent
	_, _ = s.commandClient.TrashedCategory(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	_, err = s.commandClient.DeleteCategoryPermanent(ctx, &pbcategory.FindByIdCategoryRequest{Id: catID})
	s.NoError(err)

	// 10. RestoreAll
	_, err = s.commandClient.RestoreAllCategory(ctx, &emptypb.Empty{})
	s.NoError(err)

	// 11. DeleteAll
	_, err = s.commandClient.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})
	s.NoError(err)
}

func (s *CategoryGapiTestSuite) TestGapiStats() {
	ctx := context.Background()

	_, err := s.statsClient.FindMonthlyTotalPrices(ctx, &pbcategory.FindYearMonthTotalPrices{Year: 2024, Month: 4})
	s.NoError(err)

	_, err = s.statsClient.FindYearlyTotalPrices(ctx, &pbcategory.FindYearTotalPrices{Year: 2024})
	s.NoError(err)

	_, err = s.statsByIdClient.FindMonthlyTotalPricesById(ctx, &pbcategory.FindYearMonthTotalPriceById{CategoryId: 1, Year: 2024, Month: 4})
	s.NoError(err)
}

func TestCategoryGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryGapiTestSuite))
}
