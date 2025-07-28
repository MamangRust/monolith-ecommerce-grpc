package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantBusinessCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.MerchantBusinessMapping
}

func NewMerchantBusinessCommandRepository(
	db *db.Queries,
	mapping recordmapper.MerchantBusinessMapping,
) *merchantBusinessCommandRepository {
	return &merchantBusinessCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantBusinessCommandRepository) CreateMerchantBusiness(ctx context.Context, request *requests.CreateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error) {
	req := db.CreateMerchantBusinessInformationParams{
		MerchantID:        int32(request.MerchantID),
		BusinessType:      sql.NullString{String: request.BusinessType, Valid: request.BusinessType != ""},
		TaxID:             sql.NullString{String: request.TaxID, Valid: request.TaxID != ""},
		EstablishedYear:   sql.NullInt32{Int32: int32(request.EstablishedYear), Valid: request.EstablishedYear != 0},
		NumberOfEmployees: sql.NullInt32{Int32: int32(request.NumberOfEmployees), Valid: request.NumberOfEmployees != 0},
		WebsiteUrl:        sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	merchant, err := r.db.CreateMerchantBusinessInformation(ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrCreateMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(merchant), nil
}

func (r *merchantBusinessCommandRepository) UpdateMerchantBusiness(ctx context.Context, request *requests.UpdateMerchantBusinessInformationRequest) (*record.MerchantBusinessRecord, error) {
	req := db.UpdateMerchantBusinessInformationParams{
		MerchantBusinessInfoID: int32(*request.MerchantBusinessInfoID),
		BusinessType:           sql.NullString{String: request.BusinessType, Valid: request.BusinessType != ""},
		TaxID:                  sql.NullString{String: request.TaxID, Valid: request.TaxID != ""},
		EstablishedYear:        sql.NullInt32{Int32: int32(request.EstablishedYear), Valid: request.EstablishedYear != 0},
		NumberOfEmployees:      sql.NullInt32{Int32: int32(request.NumberOfEmployees), Valid: request.NumberOfEmployees != 0},
		WebsiteUrl:             sql.NullString{String: request.WebsiteUrl, Valid: request.WebsiteUrl != ""},
	}

	merchant, err := r.db.UpdateMerchantBusinessInformation(ctx, req)
	if err != nil {
		return nil, merchantbusiness_errors.ErrUpdateMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(merchant), nil
}

func (r *merchantBusinessCommandRepository) TrashedMerchantBusiness(ctx context.Context, merchant_id int) (*record.MerchantBusinessRecord, error) {
	res, err := r.db.TrashMerchantBusinessInformation(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrTrashMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(res), nil
}

func (r *merchantBusinessCommandRepository) RestoreMerchantBusiness(ctx context.Context, merchant_id int) (*record.MerchantBusinessRecord, error) {
	res, err := r.db.RestoreMerchantBusinessInformation(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantbusiness_errors.ErrRestoreMerchantBusiness
	}

	return r.mapping.ToMerchantBusinessRecord(res), nil
}

func (r *merchantBusinessCommandRepository) DeleteMerchantBusinessPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantBusinessInformationPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantbusiness_errors.ErrDeletePermanentMerchantBusiness
	}

	return true, nil
}

func (r *merchantBusinessCommandRepository) RestoreAllMerchantBusiness(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrRestoreAllMerchantBusinesses
	}
	return true, nil
}

func (r *merchantBusinessCommandRepository) DeleteAllMerchantBusinessPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		return false, merchantbusiness_errors.ErrDeleteAllPermanentMerchantBusinesses
	}
	return true, nil
}
