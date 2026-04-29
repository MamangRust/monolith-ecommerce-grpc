package shippingaddress_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedCreateShippingAddress = errors.ErrInternal.WithMessage("Failed to create shipping address")
	ErrFailedUpdateShippingAddress = errors.ErrInternal.WithMessage("Failed to update shipping address")

	ErrFailedFindAllShippingAddresses            = errors.ErrInternal.WithMessage("Failed to fetch all shipping addresses")
	ErrFailedFindActiveShippingAddresses         = errors.ErrInternal.WithMessage("Failed to fetch active shipping addresses")
	ErrFailedFindTrashedShippingAddresses        = errors.ErrInternal.WithMessage("Failed to fetch trashed shipping addresses")
	ErrFailedFindShippingAddressByID             = errors.ErrInternal.WithMessage("Failed to find shipping address by ID")
	ErrFailedFindShippingAddressByOrder          = errors.ErrInternal.WithMessage("Failed to find shipping address by order ID")
	ErrFailedTrashShippingAddress                = errors.ErrInternal.WithMessage("Failed to trash shipping address")
	ErrFailedRestoreShippingAddress              = errors.ErrInternal.WithMessage("Failed to restore shipping address")
	ErrFailedDeleteShippingAddressPermanent      = errors.ErrInternal.WithMessage("Failed to permanently delete shipping address")
	ErrFailedRestoreAllShippingAddresses         = errors.ErrInternal.WithMessage("Failed to restore all shipping addresses")
	ErrFailedDeleteAllShippingAddressesPermanent = errors.ErrInternal.WithMessage("Failed to permanently delete all shipping addresses")
)
