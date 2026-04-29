package merchantaward_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrMerchantAwardNotFound            = errors.ErrNotFound.WithMessage("merchant award not found")
	ErrFindAllMerchantAwards            = errors.ErrInternal.WithMessage("failed to find all merchant awards")
	ErrFindByActiveMerchantAwards       = errors.ErrInternal.WithMessage("failed to find active merchant awards")
	ErrFindByTrashedMerchantAwards      = errors.ErrInternal.WithMessage("failed to find trashed merchant awards")
	ErrFindByIdMerchantAward            = errors.ErrInternal.WithMessage("failed to find merchant award by ID")
	ErrCreateMerchantAward              = errors.ErrInternal.WithMessage("failed to create merchant award")
	ErrUpdateMerchantAward              = errors.ErrInternal.WithMessage("failed to update merchant award")
	ErrTrashedMerchantAward             = errors.ErrInternal.WithMessage("failed to move merchant award to trash")
	ErrRestoreMerchantAward             = errors.ErrInternal.WithMessage("failed to restore merchant award from trash")
	ErrDeleteMerchantAwardPermanent     = errors.ErrInternal.WithMessage("failed to permanently delete merchant award")
	ErrRestoreAllMerchantAwards         = errors.ErrInternal.WithMessage("failed to restore all merchant awards")
	ErrDeleteAllMerchantAwardsPermanent = errors.ErrInternal.WithMessage("failed to permanently delete all merchant awards")

	ErrMerchantAwardInternal = errors.ErrInternal.WithMessage("merchant award internal repository error")
)
