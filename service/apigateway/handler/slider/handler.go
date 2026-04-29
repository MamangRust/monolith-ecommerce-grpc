package sliderhandler

import (
	slider_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/slider"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/cache"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/slider"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type DepsSlider struct {
	Client *grpc.ClientConn
	E      *echo.Echo
	Logger logger.LoggerInterface
	Cache  *cache.CacheStore
	Upload upload_image.ImageUploads
}

func RegisterSliderHandler(deps *DepsSlider) {
	mapper := apimapper.NewSliderResponseMapper()
	cache := slider_cache.NewSliderMencache(deps.Cache)

	NewSliderQueryHandleApi(&sliderQueryHandleDeps{
		client: pb.NewSliderQueryServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.QueryMapper(),
		cache:  cache.QueryCache(),
	})

	NewSliderCommandHandleApi(&sliderCommandHandleDeps{
		client: pb.NewSliderCommandServiceClient(deps.Client),
		router: deps.E,
		logger: deps.Logger,
		mapper: mapper.CommandMapper(),
		queryMapper: mapper.QueryMapper(),
		cache:  cache.CommandCache(),
		upload: deps.Upload,
	})
}
