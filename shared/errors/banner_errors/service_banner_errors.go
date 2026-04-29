package banner_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrBannerNotFoundRes = errors.ErrNotFound.WithMessage("Banner not found")
	ErrBannerInvalidData = errors.ErrBadRequest.WithMessage("Invalid banner data")

	ErrFailedFindAllBanners     = errors.ErrInternal.WithMessage("Failed to fetch banners")
	ErrFailedFindActiveBanners  = errors.ErrInternal.WithMessage("Failed to fetch active banners")
	ErrFailedFindTrashedBanners = errors.ErrInternal.WithMessage("Failed to fetch trashed banners")

	ErrFailedCreateBanner = errors.ErrInternal.WithMessage("Failed to create banner")
	ErrFailedUpdateBanner = errors.ErrInternal.WithMessage("Failed to update banner")

	ErrFailedTrashedBanner = errors.ErrInternal.WithMessage("Failed to move banner to trash")
	ErrFailedRestoreBanner = errors.ErrInternal.WithMessage("Failed to restore banner")
	ErrFailedDeleteBanner  = errors.ErrInternal.WithMessage("Failed to permanently delete banner")

	ErrFailedRestoreAllBanners = errors.ErrInternal.WithMessage("Failed to restore all banners")
	ErrFailedDeleteAllBanners  = errors.ErrInternal.WithMessage("Failed to permanently delete all banners")
)

