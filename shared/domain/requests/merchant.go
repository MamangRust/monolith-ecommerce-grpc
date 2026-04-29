package requests

import "github.com/go-playground/validator/v10"

type FindAllMerchant struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateMerchantRequest struct {
	UserID       int    `json:"user_id"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	Address      string `json:"address" validate:"required"`
	ContactEmail string `json:"contact_email" validate:"required,email"`
	ContactPhone string `json:"contact_phone" validate:"required"`
	Status       string `json:"status" validate:"required"`
}

type UpdateMerchantRequest struct {
	MerchantID   *int   `json:"merchant_id"`
	UserID       int    `json:"user_id"`
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description" validate:"required"`
	Address      string `json:"address" validate:"required"`
	ContactEmail string `json:"contact_email" validate:"required,email"`
	ContactPhone string `json:"contact_phone" validate:"required"`
	Status       string `json:"status" validate:"required"`
}

type UpdateMerchantStatusRequest struct {
	MerchantID *int   `json:"merchant_id"`
	Status     string `json:"status" validate:"required"`
}

func (r *CreateMerchantRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantStatusRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
