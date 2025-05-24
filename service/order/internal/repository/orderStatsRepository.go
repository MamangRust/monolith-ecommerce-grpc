package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type orderStatsRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.OrderRecordMapping
}

func NewOrderStatsRepository(db *db.Queries, ctx context.Context, mapping recordmapper.OrderRecordMapping) *orderStatsRepository {
	return &orderStatsRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *orderStatsRepository) GetMonthlyTotalRevenue(req *requests.MonthTotalRevenue) ([]*record.OrderMonthlyTotalRevenueRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenue(r.ctx, db.GetMonthlyTotalRevenueParams{
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

func (r *orderStatsRepository) GetYearlyTotalRevenue(year int) ([]*record.OrderYearlyTotalRevenueRecord, error) {
	res, err := r.db.GetYearlyTotalRevenue(r.ctx, int32(year))

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenue
	}

	fmt.Println("err", err)

	so := r.mapping.ToOrderYearlyTotalRevenues(res)

	return so, nil
}

func (r *orderStatsRepository) GetMonthlyOrder(year int) ([]*record.OrderMonthlyRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	res, err := r.db.GetMonthlyOrder(r.ctx, yearStart)

	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrder
	}

	return r.mapping.ToOrderMonthlyPrices(res), nil
}

func (r *orderStatsRepository) GetYearlyOrder(year int) ([]*record.OrderYearlyRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrder(r.ctx, yearStart)
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrder
	}

	return r.mapping.ToOrderYearlyPrices(res), nil
}
