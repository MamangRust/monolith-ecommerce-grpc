package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
)

type sliderCommandRepository struct {
	db *db.Queries
}

func NewSliderCommandRepository(db *db.Queries) *sliderCommandRepository {
	return &sliderCommandRepository{
		db: db,
	}
}

func (r *sliderCommandRepository) Create(ctx context.Context, request *requests.CreateSliderRequest) (*db.CreateSliderRow, error) {
	req := db.CreateSliderParams{
		Name:  request.Nama,
		Image: request.FilePath,
	}

	slider, err := r.db.CreateSlider(ctx, req)

	if err != nil {
		return nil, slider_errors.ErrCreateSlider
	}

	return slider, nil
}

func (r *sliderCommandRepository) Update(ctx context.Context, request *requests.UpdateSliderRequest) (*db.UpdateSliderRow, error) {
	req := db.UpdateSliderParams{
		SliderID: int32(*request.ID),
		Name:     request.Nama,
		Image:    request.FilePath,
	}

	res, err := r.db.UpdateSlider(ctx, req)

	if err != nil {
		return nil, slider_errors.ErrUpdateSlider
	}

	return res, nil
}

func (r *sliderCommandRepository) Trash(ctx context.Context, slider_id int) (*db.Slider, error) {
	res, err := r.db.TrashSlider(ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrTrashSlider
	}

	return res, nil
}

func (r *sliderCommandRepository) Restore(ctx context.Context, slider_id int) (*db.Slider, error) {
	res, err := r.db.RestoreSlider(ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}

	return res, nil
}

func (r *sliderCommandRepository) DeletePermanent(ctx context.Context, slider_id int) (bool, error) {
	err := r.db.DeleteSliderPermanently(ctx, int32(slider_id))

	if err != nil {
		return false, slider_errors.ErrDeletePermanentSlider
	}

	return true, nil
}

func (r *sliderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllSliders(ctx)

	if err != nil {
		return false, slider_errors.ErrRestoreAllSlider
	}
	return true, nil
}

func (r *sliderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentSliders(ctx)

	if err != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider
	}
	return true, nil
}
