package merchantdetail_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrMerchantDetailNotFound            = errors.ErrNotFound.WithMessage("merchant detail not found")
	ErrFindAllMerchantDetails            = errors.ErrInternal.WithMessage("failed to find all merchant details")
	ErrFindActiveMerchantDetails         = errors.ErrInternal.WithMessage("failed to find active merchant details")
	ErrFindTrashedMerchantDetails        = errors.ErrInternal.WithMessage("failed to find trashed merchant details")
	ErrFindByIdTrashedMerchantDetail     = errors.ErrInternal.WithMessage("failed to find trashed merchant detail by ID")
	ErrCreateMerchantDetail              = errors.ErrInternal.WithMessage("failed to create merchant detail")
	ErrUpdateMerchantDetail              = errors.ErrInternal.WithMessage("failed to update merchant detail")
	ErrTrashMerchantDetail               = errors.ErrInternal.WithMessage("failed to move merchant detail to trash")
	ErrRestoreMerchantDetail             = errors.ErrInternal.WithMessage("failed to restore merchant detail from trash")
	ErrDeletePermanentMerchantDetail     = errors.ErrInternal.WithMessage("failed to permanently delete merchant detail")
	ErrRestoreAllMerchantDetails         = errors.ErrInternal.WithMessage("failed to restore all merchant details")
	ErrDeleteAllPermanentMerchantDetails = errors.ErrInternal.WithMessage("failed to permanently delete all merchant details")

	ErrMerchantDetailInternal = errors.ErrInternal.WithMessage("merchant detail internal repository error")
)
