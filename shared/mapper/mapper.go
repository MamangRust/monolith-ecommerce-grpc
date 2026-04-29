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

type ResponseApiMapper struct {
	AuthResponseMapper            AuthResponseMapper
	RoleResponseMapper            RoleResponseMapper
	UserResponseMapper            UserResponseMapper
	CategoryResponseMapper        CategoryResponseMapper
	MerchantResponseMapper        MerchantResponseMapper
	OrderItemResponseMapper       OrderItemResponseMapper
	OrderResponseMapper           OrderResponseMapper
	ProductResponseMapper         ProductResponseMapper
	TransactionResponseMapper     TransactionResponseMapper
	CartResponseMapper            CartResponseMapper
	ReviewMapper                  ReviewResponseMapper
	SliderMapper                  SliderResponseMapper
	ShippingAddressResponseMapper ShippingAddressResponseMapper
	BannerResponseMapper          BannerResponseMapper
	MerchantAwardResponseMapper   MerchantAwardResponseMapper
	MerchantBusinessMapper        MerchantBusinessResponseMapper
	MerchantDetailResponseMapper  MerchantDetailResponseMapper
	MerchantPolicyResponseMapper  MerchantPolicyResponseMapper
	ReviewDetailResponseMapper    ReviewDetailResponseMapper
	MerchantDocumentProMapper     MerchantDocumentResponseMapper
	MerchantSocialLinkProtoMapper MerchantSocialLinkMapper
}

func NewResponseApiMapper() *ResponseApiMapper {
	return &ResponseApiMapper{
		AuthResponseMapper:            auth.NewAuthResponseMapper(),
		UserResponseMapper:            user.NewUserResponseMapper(),
		RoleResponseMapper:            role.NewRoleResponseMapper(),
		CategoryResponseMapper:        category.NewCategoryResponseMapper(),
		MerchantResponseMapper:        merchant.NewMerchantResponseMapper(),
		OrderItemResponseMapper:       orderitem.NewOrderItemResponseMapper(),
		OrderResponseMapper:           order.NewOrderResponseMapper(),
		ProductResponseMapper:         product.NewProductResponseMapper(),
		TransactionResponseMapper:     transaction.NewTransactionResponseMapper(),
		CartResponseMapper:            cart.NewCartResponseMapper(),
		ReviewMapper:                  review.NewReviewResponseMapper(),
		SliderMapper:                  slider.NewSliderResponseMapper(),
		ShippingAddressResponseMapper: shippingaddress.NewShippingAddressResponseMapper(),
		BannerResponseMapper:          banner.NewBannerResponseMapper(),
		MerchantAwardResponseMapper:   merchantaward.NewMerchantAwardResponseMapper(),
		MerchantBusinessMapper:        merchantbusiness.NewMerchantBusinessResponseMapper(),
		MerchantDetailResponseMapper:  merchantdetail.NewMerchantDetailResponseMapper(),
		MerchantPolicyResponseMapper:  merchantpolicy.NewMerchantPolicyResponseMapper(),
		ReviewDetailResponseMapper:    reviewdetail.NewReviewDetailResponseMapper(),
		MerchantDocumentProMapper:     merchantdocuments.NewMerchantDocumentResponseMapper(),
		MerchantSocialLinkProtoMapper: merchantsociallink.NewMerchantSocialLinkResponseMapper(),
	}
}
