package merchantdocumentsapimapper

import (
	pbmerchant_document "github.com/MamangRust/monolith-ecommerce-pb/merchant_document"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantDocumentCommandResponseMapper struct{}

func NewMerchantDocumentCommandResponseMapper() MerchantDocumentCommandResponseMapper {
	return &merchantDocumentCommandResponseMapper{}
}

func (m *merchantDocumentCommandResponseMapper) MapMerchantDocument(doc *pbmerchant_document.MerchantDocument) *response.MerchantDocumentResponse {
	if doc == nil {
		return nil
	}
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

func (m *merchantDocumentCommandResponseMapper) MapMerchantDocuments(docs []*pbmerchant_document.MerchantDocument) []*response.MerchantDocumentResponse {
	var responses []*response.MerchantDocumentResponse
	for _, doc := range docs {
		responses = append(responses, m.MapMerchantDocument(doc))
	}
	return responses
}

func (m *merchantDocumentCommandResponseMapper) ToApiResponseMerchantDocument(doc *pbmerchant_document.ApiResponseMerchantDocument) *response.ApiResponseMerchantDocument {
	return &response.ApiResponseMerchantDocument{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.MapMerchantDocument(doc.Data),
	}
}

func (m *merchantDocumentCommandResponseMapper) MapMerchantDocumentDeletedAt(doc *pbmerchant_document.MerchantDocumentDeleteAt) *response.MerchantDocumentResponseDeleteAt {
	if doc == nil {
		return nil
	}
	var deletedAt *string
	if doc.DeletedAt != nil {
		deletedAt = &doc.DeletedAt.Value
	}

	return &response.MerchantDocumentResponseDeleteAt{
		ID:           int(doc.DocumentId),
		MerchantID:   int(doc.MerchantId),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note,
		CreatedAt:    doc.UploadedAt,
		UpdatedAt:    doc.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentCommandResponseMapper) MapMerchantDocumentsDeletedAt(docs []*pbmerchant_document.MerchantDocumentDeleteAt) []*response.MerchantDocumentResponseDeleteAt {
	var responses []*response.MerchantDocumentResponseDeleteAt
	for _, doc := range docs {
		responses = append(responses, m.MapMerchantDocumentDeletedAt(doc))
	}
	return responses
}

func (m *merchantDocumentCommandResponseMapper) ToApiResponseMerchantDocumentAll(resp *pbmerchant_document.ApiResponseMerchantDocumentAll) *response.ApiResponseMerchantDocumentAll {
	return &response.ApiResponseMerchantDocumentAll{
		Status:  resp.Status,
		Message: resp.Message,
	}
}

func (m *merchantDocumentCommandResponseMapper) ToApiResponseMerchantDocumentDeleteAt(resp *pbmerchant_document.ApiResponseMerchantDocumentDelete) *response.ApiResponseMerchantDocumentDelete {
	return &response.ApiResponseMerchantDocumentDelete{
		Status:  resp.Status,
		Message: resp.Message,
	}
}
