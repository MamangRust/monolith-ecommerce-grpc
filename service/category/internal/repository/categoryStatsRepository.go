package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type categoryStatsRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryStatsRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CategoryRecordMapper) *categoryStatsRepository {
	return &categoryStatsRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *categoryStatsRepository) GetMonthlyTotalPrice(req *requests.MonthTotalPrice) ([]*record.CategoriesMonthlyTotalPriceRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalPrice(r.ctx, db.GetMonthlyTotalPriceParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
	})

	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPrice
	}

	so := r.mapping.ToCategoryMonthlyTotalPrices(res)

	return so, nil
}

func (r *categoryStatsRepository) GetYearlyTotalPrices(year int) ([]*record.CategoriesYearlyTotalPriceRecord, error) {
	res, err := r.db.GetYearlyTotalPrice(r.ctx, int32(year))

	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPrices
	}

	so := r.mapping.ToCategoryYearlyTotalPrices(res)

	return so, nil
}

func (r *categoryStatsRepository) GetMonthPrice(year int) ([]*record.CategoriesMonthPriceRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategory(r.ctx, yearStart)

	if err != nil {
		return nil, category_errors.ErrGetMonthPrice
	}

	return r.mapping.ToCategoryMonthlyPrices(res), nil
}

func (r *categoryStatsRepository) GetYearPrice(year int) ([]*record.CategoriesYearPriceRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategory(r.ctx, yearStart)

	if err != nil {
		return nil, category_errors.ErrGetYearPrice
	}

	return r.mapping.ToCategoryYearlyPrices(res), nil
}
