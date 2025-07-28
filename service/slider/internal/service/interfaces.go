package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponse, *int, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*response.SliderResponseDeleteAt, *int, *response.ErrorResponse)
}

type SliderCommandService interface {
	CreateSlider(ctx context.Context, req *requests.CreateSliderRequest) (*response.SliderResponse, *response.ErrorResponse)
	UpdateSlider(ctx context.Context, req *requests.UpdateSliderRequest) (*response.SliderResponse, *response.ErrorResponse)
	TrashedSlider(ctx context.Context, sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	RestoreSlider(ctx context.Context, sliderID int) (*response.SliderResponseDeleteAt, *response.ErrorResponse)
	DeleteSliderPermanent(ctx context.Context, sliderID int) (bool, *response.ErrorResponse)
	RestoreAllSliders(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllSlidersPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
