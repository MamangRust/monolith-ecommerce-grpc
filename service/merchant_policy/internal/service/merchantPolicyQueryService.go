package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	errorhandler                  errorhandler.MerchantPolicyQueryError
	mencache                      mencache.MerchantPolicyQueryCache
	trace                         trace.Tracer
	logger                        logger.LoggerInterface
	merchantPolicyQueryRepository repository.MerchantPoliciesQueryRepository
	mapping                       response_service.MerchantPolicyResponseMapper
	requestCounter                *prometheus.CounterVec
	requestDuration               *prometheus.HistogramVec
}

func NewMerchantPolicyQueryService(
	errorhandler errorhandler.MerchantPolicyQueryError,
	mencache mencache.MerchantPolicyQueryCache,
	logger logger.LoggerInterface, merchantPolicyQueryRepository repository.MerchantPoliciesQueryRepository, mapping response_service.MerchantPolicyResponseMapper) *merchantPolicyQueryService {
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
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantPolicyQueryService{
		errorhandler:                  errorhandler,
		mencache:                      mencache,
		trace:                         otel.Tracer("merchant-policy-query-service"),
		logger:                        logger,
		mapping:                       mapping,
		merchantPolicyQueryRepository: merchantPolicyQueryRepository,
		requestCounter:                requestCounter,
		requestDuration:               requestDuration,
	}
}

func (s *merchantPolicyQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, *response.ErrorResponse) {
	const method = "FindAllMerchantPolicy"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantPolicyAll(ctx, req); found {
		logSuccess("Successfully fetched merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantPolicyQueryRepository.FindAllMerchantPolicy(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_POLICY", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantsPolicyResponse(merchants)
	s.mencache.SetCachedMerchantPolicyAll(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantPolicyQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActiveMerchantPolicy"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantPolicyActive(ctx, req); found {
		logSuccess("Successfully fetched active merchant from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantPolicyQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE_MERCHANT_POLICY", span, &status, merchantpolicy_errors.ErrFailedFindActiveMerchantPolicies, zap.Error(err))
	}

	so := s.mapping.ToMerchantsPolicyResponseDeleteAt(merchants)
	s.mencache.SetCachedMerchantPolicyActive(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantPolicyQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashedMerchantPolicy"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantPolicyTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchant from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantPolicyQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED_MERCHANT_POLICY", span, &status, merchantpolicy_errors.ErrFailedFindTrashedMerchantPolicies, zap.Error(err))
	}

	so := s.mapping.ToMerchantsPolicyResponseDeleteAt(merchants)
	s.mencache.SetCachedMerchantPolicyTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantPolicyQueryService) FindById(ctx context.Context, merchantID int) (*response.MerchantPoliciesResponse, *response.ErrorResponse) {
	const method = "FindMerchantPolicyById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchantPolicy.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchantPolicy(ctx, merchantID); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchantPolicy.id", merchantID))

		return data, nil
	}

	merchant, err := s.merchantPolicyQueryRepository.FindById(ctx, merchantID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_POLICY_BY_ID", span, &status, merchantpolicy_errors.ErrFailedFindMerchantPolicyById, zap.Error(err))
	}

	so := s.mapping.ToMerchantPolicyResponse(merchant)

	s.mencache.SetCachedMerchantPolicy(ctx, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchantPolicy.id", merchantID))

	return so, nil
}

func (s *merchantPolicyQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *merchantPolicyQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantPolicyQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
