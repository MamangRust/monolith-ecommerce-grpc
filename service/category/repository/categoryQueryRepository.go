package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
)

type categoryQueryRepository struct {
	db *db.Queries
}

func NewCategoryQueryRepository(db *db.Queries) *categoryQueryRepository {
	return &categoryQueryRepository{
		db: db,
	}
}

func (r *categoryQueryRepository) FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategories(ctx, reqDb)

	if err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}


	return res, nil
}

func (r *categoryQueryRepository) FindActive(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesActive(ctx, reqDb)

	if err != nil {
		return nil, category_errors.ErrFindByActiveCategory.WithInternal(err)
	}


	return res, nil
}

func (r *categoryQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*db.GetCategoriesTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesTrashed(ctx, reqDb)

	if err != nil {
		return nil, category_errors.ErrFindByTrashedCategory.WithInternal(err)
	}


	return res, nil
}

func (r *categoryQueryRepository) FindByID(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error) {
	res, err := r.db.GetCategoryByID(ctx, int32(category_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, category_errors.ErrCategoryNotFound.WithInternal(err)
		}
		return nil, category_errors.ErrFindCategoryById.WithInternal(err)
	}


	return res, nil
}

func (r *categoryQueryRepository) FindByIDTrashed(ctx context.Context, category_id int) (*db.Category, error) {
	res, err := r.db.GetCategoryByIDTrashed(ctx, int32(category_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, category_errors.ErrCategoryNotFound.WithInternal(err)
		}
		return nil, category_errors.ErrFindCategoryByIdTrashed.WithInternal(err)
	}


	return res, nil
}
