package repository

import (
	"context"

	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
)

type merchantBusinessQueryRepository struct {
	db *db.Queries
}

func NewMerchantBusinessQueryRepository(
	db *db.Queries,
) *merchantBusinessQueryRepository {
	return &merchantBusinessQueryRepository{
		db: db,
	}
}

func (r *merchantBusinessQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformation(ctx, reqDb)

	if err != nil {
		return nil, merchantbusiness_errors.ErrFindAllMerchantBusinesses.WithInternal(err)
	}

	return res, nil
}

func (r *merchantBusinessQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationActive(ctx, reqDb)

	if err != nil {
		return nil, merchantbusiness_errors.ErrFindActiveMerchantBusinesses.WithInternal(err)
	}

	return res, nil
}

func (r *merchantBusinessQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsBusinessInformationTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsBusinessInformationTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchantbusiness_errors.ErrFindTrashedMerchantBusinesses.WithInternal(err)
	}

	return res, nil
}

func (r *merchantBusinessQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetMerchantBusinessInformationRow, error) {
	res, err := r.db.GetMerchantBusinessInformation(ctx, int32(user_id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, merchantbusiness_errors.ErrMerchantBusinessNotFound.WithInternal(err)
		}
		return nil, merchantbusiness_errors.ErrMerchantBusinessInternal.WithInternal(err)
	}

	return res, nil
}
