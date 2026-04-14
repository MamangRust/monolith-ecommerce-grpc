package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryQueryRepository interface {
	FindById(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, error)
	FindById(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
	FindByIdTrashed(ctx context.Context, product_id int) (*db.Product, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*db.CreateProductRow, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*db.UpdateProductRow, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error)
	TrashedProduct(ctx context.Context, product_id int) (*db.TrashProductRow, error)
	RestoreProduct(ctx context.Context, product_id int) (*db.RestoreProductRow, error)
	DeleteProductPermanent(
		ctx context.Context,
		product_id int,
	) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
