package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderStatsRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderRecordMapping
}

func NewOrderStatsRepository(db *db.Queries, mapping recordmapper.OrderRecordMapping) *orderStatsRepository {
	return &orderStatsRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderStatsRepository) GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*record.OrderMonthlyTotalRevenueRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenue(ctx, db.GetMonthlyTotalRevenueParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenue
	}

	so := r.mapping.ToOrderMonthlyTotalRevenues(res)

	return so, nil
}

func (r *orderStatsRepository) GetYearlyTotalRevenue(ctx context.Context, year int) ([]*record.OrderYearlyTotalRevenueRecord, error) {
	res, err := r.db.GetYearlyTotalRevenue(ctx, int32(year))

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenue
	}

	fmt.Println("err", err)

	so := r.mapping.ToOrderYearlyTotalRevenues(res)

	return so, nil
}

func (r *orderStatsRepository) GetMonthlyOrder(ctx context.Context, year int) ([]*record.OrderMonthlyRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyOrder(ctx, yearStart)

	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrder
	}

	return r.mapping.ToOrderMonthlyPrices(res), nil
}

func (r *orderStatsRepository) GetYearlyOrder(ctx context.Context, year int) ([]*record.OrderYearlyRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrder(ctx, yearStart)
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrder
	}

	return r.mapping.ToOrderYearlyPrices(res), nil
}
