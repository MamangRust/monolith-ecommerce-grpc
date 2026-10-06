package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDocumentQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error)
	FindByID(ctx context.Context, id int) (*db.GetMerchantDocumentRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error)
}

type MerchantDocumentCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error)
	Trash(ctx context.Context, merchant_document_id int) (*db.MerchantDocument, error)
	Restore(ctx context.Context, merchant_document_id int) (*db.MerchantDocument, error)
	DeletePermanent(ctx context.Context, merchant_document_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, error)

	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, error)

	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, error)

	FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type MerchantCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateMerchantRequest,
	) (*db.CreateMerchantRow, error)

	Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error)

	Trash(
		ctx context.Context,
		merchant_id int,
	) (*db.Merchant, error)

	Restore(
		ctx context.Context,
		merchant_id int,
	) (*db.Merchant, error)

	DeletePermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)

	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.UpdateMerchantStatusRow, error)
}

type UserQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetUserByIDRow, error)
}
