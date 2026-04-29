package merchantawardhandler

import (
	merchantaward_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_awards"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_award"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
)

type DepsMerchantAward struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterMerchantAwardHandler(deps *DepsMerchantAward) {
	mapper := apimapper.NewMerchantAwardResponseMapper()
	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	cache := merchantaward_cache.NewMerchantAward(deps.CacheStore)

	NewMerchantAwardQueryHandleApi(&merchantAwardQueryHandleDeps{
		client: pb.NewMerchantAwardQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewMerchantAwardCommandHandleApi(&merchantAwardCommandHandleDeps{
		client:         pb.NewMerchantAwardCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
	})
}
