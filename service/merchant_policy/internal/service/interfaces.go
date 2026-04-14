package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryService interface {
	FindAllMerchantPolicy(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, *int, error)
	FindById(ctx context.Context, user_id int) (*db.GetMerchantPolicyRow, error)
}

type MerchantPoliciesCommandService interface {
	CreateMerchantPolicy(
		ctx context.Context,
		request *requests.CreateMerchantPolicyRequest,
	) (*db.CreateMerchantPolicyRow, error)

	UpdateMerchantPolicy(
		ctx context.Context,
		request *requests.UpdateMerchantPolicyRequest,
	) (*db.UpdateMerchantPolicyRow, error)

	TrashedMerchantPolicy(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantPolicy, error)

	RestoreMerchantPolicy(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantPolicy, error)

	DeleteMerchantPolicyPermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAllMerchantPolicy(ctx context.Context) (bool, error)
	DeleteAllMerchantPolicyPermanent(ctx context.Context) (bool, error)
}
