package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
)

type sliderQueryRepository struct {
	db *db.Queries
}

func NewSliderQueryRepository(db *db.Queries) *sliderQueryRepository {
	return &sliderQueryRepository{
		db: db,
	}
}

func (r *sliderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSliders(ctx, reqDb)

	if err != nil {
		return nil, slider_errors.ErrFindAllSliders
	}

	return res, nil
}

func (r *sliderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersActive(ctx, reqDb)

	if err != nil {
		return nil, slider_errors.ErrFindActiveSliders
	}

	return res, nil
}

func (r *sliderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*db.GetSlidersTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersTrashed(ctx, reqDb)

	if err != nil {
		return nil, slider_errors.ErrFindTrashedSliders
	}

	return res, nil
}

func (r *sliderQueryRepository) FindByID(ctx context.Context, slider_id int) (*db.GetSliderByIDRow, error) {
	res, err := r.db.GetSliderByID(ctx, int32(slider_id))

	if err != nil {
		return nil, slider_errors.ErrFindSliderByID
	}

	return res, nil
}
