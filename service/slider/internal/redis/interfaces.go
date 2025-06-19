package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderQueryCache interface {
	GetSliderAllCache(req *requests.FindAllSlider) ([]*response.SliderResponse, *int, bool)
	SetSliderAllCache(req *requests.FindAllSlider, data []*response.SliderResponse, total *int)

	GetSliderActiveCache(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, bool)
	SetSliderActiveCache(req *requests.FindAllSlider, data []*response.SliderResponseDeleteAt, total *int)

	GetSliderTrashedCache(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, bool)
	SetSliderTrashedCache(req *requests.FindAllSlider, data []*response.SliderResponseDeleteAt, total *int)
}

type SliderCommandCache interface {
	DeleteSliderCache(review_id int)
}
