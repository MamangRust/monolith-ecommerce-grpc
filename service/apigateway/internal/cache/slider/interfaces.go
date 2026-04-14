package slider_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderQueryCache interface {
	GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) (*response.ApiResponsePaginationSlider, bool)
	SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data *response.ApiResponsePaginationSlider)

	GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) (*response.ApiResponsePaginationSliderDeleteAt, bool)
	SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data *response.ApiResponsePaginationSliderDeleteAt)

	GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) (*response.ApiResponsePaginationSliderDeleteAt, bool)
	SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data *response.ApiResponsePaginationSliderDeleteAt)
}

type SliderCommandCache interface {
	DeleteSliderCache(ctx context.Context, slider_id int)
}
