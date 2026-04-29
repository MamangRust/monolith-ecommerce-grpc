package transaction_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedPaymentStatusCannotBeModified = errors.ErrBadRequest.WithMessage("Cannot modify payment status")
	ErrFailedPaymentStatusInvalid          = errors.ErrBadRequest.WithMessage("Invalid payment status")
	ErrFailedPaymentInsufficientBalance    = errors.ErrBadRequest.WithMessage("Insufficient balance")
	ErrFailedOrderItemEmpty                = errors.ErrInternal.WithMessage("Failed to order item empty")

	ErrFailedFindMonthlyAmountSuccess = errors.ErrInternal.WithMessage("Failed to find monthly amount success")
	ErrFailedFindYearlyAmountSuccess  = errors.ErrInternal.WithMessage("Failed to find yearly amount success")
	ErrFailedFindMonthlyAmountFailed  = errors.ErrInternal.WithMessage("Failed to find monthly amount failed")
	ErrFailedFindYearlyAmountFailed   = errors.ErrInternal.WithMessage("Failed to find yearly amount failed")

	ErrFailedFindMonthlyAmountSuccessByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly amount success by merchant")
	ErrFailedFindYearlyAmountSuccessByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly amount success by merchant")
	ErrFailedFindMonthlyAmountFailedByMerchant  = errors.ErrInternal.WithMessage("Failed to find monthly amount failed by merchant")
	ErrFailedFindYearlyAmountFailedByMerchant   = errors.ErrInternal.WithMessage("Failed to find yearly amount failed by merchant")

	ErrFailedFindMonthlyMethod           = errors.ErrInternal.WithMessage("Failed to find monthly method")
	ErrFailedFindYearlyMethod            = errors.ErrInternal.WithMessage("Failed to find yearly method")
	ErrFailedFindMonthlyMethodByMerchant = errors.ErrInternal.WithMessage("Failed to find monthly method by merchant")
	ErrFailedFindYearlyMethodByMerchant  = errors.ErrInternal.WithMessage("Failed to find yearly method by merchant")

	ErrFailedFindAllTransactions        = errors.ErrInternal.WithMessage("Failed to find all transactions")
	ErrFailedFindTransactionsByMerchant = errors.ErrInternal.WithMessage("Failed to find transactions by merchant")
	ErrFailedFindTransactionsByActive   = errors.ErrInternal.WithMessage("Failed to find active transactions")
	ErrFailedFindTransactionsByTrashed  = errors.ErrInternal.WithMessage("Failed to find trashed transactions")
	ErrFailedFindTransactionById        = errors.ErrInternal.WithMessage("Failed to find transaction by ID")
	ErrFailedFindTransactionByOrderId   = errors.ErrInternal.WithMessage("Failed to find transaction by order ID")

	ErrFailedCreateTransaction             = errors.ErrInternal.WithMessage("Failed to create transaction")
	ErrFailedUpdateTransaction             = errors.ErrInternal.WithMessage("Failed to update transaction")
	ErrFailedTrashedTransaction            = errors.ErrInternal.WithMessage("Failed to trash transaction")
	ErrFailedRestoreTransaction            = errors.ErrInternal.WithMessage("Failed to restore transaction")
	ErrFailedDeleteTransactionPermanently  = errors.ErrInternal.WithMessage("Failed to permanently delete transaction")
	ErrFailedRestoreAllTransactions        = errors.ErrInternal.WithMessage("Failed to restore all transactions")
	ErrFailedDeleteAllTransactionPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all transactions")
)
