package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/services"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDocumentHandleGrpc struct {
	pb.UnimplementedMerchantDocumentServiceServer
	merchantDocumentQuery   services.MerchantDocumentQueryService
	merchantDocumentCommand services.MerchantDocumentCommandService
	logger                  logger.LoggerInterface
	mapping                 protomapper.MerchantDocumentProtoMapper
}

func NewMerchantDocumentHandleGrpc(
	service *services.Service,
	mapping protomapper.MerchantDocumentProtoMapper,
	logger logger.LoggerInterface,
) pb.MerchantDocumentServiceServer {
	return &merchantDocumentHandleGrpc{
		merchantDocumentQuery:   service.MerchantDocumentQuery,
		merchantDocumentCommand: service.MerchantDocumentCommand,
		logger:                  logger,
		mapping:                 mapping,
	}
}

func (s *merchantDocumentHandleGrpc) FindAll(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocument, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all merchant documents",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all merchant documents",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched all merchant documents",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	return s.mapping.ToProtoResponsePaginationMerchantDocument(paginationMeta, "success", "Successfully fetched merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) FindById(ctx context.Context, req *pb.FindMerchantDocumentByIdRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		s.logger.Error("Invalid document ID provided for lookup", zap.Int("document_id", id))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	s.logger.Info("Fetching merchant document by ID", zap.Int("document_id", id))

	document, err := s.merchantDocumentQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch merchant document by ID",
			zap.Int("document_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched merchant document by ID",
		zap.Int("document_id", id),
		zap.Int("merchant_id", int(document.MerchantID)),
		zap.String("doc_type", document.DocumentType),
	)

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully fetched merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) FindAllActive(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocument, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active merchant documents",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active merchant documents",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active merchant documents",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	return s.mapping.ToProtoResponsePaginationMerchantDocument(paginationMeta, "success", "Successfully fetched active merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) FindAllTrashed(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocumentAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed merchant documents",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed merchant documents",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed merchant documents",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	return s.mapping.ToProtoResponsePaginationMerchantDocumentDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) Create(ctx context.Context, req *pb.CreateMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	s.logger.Info("Creating merchant document",
		zap.Int("merchant_id", int(req.GetMerchantId())),
		zap.String("document_type", req.GetDocumentType()),
	)

	request := requests.CreateMerchantDocumentRequest{
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant document creation",
			zap.Int("merchant_id", int(req.GetMerchantId())),
			zap.String("document_type", req.GetDocumentType()),
			zap.Error(err),
		)
		return nil, merchantdocument_errors.ErrGrpcValidateCreateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.CreateMerchantDocument(ctx, &request)
	if err != nil {
		s.logger.Error("Failed to create merchant document",
			zap.Int("merchant_id", int(req.GetMerchantId())),
			zap.String("document_type", req.GetDocumentType()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant document created successfully",
		zap.Int("document_id", int(document.ID)),
		zap.Int("merchant_id", int(document.MerchantID)),
		zap.String("document_type", document.DocumentType),
	)

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully created merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) Update(ctx context.Context, req *pb.UpdateMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		s.logger.Error("Invalid document ID provided for update", zap.Int("document_id", id))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	s.logger.Info("Updating merchant document", zap.Int("document_id", id))

	request := requests.UpdateMerchantDocumentRequest{
		DocumentID:   &id,
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
		Status:       req.GetStatus(),
		Note:         req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant document update",
			zap.Int("document_id", id),
			zap.String("document_type", req.GetDocumentType()),
			zap.Error(err),
		)
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocument(ctx, &request)
	if err != nil {
		s.logger.Error("Failed to update merchant document",
			zap.Int("document_id", id),
			zap.Int("merchant_id", int(req.GetMerchantId())),
			zap.String("document_type", req.GetDocumentType()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant document updated successfully",
		zap.Int("document_id", id),
		zap.String("status", document.Status),
	)

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) UpdateStatus(ctx context.Context, req *pb.UpdateMerchantDocumentStatusRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		s.logger.Error("Invalid document ID for status update", zap.Int("document_id", id))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	s.logger.Info("Updating merchant document status",
		zap.Int("document_id", id),
		zap.String("status", req.GetStatus()),
	)

	request := requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &id,
		MerchantID: int(req.GetMerchantId()),
		Status:     req.GetStatus(),
		Note:       req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("Validation failed on document status update",
			zap.Int("document_id", id),
			zap.String("status", req.GetStatus()),
			zap.Error(err),
		)
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocumentStatus(ctx, &request)
	if err != nil {
		s.logger.Error("Failed to update merchant document status",
			zap.Int("document_id", id),
			zap.String("status", req.GetStatus()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant document status updated successfully",
		zap.Int("document_id", id),
		zap.String("status", document.Status),
		zap.String("note", document.Note),
	)

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document status", document), nil
}

func (s *merchantDocumentHandleGrpc) Trashed(ctx context.Context, req *pb.TrashedMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		s.logger.Error("Invalid document ID for trashing", zap.Int("document_id", id))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	s.logger.Info("Moving merchant document to trash", zap.Int("document_id", id))

	document, err := s.merchantDocumentCommand.TrashedMerchantDocument(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash merchant document",
			zap.Int("document_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant document moved to trash successfully",
		zap.Int("document_id", id),
		zap.Int("merchant_id", int(document.MerchantID)),
	)

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully trashed merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) Restore(ctx context.Context, req *pb.RestoreMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		s.logger.Error("Invalid document ID for restore", zap.Int("document_id", id))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	s.logger.Info("Restoring merchant document from trash", zap.Int("document_id", id))

	document, err := s.merchantDocumentCommand.RestoreMerchantDocument(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore merchant document",
			zap.Int("document_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant document restored successfully",
		zap.Int("document_id", id),
		zap.String("document_type", document.DocumentType),
	)

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully restored merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) DeletePermanent(ctx context.Context, req *pb.DeleteMerchantDocumentPermanentRequest) (*pb.ApiResponseMerchantDocumentDelete, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		s.logger.Error("Invalid document ID for permanent deletion", zap.Int("document_id", id))
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	s.logger.Info("Permanently deleting merchant document", zap.Int("document_id", id))

	_, err := s.merchantDocumentCommand.DeleteMerchantDocumentPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete merchant document",
			zap.Int("document_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant document permanently deleted", zap.Int("document_id", id))

	return s.mapping.ToProtoResponseMerchantDocumentDelete("success", "Successfully permanently deleted merchant document"), nil
}

func (s *merchantDocumentHandleGrpc) RestoreAll(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantDocumentAll, error) {
	s.logger.Info("Restoring all trashed merchant documents")

	_, err := s.merchantDocumentCommand.RestoreAllMerchantDocument(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all merchant documents", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant documents restored successfully")

	return s.mapping.ToProtoResponseMerchantDocumentAll("success", "Successfully restored all merchant documents"), nil
}

func (s *merchantDocumentHandleGrpc) DeleteAllPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantDocumentAll, error) {
	s.logger.Info("Permanently deleting all trashed merchant documents")

	_, err := s.merchantDocumentCommand.DeleteAllMerchantDocumentPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all merchant documents", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant documents permanently deleted")

	return s.mapping.ToProtoResponseMerchantDocumentAll("success", "Successfully permanently deleted all merchant documents"), nil
}
