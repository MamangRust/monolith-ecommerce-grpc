package producthandler

import (
	product_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/product"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/product"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


type DepsProduct struct {
	Client     *grpc.ClientConn
	E          *echo.Echo
	Logger     logger.LoggerInterface
	CacheStore *cache.CacheStore
	Upload     upload_image.ImageUploads
	ApiHandler errors.ApiHandler
}



func RegisterProductHandler(deps *DepsProduct) {
	mapper := apimapper.NewProductResponseMapper()
	cache := product_cache.NewProductMencache(deps.CacheStore)

	queryClient := pb.NewProductQueryServiceClient(deps.Client)

	NewProductQueryHandleApi(&productQueryHandleDeps{
		client:     queryClient,
		router:     deps.E,
		logger:     deps.Logger,
		mapper:     mapper.QueryMapper(),
		cache:      cache,
		apiHandler: deps.ApiHandler,
	})


	NewProductCommandHandleApi(&productCommandHandleDeps{
		client:       pb.NewProductCommandServiceClient(deps.Client),
		router:       deps.E,
		logger:       deps.Logger,
		mapper:       mapper.CommandMapper(),
		cache:        cache,
		upload_image: deps.Upload,
		apiHandler:   deps.ApiHandler,
	})

}
