package category_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	tests "github.com/MamangRust/monolith-ecommerce-test"
)

type CategoryStatsGapiTestSuite struct {
	tests.BaseTestSuite
	client           pbcategory.CategoryStatsServiceClient
	clientById       pbcategory.CategoryStatsByIdServiceClient
	clientByMerchant pbcategory.CategoryStatsByMerchantServiceClient
	categoryID       int
	merchantID       int
	userID           int
}

func (s *CategoryStatsGapiTestSuite) SetupSuite() {
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

	s.client = pbcategory.NewCategoryStatsServiceClient(s.Conns["category"])
	s.clientById = pbcategory.NewCategoryStatsByIdServiceClient(s.Conns["category"])
	s.clientByMerchant = pbcategory.NewCategoryStatsByMerchantServiceClient(s.Conns["category"])

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

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPrices() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearMonthTotalPrices{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.FindMonthlyTotalPrices(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPrices() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategory.FindYearTotalPrices{
		Year: int32(year),
	}

	res, err := s.client.FindYearlyTotalPrices(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPrice() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategory.FindYearCategory{
		Year: int32(year),
	}

	res, err := s.client.FindMonthPrice(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPrice() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategory.FindYearCategory{
		Year: int32(year),
	}

	res, err := s.client.FindYearPrice(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPricesById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearMonthTotalPriceById{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindMonthlyTotalPricesById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPricesById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearTotalPriceById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindYearlyTotalPricesById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearCategoryById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindMonthPriceById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearCategoryById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientById.FindYearPriceById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPricesByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearMonthTotalPriceByMerchant{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindMonthlyTotalPricesByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPricesByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearTotalPriceByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindYearlyTotalPricesByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearCategoryByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindMonthPriceByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategory.FindYearCategoryByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.FindYearPriceByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestCategoryStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsGapiTestSuite))
}
