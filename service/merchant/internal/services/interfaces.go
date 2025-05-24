package services

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type MerchantQueryService interface {
	FindAll(req *requests.FindAllMerchant) ([]*response.MerchantResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchant) ([]*response.MerchantResponseDeleteAt, *int, *response.ErrorResponse)
	FindById(merchantID int) (*response.MerchantResponse, *response.ErrorResponse)
}

type MerchantCommandService interface {
	CreateMerchant(req *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchant(req *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse)
	UpdateMerchantStatus(req *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse)
	TrashedMerchant(merchantID int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse)
	RestoreMerchant(merchantID int) (*response.MerchantResponse, *response.ErrorResponse)
	DeleteMerchantPermanent(merchantID int) (bool, *response.ErrorResponse)
	RestoreAllMerchant() (bool, *response.ErrorResponse)
	DeleteAllMerchantPermanent() (bool, *response.ErrorResponse)
}

type MerchantDocumentQueryService interface {
	FindAll(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, *response.ErrorResponse)

	FindById(document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(request *requests.CreateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse)

	TrashedMerchantDocument(document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	RestoreMerchantDocument(document_id int) (*response.MerchantDocumentResponse, *response.ErrorResponse)
	DeleteMerchantDocumentPermanent(document_id int) (bool, *response.ErrorResponse)

	RestoreAllMerchantDocument() (bool, *response.ErrorResponse)
	DeleteAllMerchantDocumentPermanent() (bool, *response.ErrorResponse)
}
