package handler

import (
	pbmerchant_policy "github.com/MamangRust/monolith-ecommerce-pb/merchant_policy"
)

type MerchantPolicyQueryHandler interface {
	pbmerchant_policy.MerchantPolicyQueryServiceServer
}

type MerchantPolicyCommandHandler interface {
	pbmerchant_policy.MerchantPolicyCommandServiceServer
}
