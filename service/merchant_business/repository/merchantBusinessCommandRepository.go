package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
)

type merchantBusinessCommandRepository struct {
	db *db.Queries
}

func NewMerchantBusinessCommandRepository(
	db *db.Queries,
) *merchantBusinessCommandRepository {
	return &merchantBusinessCommandRepository{
		db: db,
	}
}

func (r *merchantBusinessCommandRepository) Create(
	ctx context.Context,
	request *requests.CreateMerchantBusinessInformationRequest,
) (*db.CreateMerchantBusinessInformationRow, error) {

	req := db.CreateMerchantBusinessInformationParams{
		MerchantID: int32(request.MerchantID),

		BusinessType: stringPtr(request.BusinessType),
		TaxID:        stringPtr(request.TaxID),
		WebsiteUrl:   stringPtr(request.WebsiteUrl),

		EstablishedYear:   int32Ptr(request.EstablishedYear),
		NumberOfEmployees: int32Ptr(request.NumberOfEmployees),
	}

	merchant, err := r.db.CreateMerchantBusinessInformation(ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantBusinessCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*db.UpdateMerchantBusinessInformationRow, error) {
	req := db.UpdateMerchantBusinessInformationParams{
		MerchantBusinessInfoID: int32(*request.MerchantBusinessInfoID),
		BusinessType:           stringPtr(request.BusinessType),
		TaxID:                  stringPtr(request.TaxID),
		WebsiteUrl:             stringPtr(request.WebsiteUrl),
		EstablishedYear:        int32Ptr(request.EstablishedYear),
		NumberOfEmployees:      int32Ptr(request.NumberOfEmployees),
	}

	merchant, err := r.db.UpdateMerchantBusinessInformation(ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrUpdateMerchantBusiness.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantBusinessCommandRepository) Trash(ctx context.Context, merchant_business_info_id int) (*db.MerchantBusinessInformation, error) {
	res, err := r.db.TrashMerchantBusinessInformation(ctx, int32(merchant_business_info_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrTrashMerchantBusiness.WithInternal(err)
	}

	return res, nil
}

func (r *merchantBusinessCommandRepository) Restore(ctx context.Context, merchant_business_info_id int) (*db.MerchantBusinessInformation, error) {
	res, err := r.db.RestoreMerchantBusinessInformation(ctx, int32(merchant_business_info_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrRestoreMerchantBusiness.WithInternal(err)
	}

	return res, nil
}

func (r *merchantBusinessCommandRepository) DeletePermanent(ctx context.Context, merchant_business_info_id int) (bool, error) {
	err := r.db.DeleteMerchantBusinessInformationPermanently(ctx, int32(merchant_business_info_id))

	if err != nil {
		return false, merchantbusiness_errors.ErrDeletePermanentMerchantBusiness.WithInternal(err)
	}

	return true, nil
}

func (r *merchantBusinessCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantBusinessInformation(ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrRestoreAllMerchantBusinesses.WithInternal(err)
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantBusinessInformation(ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrDeleteAllPermanentMerchantBusinesses.WithInternal(err)
	}
	return true, nil
}

func int32Ptr(v int) *int32 {
	if v == 0 {
		return nil
	}
	val := int32(v)
	return &val
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
