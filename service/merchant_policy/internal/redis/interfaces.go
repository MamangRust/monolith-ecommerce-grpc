package mencache

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantPolicyQueryCache interface {
	GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, bool)
	SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponse, totalRecords *int)

	GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, bool)
	SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponseDeleteAt, totalRecords *int)

	GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, bool)
	SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponseDeleteAt, totalRecords *int)

	GetCachedMerchantPolicy(ctx context.Context, id int) (*response.MerchantPoliciesResponse, bool)
	SetCachedMerchantPolicy(ctx context.Context, data *response.MerchantPoliciesResponse)
}

type MerchantPolicyCommandCache interface {
	DeleteMerchantPolicyCache(ctx context.Context, merchantID int)
}
