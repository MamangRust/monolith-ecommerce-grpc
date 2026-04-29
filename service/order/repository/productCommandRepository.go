package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	product_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type productCommandRepository struct {
	client pb.ProductCommandServiceClient
}

func NewProductCommandRepository(client pb.ProductCommandServiceClient) *productCommandRepository {
	return &productCommandRepository{
		client: client,
	}
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error) {
	res, err := r.client.UpdateProductCountStock(ctx, &pb.UpdateProductCountStockRequest{
		ProductId: int32(product_id),
		Stock:     int32(stock),
	})
	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &db.UpdateProductCountStockRow{
		ProductID:    res.Data.Id,
		CountInStock: res.Data.CountInStock,
	}, nil
}
