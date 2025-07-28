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
	db *db.Queries

	mapping recordmapper.MerchantDetailMapping
}

func NewMerchantDetailQueryRepository(db *db.Queries, mapping recordmapper.MerchantDetailMapping) *merchantDetailQueryRepository {
	return &merchantDetailQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantDetailQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetails(ctx, reqDb)

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

func (r *merchantDetailQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsActive(ctx, reqDb)

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

func (r *merchantDetailQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*record.MerchantDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsTrashed(ctx, reqDb)

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

func (r *merchantDetailQueryRepository) FindById(ctx context.Context, user_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.GetMerchantDetail(ctx, int32(user_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByIdMerchantDetail
	}

	return r.mapping.ToMerchantDetailRelationRecord(res), nil
}

func (r *merchantDetailQueryRepository) FindByIdTrashed(ctx context.Context, user_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.GetMerchantDetailTrashed(ctx, int32(user_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrFindByIdTrashedMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}
