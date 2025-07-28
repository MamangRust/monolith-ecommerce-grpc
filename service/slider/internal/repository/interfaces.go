package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type SliderQueryRepository interface {
	FindAllSlider(ctx context.Context, req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error)
}

type SliderCommandRepository interface {
	CreateSlider(ctx context.Context, request *requests.CreateSliderRequest) (*record.SliderRecord, error)
	UpdateSlider(ctx context.Context, request *requests.UpdateSliderRequest) (*record.SliderRecord, error)
	TrashSlider(ctx context.Context, slider_id int) (*record.SliderRecord, error)
	RestoreSlider(ctx context.Context, slider_id int) (*record.SliderRecord, error)
	DeleteSliderPermanently(ctx context.Context, slider_id int) (bool, error)
	RestoreAllSlider(ctx context.Context) (bool, error)
	DeleteAllPermanentSlider(ctx context.Context) (bool, error)
}
