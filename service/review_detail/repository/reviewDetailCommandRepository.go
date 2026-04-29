package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_detail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
)

type reviewDetailCommandRepository struct {
	db *db.Queries
}

func NewReviewDetailCommandRepository(db *db.Queries) *reviewDetailCommandRepository {
	return &reviewDetailCommandRepository{
		db: db,
	}
}

func (r *reviewDetailCommandRepository) Create(ctx context.Context, request *requests.CreateReviewDetailRequest) (*db.CreateReviewDetailRow, error) {
	req := db.CreateReviewDetailParams{
		ReviewID: int32(request.ReviewID),
		Type:     request.Type,
		Url:      request.Url,
		Caption:  stringPtr(request.Caption),
	}

	reviewDetail, err := r.db.CreateReviewDetail(ctx, req)
	if err != nil {
		return nil, review_detail_errors.ErrCreateReviewDetail.WithInternal(err)
	}

	return reviewDetail, nil
}

func (r *reviewDetailCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewDetailRequest) (*db.UpdateReviewDetailRow, error) {
	req := db.UpdateReviewDetailParams{
		ReviewDetailID: int32(*request.ReviewDetailID),
		Type:           request.Type,
		Url:            request.Url,
		Caption:        stringPtr(request.Caption),
	}

	res, err := r.db.UpdateReviewDetail(ctx, req)
	if err != nil {
		return nil, review_detail_errors.ErrUpdateReviewDetail.WithInternal(err)
	}

	return res, nil
}

func (r *reviewDetailCommandRepository) Trash(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error) {
	res, err := r.db.TrashReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, review_detail_errors.ErrTrashedReviewDetail.WithInternal(err)
	}

	return res, nil
}

func (r *reviewDetailCommandRepository) Restore(ctx context.Context, ReviewDetail_id int) (*db.ReviewDetail, error) {
	res, err := r.db.RestoreReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return nil, review_detail_errors.ErrRestoreReviewDetail.WithInternal(err)
	}

	return res, nil
}

func (r *reviewDetailCommandRepository) DeletePermanent(ctx context.Context, ReviewDetail_id int) (bool, error) {
	err := r.db.DeletePermanentReviewDetail(ctx, int32(ReviewDetail_id))

	if err != nil {
		return false, review_detail_errors.ErrDeleteReviewDetailPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *reviewDetailCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllReviewDetails(ctx)

	if err != nil {
		return false, review_detail_errors.ErrRestoreAllReviewDetails.WithInternal(err)
	}
	return true, nil
}

func (r *reviewDetailCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentReviewDetails(ctx)

	if err != nil {
		return false, review_detail_errors.ErrDeleteAllReviewDetails.WithInternal(err)
	}
	return true, nil
}

func stringPtr(s string) *string {
	return &s
}
