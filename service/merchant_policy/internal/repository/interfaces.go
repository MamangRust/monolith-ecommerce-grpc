package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type MerchantPoliciesQueryRepository interface {
	FindAllMerchantPolicy(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error)
	FindById(ctx context.Context, userID int) (*record.MerchantPoliciesRecord, error)
}

type MerchantPoliciesCommandRepository interface {
	CreateMerchantPolicy(ctx context.Context, request *requests.CreateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error)
	UpdateMerchantPolicy(ctx context.Context, request *requests.UpdateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error)
	TrashedMerchantPolicy(ctx context.Context, merchantID int) (*record.MerchantPoliciesRecord, error)
	RestoreMerchantPolicy(ctx context.Context, merchantID int) (*record.MerchantPoliciesRecord, error)
	DeleteMerchantPolicyPermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAllMerchantPolicy(ctx context.Context) (bool, error)
	DeleteAllMerchantPolicyPermanent(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, userID int) (*record.MerchantRecord, error)
}
