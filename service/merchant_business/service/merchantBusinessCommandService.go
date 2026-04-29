package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantBusinessCommandService struct {
	observability              observability.TraceLoggerObservability
	cache                      mencache.MerchantBusinessCommandCache
	merchantBusinessRepository repository.MerchantBusinessCommandRepository
	logger                     logger.LoggerInterface
}

type MerchantBusinessCommandServiceDeps struct {
	Observability              observability.TraceLoggerObservability
	Cache                      mencache.MerchantBusinessCommandCache
	MerchantBusinessRepository repository.MerchantBusinessCommandRepository
	Logger                     logger.LoggerInterface
}

func NewMerchantBusinessCommandService(deps *MerchantBusinessCommandServiceDeps) MerchantBusinessCommandService {
	return &merchantBusinessCommandService{
		observability:              deps.Observability,
		cache:                      deps.Cache,
		merchantBusinessRepository: deps.MerchantBusinessRepository,
		logger:                     deps.Logger,
	}
}

func (s *merchantBusinessCommandService) Create(ctx context.Context, req *requests.CreateMerchantBusinessInformationRequest) (*db.CreateMerchantBusinessInformationRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.Create(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedCreateMerchantBusiness,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, int(merchant.MerchantBusinessInfoID))

	logSuccess("Successfully created merchant business",
		zap.Int("merchantID", int(merchant.MerchantBusinessInfoID)))

	return merchant, nil
}

func (s *merchantBusinessCommandService) Update(ctx context.Context, req *requests.UpdateMerchantBusinessInformationRequest) (*db.UpdateMerchantBusinessInformationRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", *req.MerchantBusinessInfoID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedUpdateMerchantBusiness,
			method,
			span,

			zap.Any("request", req),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, int(merchant.MerchantBusinessInfoID))

	logSuccess("Successfully updated merchant business",
		zap.Int("merchantID", int(merchant.MerchantBusinessInfoID)))

	return merchant, nil
}

func (s *merchantBusinessCommandService) Trash(ctx context.Context, merchantID int) (*db.MerchantBusinessInformation, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.Trash(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantBusinessInformation](
			s.logger,
			merchantbusiness_errors.ErrFailedTrashedMerchantBusiness,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, merchantID)

	logSuccess("Successfully trashed merchant business",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantBusinessCommandService) Restore(ctx context.Context, merchantID int) (*db.MerchantBusinessInformation, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantBusinessRepository.Restore(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantBusinessInformation](
			s.logger,
			merchantbusiness_errors.ErrFailedRestoreMerchantBusiness,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, merchantID)

	logSuccess("Successfully restored merchant business",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}

func (s *merchantBusinessCommandService) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessRepository.DeletePermanent(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantbusiness_errors.ErrFailedDeleteMerchantBusinessPermanent,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.DeleteMerchantBusinessCache(ctx, merchantID)

	logSuccess("Successfully permanently deleted merchant business",
		zap.Int("merchantID", merchantID))

	return success, nil
}

func (s *merchantBusinessCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantbusiness_errors.ErrFailedRestoreAllMerchantBusiness,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all trashed merchant businesses")

	return success, nil
}

func (s *merchantBusinessCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.merchantBusinessRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchantbusiness_errors.ErrFailedDeleteAllMerchantBusinessPermanent,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all trashed merchant businesses")

	return success, nil
}
