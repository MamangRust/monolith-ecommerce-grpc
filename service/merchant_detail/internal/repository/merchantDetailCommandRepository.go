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
	mapping recordmapper.MerchantDetailMapping
}

func NewMerchantDetailCommandRepository(db *db.Queries, mapping recordmapper.MerchantDetailMapping) *merchantDetailCommandRepository {
	return &merchantDetailCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantDetailCommandRepository) CreateMerchantDetail(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*record.MerchantDetailRecord, error) {
	req := db.CreateMerchantDetailParams{
		MerchantID:       int32(request.MerchantID),
		DisplayName:      sql.NullString{String: request.DisplayName, Valid: true},
		CoverImageUrl:    sql.NullString{String: request.CoverImageUrl, Valid: true},
		LogoUrl:          sql.NullString{String: request.LogoUrl, Valid: true},
		ShortDescription: sql.NullString{String: request.ShortDescription, Valid: true},
		WebsiteUrl:       sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	merchant, err := r.db.CreateMerchantDetail(ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(merchant), nil
}

func (r *merchantDetailCommandRepository) UpdateMerchantDetail(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*record.MerchantDetailRecord, error) {
	req := db.UpdateMerchantDetailParams{
		MerchantDetailID: int32(*request.MerchantDetailID),
		DisplayName:      sql.NullString{String: request.DisplayName, Valid: true},
		CoverImageUrl:    sql.NullString{String: request.CoverImageUrl, Valid: true},
		LogoUrl:          sql.NullString{String: request.LogoUrl, Valid: true},
		ShortDescription: sql.NullString{String: request.ShortDescription, Valid: true},
		WebsiteUrl:       sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	res, err := r.db.UpdateMerchantDetail(ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailCommandRepository) TrashedMerchantDetail(ctx context.Context, merchant_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.TrashMerchantDetail(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrTrashedMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailCommandRepository) RestoreMerchantDetail(ctx context.Context, merchant_id int) (*record.MerchantDetailRecord, error) {
	res, err := r.db.RestoreMerchantDetail(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail
	}

	return r.mapping.ToMerchantDetailRecord(res), nil
}

func (r *merchantDetailCommandRepository) DeleteMerchantDetailPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantDetailPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantdetail_errors.ErrDeleteMerchantDetailPermanent
	}

	return true, nil
}

func (r *merchantDetailCommandRepository) RestoreAllMerchantDetail(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDetails(ctx)

	if err != nil {
		return false, merchantdetail_errors.ErrRestoreAllMerchantDetails
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) DeleteAllMerchantDetailPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDetails(ctx)

	if err != nil {
		return false, merchantdetail_errors.ErrDeleteAllMerchantDetailsPermanent
	}
	return true, nil
}
