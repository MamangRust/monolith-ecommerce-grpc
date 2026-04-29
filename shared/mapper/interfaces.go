package response_api

import (
	auth "github.com/MamangRust/monolith-ecommerce-shared/mapper/auth"
	banner "github.com/MamangRust/monolith-ecommerce-shared/mapper/banner"
	cart "github.com/MamangRust/monolith-ecommerce-shared/mapper/cart"
	category "github.com/MamangRust/monolith-ecommerce-shared/mapper/category"
	merchant "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	merchantaward "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_award"
	merchantbusiness "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_business"
	merchantdetail "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_detail"
	merchantdocuments "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_documents"
	merchantpolicy "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_policy"
	merchantsociallink "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_social_link"
	order "github.com/MamangRust/monolith-ecommerce-shared/mapper/order"
	orderitem "github.com/MamangRust/monolith-ecommerce-shared/mapper/order_item"
	product "github.com/MamangRust/monolith-ecommerce-shared/mapper/product"
	review "github.com/MamangRust/monolith-ecommerce-shared/mapper/review"
	reviewdetail "github.com/MamangRust/monolith-ecommerce-shared/mapper/review_detail"
	role "github.com/MamangRust/monolith-ecommerce-shared/mapper/role"
	shippingaddress "github.com/MamangRust/monolith-ecommerce-shared/mapper/shipping_address"
	slider "github.com/MamangRust/monolith-ecommerce-shared/mapper/slider"
	transaction "github.com/MamangRust/monolith-ecommerce-shared/mapper/transaction"
	user "github.com/MamangRust/monolith-ecommerce-shared/mapper/user"
)

type RoleResponseMapper = role.RoleResponseMapper
type ShippingAddressResponseMapper = shippingaddress.ShippingAddressResponseMapper
type SliderResponseMapper = slider.SliderResponseMapper
type BannerResponseMapper = banner.BannerResponseMapper
type AuthResponseMapper = auth.AuthResponseMapper
type UserResponseMapper = user.UserResponseMapper
type CategoryResponseMapper = category.CategoryResponseMapper
type OrderResponseMapper = order.OrderResponseMapper
type OrderItemResponseMapper = orderitem.OrderItemResponseMapper
type ProductResponseMapper = product.ProductResponseMapper
type ReviewResponseMapper = review.ReviewResponseMapper
type ReviewDetailResponseMapper = reviewdetail.ReviewDetailResponseMapper
type TransactionResponseMapper = transaction.TransactionResponseMapper
type CartResponseMapper = cart.CartResponseMapper
type MerchantResponseMapper = merchant.MerchantResponseMapper
type MerchantAwardResponseMapper = merchantaward.MerchantAwardResponseMapper
type MerchantBusinessResponseMapper = merchantbusiness.MerchantBusinessResponseMapper
type MerchantDetailResponseMapper = merchantdetail.MerchantDetailResponseMapper
type MerchantDocumentResponseMapper = merchantdocuments.MerchantDocumentResponseMapper
type MerchantPolicyResponseMapper = merchantpolicy.MerchantPolicyResponseMapper
type MerchantSocialLinkMapper = merchantsociallink.MerchantSocialLinkResponseMapper
