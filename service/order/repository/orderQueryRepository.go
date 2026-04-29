package repository

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
)

type orderQueryRepository struct {
	db *db.Queries
}

func NewOrderQueryRepository(db *db.Queries) OrderQueryRepository {
	return &orderQueryRepository{
		db: db,
	}
}

func (r *orderQueryRepository) FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrders(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindAllOrders.WithInternal(err)
	}

	return res, nil
}

func (r *orderQueryRepository) FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersActive(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindByActive.WithInternal(err)
	}

	return res, nil
}

func (r *orderQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersTrashed(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}

	return res, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*db.GetOrdersByMerchantRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetOrdersByMerchantParams{
		Column1: req.Search,
		Column4: int32(req.MerchantID),
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetOrdersByMerchant(ctx, reqDb)

	if err != nil {
		return nil, order_errors.ErrFindByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *orderQueryRepository) FindByID(ctx context.Context, id int) (*db.GetOrderByIDRow, error) {
	res, err := r.db.GetOrderByID(ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, order_errors.ErrOrderNotFound.WithInternal(err)
		}
		return nil, order_errors.ErrFindById.WithInternal(err)
	}

	return res, nil
}

