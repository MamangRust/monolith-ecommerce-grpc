package merchant_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrMerchantNotFound     = errors.ErrNotFound.WithMessage("merchant not found")
	ErrFindAllMerchants     = errors.ErrInternal.WithMessage("failed to find all merchants")
	ErrFindActiveMerchants  = errors.ErrInternal.WithMessage("failed to find active merchants")
	ErrFindTrashedMerchants = errors.ErrInternal.WithMessage("failed to find trashed merchants")
	ErrMerchantConflict     = errors.ErrConflict.WithMessage("merchant already exists")

	ErrCreateMerchant = errors.ErrInternal.WithMessage("failed to create merchant")
	ErrUpdateMerchant = errors.ErrInternal.WithMessage("failed to update merchant")

	ErrTrashedMerchant         = errors.ErrInternal.WithMessage("failed to move merchant to trash")
	ErrRestoreMerchant         = errors.ErrInternal.WithMessage("failed to restore merchant from trash")
	ErrDeleteMerchantPermanent = errors.ErrInternal.WithMessage("failed to permanently delete merchant")

	ErrRestoreAllMerchants = errors.ErrInternal.WithMessage("failed to restore all merchants")
	ErrDeleteAllMerchants  = errors.ErrInternal.WithMessage("failed to permanently delete all merchants")

	ErrMerchantInternal = errors.ErrInternal.WithMessage("merchant internal repository error")
)


