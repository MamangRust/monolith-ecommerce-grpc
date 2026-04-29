package merchant_policy_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedFindAllMerchantPolicies    = errors.ErrInternal.WithMessage("failed to fetch all merchant policies")
	ErrFailedFindActiveMerchantPolicies = errors.ErrInternal.WithMessage("failed to fetch active merchant policies")
	ErrFailedFindTrashedMerchantPolicies = errors.ErrInternal.WithMessage("failed to fetch trashed merchant policies")
	ErrFailedFindMerchantPolicyByID     = errors.ErrInternal.WithMessage("failed to find merchant policy by id")

	ErrFailedCreateMerchantPolicy = errors.ErrInternal.WithMessage("failed to create merchant policy")
	ErrFailedUpdateMerchantPolicy = errors.ErrInternal.WithMessage("failed to update merchant policy")

	ErrFailedTrashedReviewPolicy          = errors.ErrInternal.WithMessage("failed to trash merchant policy")
	ErrFailedRestoreReviewPolicy          = errors.ErrInternal.WithMessage("failed to restore merchant policy")
	ErrFailedDeleteReviewPolicyPermanent = errors.ErrInternal.WithMessage("failed to permanently delete merchant policy")

	ErrFailedRestoreAllReviewPolicies     = errors.ErrInternal.WithMessage("failed to restore all merchant policies")
	ErrFailedDeleteAllReviewPoliciesPermanent = errors.ErrInternal.WithMessage("failed to permanently delete all merchant policies")
)
