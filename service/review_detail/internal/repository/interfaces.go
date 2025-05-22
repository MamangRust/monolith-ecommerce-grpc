package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ReviewDetailQueryRepository interface {
	FindAllReviews(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindByActive(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindByTrashed(req *requests.FindAllReview) ([]*record.ReviewDetailRecord, *int, error)
	FindById(user_id int) (*record.ReviewDetailRecord, error)
	FindByIdTrashed(user_id int) (*record.ReviewDetailRecord, error)
}

type ReviewDetailCommandRepository interface {
	CreateReviewDetail(request *requests.CreateReviewDetailRequest) (*record.ReviewDetailRecord, error)
	UpdateReviewDetail(request *requests.UpdateReviewDetailRequest) (*record.ReviewDetailRecord, error)
	TrashedReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error)
	RestoreReviewDetail(ReviewDetail_id int) (*record.ReviewDetailRecord, error)
	DeleteReviewDetailPermanent(ReviewDetail_id int) (bool, error)
	RestoreAllReviewDetail() (bool, error)
	DeleteAllReviewDetailPermanent() (bool, error)
}
