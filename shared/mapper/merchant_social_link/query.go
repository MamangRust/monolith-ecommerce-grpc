package merchantsociallinkapimapper

import (
	pbmerchant_detail "github.com/MamangRust/monolith-ecommerce-pb/merchant_detail"
	pbmerchant_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchant_social_link"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantSocialLinkQueryResponseMapper struct{}

func NewMerchantSocialLinkQueryResponseMapper() MerchantSocialLinkQueryResponseMapper {
	return &merchantSocialLinkQueryResponseMapper{}
}

func (m *merchantSocialLinkQueryResponseMapper) MapMerchantSocialLink(doc *pbmerchant_detail.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse {
	return &response.MerchantSocialLinkResponse{
		ID:               int(doc.Id),
		MerchantDetailID: int(doc.MerchantDetailId),
		Platform:         doc.Platform,
		URL:              doc.Url,
		CreatedAt:        doc.CreatedAt,
		UpdatedAt:        doc.UpdatedAt,
	}
}

func (m *merchantSocialLinkQueryResponseMapper) ToApiResponseMerchantSocialLink(doc *pbmerchant_social_link.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink {
	return &response.ApiResponseMerchantSocialLink{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.MapMerchantSocialLink(doc.Data),
	}
}
