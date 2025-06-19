package mencache

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantPolicyQueryCache interface {
	GetCachedMerchantPolicyAll(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponse, *int, bool)
	SetCachedMerchantPolicyAll(req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponse, totalRecords *int)
	GetCachedMerchantPolicyActive(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, bool)
	SetCachedMerchantPolicyActive(req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponseDeleteAt, totalRecords *int)
	GetCachedMerchantPolicyTrashed(req *requests.FindAllMerchant) ([]*response.MerchantPoliciesResponseDeleteAt, *int, bool)
	SetCachedMerchantPolicyTrashed(req *requests.FindAllMerchant, data []*response.MerchantPoliciesResponseDeleteAt, totalRecords *int)

	GetCachedMerchantPolicy(id int) (*response.MerchantPoliciesResponse, bool)
	SetCachedMerchantPolicy(data *response.MerchantPoliciesResponse)
}

type MerchantPolicyCommandCache interface {
	DeleteMerchantPolicyCache(merchantID int)
}
