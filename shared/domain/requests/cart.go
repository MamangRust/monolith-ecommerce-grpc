package requests

import "github.com/go-playground/validator/v10"

type FindAllCarts struct {
	UserID   int    `json:"user_id"`
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type CartCreateRecord struct {
	ProductID    int    `json:"product_id"`
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	ImageProduct string `json:"image_product"`
	Quantity     int    `json:"quantity"`
	Weight       int    `json:"weight"`
}

type CreateCartRequest struct {
	Quantity  int `json:"quantity" validate:"required,gt=0"`
	ProductID int `json:"product_id" validate:"required"`
	UserID    int `json:"user_id,omitempty"`
}

type DeleteAllCartRequest struct {
	UserID  int   `json:"user_id" validate:"required"`
	CartIds []int `json:"cart_ids" validate:"required"`
}

type DeleteCartRequest struct {
	CartID int `json:"cart_id" validate:"required"`
	UserID int `json:"user_id" validate:"required"`
}

func (l *CartCreateRecord) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}
func (l *CreateCartRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}

func (l *DeleteCartRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}

func (l *DeleteAllCartRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(l)

	if err != nil {
		return err
	}

	return nil
}
