package slider_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)


var (
	ErrFailedFindAllSliders            = errors.ErrInternal.WithMessage("failed to fetch sliders")
	ErrFailedFindActiveSliders         = errors.ErrInternal.WithMessage("failed to fetch active sliders")
	ErrFailedFindTrashedSliders        = errors.ErrInternal.WithMessage("failed to fetch trashed sliders")
	ErrFailedCreateSlider              = errors.ErrInternal.WithMessage("failed to create slider")
	ErrFailedUpdateSlider              = errors.ErrInternal.WithMessage("failed to update slider")
	ErrFailedTrashSlider               = errors.ErrInternal.WithMessage("failed to trash slider")
	ErrFailedRestoreSlider             = errors.ErrInternal.WithMessage("failed to restore slider")
	ErrFailedDeletePermanentSlider     = errors.ErrInternal.WithMessage("failed to permanently delete slider")
	ErrFailedFindSliderByID            = errors.ErrInternal.WithMessage("failed to fetch slider by ID")
	ErrFailedRestoreAllSliders         = errors.ErrInternal.WithMessage("failed to restore all sliders")
	ErrFailedDeleteAllPermanentSliders = errors.ErrInternal.WithMessage("failed to permanently delete all sliders")
)
