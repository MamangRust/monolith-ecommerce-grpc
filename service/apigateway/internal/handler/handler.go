package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type Deps struct {
	Conn    *grpc.ClientConn
	Token   auth.TokenManager
	E       *echo.Echo
	Logger  logger.LoggerInterface
	Mapping response_api.ResponseApiMapper
	Image   upload_image.ImageUploads
}

func NewHandler(deps Deps) {

	clientAuth := pb.NewAuthServiceClient(deps.Conn)
	clientRole := pb.NewRoleServiceClient(deps.Conn)
	clientUser := pb.NewUserServiceClient(deps.Conn)
	clientCategory := pb.NewCategoryServiceClient(deps.Conn)
	clientMerchant := pb.NewMerchantServiceClient(deps.Conn)
	clientOrderItem := pb.NewOrderItemServiceClient(deps.Conn)
	clientOrder := pb.NewOrderServiceClient(deps.Conn)
	clientProduct := pb.NewProductServiceClient(deps.Conn)
	clientTransaction := pb.NewTransactionServiceClient(deps.Conn)
	clientCart := pb.NewCartServiceClient(deps.Conn)
	clientReview := pb.NewReviewServiceClient(deps.Conn)
	clientSlider := pb.NewSliderServiceClient(deps.Conn)
	clientShipping := pb.NewShippingServiceClient(deps.Conn)
	clientBanner := pb.NewBannerServiceClient(deps.Conn)
	clientMerchantAward := pb.NewMerchantAwardServiceClient(deps.Conn)
	clientMerchantBusiness := pb.NewMerchantBusinessServiceClient(deps.Conn)
	clientMerchantDetail := pb.NewMerchantDetailServiceClient(deps.Conn)
	clientMerchantPolicy := pb.NewMerchantPoliciesServiceClient(deps.Conn)
	clientReviewDetail := pb.NewReviewDetailServiceClient(deps.Conn)

	NewHandlerAuth(deps.E, clientAuth, deps.Logger, deps.Mapping.AuthResponseMapper)
	NewHandlerRole(deps.E, clientRole, deps.Logger, deps.Mapping.RoleResponseMapper)
	NewHandlerUser(deps.E, clientUser, deps.Logger, deps.Mapping.UserResponseMapper)
	NewHandlerCategory(deps.E, clientCategory, deps.Logger, deps.Mapping.CategoryResponseMapper, deps.Image)
	NewHandlerMerchant(deps.E, clientMerchant, deps.Logger, deps.Mapping.MerchantResponseMapper)
	NewHandlerOrderItem(deps.E, clientOrderItem, deps.Logger, deps.Mapping.OrderItemResponseMapper)
	NewHandlerOrder(deps.E, clientOrder, deps.Logger, deps.Mapping.OrderResponseMapper)
	NewHandlerProduct(deps.E, clientProduct, deps.Logger, deps.Mapping.ProductResponseMapper, deps.Image)
	NewHandlerTransaction(deps.E, clientTransaction, deps.Logger, deps.Mapping.TransactionResponseMapper)
	NewHandlerCart(deps.E, clientCart, deps.Logger, deps.Mapping.CartResponseMapper)
	NewHandlerReview(deps.E, clientReview, deps.Logger, deps.Mapping.ReviewMapper)
	NewHandlerSlider(deps.E, clientSlider, deps.Logger, deps.Mapping.SliderMapper, deps.Image)
	NewHandlerShippingAddress(deps.E, clientShipping, deps.Logger, deps.Mapping.ShippingAddressResponseMapper)
	NewHandleBanner(deps.E, clientBanner, deps.Logger, deps.Mapping.BannerResponseMapper)
	NewHandlerMerchantAward(deps.E, clientMerchantAward, deps.Logger, deps.Mapping.MerchantAwardResponseMapper, deps.Mapping.MerchantResponseMapper)
	NewHandlerMerchantBusiness(deps.E, clientMerchantBusiness, deps.Logger, deps.Mapping.MerchantBusinessMapper, deps.Mapping.MerchantResponseMapper)
	NewHandlerMerchantDetail(deps.E, clientMerchantDetail, deps.Logger, deps.Mapping.MerchantDetailResponseMapper, deps.Mapping.MerchantResponseMapper, deps.Image)
	NewHandlerMerchantPolicies(deps.E, clientMerchantPolicy, deps.Logger, deps.Mapping.MerchantPolicyResponseMapper, deps.Mapping.MerchantResponseMapper)
	NewHandlerReviewDetail(deps.E, clientReviewDetail, deps.Logger, deps.Mapping.ReviewDetailResponseMapper, deps.Mapping.ReviewMapper, deps.Image)
}
