package services

import (
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	errorhandler            errorhandler.MerchantQueryErrorHandler
	mencache                mencache.MerchantQueryCache
	trace                   trace.Tracer
	merchantQueryRepository repository.MerchantQueryRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.MerchantResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewMerchantQueryService(
	errorhandler errorhandler.MerchantQueryErrorHandler,
	mencache mencache.MerchantQueryCache,
	merchantQueryRepository repository.MerchantQueryRepository, logger logger.LoggerInterface, mapping response_service.MerchantResponseMapper) *merchantQueryService {
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantQueryService{
		errorhandler:            errorhandler,
		mencache:                mencache,
		trace:                   otel.Tracer("merchant-query-service"),
		merchantQueryRepository: merchantQueryRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *merchantQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchants(ctx, req); found {
		logSuccess("Successfully fetched merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindAllMerchants(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, "FindAll", "FAILED_FIND_ALL_MERCHANTS", span, &status, zap.Error(err))
	}

	merchantResponses := s.mapping.ToMerchantsResponse(merchants)

	s.mencache.SetCachedMerchants(ctx, req, merchantResponses, totalRecords)

	logSuccess("Successfully fetched merchants from database", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return merchantResponses, totalRecords, nil
}

func (s *merchantQueryService) FindById(ctx context.Context, merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchant(ctx, merchant_id); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchant.id", merchant_id))

		return data, nil
	}

	res, err := s.merchantQueryRepository.FindById(ctx, merchant_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Int("merchant.id", merchant_id))
	}

	so := s.mapping.ToMerchantResponse(res)

	s.mencache.SetCachedMerchant(ctx, so)

	logSuccess("Successfully fetched merchant from database", zap.Int("merchant.id", merchant_id))

	return so, nil
}

func (s *merchantQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantActive(ctx, req); found {
		logSuccess("Successfully fetched active merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ALL_ACTIVE_MERCHANTS", span, &status, merchant_errors.ErrFailedFindActiveMerchants, zap.Error(err))
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	logSuccess("Successfully fetched active merchants from database", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_ALL_TRASHED_MERCHANTS", span, &status, merchant_errors.ErrFailedFindTrashedMerchants, zap.Error(err))
	}

	so := s.mapping.ToMerchantsResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed merchants from database", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *merchantQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
