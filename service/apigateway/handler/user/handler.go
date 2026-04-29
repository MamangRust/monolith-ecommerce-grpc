package userhandler

import (
	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/user"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/user"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsUser struct {
	Client *grpc.ClientConn

	E *echo.Echo

	Logger logger.LoggerInterface

	Cache *cache.CacheStore

	ApiHandler errors.ApiHandler
}

func RegisterUserHandler(deps *DepsUser) {
	mapper := apimapper.NewUserResponseMapper()

	cache := user_cache.NewUserMencache(deps.Cache)

	handlers := []func(){
		setupUserQueryHandler(deps, mapper.QueryMapper(), cache),
		setupUserCommandHandler(deps, mapper.CommandMapper(), cache),
	}

	for _, h := range handlers {
		h()
	}
}

func setupUserQueryHandler(deps *DepsUser, mapper apimapper.UserQueryResponseMapper, cache user_cache.UserMencache) func() {
	return func() {
		NewUserQueryHandleApi(&userQueryHandleDeps{
			client:     pb.NewUserQueryServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}

func setupUserCommandHandler(deps *DepsUser, mapper apimapper.UserCommandResponseMapper, cache user_cache.UserMencache) func() {
	return func() {
		NewUserCommandHandleApi(&userCommandHandleDeps{
			client:     pb.NewUserCommandServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}
