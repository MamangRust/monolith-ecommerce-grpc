package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type CategoryQueryRepository interface {
	FindByID(ctx context.Context, category_id int) (*db.GetCategoryByIDRow, error)
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}

type ProductQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, error)
	FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, error)
	FindByID(ctx context.Context, product_id int) (*db.GetProductByIDRow, error)
	FindByIDTrashed(ctx context.Context, product_id int) (*db.Product, error)
}

type ProductCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateProductRequest) (*db.CreateProductRow, error)
	Update(ctx context.Context, request *requests.UpdateProductRequest) (*db.UpdateProductRow, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*db.UpdateProductCountStockRow, error)
	Trash(ctx context.Context, product_id int) (*db.TrashProductRow, error)
	Restore(ctx context.Context, product_id int) (*db.RestoreProductRow, error)
	DeletePermanent(
		ctx context.Context,
		product_id int,
	) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
