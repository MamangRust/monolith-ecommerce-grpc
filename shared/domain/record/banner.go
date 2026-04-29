package record

type BannerRecord struct {
	BannerID  int     `json:"banner_id"`
	Name      string  `json:"name"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	IsActive  bool    `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}
