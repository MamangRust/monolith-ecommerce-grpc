package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantPoliciesQueryService struct {
	observability            observability.TraceLoggerObservability
	cache                    mencache.MerchantPoliciesQueryCache
	merchantPolicyRepository repository.MerchantPoliciesQueryRepository
	logger                   logger.LoggerInterface
}

type MerchantPoliciesQueryServiceDeps struct {
	Observability            observability.TraceLoggerObservability
	Cache                    mencache.MerchantPoliciesQueryCache
	MerchantPolicyRepository repository.MerchantPoliciesQueryRepository
	Logger                   logger.LoggerInterface
}

func NewMerchantPoliciesQueryService(deps *MerchantPoliciesQueryServiceDeps) MerchantPoliciesQueryService {
	return &merchantPoliciesQueryService{
		observability:            deps.Observability,
		cache:                    deps.Cache,
		merchantPolicyRepository: deps.MerchantPolicyRepository,
		logger:                   deps.Logger,
	}
}

func (s *merchantPoliciesQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, *int, error) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))
	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantPolicyAll(ctx, req); found {
		return data, total, nil
	}

	merchants, err := s.merchantPolicyRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantPoliciesRow](
			s.logger,
			merchant_policy_errors.ErrFailedFindAllMerchantPolicies.WithInternal(err),
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
		)
	}

	totalRecords := len(merchants)
	s.cache.SetCachedMerchantPolicyAll(ctx, req, merchants, &totalRecords)
	logSuccess("Successfully fetched all merchant policies")

	return merchants, &totalRecords, nil
}

func (s *merchantPoliciesQueryService) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, *int, error) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))
	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantPolicyActive(ctx, req); found {
		return data, total, nil
	}

	merchants, err := s.merchantPolicyRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantPoliciesActiveRow](
			s.logger,
			merchant_policy_errors.ErrFailedFindActiveMerchantPolicies.WithInternal(err),
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
		)
	}

	totalRecords := len(merchants)
	s.cache.SetCachedMerchantPolicyActive(ctx, req, merchants, &totalRecords)
	logSuccess("Successfully fetched active merchant policies")

	return merchants, &totalRecords, nil
}

func (s *merchantPoliciesQueryService) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, *int, error) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))
	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantPolicyTrashed(ctx, req); found {
		return data, total, nil
	}

	merchants, err := s.merchantPolicyRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantPoliciesTrashedRow](
			s.logger,
			merchant_policy_errors.ErrFailedFindTrashedMerchantPolicies.WithInternal(err),
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("page_size", req.PageSize),
		)
	}

	totalRecords := len(merchants)
	s.cache.SetCachedMerchantPolicyTrashed(ctx, req, merchants, &totalRecords)
	logSuccess("Successfully fetched trashed merchant policies")

	return merchants, &totalRecords, nil
}

func (s *merchantPoliciesQueryService) FindByID(ctx context.Context, id int) (*db.GetMerchantPolicyRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantPolicy.id", id))
	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchantPolicy(ctx, id); found {
		return data, nil
	}

	merchant, err := s.merchantPolicyRepository.FindByID(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantPolicyRow](
			s.logger,
			merchant_policy_errors.ErrFailedFindMerchantPolicyByID.WithInternal(err),
			method,
			span,
			zap.Int("merchantPolicy_id", id),
		)
	}

	s.cache.SetCachedMerchantPolicy(ctx, merchant)
	logSuccess("Successfully fetched merchant policy by ID", zap.Int("merchantPolicy_id", id))

	return merchant, nil
}

func (s *merchantPoliciesQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
