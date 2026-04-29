package requests

import "github.com/go-playground/validator/v10"

type ReviewDetailFormData struct {
	ReviewID int
	Type     string
	Url      string
	Caption  string
}

type CreateReviewDetailRequest struct {
	ReviewID int    `json:"review_id" validate:"required"`
	Type     string `json:"type" validate:"required"`
	Url      string `json:"url" validate:"required"`
	Caption  string `json:"caption" validate:"required"`
}

type UpdateReviewDetailRequest struct {
	ReviewDetailID *int   `json:"review_detail_id" validate:"required"`
	Type           string `json:"type" validate:"required"`
	Url            string `json:"url" validate:"required"`
	Caption        string `json:"caption" validate:"required"`
}

func (r *CreateReviewDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateReviewDetailRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
