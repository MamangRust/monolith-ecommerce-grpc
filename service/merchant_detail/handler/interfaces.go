package handler

import (
	pbmerchant_detail "github.com/MamangRust/monolith-ecommerce-pb/merchant_detail"
	pbmerchant_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchant_social_link"
)

type MerchantDetailQueryHandler interface {
	pbmerchant_detail.MerchantDetailQueryServiceServer
}

type MerchantDetailCommandHandler interface {
	pbmerchant_detail.MerchantDetailCommandServiceServer
}

type MerchantSocialLinkCommandHandler interface {
	pbmerchant_social_link.MerchantSocialCommandServiceServer
}
