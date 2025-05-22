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

type categoryStatsByMerchantRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryStatsByMerchantRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CategoryRecordMapper) *categoryStatsByMerchantRepository {
	return &categoryStatsByMerchantRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *categoryStatsByMerchantRepository) GetMonthlyTotalPriceByMerchant(req *requests.MonthTotalPriceMerchant) ([]*record.CategoriesMonthlyTotalPriceRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalPriceByMerchant(r.ctx, db.GetMonthlyTotalPriceByMerchantParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		MerchantID:  int32(req.MerchantID),
	})

	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceByMerchant
	}

	so := r.mapping.ToCategoryMonthlyTotalPricesByMerchant(res)

	return so, nil
}

func (r *categoryStatsByMerchantRepository) GetYearlyTotalPricesByMerchant(req *requests.YearTotalPriceMerchant) ([]*record.CategoriesYearlyTotalPriceRecord, error) {
	res, err := r.db.GetYearlyTotalPriceByMerchant(r.ctx, db.GetYearlyTotalPriceByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesByMerchant
	}

	so := r.mapping.ToCategoryYearlyTotalPricesByMerchant(res)

	return so, nil
}

func (r *categoryStatsByMerchantRepository) GetMonthPriceByMerchant(req *requests.MonthPriceMerchant) ([]*record.CategoriesMonthPriceRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategoryByMerchant(r.ctx, db.GetMonthlyCategoryByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, category_errors.ErrGetMonthPriceByMerchant
	}

	return r.mapping.ToCategoryMonthlyPricesByMerchant(res), nil
}

func (r *categoryStatsByMerchantRepository) GetYearPriceByMerchant(req *requests.YearPriceMerchant) ([]*record.CategoriesYearPriceRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategoryByMerchant(r.ctx, db.GetYearlyCategoryByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearPriceByMerchant
	}

	return r.mapping.ToCategoryYearlyPricesByMerchant(res), nil
}
