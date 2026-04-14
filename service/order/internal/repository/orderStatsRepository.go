package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type orderStatsRepository struct {
	db *db.Queries
}

func NewOrderStatsRepository(db *db.Queries) OrderStatsRepository {
	return &orderStatsRepository{
		db: db,
	}
}

func (r *orderStatsRepository) GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*db.GetMonthlyTotalRevenueRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenue(ctx, db.GetMonthlyTotalRevenueParams{
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
		return nil, order_errors.ErrGetMonthlyTotalRevenue.WithInternal(err)
	}

	return res, nil
}

func (r *orderStatsRepository) GetYearlyTotalRevenue(ctx context.Context, year int) ([]*db.GetYearlyTotalRevenueRow, error) {
	res, err := r.db.GetYearlyTotalRevenue(ctx, int32(year))

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenue.WithInternal(err)
	}

	return res, nil
}

func (r *orderStatsRepository) GetMonthlyOrder(ctx context.Context, year int) ([]*db.GetMonthlyOrderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyOrder(ctx, yearStart)

	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrder.WithInternal(err)
	}

	return res, nil
}

func (r *orderStatsRepository) GetYearlyOrder(ctx context.Context, year int) ([]*db.GetYearlyOrderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrder(ctx, yearStart)
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrder.WithInternal(err)
	}

	return res, nil
}

