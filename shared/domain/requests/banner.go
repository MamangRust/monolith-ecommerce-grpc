package requests

import "github.com/go-playground/validator/v10"

type BannerFormData struct {
	Name      string
	StartDate string
	EndDate   string
	StartTime string
	EndTime   string
	IsActive  bool
	ImagePath string
}

type FindAllBanner struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateBannerRequest struct {
	Name      string `json:"name" validate:"required"`
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02"`
	StartTime string `json:"start_time" validate:"required,datetime=15:04:05"`
	EndTime   string `json:"end_time" validate:"required,datetime=15:04:05"`
	IsActive  bool   `json:"is_active"`
}

type UpdateBannerRequest struct {
	BannerID  *int   `json:"banner_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02"`
	StartTime string `json:"start_time" validate:"required,datetime=15:04:05"`
	EndTime   string `json:"end_time" validate:"required,datetime=15:04:05"`
	IsActive  bool   `json:"is_active"`
}

func (r *CreateBannerRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateBannerRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
