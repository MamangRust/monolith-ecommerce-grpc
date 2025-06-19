package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerQueryCache interface {
	GetCachedBannersCache(req *requests.FindAllBanner) ([]*response.BannerResponse, *int, bool)
	SetCachedBannersCache(req *requests.FindAllBanner, data []*response.BannerResponse, total *int)

	GetCachedBannerActiveCache(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool)
	SetCachedBannerActiveCache(req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int)

	GetCachedBannerTrashedCache(req *requests.FindAllBanner) ([]*response.BannerResponseDeleteAt, *int, bool)
	SetCachedBannerTrashedCache(req *requests.FindAllBanner, data []*response.BannerResponseDeleteAt, total *int)

	GetCachedBannerCache(id int) (*response.BannerResponse, bool)
	SetCachedBannerCache(data *response.BannerResponse)
}

type BannerCommandCache interface {
	DeleteBannerCache(id int)
}
