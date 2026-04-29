package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantDocumentQueryService struct {
	observability observability.TraceLoggerObservability
	cache         cache.MerchantDocumentQueryCache
	repository    repository.MerchantDocumentQueryRepository
	logger        logger.LoggerInterface
}

type MerchantDocumentQueryServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Cache         cache.MerchantDocumentQueryCache
	Repository    repository.MerchantDocumentQueryRepository
	Logger        logger.LoggerInterface
}

func NewMerchantDocumentQueryService(deps *MerchantDocumentQueryServiceDeps) MerchantDocumentQueryService {
	return &merchantDocumentQueryService{
		observability: deps.Observability,
		cache:         deps.Cache,
		repository:    deps.Repository,
		logger:        deps.Logger,
	}
}

func (s *merchantDocumentQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error) {
	const method = "FindAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", req.Page), attribute.Int("pageSize", req.PageSize), attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	res, total, err := s.repository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantDocumentsRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	logSuccess("Successfully fetched merchant documents", zap.Int("totalCount", *total))

	return res, total, nil
}

func (s *merchantDocumentQueryService) FindByID(ctx context.Context, documentID int) (*db.GetMerchantDocumentRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", documentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.FindByID(ctx, documentID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantDocumentRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("documentID", documentID),
		)
	}

	logSuccess("Successfully fetched merchant document", zap.Int("documentID", int(res.DocumentID)))

	return res, nil
}

func (s *merchantDocumentQueryService) FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error) {
	const method = "FindActive"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", req.Page), attribute.Int("pageSize", req.PageSize), attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	res, total, err := s.repository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetActiveMerchantDocumentsRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	logSuccess("Successfully fetched active merchant documents", zap.Int("totalCount", *total))

	return res, total, nil
}

func (s *merchantDocumentQueryService) FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error) {
	const method = "FindTrashed"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("page", req.Page), attribute.Int("pageSize", req.PageSize), attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	res, total, err := s.repository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetTrashedMerchantDocumentsRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
		)
	}

	logSuccess("Successfully fetched trashed merchant documents", zap.Int("totalCount", *total))

	return res, total, nil
}
