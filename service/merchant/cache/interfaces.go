package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, *int, bool)
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsRow, total *int)

	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, *int, bool)
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsActiveRow, total *int)

	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, *int, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsTrashedRow, total *int)

	GetCachedMerchant(ctx context.Context, id int) (*db.GetMerchantByIDRow, bool)
	SetCachedMerchant(ctx context.Context, data *db.GetMerchantByIDRow)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetMerchantDocumentsRow, *int, bool)
	SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetMerchantDocumentsRow, total *int)

	GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetActiveMerchantDocumentsRow, *int, bool)
	SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetActiveMerchantDocumentsRow, total *int)

	GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*db.GetTrashedMerchantDocumentsRow, *int, bool)
	SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*db.GetTrashedMerchantDocumentsRow, total *int)

	GetCachedMerchantDocument(ctx context.Context, id int) (*db.GetMerchantDocumentRow, bool)
	SetCachedMerchantDocument(ctx context.Context, data *db.GetMerchantDocumentRow)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(ctx context.Context, id int)
}
