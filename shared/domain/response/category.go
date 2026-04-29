package response

type CategoryResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	SlugCategory  string `json:"slug_category"`
	ImageCategory string `json:"image_category"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CategoryResponseDeleteAt struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	SlugCategory  string  `json:"slug_category"`
	ImageCategory string  `json:"image_category"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at"`
}

type CategoryMonthPriceResponse struct {
	Month        string `json:"month"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	OrderCount   int    `json:"order_count"`
	ItemsSold    int    `json:"items_sold"`
	TotalRevenue int    `json:"total_revenue"`
}

type CategoryYearPriceResponse struct {
	Year               string `json:"year"`
	CategoryID         int    `json:"category_id"`
	CategoryName       string `json:"category_name"`
	OrderCount         int    `json:"order_count"`
	ItemsSold          int    `json:"items_sold"`
	TotalRevenue       int    `json:"total_revenue"`
	UniqueProductsSold int    `json:"unique_products_sold"`
}

type CategoriesMonthlyTotalPriceResponse struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalRevenue int    `json:"total_revenue"`
}

type CategoriesYearlyTotalPriceResponse struct {
	Year         string `json:"year"`
	TotalRevenue int    `json:"total_revenue"`
}

type ApiResponseCategory struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *CategoryResponse `json:"data"`
}

type ApiResponseCategoryDeleteAt struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    *CategoryResponseDeleteAt `json:"data"`
}

type ApiResponsesCategory struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []*CategoryResponse `json:"data"`
}

type ApiResponseCategoryDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseCategoryAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationCategoryDeleteAt struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*CategoryResponseDeleteAt `json:"data"`
	Pagination PaginationMeta              `json:"pagination"`
}

type ApiResponsePaginationCategory struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Data       []*CategoryResponse `json:"data"`
	Pagination PaginationMeta      `json:"pagination"`
}

type ApiResponseCategoryMonthPrice struct {
	Status  string                        `json:"status"`
	Message string                        `json:"message"`
	Data    []*CategoryMonthPriceResponse `json:"data"`
}

type ApiResponseCategoryYearPrice struct {
	Status  string                       `json:"status"`
	Message string                       `json:"message"`
	Data    []*CategoryYearPriceResponse `json:"data"`
}

type ApiResponseCategoryMonthlyTotalPrice struct {
	Status  string                                 `json:"status"`
	Message string                                 `json:"message"`
	Data    []*CategoriesMonthlyTotalPriceResponse `json:"data"`
}

type ApiResponseCategoryYearlyTotalPrice struct {
	Status  string                                `json:"status"`
	Message string                                `json:"message"`
	Data    []*CategoriesYearlyTotalPriceResponse `json:"data"`
}
