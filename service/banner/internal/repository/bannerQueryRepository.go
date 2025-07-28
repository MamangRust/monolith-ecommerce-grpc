package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type bannerQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.BannerRecordMapping
}

func NewBannerQueryRepository(db *db.Queries, mapping recordmapper.BannerRecordMapping) *bannerQueryRepository {
	return &bannerQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *bannerQueryRepository) FindAllBanners(ctx context.Context, req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBanners(ctx, reqDb)

	if err != nil {
		return nil, nil, banner_errors.ErrFindAllBanners
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToBannersRecordPagination(res), &totalCount, nil
}

func (r *bannerQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersActive(ctx, reqDb)

	if err != nil {
		return nil, nil, banner_errors.ErrFindActiveBanners
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToBannersRecordActivePagination(res), &totalCount, nil
}

func (r *bannerQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllBanner) ([]*record.BannerRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetBannersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetBannersTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, banner_errors.ErrFindTrashedBanners
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToBannersRecordTrashedPagination(res), &totalCount, nil
}

func (r *bannerQueryRepository) FindById(ctx context.Context, user_id int) (*record.BannerRecord, error) {
	res, err := r.db.GetBanner(ctx, int32(user_id))

	if err != nil {
		return nil, banner_errors.ErrBannerNotFound
	}

	return r.mapping.ToBannerRecord(res), nil
}
