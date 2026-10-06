package merchantsociallinkapimapper

import (
	pbmerchant_detail "github.com/MamangRust/monolith-ecommerce-pb/merchant_detail"
	pbmerchant_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchant_social_link"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantSocialLinkBaseResponseMapper interface {
	MapMerchantSocialLink(doc *pbmerchant_detail.MerchantSocialMediaLinkResponse) *response.MerchantSocialLinkResponse
	ToApiResponseMerchantSocialLink(doc *pbmerchant_social_link.ApiResponseMerchantSocial) *response.ApiResponseMerchantSocialLink
}

type MerchantSocialLinkQueryResponseMapper interface {
	MerchantSocialLinkBaseResponseMapper
}

type MerchantSocialLinkCommandResponseMapper interface {
	MerchantSocialLinkBaseResponseMapper
}
