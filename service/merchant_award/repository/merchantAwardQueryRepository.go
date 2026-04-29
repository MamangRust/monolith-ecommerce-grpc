package repository

import (
	"context"

	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
)

type merchantAwardQueryRepository struct {
	db *db.Queries
}

func NewMerchantAwardQueryRepository(db *db.Queries) *merchantAwardQueryRepository {
	return &merchantAwardQueryRepository{
		db: db,
	}
}

func (r *merchantAwardQueryRepository) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwards(ctx, reqDb)

	if err != nil {
		return nil, merchantaward_errors.ErrFindAllMerchantAwards.WithInternal(err)
	}

	return res, nil
}

func (r *merchantAwardQueryRepository) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsActive(ctx, reqDb)

	if err != nil {
		return nil, merchantaward_errors.ErrFindByActiveMerchantAwards.WithInternal(err)
	}

	return res, nil
}

func (r *merchantAwardQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantCertificationsAndAwardsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetMerchantCertificationsAndAwardsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetMerchantCertificationsAndAwardsTrashed(ctx, reqDb)

	if err != nil {
		return nil, merchantaward_errors.ErrFindByTrashedMerchantAwards.WithInternal(err)
	}

	return res, nil
}

func (r *merchantAwardQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetMerchantCertificationOrAwardRow, error) {
	res, err := r.db.GetMerchantCertificationOrAward(ctx, int32(user_id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, merchantaward_errors.ErrMerchantAwardNotFound.WithInternal(err)
		}
		return nil, merchantaward_errors.ErrFindByIdMerchantAward.WithInternal(err)
	}

	return res, nil
}
