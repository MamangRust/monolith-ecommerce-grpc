package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type categoryStatsRepository struct {
	db *db.Queries
}

func NewCategoryStatsRepository(db *db.Queries) *categoryStatsRepository {
	return &categoryStatsRepository{
		db: db,
	}
}

func (r *categoryStatsRepository) GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*db.GetMonthlyTotalPriceRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalPrice(ctx, db.GetMonthlyTotalPriceParams{
		Extract: pgtype.Date{
			Time:  currentMonthStart,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  currentMonthEnd,
			Valid: true,
		},
		CreatedAt_2: pgtype.Timestamp{
			Time:  prevMonthStart,
			Valid: true,
		},
		CreatedAt_3: pgtype.Timestamp{
			Time:  prevMonthEnd,
			Valid: true,
		},
	})

	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPrice.WithInternal(err)
	}


	return res, nil
}

func (r *categoryStatsRepository) GetYearlyTotalPrices(ctx context.Context, year int) ([]*db.GetYearlyTotalPriceRow, error) {
	res, err := r.db.GetYearlyTotalPrice(ctx, int32(year))

	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPrices.WithInternal(err)
	}


	return res, nil
}

func (r *categoryStatsRepository) GetMonthPrice(ctx context.Context, year int) ([]*db.GetMonthlyCategoryRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategory(ctx, yearStart)

	if err != nil {
		return nil, category_errors.ErrGetMonthPrice.WithInternal(err)
	}


	return res, nil
}

func (r *categoryStatsRepository) GetYearPrice(ctx context.Context, year int) ([]*db.GetYearlyCategoryRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategory(ctx, yearStart)

	if err != nil {
		return nil, category_errors.ErrGetYearPrice.WithInternal(err)
	}


	return res, nil
}
