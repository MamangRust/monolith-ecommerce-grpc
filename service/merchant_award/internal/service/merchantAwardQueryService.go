package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	errorhandler                 errorhandler.MerchantAwardQueryError
	mencache                     mencache.MerchantAwardQueryCache
	trace                        trace.Tracer
	merchantAwardQueryRepositroy repository.MerchantAwardQueryRepository
	logger                       logger.LoggerInterface
	mapping                      response_service.MerchantAwardResponseMapper
	requestCounter               *prometheus.CounterVec
	requestDuration              *prometheus.HistogramVec
}

func NewMerchantAwardQueryService(
	errorhandler errorhandler.MerchantAwardQueryError,
	mencache mencache.MerchantAwardQueryCache,
	merchantAwardQueryRepositroy repository.MerchantAwardQueryRepository, logger logger.LoggerInterface, mapping response_service.MerchantAwardResponseMapper) *merchantAwardQueryService {
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantAwardQueryService{
		errorhandler:                 errorhandler,
		mencache:                     mencache,
		trace:                        otel.Tracer("merchant-award-query-service"),
		merchantAwardQueryRepositroy: merchantAwardQueryRepositroy,
		logger:                       logger,
		mapping:                      mapping,
		requestCounter:               requestCounter,
		requestDuration:              requestDuration,
	}
}

func (s *merchantAwardQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantAwardAll(ctx, req); found {
		logSuccess("Successfully fetched merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantAwardQueryRepositroy.FindAllMerchants(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_TO_FIND_MERCHANTS", span, &status)
	}

	so := s.mapping.ToMerchantsAwardResponse(merchants)
	s.mencache.SetCachedMerchantAwardAll(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantAwardQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantAwardActive(ctx, req); found {
		logSuccess("Successfully fetched active merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantAwardQueryRepositroy.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_ACTIVE_MERCHANTS", span, &status, merchantaward_errors.ErrFailedFindActiveMerchantAwards, zap.Error(err))
	}
	so := s.mapping.ToMerchantsAwardResponseDeleteAt(merchants)
	s.mencache.SetCachedMerchantAwardActive(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantAwardQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantAwardTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantAwardQueryRepositroy.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_TRASHED_MERCHANTS", span, &status, merchantaward_errors.ErrFailedFindTrashedMerchantAwards, zap.Error(err))
	}

	so := s.mapping.ToMerchantsAwardResponseDeleteAt(merchants)
	s.mencache.SetCachedMerchantAwardTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantAwardQueryService) FindById(ctx context.Context, merchantID int) (*response.MerchantAwardResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchantAward(ctx, merchantID); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchantAward.id", merchantID))

		return data, nil
	}

	merchant, err := s.merchantAwardQueryRepositroy.FindById(ctx, merchantID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_TO_FIND_MERCHANT_AWARD_BY_ID", span, &status, merchantaward_errors.ErrFailedFindMerchantAwardById, zap.Error(err))
	}

	so := s.mapping.ToMerchantAwardResponse(merchant)

	s.mencache.SetCachedMerchantAward(ctx, so)

	logSuccess("Successfully fetched merchant", zap.Int("merchantAward.id", merchantID))

	return so, nil
}

func (s *merchantAwardQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *merchantAwardQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantAwardQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
