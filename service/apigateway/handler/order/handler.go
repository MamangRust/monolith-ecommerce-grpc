package orderhandler

import (
	order_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/order"
	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/order"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsOrder struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterOrderHandler(deps *DepsOrder) {
	mapper := apimapper.NewOrderResponseMapper()
	cache := order_cache.OrderNewMencache(deps.CacheStore)

	queryClient := pborder.NewOrderQueryServiceClient(deps.Client)
	statsClient := pborder.NewOrderStatsServiceClient(deps.Client)
	statsByMerchantClient := pborder.NewOrderStatsByMerchantServiceClient(deps.Client)

	NewOrderQueryHandleApi(&orderQueryHandleDeps{
		client: queryClient,
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewOrderCommandHandleApi(&orderCommandHandleDeps{
		client: pborder.NewOrderCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})

	NewOrderStatsHandleApi(&orderStatsHandleDeps{
		statsClient:           statsClient,
		statsByMerchantClient: statsByMerchantClient,
		router:                deps.E,
		logger:                deps.Logger,
		mapper:                mapper.StatsMapper(),
		cache:                 cache,
		merchantStatsCache:    cache,
	})
}
