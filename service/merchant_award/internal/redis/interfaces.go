package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantAwardQueryCache interface {
	GetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponse, *int, bool)
	SetCachedMerchantAwardAll(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantAwardResponse, totalRecords *int)

	GetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, bool)
	SetCachedMerchantAwardActive(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantAwardResponseDeleteAt, totalRecords *int)

	GetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantAwardResponseDeleteAt, *int, bool)
	SetCachedMerchantAwardTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantAwardResponseDeleteAt, totalRecords *int)

	GetCachedMerchantAward(ctx context.Context, id int) (*response.MerchantAwardResponse, bool)
	SetCachedMerchantAward(ctx context.Context, data *response.MerchantAwardResponse)
}

type MerchantAwardCommandCache interface {
	DeleteMerchantAwardCache(ctx context.Context, id int)
}
