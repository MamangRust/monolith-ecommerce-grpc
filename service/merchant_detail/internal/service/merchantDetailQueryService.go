package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantDetailQueryService struct {
	ctx                           context.Context
	errorhandler                  errorhandler.MerchantDetailQueryError
	mencache                      mencache.MerchantDetailQueryCache
	trace                         trace.Tracer
	merchantDetailQueryRepository repository.MerchantDetailQueryRepository
	mapping                       response_service.MerchantDetailResponseMapper
	logger                        logger.LoggerInterface
	requestCounter                *prometheus.CounterVec
	requestDuration               *prometheus.HistogramVec
}

func NewMerchantDetailQueryService(ctx context.Context,
	errorhandler errorhandler.MerchantDetailQueryError,
	mencache mencache.MerchantDetailQueryCache, merchantDetailQueryRepository repository.MerchantDetailQueryRepository,
	mapping response_service.MerchantDetailResponseMapper, logger logger.LoggerInterface) *merchantDetailQueryService {

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_detail_query_service_requests_total",
			Help: "Total number of requests to the MerchantDetailQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_detail_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantDetailQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantDetailQueryService{
		ctx:                           ctx,
		errorhandler:                  errorhandler,
		mencache:                      mencache,
		trace:                         otel.Tracer("merchant-detail-query-service"),
		merchantDetailQueryRepository: merchantDetailQueryRepository,
		mapping:                       mapping,
		logger:                        logger,
		requestCounter:                requestCounter,
		requestDuration:               requestDuration,
	}
}

func (s *merchantDetailQueryService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDetailAll(req); found {
		logSuccess("Successfully fetched merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantDetailQueryRepository.FindAllMerchants(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_DETAIL", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantsDetailResponse(merchants)

	s.mencache.SetCachedMerchantDetailAll(req, so, totalRecords)

	logSuccess("Successfully fetched merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantDetailQueryService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDetailActive(req); found {
		logSuccess("Successfully fetched active merchant from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantDetailQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE", span, &status, merchantdetail_errors.ErrFailedFindActiveMerchantDetail, zap.Error(err))
	}

	so := s.mapping.ToMerchantsDetailResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantDetailActive(req, so, totalRecords)

	logSuccess("Successfully fetched active merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantDetailQueryService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantDetailTrashed(req); found {
		logSuccess("Successfully fetched trashed merchant from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantDetailQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED", span, &status, merchantdetail_errors.ErrFailedFindTrashedMerchantDetail, zap.Error(err))
	}

	so := s.mapping.ToMerchantsDetailResponseDeleteAt(merchants)

	s.mencache.SetCachedMerchantDetailTrashed(req, so, totalRecords)

	logSuccess("Successfully fetched trashed merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantDetailQueryService) FindById(merchantID int) (*response.MerchantDetailResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantDetail.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchantDetail(merchantID); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchantDetail.id", merchantID))

		return data, nil
	}

	merchant, err := s.merchantDetailQueryRepository.FindById(merchantID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.MerchantDetailResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchantdetail_errors.ErrFailedFindMerchantDetailById, zap.Error(err))
	}

	so := s.mapping.ToMerchantDetailResponse(merchant)

	s.mencache.SetCachedMerchantDetail(so)

	logSuccess("Successfully fetched merchant", zap.Int("merchantDetail.id", merchantID))

	return so, nil
}

func (s *merchantDetailQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantDetailQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantDetailQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
