package response

type ProductResponse struct {
	ID           int     `json:"id"`
	MerchantID   int     `json:"merchant_id"`
	CategoryID   int     `json:"category_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        int     `json:"price"`
	CountInStock int     `json:"count_in_stock"`
	Brand        string  `json:"brand"`
	Weight       int     `json:"weight"`
	Rating       float32 `json:"rating"`
	SlugProduct  string  `json:"slug_product"`
	ImageProduct string  `json:"image_product"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type ProductResponseDeleteAt struct {
	ID           int     `json:"id"`
	MerchantID   int     `json:"merchant_id"`
	CategoryID   int     `json:"category_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        int     `json:"price"`
	CountInStock int     `json:"count_in_stock"`
	Brand        string  `json:"brand"`
	Weight       int     `json:"weight"`
	Rating       float32 `json:"rating"`
	SlugProduct  string  `json:"slug_product"`
	ImageProduct string  `json:"image_product"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
}

type ApiResponseProduct struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *ProductResponse `json:"data"`
}

type ApiResponseProductDeleteAt struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    *ProductResponseDeleteAt `json:"data"`
}

type ApiResponsesProduct struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []*ProductResponse `json:"data"`
}

type ApiResponseProductDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseProductAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationProductDeleteAt struct {
	Status     string                     `json:"status"`
	Message    string                     `json:"message"`
	Data       []*ProductResponseDeleteAt `json:"data"`
	Pagination PaginationMeta             `json:"pagination"`
}

type ApiResponsePaginationProduct struct {
	Status     string             `json:"status"`
	Message    string             `json:"message"`
	Data       []*ProductResponse `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}
