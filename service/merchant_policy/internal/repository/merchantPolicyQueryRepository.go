package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantPolicyQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantPolicyMapping
}

func NewMerchantPolicyQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantPolicyMapping) *merchantPolicyQueryRepository {
	return &merchantPolicyQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantPolicyQueryRepository) FindAllMerchantPolicy(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPolicies(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantpolicy_errors.ErrFindAllMerchantPolicy
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantPolicysRecordPagination(res), &totalCount, nil
}

func (r *merchantPolicyQueryRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantpolicy_errors.ErrFindByActiveMerchantPolicy
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantPolicysRecordActivePagination(res), &totalCount, nil
}

func (r *merchantPolicyQueryRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantPoliciesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantPoliciesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantPoliciesTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantpolicy_errors.ErrFindByTrashedMerchantPolicy
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantPolicysRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantPolicyQueryRepository) FindById(user_id int) (*record.MerchantPoliciesRecord, error) {
	res, err := r.db.GetMerchantPolicy(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchantpolicy_errors.ErrFindByIdMerchantPolicy
	}

	return r.mapping.ToMerchantPolicyRecord(res), nil
}
