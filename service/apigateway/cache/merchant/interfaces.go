package merchant_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchant, bool)
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchant)

	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDeleteAt, bool)
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDeleteAt)

	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDeleteAt, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDeleteAt)

	GetCachedMerchant(ctx context.Context, id int) (*response.ApiResponseMerchant, bool)
	SetCachedMerchant(ctx context.Context, data *response.ApiResponseMerchant)

	GetCachedMerchantsByUserId(ctx context.Context, id int) (*response.ApiResponsePaginationMerchant, bool)
	SetCachedMerchantsByUserId(ctx context.Context, userId int, data *response.ApiResponsePaginationMerchant)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
}
