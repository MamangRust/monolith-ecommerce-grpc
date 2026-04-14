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

type transactionStatsService struct {
	observability observability.TraceLoggerObservability
	cache         cache.TransactionStatsCache
	repository    repository.TransactionStatsRepository
	logger        logger.LoggerInterface
}

type TransactionStatsServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.TransactionStatsCache
	Repository    repository.TransactionStatsRepository
	Logger        logger.LoggerInterface
}

func NewTransactionStatsService(deps *TransactionStatsServiceDeps) *transactionStatsService {
	return &transactionStatsService{
		observability: deps.Observability,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
	}
}

func (s *transactionStatsService) FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, error) {
	const method = "FindMonthlyAmountSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyAmountSuccess(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyAmountSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionSuccessRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyAmountSuccess(ctx, req, res)
	return res, nil
}

func (s *transactionStatsService) FindYearlyAmountSuccess(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, error) {
	const method = "FindYearlyAmountSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyAmountSuccess(ctx, year); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyAmountSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionSuccessRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyAmountSuccess(ctx, year, res)
	return res, nil
}

func (s *transactionStatsService) FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, error) {
	const method = "FindMonthlyAmountFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyAmountFailed(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyAmountFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyAmountTransactionFailedRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyAmountFailed(ctx, req, res)
	return res, nil
}

func (s *transactionStatsService) FindYearlyAmountFailed(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, error) {
	const method = "FindYearlyAmountFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyAmountFailed(ctx, year); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyAmountFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyAmountTransactionFailedRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyAmountFailed(ctx, year, res)
	return res, nil
}

func (s *transactionStatsService) FindMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error) {
	const method = "FindMonthlyMethodSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyMethodSuccess(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionMethodSuccess(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsSuccessRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyMethodSuccess(ctx, req, res)
	return res, nil
}

func (s *transactionStatsService) FindYearlyMethodSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, error) {
	const method = "FindYearlyMethodSuccess"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyMethodSuccess(ctx, year); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionMethodSuccess(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsSuccessRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyMethodSuccess(ctx, year, res)
	return res, nil
}

func (s *transactionStatsService) FindMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, error) {
	const method = "FindMonthlyMethodFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedMonthlyMethodFailed(ctx, req); found {
		return data, nil
	}

	res, err := s.repository.GetMonthlyTransactionMethodFailed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetMonthlyTransactionMethodsFailedRow](s.logger, err, method, span)
	}

	s.cache.SetCachedMonthlyMethodFailed(ctx, req, res)
	return res, nil
}

func (s *transactionStatsService) FindYearlyMethodFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, error) {
	const method = "FindYearlyMethodFailed"
	ctx, span, end, status, _ := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer func() { end(status) }()

	if data, found := s.cache.GetCachedYearlyMethodFailed(ctx, year); found {
		return data, nil
	}

	res, err := s.repository.GetYearlyTransactionMethodFailed(ctx, year)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetYearlyTransactionMethodsFailedRow](s.logger, err, method, span)
	}

	s.cache.SetCachedYearlyMethodFailed(ctx, year, res)
	return res, nil
}
