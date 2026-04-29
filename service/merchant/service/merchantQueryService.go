package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantQueryService struct {
	observability      observability.TraceLoggerObservability
	cache              cache.MerchantQueryCache
	merchantRepository repository.MerchantQueryRepository
	logger             logger.LoggerInterface
}

type MerchantQueryServiceDeps struct {
	Observability      observability.TraceLoggerObservability
	Cache              cache.MerchantQueryCache
	MerchantRepository repository.MerchantQueryRepository
	Logger             logger.LoggerInterface
}

func NewMerchantQueryService(deps *MerchantQueryServiceDeps) MerchantQueryService {
	return &merchantQueryService{
		observability:      deps.Observability,
		cache:              deps.Cache,
		merchantRepository: deps.MerchantRepository,
		logger:             deps.Logger,
	}
}

func (s *merchantQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, *int, error) {
	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", req.Page), attribute.Int("pageSize", req.PageSize), attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchants(ctx, req); found {
		logSuccess("Successfully fetched merchants from cache")
		return data, total, nil
	}

	res, err := s.merchantRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsRow](
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
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.cache.SetCachedMerchants(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched merchants from database", zap.Int("totalCount", totalCount))

	return res, &totalCount, nil
}

func (s *merchantQueryService) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, *int, error) {
	const method = "FindActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", req.Page), attribute.Int("pageSize", req.PageSize), attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantActive(ctx, req); found {
		logSuccess("Successfully fetched active merchants from cache")
		return data, total, nil
	}

	res, err := s.merchantRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsActiveRow](
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
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.cache.SetCachedMerchantActive(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched active merchants from database", zap.Int("totalCount", totalCount))

	return res, &totalCount, nil
}

func (s *merchantQueryService) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, *int, error) {
	const method = "FindTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", req.Page), attribute.Int("pageSize", req.PageSize), attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantTrashed(ctx, req); found {
		logSuccess("Successfully fetched trashed merchants from cache")
		return data, total, nil
	}

	res, err := s.merchantRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsTrashedRow](
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
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.cache.SetCachedMerchantTrashed(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched trashed merchants from database", zap.Int("totalCount", totalCount))

	return res, &totalCount, nil
}

func (s *merchantQueryService) FindByID(ctx context.Context, merchantID int) (*db.GetMerchantByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchant(ctx, merchantID); found {
		logSuccess("Successfully fetched merchant from cache")
		return data, nil
	}

	res, err := s.merchantRepository.FindByID(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantByIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchantID", merchantID),
		)
	}

	s.cache.SetCachedMerchant(ctx, res)

	logSuccess("Successfully fetched merchant from database", zap.Int("merchantID", int(res.MerchantID)))

	return res, nil
}
