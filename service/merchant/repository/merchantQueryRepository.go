package repository

import (
	"context"

	"database/sql"
 
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
)


type merchantQueryRepository struct {
	db *db.Queries
}

func NewMerchantQueryRepository(db *db.Queries) *merchantQueryRepository {
	return &merchantQueryRepository{
		db: db,
	}
}

func (r *merchantQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchants(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindAllMerchants.WithInternal(err)
	}


	return res, nil
}

func (r *merchantQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsActive(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindActiveMerchants.WithInternal(err)
	}


	return res, nil
}

func (r *merchantQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantsTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchant_errors.ErrFindTrashedMerchants.WithInternal(err)
	}


	return res, nil
}

func (r *merchantQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetMerchantByIDRow, error) {
	res, err := r.db.GetMerchantByID(ctx, int32(user_id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, merchant_errors.ErrMerchantNotFound.WithInternal(err)
		}
		return nil, merchant_errors.ErrMerchantInternal.WithInternal(err)
	}


	return res, nil
}
