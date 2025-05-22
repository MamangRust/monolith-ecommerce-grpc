package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type merchantAwardCommandRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.MerchantAwardMapping
}

func NewMerchantAwardCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.MerchantAwardMapping) *merchantAwardCommandRepository {
	return &merchantAwardCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *merchantAwardCommandRepository) CreateMerchantAward(request *requests.CreateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error) {
	req := db.CreateMerchantCertificationOrAwardParams{
		MerchantID:     int32(request.MerchantID),
		Title:          request.Title,
		Description:    sql.NullString{String: request.Description, Valid: request.Description != ""},
		IssuedBy:       sql.NullString{String: request.IssuedBy, Valid: request.IssuedBy != ""},
		IssueDate:      parseDateToNullTime(request.IssueDate),
		ExpiryDate:     parseDateToNullTime(request.ExpiryDate),
		CertificateUrl: sql.NullString{String: request.CertificateUrl, Valid: request.CertificateUrl != ""},
	}

	award, err := r.db.CreateMerchantCertificationOrAward(r.ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(award), nil
}

func (r *merchantAwardCommandRepository) UpdateMerchantAward(request *requests.UpdateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error) {
	req := db.UpdateMerchantCertificationOrAwardParams{
		MerchantCertificationID: int32(*request.MerchantCertificationID),
		Title:                   request.Title,
		Description:             sql.NullString{String: request.Description, Valid: request.Description != ""},
		IssuedBy:                sql.NullString{String: request.IssuedBy, Valid: request.IssuedBy != ""},
		IssueDate:               parseDateToNullTime(request.IssueDate),
		ExpiryDate:              parseDateToNullTime(request.ExpiryDate),
		CertificateUrl:          sql.NullString{String: request.CertificateUrl, Valid: request.CertificateUrl != ""},
	}

	res, err := r.db.UpdateMerchantCertificationOrAward(r.ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardCommandRepository) TrashedMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.TrashMerchantCertificationOrAward(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrTrashedMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardCommandRepository) RestoreMerchantAward(merchant_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.RestoreMerchantCertificationOrAward(r.ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardCommandRepository) DeleteMerchantPermanent(Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantCertificationOrAwardPermanently(r.ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantaward_errors.ErrDeleteMerchantAwardPermanent
	}

	return true, nil
}

func (r *merchantAwardCommandRepository) RestoreAllMerchantAward() (bool, error) {
	err := r.db.RestoreAllMerchantCertificationsAndAwards(r.ctx)

	if err != nil {
		return false, merchantaward_errors.ErrRestoreAllMerchantAwards
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) DeleteAllMerchantAwardPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentMerchantCertificationsAndAwards(r.ctx)

	if err != nil {
		return false, merchantaward_errors.ErrDeleteAllMerchantAwardsPermanent
	}
	return true, nil
}

func parseDateToNullTime(dateStr string) sql.NullTime {
	if dateStr == "" {
		return sql.NullTime{Valid: false}
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: t, Valid: true}
}
