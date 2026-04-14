package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
)

type transactionStatsByMerchantService struct {
	observability observability.TraceLoggerObservability
	cache         cache.TransactionStatsByMerchantCache
	repository    repository.TransactionStatsByMerchantRepository
	logger        logger.LoggerInterface
}

type TransactionStatsByMerchantServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.TransactionStatsByMerchantCache
	Repository    repository.TransactionStatsByMerchantRepository
	Logger        logger.LoggerInterface
}

func NewTransactionStatsByMerchantService(deps *TransactionStatsByMerchantServiceDeps) *transactionStatsByMerchantService {
	return &transactionStatsByMerchantService{
		observability: deps.Observability,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
	}
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error) {
	const method = "FindMonthlyAmountSuccessByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyAmountSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionSuccessByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyAmountSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error) {
	const method = "FindYearlyAmountSuccessByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyAmountSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionSuccessByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyAmountSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error) {
	const method = "FindMonthlyAmountFailedByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyAmountFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionFailedByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyAmountFailedByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error) {
	const method = "FindYearlyAmountFailedByMerchant"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyAmountFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionFailedByMerchantRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyAmountFailedByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error) {
	const method = "FindMonthlyMethodByMerchantSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyMethodSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyMethodSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error) {
	const method = "FindYearlyMethodByMerchantSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyMethodSuccessByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsByMerchantSuccessRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyMethodSuccessByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error) {
	const method = "FindMonthlyMethodByMerchantFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyMethodFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsByMerchantFailedRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyMethodFailedByMerchant(ctx, req, res)
	return res, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error) {
	const method = "FindYearlyMethodByMerchantFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("year", req.Year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyMethodFailedByMerchant(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsByMerchantFailedRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyMethodFailedByMerchant(ctx, req, res)
	return res, nil
}
