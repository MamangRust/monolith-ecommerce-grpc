package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	pbtransaction "github.com/MamangRust/monolith-ecommerce-pb/transaction"
	tests "github.com/MamangRust/monolith-ecommerce-test"
)

type TransactionStatsGapiTestSuite struct {
	tests.BaseTestSuite
	client           pbtransaction.TransactionStatsServiceClient
	clientByMerchant pbtransaction.TransactionStatsByMerchantServiceClient
	merchantID       int
	userID           int
}

func (s *TransactionStatsGapiTestSuite) SetupSuite() {
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

	s.client = pbtransaction.NewTransactionStatsServiceClient(s.Conns["transaction"])
	s.clientByMerchant = pbtransaction.NewTransactionStatsByMerchantServiceClient(s.Conns["transaction"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Create a successful transaction for this order
	_, err := s.DBPool().Exec(ctx, `
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		orderID, s.merchantID, 10000, "credit_card", "success", time.Now())
	s.Require().NoError(err)

	// Create a failed transaction for this order
	_, err = s.DBPool().Exec(ctx, `
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		orderID, s.merchantID, 10000, "credit_card", "failed", time.Now())
	s.Require().NoError(err)
}

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyAmountSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransaction.MonthAmountTransactionRequest{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.GetMonthlyAmountSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyAmountSuccess() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransaction.YearAmountTransactionRequest{
		Year: int32(year),
	}

	res, err := s.client.GetYearlyAmountSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyAmountFailed() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransaction.MonthAmountTransactionRequest{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.GetMonthlyAmountFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyAmountFailed() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransaction.YearAmountTransactionRequest{
		Year: int32(year),
	}

	res, err := s.client.GetYearlyAmountFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyTransactionMethodSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransaction.MonthMethodTransactionRequest{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.client.GetMonthlyTransactionMethodSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyTransactionMethodSuccess() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransaction.YearMethodTransactionRequest{
		Year: int32(year),
	}

	res, err := s.client.GetYearlyTransactionMethodSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetMonthlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransaction.MonthAmountTransactionMerchantRequest{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.GetMonthlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *TransactionStatsGapiTestSuite) TestGetYearlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransaction.YearAmountTransactionMerchantRequest{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientByMerchant.GetYearlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestTransactionStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsGapiTestSuite))
}
