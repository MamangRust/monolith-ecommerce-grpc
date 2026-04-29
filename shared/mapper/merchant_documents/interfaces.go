package merchantdocumentsapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantDocumentBaseResponseMapper interface {
	MapMerchantDocument(doc *pb.MerchantDocument) *response.MerchantDocumentResponse
	MapMerchantDocuments(docs []*pb.MerchantDocument) []*response.MerchantDocumentResponse
	ToApiResponseMerchantDocument(doc *pb.ApiResponseMerchantDocument) *response.ApiResponseMerchantDocument
}

type MerchantDocumentQueryResponseMapper interface {
	MerchantDocumentBaseResponseMapper
	ToApiResponsesMerchantDocument(docs *pb.ApiResponsesMerchantDocument) *response.ApiResponsesMerchantDocument
	ToApiResponsePaginationMerchantDocument(docs *pb.ApiResponsePaginationMerchantDocument) *response.ApiResponsePaginationMerchantDocument
	ToApiResponsePaginationMerchantDocumentDeleteAt(docs *pb.ApiResponsePaginationMerchantDocumentAt) *response.ApiResponsePaginationMerchantDocumentDeleteAt
}

type MerchantDocumentCommandResponseMapper interface {
	MerchantDocumentBaseResponseMapper
	MapMerchantDocumentDeletedAt(doc *pb.MerchantDocumentDeleteAt) *response.MerchantDocumentResponseDeleteAt
	MapMerchantDocumentsDeletedAt(docs []*pb.MerchantDocumentDeleteAt) []*response.MerchantDocumentResponseDeleteAt
	ToApiResponseMerchantDocumentAll(resp *pb.ApiResponseMerchantDocumentAll) *response.ApiResponseMerchantDocumentAll
	ToApiResponseMerchantDocumentDeleteAt(resp *pb.ApiResponseMerchantDocumentDelete) *response.ApiResponseMerchantDocumentDelete
}
