package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantAwardQueryService struct {
	observability           observability.TraceLoggerObservability
	cache                   cache.MerchantAwardQueryCache
	merchantAwardRepository repository.MerchantAwardQueryRepository
	logger                  logger.LoggerInterface
}

type MerchantAwardQueryServiceDeps struct {
	Observability           observability.TraceLoggerObservability
	Cache                   cache.MerchantAwardQueryCache
	MerchantAwardRepository repository.MerchantAwardQueryRepository
	Logger                  logger.LoggerInterface
}

func NewMerchantAwardQueryService(deps *MerchantAwardQueryServiceDeps) MerchantAwardQueryService {
	return &merchantAwardQueryService{
		observability:           deps.Observability,
		cache:                   deps.Cache,
		merchantAwardRepository: deps.MerchantAwardRepository,
		logger:                  deps.Logger,
	}
}

func (s *merchantAwardQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsRow, *int, error) {
	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantAwardAll(ctx, req); found {
		logSuccess("Successfully fetched merchants from cache",
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return data, total, nil
	}

	merchants, err := s.merchantAwardRepository.FindAll(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantCertificationsAndAwardsRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	var totalCount int
	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	}

	s.cache.SetCachedMerchantAwardAll(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched merchants",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	return merchants, &totalCount, nil
}

func (s *merchantAwardQueryService) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, *int, error) {
	const method = "FindByActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantAwardActive(ctx, req); found {
		logSuccess("Successfully fetched active merchants from cache",
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return data, total, nil
	}

	merchants, err := s.merchantAwardRepository.FindActive(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantCertificationsAndAwardsActiveRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	var totalCount int
	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	}

	s.cache.SetCachedMerchantAwardActive(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched active merchants",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	return merchants, &totalCount, nil
}

func (s *merchantAwardQueryService) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, *int, error) {
	const method = "FindByTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", req.Page),
		attribute.Int("pageSize", req.PageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantAwardTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchants from cache",
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search))

		return data, total, nil
	}

	merchants, err := s.merchantAwardRepository.FindTrashed(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantCertificationsAndAwardsTrashedRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	var totalCount int
	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	}

	s.cache.SetCachedMerchantAwardTrashed(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched trashed merchants",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	return merchants, &totalCount, nil
}

func (s *merchantAwardQueryService) FindByID(ctx context.Context, merchantID int) (*db.GetMerchantCertificationOrAwardRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchantAward(ctx, merchantID); found {
		logSuccess("Successfully fetched merchant from cache", zap.Int("merchantAward.id", merchantID))

		return data, nil
	}

	merchant, err := s.merchantAwardRepository.FindByID(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantCertificationOrAwardRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.cache.SetCachedMerchantAward(ctx, merchant)

	logSuccess("Successfully fetched merchant", zap.Int("merchantAward.id", merchantID))

	return merchant, nil
}
