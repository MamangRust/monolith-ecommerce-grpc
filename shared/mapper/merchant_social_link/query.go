package merchantsociallinkapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantSocialLinkQueryResponseMapper struct{}

func NewMerchantSocialLinkQueryResponseMapper() MerchantSocialLinkQueryResponseMapper {
	return &merchantSocialLinkQueryResponseMapper{}
}

func (m *merchantSocialLinkQueryResponseMapper) MapMerchantSocialLink(doc *pb.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse {
	return &response.MerchantSocialLinkResponse{
		ID:               int(doc.Id),
		MerchantDetailID: int(doc.MerchantDetailId),
		Platform:         doc.Platform,
		URL:              doc.Url,
		CreatedAt:        doc.CreatedAt,
		UpdatedAt:        doc.UpdatedAt,
	}
}

func (m *merchantSocialLinkQueryResponseMapper) ToApiResponseMerchantSocialLink(doc *pb.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink {
	return &response.ApiResponseMerchantSocialLink{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.MapMerchantSocialLink(doc.Data),
	}
}
