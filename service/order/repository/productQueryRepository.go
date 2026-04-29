package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	product_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type productQueryRepository struct {
	client pb.ProductQueryServiceClient
}

func NewProductQueryRepository(client pb.ProductQueryServiceClient) *productQueryRepository {
	return &productQueryRepository{
		client: client,
	}
}

func (r *productQueryRepository) FindByID(ctx context.Context, product_id int) (*db.GetProductByIDRow, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdProductRequest{Id: int32(product_id)})
	if err != nil {
		return nil, product_errors.ErrProductInternal.WithInternal(err)
	}

	return &db.GetProductByIDRow{
		ProductID:   res.Data.Id,
		Name:        res.Data.Name,
		Description: &res.Data.Description,
		Price:       res.Data.Price,
		CountInStock: res.Data.CountInStock,
	}, nil
}
