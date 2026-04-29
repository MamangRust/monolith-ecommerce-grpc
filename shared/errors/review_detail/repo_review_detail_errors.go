package review_detail_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrReviewDetailNotFound     = errors.ErrNotFound.WithMessage("review detail not found")
	ErrFindAllReviewDetails     = errors.ErrInternal.WithMessage("failed to find all review details")
	ErrFindActiveReviewDetails  = errors.ErrInternal.WithMessage("failed to find active review details")
	ErrFindTrashedReviewDetails = errors.ErrInternal.WithMessage("failed to find trashed review details")
	ErrReviewDetailConflict     = errors.ErrConflict.WithMessage("failed, review detail already exists")

	ErrFindByIdReviewDetail        = errors.ErrInternal.WithMessage("failed to find review detail by ID")
	ErrFindByIdTrashedReviewDetail = errors.ErrInternal.WithMessage("failed to find trashed review detail by ID")

	ErrCreateReviewDetail = errors.ErrInternal.WithMessage("failed to create review detail")
	ErrUpdateReviewDetail = errors.ErrInternal.WithMessage("failed to update review detail")

	ErrTrashedReviewDetail         = errors.ErrInternal.WithMessage("failed to move review detail to trash")
	ErrRestoreReviewDetail         = errors.ErrInternal.WithMessage("failed to restore review detail from trash")
	ErrDeleteReviewDetailPermanent = errors.ErrInternal.WithMessage("failed to permanently delete review detail")

	ErrRestoreAllReviewDetails = errors.ErrInternal.WithMessage("failed to restore all review details")
	ErrDeleteAllReviewDetails  = errors.ErrInternal.WithMessage("failed to permanently delete all review details")

	ErrReviewDetailInternal = errors.ErrInternal.WithMessage("review detail internal repository error")
)
