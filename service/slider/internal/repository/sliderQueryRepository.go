package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type sliderQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.SliderMapping
}

func NewSliderQueryRepository(db *db.Queries, mapping recordmapper.SliderMapping) *sliderQueryRepository {
	return &sliderQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *sliderQueryRepository) FindAllSlider(ctx context.Context, req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSliders(ctx, reqDb)

	if err != nil {
		return nil, nil, slider_errors.ErrFindAllSliders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordPagination(res), &totalCount, nil
}

func (r *sliderQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersActive(ctx, reqDb)

	if err != nil {
		return nil, nil, slider_errors.ErrFindActiveSliders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordActivePagination(res), &totalCount, nil
}

func (r *sliderQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, slider_errors.ErrFindTrashedSliders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToSlidersRecordTrashedPagination(res), &totalCount, nil
}
