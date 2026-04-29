package merchantaward_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedFindAllMerchantAwards            = errors.ErrInternal.WithMessage("Failed to find all merchant awards")
	ErrFailedFindActiveMerchantAwards         = errors.ErrInternal.WithMessage("Failed to find active merchant awards")
	ErrFailedFindTrashedMerchantAwards        = errors.ErrInternal.WithMessage("Failed to find trashed merchant awards")
	ErrFailedFindMerchantAwardById            = errors.ErrInternal.WithMessage("Failed to find merchant award by ID")
	ErrFailedCreateMerchantAward              = errors.ErrInternal.WithMessage("Failed to create merchant award")
	ErrFailedUpdateMerchantAward              = errors.ErrInternal.WithMessage("Failed to update merchant award")
	ErrFailedTrashedMerchantAward             = errors.ErrInternal.WithMessage("Failed to trash merchant award")
	ErrFailedRestoreMerchantAward             = errors.ErrInternal.WithMessage("Failed to restore merchant award")
	ErrFailedDeleteMerchantAwardPermanent     = errors.ErrInternal.WithMessage("Failed to permanently delete merchant award")
	ErrFailedRestoreAllMerchantAwards         = errors.ErrInternal.WithMessage("Failed to restore all merchant awards")
	ErrFailedDeleteAllMerchantAwardsPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all merchant awards")
)
