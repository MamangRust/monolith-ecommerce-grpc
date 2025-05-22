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

type sliderCommandService struct {
	ctx                     context.Context
	trace                   trace.Tracer
	sliderCommandRepository repository.SliderCommandRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.SliderResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewSliderCommandService(ctx context.Context, sliderCommandRepository repository.SliderCommandRepository, logger logger.LoggerInterface, mapping response_service.SliderResponseMapper) *sliderCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slider_command_service_request_count",
			Help: "Total number of requests to the SliderCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "slider_command_service_request_duration",
			Help:    "Histogram of request durations for the SliderCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &sliderCommandService{
		ctx:                     ctx,
		trace:                   otel.Tracer("slider-command-service"),
		sliderCommandRepository: sliderCommandRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *sliderCommandService) CreateSlider(req *requests.CreateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateSlider", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateSlider")
	defer span.End()

	s.logger.Debug("Creating new slider")

	slider, err := s.sliderCommandRepository.CreateSlider(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_SLIDER")

		s.logger.Error("Failed to create new slider",
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create new slider")

		status = "failed_create_slider"

		return nil, slider_errors.ErrFailedCreateSlider
	}

	return s.mapping.ToSliderResponse(slider), nil
}

func (s *sliderCommandService) UpdateSlider(req *requests.UpdateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateSlider", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateSlider")
	defer span.End()

	s.logger.Debug("Updating slider", zap.Int("slider_id", *req.ID))

	slider, err := s.sliderCommandRepository.UpdateSlider(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_SLIDER")

		s.logger.Error("Failed to update slider",
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update slider")

		status = "failed_update_slider"

		return nil, slider_errors.ErrFailedUpdateSlider
	}

	return s.mapping.ToSliderResponse(slider), nil
}

func (s *sliderCommandService) TrashedSlider(slider_id int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedSlider", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedSlider")
	defer span.End()

	span.SetAttributes(
		attribute.Int("slider", slider_id),
	)

	s.logger.Debug("Trashing slider", zap.Int("slider", slider_id))

	slider, err := s.sliderCommandRepository.TrashSlider(slider_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_SLIDER")

		s.logger.Error("Failed to trash slider",
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash slider")

		status = "failed_trash_slider"

		return nil, slider_errors.ErrFailedTrashSlider
	}

	return s.mapping.ToSliderResponseDeleteAt(slider), nil
}

func (s *sliderCommandService) RestoreSlider(sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreSlider", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreSlider")
	defer span.End()

	span.SetAttributes(
		attribute.Int("sliderID", sliderID),
	)

	s.logger.Debug("Restoring slider", zap.Int("sliderID", sliderID))

	slider, err := s.sliderCommandRepository.RestoreSlider(sliderID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_SLIDER")

		s.logger.Error("Failed to restore slider",
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore slider")

		status = "failed_restore_slider"

		return nil, slider_errors.ErrFailedRestoreSlider
	}

	return s.mapping.ToSliderResponseDeleteAt(slider), nil
}

func (s *sliderCommandService) DeleteSliderPermanent(sliderID int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteSliderPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteSliderPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting slider", zap.Int("sliderID", sliderID))

	success, err := s.sliderCommandRepository.DeleteSliderPermanently(sliderID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_PERMANENT_SLIDER")

		s.logger.Error("Failed to permanently delete slider",
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete slider")

		status = "failed_delete_permanent_slider"

		return false, slider_errors.ErrFailedDeletePermanentSlider
	}

	return success, nil
}

func (s *sliderCommandService) RestoreAllSliders() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllSliders", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllSliders")
	defer span.End()

	s.logger.Debug("Restoring all trashed sliders")

	success, err := s.sliderCommandRepository.RestoreAllSlider()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_SLIDERS")

		s.logger.Error("Failed to restore all trashed sliders",
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed sliders")

		status = "failed_restore_all_sliders"

		return false, slider_errors.ErrFailedRestoreAllSliders
	}

	return success, nil
}

func (s *sliderCommandService) DeleteAllSlidersPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllSlidersPermanent", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllSlidersPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all sliders")

	success, err := s.sliderCommandRepository.DeleteAllPermanentSlider()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_PERMANENT_SLIDERS")

		s.logger.Error("Failed to permanently delete all sliders",
			zap.String("trace_id", traceID),
			zap.Error(err))

		span.SetAttributes(
			attribute.String("trace_id", traceID),
		)

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all sliders")

		status = "failed_delete_all_permanent_sliders"

		return false, slider_errors.ErrFailedDeleteAllPermanentSliders
	}

	return success, nil
}

func (s *sliderCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
