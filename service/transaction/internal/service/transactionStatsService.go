package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionStatsService struct {
	ctx                        context.Context
	trace                      trace.Tracer
	transactionStatsRepository repository.TransactionStatsRepository
	mapping                    response_service.TransactionResponseMapper
	logger                     logger.LoggerInterface
	requestCounter             *prometheus.CounterVec
	requestDuration            *prometheus.HistogramVec
}

func NewTransactionStatsService(
	ctx context.Context,
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
		trace:                      otel.Tracer("transaction-stats-service"),
		transactionStatsRepository: transactionStatsRepository,
		mapping:                    mapping,
		logger:                     logger,
		requestCounter:             requestCounter,
		requestDuration:            requestDuration,
	}
}

func (s *transactionStatsService) FindMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountSuccessResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountSuccess")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("find monthly successful transaction amounts",
		zap.Int("year", year),
		zap.Int("month", month))

	res, err := s.transactionStatsRepository.GetMonthlyAmountSuccess(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNT_SUCCESS")

		s.logger.Error("failed to get monthly successful transaction amounts",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly successful transaction amounts")

		status = "failed_find_monthly_amount_success"

		return nil, transaction_errors.ErrFailedFindMonthlyAmountSuccess
	}

	return s.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (s *transactionStatsService) FindYearlyAmountSuccess(year int) ([]*response.TransactionYearlyAmountSuccessResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountSuccess")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("find yearly successful transaction amounts",
		zap.Int("year", year))

	res, err := s.transactionStatsRepository.GetYearlyAmountSuccess(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNT_SUCCESS")

		s.logger.Error("failed to get yearly successful transaction amounts",
			zap.Int("year", year),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly successful transaction amounts")

		status = "failed_find_yearly_amount_success"

		return nil, transaction_errors.ErrFailedFindYearlyAmountSuccess
	}

	return s.mapping.ToTransactionYearlyAmountSuccess(res), nil
}

func (s *transactionStatsService) FindMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*response.TransactionMonthlyAmountFailedResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyAmountFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyAmountFailed")
	defer span.End()

	year := req.Year
	month := req.Month

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("find monthly failed transaction amounts",
		zap.Int("year", year),
		zap.Int("month", month))

	res, err := s.transactionStatsRepository.GetMonthlyAmountFailed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_AMOUNT_FAILED")

		s.logger.Error("failed to get monthly failed transaction amounts",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly failed transaction amounts")

		status = "failed_find_monthly_amount_failed"

		return nil, transaction_errors.ErrFailedFindMonthlyAmountFailed
	}

	return s.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (s *transactionStatsService) FindYearlyAmountFailed(year int) ([]*response.TransactionYearlyAmountFailedResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyAmountFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyAmountFailed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("find yearly failed transaction amounts",
		zap.Int("year", year))

	res, err := s.transactionStatsRepository.GetYearlyAmountFailed(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_AMOUNT_FAILED")

		s.logger.Error("failed to get yearly failed transaction amounts",
			zap.Int("year", year),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly failed transaction amounts")

		status = "failed_find_yearly_amount_failed"

		return nil, transaction_errors.ErrFailedFindYearlyAmountFailed
	}

	return s.mapping.ToTransactionYearlyAmountFailed(res), nil
}

func (s *transactionStatsService) FindMonthlyMethodSuccess(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyMethodSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyMethodSuccess")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
	)

	s.logger.Debug("find monthly successful transaction methods",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodSuccess(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_METHOD_SUCCESS")

		s.logger.Error("failed to get monthly transaction methods",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly transaction methods")

		status = "failed_find_monthly_method_success"

		return nil, transaction_errors.ErrFailedFindMonthlyMethod
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionStatsService) FindYearlyMethodSuccess(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyMethodSuccess", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyMethodSuccess")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("find yearly successful transaction methods",
		zap.Int("year", year))

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodSuccess(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_METHOD_SUCCESS")

		s.logger.Error("failed to get yearly transaction methods",
			zap.Int("year", year),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly transaction methods")

		status = "failed_find_yearly_method_success"

		return nil, transaction_errors.ErrFailedFindYearlyMethod
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionStatsService) FindMonthlyMethodFailed(req *requests.MonthMethodTransaction) ([]*response.TransactionMonthlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyMethodFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyMethodFailed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
	)

	s.logger.Debug("find monthly failed transaction methods",
		zap.Int("year", req.Year),
		zap.Int("month", req.Month))

	res, err := s.transactionStatsRepository.GetMonthlyTransactionMethodFailed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_METHOD_FAILED")

		s.logger.Error("failed to get monthly transaction methods",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly transaction methods")

		status = "failed_find_monthly_method_failed"

		return nil, transaction_errors.ErrFailedFindMonthlyMethod
	}

	return s.mapping.ToTransactionMonthlyMethod(res), nil
}

func (s *transactionStatsService) FindYearlyMethodFailed(year int) ([]*response.TransactionYearlyMethodResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyMethodFailed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyMethodFailed")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	s.logger.Debug("find yearly failed transaction methods",
		zap.Int("year", year))

	res, err := s.transactionStatsRepository.GetYearlyTransactionMethodFailed(year)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_METHOD_FAILED")

		s.logger.Error("failed to get yearly transaction methods",
			zap.Int("year", year),
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace.id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly transaction methods")

		status = "failed_find_yearly_method_failed"

		return nil, transaction_errors.ErrFailedFindYearlyMethod
	}

	return s.mapping.ToTransactionYearlyMethod(res), nil
}

func (s *transactionStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
