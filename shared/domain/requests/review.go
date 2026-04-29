package requests

import "github.com/go-playground/validator/v10"

type FindAllReview struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllReviewByMerchant struct {
	MerchantID int    `json:"merchant_id"`
	Rating     int    `json:"rating"`
	Search     string `json:"search" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllReviewByProduct struct {
	ProductID int    `json:"product_id"`
	Rating    int    `json:"rating"`
	Search    string `json:"search" validate:"required"`
	Page      int    `json:"page" validate:"min=1"`
	PageSize  int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateReviewRequest struct {
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

type UpdateReviewRequest struct {
	ReviewID *int   `json:"review_id"`
	Name     string `json:"name"`
	Rating   int    `json:"rating" validate:"required"`
	Comment  string `json:"comment" validate:"required"`
}

func (l *CreateReviewRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}

func (l *UpdateReviewRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}
