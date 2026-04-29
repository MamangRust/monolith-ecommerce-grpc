package response

type OrderItemResponse struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type OrderItemResponseDeleteAt struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     int     `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type ApiResponseOrderItem struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    *OrderItemResponse `json:"data"`
}

type ApiResponsesOrderItem struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    []*OrderItemResponse `json:"data"`
}

type ApiResponseOrderItemDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseOrderItemAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationOrderItemDeleteAt struct {
	Status     string                       `json:"status"`
	Message    string                       `json:"message"`
	Data       []*OrderItemResponseDeleteAt `json:"data"`
	Pagination PaginationMeta               `json:"pagination"`
}

type ApiResponsePaginationOrderItem struct {
	Status     string               `json:"status"`
	Message    string               `json:"message"`
	Data       []*OrderItemResponse `json:"data"`
	Pagination PaginationMeta       `json:"pagination"`
}
