package merchantdetail_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedImageNotFound             = errors.ErrNotFound.WithMessage("Image not found")
	ErrFailedRemoveImageMerchantDetail = errors.ErrInternal.WithMessage("Failed to remove image merchant detail")
	ErrFailedLogoNotFound              = errors.ErrInternal.WithMessage("Failed to upload logo merchant detail")
	ErrFailedRemoveLogoMerchantDetail  = errors.ErrInternal.WithMessage("Failed to remove logo merchant detail")

	ErrFailedFindAllMerchantDetail            = errors.ErrInternal.WithMessage("Failed to find all merchant details")
	ErrFailedFindActiveMerchantDetail         = errors.ErrInternal.WithMessage("Failed to find active merchant details")
	ErrFailedFindTrashedMerchantDetail        = errors.ErrInternal.WithMessage("Failed to find trashed merchant details")
	ErrFailedFindMerchantDetailById           = errors.ErrInternal.WithMessage("Failed to find merchant detail by ID")
	ErrFailedCreateMerchantDetail             = errors.ErrInternal.WithMessage("Failed to create merchant detail")
	ErrFailedUpdateMerchantDetail             = errors.ErrInternal.WithMessage("Failed to update merchant detail")
	ErrFailedTrashedMerchantDetail            = errors.ErrInternal.WithMessage("Failed to trash merchant detail")
	ErrFailedRestoreMerchantDetail            = errors.ErrInternal.WithMessage("Failed to restore merchant detail")
	ErrFailedDeleteMerchantDetailPermanent    = errors.ErrInternal.WithMessage("Failed to permanently delete merchant detail")
	ErrFailedRestoreAllMerchantDetail         = errors.ErrInternal.WithMessage("Failed to restore all merchant details")
	ErrFailedDeleteAllMerchantDetailPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all merchant details")
)
