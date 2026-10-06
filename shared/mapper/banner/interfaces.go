package bannerapimapper

import (
	pbbanner "github.com/MamangRust/monolith-ecommerce-pb/banner"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerBaseResponseMapper interface {
	ToResponseBanner(banner *pbbanner.BannerResponse) *response.BannerResponse
	ToResponsesBanner(banners []*pbbanner.BannerResponse) []*response.BannerResponse
}

type BannerQueryResponseMapper interface {
	BannerBaseResponseMapper
	ToApiResponseBanner(pbResponse *pbbanner.ApiResponseBanner) *response.ApiResponseBanner
	ToApiResponsesBanner(pbResponse *pbbanner.ApiResponsesBanner) *response.ApiResponsesBanner
	ToApiResponsePaginationBanner(pbResponse *pbbanner.ApiResponsePaginationBanner) *response.ApiResponsePaginationBanner
	ToApiResponsePaginationBannerDeleteAt(pbResponse *pbbanner.ApiResponsePaginationBannerDeleteAt) *response.ApiResponsePaginationBannerDeleteAt
}

type BannerCommandResponseMapper interface {
	BannerBaseResponseMapper
	ToApiResponseBanner(pbResponse *pbbanner.ApiResponseBanner) *response.ApiResponseBanner
	ToResponseBannerDeleteAt(banner *pbbanner.BannerResponseDeleteAt) *response.BannerResponseDeleteAt
	ToResponsesBannerDeleteAt(banners []*pbbanner.BannerResponseDeleteAt) []*response.BannerResponseDeleteAt
	ToApiResponseBannerDeleteAt(pbResponse *pbbanner.ApiResponseBannerDeleteAt) *response.ApiResponseBannerDeleteAt
	ToApiResponseBannerDelete(pbResponse *pbbanner.ApiResponseBannerDelete) *response.ApiResponseBannerDelete
	ToApiResponseBannerAll(pbResponse *pbbanner.ApiResponseBannerAll) *response.ApiResponseBannerAll
}
