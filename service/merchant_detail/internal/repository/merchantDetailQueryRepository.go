package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantDetailQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantDetailMapping
}

func NewMerchantDetailQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantDetailMapping) *merchantDetailQueryRepository {
	return &merchantDetailQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantDetailQueryRepository) FindAllMerchants(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetails(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantdetail_errors.ErrFindAllMerchantDetails
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantDetailsRecordPagination(res), &totalCount, nil
}

func (r *merchantDetailQueryRepository) FindByActive(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantdetail_errors.ErrFindByActiveMerchantDetails
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantDetailsRecordActivePagination(res), &totalCount, nil
}

func (r *merchantDetailQueryRepository) FindByTrashed(req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, merchantdetail_errors.ErrFindByTrashedMerchantDetails
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToMerchantDetailsRecordTrashedPagination(res), &totalCount, nil
}

func (r *merchantDetailQueryRepository) FindById(user_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.GetMerchantDetail(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByIdMerchantDetail
	}

	return r.mapping.ToMerchantDetailRelationRecord(res), nil
}

func (r *merchantDetailQueryRepository) FindByIdTrashed(user_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.GetMerchantDetailTrashed(r.ctx, int32(user_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByIdTrashedMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}
