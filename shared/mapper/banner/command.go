package bannerapimapper

import (
	pbbanner "github.com/MamangRust/monolith-ecommerce-pb/banner"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type bannerCommandResponseMapper struct{}

func NewBannerCommandResponseMapper() BannerCommandResponseMapper {
	return &bannerCommandResponseMapper{}
}

func (m *bannerCommandResponseMapper) ToResponseBanner(banner *pbbanner.BannerResponse) *response.BannerResponse {
	if banner == nil {
		return nil
	}
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

func (m *bannerCommandResponseMapper) ToResponsesBanner(banners []*pbbanner.BannerResponse) []*response.BannerResponse {
	var mappedBanners []*response.BannerResponse
	for _, banner := range banners {
		mappedBanners = append(mappedBanners, m.ToResponseBanner(banner))
	}
	return mappedBanners
}

func (m *bannerCommandResponseMapper) ToApiResponseBanner(pbResponse *pbbanner.ApiResponseBanner) *response.ApiResponseBanner {
	return &response.ApiResponseBanner{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseBanner(pbResponse.Data),
	}
}

func (m *bannerCommandResponseMapper) ToResponseBannerDeleteAt(banner *pbbanner.BannerResponseDeleteAt) *response.BannerResponseDeleteAt {
	if banner == nil {
		return nil
	}
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

func (m *bannerCommandResponseMapper) ToResponsesBannerDeleteAt(banners []*pbbanner.BannerResponseDeleteAt) []*response.BannerResponseDeleteAt {
	var mappedBanners []*response.BannerResponseDeleteAt
	for _, banner := range banners {
		mappedBanners = append(mappedBanners, m.ToResponseBannerDeleteAt(banner))
	}
	return mappedBanners
}

func (m *bannerCommandResponseMapper) ToApiResponseBannerDeleteAt(pbResponse *pbbanner.ApiResponseBannerDeleteAt) *response.ApiResponseBannerDeleteAt {
	return &response.ApiResponseBannerDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseBannerDeleteAt(pbResponse.Data),
	}
}

func (m *bannerCommandResponseMapper) ToApiResponseBannerDelete(pbResponse *pbbanner.ApiResponseBannerDelete) *response.ApiResponseBannerDelete {
	return &response.ApiResponseBannerDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (m *bannerCommandResponseMapper) ToApiResponseBannerAll(pbResponse *pbbanner.ApiResponseBannerAll) *response.ApiResponseBannerAll {
	return &response.ApiResponseBannerAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}
