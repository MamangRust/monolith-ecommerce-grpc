package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantBusinessQueryService struct {
	ctx                             context.Context
	mencache                        mencache.MerchantBusinessQueryCache
	errorhandler                    errorhandler.MerchantBusinessQueryError
	trace                           trace.Tracer
	merchantBusinessQueryRepository repository.MerchantBusinessQueryRepository
	logger                          logger.LoggerInterface
	mapping                         response_service.MerchantBusinessResponseMapper
	requestCounter                  *prometheus.CounterVec
	requestDuration                 *prometheus.HistogramVec
}

func NewMerchantBusinessQueryService(
	ctx context.Context,
	mencache mencache.MerchantBusinessQueryCache,
	errorhandler errorhandler.MerchantBusinessQueryError,
	merchantBusinessQueryRepository repository.MerchantBusinessQueryRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantBusinessResponseMapper,
) *merchantBusinessQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_business_query_service_requests_total",
			Help: "Total number of requests to the MerchantBusinessQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_business_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantBusinessQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantBusinessQueryService{
		ctx:                             ctx,
		mencache:                        mencache,
		errorhandler:                    errorhandler,
		trace:                           otel.Tracer("merchant-business-query-service"),
		merchantBusinessQueryRepository: merchantBusinessQueryRepository,
		logger:                          logger,
		mapping:                         mapping,
		requestCounter:                  requestCounter,
		requestDuration:                 requestDuration,
	}
}

func (s *merchantBusinessQueryService) FindAll(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantBusinessAll(req); found {
		logSuccess("Successfully fetched merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantBusinessQueryRepository.FindAllMerchants(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_MERCHANT_BUSINESS", span, &status, zap.Error(err))
	}

	so := s.mapping.ToMerchantsBusinessResponse(merchants)

	s.mencache.SetCachedMerchantBusinessAll(req, so, totalRecords)

	logSuccess("Successfully fetched merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantBusinessQueryService) FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantBusinessActive(req); found {
		logSuccess("Successfully fetched active merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantBusinessQueryRepository.FindByActive(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE_MERCHANT_BUSINESS", span, &status, merchantbusiness_errors.ErrFailedFindActiveMerchantBusiness, zap.Error(err))
	}

	so := s.mapping.ToMerchantsBusinessResponseDeleteAt(merchants)
	s.mencache.SetCachedMerchantBusinessActive(req, so, totalRecords)

	logSuccess("Successfully fetched active merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil

}

func (s *merchantBusinessQueryService) FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedMerchantBusinessTrashed(req); found {
		logSuccess("Successfully fetched trashed merchants from cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	merchants, totalRecords, err := s.merchantBusinessQueryRepository.FindByTrashed(req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED_MERCHANT_BUSINESS", span, &status, merchantbusiness_errors.ErrFailedFindTrashedMerchantBusiness, zap.Error(err))
	}

	so := s.mapping.ToMerchantsBusinessResponseDeleteAt(merchants)
	s.mencache.SetCachedMerchantBusinessTrashed(req, so, totalRecords)

	logSuccess("Successfully fetched trashed merchants", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *merchantBusinessQueryService) FindById(merchantID int) (*response.MerchantBusinessResponse, *response.ErrorResponse) {
	const method = "FindById"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("merchantBusiness.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMerchantBusiness(merchantID); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchantBusiness.id", merchantID))
		return data, nil
	}

	merchant, err := s.merchantBusinessQueryRepository.FindById(merchantID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchantbusiness_errors.ErrFailedFindMerchantBusinessById, zap.Error(err))
	}

	so := s.mapping.ToMerchantBusinessResponse(merchant)
	s.mencache.SetCachedMerchantBusiness(so)

	logSuccess("Successfully fetched merchant", zap.Int("merchantBusiness.id", merchantID))

	return so, nil
}

func (s *merchantBusinessQueryService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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

func (s *merchantBusinessQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *merchantBusinessQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
