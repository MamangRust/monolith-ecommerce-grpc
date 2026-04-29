package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type categoryQueryRepository struct {
	client pb.CategoryQueryServiceClient
}

func NewCategoryQueryRepository(client pb.CategoryQueryServiceClient) *categoryQueryRepository {
	return &categoryQueryRepository{
		client: client,
	}
}

func (r *categoryQueryRepository) FindByID(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error) {
	res, err := r.client.FindById(ctx, &pb.FindByIdCategoryRequest{Id: int32(category_id)})
	if err != nil {
		return nil, category_errors.ErrFindCategoryById.WithInternal(err)
	}

	return &db.GetCategoryByIDRow{
		CategoryID:  res.Data.Id,
		Name:        res.Data.Name,
		Description: &res.Data.Description,
	}, nil
}
