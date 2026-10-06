package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
)

type reviewCommandRepository struct {
	db *db.Queries
}

func NewReviewCommandRepository(db *db.Queries) *reviewCommandRepository {
	return &reviewCommandRepository{
		db: db,
	}
}

func (r *reviewCommandRepository) Create(ctx context.Context, request *requests.CreateReviewRequest) (*db.CreateReviewRow, error) {
	req := db.CreateReviewParams{
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),
		Rating:    int32(request.Rating),
		Comment:   request.Comment,
	}

	review, err := r.db.CreateReview(ctx, req)

	if err != nil {
		return nil, review_errors.ErrCreateReview.WithInternal(err)
	}

	return review, nil
}

func (r *reviewCommandRepository) Update(ctx context.Context, request *requests.UpdateReviewRequest) (*db.UpdateReviewRow, error) {
	req := db.UpdateReviewParams{
		ReviewID: int32(*request.ReviewID),
		Name:     request.Name,
		Rating:   int32(request.Rating),
		Comment:  request.Comment,
	}

	res, err := r.db.UpdateReview(ctx, req)

	if err != nil {
		return nil, review_errors.ErrUpdateReview.WithInternal(err)
	}

	return res, nil
}

func (r *reviewCommandRepository) Trash(ctx context.Context, review_id int) (*db.Review, error) {
	res, err := r.db.TrashReview(ctx, int32(review_id))

	if err != nil {
		return nil, review_errors.ErrTrashReview.WithInternal(err)
	}

	return res, nil
}

func (r *reviewCommandRepository) Restore(ctx context.Context, review_id int) (*db.Review, error) {
	res, err := r.db.RestoreReview(ctx, int32(review_id))

	if err != nil {
		return nil, review_errors.ErrRestoreReview.WithInternal(err)
	}

	return res, nil
}

func (r *reviewCommandRepository) DeletePermanent(ctx context.Context, review_id int) (bool, error) {
	err := r.db.DeleteReviewPermanently(ctx, int32(review_id))

	if err != nil {
		return false, review_errors.ErrDeleteReviewPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *reviewCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllReviews(ctx)

	if err != nil {
		return false, review_errors.ErrRestoreAllReviews.WithInternal(err)
	}
	return true, nil
}

func (r *reviewCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentReviews(ctx)

	if err != nil {
		return false, review_errors.ErrDeleteAllPermanentReview.WithInternal(err)
	}
	return true, nil
}
