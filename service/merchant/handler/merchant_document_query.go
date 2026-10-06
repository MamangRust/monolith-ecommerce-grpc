package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-merchant/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"

	pbmerchant_document "github.com/MamangRust/monolith-ecommerce-pb/merchant_document"
)

type merchantDocumentQueryHandler struct {
	pbmerchant_document.UnimplementedMerchantDocumentQueryServiceServer
	merchantDocumentQuery service.MerchantDocumentQueryService
	logger                logger.LoggerInterface
}

func NewMerchantDocumentQueryHandler(svc service.MerchantDocumentQueryService, logger logger.LoggerInterface) pbmerchant_document.MerchantDocumentQueryServiceServer {
	return &merchantDocumentQueryHandler{
		merchantDocumentQuery: svc,
		logger:                logger,
	}
}

func (s *merchantDocumentQueryHandler) FindAll(ctx context.Context, req *pbmerchant_document.FindAllMerchantDocumentsRequest) (*pbmerchant_document.ApiResponsePaginationMerchantDocument, error) {
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

	pbDocuments := make([]*pbmerchant_document.MerchantDocument, len(documents))
	for i, d := range documents {
		pbDocuments[i] = mapToProtoMerchantDocumentResponse(d)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchant_document.ApiResponsePaginationMerchantDocument{
		Status:     "success",
		Message:    "Successfully fetched merchant documents",
		Data:       pbDocuments,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDocumentQueryHandler) FindById(ctx context.Context, req *pbmerchant_document.FindMerchantDocumentByIdRequest) (*pbmerchant_document.ApiResponseMerchantDocument, error) {
	id := int(req.GetDocumentId())
	if id == 0 {
		return nil, merchant_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully fetched merchant document",
		Data:    mapToProtoMerchantDocumentResponse(document),
	}, nil
}

func (s *merchantDocumentQueryHandler) FindByActive(ctx context.Context, req *pbmerchant_document.FindAllMerchantDocumentsRequest) (*pbmerchant_document.ApiResponsePaginationMerchantDocument, error) {
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

	pbDocuments := make([]*pbmerchant_document.MerchantDocument, len(documents))
	for i, d := range documents {
		pbDocuments[i] = mapToProtoMerchantDocumentResponse(d)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchant_document.ApiResponsePaginationMerchantDocument{
		Status:     "success",
		Message:    "Successfully fetched active merchant documents",
		Data:       pbDocuments,
		Pagination: paginationMeta,
	}, nil
}

func (s *merchantDocumentQueryHandler) FindByTrashed(ctx context.Context, req *pbmerchant_document.FindAllMerchantDocumentsRequest) (*pbmerchant_document.ApiResponsePaginationMerchantDocumentAt, error) {
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

	pbDocuments := make([]*pbmerchant_document.MerchantDocumentDeleteAt, len(documents))
	for i, d := range documents {
		pbDocuments[i] = mapToProtoMerchantDocumentResponseAt(d)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pbmerchant_document.ApiResponsePaginationMerchantDocumentAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant documents",
		Data:       pbDocuments,
		Pagination: paginationMeta,
	}, nil
}
