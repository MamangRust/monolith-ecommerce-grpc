package slider_errors

import (
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
)

var (
	ErrFindAllSliders           = errors.ErrNotFound.WithMessage("failed to find all sliders")
	ErrFindActiveSliders        = errors.ErrNotFound.WithMessage("failed to find active sliders")
	ErrFindTrashedSliders       = errors.ErrNotFound.WithMessage("failed to find trashed sliders")
	ErrFindSliderByID           = errors.ErrNotFound.WithMessage("failed to find slider by ID")
	ErrSliderNotFound           = errors.ErrNotFound.WithMessage("slider not found")
	ErrCreateSlider             = errors.ErrInternal.WithMessage("failed to create slider")
	ErrUpdateSlider             = errors.ErrInternal.WithMessage("failed to update slider")
	ErrTrashSlider              = errors.ErrInternal.WithMessage("failed to trash slider")
	ErrRestoreSlider            = errors.ErrInternal.WithMessage("failed to restore slider")
	ErrDeletePermanentSlider    = errors.ErrInternal.WithMessage("failed to permanently delete slider")
	ErrRestoreAllSlider         = errors.ErrInternal.WithMessage("failed to restore all sliders")
	ErrDeleteAllPermanentSlider = errors.ErrInternal.WithMessage("failed to permanently delete all sliders")
)
