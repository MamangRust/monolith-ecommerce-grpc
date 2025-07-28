package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type transactionCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionCommandRepository(db *db.Queries, mapping recordmapper.TransactionRecordMapping) *transactionCommandRepository {
	return &transactionCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *transactionCommandRepository) CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.CreateTransactionParams{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
	}

	transaction, err := r.db.CreateTransaction(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrCreateTransaction
	}

	return r.mapping.ToTransactionRecord(transaction), nil
}

func (r *transactionCommandRepository) UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error) {
	req := db.UpdateTransactionParams{
		TransactionID: int32(*request.TransactionID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		OrderID:       int32(request.OrderID),
		PaymentStatus: *request.PaymentStatus,
	}

	res, err := r.db.UpdateTransaction(ctx, req)

	if err != nil {
		return nil, transaction_errors.ErrUpdateTransaction
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) TrashTransaction(ctx context.Context, transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.TrashTransaction(ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrTrashTransaction
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) RestoreTransaction(ctx context.Context, transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.RestoreTransaction(ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrRestoreTransaction
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionCommandRepository) DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(ctx, int32(transaction_id))

	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently
	}

	return true, nil
}

func (r *transactionCommandRepository) RestoreAllTransactions(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransactions(ctx)

	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactions
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(ctx)

	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionPermanent
	}
	return true, nil
}
