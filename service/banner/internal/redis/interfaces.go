package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerQueryCache interface {
	GetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponse, *int, bool)
	SetCachedBannersCache(ctx context.Context, req *requests.FindAllBanner, data []*response.BannerResponse, total *int)

	GetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool)
	SetCachedBannerActiveCache(ctx context.Context, req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int)

	GetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool)
	SetCachedBannerTrashedCache(ctx context.Context, req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int)

	GetCachedBannerCache(ctx context.Context, id int) (*response.BannerResponse, bool)
	SetCachedBannerCache(ctx context.Context, data *response.BannerResponse)
}

type BannerCommandCache interface {
	DeleteBannerCache(ctx context.Context, id int)
}
