package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderItemQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderItemRecordMapping
}

func NewOrderItemQueryRepository(db *db.Queries, mapping recordmapper.OrderItemRecordMapping) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderItemQueryRepository) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItems(ctx, reqDb)

	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordPagination(res), &totalCount, nil
}

func (r *orderItemQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsActive(ctx, reqDb)

	if err != nil {
		return nil, nil, orderitem_errors.ErrFindByActive
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordActivePagination(res), &totalCount, nil
}

func (r *orderItemQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, orderitem_errors.ErrFindByTrashed
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToOrderItemsRecordTrashedPagination(res), &totalCount, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*record.OrderItemRecord, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	return r.mapping.ToOrderItemsRecord(res), nil
}
