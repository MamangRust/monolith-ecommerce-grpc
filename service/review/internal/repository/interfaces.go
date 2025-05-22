package repository

import (
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindById(user_id int) (*record.UserRecord, error)
}

type ProductQueryRepository interface {
	FindById(user_id int) (*record.ProductRecord, error)
}

type ReviewQueryRepository interface {
	FindAllReview(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindByProduct(req *requests.FindAllReviewByProduct) ([]*record.ReviewsDetailRecord, *int, error)
	FindByMerchant(req *requests.FindAllReviewByMerchant) ([]*record.ReviewsDetailRecord, *int, error)
	FindByActive(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindByTrashed(req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindById(id int) (*record.ReviewRecord, error)
}

type ReviewCommandRepository interface {
	CreateReview(request *requests.CreateReviewRequest) (*record.ReviewRecord, error)
	UpdateReview(request *requests.UpdateReviewRequest) (*record.ReviewRecord, error)
	TrashReview(shipping_id int) (*record.ReviewRecord, error)
	RestoreReview(category_id int) (*record.ReviewRecord, error)
	DeleteReviewPermanently(category_id int) (bool, error)
	RestoreAllReview() (bool, error)
	DeleteAllPermanentReview() (bool, error)
}
