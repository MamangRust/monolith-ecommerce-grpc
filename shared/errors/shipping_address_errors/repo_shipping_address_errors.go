package shippingaddress_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrFindAllShippingAddress            = errors.ErrInternal.WithMessage("failed to find all shipping addresses")
	ErrFindActiveShippingAddress         = errors.ErrInternal.WithMessage("failed to find active shipping addresses")
	ErrFindTrashedShippingAddress        = errors.ErrInternal.WithMessage("failed to find trashed shipping addresses")
	ErrFindShippingAddressByID           = errors.ErrNotFound.WithMessage("failed to find shipping address by ID")
	ErrFindShippingAddressByOrder        = errors.ErrNotFound.WithMessage("failed to find shipping address by order ID")
	ErrCreateShippingAddress             = errors.ErrInternal.WithMessage("failed to create shipping address")
	ErrUpdateShippingAddress             = errors.ErrInternal.WithMessage("failed to update shipping address")
	ErrTrashShippingAddress              = errors.ErrInternal.WithMessage("failed to trash shipping address")
	ErrRestoreShippingAddress            = errors.ErrInternal.WithMessage("failed to restore shipping address")
	ErrDeleteShippingAddressPermanent    = errors.ErrInternal.WithMessage("failed to permanently delete shipping address")
	ErrRestoreAllShippingAddresses       = errors.ErrInternal.WithMessage("failed to restore all shipping addresses")
	ErrDeleteAllPermanentShippingAddress = errors.ErrInternal.WithMessage("failed to permanently delete all shipping addresses")

	ErrShippingAddressNotFound = errors.ErrNotFound.WithMessage("shipping address not found")
)

