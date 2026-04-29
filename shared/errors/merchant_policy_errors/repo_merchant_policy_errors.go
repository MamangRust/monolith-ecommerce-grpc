package merchant_policy_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrMerchantPolicyNotFound = errors.ErrNotFound.WithMessage("merchant policy not found")

	ErrFindAllMerchantPolicies    = errors.ErrInternal.WithMessage("failed to find all merchant policies")
	ErrFindActiveMerchantPolicies = errors.ErrInternal.WithMessage("failed to find active merchant policies")
	ErrFindTrashedMerchantPolicies = errors.ErrInternal.WithMessage("failed to find trashed merchant policies")
	ErrFindMerchantPolicyByID     = errors.ErrInternal.WithMessage("failed to find merchant policy by id")

	ErrCreateMerchantPolicy = errors.ErrInternal.WithMessage("failed to create merchant policy")
	ErrUpdateMerchantPolicy = errors.ErrInternal.WithMessage("failed to update merchant policy")

	ErrTrashedMerchantPolicy          = errors.ErrInternal.WithMessage("failed to move merchant policy to trash")
	ErrRestoreMerchantPolicy          = errors.ErrInternal.WithMessage("failed to restore merchant policy from trash")
	ErrDeleteMerchantPolicyPermanent = errors.ErrInternal.WithMessage("failed to permanently delete merchant policy")

	ErrRestoreAllMerchantPolicies     = errors.ErrInternal.WithMessage("failed to restore all merchant policies")
	ErrDeleteAllMerchantPoliciesPermanent = errors.ErrInternal.WithMessage("failed to permanently delete all merchant policies")
)
