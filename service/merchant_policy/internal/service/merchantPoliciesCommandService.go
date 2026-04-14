package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantPoliciesCommandService struct {
	observability            observability.TraceLoggerObservability
	cache                    mencache.MerchantPoliciesCommandCache
	merchantPolicyRepository repository.MerchantPoliciesCommandRepository
	logger                   logger.LoggerInterface
}

type MerchantPoliciesCommandServiceDeps struct {
	Observability            observability.TraceLoggerObservability
	Cache                    mencache.MerchantPoliciesCommandCache
	MerchantPolicyRepository repository.MerchantPoliciesCommandRepository
	Logger                   logger.LoggerInterface
}

func NewMerchantPoliciesCommandService(deps *MerchantPoliciesCommandServiceDeps) MerchantPoliciesCommandService {
	return &merchantPoliciesCommandService{
		observability:            deps.Observability,
		cache:                    deps.Cache,
		merchantPolicyRepository: deps.MerchantPolicyRepository,
		logger:                   deps.Logger,
	}
}

func (s *merchantPoliciesCommandService) CreateMerchantPolicy(ctx context.Context, req *requests.CreateMerchantPolicyRequest) (*db.CreateMerchantPolicyRow, error) {
	const method = "CreateMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", req.MerchantID))
	defer func() {
		end(status)
	}()

	policy, err := s.merchantPolicyRepository.CreateMerchantPolicy(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantPolicyRow](
			s.logger,
			merchant_policy_errors.ErrFailedCreateMerchantPolicy.WithInternal(err),
			method,
			span,
			zap.Int("merchant_id", req.MerchantID),
		)
	}

	logSuccess("Successfully created merchant policy", zap.Int("merchant_id", req.MerchantID))
	return policy, nil
}

func (s *merchantPoliciesCommandService) UpdateMerchantPolicy(ctx context.Context, req *requests.UpdateMerchantPolicyRequest) (*db.UpdateMerchantPolicyRow, error) {
	const method = "UpdateMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantPolicy.id", *req.MerchantPolicyID))
	defer func() {
		end(status)
	}()

	policy, err := s.merchantPolicyRepository.UpdateMerchantPolicy(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantPolicyRow](
			s.logger,
			merchant_policy_errors.ErrFailedUpdateMerchantPolicy.WithInternal(err),
			method,
			span,
			zap.Int("merchantPolicy_id", *req.MerchantPolicyID),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, *req.MerchantPolicyID)

	logSuccess("Successfully updated merchant policy", zap.Int("merchantPolicy_id", *req.MerchantPolicyID))
	return policy, nil
}

func (s *merchantPoliciesCommandService) TrashedMerchantPolicy(ctx context.Context, id int) (*db.MerchantPolicy, error) {
	const method = "TrashedMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantPolicy.id", id))
	defer func() {
		end(status)
	}()

	policy, err := s.merchantPolicyRepository.TrashedMerchantPolicy(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantPolicy](
			s.logger,
			merchant_policy_errors.ErrFailedTrashedReviewPolicy.WithInternal(err),
			method,
			span,
			zap.Int("merchantPolicy_id", id),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, id)

	logSuccess("Successfully trashed merchant policy", zap.Int("merchantPolicy_id", id))
	return policy, nil
}

func (s *merchantPoliciesCommandService) RestoreMerchantPolicy(ctx context.Context, id int) (*db.MerchantPolicy, error) {
	const method = "RestoreMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantPolicy.id", id))
	defer func() {
		end(status)
	}()

	policy, err := s.merchantPolicyRepository.RestoreMerchantPolicy(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantPolicy](
			s.logger,
			merchant_policy_errors.ErrFailedRestoreReviewPolicy.WithInternal(err),
			method,
			span,
			zap.Int("merchantPolicy_id", id),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, id)

	logSuccess("Successfully restored merchant policy", zap.Int("merchantPolicy_id", id))
	return policy, nil
}

func (s *merchantPoliciesCommandService) DeleteMerchantPolicyPermanent(ctx context.Context, id int) (bool, error) {
	const method = "DeleteMerchantPolicyPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantPolicy.id", id))
	defer func() {
		end(status)
	}()

	success, err := s.merchantPolicyRepository.DeleteMerchantPolicyPermanent(ctx, id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_policy_errors.ErrFailedDeleteReviewPolicyPermanent.WithInternal(err),
			method,
			span,
			zap.Int("merchantPolicy_id", id),
		)
	}

	s.cache.DeleteMerchantPolicyCache(ctx, id)

	logSuccess("Successfully permanently deleted merchant policy", zap.Int("merchantPolicy_id", id))
	return success, nil
}

func (s *merchantPoliciesCommandService) RestoreAllMerchantPolicy(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchantPolicy"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantPolicyRepository.RestoreAllMerchantPolicy(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_policy_errors.ErrFailedRestoreAllReviewPolicies.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully restored all merchant policies")
	return success, nil
}

func (s *merchantPoliciesCommandService) DeleteAllMerchantPolicyPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantPolicyPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	success, err := s.merchantPolicyRepository.DeleteAllMerchantPolicyPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_policy_errors.ErrFailedDeleteAllReviewPoliciesPermanent.WithInternal(err),
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all merchant policies")
	return success, nil
}
