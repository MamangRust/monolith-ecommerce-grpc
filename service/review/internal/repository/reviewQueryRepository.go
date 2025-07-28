package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type reviewQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.ReviewRecordMapping
}

func NewReviewQueryRepository(db *db.Queries, mapping recordmapper.ReviewRecordMapping) *reviewQueryRepository {
	return &reviewQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *reviewQueryRepository) FindAllReview(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviews(ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindAllReviews
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordPagination(res), &totalCount, nil
}

func (r *reviewQueryRepository) FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*record.ReviewsDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByProductIdParams{
		ProductID: int32(req.ProductID),
		Column2:   int32(req.Rating),
		Limit:     int32(req.PageSize),
		Offset:    int32(offset),
	}

	res, err := r.db.GetReviewByProductId(ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindReviewsByProduct
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsProductRecordPagination(res), &totalCount, nil
}

func (r *reviewQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*record.ReviewsDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewByMerchantIdParams{
		MerchantID: int32(req.MerchantID),
		Column2:    int32(req.Rating),
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetReviewByMerchantId(ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindReviewsByMerchant
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsMerchantRecordPagination(res), &totalCount, nil
}

func (r *reviewQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsActive(ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindActiveReviews
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordActivePagination(res), &totalCount, nil
}

func (r *reviewQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewsTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, review_errors.ErrFindTrashedReviews
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewsRecordTrashedPagination(res), &totalCount, nil
}

func (r *reviewQueryRepository) FindById(ctx context.Context, id int) (*record.ReviewRecord, error) {
	res, err := r.db.GetReviewByID(ctx, int32(id))

	if err != nil {
		return nil, review_errors.ErrFindReviewByID
	}

	return r.mapping.ToReviewRecord(res), nil
}
