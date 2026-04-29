package merchantbusinesshandler

import (
	merchantbusiness_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_business"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_business"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
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
		client: pb.NewMerchantBusinessQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache,
	})

	NewMerchantBusinessCommandHandleApi(&merchantBusinessCommandHandleDeps{
		client:         pb.NewMerchantBusinessCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
	})
}
