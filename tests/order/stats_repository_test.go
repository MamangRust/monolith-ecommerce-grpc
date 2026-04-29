package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderStatsRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo          repository.OrderStatsRepository
	repoByMerchant repository.OrderStatsByMerchantRepository
	testYear      int
	testMonth     int
	userID        int
	merchantID    int
	categoryID    int
	productID     int
}

func (s *OrderStatsRepositoryTestSuite) SetupSuite() {
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
	s.repo = repository.NewOrderStatsRepository(queries)
	s.repoByMerchant = repository.NewOrderStatsByMerchantRepository(queries)

	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())

	ctx := context.Background()

	// Seed data
	s.userID = s.SeedUser(ctx)
	s.categoryID = s.SeedCategory(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.productID = s.SeedProduct(ctx, s.merchantID, s.categoryID)

	// Create an order for the test month
	s.SeedOrder(ctx, s.userID, s.merchantID, s.productID)
}

func (s *OrderStatsRepositoryTestSuite) TestGetMonthlyTotalRevenue() {
	ctx := context.Background()
	req := &requests.MonthTotalRevenue{
		Year:  s.testYear,
		Month: s.testMonth,
	}

	res, err := s.repo.GetMonthlyTotalRevenue(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsRepositoryTestSuite) TestGetYearlyTotalRevenue() {
	ctx := context.Background()
	res, err := s.repo.GetYearlyTotalRevenue(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsRepositoryTestSuite) TestGetMonthlyOrder() {
	ctx := context.Background()
	res, err := s.repo.GetMonthlyOrder(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsRepositoryTestSuite) TestGetYearlyOrder() {
	ctx := context.Background()
	res, err := s.repo.GetYearlyOrder(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsRepositoryTestSuite) TestGetMonthlyTotalRevenueByMerchant() {
	ctx := context.Background()
	req := &requests.MonthTotalRevenueMerchant{
		Year:       s.testYear,
		Month:      s.testMonth,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsRepositoryTestSuite) TestGetYearlyTotalRevenueByMerchant() {
	ctx := context.Background()
	req := &requests.YearTotalRevenueMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsRepositoryTestSuite) TestGetMonthlyOrderByMerchant() {
	ctx := context.Background()
	req := &requests.MonthOrderMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthlyOrderByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *OrderStatsRepositoryTestSuite) TestGetYearlyOrderByMerchant() {
	ctx := context.Background()
	req := &requests.YearOrderMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearlyOrderByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestOrderStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsRepositoryTestSuite))
}
