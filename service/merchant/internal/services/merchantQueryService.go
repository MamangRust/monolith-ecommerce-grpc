package services

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type merchantQueryService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	merchantQueryRepository repository.MerchantQueryRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.MerchantResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewMerchantQueryService(ctx context.Context, merchantQueryRepository repository.MerchantQueryRepository, logger logger.LoggerInterface, mapping response_service.MerchantResponseMapper) *merchantQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_query_service_requests_total",
			Help: "Total number of requests to the MerchantQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantQueryService{
		ctx:                     ctx,
		trace:                   otel.Tracer("merchant-query-service"),
		merchantQueryRepository: merchantQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *merchantQueryService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, startTime)
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

	s.logger.Debug("Fetching all merchant records",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindAllMerchants(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_MERCHANTS")

		s.logger.Error("Failed to retrieve all merchants",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve all merchants")
		status = "failed_to_find_all_merchants"

		return nil, nil, merchant_errors.ErrFailedFindAllMerchants
	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.logger.Debug("Successfully all merchant records",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindById(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Finding merchant by ID", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantQueryRepository.FindById(merchant_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_BY_ID")

		s.logger.Error("Failed to find merchant by ID",
			zap.Error(err),
			zap.Int("merchant_id", merchant_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant by ID")
		status = "failed_to_find_merchant_by_id"

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	so := s.mapping.ToMerchantResponse(res)

	return so, nil
}

func (s *merchantQueryService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, startTime)
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

	s.logger.Debug("Fetching all merchant active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_MERCHANTS")

		s.logger.Error("Failed to retrieve active merchant",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active merchant")
		status = "failed_to_find_active_merchants"

		return nil, nil, merchant_errors.ErrFailedFindActiveMerchants
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched active merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantQueryService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()

	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, startTime)
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

	s.logger.Debug("Fetching fetched trashed merchants",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_MERCHANTS")

		s.logger.Error("Failed to retrieve trashed merchant",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve trashed merchant")
		status = "failed_to_find_trashed_merchants"

		return nil, nil, merchant_errors.ErrFailedFindTrashedMerchants
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.logger.Debug("Successfully fetched trashed merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *merchantQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
