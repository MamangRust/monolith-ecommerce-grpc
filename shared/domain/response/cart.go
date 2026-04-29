package response

type CartResponse struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Image     string `json:"image"`
	Quantity  int    `json:"quantity"`
	Weight    int    `json:"weight"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ApiResponseCartPagination struct {
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Data       []*CartResponse `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
}

type ApiResponseCart struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Data    CartResponse `json:"data"`
}

type ApiResponseCartAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseCartDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
