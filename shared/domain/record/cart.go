package record

type CartRecord struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     int     `json:"price"`
	Image     string  `json:"image"`
	Quantity  int     `json:"quantity"`
	Weight    int     `json:"weight"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}
