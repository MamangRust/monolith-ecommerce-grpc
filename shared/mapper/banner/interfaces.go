package bannerapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type BannerBaseResponseMapper interface {
	ToResponseBanner(banner *pb.BannerResponse) *response.BannerResponse
	ToResponsesBanner(banners []*pb.BannerResponse) []*response.BannerResponse
}

type BannerQueryResponseMapper interface {
	BannerBaseResponseMapper
	ToApiResponseBanner(pbResponse *pb.ApiResponseBanner) *response.ApiResponseBanner
	ToApiResponsesBanner(pbResponse *pb.ApiResponsesBanner) *response.ApiResponsesBanner
	ToApiResponsePaginationBanner(pbResponse *pb.ApiResponsePaginationBanner) *response.ApiResponsePaginationBanner
	ToApiResponsePaginationBannerDeleteAt(pbResponse *pb.ApiResponsePaginationBannerDeleteAt) *response.ApiResponsePaginationBannerDeleteAt
}

type BannerCommandResponseMapper interface {
	BannerBaseResponseMapper
	ToApiResponseBanner(pbResponse *pb.ApiResponseBanner) *response.ApiResponseBanner
	ToResponseBannerDeleteAt(banner *pb.BannerResponseDeleteAt) *response.BannerResponseDeleteAt
	ToResponsesBannerDeleteAt(banners []*pb.BannerResponseDeleteAt) []*response.BannerResponseDeleteAt
	ToApiResponseBannerDeleteAt(pbResponse *pb.ApiResponseBannerDeleteAt) *response.ApiResponseBannerDeleteAt
	ToApiResponseBannerDelete(pbResponse *pb.ApiResponseBannerDelete) *response.ApiResponseBannerDelete
	ToApiResponseBannerAll(pbResponse *pb.ApiResponseBannerAll) *response.ApiResponseBannerAll
}
