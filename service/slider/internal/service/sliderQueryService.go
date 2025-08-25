package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type sliderQueryService struct {
	errorhandler          errorhandler.SliderQueryError
	mencache              mencache.SliderQueryCache
	trace                 trace.Tracer
	sliderQueryRepository repository.SliderQueryRepository
	logger                logger.LoggerInterface
	mapping               response_service.SliderResponseMapper
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewSliderQueryService(
	errorhandler errorhandler.SliderQueryError,
	mencache mencache.SliderQueryCache,
	sliderQueryRepository repository.SliderQueryRepository, logger logger.LoggerInterface, mapping response_service.SliderResponseMapper) *sliderQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slider_query_service_request_count",
			Help: "Total number of requests to the SliderQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "slider_query_service_request_duration",
			Help:    "Histogram of request durations for the SliderQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &sliderQueryService{
		errorhandler:          errorhandler,
		mencache:              mencache,
		trace:                 otel.Tracer("slider-query-service"),
		sliderQueryRepository: sliderQueryRepository,
		logger:                logger,
		mapping:               mapping,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *sliderQueryService) FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetSliderAllCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	sliders, totalRecords, err := s.sliderQueryRepository.FindAllSlider(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(
			err,
			method,
			"FAILED_TO_FIND_SLIDERS",
			span,
			&status,
			zap.Error(err),
		)
	}

	slidersResponse := s.mapping.ToSlidersResponse(sliders)

	s.mencache.SetSliderAllCache(ctx, req, slidersResponse, totalRecords)

	logSuccess("Successfully fetched all sliders", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return slidersResponse, totalRecords, nil
}

func (s *sliderQueryService) FindByActive(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetSliderActiveCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	sliders, totalRecords, err := s.sliderQueryRepository.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_ACTIVE_SLIDERS", span, &status, slider_errors.ErrFailedFindActiveSliders, zap.Error(err))
	}

	so := s.mapping.ToSlidersResponseDeleteAt(sliders)

	s.mencache.SetSliderActiveCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched active sliders", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *sliderQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetSliderTrashedCache(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

		return data, total, nil
	}

	sliders, totalRecords, err := s.sliderQueryRepository.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_TO_FIND_TRASHED_SLIDERS", span, &status, slider_errors.ErrFailedFindTrashedSliders)
	}

	so := s.mapping.ToSlidersResponseDeleteAt(sliders)
	s.mencache.SetSliderTrashedCache(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched trashed sliders", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return so, totalRecords, nil
}

func (s *sliderQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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
func (s *sliderQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
func (s *sliderQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
