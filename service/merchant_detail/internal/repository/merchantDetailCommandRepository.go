package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantDetailCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantDetailMapping
}

func NewMerchantDetailCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantDetailMapping) *merchantDetailCommandRepository {
	return &merchantDetailCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantDetailCommandRepository) CreateMerchantDetail(request *requests.CreateMerchantDetailRequest) (*record.MerchantDetailRecord, error) {
	req := db.CreateMerchantDetailParams{
		MerchantID:       int32(request.MerchantID),
		DisplayName:      sql.NullString{String: request.DisplayName, Valid: true},
		CoverImageUrl:    sql.NullString{String: request.CoverImageUrl, Valid: true},
		LogoUrl:          sql.NullString{String: request.LogoUrl, Valid: true},
		ShortDescription: sql.NullString{String: request.ShortDescription, Valid: true},
		WebsiteUrl:       sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	merchant, err := r.db.CreateMerchantDetail(r.ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(merchant), nil
}

func (r *merchantDetailCommandRepository) UpdateMerchantDetail(request *requests.UpdateMerchantDetailRequest) (*record.MerchantDetailRecord, error) {
	req := db.UpdateMerchantDetailParams{
		MerchantDetailID: int32(*request.MerchantDetailID),
		DisplayName:      sql.NullString{String: request.DisplayName, Valid: true},
		CoverImageUrl:    sql.NullString{String: request.CoverImageUrl, Valid: true},
		LogoUrl:          sql.NullString{String: request.LogoUrl, Valid: true},
		ShortDescription: sql.NullString{String: request.ShortDescription, Valid: true},
		WebsiteUrl:       sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	res, err := r.db.UpdateMerchantDetail(r.ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailCommandRepository) TrashedMerchantDetail(merchant_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.TrashMerchantDetail(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrTrashedMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailCommandRepository) RestoreMerchantDetail(merchant_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.RestoreMerchantDetail(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailCommandRepository) DeleteMerchantDetailPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantDetailPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantdetail_errors.ErrDeleteMerchantDetailPermanent
	}

	return true, nil
}

func (r *merchantDetailCommandRepository) RestoreAllMerchantDetail() (bool, error) {
	err := r.db.RestoreAllMerchantDetails(r.ctx)

	if err != nil {
		return false, merchantdetail_errors.ErrRestoreAllMerchantDetails
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) DeleteAllMerchantDetailPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDetails(r.ctx)

	if err != nil {
		return false, merchantdetail_errors.ErrDeleteAllMerchantDetailsPermanent
	}
	return true, nil
}
