package record

type ReviewRecord struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Comment   string  `json:"comment"`
	Rating    int     `json:"rating"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type ReviewsDetailRecord struct {
	ID           int                 `json:"id"`
	UserID       int                 `json:"user_id"`
	ProductID    int                 `json:"product_id"`
	Name         string              `json:"name"`
	Comment      string              `json:"comment"`
	Rating       int                 `json:"rating"`
	ReviewDetail *ReviewDetailRecord `json:"review_detail"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
	DeletedAt    *string             `json:"deleted_at"`
}

type ReviewDetailRecord struct {
	ID        int     `json:"id"`
	ReviewID  *int    `json:"review_id"`
	Type      string  `json:"type"`
	Url       string  `json:"url"`
	Caption   string  `json:"string"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}
