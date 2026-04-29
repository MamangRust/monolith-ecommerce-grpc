package review_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedReviewNotFound = errors.ErrNotFound.WithMessage("Review not found")

	ErrFailedFindAllReviews        = errors.ErrInternal.WithMessage("Failed to fetch all reviews")
	ErrFailedFindActiveReviews     = errors.ErrInternal.WithMessage("Failed to fetch active reviews")
	ErrFailedFindTrashedReviews    = errors.ErrInternal.WithMessage("Failed to fetch trashed reviews")
	ErrFailedFindByProductReviews  = errors.ErrInternal.WithMessage("Failed to fetch reviews by product")
	ErrFailedFindByMerchantReviews = errors.ErrInternal.WithMessage("Failed to fetch reviews by merchant")

	ErrFailedCreateReview = errors.ErrInternal.WithMessage("Failed to create review")
	ErrFailedUpdateReview = errors.ErrInternal.WithMessage("Failed to update review")

	ErrFailedTrashedReview         = errors.ErrInternal.WithMessage("Failed to move review to trash")
	ErrFailedRestoreReview         = errors.ErrInternal.WithMessage("Failed to restore review from trash")
	ErrFailedDeletePermanentReview = errors.ErrInternal.WithMessage("Failed to permanently delete review")

	ErrFailedRestoreAllReviews         = errors.ErrInternal.WithMessage("Failed to restore all reviews")
	ErrFailedDeleteAllPermanentReviews = errors.ErrInternal.WithMessage("Failed to permanently delete all reviews")
)
