package cache

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryCache interface {
	GetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, *int, bool)
	SetCachedMerchantPolicyAll(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantPoliciesRow, total *int)

	GetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, *int, bool)
	SetCachedMerchantPolicyActive(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantPoliciesActiveRow, total *int)

	GetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, *int, bool)
	SetCachedMerchantPolicyTrashed(ctx context.Context, req *requests.FindAllMerchant, data []*db.GetMerchantPoliciesTrashedRow, total *int)

	GetCachedMerchantPolicy(ctx context.Context, id int) (*db.GetMerchantPolicyRow, bool)
	SetCachedMerchantPolicy(ctx context.Context, data *db.GetMerchantPolicyRow)
}

type MerchantPoliciesCommandCache interface {
	DeleteMerchantPolicyCache(ctx context.Context, merchantID int)
}
