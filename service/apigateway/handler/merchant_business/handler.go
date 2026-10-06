package merchantbusinesshandler

import (
	merchantbusiness_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_business"
	pbmerchant_business "github.com/MamangRust/monolith-ecommerce-pb/merchant_business"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_business"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantBusiness struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
}

func RegisterMerchantBusinessHandler(deps *DepsMerchantBusiness) {
	mapper := apimapper.NewMerchantBusinessResponseMapper()
	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	cache := merchantbusiness_cache.NewMerchantBusinessMencache(deps.CacheStore)

	NewMerchantBusinessQueryHandleApi(&merchantBusinessQueryHandleDeps{
		client: pbmerchant_business.NewMerchantBusinessQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewMerchantBusinessCommandHandleApi(&merchantBusinessCommandHandleDeps{
		client:         pbmerchant_business.NewMerchantBusinessCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
	})
}
