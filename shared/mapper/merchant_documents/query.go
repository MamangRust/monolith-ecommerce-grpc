package merchantdocumentsapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type merchantDocumentQueryResponseMapper struct{}

func NewMerchantDocumentQueryResponseMapper() MerchantDocumentQueryResponseMapper {
	return &merchantDocumentQueryResponseMapper{}
}

func (m *merchantDocumentQueryResponseMapper) MapMerchantDocument(doc *pb.MerchantDocument) *response.MerchantDocumentResponse {
	if doc == nil { return nil }
	return &response.MerchantDocumentResponse{
		ID:           int(doc.DocumentId),
		MerchantID:   int(doc.MerchantId),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note,
		CreatedAt:    doc.UploadedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}

func (m *merchantDocumentQueryResponseMapper) MapMerchantDocuments(docs []*pb.MerchantDocument) []*response.MerchantDocumentResponse {
	var responses []*response.MerchantDocumentResponse
	for _, doc := range docs {
		responses = append(responses, m.MapMerchantDocument(doc))
	}
	return responses
}

func (m *merchantDocumentQueryResponseMapper) ToApiResponseMerchantDocument(doc *pb.ApiResponseMerchantDocument) *response.ApiResponseMerchantDocument {
	return &response.ApiResponseMerchantDocument{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.MapMerchantDocument(doc.Data),
	}
}

func (m *merchantDocumentQueryResponseMapper) ToApiResponsesMerchantDocument(docs *pb.ApiResponsesMerchantDocument) *response.ApiResponsesMerchantDocument {
	return &response.ApiResponsesMerchantDocument{
		Status:  docs.Status,
		Message: docs.Message,
		Data:    m.MapMerchantDocuments(docs.Data),
	}
}

func (m *merchantDocumentQueryResponseMapper) ToApiResponsePaginationMerchantDocument(docs *pb.ApiResponsePaginationMerchantDocument) *response.ApiResponsePaginationMerchantDocument {
	return &response.ApiResponsePaginationMerchantDocument{
		Status:     docs.Status,
		Message:    docs.Message,
		Data:       m.MapMerchantDocuments(docs.Data),
		Pagination: paginationapimapper.MapPaginationMeta(docs.Pagination),
	}
}

func (m *merchantDocumentQueryResponseMapper) ToApiResponsePaginationMerchantDocumentDeleteAt(docs *pb.ApiResponsePaginationMerchantDocumentAt) *response.ApiResponsePaginationMerchantDocumentDeleteAt {
	var data []*response.MerchantDocumentResponseDeleteAt
	for _, doc := range docs.Data {
		var deletedAt *string
		if doc.DeletedAt != nil { deletedAt = &doc.DeletedAt.Value }
		data = append(data, &response.MerchantDocumentResponseDeleteAt{
			ID:           int(doc.DocumentId),
			MerchantID:   int(doc.MerchantId),
			DocumentType: doc.DocumentType,
			DocumentURL:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         doc.Note,
			CreatedAt:    doc.UploadedAt,
			UpdatedAt:    doc.UpdatedAt,
			DeletedAt:    deletedAt,
		})
	}
	return &response.ApiResponsePaginationMerchantDocumentDeleteAt{
		Status:     docs.Status,
		Message:    docs.Message,
		Data:       data,
		Pagination: paginationapimapper.MapPaginationMeta(docs.Pagination),
	}
}
