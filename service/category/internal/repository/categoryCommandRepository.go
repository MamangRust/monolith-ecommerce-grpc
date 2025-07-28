package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type categoryCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryCommandRepository(db *db.Queries, mapping recordmapper.CategoryRecordMapper) *categoryCommandRepository {
	return &categoryCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *categoryCommandRepository) CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error) {
	req := db.CreateCategoryParams{
		Name: request.Name,
		Description: sql.NullString{
			String: request.Description,
			Valid:  true,
		},
		SlugCategory: sql.NullString{
			String: *request.SlugCategory,
			Valid:  true,
		},
		ImageCategory: sql.NullString{
			String: request.ImageCategory,
			Valid:  true,
		},
	}

	category, err := r.db.CreateCategory(ctx, req)
	if err != nil {
		return nil, category_errors.ErrCreateCategory
	}

	return r.mapping.ToCategoryRecord(category), nil
}

func (r *categoryCommandRepository) UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error) {
	req := db.UpdateCategoryParams{
		CategoryID: int32(*request.CategoryID),
		Name:       request.Name,
		Description: sql.NullString{
			String: request.Description,
			Valid:  true,
		},
		SlugCategory: sql.NullString{
			String: *request.SlugCategory,
			Valid:  true,
		},
		ImageCategory: sql.NullString{
			String: request.ImageCategory,
			Valid:  true,
		},
	}

	res, err := r.db.UpdateCategory(ctx, req)

	if err != nil {
		return nil, category_errors.ErrUpdateCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryCommandRepository) TrashedCategory(ctx context.Context, category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.TrashCategory(ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrTrashedCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryCommandRepository) RestoreCategory(ctx context.Context, category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.RestoreCategory(ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrRestoreCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryCommandRepository) DeleteCategoryPermanently(ctx context.Context, category_id int) (bool, error) {
	err := r.db.DeleteCategoryPermanently(ctx, int32(category_id))

	if err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently
	}

	return true, nil
}

func (r *categoryCommandRepository) RestoreAllCategories(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllCategories(ctx)

	if err != nil {
		return false, category_errors.ErrRestoreAllCategories
	}
	return true, nil
}

func (r *categoryCommandRepository) DeleteAllPermanentCategories(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentCategories(ctx)

	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories
	}
	return true, nil
}
