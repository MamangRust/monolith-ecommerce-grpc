package bannerapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type bannerCommandResponseMapper struct{}

func NewBannerCommandResponseMapper() BannerCommandResponseMapper {
	return &bannerCommandResponseMapper{}
}

func (m *bannerCommandResponseMapper) ToResponseBanner(banner *pb.BannerResponse) *response.BannerResponse {
	if banner == nil { return nil }
	return &response.BannerResponse{
		ID:        banner.BannerId,
		Name:      banner.Name,
		StartDate: banner.StartDate,
		EndDate:   banner.EndDate,
		StartTime: banner.StartTime,
		EndTime:   banner.EndTime,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
	}
}

func (m *bannerCommandResponseMapper) ToResponsesBanner(banners []*pb.BannerResponse) []*response.BannerResponse {
	var mappedBanners []*response.BannerResponse
	for _, banner := range banners {
		mappedBanners = append(mappedBanners, m.ToResponseBanner(banner))
	}
	return mappedBanners
}

func (m *bannerCommandResponseMapper) ToApiResponseBanner(pbResponse *pb.ApiResponseBanner) *response.ApiResponseBanner {
	return &response.ApiResponseBanner{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseBanner(pbResponse.Data),
	}
}

func (m *bannerCommandResponseMapper) ToResponseBannerDeleteAt(banner *pb.BannerResponseDeleteAt) *response.BannerResponseDeleteAt {
	if banner == nil { return nil }
	var deletedAt string
	if banner.DeletedAt != nil {
		deletedAt = banner.DeletedAt.Value
	}

	return &response.BannerResponseDeleteAt{
		ID:        banner.BannerId,
		Name:      banner.Name,
		StartDate: banner.StartDate,
		EndDate:   banner.EndDate,
		StartTime: banner.StartTime,
		EndTime:   banner.EndTime,
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (m *bannerCommandResponseMapper) ToResponsesBannerDeleteAt(banners []*pb.BannerResponseDeleteAt) []*response.BannerResponseDeleteAt {
	var mappedBanners []*response.BannerResponseDeleteAt
	for _, banner := range banners {
		mappedBanners = append(mappedBanners, m.ToResponseBannerDeleteAt(banner))
	}
	return mappedBanners
}

func (m *bannerCommandResponseMapper) ToApiResponseBannerDeleteAt(pbResponse *pb.ApiResponseBannerDeleteAt) *response.ApiResponseBannerDeleteAt {
	return &response.ApiResponseBannerDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseBannerDeleteAt(pbResponse.Data),
	}
}

func (m *bannerCommandResponseMapper) ToApiResponseBannerDelete(pbResponse *pb.ApiResponseBannerDelete) *response.ApiResponseBannerDelete {
	return &response.ApiResponseBannerDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *bannerCommandResponseMapper) ToApiResponseBannerAll(pbResponse *pb.ApiResponseBannerAll) *response.ApiResponseBannerAll {
	return &response.ApiResponseBannerAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}
