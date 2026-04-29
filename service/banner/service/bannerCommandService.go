package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-banner/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type bannerCommandService struct {
	observability    observability.TraceLoggerObservability
	cache            cache.BannerCommandCache
	bannerRepository repository.BannerCommandRepository
	logger           logger.LoggerInterface
}

type BannerCommandServiceDeps struct {
	Observability    observability.TraceLoggerObservability
	Cache            cache.BannerCommandCache
	BannerRepository repository.BannerCommandRepository
	Logger           logger.LoggerInterface
}

func NewBannerCommandService(deps *BannerCommandServiceDeps) *bannerCommandService {
	return &bannerCommandService{
		observability:    deps.Observability,
		cache:            deps.Cache,
		bannerRepository: deps.BannerRepository,
		logger:           deps.Logger,
	}
}

func (s *bannerCommandService) Create(ctx context.Context, req *requests.CreateBannerRequest) (*db.CreateBannerRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.Create(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateBannerRow](
			s.logger,
			banner_errors.ErrFailedCreateBanner.WithInternal(err),
			method,
			span,
			zap.Any("request", req),
		)
	}

	logSuccess("Successfully created banner", zap.Int("bannerID", int(banner.BannerID)))
	return banner, nil
}

func (s *bannerCommandService) Update(ctx context.Context, req *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", *req.BannerID))

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateBannerRow](
			s.logger,
			banner_errors.ErrFailedUpdateBanner.WithInternal(err),
			method,
			span,

			zap.Int("banner_id", *req.BannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, *req.BannerID)

	logSuccess("Successfully updated banner", zap.Int("bannerID", *req.BannerID))
	return banner, nil
}

func (s *bannerCommandService) Trash(ctx context.Context, bannerID int) (*db.Banner, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.Trash(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Banner](
			s.logger,
			banner_errors.ErrFailedTrashedBanner.WithInternal(err),
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, bannerID)

	logSuccess("Successfully trashed banner", zap.Int("bannerID", bannerID))
	return banner, nil
}

func (s *bannerCommandService) Restore(ctx context.Context, bannerID int) (*db.Banner, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	banner, err := s.bannerRepository.Restore(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Banner](
			s.logger,
			banner_errors.ErrFailedRestoreBanner.WithInternal(err),
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, bannerID)

	logSuccess("Successfully restored banner", zap.Int("bannerID", bannerID))
	return banner, nil
}

func (s *bannerCommandService) DeletePermanent(ctx context.Context, bannerID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	success, err := s.bannerRepository.DeletePermanent(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			banner_errors.ErrFailedDeleteBanner.WithInternal(err),
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.DeleteBannerCache(ctx, bannerID)

	logSuccess("Successfully deleted banner permanently", zap.Int("bannerID", bannerID))
	return success, nil
}

func (s *bannerCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.bannerRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			banner_errors.ErrFailedRestoreAllBanners.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed banners")
	return success, nil
}

func (s *bannerCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.bannerRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			banner_errors.ErrFailedDeleteAllBanners.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all trashed banners permanently")
	return success, nil
}
