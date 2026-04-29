package carthandler

import (
	cart_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/cart"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/cart"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
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
		client: pb.NewCartQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewCartCommandHandleApi(&cartCommandHandleDeps{
		client: pb.NewCartCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
	})
}
