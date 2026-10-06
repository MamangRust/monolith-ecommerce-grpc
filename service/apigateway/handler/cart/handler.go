package carthandler

import (
	cart_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/cart"
	pbcart "github.com/MamangRust/monolith-ecommerce-pb/cart"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/cart"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsCart struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterCartHandler(deps *DepsCart) {
	mapper := apimapper.NewCartResponseMapper()
	cache := cart_cache.NewCartMencache(deps.CacheStore)

	NewCartQueryHandleApi(&cartQueryHandleDeps{
		client: pbcart.NewCartQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewCartCommandHandleApi(&cartCommandHandleDeps{
		client: pbcart.NewCartCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})
}
