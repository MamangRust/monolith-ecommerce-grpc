package category_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type CategoryStatsRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo          repository.CategoryStatsRepository
	repoById      repository.CategoryStatsByIdRepository
	repoByMerchant repository.CategoryStatsByMerchantRepository
	testYear      int
	testMonth     int
	userID        int
	merchantID    int
	categoryID    int
	productID     int
}

func (s *CategoryStatsRepositoryTestSuite) SetupSuite() {
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
	s.repo = repository.NewCategoryStatsRepository(queries)
	s.repoById = repository.NewCategoryStatsByIdRepository(queries)
	s.repoByMerchant = repository.NewCategoryStatsByMerchantRepository(queries)

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

func (s *CategoryStatsRepositoryTestSuite) TestGetMonthlyTotalPrice() {
	ctx := context.Background()
	req := &requests.MonthTotalPrice{
		Year:  s.testYear,
		Month: s.testMonth,
	}

	res, err := s.repo.GetMonthlyTotalPrice(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetYearlyTotalPrices() {
	ctx := context.Background()
	res, err := s.repo.GetYearlyTotalPrices(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetMonthPrice() {
	ctx := context.Background()
	res, err := s.repo.GetMonthPrice(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetYearPrice() {
	ctx := context.Background()
	res, err := s.repo.GetYearPrice(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetMonthlyTotalPriceById() {
	ctx := context.Background()
	req := &requests.MonthTotalPriceCategory{
		Year:       s.testYear,
		Month:      s.testMonth,
		CategoryID: s.categoryID,
	}

	res, err := s.repoById.GetMonthlyTotalPriceById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetYearlyTotalPricesById() {
	ctx := context.Background()
	req := &requests.YearTotalPriceCategory{
		Year:       s.testYear,
		CategoryID: s.categoryID,
	}

	res, err := s.repoById.GetYearlyTotalPricesById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetMonthPriceById() {
	ctx := context.Background()
	req := &requests.MonthPriceId{
		Year:       s.testYear,
		CategoryID: s.categoryID,
	}

	res, err := s.repoById.GetMonthPriceById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetYearPriceById() {
	ctx := context.Background()
	req := &requests.YearPriceId{
		Year:       s.testYear,
		CategoryID: s.categoryID,
	}

	res, err := s.repoById.GetYearPriceById(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetMonthlyTotalPriceByMerchant() {
	ctx := context.Background()
	req := &requests.MonthTotalPriceMerchant{
		Year:       s.testYear,
		Month:      s.testMonth,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthlyTotalPriceByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetYearlyTotalPricesByMerchant() {
	ctx := context.Background()
	req := &requests.YearTotalPriceMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearlyTotalPricesByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetMonthPriceByMerchant() {
	ctx := context.Background()
	req := &requests.MonthPriceMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthPriceByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *CategoryStatsRepositoryTestSuite) TestGetYearPriceByMerchant() {
	ctx := context.Background()
	req := &requests.YearPriceMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearPriceByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestCategoryStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsRepositoryTestSuite))
}
