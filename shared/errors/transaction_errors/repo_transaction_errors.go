package transaction_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrGetMonthlyAmountSuccess = errors.ErrInternal.WithMessage("failed to get monthly amount success")
	ErrGetYearlyAmountSuccess  = errors.ErrInternal.WithMessage("failed to get yearly amount success")
	ErrGetMonthlyAmountFailed  = errors.ErrInternal.WithMessage("failed to get monthly amount failed")
	ErrGetYearlyAmountFailed   = errors.ErrInternal.WithMessage("failed to get yearly amount failed")

	ErrGetMonthlyAmountSuccessByMerchant = errors.ErrInternal.WithMessage("failed to get monthly amount success by merchant")
	ErrGetYearlyAmountSuccessByMerchant  = errors.ErrInternal.WithMessage("failed to get yearly amount success by merchant")
	ErrGetMonthlyAmountFailedByMerchant  = errors.ErrInternal.WithMessage("failed to get monthly amount failed by merchant")
	ErrGetYearlyAmountFailedByMerchant   = errors.ErrInternal.WithMessage("failed to get yearly amount failed by merchant")

	ErrGetMonthlyTransactionMethod           = errors.ErrInternal.WithMessage("failed to get monthly transaction method")
	ErrGetYearlyTransactionMethod            = errors.ErrInternal.WithMessage("failed to get yearly transaction method")
	ErrGetMonthlyTransactionMethodByMerchant = errors.ErrInternal.WithMessage("failed to get monthly transaction method by merchant")
	ErrGetYearlyTransactionMethodByMerchant  = errors.ErrInternal.WithMessage("failed to get yearly transaction method by merchant")

	ErrFindAllTransactions = errors.ErrInternal.WithMessage("failed to find all transactions")
	ErrFindByActive        = errors.ErrInternal.WithMessage("failed to find active transactions")
	ErrFindByTrashed       = errors.ErrInternal.WithMessage("failed to find trashed transactions")
	ErrFindByMerchant      = errors.ErrInternal.WithMessage("failed to find transactions by merchant")
	ErrFindById            = errors.ErrNotFound.WithMessage("failed to find transaction by ID")
	ErrFindByOrderId       = errors.ErrNotFound.WithMessage("failed to find transaction by order ID")

	ErrCreateTransaction             = errors.ErrInternal.WithMessage("failed to create transaction")
	ErrUpdateTransaction             = errors.ErrInternal.WithMessage("failed to update transaction")
	ErrTrashTransaction              = errors.ErrInternal.WithMessage("failed to move transaction to trash")
	ErrRestoreTransaction            = errors.ErrInternal.WithMessage("failed to restore transaction")
	ErrDeleteTransactionPermanently  = errors.ErrInternal.WithMessage("failed to permanently delete transaction")
	ErrRestoreAllTransactions        = errors.ErrInternal.WithMessage("failed to restore all transactions")
	ErrDeleteAllTransactionPermanent = errors.ErrInternal.WithMessage("failed to permanently delete all transactions")

	ErrTransactionNotFound = errors.ErrNotFound.WithMessage("transaction not found")
)

