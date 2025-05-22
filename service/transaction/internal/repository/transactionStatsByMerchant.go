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

type transactionStatsByMerchantRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionStatsByMerchantRepository(db *db.Queries, ctx context.Context, mapping recordmapper.TransactionRecordMapping) *transactionStatsByMerchantRepository {
	return &transactionStatsByMerchantRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountSuccessRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionSuccessByMerchant(r.ctx, db.GetMonthlyAmountTransactionSuccessByMerchantParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant
	}

	return r.mapping.ToTransactionMonthlyAmountSuccessByMerchant(res), nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountSuccessRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionSuccessByMerchant(r.ctx, db.GetYearlyAmountTransactionSuccessByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant
	}

	return r.mapping.ToTransactionYearlyAmountSuccessByMerchant(res), nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountFailedRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyAmountTransactionFailedByMerchant(r.ctx, db.GetMonthlyAmountTransactionFailedByMerchantParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant
	}

	return r.mapping.ToTransactionMonthlyAmountFailedByMerchant(res), nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountFailedRecord, error) {
	res, err := r.db.GetYearlyAmountTransactionFailedByMerchant(r.ctx, db.GetYearlyAmountTransactionFailedByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant
	}

	return r.mapping.ToTransactionYearlyAmountFailedByMerchant(res), nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsByMerchantSuccess(r.ctx, db.GetMonthlyTransactionMethodsByMerchantSuccessParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionMonthlyByMerchantMethodSuccess(res), nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsByMerchantSuccess(r.ctx, db.GetYearlyTransactionMethodsByMerchantSuccessParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionYearlyMethodByMerchantSuccess(res), nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevDate := currentDate.AddDate(0, -1, 0)

	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)
	lastDayPrevMonth := prevDate.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTransactionMethodsByMerchantFailed(r.ctx, db.GetMonthlyTransactionMethodsByMerchantFailedParams{
		Column1:    currentDate,
		Column2:    lastDayCurrentMonth,
		Column3:    prevDate,
		Column4:    lastDayPrevMonth,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionMonthlyByMerchantMethodFailed(res), nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyTransactionMethodsByMerchantFailed(r.ctx, db.GetYearlyTransactionMethodsByMerchantFailedParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}

	return r.mapping.ToTransactionYearlyMethodByMerchantFailed(res), nil
}
