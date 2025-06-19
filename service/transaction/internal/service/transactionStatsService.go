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
	ctx                        context.Context
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
	ctx context.Context,
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
		ctx:                        ctx,
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

func (s *transactionStatsService) FindMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmountSuccess"

	year := req.Year
	month := req.Month

	span, end, status, logSuccess := s.startTraceWithLogging(method, year, &month)

	defer end()

	if data, found := s.mencache.GetCachedMonthAmountSuccessCached(req); found {
		logSuccess("Successfully fetched monthly successful transaction amounts from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyAmountSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyAmountSuccessError(err, method, "FAILED_FIND_MONTHLY_AMOUNT_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyAmountSuccess(res)

	s.mencache.SetCachedMonthAmountSuccessCached(req, so)

	logSuccess("Successfully fetched monthly successful transaction amounts", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyAmountSuccess(year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmountSuccess"

	span, end, status, logSuccess := s.startTraceWithLogging(method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearAmountSuccessCached(year); found {
		logSuccess("Successfully fetched yearly successful transaction amounts from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyAmountSuccess(year)
	if err != nil {
		return s.errorhandler.HandleYearlyAmountSuccessError(err, method, "FAILED_FIND_YEARLY_AMOUNT_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyAmountSuccess(res)

	s.mencache.SetCachedYearAmountSuccessCached(year, so)

	logSuccess("Successfully fetched yearly successful transaction amounts", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) FindMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmountFailed"

	year := req.Year
	month := req.Month

	span, end, status, logSuccess := s.startTraceWithLogging(method, year, &month)

	defer end()

	if data, found := s.mencache.GetCachedMonthAmountFailedCached(req); found {
		logSuccess("Successfully fetched monthly failed transaction amounts from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyAmountFailed(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyAmountFailedError(err, method, "FAILED_FIND_MONTHLY_AMOUNT_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyAmountFailed(res)

	s.mencache.SetCachedMonthAmountFailedCached(req, so)

	logSuccess("Successfully fetched monthly failed transaction amounts", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyAmountFailed(year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmountFailed"

	span, end, status, logSuccess := s.startTraceWithLogging(method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearAmountFailedCached(year); found {
		logSuccess("Successfully fetched yearly failed transaction amounts from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyAmountFailed(year)

	if err != nil {
		return s.errorhandler.HandleYearlyAmountFailedError(err, method, "FAILED_FIND_YEARLY_AMOUNT_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyAmountFailed(res)

	s.mencache.SetCachedYearAmountFailedCached(year, so)

	logSuccess("Successfully fetched yearly failed transaction amounts", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) FindMonthlyMethodSuccess(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyMethodSuccess"

	span, end, status, logSuccess := s.startTraceWithLogging(method, req.Year, &req.Month)

	defer end()

	if data, found := s.mencache.GetCachedMonthMethodSuccessCached(req); found {
		logSuccess("Successfully fetched monthly successful transaction methods from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyMethodSuccessError(err, method, "FAILED_FIND_MONTHLY_METHOD_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyMethod(res)

	s.mencache.SetCachedMonthMethodSuccessCached(req, so)

	logSuccess("Successfully fetched monthly successful transaction methods", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyMethodSuccess(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyMethodSuccess"

	span, end, status, logSuccess := s.startTraceWithLogging(method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearMethodSuccessCached(year); found {
		logSuccess("Successfully fetched yearly successful transaction methods from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodSuccess(year)

	if err != nil {
		return s.errorhandler.HandleYearlyMethodSuccessError(err, method, "FAILED_FIND_YEARLY_METHOD_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyMethod(res)

	s.mencache.SetCachedYearMethodSuccessCached(year, so)

	logSuccess("Successfully fetched yearly successful transaction methods", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) FindMonthlyMethodFailed(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyMethodFailed"

	span, end, status, logSuccess := s.startTraceWithLogging(method, req.Year, &req.Month)

	defer end()

	if data, found := s.mencache.GetCachedMonthMethodFailedCached(req); found {
		logSuccess("Successfully fetched monthly failed transaction methods from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodFailed(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyMethodFailedError(err, method, "FAILED_FIND_MONTHLY_METHOD_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyMethod(res)

	s.mencache.SetCachedMonthMethodFailedCached(req, so)

	logSuccess("Successfully fetched monthly failed transaction methods", zap.Int("year", req.Year), zap.Int("month", req.Month))

	return so, nil
}

func (s *transactionStatsService) FindYearlyMethodFailed(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyMethodFailed"

	span, end, status, logSuccess := s.startTraceWithLogging(method, year, nil)

	defer end()

	if data, found := s.mencache.GetCachedYearMethodFailedCached(year); found {
		logSuccess("Successfully fetched yearly failed transaction methods from cache", zap.Int("year", year))

		return data, nil
	}

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodFailed(year)

	if err != nil {
		return s.errorhandler.HandleYearlyMethodFailedError(err, method, "FAILED_FIND_YEARLY_METHOD_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyMethod(res)

	s.mencache.SetCachedYearMethodFailedCached(year, so)

	logSuccess("Successfully fetched yearly failed transaction methods", zap.Int("year", year))

	return so, nil
}

func (s *transactionStatsService) startTraceWithLogging(method string, year int, month *int) (trace.Span, func(), string, func(string, ...zap.Field)) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

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

	return span, end, status, logSuccess
}

func (s *transactionStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
