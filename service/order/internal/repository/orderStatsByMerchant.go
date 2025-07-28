package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderStatsByMerchantRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderRecordMapping
}

func NewOrderStatsByMerchantRepository(db *db.Queries, mapping recordmapper.OrderRecordMapping) *orderStatsByMerchantRepository {
	return &orderStatsByMerchantRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderStatsByMerchantRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*record.OrderMonthlyTotalRevenueRecord, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalRevenueByMerchant(ctx, db.GetMonthlyTotalRevenueByMerchantParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		MerchantID:  int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant
	}

	so := r.mapping.ToOrderMonthlyTotalRevenuesByMerchant(res)

	return so, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*record.OrderYearlyTotalRevenueRecord, error) {
	res, err := r.db.GetYearlyTotalRevenueByMerchant(ctx, db.GetYearlyTotalRevenueByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant
	}

	so := r.mapping.ToOrderYearlyTotalRevenuesByMerchant(res)

	return so, nil
}

func (r *orderStatsByMerchantRepository) GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*record.OrderMonthlyRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyOrderByMerchant(ctx, db.GetMonthlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant
	}

	return r.mapping.ToOrderMonthlyPricesByMerchant(res), nil
}

func (r *orderStatsByMerchantRepository) GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*record.OrderYearlyRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyOrderByMerchant(ctx, db.GetYearlyOrderByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant
	}

	return r.mapping.ToOrderYearlyPricesByMerchant(res), nil
}
