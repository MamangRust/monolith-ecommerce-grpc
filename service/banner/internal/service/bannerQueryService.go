package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
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
	ctx                   context.Context
	trace                 trace.Tracer
	bannerQueryRepository repository.BannerQueryRepository
	logger                logger.LoggerInterface
	mapping               response_service.BannerResponseMapper
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewBannerQueryService(ctx context.Context, bannerQuery repository.BannerQueryRepository, logger logger.LoggerInterface, mapping response_service.BannerResponseMapper) *bannerQueryService {
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
		ctx:                   ctx,
		trace:                 otel.Tracer("banner-query-service"),
		bannerQueryRepository: bannerQuery,
		logger:                logger,
		mapping:               mapping,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *bannerQueryService) FindAll(req *requests.FindAllBanner) ([]*response.BannerResponse, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all Banners",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Banners, totalRecords, err := s.bannerQueryRepository.FindAllBanners(req)

	if err != nil {
		traceunic := traceunic.GenerateTraceID("FAILED_FIND_ALL_BANNERS")

		s.logger.Error("Failed to fetch Banners",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceunic))

		span.SetAttributes(attribute.String("trace.id", traceunic))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch Banners")

		status = "failed_to_fetch_banners"

		return nil, nil, banner_errors.ErrFailedFindAllBanners
	}

	s.logger.Debug("Successfully fetched Banners",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToBannersResponse(Banners), totalRecords, nil
}

func (s *bannerQueryService) FindByActive(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all Banners active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Banners, totalRecords, err := s.bannerQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_BANNERS")

		s.logger.Error("Failed to fetch active Banners",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch active Banners")

		status = "failed_to_fetch_active_banners"

		return nil, nil, banner_errors.ErrFailedFindActiveBanners
	}

	s.logger.Debug("Successfully fetched active Banner",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToBannersResponseDeleteAt(Banners), totalRecords, nil
}

func (s *bannerQueryService) FindByTrashed(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching all Banners trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	Banners, totalRecords, err := s.bannerQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_BANNERS")

		s.logger.Error("Failed to fetch trashed Banners",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
			zap.String("search", search),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch trashed Banners")

		status = "failed_to_fetch_trashed_banners"

		return nil, nil, banner_errors.ErrFailedFindTrashedBanners
	}

	s.logger.Debug("Successfully fetched trashed Banner",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToBannersResponseDeleteAt(Banners), totalRecords, nil
}

func (s *bannerQueryService) FindById(BannerID int) (*response.BannerResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindById")
	defer span.End()

	s.logger.Debug("Fetching Banner by ID", zap.Int("BannerID", BannerID))

	Banner, err := s.bannerQueryRepository.FindById(BannerID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_BANNER_BY_ID")

		s.logger.Error("Failed to retrieve Banner details",
			zap.Error(err),
			zap.Int("Banner_id", BannerID),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve Banner details")

		status = "failed_to_retrieve_banner_details"

		return nil, banner_errors.ErrBannerNotFoundRes
	}

	return s.mapping.ToBannerResponse(Banner), nil
}

func (s *bannerQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
