package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderRecordMapping
}

func NewOrderQueryRepository(db *db.Queries, mapping recordmapper.OrderRecordMapping) *orderQueryRepository {
	return &orderQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderQueryRepository) FindAllOrders(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrders(ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindAllOrders
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordPagination(res), &totalCount, nil
}

func (r *orderQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersActive(ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordActivePagination(res), &totalCount, nil
}

func (r *orderQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordTrashedPagination(res), &totalCount, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*record.OrderRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersByMerchantParams{
		Column1: req.Search,
		Column4: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersByMerchant(ctx, reqDb)

	if err != nil {
		return nil, nil, order_errors.ErrFindByMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrdersRecordByMerchantPagination(res), &totalCount, nil
}

func (r *orderQueryRepository) FindById(ctx context.Context, user_id int) (*record.OrderRecord, error) {
	res, err := r.db.GetOrderByID(ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return r.mapping.ToOrderRecord(res), nil
}
