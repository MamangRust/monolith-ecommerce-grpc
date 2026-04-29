package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
)

type merchantDetailCommandRepository struct {
	db *db.Queries
}

func NewMerchantDetailCommandRepository(db *db.Queries) *merchantDetailCommandRepository {
	return &merchantDetailCommandRepository{
		db: db,
	}
}

func (r *merchantDetailCommandRepository) Create(ctx context.Context, request *requests.CreateMerchantDetailRequest) (*db.CreateMerchantDetailRow, error) {
	req := db.CreateMerchantDetailParams{
		MerchantID:       int32(request.MerchantID),
		DisplayName:      &request.DisplayName,
		CoverImageUrl:    &request.CoverImageUrl,
		LogoUrl:          &request.LogoUrl,
		ShortDescription: &request.ShortDescription,
		WebsiteUrl:       &request.WebsiteUrl,
	}

	merchant, err := r.db.CreateMerchantDetail(ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrCreateMerchantDetail.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantDetailCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantDetailRequest) (*db.UpdateMerchantDetailRow, error) {
	req := db.UpdateMerchantDetailParams{
		MerchantDetailID: int32(*request.MerchantDetailID),
		DisplayName:      &request.DisplayName,
		CoverImageUrl:    &request.CoverImageUrl,
		LogoUrl:          &request.LogoUrl,
		ShortDescription: &request.ShortDescription,
		WebsiteUrl:       &request.WebsiteUrl,
	}

	res, err := r.db.UpdateMerchantDetail(ctx, req)
	if err != nil {
		return nil, merchantdetail_errors.ErrUpdateMerchantDetail.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDetailCommandRepository) Trash(ctx context.Context, merchant_detail_id int) (*db.MerchantDetail, error) {
	res, err := r.db.TrashMerchantDetail(ctx, int32(merchant_detail_id))
	if err != nil {
		return nil, merchantdetail_errors.ErrTrashMerchantDetail.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDetailCommandRepository) Restore(ctx context.Context, merchant_detail_id int) (*db.MerchantDetail, error) {
	res, err := r.db.RestoreMerchantDetail(ctx, int32(merchant_detail_id))
	if err != nil {
		return nil, merchantdetail_errors.ErrRestoreMerchantDetail.WithInternal(err)
	}

	return res, nil
}

func (r *merchantDetailCommandRepository) DeletePermanent(ctx context.Context, merchant_detail_id int) (bool, error) {
	err := r.db.DeleteMerchantDetailPermanently(ctx, int32(merchant_detail_id))
	if err != nil {
		return false, merchantdetail_errors.ErrDeletePermanentMerchantDetail.WithInternal(err)
	}

	return true, nil
}

func (r *merchantDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantDetails(ctx)
	if err != nil {
		return false, merchantdetail_errors.ErrRestoreAllMerchantDetails.WithInternal(err)
	}
	return true, nil
}

func (r *merchantDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantDetails(ctx)
	if err != nil {
		return false, merchantdetail_errors.ErrDeleteAllPermanentMerchantDetails.WithInternal(err)
	}
	return true, nil
}
