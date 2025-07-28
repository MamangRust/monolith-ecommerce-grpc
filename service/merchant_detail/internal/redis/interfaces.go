package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantDetailQueryCache interface {
	GetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantDetailResponse, *int, bool)
	SetCachedMerchantDetailAll(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantDetailResponse, totalRecords *int)

	GetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, bool)
	SetCachedMerchantDetailActive(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantDetailResponseDeleteAt, totalRecords *int)

	GetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantDetailResponseDeleteAt, *int, bool)
	SetCachedMerchantDetailTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantDetailResponseDeleteAt, totalRecords *int)

	GetCachedMerchantDetail(ctx context.Context, id int) (*response.MerchantDetailResponse, bool)
	SetCachedMerchantDetail(ctx context.Context, data *response.MerchantDetailResponse)
}

type MerchantDetailCommandCache interface {
	DeleteMerchantDetailCache(ctx context.Context, id int)
}
