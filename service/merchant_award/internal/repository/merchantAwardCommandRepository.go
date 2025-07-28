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
	mapping recordmapper.MerchantAwardMapping
}

func NewMerchantAwardCommandRepository(db *db.Queries, mapping recordmapper.MerchantAwardMapping) *merchantAwardCommandRepository {
	return &merchantAwardCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantAwardCommandRepository) CreateMerchantAward(ctx context.Context, request *requests.CreateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error) {
	req := db.CreateMerchantCertificationOrAwardParams{
		MerchantID:     int32(request.MerchantID),
		Title:          request.Title,
		Description:    sql.NullString{String: request.Description, Valid: request.Description != ""},
		IssuedBy:       sql.NullString{String: request.IssuedBy, Valid: request.IssuedBy != ""},
		IssueDate:      parseDateToNullTime(request.IssueDate),
		ExpiryDate:     parseDateToNullTime(request.ExpiryDate),
		CertificateUrl: sql.NullString{String: request.CertificateUrl, Valid: request.CertificateUrl != ""},
	}

	award, err := r.db.CreateMerchantCertificationOrAward(ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(award), nil
}

func (r *merchantAwardCommandRepository) UpdateMerchantAward(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*record.MerchantAwardRecord, error) {
	req := db.UpdateMerchantCertificationOrAwardParams{
		MerchantCertificationID: int32(*request.MerchantCertificationID),
		Title:                   request.Title,
		Description:             sql.NullString{String: request.Description, Valid: request.Description != ""},
		IssuedBy:                sql.NullString{String: request.IssuedBy, Valid: request.IssuedBy != ""},
		IssueDate:               parseDateToNullTime(request.IssueDate),
		ExpiryDate:              parseDateToNullTime(request.ExpiryDate),
		CertificateUrl:          sql.NullString{String: request.CertificateUrl, Valid: request.CertificateUrl != ""},
	}

	res, err := r.db.UpdateMerchantCertificationOrAward(ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardCommandRepository) TrashedMerchantAward(ctx context.Context, merchant_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.TrashMerchantCertificationOrAward(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrTrashedMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardCommandRepository) RestoreMerchantAward(ctx context.Context, merchant_id int) (*record.MerchantAwardRecord, error) {
	res, err := r.db.RestoreMerchantCertificationOrAward(ctx, int32(merchant_id))

	if err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward
	}

	return r.mapping.ToMerchantAwardRecord(res), nil
}

func (r *merchantAwardCommandRepository) DeleteMerchantPermanent(ctx context.Context, Merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantCertificationOrAwardPermanently(ctx, int32(Merchant_id))

	if err != nil {
		return false, merchantaward_errors.ErrDeleteMerchantAwardPermanent
	}

	return true, nil
}

func (r *merchantAwardCommandRepository) RestoreAllMerchantAward(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantCertificationsAndAwards(ctx)

	if err != nil {
		return false, merchantaward_errors.ErrRestoreAllMerchantAwards
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) DeleteAllMerchantAwardPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantCertificationsAndAwards(ctx)

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
