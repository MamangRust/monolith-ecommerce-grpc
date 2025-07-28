package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type transactionQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.TransactionRecordMapping
}

func NewTransactionQueryRepository(db *db.Queries, mapping recordmapper.TransactionRecordMapping) *transactionQueryRepository {
	return &transactionQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *transactionQueryRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactions(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordPagination(res), &totalCount, nil
}

func (r *transactionQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsActive(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordActivePagination(res), &totalCount, nil
}

func (r *transactionQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionsTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionsRecordTrashedPagination(res), &totalCount, nil
}

func (r *transactionQueryRepository) FindByMerchant(
	ctx context.Context,
	req *requests.FindAllTransactionByMerchant,
) ([]*record.TransactionRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTransactionByMerchantParams{
		Column1: req.Search,
		Column2: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTransactionByMerchant(ctx, reqDb)

	if err != nil {
		return nil, nil, transaction_errors.ErrFindByMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToTransactionMerchantsRecordPagination(res), &totalCount, nil
}

func (r *transactionQueryRepository) FindById(ctx context.Context, transaction_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByID(ctx, int32(transaction_id))

	if err != nil {
		return nil, transaction_errors.ErrFindById
	}

	return r.mapping.ToTransactionRecord(res), nil
}

func (r *transactionQueryRepository) FindByOrderId(ctx context.Context, order_id int) (*record.TransactionRecord, error) {
	res, err := r.db.GetTransactionByOrderID(ctx, int32(order_id))

	if err != nil {
		return nil, transaction_errors.ErrFindByOrderId
	}

	return r.mapping.ToTransactionRecord(res), nil
}
