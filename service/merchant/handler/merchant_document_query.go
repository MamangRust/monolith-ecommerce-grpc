package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type merchantDocumentQueryHandler struct {
	pb.UnimplementedMerchantDocumentQueryServiceServer
	merchantDocumentQuery service.MerchantDocumentQueryService
	logger                logger.LoggerInterface
}

func NewMerchantDocumentQueryHandler(svc service.MerchantDocumentQueryService, logger logger.LoggerInterface) pb.MerchantDocumentQueryServiceServer {
	return &merchantDocumentQueryHandler{
		merchantDocumentQuery: svc,
		logger:                logger,
	}
}

func (s *merchantDocumentQueryHandler) FindAll(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocument, error) {
	page, pageSize := normalizePage(int(req.GetPage()), int(req.GetPageSize()))
	search := req.GetSearch()

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbDocuments := make([]*pb.MerchantDocument, len(documents))
	for i, d := range documents {
		pbDocuments[i] = mapToProtoMerchantDocumentResponse(d)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantDocument{
		Status:     "success",
		Message:    "Successfully fetched merchant documents",
		Data:       pbDocuments,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDocumentQueryHandler) FindById(ctx context.Context, req *pb.FindMerchantDocumentByIdRequest) (*pb.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully fetched merchant document",
		Data:    mapToProtoMerchantDocumentResponse(document),
	}, nil
}

func (s *merchantDocumentQueryHandler) FindByActive(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocument, error) {
	page, pageSize := normalizePage(int(req.GetPage()), int(req.GetPageSize()))
	search := req.GetSearch()

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbDocuments := make([]*pb.MerchantDocument, len(documents))
	for i, d := range documents {
		pbDocuments[i] = mapToProtoMerchantDocumentResponse(d)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantDocument{
		Status:     "success",
		Message:    "Successfully fetched active merchant documents",
		Data:       pbDocuments,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDocumentQueryHandler) FindByTrashed(ctx context.Context, req *pb.FindAllMerchantDocumentsRequest) (*pb.ApiResponsePaginationMerchantDocumentAt, error) {
	page, pageSize := normalizePage(int(req.GetPage()), int(req.GetPageSize()))
	search := req.GetSearch()

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbDocuments := make([]*pb.MerchantDocumentDeleteAt, len(documents))
	for i, d := range documents {
		pbDocuments[i] = mapToProtoMerchantDocumentResponseAt(d)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationMerchantDocumentAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant documents",
		Data:       pbDocuments,
		Pagination: paginationMeta,
	}, nil
}
