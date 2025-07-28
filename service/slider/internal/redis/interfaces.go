package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderQueryCache interface {
	GetSliderAllCache(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponse, *int, bool)
	SetSliderAllCache(ctx context.Context, req *requests.FindAllSlider, data []*response.SliderResponse, total *int)

	GetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, bool)
	SetSliderActiveCache(ctx context.Context, req *requests.FindAllSlider, data []*response.SliderResponseDeleteAt, total *int)

	GetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, bool)
	SetSliderTrashedCache(ctx context.Context, req *requests.FindAllSlider, data []*response.SliderResponseDeleteAt, total *int)
}

type SliderCommandCache interface {
	DeleteSliderCache(ctx context.Context, review_id int)
}
