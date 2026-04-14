package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/repository"
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

func (s *merchantAwardCommandService) CreateMerchantAward(ctx context.Context, req *requests.CreateMerchantCertificationOrAwardRequest) (*db.CreateMerchantCertificationOrAwardRow, error) {
	const method = "CreateMerchantAward"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", req.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.CreateMerchantAward(ctx, req)

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

func (s *merchantAwardCommandService) UpdateMerchantAward(ctx context.Context, req *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error) {
	const method = "UpdateMerchantAward"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", *req.MerchantCertificationID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.UpdateMerchantAward(ctx, req)

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

func (s *merchantAwardCommandService) TrashedMerchantAward(ctx context.Context, merchantID int) (*db.MerchantCertificationsAndAward, error) {
	const method = "TrashedMerchantAward"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.TrashedMerchantAward(ctx, merchantID)

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

func (s *merchantAwardCommandService) RestoreMerchantAward(ctx context.Context, merchantID int) (*db.MerchantCertificationsAndAward, error) {
	const method = "RestoreMerchantAward"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantAwardRepository.RestoreMerchantAward(ctx, merchantID)

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

func (s *merchantAwardCommandService) DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeleteMerchantPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantAward.id", merchantID))

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.DeleteMerchantPermanent(ctx, merchantID)

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

func (s *merchantAwardCommandService) RestoreAllMerchantAward(ctx context.Context) (bool, error) {
	const method = "RestoreAllMerchantAward"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.RestoreAllMerchantAward(ctx)

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

func (s *merchantAwardCommandService) DeleteAllMerchantAwardPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllMerchantAwardPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantAwardRepository.DeleteAllMerchantAwardPermanent(ctx)

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
