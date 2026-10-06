package merchantawardhandler

import (
	merchantaward_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_awards"
	pbmerchant_award "github.com/MamangRust/monolith-ecommerce-pb/merchant_award"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_award"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
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
		client: pbmerchant_award.NewMerchantAwardQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewMerchantAwardCommandHandleApi(&merchantAwardCommandHandleDeps{
		client:         pbmerchant_award.NewMerchantAwardCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
	})
}
