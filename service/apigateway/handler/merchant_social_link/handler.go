package merchantsociallinkhandler

import (
	pbmerchant_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchant_social_link"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_social_link"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantSocialLink struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	ApiHandler sharedErrors.ApiHandler
}

func RegisterMerchantSocialLinkHandler(deps *DepsMerchantSocialLink) {
	mapper := apimapper.NewMerchantSocialLinkResponseMapper()

	NewMerchantSocialLinkCommandHandleApi(&merchantSocialLinkCommandHandleDeps{
		client:     pbmerchant_social_link.NewMerchantSocialCommandServiceClient(deps.Client),
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		apiHandler: deps.ApiHandler,
	})
}
