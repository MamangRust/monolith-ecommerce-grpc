package record

type SliderRecord struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}
