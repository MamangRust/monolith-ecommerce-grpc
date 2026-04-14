package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantBusinessQueryCache interface {
	GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationRow, *int, bool)
	SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsBusinessInformationRow, total *int)

	GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationActiveRow, *int, bool)
	SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsBusinessInformationActiveRow, total *int)

	GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationTrashedRow, *int, bool)
	SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantsBusinessInformationTrashedRow, total *int)

	GetCachedMerchantBusiness(ctx context.Context, id int) (*db.GetMerchantBusinessInformationRow, bool)
	SetCachedMerchantBusiness(ctx context.Context, data *db.GetMerchantBusinessInformationRow)
}

type MerchantBusinessCommandCache interface {
	DeleteMerchantBusinessCache(ctx context.Context, id int)
}
