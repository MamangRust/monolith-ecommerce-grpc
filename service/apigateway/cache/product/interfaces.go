package product_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProduct, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProduct)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) (*response.ApiResponsePaginationProduct, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant, data *response.ApiResponsePaginationProduct)

	GetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory) (*response.ApiResponsePaginationProduct, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.FindAllProductByCategory, data *response.ApiResponsePaginationProduct)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProductDeleteAt, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProductDeleteAt)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct) (*response.ApiResponsePaginationProductDeleteAt, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProduct, data *response.ApiResponsePaginationProductDeleteAt)

	GetCachedProduct(ctx context.Context, productID int) (*response.ApiResponseProduct, bool)
	SetCachedProduct(ctx context.Context, data *response.ApiResponseProduct)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
}
