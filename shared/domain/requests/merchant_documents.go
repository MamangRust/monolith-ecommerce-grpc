package requests

import "github.com/go-playground/validator/v10"

type FindAllMerchantDocuments struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateMerchantDocumentRequest struct {
	MerchantID   int    `json:"merchant_id" validate:"required,min=1"`
	DocumentType string `json:"document_type" validate:"required"`
	DocumentUrl  string `json:"document_url" validate:"required"`
}

type UpdateMerchantDocumentRequest struct {
	DocumentID   *int   `json:"document_id" `
	MerchantID   int    `json:"merchant_id" validate:"required,min=1"`
	DocumentType string `json:"document_type" validate:"required"`
	DocumentUrl  string `json:"document_url" validate:"required"`
	Status       string `json:"status" validate:"required"`
	Note         string `json:"note" validate:"required"`
}

type UpdateMerchantDocumentStatusRequest struct {
	DocumentID *int   `json:"document_id" `
	MerchantID int    `json:"merchant_id" validate:"required,min=1"`
	Status     string `json:"status" validate:"required"`
	Note       string `json:"note" validate:"required"`
}

func (r *CreateMerchantDocumentRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantDocumentRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantDocumentStatusRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}
	return nil
}
