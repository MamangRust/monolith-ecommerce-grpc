package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	merchant_social_link_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantSocialLinkCommandService struct {
	observability    observability.TraceLoggerObservability
	repository       repository.MerchantSocialLinkCommandRepository
	logger           logger.LoggerInterface
}

type MerchantSocialLinkCommandServiceDeps struct {
	Observability    observability.TraceLoggerObservability
	Repository       repository.MerchantSocialLinkCommandRepository
	Logger           logger.LoggerInterface
}

func NewMerchantSocialLinkCommandService(deps *MerchantSocialLinkCommandServiceDeps) *merchantSocialLinkCommandService {
	return &merchantSocialLinkCommandService{
		observability:    deps.Observability,
		repository:       deps.Repository,
		logger:           deps.Logger,
	}
}

func (s *merchantSocialLinkCommandService) Create(ctx context.Context, req *requests.CreateMerchantSocialRequest) (*db.CreateMerchantSocialMediaLinkRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantDetailID", *req.MerchantDetailID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Create(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantSocialMediaLinkRow](
			s.logger,
			merchant_social_link_errors.ErrCreateMerchantSocialLink,
			method,
			span,
			zap.Int("merchantDetailID", *req.MerchantDetailID),
		)
	}

	logSuccess("Successfully created merchant social link", zap.Int("socialLinkID", int(res.MerchantSocialID)))
	return res, nil
}

func (s *merchantSocialLinkCommandService) Update(ctx context.Context, req *requests.UpdateMerchantSocialRequest) (*db.UpdateMerchantSocialMediaLinkRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("socialLinkID", req.ID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantSocialMediaLinkRow](
			s.logger,
			merchant_social_link_errors.ErrUpdateMerchantSocialLink,
			method,
			span,
			zap.Int("socialLinkID", req.ID),
		)
	}

	logSuccess("Successfully updated merchant social link", zap.Int("socialLinkID", req.ID))
	return res, nil
}

func (s *merchantSocialLinkCommandService) Trash(ctx context.Context, socialID int) (bool, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("socialLinkID", socialID))

	defer func() {
		end(status)
	}()

	success, err := s.repository.Trash(ctx, socialID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_social_link_errors.ErrTrashMerchantSocialLink,
			method,
			span,
			zap.Int("socialLinkID", socialID),
		)
	}

	logSuccess("Successfully trashed merchant social link", zap.Int("socialLinkID", socialID))
	return success, nil
}

func (s *merchantSocialLinkCommandService) Restore(ctx context.Context, socialID int) (bool, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("socialLinkID", socialID))

	defer func() {
		end(status)
	}()

	success, err := s.repository.Restore(ctx, socialID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_social_link_errors.ErrRestoreMerchantSocialLink,
			method,
			span,
			zap.Int("socialLinkID", socialID),
		)
	}

	logSuccess("Successfully restored merchant social link", zap.Int("socialLinkID", socialID))
	return success, nil
}

func (s *merchantSocialLinkCommandService) DeletePermanent(ctx context.Context, socialID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("socialLinkID", socialID))

	defer func() {
		end(status)
	}()

	success, err := s.repository.DeletePermanent(ctx, socialID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_social_link_errors.ErrDeletePermanentMerchantSocialLink,
			method,
			span,
			zap.Int("socialLinkID", socialID),
		)
	}

	logSuccess("Successfully deleted merchant social link permanently", zap.Int("socialLinkID", socialID))
	return success, nil
}

func (s *merchantSocialLinkCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.repository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_social_link_errors.ErrRestoreAllMerchantSocialLinks,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all merchant social links")
	return success, nil
}

func (s *merchantSocialLinkCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.repository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			merchant_social_link_errors.ErrDeleteAllPermanentMerchantSocialLinks,
			method,
			span,
		)
	}

	logSuccess("Successfully deleted all merchant social links permanently")
	return success, nil
}
