package merchantpolicies_cache

import (
	"context"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantPolicyQueryCache interface {
	GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantPolicies, bool)
	SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantPolicies)

	GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantPoliciesDeleteAt, bool)
	SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantPoliciesDeleteAt)

	GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) (*response.ApiResponsePaginationMerchantPoliciesDeleteAt, bool)
	SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data *response.ApiResponsePaginationMerchantPoliciesDeleteAt)

	GetCachedMerchantPolicy(ctx context.Context, id int) (*response.ApiResponseMerchantPolicies, bool)
	SetCachedMerchantPolicy(ctx context.Context, data *response.ApiResponseMerchantPolicies)
}

type MerchantPolicyCommandCache interface {
	DeleteMerchantPolicyCache(ctx context.Context, merchantID int)
}
