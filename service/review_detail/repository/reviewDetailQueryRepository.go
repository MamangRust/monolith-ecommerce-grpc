package repository

import (
	"context"

	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_detail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
)

type reviewDetailQueryRepository struct {
	db *db.Queries
}

func NewReviewDetailQueryRepository(db *db.Queries) *reviewDetailQueryRepository {
	return &reviewDetailQueryRepository{
		db: db,
	}
}

func (r *reviewDetailQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetails(ctx, reqDb)

	if err != nil {
		return nil, review_detail_errors.ErrFindAllReviewDetails.WithInternal(err)
	}

	return res, nil
}

func (r *reviewDetailQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsActive(ctx, reqDb)

	if err != nil {
		return nil, review_detail_errors.ErrFindActiveReviewDetails.WithInternal(err)
	}

	return res, nil
}

func (r *reviewDetailQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewDetailsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsTrashed(ctx, reqDb)

	if err != nil {
		return nil, review_detail_errors.ErrFindTrashedReviewDetails.WithInternal(err)
	}

	return res, nil
}

func (r *reviewDetailQueryRepository) FindByID(ctx context.Context, user_id int) (*db.GetReviewDetailRow, error) {
	res, err := r.db.GetReviewDetail(ctx, int32(user_id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, review_detail_errors.ErrReviewDetailNotFound
		}
		return nil, review_detail_errors.ErrFindByIdReviewDetail.WithInternal(err)
	}

	return res, nil
}

func (r *reviewDetailQueryRepository) FindByIDTrashed(ctx context.Context, user_id int) (*db.ReviewDetail, error) {
	res, err := r.db.GetReviewDetailTrashed(ctx, int32(user_id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, review_detail_errors.ErrReviewDetailNotFound
		}
		return nil, review_detail_errors.ErrFindByIdTrashedReviewDetail.WithInternal(err)
	}

	return res, nil
}
