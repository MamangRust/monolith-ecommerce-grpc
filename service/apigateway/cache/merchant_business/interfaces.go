package merchantbusiness_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantBusinessQueryCache interface {
	GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantBusiness, bool)
	SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantBusiness)

	GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantBusinessDeleteAt, bool)
	SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantBusinessDeleteAt)

	GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantBusinessDeleteAt, bool)
	SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantBusinessDeleteAt)

	GetCachedMerchantBusiness(ctx context.Context, id int) (*response.ApiResponseMerchantBusiness, bool)
	SetCachedMerchantBusiness(ctx context.Context, data *response.ApiResponseMerchantBusiness)
}

type MerchantBusinessCommandCache interface {
	DeleteMerchantBusinessCache(ctx context.Context, id int)
}
