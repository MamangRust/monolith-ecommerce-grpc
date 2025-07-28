package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type transactonStatsRepository struct {
	db      *db.Queries
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionStatsRepository(db *db.Queries, mapping recordmapper.TransactionRecordMapping) *transactonStatsRepository {
	return &transactonStatsRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *transactonStatsRepository) GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountSuccessRecord, error) {
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
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess
	}

	return r.mapping.ToTransactionMonthlyAmountSuccess(res), nil
}

func (r *transactonStatsRepository) GetYearlyAmountSuccess(ctx context.Context, year int) ([]*record.TransactionYearlyAmountSuccessRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccess(ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess
	}

	return r.mapping.ToTransactionYearlyAmountSuccess(res), nil
}

func (r *transactonStatsRepository) GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountFailedRecord, error) {
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
		return nil, transaction_errors.ErrGetMonthlyAmountFailed
	}

	return r.mapping.ToTransactionMonthlyAmountFailed(res), nil
}

func (r *transactonStatsRepository) GetYearlyAmountFailed(ctx context.Context, year int) ([]*record.TransactionYearlyAmountFailedRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionFailed(ctx, int32(year))

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed
	}

	return r.mapping.ToTransactionYearlyAmountFailed(res), nil
}

func (r *transactonStatsRepository) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error) {
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
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	return r.mapping.ToTransactionMonthlyMethodSuccess(res), nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsSuccess(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	return r.mapping.ToTransactionYearlyMethodSuccess(res), nil
}

func (r *transactonStatsRepository) GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error) {
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
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	return r.mapping.ToTransactionMonthlyMethodFailed(res), nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsFailed(ctx, yearStart)

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	return r.mapping.ToTransactionYearlyMethodFailed(res), nil
}
