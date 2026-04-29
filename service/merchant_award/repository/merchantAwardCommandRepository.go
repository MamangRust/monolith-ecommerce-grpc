package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	"github.com/jackc/pgx/v5/pgtype"
)

type merchantAwardCommandRepository struct {
	db *db.Queries
}

func NewMerchantAwardCommandRepository(db *db.Queries) *merchantAwardCommandRepository {
	return &merchantAwardCommandRepository{
		db: db,
	}
}

func (r *merchantAwardCommandRepository) Create(
	ctx context.Context,
	request *requests.CreateMerchantCertificationOrAwardRequest,
) (*db.CreateMerchantCertificationOrAwardRow, error) {

	req := db.CreateMerchantCertificationOrAwardParams{
		MerchantID: int32(request.MerchantID),
		Title:      request.Title,

		Description:    stringPtr(request.Description),
		IssuedBy:       stringPtr(request.IssuedBy),
		CertificateUrl: stringPtr(request.CertificateUrl),

		IssueDate:  parseDateToPgDate(request.IssueDate),
		ExpiryDate: parseDateToPgDate(request.ExpiryDate),
	}

	award, err := r.db.CreateMerchantCertificationOrAward(ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrCreateMerchantAward.WithInternal(err)
	}

	return award, nil
}

func (r *merchantAwardCommandRepository) Update(ctx context.Context, request *requests.UpdateMerchantCertificationOrAwardRequest) (*db.UpdateMerchantCertificationOrAwardRow, error) {
	req := db.UpdateMerchantCertificationOrAwardParams{
		MerchantCertificationID: int32(*request.MerchantCertificationID),
		Title:                   request.Title,
		Description:             stringPtr(request.Description),
		IssuedBy:                stringPtr(request.IssuedBy),
		CertificateUrl:          stringPtr(request.CertificateUrl),
		IssueDate:               parseDateToPgDate(request.IssueDate),
		ExpiryDate:              parseDateToPgDate(request.ExpiryDate),
	}

	res, err := r.db.UpdateMerchantCertificationOrAward(ctx, req)
	if err != nil {
		return nil, merchantaward_errors.ErrUpdateMerchantAward.WithInternal(err)
	}

	return res, nil
}

func (r *merchantAwardCommandRepository) Trash(ctx context.Context, award_id int) (*db.MerchantCertificationsAndAward, error) {
	res, err := r.db.TrashMerchantCertificationOrAward(ctx, int32(award_id))

	if err != nil {
		return nil, merchantaward_errors.ErrTrashedMerchantAward.WithInternal(err)
	}

	return res, nil
}

func (r *merchantAwardCommandRepository) Restore(ctx context.Context, award_id int) (*db.MerchantCertificationsAndAward, error) {
	res, err := r.db.RestoreMerchantCertificationOrAward(ctx, int32(award_id))

	if err != nil {
		return nil, merchantaward_errors.ErrRestoreMerchantAward.WithInternal(err)
	}

	return res, nil
}

func (r *merchantAwardCommandRepository) DeletePermanent(ctx context.Context, award_id int) (bool, error) {
	err := r.db.DeleteMerchantCertificationOrAwardPermanently(ctx, int32(award_id))

	if err != nil {
		return false, merchantaward_errors.ErrDeleteMerchantAwardPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *merchantAwardCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchantCertificationsAndAwards(ctx)

	if err != nil {
		return false, merchantaward_errors.ErrRestoreAllMerchantAwards.WithInternal(err)
	}
	return true, nil
}

func (r *merchantAwardCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchantCertificationsAndAwards(ctx)

	if err != nil {
		return false, merchantaward_errors.ErrDeleteAllMerchantAwardsPermanent.WithInternal(err)
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

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseDateToPgDate(dateStr string) pgtype.Date {
	if dateStr == "" {
		return pgtype.Date{Valid: false}
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return pgtype.Date{Valid: false}
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}
