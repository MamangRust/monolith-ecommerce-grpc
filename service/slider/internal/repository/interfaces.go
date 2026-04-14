package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type SliderQueryRepository interface {
	FindAllSlider(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersRow, error)

	FindByActive(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersActiveRow, error)

	FindByTrashed(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersTrashedRow, error)

	FindById(
		ctx context.Context,
		slider_id int,
	) (*db.GetSliderByIDRow, error)
}

type SliderCommandRepository interface {
	CreateSlider(
		ctx context.Context,
		request *requests.CreateSliderRequest,
	) (*db.CreateSliderRow, error)

	UpdateSlider(
		ctx context.Context,
		request *requests.UpdateSliderRequest,
	) (*db.UpdateSliderRow, error)

	TrashSlider(
		ctx context.Context,
		slider_id int,
	) (*db.Slider, error)

	RestoreSlider(
		ctx context.Context,
		slider_id int,
	) (*db.Slider, error)

	DeleteSliderPermanently(
		ctx context.Context,
		slider_id int,
	) (bool, error)

	RestoreAllSlider(ctx context.Context) (bool, error)
	DeleteAllPermanentSlider(ctx context.Context) (bool, error)
}
