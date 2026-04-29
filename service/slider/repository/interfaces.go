package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type SliderQueryRepository interface {
	FindAll(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersRow, error)

	FindActive(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersActiveRow, error)

	FindTrashed(
		ctx context.Context,
		req *requests.FindAllSlider,
	) ([]*db.GetSlidersTrashedRow, error)

	FindByID(
		ctx context.Context,
		slider_id int,
	) (*db.GetSliderByIDRow, error)
}

type SliderCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateSliderRequest,
	) (*db.CreateSliderRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateSliderRequest,
	) (*db.UpdateSliderRow, error)

	Trash(
		ctx context.Context,
		slider_id int,
	) (*db.Slider, error)

	Restore(
		ctx context.Context,
		slider_id int,
	) (*db.Slider, error)

	DeletePermanent(
		ctx context.Context,
		slider_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
