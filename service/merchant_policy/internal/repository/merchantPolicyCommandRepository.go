package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_policy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
)

type merchantPolicyCommandRepository struct {
	db *db.Queries
}

func NewMerchantPolicyCommandRepository(db *db.Queries) *merchantPolicyCommandRepository {
	return &merchantPolicyCommandRepository{
		db: db,
	}
}

func (r *merchantPolicyCommandRepository) CreateMerchantPolicy(ctx context.Context, request *requests.CreateMerchantPolicyRequest) (*db.CreateMerchantPolicyRow, error) {
	req := db.CreateMerchantPolicyParams{
		MerchantID:  int32(request.MerchantID),
		PolicyType:  request.PolicyType,
		Title:       request.Title,
		Description: request.Description,
	}

	policy, err := r.db.CreateMerchantPolicy(ctx, req)
	if err != nil {
		return nil, merchant_policy_errors.ErrCreateMerchantPolicy.WithInternal(err)
	}

	return policy, nil
}

func (r *merchantPolicyCommandRepository) UpdateMerchantPolicy(ctx context.Context, request *requests.UpdateMerchantPolicyRequest) (*db.UpdateMerchantPolicyRow, error) {
	req := db.UpdateMerchantPolicyParams{
		MerchantPolicyID: int32(*request.MerchantPolicyID),
		PolicyType:       request.PolicyType,
		Title:            request.Title,
		Description:      request.Description,
	}

	res, err := r.db.UpdateMerchantPolicy(ctx, req)
	if err != nil {
		return nil, merchant_policy_errors.ErrUpdateMerchantPolicy.WithInternal(err)
	}

	return res, nil
}

func (r *merchantPolicyCommandRepository) TrashedMerchantPolicy(ctx context.Context, merchant_id int) (*db.MerchantPolicy, error) {
	res, err := r.db.TrashMerchantPolicy(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_policy_errors.ErrTrashedMerchantPolicy.WithInternal(err)
	}

	return res, nil
}

func (r *merchantPolicyCommandRepository) RestoreMerchantPolicy(ctx context.Context, merchant_id int) (*db.MerchantPolicy, error) {
	res, err := r.db.RestoreMerchantPolicy(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_policy_errors.ErrRestoreMerchantPolicy.WithInternal(err)
	}

	return res, nil
}

func (r *merchantPolicyCommandRepository) DeleteMerchantPolicyPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchant_policy_errors.ErrDeleteMerchantPolicyPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *merchantPolicyCommandRepository) RestoreAllMerchantPolicy(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		return false, merchant_policy_errors.ErrRestoreAllMerchantPolicies.WithInternal(err)
	}
	return true, nil
}

func (r *merchantPolicyCommandRepository) DeleteAllMerchantPolicyPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		return false, merchant_policy_errors.ErrDeleteAllMerchantPoliciesPermanent.WithInternal(err)
	}
	return true, nil
}
