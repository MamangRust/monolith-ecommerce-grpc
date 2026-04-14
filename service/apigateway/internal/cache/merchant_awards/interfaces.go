package merchantawards_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantAwardQueryCache interface {
	GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantAward, bool)
	SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantAward)

	GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantAwardDeleteAt, bool)
	SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantAwardDeleteAt)

	GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantAwardDeleteAt, bool)
	SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantAwardDeleteAt)

	GetCachedMerchantAward(ctx context.Context, id int) (*response.ApiResponseMerchantAward, bool)
	SetCachedMerchantAward(ctx context.Context, data *response.ApiResponseMerchantAward)
}

type MerchantAwardCommandCache interface {
	DeleteMerchantAwardCache(ctx context.Context, id int)
}
