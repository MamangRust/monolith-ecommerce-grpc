package merchant_social_link_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedCreateMerchantSocialLink              = errors.ErrInternal.WithMessage("failed to create merchant social link")
	ErrFailedUpdateMerchantSocialLink              = errors.ErrInternal.WithMessage("failed to update merchant social link")
	ErrFailedTrashMerchantSocialLink               = errors.ErrInternal.WithMessage("failed to trash merchant social link")
	ErrFailedRestoreMerchantSocialLink             = errors.ErrInternal.WithMessage("failed to restore merchant social link")
	ErrFailedDeletePermanentMerchantSocialLink     = errors.ErrInternal.WithMessage("failed to permanently delete merchant social link")
	ErrFailedRestoreAllMerchantSocialLinks         = errors.ErrInternal.WithMessage("failed to restore all merchant social links")
	ErrFailedDeleteAllPermanentMerchantSocialLinks = errors.ErrInternal.WithMessage("failed to permanently delete all merchant social links")
	ErrFailedFindMonthlyTotalPriceByMerchant       = errors.ErrInternal.WithMessage("Failed to find monthly total price by merchant")
	ErrFailedFindYearlyTotalPriceByMerchant        = errors.ErrInternal.WithMessage("Failed to find yearly total price by merchant")
)
