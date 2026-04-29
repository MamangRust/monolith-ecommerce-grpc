package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
)

type transactionStatsByMerchantRepository struct {
	db *db.Queries
}

func NewTransactionStatsByMerchantRepository(db *db.Queries) *transactionStatsByMerchantRepository {
	return &transactionStatsByMerchantRepository{
		db: db,
	}
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionSuccessByMerchantRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionSuccessByMerchant(ctx, db.GetMonthlyAmountTransactionSuccessByMerchantParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionSuccessByMerchantRow, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccessByMerchant(ctx, db.GetYearlyAmountTransactionSuccessByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*db.GetMonthlyAmountTransactionFailedByMerchantRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionFailedByMerchant(ctx, db.GetMonthlyAmountTransactionFailedByMerchantParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*db.GetYearlyAmountTransactionFailedByMerchantRow, error) {
	res, err := r.db.GetYearlyAmountTransactionFailedByMerchant(ctx, db.GetYearlyAmountTransactionFailedByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsByMerchantSuccess(ctx, db.GetMonthlyTransactionMethodsByMerchantSuccessParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantSuccessRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsByMerchantSuccess(ctx, db.GetYearlyTransactionMethodsByMerchantSuccessParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*db.GetMonthlyTransactionMethodsByMerchantFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsByMerchantFailed(ctx, db.GetMonthlyTransactionMethodsByMerchantFailedParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*db.GetYearlyTransactionMethodsByMerchantFailedRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsByMerchantFailed(ctx, db.GetYearlyTransactionMethodsByMerchantFailedParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}

	return res, nil
}

