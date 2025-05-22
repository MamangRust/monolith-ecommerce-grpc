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
	ctx     context.Context
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantRecordMapping) *merchantCommandRepository {
	return &merchantCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantCommandRepository) CreateMerchant(request *requests.CreateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Address:      sql.NullString{String: request.Address, Valid: true},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: true},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: true},
		Status:       "inactive",
	}

	merchant, err := r.db.CreateMerchant(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrCreateMerchant
	}

	return r.mapping.ToMerchantRecord(merchant), nil
}

func (r *merchantCommandRepository) UpdateMerchant(request *requests.UpdateMerchantRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantParams{
		MerchantID:   int32(*request.MerchantID),
		Name:         request.Name,
		Description:  sql.NullString{String: request.Description, Valid: true},
		Address:      sql.NullString{String: request.Address, Valid: true},
		ContactEmail: sql.NullString{String: request.ContactEmail, Valid: true},
		ContactPhone: sql.NullString{String: request.ContactPhone, Valid: true},
		Status:       request.Status,
	}

	res, err := r.db.UpdateMerchant(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*record.MerchantRecord, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(r.ctx, req)

	if err != nil {
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) TrashedMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.TrashMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrTrashedMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) RestoreMerchant(merchant_id int) (*record.MerchantRecord, error) {
	res, err := r.db.RestoreMerchant(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchant_errors.ErrRestoreMerchant
	}

	return r.mapping.ToMerchantRecord(res), nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanent
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant() (bool, error) {
	err := r.db.RestoreAllMerchants(r.ctx)

	if err != nil {
		return false, merchant_errors.ErrRestoreAllMerchant
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(r.ctx)

	if err != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanent
	}
	return true, nil
}
