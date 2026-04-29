package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
)

type transactionStatsRepository struct {
	db *db.Queries
}

func NewTransactionStatsRepository(db *db.Queries) *transactionStatsRepository {
	return &transactionStatsRepository{
		db: db,
	}
}

func (r *transactionStatsRepository) GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionSuccess(ctx, db.GetMonthlyAmountTransactionSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsRepository) GetYearlyAmountSuccess(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionSuccessRow, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccess(ctx, int32(year))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess.WithInternal(err)
	}
	return res, nil
}

func (r *transactionStatsRepository) GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*db.GetMonthlyAmountTransactionFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionFailed(ctx, db.GetMonthlyAmountTransactionFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailed.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsRepository) GetYearlyAmountFailed(ctx context.Context, year int) ([]*db.GetYearlyAmountTransactionFailedRow, error) {
	res, err := r.db.GetYearlyAmountTransactionFailed(ctx, int32(year))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed.WithInternal(err)
	}
	return res, nil
}

func (r *transactionStatsRepository) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsSuccess(ctx, db.GetMonthlyTransactionMethodsSuccessParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsRepository) GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsSuccessRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsSuccess(ctx, yearStart)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	return res, nil
}

func (r *transactionStatsRepository) GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*db.GetMonthlyTransactionMethodsFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsFailed(ctx, db.GetMonthlyTransactionMethodsFailedParams{
		Column1: currentDate,
		Column2: lastDayCurrentMonth,
		Column3: prevDate,
		Column4: lastDayPrevMonth,
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsRepository) GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*db.GetYearlyTransactionMethodsFailedRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsFailed(ctx, yearStart)
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	return res, nil
}

