package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantBusinessQueryCache interface {
	GetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponse, *int, bool)
	SetCachedMerchantBusinessAll(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantBusinessResponse, totalRecords *int)

	GetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, bool)
	SetCachedMerchantBusinessActive(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantBusinessResponseDeleteAt, totalRecords *int)

	GetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantBusinessResponseDeleteAt, *int, bool)
	SetCachedMerchantBusinessTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantBusinessResponseDeleteAt, totalRecords *int)

	GetCachedMerchantBusiness(ctx context.Context, id int) (*response.MerchantBusinessResponse, bool)
	SetCachedMerchantBusiness(ctx context.Context, data *response.MerchantBusinessResponse)
}

type MerchantBusinessCommandCache interface {
	DeleteMerchantBusinessCache(ctx context.Context, id int)
}
