package service

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"golang.org/x/net/context"
)

type MerchantQueryService interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, *int, error)

	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, *int, error)

	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, *int, error)

	FindById(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type MerchantCommandService interface {
	CreateMerchant(
		ctx context.Context,
		request *requests.CreateMerchantRequest,
	) (*db.CreateMerchantRow, error)

	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error)

	TrashedMerchant(
		ctx context.Context,
		merchant_id int,
	) (*db.Merchant, error)

	RestoreMerchant(
		ctx context.Context,
		merchant_id int,
	) (*db.Merchant, error)

	DeleteMerchantPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.UpdateMerchantStatusRow, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error)
	FindById(ctx context.Context, documentID int) (*db.GetMerchantDocumentRow, error)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error)
	TrashedMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, documentID int) (*db.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}
