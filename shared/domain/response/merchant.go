package response

type MerchantResponse struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Address      string `json:"address"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type MerchantResponseDeleteAt struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Address      string  `json:"address"`
	ContactEmail string  `json:"contact_email"`
	ContactPhone string  `json:"contact_phone"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
}

type ApiResponseMerchant struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *MerchantResponse `json:"data"`
}

type ApiResponseMerchantDeleteAt struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    *MerchantResponseDeleteAt `json:"data"`
}

type ApiResponsesMerchant struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []*MerchantResponse `json:"data"`
}

type ApiResponseMerchantDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseMerchantAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationMerchantDeleteAt struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*MerchantResponseDeleteAt `json:"data"`
	Pagination PaginationMeta              `json:"pagination"`
}

type ApiResponsePaginationMerchant struct {
	Status     string              `json:"status"`
	Message    string              `json:"message"`
	Data       []*MerchantResponse `json:"data"`
	Pagination PaginationMeta      `json:"pagination"`
}
