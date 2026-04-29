package order_test

import (
	"context"
	"testing"
	"time"

	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderStatsServiceTestSuite struct {
	tests.BaseTestSuite
	svc             service.OrderStatsService
	svcByMerchant   service.OrderStatsByMerchantService
	merchantID      int
	userID          int
}

func (s *OrderStatsServiceTestSuite) SetupSuite() {
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
	repos := repository.NewRepositories(&repository.Deps{
		DB: queries,
	})
	cacheStore := s.GetCacheStore()
	mencache := order_cache.NewMencache(cacheStore)

	deps := &service.OrderStatsServiceDeps{
		Cache:                mencache.OrderStatsCache,
		OrderStatsRepository: repos.OrderStats,
		Logger:               s.Log,
		Observability:        s.Obs,
	}
	s.svc = service.NewOrderStatsService(deps)

	depsMerchant := &service.OrderStatsByMerchantServiceDeps{
		Cache:                          mencache.OrderStatsByMerchantCache,
		OrderStatsByMerchantRepository: repos.OrderStatsByMerchant,
		Logger:                         s.Log,
		Observability:                  s.Obs,
	}
	s.svcByMerchant = service.NewOrderStatsByMerchantService(depsMerchant)

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	_, err := s.DBPool().Exec(ctx, "UPDATE orders SET created_at = $1 WHERE order_id = $2", 
		time.Now(), orderID)
	s.Require().NoError(err)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyTotalRevenue() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthTotalRevenue{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyTotalRevenue(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyTotalRevenue() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyTotalRevenue(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyOrder() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindMonthlyOrder(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyOrder() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyOrder(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthTotalRevenueMerchant{
		Year:       now.Year(),
		Month:      int(now.Month()),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearTotalRevenueMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyOrderByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthOrderMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthlyOrderByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyOrderByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearOrderMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearlyOrderByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestOrderStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsServiceTestSuite))
}
