package repository

import (
	"context"

	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
)

type merchantDetailQueryRepository struct {
	db *db.Queries
}

func NewMerchantDetailQueryRepository(db *db.Queries) *merchantDetailQueryRepository {
	return &merchantDetailQueryRepository{
		db: db,
	}
}

func (r *merchantDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetails(ctx, reqDb)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindAllMerchantDetails.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsActive(ctx, reqDb)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindActiveMerchantDetails.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantDetailsTrashed(ctx, reqDb)
	if err != nil {
		return nil, merchantdetail_errors.ErrFindTrashedMerchantDetails.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDetailQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetMerchantDetailRow, error) {
	res, err := r.db.GetMerchantDetail(ctx, int32(user_id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, merchantdetail_errors.ErrMerchantDetailNotFound.WithInternal(err)
		}
		return nil, merchantdetail_errors.ErrMerchantDetailInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDetailQueryRepository) FindByIDTrashed(ctx context.Context, user_id int) (*db.GetMerchantDetailTrashedRow, error) {
	res, err := r.db.GetMerchantDetailTrashed(ctx, int32(user_id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, merchantdetail_errors.ErrMerchantDetailNotFound.WithInternal(err)
		}
		return nil, merchantdetail_errors.ErrFindByIdTrashedMerchantDetail.WithInternal(err)
	}

	return res, nil
}
