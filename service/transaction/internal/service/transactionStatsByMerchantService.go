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

type transactionStatsByMerchantService struct {
	ctx                                  context.Context
	trace                                trace.Tracer
	transactionStatsByMerchantRepository repository.TransactonStatsByMerchantRepository
	mapping                              response_service.TransactionResponseMapper
	logger                               logger.LoggerInterface
	requestCounter                       *prometheus.CounterVec
	requestDuration                      *prometheus.HistogramVec
}

func NewTransactionStatsByMerchantService(
	ctx context.Context,
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
		trace:                                otel.Tracer("transaction-stats-by-merchant-service"),
		transactionStatsByMerchantRepository: transactionStatsByMerchantRepository,
		mapping:                              mapping,
		logger:                               logger,
		requestCounter:                       requestCounter,
		requestDuration:                      requestDuration,
	}
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountSuccessByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountSuccessByMerchant")
	defer span.End()

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find monthly successful transactions by merchant",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("merchantID", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountSuccessByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNT_SUCCESS_BY_MERCHANT")

		s.logger.Error("failed to get monthly successful transactions by merchant",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Int("merchantID", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly successful transactions by merchant")

		status = "failed_find_monthly_amount_success_by_merchant"

		return nil, transaction_errors.ErrFailedFindMonthlyAmountSuccessByMerchant
	}

	return s.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountSuccessByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountSuccessByMerchant")
	defer span.End()

	year := req.Year
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find yearly successful transactions by merchant",
		zap.Int("year", year),
		zap.Int("merchantID", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountSuccessByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNT_SUCCESS_BY_MERCHANT")

		s.logger.Error("failed to get yearly successful transactions by merchant",
			zap.Int("year", year),
			zap.Int("merchantID", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly successful transactions by merchant")

		status = "failed_find_yearly_amount_success_by_merchant"

		return nil, transaction_errors.ErrFailedFindYearlyAmountSuccessByMerchant
	}

	return s.mapping.ToTransactionYearlyAmountSuccess(res), nil
}

func (s *transactionStatsByMerchantService) FindMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountFailedByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountFailedByMerchant")
	defer span.End()

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find monthly failed transactions by merchant",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("merchantID", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountFailedByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNT_FAILED_BY_MERCHANT")

		s.logger.Error("failed to get monthly failed transactions by merchant",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Int("merchantID", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly failed transactions by merchant")

		status = "failed_find_monthly_amount_failed_by_merchant"

		return nil, transaction_errors.ErrFailedFindMonthlyAmountFailedByMerchant
	}

	return s.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (s *transactionStatsByMerchantService) FindYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountFailedByMerchant", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountFailedByMerchant")
	defer span.End()

	year := req.Year
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find yearly failed transactions by merchant",
		zap.Int("year", year),
		zap.Int("merchantID", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountFailedByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNT_FAILED_BY_MERCHANT")

		s.logger.Error("failed to get yearly failed transactions by merchant",
			zap.Int("year", year),
			zap.Int("merchantID", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly failed transactions by merchant")

		status = "failed_find_yearly_amount_failed_by_merchant"

		return nil, transaction_errors.ErrFailedFindYearlyAmountFailedByMerchant
	}

	return s.mapping.ToTransactionYearlyAmountFailed(res), nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyMethodByMerchantSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyMethodByMerchantSuccess")
	defer span.End()

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find monthly successful transaction methods by merchant",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("merchant_id", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantSuccess(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_METHOD_BY_MERCHANT_SUCCESS")

		s.logger.Error("failed to get monthly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Int("merchant_id", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly transaction methods by merchant")

		status = "failed_find_monthly_method_by_merchant_success"

		return nil, transaction_errors.ErrFailedFindMonthlyMethodByMerchant
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyMethodByMerchantSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyMethodByMerchantSuccess")
	defer span.End()

	year := req.Year
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find yearly successful transaction methods by merchant",
		zap.Int("year", year),
		zap.Int("merchant_id", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantSuccess(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_METHOD_BY_MERCHANT_SUCCESS")

		s.logger.Error("failed to get yearly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly transaction methods by merchant")

		status = "failed_find_yearly_method_by_merchant_success"

		return nil, transaction_errors.ErrFailedFindYearlyMethodByMerchant
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionStatsByMerchantService) FindMonthlyMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyMethodByMerchantFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyMethodByMerchantFailed")
	defer span.End()

	year := req.Year
	month := req.Month
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find monthly failed transaction methods by merchant",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("merchant_id", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantFailed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_METHOD_BY_MERCHANT_FAILED")

		s.logger.Error("failed to get monthly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Int("merchant_id", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly transaction methods by merchant")

		status = "failed_find_monthly_method_by_merchant_failed"

		return nil, transaction_errors.ErrFailedFindMonthlyMethodByMerchant
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionStatsByMerchantService) FindYearlyMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyMethodByMerchantFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyMethodByMerchantFailed")
	defer span.End()

	year := req.Year
	merchantId := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant.id", merchantId),
	)

	s.logger.Debug("find yearly failed transaction methods by merchant",
		zap.Int("year", year),
		zap.Int("merchant_id", merchantId))

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantFailed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_METHOD_BY_MERCHANT_FAILED")

		s.logger.Error("failed to get yearly transaction methods by merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchantId),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly transaction methods by merchant")

		status = "failed_find_yearly_method_by_merchant_failed"

		return nil, transaction_errors.ErrFailedFindYearlyMethodByMerchant
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}
func (s *transactionStatsByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
