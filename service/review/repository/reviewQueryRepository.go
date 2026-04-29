package repository

import (
	"context"

	"database/sql"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
)

type reviewQueryRepository struct {
	db *db.Queries
}

func NewReviewQueryRepository(db *db.Queries) *reviewQueryRepository {
	return &reviewQueryRepository{
		db: db,
	}
}

func (r *reviewQueryRepository) FindAll(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviews(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindAllReviews.WithInternal(err)
	}

	return res, nil
}

func (r *reviewQueryRepository) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*db.GetReviewByProductIdRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByProductIdParams{
		ProductID: int32(req.ProductID),
		Column2:   int32(req.Rating),
		Limit:     int32(req.PageSize),
		Offset:    int32(offset),
	}

	res, err := r.db.GetReviewByProductId(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindReviewsByProduct.WithInternal(err)
	}

	return res, nil
}

func (r *reviewQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*db.GetReviewByMerchantIdRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByMerchantIdParams{
		MerchantID: int32(req.MerchantID),
		Column2:    int32(req.Rating),
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetReviewByMerchantId(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindReviewsByMerchant.WithInternal(err)
	}

	return res, nil
}

func (r *reviewQueryRepository) FindActive(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsActive(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindActiveReviews.WithInternal(err)
	}

	return res, nil
}

func (r *reviewQueryRepository) FindTrashed(ctx context.Context, req *requests.FindAllReview) ([]*db.GetReviewsTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsTrashed(ctx, reqDb)

	if err != nil {
		return nil, review_errors.ErrFindTrashedReviews.WithInternal(err)
	}

	return res, nil
}

func (r *reviewQueryRepository) FindByID(ctx context.Context, id int) (*db.GetReviewByIDRow, error) {
	res, err := r.db.GetReviewByID(ctx, int32(id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, review_errors.ErrReviewNotFound
		}
		return nil, review_errors.ErrFindReviewByID.WithInternal(err)
	}

	return res, nil
}
