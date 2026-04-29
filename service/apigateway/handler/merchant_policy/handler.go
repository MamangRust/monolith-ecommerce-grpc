package merchantpolicyhandler

import (
	merchantpolicy_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_policies"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_policy"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchantPolicy struct {
	Client        *grpc.ClientConn
	E             *echo.Echo
	Logger        logger.LoggerInterface
	CacheStore    *cache.CacheStore
	Observability observability.TraceLoggerObservability
}

func RegisterMerchantPolicyHandler(deps *DepsMerchantPolicy) {
	mapper := apimapper.NewMerchantPolicyResponseMapper()
	merchantMapper := merchantapimapper.NewMerchantResponseMapper()
	cache := merchantpolicy_cache.NewMerchantPoliciesMencache(deps.CacheStore)

	NewMerchantPolicyQueryHandleApi(&merchantPolicyQueryHandleDeps{
		client:        pb.NewMerchantPolicyQueryServiceClient(deps.Client),
		router:        deps.E,
		logger:        deps.Logger,
		mapper:        mapper.QueryMapper(),
		cache:         cache,
		observability: deps.Observability,
	})

	NewMerchantPolicyCommandHandleApi(&merchantPolicyCommandHandleDeps{
		client:         pb.NewMerchantPolicyCommandServiceClient(deps.Client),
		router:         deps.E,
		logger:         deps.Logger,
		mapper:         mapper.CommandMapper(),
		merchantMapper: merchantMapper.CommandMapper(),
		cache:          cache,
		observability:  deps.Observability,
	})
}
