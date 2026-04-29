package requests

import "github.com/go-playground/validator/v10"

type MonthAmountTransaction struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type MonthAmountTransactionMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
	Month      int `json:"month" validate:"required"`
}

type YearAmountTransactionMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type MonthMethodTransaction struct {
	Year  int `json:"year" validate:"required"`
	Month int `json:"month" validate:"required"`
}

type MonthMethodTransactionMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
	Month      int `json:"month" validate:"required"`
}

type YearMethodTransactionMerchant struct {
	MerchantID int `json:"merchant_id" validate:"required"`
	Year       int `json:"year" validate:"required"`
}

type FindAllTransaction struct {
	Search   string `json:"search" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
}

type FindAllTransactionByMerchant struct {
	MerchantID int    `json:"merchant_id" validate:"required"`
	Search     string `json:"search" validate:"required"`
	Page       int    `json:"page" validate:"min=1"`
	PageSize   int    `json:"page_size" validate:"min=1,max=100"`
}

type CreateTransactionRequest struct {
	UserID        int     `json:"user_id"`
	OrderID       int     `json:"order_id" validate:"required"`
	MerchantID    int     `json:"merchant_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Amount        int     `json:"amount" validate:"required"`
	PaymentStatus *string `json:"payment_status"`
}

type UpdateTransactionRequest struct {
	TransactionID *int    `json:"transaction_id"`
	OrderID       int     `json:"order_id" validate:"required"`
	MerchantID    int     `json:"merchant_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Amount        int     `json:"amount" validate:"required"`
	PaymentStatus *string `json:"payment_status"`
}

func (r *CreateTransactionRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *UpdateTransactionRequest) Validate() error {
	validate := validator.New()
	err := validate.Struct(r)
	if err != nil {
		return err
	}
	return nil
}
