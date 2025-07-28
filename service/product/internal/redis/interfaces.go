package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponse, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponse, total *int)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant, data []*response.ProductResponse, total *int)

	GetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory, data []*response.ProductResponse, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*response.ProductResponse, bool)
	SetCachedProduct(ctx context.Context, data *response.ProductResponse)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
}
