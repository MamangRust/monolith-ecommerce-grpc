package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
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
	ctx                   context.Context
	trace                 trace.Tracer
	sliderQueryRepository repository.SliderQueryRepository
	logger                logger.LoggerInterface
	mapping               response_service.SliderResponseMapper
	requestCounter        *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
}

func NewSliderQueryService(ctx context.Context, sliderQueryRepository repository.SliderQueryRepository, logger logger.LoggerInterface, mapping response_service.SliderResponseMapper) *sliderQueryService {
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
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &sliderQueryService{
		ctx:                   ctx,
		trace:                 otel.Tracer("slider-query-service"),
		sliderQueryRepository: sliderQueryRepository,
		logger:                logger,
		mapping:               mapping,
		requestCounter:        requestCounter,
		requestDuration:       requestDuration,
	}
}

func (s *sliderQueryService) FindAll(req *requests.FindAllSlider) ([]*response.SliderResponse, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching sliders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	sliders, totalRecords, err := s.sliderQueryRepository.FindAllSlider(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_SLIDERS")

		s.logger.Error("Failed to retrieve sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve sliders")

		status = "failed_find_all_sliders"

		return nil, nil, slider_errors.ErrFailedFindAllSliders
	}

	slidersResponse := s.mapping.ToSlidersResponse(sliders)

	s.logger.Debug("Successfully fetched sliders",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return slidersResponse, totalRecords, nil
}

func (s *sliderQueryService) FindByActive(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse) {
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

	s.logger.Debug("Fetching sliders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	sliders, totalRecords, err := s.sliderQueryRepository.FindByActive(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_SLIDERS")

		s.logger.Error("Failed to retrieve active sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve active sliders")

		status = "failed_find_active_sliders"

		return nil, nil, slider_errors.ErrFailedFindActiveSliders
	}

	s.logger.Debug("Successfully fetched sliders",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToSlidersResponseDeleteAt(sliders), totalRecords, nil
}

func (s *sliderQueryService) FindByTrashed(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindByTrashed", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindByTrashed")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching sliders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	sliders, totalRecords, err := s.sliderQueryRepository.FindByTrashed(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_SLIDERS")

		s.logger.Error("Failed to retrieve trashed sliders",
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.String("trace_id", traceID))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve trashed sliders")

		status = "failed_find_trashed_sliders"

		return nil, nil, slider_errors.ErrFailedFindTrashedSliders
	}

	s.logger.Debug("Successfully fetched slider",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return s.mapping.ToSlidersResponseDeleteAt(sliders), totalRecords, nil
}

func (s *sliderQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
