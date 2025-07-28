package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDocumentQueryRepository interface {
	FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
	FindById(ctx context.Context, id int) (*record.MerchantDocumentRecord, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*record.MerchantDocumentRecord, *int, error)
}

type MerchantDocumentCommandRepository interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*record.MerchantDocumentRecord, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*record.MerchantDocumentRecord, error)
	TrashedMerchantDocument(ctx context.Context, merchant_document_id int) (*record.MerchantDocumentRecord, error)
	RestoreMerchantDocument(ctx context.Context, merchant_document_id int) (*record.MerchantDocumentRecord, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, merchant_document_id int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.MerchantRecord, error)
}

type MerchantCommandRepository interface {
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error)
	TrashedMerchant(ctx context.Context, merchant_id int) (*record.MerchantRecord, error)
	RestoreMerchant(ctx context.Context, merchant_id int) (*record.MerchantRecord, error)
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type UserQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*record.UserRecord, error)
}
