package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	tests "github.com/MamangRust/monolith-ecommerce-test"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo          repository.TransactionStatsRepository
	repoByMerchant repository.TransactionStatsByMerchantRepository
	testYear      int
	testMonth     int
	userID        int
	merchantID    int
	categoryID    int
	productID     int
	orderID       int
}

func (s *TransactionStatsRepositoryTestSuite) SetupSuite() {
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
	s.repo = repository.NewTransactionStatsRepository(queries)
	s.repoByMerchant = repository.NewTransactionStatsByMerchantRepository(queries)

	s.testYear = time.Now().Year()
	s.testMonth = int(time.Now().Month())

	ctx := context.Background()

	// Seed data
	s.userID = s.SeedUser(ctx)
	s.categoryID = s.SeedCategory(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.productID = s.SeedProduct(ctx, s.merchantID, s.categoryID)
	s.orderID = s.SeedOrder(ctx, s.userID, s.merchantID, s.productID)

	// Seed successful transaction
	_, err := s.DBPool().Exec(ctx, `
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, s.orderID, s.merchantID, 10000, "credit_card", "success", time.Date(s.testYear, time.Month(s.testMonth), 10, 10, 0, 0, 0, time.UTC))
	s.Require().NoError(err)

	// Seed failed transaction
	_, err = s.DBPool().Exec(ctx, `
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, s.orderID, s.merchantID, 5000, "bank_transfer", "failed", time.Date(s.testYear, time.Month(s.testMonth), 11, 10, 0, 0, 0, time.UTC))
	s.Require().NoError(err)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyAmountSuccess() {
	ctx := context.Background()
	req := &requests.MonthAmountTransaction{
		Year:  s.testYear,
		Month: s.testMonth,
	}

	res, err := s.repo.GetMonthlyAmountSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyAmountSuccess() {
	ctx := context.Background()
	res, err := s.repo.GetYearlyAmountSuccess(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyAmountFailed() {
	ctx := context.Background()
	req := &requests.MonthAmountTransaction{
		Year:  s.testYear,
		Month: s.testMonth,
	}

	res, err := s.repo.GetMonthlyAmountFailed(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyAmountFailed() {
	ctx := context.Background()
	res, err := s.repo.GetYearlyAmountFailed(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyTransactionMethodSuccess() {
	ctx := context.Background()
	req := &requests.MonthMethodTransaction{
		Year:  s.testYear,
		Month: s.testMonth,
	}

	res, err := s.repo.GetMonthlyTransactionMethodSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyTransactionMethodSuccess() {
	ctx := context.Background()
	res, err := s.repo.GetYearlyTransactionMethodSuccess(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyTransactionMethodFailed() {
	ctx := context.Background()
	req := &requests.MonthMethodTransaction{
		Year:  s.testYear,
		Month: s.testMonth,
	}

	res, err := s.repo.GetMonthlyTransactionMethodFailed(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyTransactionMethodFailed() {
	ctx := context.Background()
	res, err := s.repo.GetYearlyTransactionMethodFailed(ctx, s.testYear)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyAmountSuccessByMerchant() {
	ctx := context.Background()
	req := &requests.MonthAmountTransactionMerchant{
		Year:       s.testYear,
		Month:      s.testMonth,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyAmountSuccessByMerchant() {
	ctx := context.Background()
	req := &requests.YearAmountTransactionMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyAmountFailedByMerchant() {
	ctx := context.Background()
	req := &requests.MonthAmountTransactionMerchant{
		Year:       s.testYear,
		Month:      s.testMonth,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthlyAmountFailedByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyAmountFailedByMerchant() {
	ctx := context.Background()
	req := &requests.YearAmountTransactionMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearlyAmountFailedByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyTransactionMethodByMerchantSuccess() {
	ctx := context.Background()
	req := &requests.MonthMethodTransactionMerchant{
		Year:       s.testYear,
		Month:      s.testMonth,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthlyTransactionMethodByMerchantSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyTransactionMethodByMerchantSuccess() {
	ctx := context.Background()
	req := &requests.YearMethodTransactionMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearlyTransactionMethodByMerchantSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetMonthlyTransactionMethodByMerchantFailed() {
	ctx := context.Background()
	req := &requests.MonthMethodTransactionMerchant{
		Year:       s.testYear,
		Month:      s.testMonth,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetMonthlyTransactionMethodByMerchantFailed(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsRepositoryTestSuite) TestGetYearlyTransactionMethodByMerchantFailed() {
	ctx := context.Background()
	req := &requests.YearMethodTransactionMerchant{
		Year:       s.testYear,
		MerchantID: s.merchantID,
	}

	res, err := s.repoByMerchant.GetYearlyTransactionMethodByMerchantFailed(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestTransactionStatsRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsRepositoryTestSuite))
}
