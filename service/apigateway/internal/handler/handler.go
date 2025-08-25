package handler

import (
	"strconv"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type ServiceConnections struct {
	Auth             *grpc.ClientConn
	Role             *grpc.ClientConn
	User             *grpc.ClientConn
	Category         *grpc.ClientConn
	Merchant         *grpc.ClientConn
	OrderItem        *grpc.ClientConn
	Order            *grpc.ClientConn
	Product          *grpc.ClientConn
	Transaction      *grpc.ClientConn
	Cart             *grpc.ClientConn
	Review           *grpc.ClientConn
	Slider           *grpc.ClientConn
	Shipping         *grpc.ClientConn
	Banner           *grpc.ClientConn
	MerchantAward    *grpc.ClientConn
	MerchantBusiness *grpc.ClientConn
	MerchantDetail   *grpc.ClientConn
	MerchantPolicy   *grpc.ClientConn
	ReviewDetail     *grpc.ClientConn
}

type Deps struct {
	Caceh              mencache.RoleCache
	Kafka              *kafka.Kafka
	Token              auth.TokenManager
	E                  *echo.Echo
	Logger             logger.LoggerInterface
	Mapping            *response_api.ResponseApiMapper
	Image              upload_image.ImageUploads
	ServiceConnections *ServiceConnections
}

func NewHandler(deps *Deps) {

	clientAuth := pb.NewAuthServiceClient(deps.ServiceConnections.Auth)
	clientRole := pb.NewRoleServiceClient(deps.ServiceConnections.Role)
	clientUser := pb.NewUserServiceClient(deps.ServiceConnections.User)
	clientCategory := pb.NewCategoryServiceClient(deps.ServiceConnections.Category)
	clientMerchant := pb.NewMerchantServiceClient(deps.ServiceConnections.Merchant)
	clientOrderItem := pb.NewOrderItemServiceClient(deps.ServiceConnections.OrderItem)
	clientOrder := pb.NewOrderServiceClient(deps.ServiceConnections.Order)
	clientProduct := pb.NewProductServiceClient(deps.ServiceConnections.Product)
	clientMerchantDocument := pb.NewMerchantDocumentServiceClient(deps.ServiceConnections.Merchant)
	clientTransaction := pb.NewTransactionServiceClient(deps.ServiceConnections.Transaction)
	clientCart := pb.NewCartServiceClient(deps.ServiceConnections.Cart)
	clientReview := pb.NewReviewServiceClient(deps.ServiceConnections.Review)
	clientSlider := pb.NewSliderServiceClient(deps.ServiceConnections.Slider)
	clientShipping := pb.NewShippingServiceClient(deps.ServiceConnections.Shipping)
	clientBanner := pb.NewBannerServiceClient(deps.ServiceConnections.Banner)
	clientMerchantAward := pb.NewMerchantAwardServiceClient(deps.ServiceConnections.MerchantAward)
	clientMerchantBusiness := pb.NewMerchantBusinessServiceClient(deps.ServiceConnections.MerchantBusiness)
	clientMerchantDetail := pb.NewMerchantDetailServiceClient(deps.ServiceConnections.MerchantDetail)
	clientMerchantPolicy := pb.NewMerchantPoliciesServiceClient(deps.ServiceConnections.MerchantPolicy)
	clientReviewDetail := pb.NewReviewDetailServiceClient(deps.ServiceConnections.ReviewDetail)
	clientMerchantSocial := pb.NewMerchantSocialServiceClient(deps.ServiceConnections.MerchantDetail)

	NewHandlerAuth(deps.E, clientAuth, deps.Logger, deps.Mapping.AuthResponseMapper)
	NewHandlerRole(deps.Caceh, deps.E, clientRole, deps.Logger, deps.Mapping.RoleResponseMapper, deps.Kafka)
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
	NewHandlerMerchantSocialLink(deps.E, clientMerchantSocial, deps.Logger, deps.Mapping.MerchantSocialLinkProtoMapper)
	NewHandlerMerchantDocument(deps.E, clientMerchantDocument, deps.Logger, deps.Mapping.MerchantDocumentProMapper, deps.Image)
}

func parseQueryInt(c echo.Context, key string, defaultValue int) int {
	val, err := strconv.Atoi(c.QueryParam(key))
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}
