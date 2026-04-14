package handler

import (
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/banner"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/cart"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/category"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/merchant"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/merchant_award"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/merchant_business"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/merchant_detail"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/merchant_document"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/merchant_policy"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/merchant_social_link"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/order"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/order_item"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/product"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/role"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/transaction"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/review"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/review_detail"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/shipping_address"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/slider"
	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/handler/user"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

// ServiceConnections aggregates gRPC connections to backend services.
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
	MerchantDocument *grpc.ClientConn
	MerchantSocial   *grpc.ClientConn
	MerchantPolicy   *grpc.ClientConn
	ReviewDetail     *grpc.ClientConn
	Card             *grpc.ClientConn
	Saldo            *grpc.ClientConn
	Topup            *grpc.ClientConn
	Transfer         *grpc.ClientConn
	Withdraw         *grpc.ClientConn
}

type Deps struct {
	E                  *echo.Echo
	Logger             logger.LoggerInterface
	ServiceConnections *ServiceConnections
	Cache              *cache.CacheStore
	Image              upload_image.ImageUploads
	Kafka              *kafka.Kafka
	Token              auth.TokenManager
}

func NewHandler(deps *Deps) {
	observability, _ := observability.NewObservability("apigateway", deps.Logger)
	apiHandler := errors.NewApiHandler(observability, deps.Logger)
	// _ = mencache.NewCacheApiGateway(deps.Cache) // Removed for now as it cause mismatch

	bannerhandler.RegisterBannerHandler(&bannerhandler.DepsBanner{
		Client:      deps.ServiceConnections.Banner,
		E:           deps.E,
		Logger:      deps.Logger,
		CacheStore:  deps.Cache,
	})

	carthandler.RegisterCartHandler(&carthandler.DepsCart{
		Client:     deps.ServiceConnections.Cart,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
	})

	categoryhandler.RegisterCategoryHandler(&categoryhandler.DepsCategory{
		Client:      deps.ServiceConnections.Category,
		E:           deps.E,
		Logger:      deps.Logger,
		CacheStore:  deps.Cache,
		UploadImage: deps.Image,
		ApiHandler:  apiHandler,
	})

	merchanthandler.RegisterMerchantHandler(&merchanthandler.DepsMerchant{
		Client:      deps.ServiceConnections.Merchant,
		E:           deps.E,
		Logger:      deps.Logger,
		CacheStore:  deps.Cache,
		UploadImage: deps.Image,
		ApiHandler:  apiHandler,
	})

	merchantawardhandler.RegisterMerchantAwardHandler(&merchantawardhandler.DepsMerchantAward{
		Client:     deps.ServiceConnections.MerchantAward,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
	})

	merchantbusinesshandler.RegisterMerchantBusinessHandler(&merchantbusinesshandler.DepsMerchantBusiness{
		Client:     deps.ServiceConnections.MerchantBusiness,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
	})

	merchantdocumenthandler.RegisterMerchantDocumentHandler(&merchantdocumenthandler.DepsMerchantDocument{
		Client:      deps.ServiceConnections.MerchantDocument,
		E:           deps.E,
		Logger:      deps.Logger,
		UploadImage: deps.Image,
	})
	
	merchantpolicyhandler.RegisterMerchantPolicyHandler(&merchantpolicyhandler.DepsMerchantPolicy{
		Client:     deps.ServiceConnections.MerchantPolicy,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
	})

	merchantsociallinkhandler.RegisterMerchantSocialLinkHandler(&merchantsociallinkhandler.DepsMerchantSocialLink{
		Client: deps.ServiceConnections.MerchantSocial,
		E:      deps.E,
		Logger: deps.Logger,
	})

	orderhandler.RegisterOrderHandler(&orderhandler.DepsOrder{
		Client:     deps.ServiceConnections.Order,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
	})

	orderitemhandler.RegisterOrderItemHandler(&orderitemhandler.DepsOrderItem{
		Client:     deps.ServiceConnections.OrderItem,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
	})

	producthandler.RegisterProductHandler(&producthandler.DepsProduct{
		Client:     deps.ServiceConnections.Product,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
		Upload:     deps.Image,
		ApiHandler: apiHandler,
	})


	transactionhandler.RegisterTransactionHandler(&transactionhandler.DepsTransaction{
		Client:     deps.ServiceConnections.Transaction,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
	})

	merchantdetailhandler.RegisterMerchantDetailHandler(&merchantdetailhandler.DepsMerchantDetail{
		Client:      deps.ServiceConnections.MerchantDetail,
		E:           deps.E,
		Logger:      deps.Logger,
		CacheStore:  deps.Cache,
		UploadImage: deps.Image,
		ApiHandler:  apiHandler,
	})

	rolehandler.RegisterRoleHandler(&rolehandler.DepsRole{
		Client:     deps.ServiceConnections.Role,
		E:          deps.E,
		Logger:     deps.Logger,
		CacheStore: deps.Cache,
		ApiHandler: apiHandler,
	})

	sliderhandler.RegisterSliderHandler(&sliderhandler.DepsSlider{
		Client: deps.ServiceConnections.Slider,
		E:      deps.E,
		Logger: deps.Logger,
		Cache:  deps.Cache,
		Upload: deps.Image,
	})

	reviewhandler.RegisterReviewHandler(&reviewhandler.DepsReview{
		Client: deps.ServiceConnections.Review,
		E:      deps.E,
		Logger: deps.Logger,
		Cache:  deps.Cache,
	})

	reviewdetailhandler.RegisterReviewDetailHandler(&reviewdetailhandler.DepsReviewDetail{
		Client: deps.ServiceConnections.ReviewDetail,
		E:      deps.E,
		Logger: deps.Logger,
		Cache:  deps.Cache,
		Upload: deps.Image,
	})

	shippingaddresshandler.RegisterShippingAddressHandler(&shippingaddresshandler.DepsShippingAddress{
		Client: deps.ServiceConnections.Shipping,
		E:      deps.E,
		Logger: deps.Logger,
		Cache:  deps.Cache,
	})

	userhandler.RegisterUserHandler(&userhandler.DepsUser{
		Client:     deps.ServiceConnections.User,
		E:          deps.E,
		Logger:     deps.Logger,
		Cache:      deps.Cache,
		ApiHandler: apiHandler,
	})
}
