package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, error)
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantPolicyRow, error)
}

type MerchantPoliciesCommandRepository interface {
	Create(
		ctx context.Context,
		request *requests.CreateMerchantPolicyRequest,
	) (*db.CreateMerchantPolicyRow, error)

	Update(
		ctx context.Context,
		request *requests.UpdateMerchantPolicyRequest,
	) (*db.UpdateMerchantPolicyRow, error)

	Trash(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantPolicy, error)

	Restore(
		ctx context.Context,
		merchant_id int,
	) (*db.MerchantPolicy, error)

	DeletePermanent(
		ctx context.Context,
		merchant_id int,
	) (bool, error)

	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error)
}
