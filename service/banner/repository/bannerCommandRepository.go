package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/utils"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
)

type bannerCommandRepository struct {
	db *db.Queries
}

func NewBannerCommandRepository(db *db.Queries) *bannerCommandRepository {
	return &bannerCommandRepository{
		db: db,
	}
}

func (r *bannerCommandRepository) Create(ctx context.Context, request *requests.CreateBannerRequest) (*db.CreateBannerRow, error) {
	startDate, err := utils.ParseDate(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}

	endDate, err := utils.ParseDate(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}

	startTime, err := utils.ParseTime(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}

	endTime, err := utils.ParseTime(request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	req := db.CreateBannerParams{
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  &request.IsActive,
	}

	result, err := r.db.CreateBanner(ctx, req)
	if err != nil {
		return nil, banner_errors.ErrCreateBanner.WithInternal(err)
	}

	return result, nil
}

func (r *bannerCommandRepository) Update(ctx context.Context, request *requests.UpdateBannerRequest) (*db.UpdateBannerRow, error) {
	startDate, err := utils.ParseDate(request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate.WithInternal(err)
	}

	endDate, err := utils.ParseDate(request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate.WithInternal(err)
	}

	startTime, err := utils.ParseTime(request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime.WithInternal(err)
	}

	endTime, err := utils.ParseTime(request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime.WithInternal(err)
	}

	req := db.UpdateBannerParams{
		BannerID:  int32(*request.BannerID),
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  &request.IsActive,
	}

	result, err := r.db.UpdateBanner(ctx, req)
	if err != nil {
		return nil, banner_errors.ErrUpdateBanner.WithInternal(err)
	}

	return result, nil
}

func (r *bannerCommandRepository) Trash(ctx context.Context, bannerID int) (*db.Banner, error) {
	res, err := r.db.TrashBanner(ctx, int32(bannerID))

	if err != nil {
		return nil, banner_errors.ErrTrashedBanner.WithInternal(err)
	}

	return res, nil
}

func (r *bannerCommandRepository) Restore(ctx context.Context, bannerID int) (*db.Banner, error) {
	res, err := r.db.RestoreBanner(ctx, int32(bannerID))

	if err != nil {
		return nil, banner_errors.ErrRestoreBanner.WithInternal(err)
	}

	return res, nil
}

func (r *bannerCommandRepository) DeletePermanent(ctx context.Context, bannerID int) (bool, error) {
	err := r.db.DeleteBannerPermanently(ctx, int32(bannerID))

	if err != nil {
		return false, banner_errors.ErrDeleteBannerPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *bannerCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllBanners(ctx)

	if err != nil {
		return false, banner_errors.ErrRestoreAllBanners.WithInternal(err)
	}
	return true, nil
}

func (r *bannerCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentBanners(ctx)

	if err != nil {
		return false, banner_errors.ErrDeleteAllBanners.WithInternal(err)
	}
	return true, nil
}

