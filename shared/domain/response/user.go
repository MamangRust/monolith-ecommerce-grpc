package response

type UserResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserResponseDeleteAt struct {
	ID        int     `json:"id"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type ApiResponseUser struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *UserResponse `json:"data"`
}

type ApiResponseUserDeleteAt struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    *UserResponseDeleteAt `json:"data"`
}

type ApiResponsesUser struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    []*UserResponse `json:"data"`
}

type ApiResponseUserDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseUserAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationUserDeleteAt struct {
	Status     string                  `json:"status"`
	Message    string                  `json:"message"`
	Data       []*UserResponseDeleteAt `json:"data"`
	Pagination PaginationMeta          `json:"pagination"`
}

type ApiResponsePaginationUser struct {
	Status     string          `json:"status"`
	Message    string          `json:"message"`
	Data       []*UserResponse `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
}
