package repository

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type reviewDetailCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.ReviewDetailRecordMapping
}

func NewReviewDetailCommandRepository(db *db.Queries, mapping recordmapper.ReviewDetailRecordMapping) *reviewDetailCommandRepository {
	return &reviewDetailCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *reviewDetailCommandRepository) CreateReviewDetail(ctx context.Context, request *requests.CreateReviewDetailRequest) (*record.ReviewDetailRecord, error) {
	req := db.CreateReviewDetailParams{
		ReviewID: int32(request.ReviewID),
		Type:     request.Type,
		Url:      request.Url,
		Caption:  sql.NullString{String: request.Caption, Valid: request.Caption != ""},
	}

	reviewDetail, err := r.db.CreateReviewDetail(ctx, req)
	if err != nil {
		return nil, reviewdetail_errors.ErrCreateReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(reviewDetail), nil
}

func (r *reviewDetailCommandRepository) UpdateReviewDetail(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*record.ReviewDetailRecord, error) {
	req := db.UpdateReviewDetailParams{
		ReviewDetailID: int32(*request.ReviewDetailID),
		Type:           request.Type,
		Url:            request.Url,
		Caption:        sql.NullString{String: request.Caption, Valid: request.Caption != ""},
	}

	res, err := r.db.UpdateReviewDetail(ctx, req)
	if err != nil {
		return nil, reviewdetail_errors.ErrUpdateReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailCommandRepository) TrashedReviewDetail(ctx context.Context, ReviewDetail_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.TrashReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrTrashedReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailCommandRepository) RestoreReviewDetail(ctx context.Context, ReviewDetail_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.RestoreReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrRestoreReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailCommandRepository) DeleteReviewDetailPermanent(ctx context.Context, ReviewDetail_id int) (bool, error) {
	err := r.db.DeletePermanentReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return false, reviewdetail_errors.ErrDeleteReviewDetailPermanent
	}

	return true, nil
}

func (r *reviewDetailCommandRepository) RestoreAllReviewDetail(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllReviewDetails(ctx)

	if err != nil {
		return false, reviewdetail_errors.ErrRestoreAllReviewDetails
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) DeleteAllReviewDetailPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentReviewDetails(ctx)

	if err != nil {
		return false, reviewdetail_errors.ErrDeleteAllReviewDetails
	}
	return true, nil
}
