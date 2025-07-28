package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantAwardQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.MerchantAwardMapping
}

func NewMerchantAwardQueryRepository(db *db.Queries, mapping recordmapper.MerchantAwardMapping) *merchantAwardQueryRepository {
	return &merchantAwardQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantAwardQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwards(ctx, reqDb)

	if err != nil {
		return nil, nil, merchantaward_errors.ErrFindAllMerchantAwards
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantAwardsRecordPagination(res), &totalCount, nil
}

func (r *merchantAwardQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsActive(ctx, reqDb)

	if err != nil {
		return nil, nil, merchantaward_errors.ErrFindByActiveMerchantAwards
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantAwardsRecordActivePagination(res), &totalCount, nil
}

func (r *merchantAwardQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantAwardRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, merchantaward_errors.ErrFindByTrashedMerchantAwards
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantAwardsRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantAwardQueryRepository) FindById(ctx context.Context, user_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.GetMerchantCertificationOrAward(ctx, int32(user_id))

	if err != nil {
		return nil, merchantaward_errors.ErrFindByIdMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}
