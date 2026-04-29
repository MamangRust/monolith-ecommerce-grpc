package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type productQueryRepository struct {
	client pb.ProductQueryServiceClient
}

func NewProductQueryRepository(client pb.ProductQueryServiceClient) ProductQueryRepository {
	return &productQueryRepository{
		client: client,
	}
}

func (r *productQueryRepository) FindById(ctx context.Context, id int) (*db.GetProductByIDRow, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdProductRequest{Id: int32(id)})
	if err != nil {
		return nil, product_errors.ErrProductNotFound.WithInternal(err)
	}

	rating := float64(res.Data.Rating)
	
	return &db.GetProductByIDRow{
		ProductID:    res.Data.Id,
		MerchantID:   res.Data.MerchantId,
		CategoryID:   res.Data.CategoryId,
		Name:         res.Data.Name,
		Description:  &res.Data.Description,
		Price:        res.Data.Price,
		CountInStock: res.Data.CountInStock,
		Brand:        &res.Data.Brand,
		Weight:       &res.Data.Weight,
		Rating:       &rating,
		SlugProduct:  &res.Data.SlugProduct,
		ImageProduct: &res.Data.ImageProduct,
	}, nil
}
