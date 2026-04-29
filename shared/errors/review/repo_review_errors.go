package review_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrReviewNotFound            = errors.ErrNotFound.WithMessage("review not found")
	ErrFindAllReviews            = errors.ErrInternal.WithMessage("failed to find all reviews")
	ErrFindReviewsByProduct      = errors.ErrInternal.WithMessage("failed to find reviews by product")
	ErrFindReviewsByMerchant     = errors.ErrInternal.WithMessage("failed to find reviews by merchant")
	ErrFindActiveReviews         = errors.ErrInternal.WithMessage("failed to find active reviews")
	ErrFindTrashedReviews        = errors.ErrInternal.WithMessage("failed to find trashed reviews")
	ErrFindReviewByID            = errors.ErrInternal.WithMessage("failed to find review by ID")
	ErrCreateReview              = errors.ErrInternal.WithMessage("failed to create review")
	ErrUpdateReview              = errors.ErrInternal.WithMessage("failed to update review")
	ErrTrashReview               = errors.ErrInternal.WithMessage("failed to move review to trash")
	ErrRestoreReview             = errors.ErrInternal.WithMessage("failed to restore review from trash")
	ErrDeleteReviewPermanent     = errors.ErrInternal.WithMessage("failed to permanently delete review")
	ErrRestoreAllReviews         = errors.ErrInternal.WithMessage("failed to restore all reviews")
	ErrDeleteAllPermanentReview = errors.ErrInternal.WithMessage("failed to permanently delete all reviews")

	ErrReviewInternal = errors.ErrInternal.WithMessage("review internal repository error")
)
