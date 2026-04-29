package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionQueryService struct {
	observability observability.TraceLoggerObservability
	cache         cache.TransactionQueryCache
	repository    repository.TransactionQueryRepository
	logger        logger.LoggerInterface
}

type TransactionQueryServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.TransactionQueryCache
	Repository    repository.TransactionQueryRepository
	Logger        logger.LoggerInterface
}

func NewTransactionQueryService(deps *TransactionQueryServiceDeps) TransactionQueryService {
	return &transactionQueryService{
		observability: deps.Observability,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
	}
}

func (s *transactionQueryService) FindAll(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, *int, error) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionsCache(ctx, req); found {
		logSuccess("Successfully fetched transactions from cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.repository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	}

	s.cache.SetCachedTransactionsCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched all transactions", zap.Int("page", page), zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*db.GetTransactionByMerchantRow, *int, error) {
	const method = "FindByMerchant"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.Int("merchant_id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionByMerchant(ctx, req); found {
		logSuccess("Successfully fetched transactions from cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.repository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant_id", req.MerchantID),
		)
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	}

	s.cache.SetCachedTransactionByMerchant(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched all transactions by merchant", zap.Int("page", page), zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionQueryService) FindActive(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, *int, error) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionActiveCache(ctx, req); found {
		logSuccess("Successfully fetched transactions from cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.repository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsActiveRow](
			s.logger,
			err,
			method,
			span,
		)
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	}

	s.cache.SetCachedTransactionActiveCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched active transactions", zap.Int("page", page), zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionQueryService) FindTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, *int, error) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedTransactionTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched transactions from cache", zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	transactions, err := s.repository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTransactionsTrashedRow](
			s.logger,
			err,
			method,
			span,
		)
	}

	var totalCount int
	if len(transactions) > 0 {
		totalCount = int(transactions[0].TotalCount)
	}

	s.cache.SetCachedTransactionTrashedCache(ctx, req, transactions, &totalCount)

	logSuccess("Successfully fetched trashed transactions", zap.Int("page", page), zap.Int("pageSize", pageSize))

	return transactions, &totalCount, nil
}

func (s *transactionQueryService) FindByID(ctx context.Context, transactionID int) (*db.GetTransactionByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransactionCache(ctx, transactionID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("transaction_id", transactionID))
		return data, nil
	}

	res, err := s.repository.FindByID(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetTransactionByIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("transaction_id", transactionID),
		)
	}

	s.cache.SetCachedTransactionCache(ctx, res)

	logSuccess("Successfully fetched transaction", zap.Int("transaction_id", transactionID))

	return res, nil
}

func (s *transactionQueryService) FindByOrderID(ctx context.Context, orderID int) (*db.GetTransactionByOrderIDRow, error) {
	const method = "FindByOrderID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedTransactionByOrderId(ctx, orderID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("order_id", orderID))
		return data, nil
	}

	res, err := s.repository.FindByOrderID(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetTransactionByOrderIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("order_id", orderID),
		)
	}

	s.cache.SetCachedTransactionByOrderId(ctx, orderID, res)

	logSuccess("Successfully fetched transaction by order id", zap.Int("order_id", orderID))

	return res, nil
}

func (s *transactionQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
