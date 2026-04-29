package record

type UserRecord struct {
	ID              int     `json:"id"`
	FirstName       string  `json:"firstname"`
	LastName        string  `json:"lastname"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       *string `json:"deleted_at"`
}
