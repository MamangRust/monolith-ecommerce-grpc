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

type transactionStatsByMerchantService struct {
	ctx                                  context.Context
	errorhandler                         errorhandler.TransactionStatsByMerchantError
	mencache                             mencache.TransactionStatsByMerchantCache
	trace                                trace.Tracer
	transactionStatsByMerchantRepository repository.TransactonStatsByMerchantRepository
	mapping                              response_service.TransactionResponseMapper
	logger                               logger.LoggerInterface
	requestCounter                       *prometheus.CounterVec
	requestDuration                      *prometheus.HistogramVec
}

func NewTransactionStatsByMerchantService(
	ctx context.Context,
	errorhandler errorhandler.TransactionStatsByMerchantError,
	mencache mencache.TransactionStatsByMerchantCache,
	transactionStatsByMerchantRepository repository.TransactonStatsByMerchantRepository,
	mapping response_service.TransactionResponseMapper,
	logger logger.LoggerInterface,
) *transactionStatsByMerchantService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_stats_by_merchant_service_request_total",
			Help: "Total number of requests to the TransactionStatsByMerchantService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_stats_by_merchant_service_request_duration",
			Help:    "Histogram of request durations for the TransactionStatsByMerchantService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionStatsByMerchantService{
		ctx:                                  ctx,
		errorhandler:                         errorhandler,
		mencache:                             mencache,
		trace:                                otel.Tracer("transaction-stats-by-merchant-service"),
		transactionStatsByMerchantRepository: transactionStatsByMerchantRepository,
		mapping:                              mapping,
		logger:                               logger,
		requestCounter:                       requestCounter,
		requestDuration:                      requestDuration,
	}
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmountSuccessByMerchant"

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountSuccessCached(req); found {
		logSuccess("Successfully fetched monthly successful transactions by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountSuccessByMerchant(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyAmountSuccessByMerchantError(err, method, "FAILED_FIND_MONTHLY_AMOUNT_SUCCESS_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyAmountSuccess(res)

	s.mencache.SetCachedMonthAmountSuccessCached(req, so)

	logSuccess("Successfully fetched monthly successful transactions by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmountSuccessByMerchant"

	year := req.Year
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountSuccessCached(req); found {
		logSuccess("Successfully fetched yearly successful transactions by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountSuccessByMerchant(req)

	if err != nil {
		return s.errorhandler.HandleYearlyAmountSuccessByMerchantError(err, method, "FAILED_FIND_YEARLY_AMOUNT_SUCCESS_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyAmountSuccess(res)

	s.mencache.SetCachedYearAmountSuccessCached(req, so)

	logSuccess("Successfully fetched yearly successful transactions by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	const method = "FindMonthlyAmountFailedByMerchant"

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountFailedCached(req); found {
		logSuccess("Successfully fetched monthly failed transactions by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountFailedByMerchant(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyAmountFailedByMerchantError(err, method, "FAILED_FIND_MONTHLY_AMOUNT_FAILED_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyAmountFailed(res)

	s.mencache.SetCachedMonthAmountFailedCached(req, so)

	logSuccess("Successfully fetched monthly failed transactions by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	const method = "FindYearlyAmountFailedByMerchant"

	year := req.Year
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountFailedCached(req); found {
		logSuccess("Successfully fetched yearly failed transactions by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountFailedByMerchant(req)

	if err != nil {
		return s.errorhandler.HandleYearlyAmountFailedByMerchantError(err, method, "FAILED_FIND_YEARLY_AMOUNT_FAILED_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyAmountFailed(res)

	s.mencache.SetCachedYearAmountFailedCached(req, so)

	logSuccess("Successfully fetched yearly failed transactions by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyMethodByMerchantSuccess"

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodSuccessCached(req); found {
		logSuccess("Successfully found monthly successful transaction methods by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantSuccess(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyMethodSuccessByMerchantError(err, method, "FAILED_FIND_MONTHLY_METHOD_BY_MERCHANT_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyMethod(res)

	s.mencache.SetCachedMonthMethodSuccessCached(req, so)

	logSuccess("Successfully found monthly successful transaction methods by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyMethodByMerchantSuccess"

	year := req.Year
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodSuccessCached(req); found {
		logSuccess("Successfully found yearly successful transaction methods by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantSuccess(req)

	if err != nil {
		return s.errorhandler.HandleYearlyMethodSuccessByMerchantError(err, method, "FAILED_FIND_YEARLY_METHOD_BY_MERCHANT_SUCCESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyMethod(res)

	s.mencache.SetCachedYearMethodSuccessCached(req, so)

	logSuccess("Successfully found yearly successful transaction methods by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	const method = "FindMonthlyMethodByMerchantFailed"

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodFailedCached(req); found {
		logSuccess("Successfully found monthly failed transaction methods by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantFailed(req)
	if err != nil {
		return s.errorhandler.HandleMonthlyMethodFailedByMerchantError(err, method, "FAILED_FIND_MONTHLY_METHOD_BY_MERCHANT_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionMonthlyMethod(res)

	s.mencache.SetCachedMonthMethodFailedCached(req, so)

	logSuccess("Successfully found monthly failed transaction methods by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	const method = "FindYearlyMethodByMerchantFailed"

	year := req.Year
	merchantId := req.MerchantID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("merchant.id", merchantId))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodFailedCached(req); found {
		logSuccess("Successfully found yearly failed transaction methods by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantId))

		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantFailed(req)

	if err != nil {
		return s.errorhandler.HandleYearlyMethodFailedByMerchantError(err, method, "FAILED_FIND_YEARLY_METHOD_BY_MERCHANT_FAILED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionYearlyMethod(res)

	s.mencache.SetCachedYearMethodFailedCached(req, so)

	logSuccess("Successfully found yearly failed transaction methods by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantId))

	return so, nil
}

func (s *transactionStatsByMerchantService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

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

	return span, end, status, logSuccess
}

func (s *transactionStatsByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
