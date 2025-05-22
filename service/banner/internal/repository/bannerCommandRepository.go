package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type bannerCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.BannerRecordMapping
}

func NewBannerCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.BannerRecordMapping) *bannerCommandRepository {
	return &bannerCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *bannerCommandRepository) CreateBanner(request *requests.CreateBannerRequest) (*record.BannerRecord, error) {
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	req := db.CreateBannerParams{
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  sql.NullBool{Bool: request.IsActive, Valid: true},
	}

	result, err := r.db.CreateBanner(r.ctx, req)
	if err != nil {
		return nil, banner_errors.ErrCreateBanner
	}

	return r.mapping.ToBannerRecord(result), nil
}

func (r *bannerCommandRepository) UpdateBanner(request *requests.UpdateBannerRequest) (*record.BannerRecord, error) {
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, banner_errors.ErrBannerStartDate
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		return nil, banner_errors.ErrBannerEndDate
	}

	startTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return nil, banner_errors.ErrBannerStartTime
	}

	endTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return nil, banner_errors.ErrBannerEndTime
	}

	req := db.UpdateBannerParams{
		BannerID:  int32(*request.BannerID),
		Name:      request.Name,
		StartDate: startDate,
		EndDate:   endDate,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  sql.NullBool{Bool: request.IsActive, Valid: true},
	}

	result, err := r.db.UpdateBanner(r.ctx, req)
	if err != nil {
		return nil, banner_errors.ErrUpdateBanner
	}

	return r.mapping.ToBannerRecord(result), nil
}

func (r *bannerCommandRepository) TrashedBanner(Banner_id int) (*record.BannerRecord, error) {
	res, err := r.db.TrashBanner(r.ctx, int32(Banner_id))

	if err != nil {
		return nil, banner_errors.ErrTrashedBanner
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerCommandRepository) RestoreBanner(Banner_id int) (*record.BannerRecord, error) {
	res, err := r.db.RestoreBanner(r.ctx, int32(Banner_id))

	if err != nil {
		return nil, banner_errors.ErrRestoreBanner
	}

	return r.mapping.ToBannerRecord(res), nil
}

func (r *bannerCommandRepository) DeleteBannerPermanent(Banner_id int) (bool, error) {
	err := r.db.DeleteBannerPermanently(r.ctx, int32(Banner_id))

	if err != nil {
		return false, banner_errors.ErrDeleteBannerPermanent
	}

	return true, nil
}

func (r *bannerCommandRepository) RestoreAllBanner() (bool, error) {
	err := r.db.RestoreAllBanners(r.ctx)

	if err != nil {
		return false, banner_errors.ErrRestoreAllBanners
	}
	return true, nil
}

func (r *bannerCommandRepository) DeleteAllBannerPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentBanners(r.ctx)

	if err != nil {
		return false, banner_errors.ErrDeleteAllBanners
	}
	return true, nil
}
