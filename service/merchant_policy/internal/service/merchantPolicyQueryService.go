package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantPolicyQueryService struct {
	ctx                           context.Context
	trace                         trace.Tracer
	logger                        logger.LoggerInterface
	merchantPolicyQueryRepository repository.MerchantPoliciesQueryRepository
	mapping                       response_service.MerchantPolicyResponseMapper
	requestCounter                *prometheus.CounterVec
	requestDuration               *prometheus.HistogramVec
}

func NewMerchantPolicyQueryService(ctx context.Context, logger logger.LoggerInterface, merchantPolicyQueryRepository repository.MerchantPoliciesQueryRepository, mapping response_service.MerchantPolicyResponseMapper) *merchantPolicyQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_policy_query_service_requests_total",
			Help: "Total number of requests to the MerchantPolicyQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_policy_query_service_request_duration_seconds",
			Help:    "Histogram of request duration for the MerchantPolicyQueryService",
			Buckets: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantPolicyQueryService{
		ctx:                           ctx,
		trace:                         otel.Tracer("merchant-policy-query-service"),
		logger:                        logger,
		mapping:                       mapping,
		merchantPolicyQueryRepository: merchantPolicyQueryRepository,
		requestCounter:                requestCounter,
		requestDuration:               requestDuration,
	}
}

func (s *merchantPolicyQueryService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllMerchantPolicy")
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

	merchants, totalRecords, err := s.merchantPolicyQueryRepository.FindAllMerchantPolicy(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_MERCHANT_POLICY")

		s.logger.Error("Failed to retrieve merchants",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve merchants")

		status = "failed_find_all_merchant_policies"

		return nil, nil, merchantpolicy_errors.ErrFailedFindAllMerchantPolicies
	}

	s.logger.Debug("Successfully fetched merchants",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToMerchantsPolicyResponse(merchants), totalRecords, nil
}

func (s *merchantPolicyQueryService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByActive", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByActiveMerchantPolicy")
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

	merchants, totalRecords, err := s.merchantPolicyQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_MERCHANT_POLICY")

		s.logger.Error("Failed to retrieve active merchants",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve active merchants")

		status = "failed_find_active_merchant_policies"

		return nil, nil, merchantpolicy_errors.ErrFailedFindActiveMerchantPolicies
	}

	s.logger.Debug("Successfully fetched active merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsPolicyResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantPolicyQueryService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashedMerchantPolicy")
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

	merchants, totalRecords, err := s.merchantPolicyQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_MERCHANT_POLICY")

		s.logger.Error("Failed to retrieve trashed merchants",
			zap.Error(err),
			zap.String("traceID", traceID),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve trashed merchants")

		status = "failed_find_trashed_merchant_policies"

		return nil, nil, merchantpolicy_errors.ErrFailedFindTrashedMerchantPolicies
	}

	s.logger.Debug("Successfully fetched trashed merchant",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToMerchantsPolicyResponseDeleteAt(merchants), totalRecords, nil
}

func (s *merchantPolicyQueryService) FindById(merchantID int) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByIdMerchantPolicy")
	defer span.End()

	s.logger.Debug("Fetching merchant by ID", zap.Int("merchantID", merchantID))

	merchant, err := s.merchantPolicyQueryRepository.FindById(merchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT_POLICY_BY_ID")

		s.logger.Error("Failed to retrieve merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)

		span.RecordError(err)

		span.SetStatus(codes.Error, "Failed to retrieve merchant")

		status = "failed_find_merchant_policy_by_id"

		return nil, merchantpolicy_errors.ErrFailedFindMerchantPolicyById
	}

	return s.mapping.ToMerchantPolicyResponse(merchant), nil
}

func (s *merchantPolicyQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
