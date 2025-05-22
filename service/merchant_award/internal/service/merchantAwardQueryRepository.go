package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantAwardQueryService struct {
	ctx                          context.Context
	trace                        trace.Tracer
	merchantAwardQueryRepositroy repository.MerchantAwardQueryRepository
	logger                       logger.LoggerInterface
	mapping                      response_service.MerchantAwardResponseMapper
	requestCounter               *prometheus.CounterVec
	requestDuration              *prometheus.HistogramVec
}

func NewMerchantAwardQueryService(ctx context.Context, merchantAwardQueryRepositroy repository.MerchantAwardQueryRepository, logger logger.LoggerInterface, mapping response_service.MerchantAwardResponseMapper) *merchantAwardQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_award_query_service_request_total",
			Help: "Total number of requests to the MerchantAwardQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_award_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantAwardQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantAwardQueryService{
		ctx:                          ctx,
		trace:                        otel.Tracer("merchant-award-query-service"),
		merchantAwardQueryRepositroy: merchantAwardQueryRepositroy,
		logger:                       logger,
		mapping:                      mapping,
		requestCounter:               requestCounter,
		requestDuration:              requestDuration,
	}
}

func (s *merchantAwardQueryService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAll")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all merchants",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantAwardQueryRepositroy.FindAllMerchants(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_MERCHANT_AWARDS")

		s.logger.Error("Failed to retrieve merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve merchants")

		status = "failed_find_all_merchant_awards"

		return nil, nil, merchantaward_errors.ErrFailedFindAllMerchantAwards
	}

	s.logger.Debug("Successfully fetched merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToMerchantsAwardResponse(merchants), totalRecords, nil
}

func (s *merchantAwardQueryService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all merchants active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantAwardQueryRepositroy.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_MERCHANT_AWARDS")

		s.logger.Error("Failed to retrieve active merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active merchants")

		status = "failed_find_active_merchant_awards"

		return nil, nil, merchantaward_errors.ErrFailedFindActiveMerchantAwards
	}

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsAwardResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantAwardQueryService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all merchants trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantAwardQueryRepositroy.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_MERCHANT_AWARDS")

		s.logger.Error("Failed to retrieve trashed merchants",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve trashed merchants")

		status = "failed_find_trashed_merchant_awards"

		return nil, nil, merchantaward_errors.ErrFailedFindTrashedMerchantAwards
	}

	s.logger.Debug("Successfully fetched trashed merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsAwardResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantAwardQueryService) FindById(merchantID int) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchantID", merchantID),
	)

	s.logger.Debug("Fetching merchant by ID", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantAwardQueryRepositroy.FindById(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_AWARD_BY_ID")

		s.logger.Error("Failed to retrieve merchant by ID",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve merchant by ID")

		status = "failed_find_merchant_award_by_id"

		return nil, merchantaward_errors.ErrFailedFindMerchantAwardById
	}

	return s.mapping.ToMerchantAwardResponse(merchant), nil
}

func (s *merchantAwardQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
