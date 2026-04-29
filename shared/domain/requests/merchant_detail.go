package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantDetailFormData struct {
	MerchantID       int
	DisplayName      string
	CoverImageUrl    string
	LogoUrl          string
	ShortDescription string
	WebsiteUrl       string
}

type UpdateMerchantDetailFormData struct {
	DisplayName      string
	CoverImageUrl    string
	LogoUrl          string
	ShortDescription string
	WebsiteUrl       string
}

type CreateMerchantDetailRequest struct {
	MerchantID       int    `json:"merchant_id" validate:"required"`
	DisplayName      string `json:"display_name" validate:"required"`
	CoverImageUrl    string `json:"cover_image_url"`
	LogoUrl          string `json:"logo_url"`
	ShortDescription string `json:"short_description" validate:"required"`
	WebsiteUrl       string `json:"website_url" validate:"omitempty,url"`
}

type UpdateMerchantDetailRequest struct {
	MerchantDetailID *int   `json:"merchant_detail_id" validate:"required"`
	DisplayName      string `json:"display_name" validate:"required"`
	CoverImageUrl    string `json:"cover_image_url" `
	LogoUrl          string `json:"logo_url"`
	ShortDescription string `json:"short_description" validate:"required"`
	WebsiteUrl       string `json:"website_url" validate:"omitempty,url"`
}

func (r *CreateMerchantDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
