package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type orderStatsByMerchantRepository struct {
	db *db.Queries
}

func NewOrderStatsByMerchantRepository(db *db.Queries) OrderStatsByMerchantRepository {
	return &orderStatsByMerchantRepository{
		db: db,
	}
}

func (r *orderStatsByMerchantRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*db.GetMonthlyTotalRevenueByMerchantRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenueByMerchant(ctx, db.GetMonthlyTotalRevenueByMerchantParams{
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
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*db.GetYearlyTotalRevenueByMerchantRow, error) {
	res, err := r.db.GetYearlyTotalRevenueByMerchant(ctx, db.GetYearlyTotalRevenueByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *orderStatsByMerchantRepository) GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*db.GetMonthlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyOrderByMerchant(ctx, db.GetMonthlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*db.GetYearlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrderByMerchant(ctx, db.GetYearlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant.WithInternal(err)
	}

	return res, nil
}

