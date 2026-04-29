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

type bannerQueryService struct {
	observability    observability.TraceLoggerObservability
	cache            cache.BannerQueryCache
	bannerRepository repository.BannerQueryRepository
	logger           logger.LoggerInterface
}

type BannerQueryServiceDeps struct {
	Observability    observability.TraceLoggerObservability
	Cache            cache.BannerQueryCache
	BannerRepository repository.BannerQueryRepository
	Logger           logger.LoggerInterface
}

func NewBannerQueryService(deps *BannerQueryServiceDeps) *bannerQueryService {
	return &bannerQueryService{
		observability:    deps.Observability,
		cache:            deps.Cache,
		bannerRepository: deps.BannerRepository,
		logger:           deps.Logger,
	}
}

func (s *bannerQueryService) FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedBannersCache(ctx, req); found {
		logSuccess("Successfully retrieved all banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	banners, err := s.bannerRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetBannersRow](
			s.logger,
			banner_errors.ErrFailedFindAllBanners.WithInternal(err),
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(banners) > 0 {
		totalCount = int(banners[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedBannersCache(ctx, req, banners, &totalCount)

	logSuccess("Successfully fetched all banners",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return banners, &totalCount, nil
}

func (s *bannerQueryService) FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, *int, error) {
	const method = "FindActive"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedBannerActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	banners, err := s.bannerRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetBannersActiveRow](
			s.logger,
			banner_errors.ErrFailedFindActiveBanners.WithInternal(err),
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(banners) > 0 {
		totalCount = int(banners[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedBannerActiveCache(ctx, req, banners, &totalCount)

	logSuccess("Successfully fetched active banners",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return banners, &totalCount, nil
}

func (s *bannerQueryService) FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, *int, error) {
	const method = "FindTrashed"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedBannerTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed banner records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	banners, err := s.bannerRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetBannersTrashedRow](
			s.logger,
			banner_errors.ErrFailedFindTrashedBanners.WithInternal(err),
			method,
			span,

			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	var totalCount int

	if len(banners) > 0 {
		totalCount = int(banners[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedBannerTrashedCache(ctx, req, banners, &totalCount)

	logSuccess("Successfully fetched trashed banners",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return banners, &totalCount, nil
}

func (s *bannerQueryService) FindByID(ctx context.Context, bannerID int) (*db.GetBannerRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("bannerID", bannerID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedBannerCache(ctx, bannerID); found {
		logSuccess("Successfully retrieved banner from cache", zap.Int("bannerID", bannerID))
		return data, nil
	}

	res, err := s.bannerRepository.FindByID(ctx, bannerID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetBannerRow](
			s.logger,
			banner_errors.ErrBannerNotFoundRes.WithInternal(err),
			method,
			span,

			zap.Int("banner_id", bannerID),
		)
	}

	s.cache.SetCachedBannerCache(ctx, res)

	logSuccess("Successfully fetched banner", zap.Int("bannerID", bannerID))
	return res, nil
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
