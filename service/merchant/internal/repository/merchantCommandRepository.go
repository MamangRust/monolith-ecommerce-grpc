package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantCommandRepository(db *db.Queries, mapping recordmapper.MerchantRecordMapping) *merchantCommandRepository {
	return &merchantCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Address:      sql.NullString{String: request.Address, Valid: true},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: true},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: true},
		Status:       "inactive",
	}

	merchant, err := r.db.CreateMerchant(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrCreateMerchant
	}

	return r.mapping.ToMerchantRecord(merchant), nil
}

func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID:   int32(*request.MerchantID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Address:      sql.NullString{String: request.Address, Valid: true},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: true},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: true},
		Status:       request.Status,
	}

	res, err := r.db.UpdateMerchant(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.TrashMerchant(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrTrashedMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.RestoreMerchant(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanent
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)

	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchant
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)

	if err != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanent
	}
	return true, nil
}
