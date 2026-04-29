package repository

import (
	"context"
	"database/sql"
	errorsstd "errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
)

type bannerQueryRepository struct {
	db *db.Queries
}

func NewBannerQueryRepository(db *db.Queries) *bannerQueryRepository {
	return &bannerQueryRepository{
		db: db,
	}
}

func (r *bannerQueryRepository) FindAll(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBanners(ctx, reqDb)

	if err != nil {
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err)
	}

	return res, nil
}

func (r *bannerQueryRepository) FindActive(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersActive(ctx, reqDb)

	if err != nil {
		return nil, banner_errors.ErrFindActiveBanners.WithInternal(err)
	}

	return res, nil
}

func (r *bannerQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*db.GetBannersTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersTrashed(ctx, reqDb)

	if err != nil {
		return nil, banner_errors.ErrFindTrashedBanners.WithInternal(err)
	}

	return res, nil
}

func (r *bannerQueryRepository) FindByID(ctx context.Context, bannerID int) (*db.GetBannerRow, error) {
	res, err := r.db.GetBanner(ctx, int32(bannerID))

	if err != nil {
		if errorsstd.Is(err, sql.ErrNoRows) {
			return nil, banner_errors.ErrBannerNotFound.WithInternal(err)
		}
		return nil, banner_errors.ErrFindAllBanners.WithInternal(err) // Generic find error
	}

	return res, nil
}

