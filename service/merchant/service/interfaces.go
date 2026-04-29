package service

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"golang.org/x/net/context"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, *int, error)
	FindByID(ctx context.Context, merchantID int) (*db.GetMerchantByIDRow, error)
}

type MerchantCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantRequest) (*db.CreateMerchantRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error)
	Trash(ctx context.Context, merchantID int) (*db.Merchant, error)
	Restore(ctx context.Context, merchantID int) (*db.Merchant, error)
	DeletePermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.UpdateMerchantStatusRow, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, error)
	FindByID(ctx context.Context, documentID int) (*db.GetMerchantDocumentRow, error)
}

type MerchantDocumentCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error)
	UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error)
	Trash(ctx context.Context, documentID int) (*db.MerchantDocument, error)
	Restore(ctx context.Context, documentID int) (*db.MerchantDocument, error)
	DeletePermanent(ctx context.Context, documentID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
