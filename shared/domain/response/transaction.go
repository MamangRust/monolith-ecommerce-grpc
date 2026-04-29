package response

type TransactionResponse struct {
	ID            int    `json:"id"`
	OrderID       int    `json:"order_id"`
	MerchantID    int    `json:"merchant_id"`
	PaymentMethod string `json:"payment_method"`
	Amount        int    `json:"amount"`
	PaymentStatus string `json:"payment_status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type TransactionResponseDeleteAt struct {
	ID            int     `json:"id"`
	OrderID       int     `json:"order_id"`
	MerchantID    int     `json:"merchant_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        int     `json:"amount"`
	PaymentStatus string  `json:"payment_status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at"`
}

type TransactionMonthlyAmountSuccessResponse struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TransactionMonthlyAmountFailedResponse struct {
	Year        string `json:"year"`
	Month       string `json:"month"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TransactionYearlyAmountSuccessResponse struct {
	Year         string `json:"year"`
	TotalSuccess int    `json:"total_success"`
	TotalAmount  int    `json:"total_amount"`
}

type TransactionYearlyAmountFailedResponse struct {
	Year        string `json:"year"`
	TotalFailed int    `json:"total_failed"`
	TotalAmount int    `json:"total_amount"`
}

type TransactionMonthlyMethodResponse struct {
	Month             string `json:"month"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int    `json:"total_transactions"`
	TotalAmount       int    `json:"total_amount"`
}

type TransactionYearlyMethodResponse struct {
	Year              string `json:"year"`
	PaymentMethod     string `json:"payment_method"`
	TotalTransactions int    `json:"total_transactions"`
	TotalAmount       int    `json:"total_amount"`
}

type ApiResponseTransaction struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    *TransactionResponse `json:"data"`
}

type ApiResponseTransactionDeleteAt struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    *TransactionResponseDeleteAt `json:"data"`
}

type ApiResponsesTransaction struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []*TransactionResponse `json:"data"`
}

type ApiResponseTransactionDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseTransactionAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationTransactionDeleteAt struct {
	Status     string                         `json:"status"`
	Message    string                         `json:"message"`
	Data       []*TransactionResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                 `json:"pagination"`
}

type ApiResponsePaginationTransaction struct {
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Data       []*TransactionResponse `json:"data"`
	Pagination PaginationMeta         `json:"pagination"`
}

type ApiResponsesTransactionMonthSuccess struct {
	Status  string                                     `json:"status"`
	Message string                                     `json:"message"`
	Data    []*TransactionMonthlyAmountSuccessResponse `json:"data"`
}

type ApiResponsesTransactionMonthFailed struct {
	Status  string                                    `json:"status"`
	Message string                                    `json:"message"`
	Data    []*TransactionMonthlyAmountFailedResponse `json:"data"`
}

type ApiResponsesTransactionYearSuccess struct {
	Status  string                                    `json:"status"`
	Message string                                    `json:"message"`
	Data    []*TransactionYearlyAmountSuccessResponse `json:"data"`
}

type ApiResponsesTransactionYearFailed struct {
	Status  string                                   `json:"status"`
	Message string                                   `json:"message"`
	Data    []*TransactionYearlyAmountFailedResponse `json:"data"`
}

type ApiResponsesTransactionMonthMethod struct {
	Status  string                              `json:"status"`
	Message string                              `json:"message"`
	Data    []*TransactionMonthlyMethodResponse `json:"data"`
}

type ApiResponsesTransactionYearMethod struct {
	Status  string                             `json:"status"`
	Message string                             `json:"message"`
	Data    []*TransactionYearlyMethodResponse `json:"data"`
}
