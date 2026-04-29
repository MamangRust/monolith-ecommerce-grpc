package merchantdetailapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type merchantDetailCommandResponseMapper struct{}

func NewMerchantDetailCommandResponseMapper() MerchantDetailCommandResponseMapper {
	return &merchantDetailCommandResponseMapper{}
}

func (m *merchantDetailCommandResponseMapper) ToResponseMerchantDetail(merchant *pb.MerchantDetailResponse) *response.MerchantDetailResponse {
	if merchant == nil { return nil }
	return &response.MerchantDetailResponse{
		ID:               int(merchant.Id),
		MerchantID:       int(merchant.MerchantId),
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: nil,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}
}

func (m *merchantDetailCommandResponseMapper) ToResponseMerchantDetailRelation(merchant *pb.MerchantDetailResponse) *response.MerchantDetailResponse {
	if merchant == nil { return nil }
	var socialMediaLinks []*response.MerchantSocialMediaLinkResponse
	for _, sm := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &response.MerchantSocialMediaLinkResponse{
			ID:       int(sm.Id),
			Platform: sm.Platform,
			Url:      sm.Url,
		})
	}

	return &response.MerchantDetailResponse{
		ID:               int(merchant.Id),
		MerchantID:       int(merchant.MerchantId),
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
	}
}

func (m *merchantDetailCommandResponseMapper) ToResponsesMerchantDetail(merchants []*pb.MerchantDetailResponse) []*response.MerchantDetailResponse {
	var mappedMerchants []*response.MerchantDetailResponse
	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantDetailRelation(merchant))
	}
	return mappedMerchants
}

func (m *merchantDetailCommandResponseMapper) ToResponseMerchantDetailDeleteAt(merchant *pb.MerchantDetailResponseDeleteAt) *response.MerchantDetailResponseDeleteAt {
	if merchant == nil { return nil }
	var deletedAt string
	if merchant.DeletedAt != nil {
		deletedAt = merchant.DeletedAt.Value
	}

	var socialMediaLinks []*response.MerchantSocialMediaLinkResponse
	for _, sm := range merchant.SocialMediaLinks {
		socialMediaLinks = append(socialMediaLinks, &response.MerchantSocialMediaLinkResponse{
			ID:       int(sm.Id),
			Platform: sm.Platform,
			Url:      sm.Url,
		})
	}

	return &response.MerchantDetailResponseDeleteAt{
		ID:               int(merchant.Id),
		MerchantID:       int(merchant.MerchantId),
		DisplayName:      merchant.DisplayName,
		CoverImageUrl:    merchant.CoverImageUrl,
		LogoUrl:          merchant.LogoUrl,
		ShortDescription: merchant.ShortDescription,
		WebsiteUrl:       merchant.WebsiteUrl,
		SocialMediaLinks: socialMediaLinks,
		CreatedAt:        merchant.CreatedAt,
		UpdatedAt:        merchant.UpdatedAt,
		DeletedAt:        &deletedAt,
	}
}

func (m *merchantDetailCommandResponseMapper) ToResponsesMerchantDetailDeleteAt(merchants []*pb.MerchantDetailResponseDeleteAt) []*response.MerchantDetailResponseDeleteAt {
	var mappedMerchants []*response.MerchantDetailResponseDeleteAt
	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantDetailDeleteAt(merchant))
	}
	return mappedMerchants
}

func (m *merchantDetailCommandResponseMapper) ToApiResponseMerchantDetail(pbResponse *pb.ApiResponseMerchantDetail) *response.ApiResponseMerchantDetail {
	return &response.ApiResponseMerchantDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDetail(pbResponse.Data),
	}
}

func (m *merchantDetailCommandResponseMapper) ToApiResponseMerchantDetailDeleteAt(pbResponse *pb.ApiResponseMerchantDetailDeleteAt) *response.ApiResponseMerchantDetailDeleteAt {
	return &response.ApiResponseMerchantDetailDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDetailDeleteAt(pbResponse.Data),
	}
}

func (m *merchantDetailCommandResponseMapper) ToApiResponsePaginationMerchantDetailDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantDetailDeleteAt) *response.ApiResponsePaginationMerchantDetailDeleteAt {
	return &response.ApiResponsePaginationMerchantDetailDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantDetailDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
