package service

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderQueryService interface {
	FindAll(req *requests.FindAllSlider) ([]*response.SliderResponse, *int, *response.ErrorResponse)
	FindByActive(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse)
}

type SliderCommandService interface {
	CreateSlider(req *requests.CreateSliderRequest) (*response.SliderResponse, *response.ErrorResponse)
	UpdateSlider(req *requests.UpdateSliderRequest) (*response.SliderResponse, *response.ErrorResponse)
	TrashedSlider(slider_id int) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	RestoreSlider(sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	DeleteSliderPermanent(sliderID int) (bool, *response.ErrorResponse)
	RestoreAllSliders() (bool, *response.ErrorResponse)
	DeleteAllSlidersPermanent() (bool, *response.ErrorResponse)
}
