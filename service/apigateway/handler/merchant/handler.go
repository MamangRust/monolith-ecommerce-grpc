package merchanthandler

import (
	merchant_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsMerchant struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger logger.LoggerInterface
	CacheStore *cache.CacheStore
	UploadImage upload_image.ImageUploads
	ApiHandler errors.ApiHandler
}

func RegisterMerchantHandler(deps *DepsMerchant) {
	mapper := apimapper.NewMerchantResponseMapper()
	cache := merchant_cache.NewMerchantMencache(deps.CacheStore)

	handlers := []func(){
		setupMerchantQueryHandler(deps, mapper.QueryMapper(), cache),
		setupMerchantCommandHandler(deps, mapper.CommandMapper(), cache),
	}

	for _, h := range handlers {
		h()
	}
}

func setupMerchantQueryHandler(deps *DepsMerchant, mapper apimapper.MerchantQueryResponseMapper, cache merchant_cache.MerchantQueryCache) func() {
	return func() {
		NewMerchantQueryHandleApi(&merchantQueryHandleDeps{
			client:     pb.NewMerchantQueryServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			apiHandler: deps.ApiHandler,
		})
	}
}

func setupMerchantCommandHandler(deps *DepsMerchant, mapper apimapper.MerchantCommandResponseMapper, cache merchant_cache.MerchantCommandCache) func() {
	return func() {
		NewMerchantCommandHandleApi(&merchantCommandHandleDeps{
			client:     pb.NewMerchantCommandServiceClient(deps.Client),
			router:     deps.E,
			logger:     deps.Logger,
			mapper:     mapper,
			cache:      cache,
			upload_image: deps.UploadImage,
			apiHandler: deps.ApiHandler,
		})
	}
}
