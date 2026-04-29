package requests

import "github.com/go-playground/validator/v10"

type FindAllCategory struct {
	Search   string `json:"search" form:"search" validate:"required"`
	Page     int    `json:"page" form:"page" validate:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

type MonthTotalPrice struct {
	Year  int `json:"year" form:"year" validate:"required"`
	Month int `json:"month" form:"month" validate:"required"`
}

type MonthTotalPriceCategory struct {
	CategoryID int `json:"category_id" form:"category_id"`
	Year       int `json:"year" form:"year" validate:"required"`
	Month      int `json:"month" form:"month" validate:"required"`
}

type YearTotalPriceCategory struct {
	CategoryID int `json:"category_id" form:"category_id"`
	Year       int `json:"year" form:"year" validate:"required"`
}

type MonthTotalPriceMerchant struct {
	MerchantID int `json:"merchant_id" form:"merchant_id" validate:"required"`
	Year       int `json:"year" form:"year" validate:"required"`
	Month      int `json:"month" form:"month" validate:"required"`
}

type YearTotalPriceMerchant struct {
	MerchantID int `json:"merchant_id" form:"merchant_id" validate:"required"`
	Year       int `json:"year" form:"year" validate:"required"`
}

type MonthPriceMerchant struct {
	MerchantID int `json:"merchant_id" form:"merchant_id" validate:"required"`
	Year       int `json:"year" form:"year" validate:"required"`
}

type YearPriceMerchant struct {
	MerchantID int `json:"merchant_id" form:"merchant_id" validate:"required"`
	Year       int `json:"year" form:"year" validate:"required"`
}

type MonthPriceId struct {
	CategoryID int `json:"category_id" form:"category_id"`
	Year       int `json:"year" form:"year" validate:"required"`
}

type YearPriceId struct {
	CategoryID int `json:"category_id" form:"category_id"`
	Year       int `json:"year" form:"year" validate:"required"`
}

type CategoryNameAndId struct {
	Name       string `json:"name" form:"name"`
	CategoryID int    `json:"category_id" form:"category_id"`
}

type CategoryFormData struct {
	Name          string  `json:"name" form:"name" validate:"required"`
	Description   string  `json:"description" form:"description" validate:"required"`
	SlugCategory  *string `json:"slug_category" form:"slug_category" validate:"required"`
	ImageCategory string  `json:"image_category" form:"image_category" validate:"required"`
}

type CreateCategoryRequest struct {
	Name          string  `json:"name" form:"name" validate:"required"`
	Description   string  `json:"description" form:"description" validate:"required"`
	SlugCategory  *string `json:"slug_category" form:"slug_category" validate:"required"`
	ImageCategory string  `json:"image_category" form:"image_category" validate:"required"`
}

type UpdateCategoryRequest struct {
	CategoryID    *int    `json:"category_id" form:"category_id" `
	Name          string  `json:"name" form:"name" validate:"required"`
	Description   string  `json:"description" form:"description" validate:"required"`
	SlugCategory  *string `json:"slug_category" form:"slug_category" validate:"required"`
	ImageCategory string  `json:"image_category" form:"image_category" validate:"required"`
}

func (r *CreateCategoryRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateCategoryRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
