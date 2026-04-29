package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsRow, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data []*db.GetProductsRow, total *int)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*db.GetProductsByMerchantRow, *int, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant, data []*db.GetProductsByMerchantRow, total *int)

	GetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*db.GetProductsByCategoryNameRow, *int, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory, data []*db.GetProductsByCategoryNameRow, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsActiveRow, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data []*db.GetProductsActiveRow, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*db.GetProductsTrashedRow, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data []*db.GetProductsTrashedRow, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*db.GetProductByIDRow, bool)
	SetCachedProduct(ctx context.Context, data *db.GetProductByIDRow)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
}
