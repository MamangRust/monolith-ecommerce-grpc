package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
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
	errorhandler            errorhandler.SliderCommandError
	trace                   trace.Tracer
	sliderCommandRepository repository.SliderCommandRepository
	logger                  logger.LoggerInterface
	mapping                 response_service.SliderResponseMapper
	requestCounter          *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
}

func NewSliderCommandService(ctx context.Context,
	errorhandler errorhandler.SliderCommandError,
	sliderCommandRepository repository.SliderCommandRepository, logger logger.LoggerInterface, mapping response_service.SliderResponseMapper) *sliderCommandService {
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
		errorhandler:            errorhandler,
		trace:                   otel.Tracer("slider-command-service"),
		sliderCommandRepository: sliderCommandRepository,
		logger:                  logger,
		mapping:                 mapping,
		requestCounter:          requestCounter,
		requestDuration:         requestDuration,
	}
}

func (s *sliderCommandService) CreateSlider(req *requests.CreateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	const method = "CreateSlider"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.String("slider.name", req.Nama))

	defer func() {
		end(status)
	}()

	s.logger.Debug("Creating new slider")

	slider, err := s.sliderCommandRepository.CreateSlider(req)

	if err != nil {
		return s.errorhandler.HandleCreateSliderError(
			err, method, "FAILED_CREATE_SLIDER", span, &status, zap.Error(err),
		)
	}

	so := s.mapping.ToSliderResponse(slider)

	logSuccess("Successfully created new slider", zap.Int("slider.id", slider.ID), zap.Bool("success", true))

	return so, nil
}

func (s *sliderCommandService) UpdateSlider(req *requests.UpdateSliderRequest) (*response.SliderResponse, *response.ErrorResponse) {
	const method = "UpdateSlider"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.String("slider.name", req.Nama), attribute.Int("slider.id", *req.ID))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderCommandRepository.UpdateSlider(req)

	if err != nil {
		return s.errorhandler.HandleUpdateSliderError(
			err, method, "FAILED_UPDATE_SLIDER", span, &status, zap.Error(err),
		)
	}

	so := s.mapping.ToSliderResponse(slider)

	logSuccess("Successfully updated slider", zap.Int("slider.id", slider.ID), zap.Bool("success", true))

	return so, nil
}

func (s *sliderCommandService) TrashedSlider(slider_id int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedSlider"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("slider.id", slider_id))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderCommandRepository.TrashSlider(slider_id)

	if err != nil {
		return s.errorhandler.HandleTrashedSliderError(err, method, "FAILED_TRASH_SLIDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToSliderResponseDeleteAt(slider)

	logSuccess("Successfully trashed slider", zap.Int("slider.id", slider.ID), zap.Bool("success", true))

	return so, nil
}

func (s *sliderCommandService) RestoreSlider(sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreSlider"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("slider.id", sliderID))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderCommandRepository.RestoreSlider(sliderID)

	if err != nil {
		return s.errorhandler.HandleRestoreSliderError(err, method, "FAILED_RESTORE_SLIDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToSliderResponseDeleteAt(slider)

	logSuccess("Successfully restored slider", zap.Int("slider.id", slider.ID), zap.Bool("success", true))

	return so, nil
}

func (s *sliderCommandService) DeleteSliderPermanent(sliderID int) (bool, *response.ErrorResponse) {
	const method = "DeleteSliderPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("slider.id", sliderID))

	defer func() {
		end(status)
	}()

	success, err := s.sliderCommandRepository.DeleteSliderPermanently(sliderID)
	if err != nil {
		return s.errorhandler.HandleDeleteSliderError(
			err,
			method,
			"FAILED_DELETE_PERMANENT_SLIDER",
			span,
			&status,
			zap.Error(err),
		)
	}

	logSuccess("Successfully permanently deleted slider", zap.Int("slider.id", sliderID), zap.Bool("success", success))

	return success, nil
}

func (s *sliderCommandService) RestoreAllSliders() (bool, *response.ErrorResponse) {
	const method = "RestoreAllSliders"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderCommandRepository.RestoreAllSlider()

	if err != nil {
		return s.errorhandler.HandleRestoreAllSliderError(
			err,
			method,
			"FAILED_RESTORE_ALL_SLIDERS",
			span,
			&status,
			zap.Error(err),
		)
	}

	logSuccess("All trashed sliders restored successfully", zap.Bool("success", success))

	return success, nil
}

func (s *sliderCommandService) DeleteAllSlidersPermanent() (bool, *response.ErrorResponse) {
	const method = "DeleteAllSlidersPermanent"

	span, end, status, logSuccess := s.startTracingAndLogging(method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderCommandRepository.DeleteAllPermanentSlider()

	if err != nil {
		return s.errorhandler.HandleDeleteAllSliderError(
			err,
			method,
			"FAILED_DELETE_ALL_PERMANENT_SLIDERS",
			span,
			&status,
			zap.Error(err),
		)
	}

	logSuccess("All trashed sliders permanently deleted successfully", zap.Bool("success", success))

	return success, nil
}

func (s *sliderCommandService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
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
func (s *sliderCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
