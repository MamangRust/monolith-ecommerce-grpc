package merchantsociallinkhandler

import (
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_social_link"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantSocialLink struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger     logger.LoggerInterface
	ApiHandler sharedErrors.ApiHandler
}

func RegisterMerchantSocialLinkHandler(deps *DepsMerchantSocialLink) {
	mapper := apimapper.NewMerchantSocialLinkResponseMapper()

	NewMerchantSocialLinkCommandHandleApi(&merchantSocialLinkCommandHandleDeps{
		client: pb.NewMerchantSocialCommandServiceClient(deps.Client),
		router: deps.E,
		logger:     deps.Logger,
		mapper:     mapper.CommandMapper(),
		apiHandler: deps.ApiHandler,
	})
}
