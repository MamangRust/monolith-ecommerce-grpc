package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantQueryRepository(db *db.Queries, mapping recordmapper.MerchantRecordMapping) *merchantQueryRepository {
	return &merchantQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchants(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindAllMerchants
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordPagination(res), &totalCount, nil
}

func (r *merchantQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsActive(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindByActiveMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordActivePagination(res), &totalCount, nil
}

func (r *merchantQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, merchant_errors.ErrFindByTrashedMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantsRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantQueryRepository) FindById(ctx context.Context, user_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(ctx, int32(user_id))

	if err != nil {
		return nil, merchant_errors.ErrFindByIdMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}
