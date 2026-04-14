package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type BannerQueryCache interface {
	GetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, *int, bool)
	SetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner, data []*db.GetBannersRow, total *int)

	GetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, *int, bool)
	SetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner, data []*db.GetBannersActiveRow, total *int)

	GetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, *int, bool)
	SetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner, data []*db.GetBannersTrashedRow, total *int)

	GetCachedBannerCache(ctx context.Context, id int) (*db.GetBannerRow, bool)
	SetCachedBannerCache(ctx context.Context, data *db.GetBannerRow)
}

type BannerCommandCache interface {
	DeleteBannerCache(ctx context.Context, id int)
}
