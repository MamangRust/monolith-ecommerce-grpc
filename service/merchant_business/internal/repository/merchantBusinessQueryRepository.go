package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantBusinessQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.MerchantBusinessMapping
}

func NewMerchantBusinessQueryRepository(
	db *db.Queries,
	mapping recordmapper.MerchantBusinessMapping,
) *merchantBusinessQueryRepository {
	return &merchantBusinessQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantBusinessQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformation(ctx, reqDb)

	if err != nil {
		return nil, nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantBusinesssRecordPagination(res), &totalCount, nil
}

func (r *merchantBusinessQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationActive(ctx, reqDb)

	if err != nil {
		return nil, nil, merchantbusiness_errors.ErrFindActiveMerchantBusinesses
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantBusinesssRecordActivePagination(res), &totalCount, nil
}

func (r *merchantBusinessQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, merchantbusiness_errors.ErrFindTrashedMerchantBusinesses
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantBusinesssRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantBusinessQueryRepository) FindById(ctx context.Context, user_id int) (*record.MerchantBusinessRecord, error) {
	res, err := r.db.GetMerchantBusinessInformation(ctx, int32(user_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}

	return r.mapping.ToMerchantBusinessRecord(res), nil
}
