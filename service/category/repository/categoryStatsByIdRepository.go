package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type categoryStatsByIdRepository struct {
	db *db.Queries
}

func NewCategoryStatsByIdRepository(db *db.Queries) *categoryStatsByIdRepository {
	return &categoryStatsByIdRepository{
		db: db,
	}
}

func (r *categoryStatsByIdRepository) GetMonthlyTotalPriceById(
	ctx context.Context,
	req *requests.MonthTotalPriceCategory,
) ([]*db.GetMonthlyTotalPriceByIdRow, error) {

	currentMonthStart := time.Date(
		req.Year,
		time.Month(req.Month),
		1, 0, 0, 0, 0,
		time.UTC,
	)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	params := db.GetMonthlyTotalPriceByIdParams{
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
		CategoryID: int32(req.CategoryID),
	}

	res, err := r.db.GetMonthlyTotalPriceById(ctx, params)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceById.WithInternal(err)
	}


	return res, nil
}

func (r *categoryStatsByIdRepository) GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*db.GetYearlyTotalPriceByIdRow, error) {
	res, err := r.db.GetYearlyTotalPriceById(ctx, db.GetYearlyTotalPriceByIdParams{
		Column1:    int32(req.Year),
		CategoryID: int32(req.CategoryID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesById.WithInternal(err)
	}


	return res, nil
}

func (r *categoryStatsByIdRepository) GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*db.GetMonthlyCategoryByIdRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategoryById(ctx, db.GetMonthlyCategoryByIdParams{
		Column1:    yearStart,
		CategoryID: int32(req.CategoryID),
	})
	if err != nil {
		return nil, category_errors.ErrGetMonthPriceById.WithInternal(err)
	}


	return res, nil
}

func (r *categoryStatsByIdRepository) GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*db.GetYearlyCategoryByIdRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategoryById(ctx, db.GetYearlyCategoryByIdParams{
		Column1:    yearStart,
		CategoryID: int32(req.CategoryID),
	})

	if err != nil {
		return nil, category_errors.ErrGetYearPriceById.WithInternal(err)
	}


	return res, nil
}
