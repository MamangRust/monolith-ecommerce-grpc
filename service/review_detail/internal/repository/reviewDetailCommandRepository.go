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
	ctx     context.Context
	mapping recordmapper.ReviewDetailRecordMapping
}

func NewReviewDetailCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ReviewDetailRecordMapping) *reviewDetailCommandRepository {
	return &reviewDetailCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *reviewDetailCommandRepository) CreateReviewDetail(request *requests.CreateReviewDetailRequest) (*record.ReviewDetailRecord, error) {
	req := db.CreateReviewDetailParams{
		ReviewID: int32(request.ReviewID),
		Type:     request.Type,
		Url:      request.Url,
		Caption:  sql.NullString{String: request.Caption, Valid: request.Caption != ""},
	}

	reviewDetail, err := r.db.CreateReviewDetail(r.ctx, req)
	if err != nil {
		return nil, reviewdetail_errors.ErrCreateReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(reviewDetail), nil
}

func (r *reviewDetailCommandRepository) UpdateReviewDetail(request *requests.UpdateReviewDetailRequest) (*record.ReviewDetailRecord, error) {
	req := db.UpdateReviewDetailParams{
		ReviewDetailID: int32(*request.ReviewDetailID),
		Type:           request.Type,
		Url:            request.Url,
		Caption:        sql.NullString{String: request.Caption, Valid: request.Caption != ""},
	}

	res, err := r.db.UpdateReviewDetail(r.ctx, req)
	if err != nil {
		return nil, reviewdetail_errors.ErrUpdateReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailCommandRepository) TrashedReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.TrashReviewDetail(r.ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrTrashedReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailCommandRepository) RestoreReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.RestoreReviewDetail(r.ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrRestoreReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailCommandRepository) DeleteReviewDetailPermanent(ReviewDetail_id int) (bool, error) {
	err := r.db.DeletePermanentReviewDetail(r.ctx, int32(ReviewDetail_id))

	if err != nil {
		return false, reviewdetail_errors.ErrDeleteReviewDetailPermanent
	}

	return true, nil
}

func (r *reviewDetailCommandRepository) RestoreAllReviewDetail() (bool, error) {
	err := r.db.RestoreAllReviewDetails(r.ctx)

	if err != nil {
		return false, reviewdetail_errors.ErrRestoreAllReviewDetails
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) DeleteAllReviewDetailPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentReviewDetails(r.ctx)

	if err != nil {
		return false, reviewdetail_errors.ErrDeleteAllReviewDetails
	}
	return true, nil
}
