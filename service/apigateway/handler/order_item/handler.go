package orderitemhandler

import (
	orderitem_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/order_item"
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/order_item"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsOrderItem struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterOrderItemHandler(deps *DepsOrderItem) {
	mapper := apimapper.NewOrderItemResponseMapper()
	cache := orderitem_cache.NewOrderItemMencache(deps.CacheStore)

	queryClient := pborder_item.NewOrderItemQueryServiceClient(deps.Client)

	NewOrderItemQueryHandleApi(&orderItemQueryHandleDeps{
		client: queryClient,
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewOrderItemCommandHandleApi(&orderItemCommandHandleDeps{
		client: pborder_item.NewOrderItemCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		cache:  cache,
	})
}
