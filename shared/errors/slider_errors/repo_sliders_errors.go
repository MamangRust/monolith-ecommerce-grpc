package slider_errors

import "errors"

var (
	ErrFindAllSliders           = errors.New("failed to find all sliders")
	ErrFindActiveSliders        = errors.New("failed to find active sliders")
	ErrFindTrashedSliders       = errors.New("failed to find trashed sliders")
	ErrFindSliderByID           = errors.New("failed to find slider by ID")
	ErrCreateSlider             = errors.New("failed to create slider")
	ErrUpdateSlider             = errors.New("failed to update slider")
	ErrTrashSlider              = errors.New("failed to trash slider")
	ErrRestoreSlider            = errors.New("failed to restore slider")
	ErrDeletePermanentSlider    = errors.New("failed to permanently delete slider")
	ErrRestoreAllSlider         = errors.New("failed to restore all sliders")
	ErrDeleteAllPermanentSlider = errors.New("failed to permanently delete all sliders")
)
