package record

type OrderRecord struct {
	ID         int     `json:"id"`
	MerchantID int     `json:"merchant_id"`
	UserID     int     `json:"user_id"`
	TotalPrice int     `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

type OrderMonthlyRecord struct {
	Month          string `json:"month"`
	OrderCount     int    `json:"order_count"`
	TotalRevenue   int    `json:"total_revenue"`
	TotalItemsSold int    `json:"total_items_sold"`
}

type OrderYearlyRecord struct {
	Year               string `json:"year"`
	OrderCount         int    `json:"order_count"`
	TotalRevenue       int    `json:"total_revenue"`
	TotalItemsSold     int    `json:"total_items_sold"`
	ActiveCashiers     int    `json:"active_cashiers"`
	UniqueProductsSold int    `json:"unique_products_sold"`
}

type OrderMonthlyTotalRevenueRecord struct {
	Year         string `json:"year"`
	Month        string `json:"month"`
	TotalRevenue int    `json:"total_revenue"`
}

type OrderYearlyTotalRevenueRecord struct {
	Year         string `json:"year"`
	TotalRevenue int    `json:"total_revenue"`
}
