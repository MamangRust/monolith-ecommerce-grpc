package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionQueryService struct {
	mencache                   mencache.TransactionQueryCache
	errorhandler               errorhandler.TransactionQueryError
	trace                      trace.Tracer
	transactionQueryRepository repository.TransactionQueryRepository
	mapping                    response_service.TransactionResponseMapper
	logger                     logger.LoggerInterface
	requestCounter             *prometheus.CounterVec
	requestDuration            *prometheus.HistogramVec
}

func NewTransactionQueryService(
	mencache mencache.TransactionQueryCache,
	errorhandler errorhandler.TransactionQueryError,
	transactionQueryRepository repository.TransactionQueryRepository,
	mapping response_service.TransactionResponseMapper,
	logger logger.LoggerInterface,
) *transactionQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_query_service_request_total",
			Help: "Total number of requests to the TransactionQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_query_service_request_duration",
			Help:    "Histogram of request durations for the TransactionQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionQueryService{
		mencache:                   mencache,
		errorhandler:               errorhandler,
		trace:                      otel.Tracer("transaction-query-service"),
		transactionQueryRepository: transactionQueryRepository,
		mapping:                    mapping,
		logger:                     logger,
		requestCounter:             requestCounter,
		requestDuration:            requestDuration,
	}
}

func (s *transactionQueryService) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	const method = "FindAllTransactions"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionsCache(ctx, req); found {
		logSuccess("Successfully fetched transactions from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactions(ctx, req)
	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_TO_FIND_ALL_TRANSACTIONS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionsCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched all transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	const method = "FindByMerchant"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionByMerchant(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_TO_FIND_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponse(transactions)

	s.mencache.SetCachedTransactionByMerchant(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched all transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionActiveCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRANSACTIONS_BY_ACTIVE", span, &status, transaction_errors.ErrFailedFindTransactionsByActive, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedTransactionTrashedCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_TRANSACTIONS_BY_TRASHED", span, &status, transaction_errors.ErrFailedFindTransactionsByTrashed, zap.Error(err))
	}

	so := s.mapping.ToTransactionsResponseDeleteAt(transactions)

	s.mencache.SetCachedTransactionTrashedCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed transactions", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *transactionQueryService) FindById(ctx context.Context, transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransactionCache(ctx, transactionID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("transaction.id", transactionID))

		return data, nil
	}

	transaction, err := s.transactionQueryRepository.FindById(ctx, transactionID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSACTION_BY_ID", span, &status, transaction_errors.ErrFailedFindTransactionById, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.mencache.SetCachedTransactionCache(ctx, so)

	logSuccess("Successfully fetched transaction", zap.Int("transaction.id", transactionID))

	return so, nil
}

func (s *transactionQueryService) FindByOrderId(ctx context.Context, orderID int) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "FindByOrderId"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedTransactionByOrderId(ctx, orderID); found {
		logSuccess("Successfully fetched transaction from cache", zap.Int("order.id", orderID))

		return data, nil
	}

	transaction, err := s.transactionQueryRepository.FindByOrderId(ctx, orderID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSACTION_BY_ORDER_ID", span, &status, transaction_errors.ErrFailedFindTransactionByOrderId, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.mencache.SetCachedTransactionByOrderId(ctx, orderID, so)

	logSuccess("Successfully fetched transaction", zap.Int("order.id", orderID))

	return so, nil
}

func (s *transactionQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
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

func (s *transactionQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
