package merchant_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedFindAllMerchants            = errors.ErrInternal.WithMessage("Failed to find all merchants")
	ErrFailedFindActiveMerchants         = errors.ErrInternal.WithMessage("Failed to find active merchants")
	ErrFailedFindTrashedMerchants        = errors.ErrInternal.WithMessage("Failed to find trashed merchants")
	ErrFailedFindMerchantById            = errors.ErrInternal.WithMessage("Failed to find merchant by ID")
	ErrFailedCreateMerchant              = errors.ErrInternal.WithMessage("Failed to create merchant")
	ErrFailedUpdateMerchant              = errors.ErrInternal.WithMessage("Failed to update merchant")
	ErrFailedTrashedMerchant             = errors.ErrInternal.WithMessage("Failed to trash merchant")
	ErrFailedRestoreMerchant             = errors.ErrInternal.WithMessage("Failed to restore merchant")
	ErrFailedDeleteMerchantPermanent     = errors.ErrInternal.WithMessage("Failed to permanently delete merchant")
	ErrFailedRestoreAllMerchants         = errors.ErrInternal.WithMessage("Failed to restore all merchants")
	ErrFailedDeleteAllMerchantsPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all merchants")

	ErrFailedFindAllMerchantDocuments            = errors.ErrInternal.WithMessage("Failed to find all merchant documents")
	ErrFailedFindActiveMerchantDocuments         = errors.ErrInternal.WithMessage("Failed to find active merchant documents")
	ErrFailedFindTrashedMerchantDocuments        = errors.ErrInternal.WithMessage("Failed to find trashed merchant documents")
	ErrFailedFindMerchantDocumentById            = errors.ErrInternal.WithMessage("Failed to find merchant document by ID")
	ErrFailedCreateMerchantDocument              = errors.ErrInternal.WithMessage("Failed to create merchant document")
	ErrFailedUpdateMerchantDocument              = errors.ErrInternal.WithMessage("Failed to update merchant document")
	ErrFailedUpdateMerchantDocumentStatus        = errors.ErrInternal.WithMessage("Failed to update merchant document status")
	ErrFailedTrashedMerchantDocument             = errors.ErrInternal.WithMessage("Failed to trash merchant document")
	ErrFailedRestoreMerchantDocument             = errors.ErrInternal.WithMessage("Failed to restore merchant document")
	ErrFailedDeleteMerchantDocumentPermanent     = errors.ErrInternal.WithMessage("Failed to permanently delete merchant document")
	ErrFailedRestoreAllMerchantDocuments         = errors.ErrInternal.WithMessage("Failed to restore all merchant documents")
	ErrFailedDeleteAllMerchantDocumentsPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all merchant documents")

	ErrMerchantDocumentNotFoundRes = errors.ErrNotFound.WithMessage("Merchant document not found")
)
