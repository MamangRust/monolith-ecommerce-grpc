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
	ctx     context.Context
	mapping recordmapper.MerchantBusinessMapping
}

func NewMerchantBusinessQueryRepository(
	db *db.Queries,
	ctx context.Context,
	mapping recordmapper.MerchantBusinessMapping,
) *merchantBusinessQueryRepository {
	return &merchantBusinessQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantBusinessQueryRepository) FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformation(r.ctx, reqDb)

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

func (r *merchantBusinessQueryRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationActive(r.ctx, reqDb)

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

func (r *merchantBusinessQueryRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantBusinessRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationTrashed(r.ctx, reqDb)

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

func (r *merchantBusinessQueryRepository) FindById(user_id int) (*record.MerchantBusinessRecord, error) {
	res, err := r.db.GetMerchantBusinessInformation(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound
	}

	return r.mapping.ToMerchantBusinessRecord(res), nil
}
