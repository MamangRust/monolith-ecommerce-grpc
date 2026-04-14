package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ProductQueryService interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, *int, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, *int, error)
	FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
}

type ProductCommandService interface {
	CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*db.CreateProductRow, error)
	UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*db.UpdateProductRow, error)
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*db.UpdateProductCountStockRow, error)
	TrashedProduct(ctx context.Context, productID int) (interface{}, error)
	RestoreProduct(ctx context.Context, productID int) (interface{}, error)
	DeleteProductPermanent(ctx context.Context, productID int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
