package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	order_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type orderQueryRepository struct {
	client pb.OrderQueryServiceClient
}

func NewOrderQueryRepository(client pb.OrderQueryServiceClient) *orderQueryRepository {
	return &orderQueryRepository{
		client: client,
	}
}

func (r *orderQueryRepository) FindByID(ctx context.Context, order_id int) (*db.GetOrderByIDRow, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdOrderRequest{Id: int32(order_id)})
	if err != nil {
		return nil, order_errors.ErrFindById.WithInternal(err)
	}

	return &db.GetOrderByIDRow{
		OrderID:    res.Data.Id,
		UserID:     res.Data.UserId,
		MerchantID: res.Data.MerchantId,
		TotalPrice: res.Data.TotalPrice,
	}, nil
}
