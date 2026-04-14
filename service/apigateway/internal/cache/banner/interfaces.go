package banner_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerQueryCache interface {
	GetCachedBanners(
		ctx context.Context,
		req *requests.FindAllBanner,
	) (*response.ApiResponsePaginationBanner, bool)
	SetCachedBanners(
		ctx context.Context,
		req *requests.FindAllBanner,
		data *response.ApiResponsePaginationBanner,
	)
	GetCachedActiveBanners(
		ctx context.Context,
		req *requests.FindAllBanner,
	) (*response.ApiResponsePaginationBannerDeleteAt, bool)
	SetCachedActiveBanners(
		ctx context.Context,
		req *requests.FindAllBanner,
		data *response.ApiResponsePaginationBannerDeleteAt,
	)
	GetCachedTrashedBanners(
		ctx context.Context,
		req *requests.FindAllBanner,
	) (*response.ApiResponsePaginationBannerDeleteAt, bool)
	SetCachedTrashedBanners(
		ctx context.Context,
		req *requests.FindAllBanner,
		data *response.ApiResponsePaginationBannerDeleteAt,
	)
	GetCachedBanner(
		ctx context.Context,
		id int,
	) (*response.ApiResponseBanner, bool)
	SetCachedBanner(
		ctx context.Context,
		data *response.ApiResponseBanner,
	)
}

type BannerCommandCache interface {
	DeleteBannerCache(ctx context.Context, id int)
}
