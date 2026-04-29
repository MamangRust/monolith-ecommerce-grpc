package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-slider/repository"
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

func (s *sliderCommandService) Create(ctx context.Context, req *requests.CreateSliderRequest) (*db.CreateSliderRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("slider", req.Nama))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.Create(ctx, req)
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

func (s *sliderCommandService) Update(ctx context.Context, req *requests.UpdateSliderRequest) (*db.UpdateSliderRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", *req.ID),
		attribute.String("new_name", req.Nama))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.Update(ctx, req)
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

func (s *sliderCommandService) Trash(ctx context.Context, slider_id int) (*db.Slider, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", slider_id))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.Trash(ctx, slider_id)
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

func (s *sliderCommandService) Restore(ctx context.Context, sliderID int) (*db.Slider, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("sliderID", sliderID))

	defer func() {
		end(status)
	}()

	slider, err := s.sliderRepository.Restore(ctx, sliderID)
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

func (s *sliderCommandService) DeletePermanent(ctx context.Context, sliderID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("sliderID", sliderID))

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.DeletePermanent(ctx, sliderID)
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

func (s *sliderCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.RestoreAll(ctx)
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

func (s *sliderCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.sliderRepository.DeleteAll(ctx)
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
