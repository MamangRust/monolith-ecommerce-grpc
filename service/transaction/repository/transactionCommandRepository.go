package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
)

type transactionCommandRepository struct {
	db *db.Queries
}

func NewTransactionCommandRepository(db *db.Queries) *transactionCommandRepository {
	return &transactionCommandRepository{
		db: db,
	}
}

func transactionCreateParams(request *requests.CreateTransactionRequest) db.CreateTransactionParams {
	return db.CreateTransactionParams{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
	}
}

func (r *transactionCommandRepository) Create(ctx context.Context, request *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error) {
	res, err := r.db.CreateTransaction(ctx, transactionCreateParams(request))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transaction_errors.ErrTransactionNotFound
		}
		return nil, transaction_errors.ErrCreateTransaction.WithInternal(err)
	}

	return res, nil
}

// CreateInTx runs the transaction insert inside the given database transaction so
// the caller can commit the business row and its outbox event atomically.
func (r *transactionCommandRepository) CreateInTx(ctx context.Context, tx pgx.Tx, request *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error) {
	res, err := r.db.WithTx(tx).CreateTransaction(ctx, transactionCreateParams(request))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transaction_errors.ErrTransactionNotFound
		}
		return nil, transaction_errors.ErrCreateTransaction.WithInternal(err)
	}

	return res, nil
}

func (r *transactionCommandRepository) Update(ctx context.Context, request *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transaction_errors.ErrTransactionNotFound
		}
		return nil, transaction_errors.ErrUpdateTransaction.WithInternal(err)
	}

	return res, nil
}

func (r *transactionCommandRepository) Trash(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	res, err := r.db.TrashTransaction(ctx, int32(transaction_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transaction_errors.ErrTransactionNotFound
		}
		return nil, transaction_errors.ErrTrashTransaction.WithInternal(err)
	}

	return res, nil
}

func (r *transactionCommandRepository) Restore(ctx context.Context, transaction_id int) (*db.Transaction, error) {
	res, err := r.db.RestoreTransaction(ctx, int32(transaction_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transaction_errors.ErrTransactionNotFound
		}
		return nil, transaction_errors.ErrRestoreTransaction.WithInternal(err)
	}

	return res, nil
}

func (r *transactionCommandRepository) DeletePermanent(ctx context.Context, transaction_id int) (bool, error) {
	err := r.db.DeleteTransactionPermanently(ctx, int32(transaction_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, transaction_errors.ErrTransactionNotFound
		}
		return false, transaction_errors.ErrDeleteTransactionPermanently.WithInternal(err)
	}

	return true, nil
}

func (r *transactionCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error) {
	err := r.db.DeleteTransactionByOrderPermanent(ctx, int32(order_id))

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, transaction_errors.ErrTransactionNotFound
		}
		return false, transaction_errors.ErrDeleteTransactionPermanently.WithInternal(err)
	}

	return true, nil
}

func (r *transactionCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllTransactions(ctx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, transaction_errors.ErrTransactionNotFound
		}
		return false, transaction_errors.ErrRestoreAllTransactions.WithInternal(err)
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentTransactions(ctx)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, transaction_errors.ErrTransactionNotFound
		}
		return false, transaction_errors.ErrDeleteAllTransactionPermanent.WithInternal(err)
	}
	return true, nil
}
