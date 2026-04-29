package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantAwardCommandService struct {
	observability           observability.TraceLoggerObservability
	cache                   cache.MerchantAwardCommandCache
	merchantAwardRepository repository.MerchantAwardCommandRepository
	logger                  logger.LoggerInterface
}

type MerchantAwardCommandServiceDeps struct {
	Observability           observability.TraceLoggerObservability
	Cache                   cache.MerchantAwardCommandCache
	MerchantAwardRepository repository.MerchantAwardCommandRepository
	Logger                  logger.LoggerInterface
}

func NewMerchantAwardCommandService(deps *MerchantAwardCommandServiceDeps) MerchantAwardCommandService {
	return &merchantAwardCommandService{
		observability:           deps.Observability,
		cache:                   deps.Cache,
		merchantAwardRepository: deps.MerchantAwardRepository,
		logger:                  deps.Logger,
	}
}

func (s *merchantAwardCommandService) Create(ctx context.Context, req *requests.CreateMerchantCertificationOrAwardRequest) (*db.CreateMerchantCertificationOrAwardRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", req.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.Create(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantCertificationOrAwardRow](
			s.logger,
			merchantaward_errors.ErrFailedCreateMerchantAward,
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Merchant Award created", zap.Int("merchantAward.id", int(merchant.MerchantCertificationID)))

	return merchant, nil
}

func (s *merchantAwardCommandService) Update(ctx context.Context, req *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", *req.MerchantCertificationID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.Update(ctx, req)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantCertificationOrAwardRow](
			s.logger,
			merchantaward_errors.ErrFailedUpdateMerchantAward,
			method,
			span,
			zap.Error(err),
		)
	}

	s.cache.DeleteMerchantAwardCache(ctx, *req.MerchantCertificationID)

	logSuccess("Merchant Award updated", zap.Int("merchantAward.id", int(merchant.MerchantCertificationID)))

	return merchant, nil
}

func (s *merchantAwardCommandService) Trash(ctx context.Context, merchantID int) (*db.MerchantCertificationsAndAward, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.Trash(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantCertificationsAndAward](
			s.logger,
			merchantaward_errors.ErrFailedTrashedMerchantAward,
			method,
			span,
			zap.Error(err),
		)
	}

	s.cache.DeleteMerchantAwardCache(ctx, merchantID)

	logSuccess("Merchant Award trashed", zap.Int("merchantAward.id", int(merchant.MerchantCertificationID)))

	return merchant, nil
}

func (s *merchantAwardCommandService) Restore(ctx context.Context, merchantID int) (*db.MerchantCertificationsAndAward, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.Restore(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantCertificationsAndAward](
			s.logger,
			merchantaward_errors.ErrFailedRestoreMerchantAward,
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Merchant Award restored", zap.Int("merchantAward.id", int(merchant.MerchantCertificationID)))

	return merchant, nil
}

func (s *merchantAwardCommandService) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.DeletePermanent(ctx, merchantID)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantaward_errors.ErrFailedDeleteMerchantAwardPermanent,
			method,
			span,
			zap.Error(err),
		)
	}

	s.cache.DeleteMerchantAwardCache(ctx, merchantID)

	logSuccess("Successfully deleted merchant permanently", zap.Int("merchantAward.id", merchantID), zap.Bool("success", success))

	return success, nil
}

func (s *merchantAwardCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.RestoreAll(ctx)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantaward_errors.ErrFailedRestoreAllMerchantAwards,
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("All trashed merchants restored", zap.Bool("success", success))

	return success, nil
}

func (s *merchantAwardCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.DeleteAll(ctx)

	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantaward_errors.ErrFailedDeleteAllMerchantAwardsPermanent,
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully deleted all merchants permanently", zap.Bool("success", success))

	return success, nil
}
