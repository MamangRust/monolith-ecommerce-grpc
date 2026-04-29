package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantSocialRequest struct {
	MerchantDetailID *int   `json:"merchant_detail_id" validate:"required"`
	Platform         string `json:"platform" validate:"required"`
	Url              string `json:"url" validate:"required,url"`
}

type UpdateMerchantSocialRequest struct {
	ID               int    `json:"id"`
	MerchantDetailID *int   `json:"merchant_detail_id" validate:"required"`
	Platform         string `json:"platform" validate:"required"`
	Url              string `json:"url" validate:"required,url"`
}

func (r *CreateMerchantSocialRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantSocialRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
