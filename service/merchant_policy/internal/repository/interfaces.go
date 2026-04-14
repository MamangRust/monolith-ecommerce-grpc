package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryRepository interface {
	FindAllMerchantPolicy(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, error)
	FindById(ctx context.Context, user_id int) (*db.GetMerchantPolicyRow, error)
}

type MerchantPoliciesCommandRepository interface {
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

type MerchantQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}
