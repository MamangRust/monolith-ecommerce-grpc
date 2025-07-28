package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*record.UserRecord, error)
}

type ProductQueryRepository interface {
	FindById(ctx context.Context, user_id int) (*record.ProductRecord, error)
}

type ReviewQueryRepository interface {
	FindAllReview(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindByProduct(ctx context.Context, req *requests.FindAllReviewByProduct) ([]*record.ReviewsDetailRecord, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllReviewByMerchant) ([]*record.ReviewsDetailRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllReview) ([]*record.ReviewRecord, *int, error)
	FindById(ctx context.Context, id int) (*record.ReviewRecord, error)
}

type ReviewCommandRepository interface {
	CreateReview(ctx context.Context, request *requests.CreateReviewRequest) (*record.ReviewRecord, error)
	UpdateReview(ctx context.Context, request *requests.UpdateReviewRequest) (*record.ReviewRecord, error)
	TrashReview(ctx context.Context, shipping_id int) (*record.ReviewRecord, error)
	RestoreReview(ctx context.Context, category_id int) (*record.ReviewRecord, error)
	DeleteReviewPermanently(ctx context.Context, category_id int) (bool, error)
	RestoreAllReview(ctx context.Context) (bool, error)
	DeleteAllPermanentReview(ctx context.Context) (bool, error)
}
