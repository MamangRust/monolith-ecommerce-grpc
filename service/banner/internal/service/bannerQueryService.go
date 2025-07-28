package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type bannerQueryService struct {
	errorhandler          errorhandler.BannerQueryError
	mencache              mencache.BannerQueryCache
	trace                 trace.Tracer
	bannerQueryRepository repository.BannerQueryRepository
	logger                logger.LoggerInterface
	mapping               response_service.BannerResponseMapper
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewBannerQueryService(
	errorhandler errorhandler.BannerQueryError,
	mencache mencache.BannerQueryCache,
	bannerQuery repository.BannerQueryRepository, logger logger.LoggerInterface, mapping response_service.BannerResponseMapper) *bannerQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "banner_query_service_request_total",
			Help: "Total number of requests to the BannerQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "banner_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the BannerQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &bannerQueryService{
		errorhandler:          errorhandler,
		mencache:              mencache,
		trace:                 otel.Tracer("banner-query-service"),
		bannerQueryRepository: bannerQuery,
		logger:                logger,
		mapping:               mapping,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *bannerQueryService) FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedBannersCache(ctx, req); found {
		logSuccess("Successfully fetched Banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	Banners, totalRecords, err := s.bannerQueryRepository.FindAllBanners(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_BANNERS", span, &status, zap.Error(err))
	}
	so := s.mapping.ToBannersResponse(Banners)

	s.mencache.SetCachedBannersCache(ctx, req, so, totalRecords)

	return so, totalRecords, nil
}

func (s *bannerQueryService) FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedBannerTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched active Banners from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	Banners, totalRecords, err := s.bannerQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE", span, &status, zap.Error(err))
	}

	so := s.mapping.ToBannersResponseDeleteAt(Banners)

	s.mencache.SetCachedBannerActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active Banners", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *bannerQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedBannerTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed Banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	Banners, totalRecords, err := s.bannerQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED", span, &status, zap.Error(err))
	}

	so := s.mapping.ToBannersResponseDeleteAt(Banners)

	s.mencache.SetCachedBannerTrashedCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed Banner records", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *bannerQueryService) FindById(ctx context.Context, BannerID int) (*response.BannerResponse, *response.ErrorResponse) {
	const method = "FindByCardNumber"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("banner.id", BannerID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedBannerCache(ctx, BannerID); found {
		logSuccess("Successfully fetched Banner from cache", zap.Int("banner.id", BannerID))
		return data, nil
	}

	Banner, err := s.bannerQueryRepository.FindById(ctx, BannerID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_BY_ID", span, &status, banner_errors.ErrBannerNotFoundRes, zap.Error(err))
	}

	so := s.mapping.ToBannerResponse(Banner)

	s.mencache.SetCachedBannerCache(ctx, so)

	logSuccess("Successfully fetched Banner", zap.Int("banner.id", int(so.ID)))

	return so, nil
}

func (s *bannerQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *bannerQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *bannerQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
