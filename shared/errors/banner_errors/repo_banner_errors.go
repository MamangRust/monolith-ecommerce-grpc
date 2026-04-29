package banner_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrBannerStartDate = errors.ErrBadRequest.WithMessage("banner start date must be less than end date")
	ErrBannerEndDate   = errors.ErrBadRequest.WithMessage("banner end date must be greater than start date")
	ErrBannerStartTime = errors.NewBadRequestError("banner start time must be less than end time")
	ErrBannerEndTime   = errors.NewBadRequestError("banner end time must be greater than start time")

	ErrBannerNotFound     = errors.ErrNotFound.WithMessage("banner not found")
	ErrFindAllBanners     = errors.ErrInternal.WithMessage("failed to find all banners")
	ErrFindActiveBanners  = errors.ErrInternal.WithMessage("failed to find active banners")
	ErrFindTrashedBanners = errors.ErrInternal.WithMessage("failed to find trashed banners")
	ErrBannerConflict     = errors.ErrConflict.WithMessage("failed banner already exists")

	ErrCreateBanner = errors.ErrInternal.WithMessage("failed to create banner")
	ErrUpdateBanner = errors.ErrInternal.WithMessage("failed to update banner")

	ErrTrashedBanner         = errors.ErrInternal.WithMessage("failed to move banner to trash")
	ErrRestoreBanner         = errors.ErrInternal.WithMessage("failed to restore banner from trash")
	ErrDeleteBannerPermanent = errors.ErrInternal.WithMessage("failed to permanently delete banner")

	ErrRestoreAllBanners = errors.ErrInternal.WithMessage("failed to restore all banners")
	ErrDeleteAllBanners  = errors.ErrInternal.WithMessage("failed to permanently delete all banners")
)

