package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type categoryQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryQueryRepository(db *db.Queries, mapping recordmapper.CategoryRecordMapper) *categoryQueryRepository {
	return &categoryQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByID(ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrFindCategoryById
	}

	return r.mapping.ToCategoryRecord(res), nil
}
