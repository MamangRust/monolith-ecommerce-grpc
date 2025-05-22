package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantPolicyCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantPolicyMapping
}

func NewMerchantPolicyCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantPolicyMapping) *merchantPolicyCommandRepository {
	return &merchantPolicyCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantPolicyCommandRepository) CreateMerchantPolicy(request *requests.CreateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error) {
	req := db.CreateMerchantPolicyParams{
		MerchantID:  int32(request.MerchantID),
		PolicyType:  request.PolicyType,
		Title:       request.Title,
		Description: request.Description,
	}

	policy, err := r.db.CreateMerchantPolicy(r.ctx, req)
	if err != nil {
		return nil, merchantpolicy_errors.ErrCreateMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(policy), nil
}

func (r *merchantPolicyCommandRepository) UpdateMerchantPolicy(request *requests.UpdateMerchantPolicyRequest) (*record.MerchantPoliciesRecord, error) {
	req := db.UpdateMerchantPolicyParams{
		MerchantPolicyID: int32(*request.MerchantPolicyID),
		PolicyType:       request.PolicyType,
		Title:            request.Title,
		Description:      request.Description,
	}

	res, err := r.db.UpdateMerchantPolicy(r.ctx, req)
	if err != nil {
		return nil, merchantpolicy_errors.ErrUpdateMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}

func (r *merchantPolicyCommandRepository) TrashedMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error) {
	res, err := r.db.TrashMerchantPolicy(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrTrashedMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}

func (r *merchantPolicyCommandRepository) RestoreMerchantPolicy(merchant_id int) (*record.MerchantPoliciesRecord, error) {
	res, err := r.db.RestoreMerchantPolicy(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrRestoreMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}

func (r *merchantPolicyCommandRepository) DeleteMerchantPolicyPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteMerchantPolicyPermanent
	}

	return true, nil
}

func (r *merchantPolicyCommandRepository) RestoreAllMerchantPolicy() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, merchantpolicy_errors.ErrRestoreAllMerchantPolicy
	}
	return true, nil
}

func (r *merchantPolicyCommandRepository) DeleteAllMerchantPolicyPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, merchantpolicy_errors.ErrDeleteAllMerchantPolicyPermanent
	}
	return true, nil
}
