package order_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type OrderStatsGapiTestSuite struct {
	tests.BaseTestSuite
	client      pb.OrderStatsServiceClient
	merchantID  int
	userID      int
}

func (s *OrderStatsGapiTestSuite) SetupSuite() {
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

	s.client = pb.NewOrderStatsServiceClient(s.Conns["order"])

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

func (s *OrderStatsGapiTestSuite) TestFindMonthlyTotalRevenue() {
	ctx := context.Background()
	now := time.Now()
	req := &pb.FindYearMonthTotalRevenue{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.FindMonthlyTotalRevenue(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsGapiTestSuite) TestFindYearlyTotalRevenue() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pb.FindYearTotalRevenue{
		Year: int32(year),
	}

	res, err := s.client.FindYearlyTotalRevenue(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsGapiTestSuite) TestFindMonthlyRevenue() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pb.FindYearOrder{
		Year: int32(year),
	}

	res, err := s.client.FindMonthlyRevenue(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsGapiTestSuite) TestFindYearlyRevenue() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pb.FindYearOrder{
		Year: int32(year),
	}

	res, err := s.client.FindYearlyRevenue(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsGapiTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pb.FindYearMonthTotalRevenueByMerchant{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.client.FindMonthlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsGapiTestSuite) TestFindYearlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pb.FindYearTotalRevenueByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.client.FindYearlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsGapiTestSuite) TestFindMonthlyRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pb.FindYearOrderByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.client.FindMonthlyRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsGapiTestSuite) TestFindYearlyRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pb.FindYearOrderByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.client.FindYearlyRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestOrderStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsGapiTestSuite))
}
