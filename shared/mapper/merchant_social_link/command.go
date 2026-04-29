package merchantsociallinkapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantSocialLinkCommandResponseMapper struct{}

func NewMerchantSocialLinkCommandResponseMapper() MerchantSocialLinkCommandResponseMapper {
	return &merchantSocialLinkCommandResponseMapper{}
}

func (m *merchantSocialLinkCommandResponseMapper) MapMerchantSocialLink(doc *pb.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse {
	if doc == nil { return nil }
	return &response.MerchantSocialLinkResponse{
		ID:               int(doc.Id),
		MerchantDetailID: int(doc.MerchantDetailId),
		Platform:         doc.Platform,
		URL:              doc.Url,
	}
}

func (m *merchantSocialLinkCommandResponseMapper) ToApiResponseMerchantSocialLink(doc *pb.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink {
	return &response.ApiResponseMerchantSocialLink{
		Status:  doc.Status,
		Message: doc.Message,
		Data:    m.MapMerchantSocialLink(doc.Data),
	}
}
