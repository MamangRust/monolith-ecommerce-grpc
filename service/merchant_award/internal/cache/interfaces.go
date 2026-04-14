package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantAwardQueryCache interface {
	GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsRow, *int, bool)
	SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantCertificationsAndAwardsRow, total *int)

	GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, *int, bool)
	SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantCertificationsAndAwardsActiveRow, total *int)

	GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, *int, bool)
	SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantCertificationsAndAwardsTrashedRow, total *int)

	GetCachedMerchantAward(ctx context.Context, id int) (*db.GetMerchantCertificationOrAwardRow, bool)
	SetCachedMerchantAward(ctx context.Context, data *db.GetMerchantCertificationOrAwardRow)
}

type MerchantAwardCommandCache interface {
	DeleteMerchantAwardCache(ctx context.Context, id int)
}
