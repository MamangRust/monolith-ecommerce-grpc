package merchantbusiness_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedFindAllMerchantBusiness            = errors.ErrInternal.WithMessage("Failed to fetch all merchant businesses")
	ErrFailedFindActiveMerchantBusiness         = errors.ErrInternal.WithMessage("Failed to fetch active merchant businesses")
	ErrFailedFindTrashedMerchantBusiness        = errors.ErrInternal.WithMessage("Failed to fetch trashed merchant businesses")
	ErrFailedFindMerchantBusinessById           = errors.ErrInternal.WithMessage("Failed to find merchant business by ID")
	ErrFailedCreateMerchantBusiness             = errors.ErrInternal.WithMessage("Failed to create merchant business")
	ErrFailedUpdateMerchantBusiness             = errors.ErrInternal.WithMessage("Failed to update merchant business")
	ErrFailedTrashedMerchantBusiness            = errors.ErrInternal.WithMessage("Failed to trash merchant business")
	ErrFailedRestoreMerchantBusiness            = errors.ErrInternal.WithMessage("Failed to restore merchant business")
	ErrFailedDeleteMerchantBusinessPermanent    = errors.ErrInternal.WithMessage("Failed to permanently delete merchant business")
	ErrFailedRestoreAllMerchantBusiness         = errors.ErrInternal.WithMessage("Failed to restore all merchant businesses")
	ErrFailedDeleteAllMerchantBusinessPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all merchant businesses")
)
