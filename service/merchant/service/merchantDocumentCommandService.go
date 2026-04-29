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

type merchantDocumentCommandService struct {
	observability observability.TraceLoggerObservability
	cache         cache.MerchantDocumentCommandCache
	repository    repository.MerchantDocumentCommandRepository
	logger        logger.LoggerInterface
}

type MerchantDocumentCommandServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.MerchantDocumentCommandCache
	Repository    repository.MerchantDocumentCommandRepository
	Logger        logger.LoggerInterface
}

func NewMerchantDocumentCommandService(deps *MerchantDocumentCommandServiceDeps) MerchantDocumentCommandService {
	return &merchantDocumentCommandService{
		observability: deps.Observability,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
	}
}

func (s *merchantDocumentCommandService) Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", request.MerchantID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Create(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantDocumentRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchantID", request.MerchantID),
		)
	}

	logSuccess("Successfully created merchant document", zap.Int("merchantID", request.MerchantID))

	return res, nil
}

func (s *merchantDocumentCommandService) Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("document_id", *request.DocumentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Update(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantDocumentRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", *request.DocumentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(res.DocumentID))

	logSuccess("Successfully updated merchant document", zap.Int("document_id", *request.DocumentID))

	return res, nil
}

func (s *merchantDocumentCommandService) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error) {
	const method = "UpdateStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", *request.DocumentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.UpdateStatus(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantDocumentStatusRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", *request.DocumentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(res.DocumentID))

	logSuccess("Successfully updated merchant document status", zap.Int("document_id", *request.DocumentID))

	return res, nil
}

func (s *merchantDocumentCommandService) Trash(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", documentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Trash(ctx, documentID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", documentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(documentID))

	logSuccess("Successfully trashed merchant document", zap.Int("document_id", documentID))

	return res, nil
}

func (s *merchantDocumentCommandService) Restore(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", documentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Restore(ctx, documentID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", documentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(documentID))

	logSuccess("Successfully restored merchant document", zap.Int("document_id", documentID))

	return res, nil
}

func (s *merchantDocumentCommandService) DeletePermanent(ctx context.Context, documentID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", documentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.DeletePermanent(ctx, documentID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", documentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(documentID))

	logSuccess("Successfully permanently deleted merchant document", zap.Int("document_id", documentID))

	return res, nil
}

func (s *merchantDocumentCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.repository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all merchant documents")

	return res, nil
}

func (s *merchantDocumentCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.repository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all merchant documents")

	return res, nil
}
