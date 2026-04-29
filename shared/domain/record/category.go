package record

type CategoriesRecord struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	SlugCategory  string  `json:"slug_category"`
	ImageCategory string  `json:"image_category"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     *string `json:"deleted_at"`
}

type CategoriesMonthPriceRecord struct {
	Month        string `json:"month"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	OrderCount   int    `json:"order_count"`
	ItemsSold    int    `json:"items_sold"`
	TotalRevenue int    `json:"total_revenue"`
}

type CategoriesYearPriceRecord struct {
	Year               string `json:"year"`
	CategoryID         int    `json:"category_id"`
	CategoryName       string `json:"category_name"`
	OrderCount         int    `json:"order_count"`
	ItemsSold          int    `json:"items_sold"`
	TotalRevenue       int    `json:"total_revenue"`
	UniqueProductsSold int    `json:"unique_products_sold"`
}

type CategoriesMonthlyTotalPriceRecord struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalRevenue int    `json:"total_revenue"`
}

type CategoriesYearlyTotalPriceRecord struct {
	Year         string `json:"year"`
	TotalRevenue int    `json:"total_revenue"`
}
