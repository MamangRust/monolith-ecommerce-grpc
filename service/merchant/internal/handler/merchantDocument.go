package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/services"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDocumentHandleGrpc struct {
	pb.UnimplementedMerchantDocumentServiceServer
	merchantDocumentQuery   services.MerchantDocumentQueryService
	merchantDocumentCommand services.MerchantDocumentCommandService

	mapping protomapper.MerchantDocumentProtoMapper
}

func NewMerchantDocumentHandleGrpc(
	service *services.Service,
	mapping protomapper.MerchantDocumentProtoMapper,
) pb.MerchantDocumentServiceServer {
	return &merchantDocumentHandleGrpc{
		merchantDocumentQuery:   service.MerchantDocumentQuery,
		merchantDocumentCommand: service.MerchantDocumentCommand,
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

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return s.mapping.ToProtoResponsePaginationMerchantDocument(paginationMeta, "success", "Successfully fetched merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) FindById(ctx context.Context, req *pb.FindMerchantDocumentByIdRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentQuery.FindById(ctx, id)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

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

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByActive(ctx, &reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

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

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return s.mapping.ToProtoResponsePaginationMerchantDocumentDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant documents", documents), nil
}

func (s *merchantDocumentHandleGrpc) Create(ctx context.Context, req *pb.CreateMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	request := requests.CreateMerchantDocumentRequest{
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchantdocument_errors.ErrGrpcValidateCreateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.CreateMerchantDocument(ctx, &request)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully created merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) Update(ctx context.Context, req *pb.UpdateMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentRequest{
		DocumentID:   &id,
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
		Status:       req.GetStatus(),
		Note:         req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocument(ctx, &request)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) UpdateStatus(ctx context.Context, req *pb.UpdateMerchantDocumentStatusRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &id,
		MerchantID: int(req.GetMerchantId()),
		Status:     req.GetStatus(),
		Note:       req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocumentStatus(ctx, &request)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully updated merchant document status", document), nil
}

func (s *merchantDocumentHandleGrpc) Trashed(ctx context.Context, req *pb.TrashedMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.TrashedMerchantDocument(ctx, id)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully trashed merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) Restore(ctx context.Context, req *pb.RestoreMerchantDocumentRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.RestoreMerchantDocument(ctx, id)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocument("success", "Successfully restored merchant document", document), nil
}

func (s *merchantDocumentHandleGrpc) DeletePermanent(ctx context.Context, req *pb.DeleteMerchantDocumentPermanentRequest) (*pb.ApiResponseMerchantDocumentDelete, error) {
	id := int(req.GetDocumentId())

	if id == 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantDocumentCommand.DeleteMerchantDocumentPermanent(ctx, id)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocumentDelete("success", "Successfully permanently deleted merchant document"), nil
}

func (s *merchantDocumentHandleGrpc) RestoreAll(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantDocumentAll, error) {
	_, err := s.merchantDocumentCommand.RestoreAllMerchantDocument(ctx)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocumentAll("success", "Successfully restored all merchant documents"), nil
}

func (s *merchantDocumentHandleGrpc) DeleteAllPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantDocumentAll, error) {
	_, err := s.merchantDocumentCommand.DeleteAllMerchantDocumentPermanent(ctx)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMerchantDocumentAll("success", "Successfully permanently deleted all merchant documents"), nil
}
