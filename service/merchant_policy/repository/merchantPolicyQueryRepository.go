package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_policy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
)

type merchantPolicyQueryRepository struct {
	db *db.Queries
}

func NewMerchantPolicyQueryRepository(db *db.Queries) *merchantPolicyQueryRepository {
	return &merchantPolicyQueryRepository{
		db: db,
	}
}

func (r *merchantPolicyQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPolicies(ctx, reqDb)

	if err != nil {
		return nil, merchant_policy_errors.ErrFindAllMerchantPolicies.WithInternal(err)
	}

	return res, nil
}

func (r *merchantPolicyQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesActive(ctx, reqDb)

	if err != nil {
		return nil, merchant_policy_errors.ErrFindActiveMerchantPolicies.WithInternal(err)
	}

	return res, nil
}

func (r *merchantPolicyQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantPoliciesTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchant_policy_errors.ErrFindTrashedMerchantPolicies.WithInternal(err)
	}

	return res, nil
}

func (r *merchantPolicyQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetMerchantPolicyRow, error) {
	res, err := r.db.GetMerchantPolicy(ctx, int32(user_id))

	if err != nil {
		return nil, merchant_policy_errors.ErrFindMerchantPolicyByID.WithInternal(err)
	}

	return res, nil
}
