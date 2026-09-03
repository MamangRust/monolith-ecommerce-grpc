package rolehandler

import (
	api_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache"
	role_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/role"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/role"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsRole struct {
	Kafka      *kafka.Kafka
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	Cache      api_cache.RoleCache
	ApiHandler errors.ApiHandler
}

func RegisterRoleHandler(deps *DepsRole) {
	mapper := apimapper.NewRoleResponseMapper()
	cache := role_cache.NewRoleMencache(deps.CacheStore)

	NewRoleQueryHandleApi(&roleQueryHandleDeps{
		client:     pb.NewRoleQueryServiceClient(deps.Client),
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper.QueryMapper(),
		kafka:      deps.Kafka,
		cache_role: deps.Cache,
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})

	NewRoleCommandHandleApi(&roleCommandHandleDeps{
		client:      pb.NewRoleCommandServiceClient(deps.Client),
		queryClient: pb.NewRoleQueryServiceClient(deps.Client),
		router:      deps.E,
		logger:      deps.Logger,
		mapper:      mapper.CommandMapper(),
		kafka:       deps.Kafka,
		cache_role:  deps.Cache,
		cache:       cache,
		apiHandler:  deps.ApiHandler,
	})
}
