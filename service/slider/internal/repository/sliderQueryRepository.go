package repository

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
	"golang.org/x/net/context"
)

type sliderQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.SliderMapping
}

func NewSliderQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.SliderMapping) *sliderQueryRepository {
	return &sliderQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *sliderQueryRepository) FindAllSlider(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSliders(r.ctx, reqDb)

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

func (r *sliderQueryRepository) FindByActive(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersActive(r.ctx, reqDb)

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

func (r *sliderQueryRepository) FindByTrashed(req *requests.FindAllSlider) ([]*record.SliderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetSlidersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetSlidersTrashed(r.ctx, reqDb)

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
