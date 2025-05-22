package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDocumentQueryRepository interface {
	FindAllDocuments(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
	FindById(id int) (*record.MerchantDocumentRecord, error)

	FindByActive(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
}
type MerchantDocumentCommandRepository interface {
	CreateMerchantDocument(request *requests.CreateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)
	UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)
	UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*record.MerchantDocumentRecord, error)
	TrashedMerchantDocument(merchant_document_id int) (*record.MerchantDocumentRecord, error)
	RestoreMerchantDocument(merchant_document_id int) (*record.MerchantDocumentRecord, error)
	DeleteMerchantDocumentPermanent(merchant_document_id int) (bool, error)
	RestoreAllMerchantDocument() (bool, error)
	DeleteAllMerchantDocumentPermanent() (bool, error)
}

type MerchantQueryRepository interface {
	FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindById(user_id int) (*record.MerchantRecord, error)
}

type MerchantCommandRepository interface {
	CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error)
	TrashedMerchant(merchant_id int) (*record.MerchantRecord, error)
	RestoreMerchant(merchant_id int) (*record.MerchantRecord, error)
	DeleteMerchantPermanent(Merchant_id int) (bool, error)
	RestoreAllMerchant() (bool, error)
	DeleteAllMerchantPermanent() (bool, error)
}

type UserQueryRepository interface {
	FindById(user_id int) (*record.UserRecord, error)
}
