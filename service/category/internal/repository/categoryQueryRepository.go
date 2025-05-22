package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type categoryQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CategoryRecordMapper) *categoryQueryRepository {
	return &categoryQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *categoryQueryRepository) FindAllCategory(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategories(r.ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordPagination(res), &totalCount, nil
}

func (r *categoryQueryRepository) FindByActive(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindByActiveCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordActivePagination(res), &totalCount, nil
}

func (r *categoryQueryRepository) FindByTrashed(req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindByTrashedCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordTrashedPagination(res), &totalCount, nil
}

func (r *categoryQueryRepository) FindById(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByID(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrFindCategoryById
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryQueryRepository) FindByIdTrashed(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByIDTrashed(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrFindCategoryByIdTrashed
	}

	return r.mapping.ToCategoryRecord(res), nil
}
