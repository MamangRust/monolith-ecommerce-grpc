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
	ctx     context.Context
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CategoryRecordMapper) *categoryCommandRepository {
	return &categoryCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *categoryCommandRepository) CreateCategory(request *requests.CreateCategoryRequest) (*record.CategoriesRecord, error) {
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

	category, err := r.db.CreateCategory(r.ctx, req)
	if err != nil {
		return nil, category_errors.ErrCreateCategory
	}

	return r.mapping.ToCategoryRecord(category), nil
}

func (r *categoryCommandRepository) UpdateCategory(request *requests.UpdateCategoryRequest) (*record.CategoriesRecord, error) {
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

	res, err := r.db.UpdateCategory(r.ctx, req)

	if err != nil {
		return nil, category_errors.ErrUpdateCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryCommandRepository) TrashedCategory(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.TrashCategory(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrTrashedCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryCommandRepository) RestoreCategory(category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.RestoreCategory(r.ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrRestoreCategory
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryCommandRepository) DeleteCategoryPermanently(category_id int) (bool, error) {
	err := r.db.DeleteCategoryPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently
	}

	return true, nil
}

func (r *categoryCommandRepository) RestoreAllCategories() (bool, error) {
	err := r.db.RestoreAllCategories(r.ctx)

	if err != nil {
		return false, category_errors.ErrRestoreAllCategories
	}
	return true, nil
}

func (r *categoryCommandRepository) DeleteAllPermanentCategories() (bool, error) {
	err := r.db.DeleteAllPermanentCategories(r.ctx)

	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories
	}
	return true, nil
}
