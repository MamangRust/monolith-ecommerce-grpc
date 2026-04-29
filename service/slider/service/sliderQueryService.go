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

type sliderQueryService struct {
	observability    observability.TraceLoggerObservability
	cache            cache.SliderQueryCache
	sliderRepository repository.SliderQueryRepository
	logger           logger.LoggerInterface
}

type SliderQueryServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.SliderQueryCache
	Repositories  repository.SliderQueryRepository
	Logger        logger.LoggerInterface
}

func NewSliderQueryService(deps *SliderQueryServiceDeps) SliderQueryService {
	return &sliderQueryService{
		observability:    deps.Observability,
		cache:            deps.Cache,
		sliderRepository: deps.Repositories,
		logger:           deps.Logger,
	}
}

func (s *sliderQueryService) FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersRow, *int, error) {
	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetSliderAllCache(ctx, req); found {
		logSuccess("Successfully retrieved sliders from cache")
		return data, total, nil
	}

	sliders, err := s.sliderRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetSlidersRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int
	if len(sliders) > 0 {
		totalCount = int(sliders[0].TotalCount)
	}

	s.cache.SetSliderAllCache(ctx, req, sliders, &totalCount)

	logSuccess("Successfully fetched sliders from repository", zap.Int("totalCount", totalCount))

	return sliders, &totalCount, nil
}

func (s *sliderQueryService) FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersActiveRow, *int, error) {
	const method = "FindActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetSliderActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active sliders from cache")
		return data, total, nil
	}

	sliders, err := s.sliderRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetSlidersActiveRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int
	if len(sliders) > 0 {
		totalCount = int(sliders[0].TotalCount)
	}

	s.cache.SetSliderActiveCache(ctx, req, sliders, &totalCount)

	logSuccess("Successfully fetched active sliders from repository", zap.Int("totalCount", totalCount))

	return sliders, &totalCount, nil
}

func (s *sliderQueryService) FindByID(ctx context.Context, slider_id int) (*db.GetSliderByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("slider_id", slider_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetSliderCache(ctx, slider_id); found {
		logSuccess("Successfully retrieved slider by ID from cache")
		return data, nil
	}

	slider, err := s.sliderRepository.FindByID(ctx, slider_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetSliderByIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("slider_id", slider_id),
		)
	}

	s.cache.SetSliderCache(ctx, slider)

	logSuccess("Successfully fetched slider by ID from repository", zap.Int("slider_id", slider_id))

	return slider, nil
}

func (s *sliderQueryService) FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersTrashedRow, *int, error) {
	const method = "FindTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetSliderTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed sliders from cache")
		return data, total, nil
	}

	sliders, err := s.sliderRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetSlidersTrashedRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int
	if len(sliders) > 0 {
		totalCount = int(sliders[0].TotalCount)
	}

	s.cache.SetSliderTrashedCache(ctx, req, sliders, &totalCount)

	logSuccess("Successfully fetched trashed sliders from repository", zap.Int("totalCount", totalCount))

	return sliders, &totalCount, nil
}
