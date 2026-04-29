package merchantdetail_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantDetailQueryCache interface {
	GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetail, bool)
	SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetail)

	GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetailDeleteAt, bool)
	SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetailDeleteAt)

	GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantDetailDeleteAt, bool)
	SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantDetailDeleteAt)

	GetCachedMerchantDetail(ctx context.Context, id int) (*response.ApiResponseMerchantDetail, bool)
	SetCachedMerchantDetail(ctx context.Context, data *response.ApiResponseMerchantDetail)

	GetCachedMerchantDetailRelation(
		ctx context.Context,
		merchantID int,
	) (*response.ApiResponseMerchantDetailRelation, bool)

	SetCachedMerchantDetailRelation(
		ctx context.Context,
		merchantID int,
		data *response.ApiResponseMerchantDetailRelation,
	)
}

type MerchantDetailCommandCache interface {
	DeleteMerchantDetailCache(ctx context.Context, id int)
}
