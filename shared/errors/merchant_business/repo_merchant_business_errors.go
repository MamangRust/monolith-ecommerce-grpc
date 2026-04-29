package merchantbusiness_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrMerchantBusinessNotFound             = errors.ErrNotFound.WithMessage("merchant business not found")
	ErrFindAllMerchantBusinesses            = errors.ErrInternal.WithMessage("failed to find all merchant businesses")
	ErrFindActiveMerchantBusinesses         = errors.ErrInternal.WithMessage("failed to find active merchant businesses")
	ErrFindTrashedMerchantBusinesses        = errors.ErrInternal.WithMessage("failed to find trashed merchant businesses")
	ErrCreateMerchantBusiness               = errors.ErrInternal.WithMessage("failed to create merchant business")
	ErrUpdateMerchantBusiness               = errors.ErrInternal.WithMessage("failed to update merchant business")
	ErrTrashMerchantBusiness                = errors.ErrInternal.WithMessage("failed to move merchant business to trash")
	ErrRestoreMerchantBusiness              = errors.ErrInternal.WithMessage("failed to restore merchant business from trash")
	ErrDeletePermanentMerchantBusiness      = errors.ErrInternal.WithMessage("failed to permanently delete merchant business")
	ErrRestoreAllMerchantBusinesses         = errors.ErrInternal.WithMessage("failed to restore all merchant businesses")
	ErrDeleteAllPermanentMerchantBusinesses = errors.ErrInternal.WithMessage("failed to permanently delete all merchant businesses")

	ErrMerchantBusinessInternal = errors.ErrInternal.WithMessage("merchant business internal repository error")
)
