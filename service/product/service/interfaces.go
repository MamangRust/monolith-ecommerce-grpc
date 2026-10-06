package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, *int, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, *int, error)
	FindByID(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
}

type ProductCommandService interface {
	Create(ctx context.Context, req *requests.CreateProductRequest) (*db.CreateProductRow, error)
	Update(ctx context.Context, req *requests.UpdateProductRequest) (*db.UpdateProductRow, error)
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*db.UpdateProductCountStockRow, error)
	Trash(ctx context.Context, productID int) (interface{}, error)
	Restore(ctx context.Context, productID int) (interface{}, error)
	DeletePermanent(ctx context.Context, productID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
