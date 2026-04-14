package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type sliderCommandService struct {
	observability    observability.TraceLoggerObservability
	cache            cache.SliderCommandCache
	sliderRepository repository.SliderCommandRepository
	logger           logger.LoggerInterface
}

type SliderCommandServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.SliderCommandCache
	Repositories  repository.SliderCommandRepository
	Logger        logger.LoggerInterface
}

func NewSliderCommandService(deps *SliderCommandServiceDeps) SliderCommandService {
	return &sliderCommandService{
		observability:    deps.Observability,
		cache:            deps.Cache,
		sliderRepository: deps.Repositories,
		logger:           deps.Logger,
	}
}

func (s *sliderCommandService) CreateSlider(ctx context.Context, req *requests.CreateSliderRequest) (*db.CreateSliderRow, error) {
	const method = "CreateSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("slider", req.Nama))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.CreateSlider(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateSliderRow](
			s.logger,
			err,
			method,
			span,
			zap.String("slider", req.Nama),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully created slider",
		zap.Int("slider_id", int(slider.SliderID)),
		zap.String("slider_name", slider.Name))

	return slider, nil
}

func (s *sliderCommandService) UpdateSlider(ctx context.Context, req *requests.UpdateSliderRequest) (*db.UpdateSliderRow, error) {
	const method = "UpdateSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", *req.ID),
		attribute.String("new_name", req.Nama))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.UpdateSlider(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateSliderRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("slider_id", *req.ID),
			zap.String("new_name", req.Nama),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully updated slider",
		zap.Int("slider_id", int(slider.SliderID)),
		zap.String("slider_name", slider.Name))

	return slider, nil
}

func (s *sliderCommandService) TrashSlider(ctx context.Context, slider_id int) (*db.Slider, error) {
	const method = "TrashSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", slider_id))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.TrashSlider(ctx, slider_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Slider](
			s.logger,
			err,
			method,
			span,
			zap.Int("slider_id", slider_id),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully trashed slider",
		zap.Int("slider_id", int(slider.SliderID)))

	return slider, nil
}

func (s *sliderCommandService) RestoreSlider(ctx context.Context, sliderID int) (*db.Slider, error) {
	const method = "RestoreSlider"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("sliderID", sliderID))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.RestoreSlider(ctx, sliderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Slider](
			s.logger,
			err,
			method,
			span,
			zap.Int("sliderID", sliderID),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully restored slider",
		zap.Int("slider_id", int(slider.SliderID)))

	return slider, nil
}

func (s *sliderCommandService) DeleteSliderPermanently(ctx context.Context, sliderID int) (bool, error) {
	const method = "DeleteSliderPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("sliderID", sliderID))

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.DeleteSliderPermanently(ctx, sliderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("sliderID", sliderID),
		)
	}

	s.cache.InvalidateSliderCache(ctx)

	logSuccess("Successfully permanently deleted slider",
		zap.Int("sliderID", sliderID))

	return success, nil
}

func (s *sliderCommandService) RestoreAllSliders(ctx context.Context) (bool, error) {
	const method = "RestoreAllSliders"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.RestoreAllSlider(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	s.cache.InvalidateSliderCache(ctx)
	logSuccess("Successfully restored all trashed sliders")

	return success, nil
}

func (s *sliderCommandService) DeleteAllPermanentSlider(ctx context.Context) (bool, error) {
	const method = "DeleteAllSlidersPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.DeleteAllPermanentSlider(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	s.cache.InvalidateSliderCache(ctx)
	logSuccess("Successfully permanently deleted all trashed sliders")

	return success, nil
}
