package category_test

import (
	"context"
	"testing"
	"time"

	category_cache "github.com/MamangRust/monolith-ecommerce-grpc-category/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CategoryStatsServiceTestSuite struct {
	tests.BaseTestSuite
	svc             service.CategoryStatsService
	svcById         service.CategoryStatsByIdService
	svcByMerchant   service.CategoryStatsByMerchantService
	categoryID      int
	merchantID      int
	userID          int
}

func (s *CategoryStatsServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupShippingAddressService()
	s.SetupTransactionService()
	s.SetupOrderService()

	queries := db.New(s.DBPool())
	repos := repository.NewRepositories(queries)
	cacheStore := s.GetCacheStore()
	mencache := category_cache.NewMencache(cacheStore)

	deps := &service.Deps{
		Cache:         mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	}
	allServices := service.NewService(deps)
	s.svc = allServices.CategoryStats
	s.svcById = allServices.CategoryStatsById
	s.svcByMerchant = allServices.CategoryStatsByMerchant

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.categoryID = s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, s.categoryID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	_, err := s.DBPool().Exec(ctx, "UPDATE orders SET created_at = $1 WHERE order_id = $2", 
		time.Now(), orderID)
	s.Require().NoError(err)
}

func (s *CategoryStatsServiceTestSuite) TestFindMonthlyTotalPrice() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthTotalPrice{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	// Test first call (repository)
	res, err := s.svc.FindMonthlyTotalPrice(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)

	// Test second call (cache)
	res2, err := s.svc.FindMonthlyTotalPrice(ctx, req)
	s.NoError(err)
	s.Equal(res, res2)
}

func (s *CategoryStatsServiceTestSuite) TestFindYearlyTotalPrice() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyTotalPrice(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindMonthPrice() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindMonthPrice(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindYearPrice() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearPrice(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindMonthlyTotalPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthTotalPriceCategory{
		Year:       now.Year(),
		Month:      int(now.Month()),
		CategoryID: s.categoryID,
	}

	res, err := s.svcById.FindMonthlyTotalPriceById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindYearlyTotalPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearTotalPriceCategory{
		Year:       now.Year(),
		CategoryID: s.categoryID,
	}

	res, err := s.svcById.FindYearlyTotalPriceById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindMonthPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthPriceId{
		Year:       now.Year(),
		CategoryID: s.categoryID,
	}

	res, err := s.svcById.FindMonthPriceById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindYearPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearPriceId{
		Year:       now.Year(),
		CategoryID: s.categoryID,
	}

	res, err := s.svcById.FindYearPriceById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindMonthlyTotalPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthTotalPriceMerchant{
		Year:       now.Year(),
		Month:      int(now.Month()),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthlyTotalPriceByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindYearlyTotalPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearTotalPriceMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearlyTotalPriceByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindMonthPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthPriceMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthPriceByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsServiceTestSuite) TestFindYearPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearPriceMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearPriceByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestCategoryStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsServiceTestSuite))
}
