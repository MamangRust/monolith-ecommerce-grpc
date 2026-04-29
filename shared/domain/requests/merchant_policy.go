package requests

import "github.com/go-playground/validator/v10"

type CreateMerchantPolicyRequest struct {
	MerchantID  int    `json:"merchant_id" validate:"required"`
	PolicyType  string `json:"policy_type" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateMerchantPolicyRequest struct {
	MerchantPolicyID *int   `json:"merchant_policy_id" validate:"required"`
	PolicyType       string `json:"policy_type" validate:"required"`
	Title            string `json:"title" validate:"required"`
	Description      string `json:"description" validate:"required"`
}

func (r *CreateMerchantPolicyRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateMerchantPolicyRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
