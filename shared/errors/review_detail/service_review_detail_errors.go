package review_detail_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedImageNotFound = errors.ErrNotFound.WithMessage("image not found")
	ErrFailedRemoveImage   = errors.ErrInternal.WithMessage("failed to remove image")

	ErrFailedFindReviewDetail = errors.ErrNotFound.WithMessage("review detail not found")
	ErrFailedFindAllReviewDetails     = errors.ErrInternal.WithMessage("failed to find all review details")
	ErrFailedFindActiveReviewDetails  = errors.ErrInternal.WithMessage("failed to find active review details")
	ErrFailedFindTrashedReviewDetails = errors.ErrInternal.WithMessage("failed to find trashed review details")

	ErrFailedCreateReviewDetail = errors.ErrInternal.WithMessage("failed to create review detail")
	ErrFailedUpdateReviewDetail = errors.ErrInternal.WithMessage("failed to update review detail")

	ErrFailedTrashedReviewDetail         = errors.ErrInternal.WithMessage("failed to move review detail to trash")
	ErrFailedRestoreReviewDetail         = errors.ErrInternal.WithMessage("failed to restore review detail from trash")
	ErrFailedDeleteReviewDetailPermanent = errors.ErrInternal.WithMessage("failed to permanently delete review detail")

	ErrFailedRestoreAllReviewDetails = errors.ErrInternal.WithMessage("failed to restore all review details")
	ErrFailedDeleteAllReviewDetails  = errors.ErrInternal.WithMessage("failed to permanently delete all review details")
)
