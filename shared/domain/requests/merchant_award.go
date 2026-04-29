package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantCertificationOrAwardRequest struct {
	MerchantID     int    `json:"merchant_id" validate:"required"`
	Title          string `json:"title" validate:"required"`
	Description    string `json:"description" validate:"required"`
	IssuedBy       string `json:"issued_by" validate:"required"`
	IssueDate      string `json:"issue_date" validate:"required,datetime=2006-01-02"`
	ExpiryDate     string `json:"expiry_date" validate:"omitempty,datetime=2006-01-02"`
	CertificateUrl string `json:"certificate_url" validate:"omitempty,url"`
}

type UpdateMerchantCertificationOrAwardRequest struct {
	MerchantCertificationID *int   `json:"merchant_certification_id" validate:"required"`
	Title                   string `json:"title" validate:"required"`
	Description             string `json:"description" validate:"required"`
	IssuedBy                string `json:"issued_by" validate:"required"`
	IssueDate               string `json:"issue_date" validate:"required,datetime=2006-01-02"`
	ExpiryDate              string `json:"expiry_date" validate:"omitempty,datetime=2006-01-02"`
	CertificateUrl          string `json:"certificate_url" validate:"omitempty,url"`
}

func (r *CreateMerchantCertificationOrAwardRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantCertificationOrAwardRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
