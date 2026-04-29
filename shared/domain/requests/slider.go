package requests

import "github.com/go-playground/validator/v10"

type FindAllSlider struct {
	Search   string `json:"search" form:"search" validate:"required"`
	Page     int    `json:"page" form:"page" validate:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

type SliderFormData struct {
	Nama     string `json:"nama" form:"nama"`
	FilePath string `json:"file_path" form:"file_path"`
}

type CreateSliderRequest struct {
	Nama     string `json:"nama" form:"nama"`
	FilePath string `json:"file_path" form:"file_path"`
}

type UpdateSliderRequest struct {
	ID       *int   `json:"id" form:"id"`
	Nama     string `json:"nama" form:"nama"`
	FilePath string `json:"file_path" form:"file_path"`
}

func (l *CreateSliderRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}

func (l *UpdateSliderRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}
