package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, *int, error)
	FindByID(ctx context.Context, user_id int) (*db.GetMerchantPolicyRow, error)
}

type MerchantPoliciesCommandService interface {
	Create(ctx context.Context, request *requests.CreateMerchantPolicyRequest) (*db.CreateMerchantPolicyRow, error)
	Update(ctx context.Context, request *requests.UpdateMerchantPolicyRequest) (*db.UpdateMerchantPolicyRow, error)
	Trash(ctx context.Context, merchant_id int) (*db.MerchantPolicy, error)
	Restore(ctx context.Context, merchant_id int) (*db.MerchantPolicy, error)
	DeletePermanent(ctx context.Context, merchant_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
