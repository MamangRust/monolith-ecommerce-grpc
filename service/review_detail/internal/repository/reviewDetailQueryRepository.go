package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type reviewDetailQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.ReviewDetailRecordMapping
}

func NewReviewDetailQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.ReviewDetailRecordMapping) *reviewDetailQueryRepository {
	return &reviewDetailQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *reviewDetailQueryRepository) FindAllReviews(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetails(r.ctx, reqDb)

	if err != nil {
		return nil, nil, reviewdetail_errors.ErrFindAllReviewDetails
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewDetailsRecordPagination(res), &totalCount, nil
}

func (r *reviewDetailQueryRepository) FindByActive(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, reviewdetail_errors.ErrFindActiveReviewDetails
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewDetailsRecordActivePagination(res), &totalCount, nil
}

func (r *reviewDetailQueryRepository) FindByTrashed(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetReviewDetailsTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetReviewDetailsTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, reviewdetail_errors.ErrFindTrashedReviewDetails
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToReviewDetailsRecordTrashedPagination(res), &totalCount, nil
}

func (r *reviewDetailQueryRepository) FindById(user_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.GetReviewDetail(r.ctx, int32(user_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrFindByIdReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}

func (r *reviewDetailQueryRepository) FindByIdTrashed(user_id int) (*record.ReviewDetailRecord, error) {
	res, err := r.db.GetReviewDetailTrashed(r.ctx, int32(user_id))

	if err != nil {
		return nil, reviewdetail_errors.ErrFindByIdTrashedReviewDetail
	}

	return r.mapping.ToReviewDetailRecord(res), nil
}
