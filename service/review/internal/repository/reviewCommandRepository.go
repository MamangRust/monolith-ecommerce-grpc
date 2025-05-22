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
	ctx     context.Context
	mapping recordmapper.ReviewRecordMapping
}

func NewReviewCommandRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ReviewRecordMapping) *reviewCommandRepository {
	return &reviewCommandRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *reviewCommandRepository) CreateReview(request *requests.CreateReviewRequest) (*record.ReviewRecord, error) {
	req := db.CreateReviewParams{
		UserID:    int32(request.UserID),
		ProductID: int32(request.ProductID),
		Rating:    int32(request.Rating),
		Comment:   request.Comment,
	}

	review, err := r.db.CreateReview(r.ctx, req)

	if err != nil {
		return nil, review_errors.ErrCreateReview
	}

	return r.mapping.ToReviewRecord(review), nil
}

func (r *reviewCommandRepository) UpdateReview(request *requests.UpdateReviewRequest) (*record.ReviewRecord, error) {
	req := db.UpdateReviewParams{
		ReviewID: int32(*request.ReviewID),
		Name:     request.Name,
		Rating:   int32(request.Rating),
		Comment:  request.Comment,
	}

	res, err := r.db.UpdateReview(r.ctx, req)

	if err != nil {
		return nil, review_errors.ErrUpdateReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewCommandRepository) TrashReview(shipping_id int) (*record.ReviewRecord, error) {
	res, err := r.db.TrashReview(r.ctx, int32(shipping_id))

	if err != nil {
		return nil, review_errors.ErrTrashReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewCommandRepository) RestoreReview(category_id int) (*record.ReviewRecord, error) {
	res, err := r.db.RestoreReview(r.ctx, int32(category_id))

	if err != nil {
		return nil, review_errors.ErrRestoreReview
	}

	return r.mapping.ToReviewRecord(res), nil
}

func (r *reviewCommandRepository) DeleteReviewPermanently(category_id int) (bool, error) {
	err := r.db.DeleteReviewPermanently(r.ctx, int32(category_id))

	if err != nil {
		return false, review_errors.ErrDeleteReviewPermanent
	}

	return true, nil
}

func (r *reviewCommandRepository) RestoreAllReview() (bool, error) {
	err := r.db.RestoreAllReviews(r.ctx)

	if err != nil {
		return false, review_errors.ErrRestoreAllReviews
	}
	return true, nil
}

func (r *reviewCommandRepository) DeleteAllPermanentReview() (bool, error) {
	err := r.db.DeleteAllPermanentReviews(r.ctx)

	if err != nil {
		return false, review_errors.ErrDeleteAllPermanentReview
	}
	return true, nil
}
