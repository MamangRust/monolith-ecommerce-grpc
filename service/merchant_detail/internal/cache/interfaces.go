package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantDetailQueryCache interface {
	GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, *int, bool)
	SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantDetailsRow, total *int)

	GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, *int, bool)
	SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantDetailsActiveRow, total *int)

	GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, *int, bool)
	SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantDetailsTrashedRow, total *int)

	GetCachedMerchantDetail(ctx context.Context, id int) (*db.GetMerchantDetailRow, bool)
	SetCachedMerchantDetail(ctx context.Context, data *db.GetMerchantDetailRow)
}

type MerchantDetailCommandCache interface {
	DeleteMerchantDetailCache(ctx context.Context, id int)
}
