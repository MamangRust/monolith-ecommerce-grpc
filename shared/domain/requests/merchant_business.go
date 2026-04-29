package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantBusinessInformationRequest struct {
	MerchantID        int    `json:"merchant_id" validate:"required"`
	BusinessType      string `json:"business_type" validate:"required"`
	TaxID             string `json:"tax_id" validate:"required"`
	EstablishedYear   int    `json:"established_year" validate:"required,min=1900,max=2100"`
	NumberOfEmployees int    `json:"number_of_employees" validate:"required,min=1"`
	WebsiteUrl        string `json:"website_url" validate:"omitempty,url"`
}

type UpdateMerchantBusinessInformationRequest struct {
	MerchantBusinessInfoID *int   `json:"merchant_business_info_id" validate:"required"`
	BusinessType           string `json:"business_type" validate:"required"`
	TaxID                  string `json:"tax_id" validate:"required"`
	EstablishedYear        int    `json:"established_year" validate:"required,min=1900,max=2100"`
	NumberOfEmployees      int    `json:"number_of_employees" validate:"required,min=1"`
	WebsiteUrl             string `json:"website_url" validate:"omitempty,url"`
}

func (r *CreateMerchantBusinessInformationRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantBusinessInformationRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
