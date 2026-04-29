package shippingaddresshandler

import (
	shippingaddress_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/shipping_address"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/shipping_address"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsShippingAddress struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger logger.LoggerInterface
	Cache  *cache.CacheStore
}

func RegisterShippingAddressHandler(deps *DepsShippingAddress) {
	mapper := apimapper.NewShippingAddressResponseMapper()
	cache := shippingaddress_cache.NewShippingAddressMencache(deps.Cache)

	NewShippingAddressQueryHandleApi(&shippingAddressQueryHandleDeps{
		client: pb.NewShippingQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache.QueryCache(),
	})

	NewShippingAddressCommandHandleApi(&shippingAddressCommandHandleDeps{
		client: pb.NewShippingCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache.CommandCache(),
	})
}
