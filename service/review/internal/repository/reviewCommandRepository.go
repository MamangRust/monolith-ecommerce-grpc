package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type reviewCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.ReviewRecordMapping
}

func NewReviewCommandRepository(db *db.Queries, mapping recordmapper.ReviewRecordMapping) *reviewCommandRepository {
	return &reviewCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *reviewCommandRepository) CreateReview(ctx context.Context, request *requests.CreateReviewRequest) (*record.ReviewRecord, error) {
	req := db.CreateReviewParams{
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),
		Rating:    int32(request.Rating),
		Comment:   request.Comment,
	}

	review, err := r.db.CreateReview(ctx, req)

	if err != nil {
		return nil, review_errors.ErrCreateReview
	}

	return r.mapping.ToReviewRecord(review), nil
}

func (r *reviewCommandRepository) UpdateReview(ctx context.Context, request *requests.UpdateReviewRequest) (*record.ReviewRecord, error) {
	req := db.UpdateReviewParams{
		ReviewID: int32(*request.ReviewID),
		Name:     request.Name,
		Rating:   int32(request.Rating),
		Comment:  request.Comment,
	}

	res, err := r.db.UpdateReview(ctx, req)

	if err != nil {
		return nil, review_errors.ErrUpdateReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewCommandRepository) TrashReview(ctx context.Context, shipping_id int) (*record.ReviewRecord, error) {
	res, err := r.db.TrashReview(ctx, int32(shipping_id))

	if err != nil {
		return nil, review_errors.ErrTrashReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewCommandRepository) RestoreReview(ctx context.Context, category_id int) (*record.ReviewRecord, error) {
	res, err := r.db.RestoreReview(ctx, int32(category_id))

	if err != nil {
		return nil, review_errors.ErrRestoreReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewCommandRepository) DeleteReviewPermanently(ctx context.Context, category_id int) (bool, error) {
	err := r.db.DeleteReviewPermanently(ctx, int32(category_id))

	if err != nil {
		return false, review_errors.ErrDeleteReviewPermanent
	}

	return true, nil
}

func (r *reviewCommandRepository) RestoreAllReview(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllReviews(ctx)

	if err != nil {
		return false, review_errors.ErrRestoreAllReviews
	}
	return true, nil
}

func (r *reviewCommandRepository) DeleteAllPermanentReview(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentReviews(ctx)

	if err != nil {
		return false, review_errors.ErrDeleteAllPermanentReview
	}
	return true, nil
}
