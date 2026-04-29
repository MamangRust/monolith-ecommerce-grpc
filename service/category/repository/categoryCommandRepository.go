package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
)

type categoryCommandRepository struct {
	db *db.Queries
}

func NewCategoryCommandRepository(db *db.Queries) *categoryCommandRepository {
	return &categoryCommandRepository{
		db: db,
	}
}

func (r *categoryCommandRepository) Create(ctx context.Context, request *requests.CreateCategoryRequest) (*db.CreateCategoryRow, error) {
	req := db.CreateCategoryParams{
		Name:          request.Name,
		Description:   &request.Description,
		SlugCategory:  request.SlugCategory,
		ImageCategory: &request.ImageCategory,
	}

	category, err := r.db.CreateCategory(ctx, req)
	if err != nil {
		return nil, category_errors.ErrCreateCategory.WithInternal(err)
	}


	return category, nil
}

func (r *categoryCommandRepository) Update(ctx context.Context, request *requests.UpdateCategoryRequest) (*db.UpdateCategoryRow, error) {
	req := db.UpdateCategoryParams{
		CategoryID:    int32(*request.CategoryID),
		Name:          request.Name,
		Description:   &request.Description,
		SlugCategory:  request.SlugCategory,
		ImageCategory: &request.ImageCategory,
	}

	res, err := r.db.UpdateCategory(ctx, req)

	if err != nil {
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}


	return res, nil
}

func (r *categoryCommandRepository) Trash(ctx context.Context, category_id int) (*db.Category, error) {
	res, err := r.db.TrashCategory(ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}


	return res, nil
}

func (r *categoryCommandRepository) Restore(ctx context.Context, category_id int) (*db.Category, error) {
	res, err := r.db.RestoreCategory(ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}


	return res, nil
}

func (r *categoryCommandRepository) DeletePermanent(ctx context.Context, category_id int) (bool, error) {
	err := r.db.DeleteCategoryPermanently(ctx, int32(category_id))

	if err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently.WithInternal(err)
	}


	return true, nil
}

func (r *categoryCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllCategories(ctx)

	if err != nil {
		return false, category_errors.ErrRestoreAllCategories.WithInternal(err)
	}

	return true, nil
}

func (r *categoryCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentCategories(ctx)

	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories.WithInternal(err)
	}

	return true, nil
}
