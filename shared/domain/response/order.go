package response

type OrderResponse struct {
	ID         int    `json:"id"`
	MerchantID int    `json:"merchant_id"`
	UserID     int    `json:"user_id"`
	TotalPrice int    `json:"total_price"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type OrderResponseDeleteAt struct {
	ID         int     `json:"id"`
	MerchantID int     `json:"merchant_id"`
	UserID     int     `json:"user_id"`
	TotalPrice int     `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

type OrderMonthlyResponse struct {
	Month          string `json:"month"`
	OrderCount     int    `json:"order_count"`
	TotalRevenue   int    `json:"total_revenue"`
	TotalItemsSold int    `json:"total_items_sold"`
}

type OrderYearlyResponse struct {
	Year               string `json:"year"`
	OrderCount         int    `json:"order_count"`
	TotalRevenue       int    `json:"total_revenue"`
	TotalItemsSold     int    `json:"total_items_sold"`
	ActiveCashiers     int    `json:"active_cashiers"`
	UniqueProductsSold int    `json:"unique_products_sold"`
}

type OrderMonthlyTotalRevenueResponse struct {
	Year           string `json:"year"`
	Month          string `json:"month"`
	OrderCount     int    `json:"order_count"`
	TotalRevenue   int    `json:"total_revenue"`
	TotalItemsSold int    `json:"total_items_sold"`
}

type OrderYearlyTotalRevenueResponse struct {
	Year         string `json:"year"`
	OrderCount   int    `json:"order_count"`
	TotalRevenue int    `json:"total_revenue"`
}

type ApiResponseOrder struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *OrderResponse `json:"data"`
}

type ApiResponseOrderDeleteAt struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    *OrderResponseDeleteAt `json:"data"`
}

type ApiResponsesOrder struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*OrderResponse `json:"data"`
}

type ApiResponseOrderDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseOrderAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationOrderDeleteAt struct {
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Data       []*OrderResponseDeleteAt `json:"data"`
	Pagination PaginationMeta           `json:"pagination"`
}

type ApiResponsePaginationOrder struct {
	Status     string           `json:"status"`
	Message    string           `json:"message"`
	Data       []*OrderResponse `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

type ApiResponseOrderMonthly struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    []*OrderMonthlyResponse `json:"data"`
}

type ApiResponseOrderYearly struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []*OrderYearlyResponse `json:"data"`
}

type ApiResponseOrderMonthlyTotalRevenue struct {
	Status  string                              `json:"status"`
	Message string                              `json:"message"`
	Data    []*OrderMonthlyTotalRevenueResponse `json:"data"`
}

type ApiResponseOrderYearlyTotalRevenue struct {
	Status  string                             `json:"status"`
	Message string                             `json:"message"`
	Data    []*OrderYearlyTotalRevenueResponse `json:"data"`
}
