package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
)

type transactionQueryRepository struct {
	db *db.Queries
}

func NewTransactionQueryRepository(db *db.Queries) *transactionQueryRepository {
	return &transactionQueryRepository{
		db: db,
	}
}

func (r *transactionQueryRepository) FindAll(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactions(ctx, reqDb)
	if err != nil {
		return nil, transaction_errors.ErrFindAllTransactions.WithInternal(err)
	}

	return res, nil
}

func (r *transactionQueryRepository) FindActive(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsActive(ctx, reqDb)
	if err != nil {
		return nil, transaction_errors.ErrFindByActive.WithInternal(err)
	}

	return res, nil
}

func (r *transactionQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*db.GetTransactionsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsTrashed(ctx, reqDb)
	if err != nil {
		return nil, transaction_errors.ErrFindByTrashed.WithInternal(err)
	}

	return res, nil
}

func (r *transactionQueryRepository) FindByMerchant(
	ctx context.Context,
	req *requests.FindAllTransactionByMerchant,
) ([]*db.GetTransactionByMerchantRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionByMerchantParams{
		Column1: req.Search,
		Column2: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionByMerchant(ctx, reqDb)
	if err != nil {
		return nil, transaction_errors.ErrFindByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *transactionQueryRepository) FindByID(ctx context.Context, transaction_id int) (*db.GetTransactionByIDRow, error) {
	res, err := r.db.GetTransactionByID(ctx, int32(transaction_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, transaction_errors.ErrTransactionNotFound.WithInternal(err)
		}
		return nil, transaction_errors.ErrFindById.WithInternal(err)
	}

	return res, nil
}

func (r *transactionQueryRepository) FindByOrderID(ctx context.Context, order_id int) (*db.GetTransactionByOrderIDRow, error) {
	res, err := r.db.GetTransactionByOrderID(ctx, int32(order_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, transaction_errors.ErrTransactionNotFound.WithInternal(err)
		}
		return nil, transaction_errors.ErrFindByOrderId.WithInternal(err)
	}

	return res, nil
}

