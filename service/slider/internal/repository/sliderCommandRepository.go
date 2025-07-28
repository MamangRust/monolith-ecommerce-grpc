package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type sliderCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.SliderMapping
}

func NewSliderCommandRepository(db *db.Queries, mapping recordmapper.SliderMapping) *sliderCommandRepository {
	return &sliderCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *sliderCommandRepository) CreateSlider(ctx context.Context, request *requests.CreateSliderRequest) (*record.SliderRecord, error) {
	req := db.CreateSliderParams{
		Name:  request.Nama,
		Image: request.FilePath,
	}

	slider, err := r.db.CreateSlider(ctx, req)

	if err != nil {
		return nil, slider_errors.ErrCreateSlider
	}

	return r.mapping.ToSliderRecord(slider), nil
}

func (r *sliderCommandRepository) UpdateSlider(ctx context.Context, request *requests.UpdateSliderRequest) (*record.SliderRecord, error) {
	req := db.UpdateSliderParams{
		SliderID: int32(*request.ID),
		Name:     request.Nama,
		Image:    request.FilePath,
	}

	res, err := r.db.UpdateSlider(ctx, req)

	if err != nil {
		return nil, slider_errors.ErrUpdateSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderCommandRepository) TrashSlider(ctx context.Context, slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.TrashSlider(ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrTrashSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderCommandRepository) RestoreSlider(ctx context.Context, slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.RestoreSlider(ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderCommandRepository) DeleteSliderPermanently(ctx context.Context, slider_id int) (bool, error) {
	err := r.db.DeleteSliderPermanently(ctx, int32(slider_id))

	if err != nil {
		return false, slider_errors.ErrDeletePermanentSlider
	}

	return true, nil
}

func (r *sliderCommandRepository) RestoreAllSlider(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllSliders(ctx)

	if err != nil {
		return false, slider_errors.ErrRestoreAllSlider
	}
	return true, nil
}

func (r *sliderCommandRepository) DeleteAllPermanentSlider(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentSliders(ctx)

	if err != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider
	}
	return true, nil
}
