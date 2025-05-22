package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
	"golang.org/x/net/context"
)

type sliderCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SliderMapping
}

func NewSliderCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.SliderMapping) *sliderCommandRepository {
	return &sliderCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *sliderCommandRepository) CreateSlider(request *requests.CreateSliderRequest) (*record.SliderRecord, error) {
	req := db.CreateSliderParams{
		Name:  request.Nama,
		Image: request.FilePath,
	}

	slider, err := r.db.CreateSlider(r.ctx, req)

	if err != nil {
		return nil, slider_errors.ErrCreateSlider
	}

	return r.mapping.ToSliderRecord(slider), nil
}

func (r *sliderCommandRepository) UpdateSlider(request *requests.UpdateSliderRequest) (*record.SliderRecord, error) {
	req := db.UpdateSliderParams{
		SliderID: int32(*request.ID),
		Name:     request.Nama,
		Image:    request.FilePath,
	}

	res, err := r.db.UpdateSlider(r.ctx, req)

	if err != nil {
		return nil, slider_errors.ErrUpdateSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderCommandRepository) TrashSlider(slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.TrashSlider(r.ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrTrashSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderCommandRepository) RestoreSlider(slider_id int) (*record.SliderRecord, error) {
	res, err := r.db.RestoreSlider(r.ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrRestoreSlider
	}

	return r.mapping.ToSliderRecord(res), nil
}

func (r *sliderCommandRepository) DeleteSliderPermanently(slider_id int) (bool, error) {
	err := r.db.DeleteSliderPermanently(r.ctx, int32(slider_id))

	if err != nil {
		return false, slider_errors.ErrDeletePermanentSlider
	}

	return true, nil
}

func (r *sliderCommandRepository) RestoreAllSlider() (bool, error) {
	err := r.db.RestoreAllSliders(r.ctx)

	if err != nil {
		return false, slider_errors.ErrRestoreAllSlider
	}
	return true, nil
}

func (r *sliderCommandRepository) DeleteAllPermanentSlider() (bool, error) {
	err := r.db.DeleteAllPermanentSliders(r.ctx)

	if err != nil {
		return false, slider_errors.ErrDeleteAllPermanentSlider
	}
	return true, nil
}
