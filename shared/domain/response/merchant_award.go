package response

type MerchantAwardResponse struct {
	ID             int    `json:"id"`
	MerchantID     int    `json:"merchant_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	IssuedBy       string `json:"issued_by"`
	IssueDate      string `json:"issue_date"`
	ExpiryDate     string `json:"expiry_date"`
	CertificateUrl string `json:"certificate_url"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	MerchantName   string `json:"merchant_name"`
}

type MerchantAwardResponseDeleteAt struct {
	ID             int     `json:"id"`
	MerchantID     int     `json:"merchant_id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	IssuedBy       string  `json:"issued_by"`
	IssueDate      string  `json:"issue_date"`
	ExpiryDate     string  `json:"expiry_date"`
	CertificateUrl string  `json:"certificate_url"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	MerchantName   string  `json:"merchant_name"`
	DeletedAt      *string `json:"deleted_at"`
}

type ApiResponseMerchantAward struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    *MerchantAwardResponse `json:"data"`
}

type ApiResponseMerchantAwardDeleteAt struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    *MerchantAwardResponseDeleteAt `json:"data"`
}

type ApiResponsesMerchantAward struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []*MerchantAwardResponse `json:"data"`
}

type ApiResponsePaginationMerchantAwardDeleteAt struct {
	Status     string                           `json:"status"`
	Message    string                           `json:"message"`
	Data       []*MerchantAwardResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                   `json:"pagination"`
}

type ApiResponsePaginationMerchantAward struct {
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Data       []*MerchantAwardResponse `json:"data"`
	Pagination PaginationMeta           `json:"pagination"`
}
