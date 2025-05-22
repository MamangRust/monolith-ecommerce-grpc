package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
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
	ctx                        context.Context
	trace                      trace.Tracer
	transactionQueryRepository repository.TransactionQueryRepository
	mapping                    response_service.TransactionResponseMapper
	logger                     logger.LoggerInterface
	requestCounter             *prometheus.CounterVec
	requestDuration            *prometheus.HistogramVec
}

func NewTransactionQueryService(
	ctx context.Context,
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
		ctx:                        ctx,
		trace:                      otel.Tracer("transaction-query-service"),
		transactionQueryRepository: transactionQueryRepository,
		mapping:                    mapping,
		logger:                     logger,
		requestCounter:             requestCounter,
		requestDuration:            requestDuration,
	}
}

func (s *transactionQueryService) FindAllTransactions(req *requests.FindAllTransaction) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAllTransactions", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllTransactions")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindAllTransactions(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_TRANSACTIONS")

		s.logger.Error("Failed to retrieve transactions from database",
			zap.String("traceID", traceID),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve transactions from database")

		status = "failed_find_all_transactions"

		return nil, nil, transaction_errors.ErrFailedFindAllTransactions
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponse(transactions), totalRecords, nil
}

func (s *transactionQueryService) FindByMerchant(req *requests.FindAllTransactionByMerchant) ([]*response.TransactionResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByMerchant")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByMerchant(req)

	if err != nil {
		s.logger.Error("Failed to retrieve merchant's transactions from database",
			zap.Int("merchant_id", merchant_id),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Error(err))
		return nil, nil, transaction_errors.ErrFailedFindTransactionsByMerchant
	}

	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponse(transactions), totalRecords, nil
}

func (s *transactionQueryService) FindByActive(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActive")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
		span.SetAttributes(attribute.Int("adjusted.page", page))
	}

	if pageSize <= 0 {
		pageSize = 10
		span.SetAttributes(attribute.Int("adjusted.pageSize", pageSize))
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByActive(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTIONS_BY_ACTIVE")

		s.logger.Error("Failed to retrieve active transactions from database",
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active transactions")

		status = "failed_find_by_active"
		return nil, nil, transaction_errors.ErrFailedFindTransactionsByActive
	}

	span.SetAttributes(attribute.Int("totalRecords", *totalRecords))
	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponseDeleteAt(transactions), totalRecords, nil
}

func (s *transactionQueryService) FindByTrashed(req *requests.FindAllTransaction) ([]*response.TransactionResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching transactions",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
		span.SetAttributes(attribute.Int("adjusted.page", page))
	}

	if pageSize <= 0 {
		pageSize = 10
		span.SetAttributes(attribute.Int("adjusted.pageSize", pageSize))
	}

	transactions, totalRecords, err := s.transactionQueryRepository.FindByTrashed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTIONS_BY_TRASHED")

		s.logger.Error("Failed to retrieve trashed transactions",
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve trashed transactions")

		status = "failed_find_by_trashed"
		return nil, nil, transaction_errors.ErrFailedFindTransactionsByTrashed
	}

	span.SetAttributes(attribute.Int("totalRecords", *totalRecords))
	s.logger.Debug("Successfully fetched transactions",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToTransactionsResponseDeleteAt(transactions), totalRecords, nil
}

func (s *transactionQueryService) FindById(transactionID int) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(attribute.Int("transaction.id", transactionID))

	s.logger.Debug("Fetching transaction by ID",
		zap.Int("transactionID", transactionID))

	transaction, err := s.transactionQueryRepository.FindById(transactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTION_BY_ID")

		s.logger.Error("Failed to retrieve transaction details",
			zap.Int("transaction_id", transactionID),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve transaction by ID")

		status = "failed_find_by_id"
		return nil, transaction_errors.ErrFailedFindTransactionById
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionQueryService) FindByOrderId(orderID int) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindByOrderId", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByOrderId")
	defer span.End()

	span.SetAttributes(attribute.Int("order.id", orderID))

	s.logger.Debug("Fetching transaction by Order ID",
		zap.Int("orderID", orderID))

	transaction, err := s.transactionQueryRepository.FindByOrderId(orderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTION_BY_ORDER_ID")

		s.logger.Error("Failed to retrieve transaction by order ID",
			zap.Int("order_id", orderID),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve transaction by order ID")

		status = "failed_find_by_order_id"
		return nil, transaction_errors.ErrFailedFindTransactionByOrderId
	}

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
