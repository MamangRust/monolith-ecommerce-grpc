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
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatsService struct {
	errorhandler               errorhandler.TransactionStatsError
	mencache                   mencache.TransactionStatsCache
	trace                      trace.Tracer
	transactionStatsRepository repository.TransactionStatsRepository
	mapping                    response_service.TransactionResponseMapper
	logger                     logger.LoggerInterface
	requestCounter             *prometheus.CounterVec
	requestDuration            *prometheus.HistogramVec
}

func NewTransactionStatsService(
	errorhandler errorhandler.TransactionStatsError,
	mencache mencache.TransactionStatsCache,
	transactionStatsRepository repository.TransactionStatsRepository,
	mapping response_service.TransactionResponseMapper,
	logger logger.LoggerInterface,
) *transactionStatsService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_service_request_total",
			Help: "Total number of requests to the TransactionStatsService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_service_request_duration",
			Help:    "Histogram of request durations for the TransactionStatsService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionStatsService{
		errorhandler:               errorhandler,
		mencache:                   mencache,
		trace:                      otel.Tracer("transaction-stats-service"),
		transactionStatsRepository: transactionStatsRepository,
		mapping:                    mapping,
		logger:                     logger,
		requestCounter:             requestCounter,
		requestDuration:            requestDuration,
	}
}

func (s *transactionStatsService) FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmountSuccess"

	year := req.Year
	month := req.Month

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, year, &month)

	defer end()

	if data, found := s.mencache.GetCachedMonthAmountSuccessCached(ctx, req); found {
		logSuccess("Successfully fetched monthly successful transaction amounts from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyAmountSuccess(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthlyAmountSuccessError(err, method, "FAILED_FIND_MONTHLY_AMOUNT_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyAmountSuccess(res)

	s.mencache.SetCachedMonthAmountSuccessCached(ctx, req, so)

	logSuccess("Successfully fetched monthly successful transaction amounts", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyAmountSuccess(ctx context.Context, year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmountSuccess"

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearAmountSuccessCached(ctx, year); found {
		logSuccess("Successfully fetched yearly successful transaction amounts from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyAmountSuccess(ctx, year)
	if err != nil {
		return s.errorhandler.HandleYearlyAmountSuccessError(err, method, "FAILED_FIND_YEARLY_AMOUNT_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyAmountSuccess(res)

	s.mencache.SetCachedYearAmountSuccessCached(ctx, year, so)

	logSuccess("Successfully fetched yearly successful transaction amounts", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmountFailed"

	year := req.Year
	month := req.Month

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, year, &month)

	defer end()

	if data, found := s.mencache.GetCachedMonthAmountFailedCached(ctx, req); found {
		logSuccess("Successfully fetched monthly failed transaction amounts from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyAmountFailed(ctx, req)
	if err != nil {
		return s.errorhandler.HandleMonthlyAmountFailedError(err, method, "FAILED_FIND_MONTHLY_AMOUNT_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyAmountFailed(res)

	s.mencache.SetCachedMonthAmountFailedCached(ctx, req, so)

	logSuccess("Successfully fetched monthly failed transaction amounts", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyAmountFailed(ctx context.Context, year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmountFailed"

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearAmountFailedCached(ctx, year); found {
		logSuccess("Successfully fetched yearly failed transaction amounts from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyAmountFailed(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearlyAmountFailedError(err, method, "FAILED_FIND_YEARLY_AMOUNT_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyAmountFailed(res)

	s.mencache.SetCachedYearAmountFailedCached(ctx, year, so)

	logSuccess("Successfully fetched yearly failed transaction amounts", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) FindMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyMethodSuccess"

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, req.Year, &req.Month)

	defer end()

	if data, found := s.mencache.GetCachedMonthMethodSuccessCached(ctx, req); found {
		logSuccess("Successfully fetched monthly successful transaction methods from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodSuccess(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthlyMethodSuccessError(err, method, "FAILED_FIND_MONTHLY_METHOD_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyMethod(res)

	s.mencache.SetCachedMonthMethodSuccessCached(ctx, req, so)

	logSuccess("Successfully fetched monthly successful transaction methods", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyMethodSuccess(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyMethodSuccess"

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearMethodSuccessCached(ctx, year); found {
		logSuccess("Successfully fetched yearly successful transaction methods from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodSuccess(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearlyMethodSuccessError(err, method, "FAILED_FIND_YEARLY_METHOD_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyMethod(res)

	s.mencache.SetCachedYearMethodSuccessCached(ctx, year, so)

	logSuccess("Successfully fetched yearly successful transaction methods", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) FindMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyMethodFailed"

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, req.Year, &req.Month)

	defer end()

	if data, found := s.mencache.GetCachedMonthMethodFailedCached(ctx, req); found {
		logSuccess("Successfully fetched monthly failed transaction methods from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodFailed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthlyMethodFailedError(err, method, "FAILED_FIND_MONTHLY_METHOD_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyMethod(res)

	s.mencache.SetCachedMonthMethodFailedCached(ctx, req, so)

	logSuccess("Successfully fetched monthly failed transaction methods", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyMethodFailed(ctx context.Context, year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyMethodFailed"

	ctx, span, end, status, logSuccess := s.startTraceWithLogging(ctx, method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearMethodFailedCached(ctx, year); found {
		logSuccess("Successfully fetched yearly failed transaction methods from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodFailed(ctx, year)

	if err != nil {
		return s.errorhandler.HandleYearlyMethodFailedError(err, method, "FAILED_FIND_YEARLY_METHOD_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyMethod(res)

	s.mencache.SetCachedYearMethodFailedCached(ctx, year, so)

	logSuccess("Successfully fetched yearly failed transaction methods", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) startTraceWithLogging(ctx context.Context, method string, year int, month *int) (context.Context, trace.Span, func(), string, func(string, ...zap.Field)) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.Int("year", year),
	}

	logFields := []zap.Field{
		zap.String("method", method),
		zap.Int("year", year),
	}

	if month != nil {
		attrs = append(attrs, attribute.Int("month", *month))
		logFields = append(logFields, zap.Int("month", *month))
	}

	span.SetAttributes(attrs...)
	msg := "Start " + method
	s.logger.Debug(msg, logFields...)
	span.AddEvent(msg)

	end := func() {
		span.SetStatus(codes.Ok, status)
		s.recordMetrics(method, status, start)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *transactionStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
