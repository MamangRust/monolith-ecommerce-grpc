package merchant_social_link_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrMerchantSocialLinkNotFound            = errors.ErrNotFound.WithMessage("merchant social link not found")
	ErrFindAllMerchantSocialLinks            = errors.ErrInternal.WithMessage("failed to find all merchant social links")
	ErrFindActiveMerchantSocialLinks         = errors.ErrInternal.WithMessage("failed to find active merchant social links")
	ErrFindTrashedMerchantSocialLinks        = errors.ErrInternal.WithMessage("failed to find trashed merchant social links")
	ErrCreateMerchantSocialLink              = errors.ErrInternal.WithMessage("failed to create merchant social link")
	ErrUpdateMerchantSocialLink              = errors.ErrInternal.WithMessage("failed to update merchant social link")
	ErrTrashMerchantSocialLink               = errors.ErrInternal.WithMessage("failed to move merchant social link to trash")
	ErrRestoreMerchantSocialLink             = errors.ErrInternal.WithMessage("failed to restore merchant social link from trash")
	ErrDeletePermanentMerchantSocialLink     = errors.ErrInternal.WithMessage("failed to permanently delete merchant social link")
	ErrRestoreAllMerchantSocialLinks         = errors.ErrInternal.WithMessage("failed to restore all merchant social links")
	ErrDeleteAllPermanentMerchantSocialLinks = errors.ErrInternal.WithMessage("failed to permanently delete all merchant social links")

	ErrMerchantSocialLinkInternal = errors.ErrInternal.WithMessage("merchant social link internal repository error")
)
