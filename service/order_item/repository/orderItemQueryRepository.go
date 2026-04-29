package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
)

type orderItemQueryRepository struct {
	db *db.Queries
}

func NewOrderItemQueryRepository(db *db.Queries) *orderItemQueryRepository {
	return &orderItemQueryRepository{
		db: db,
	}
}

func (r *orderItemQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItems(ctx, reqDb)

	if err != nil {
		return nil, orderitem_errors.ErrFindAllOrderItems
	}

	return res, nil
}

func (r *orderItemQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsActive(ctx, reqDb)

	if err != nil {
		return nil, orderitem_errors.ErrFindByActive
	}

	return res, nil
}

func (r *orderItemQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrderItemsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrderItemsTrashed(ctx, reqDb)

	if err != nil {
		return nil, orderitem_errors.ErrFindByTrashed
	}

	return res, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.GetOrderItemsByOrderRow, error) {
	res, err := r.db.GetOrderItemsByOrder(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	return res, nil
}
