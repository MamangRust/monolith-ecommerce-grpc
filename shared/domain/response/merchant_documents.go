package response

type MerchantDocumentResponse struct {
	ID           int    `json:"id"`
	MerchantID   int    `json:"merchant_id"`
	DocumentType string `json:"document_type"`
	DocumentURL  string `json:"document_url"`
	Status       string `json:"status"`
	Note         string `json:"note"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type MerchantDocumentResponseDeleteAt struct {
	ID           int     `json:"id"`
	MerchantID   int     `json:"merchant_id"`
	DocumentType string  `json:"document_type"`
	DocumentURL  string  `json:"document_url"`
	Status       string  `json:"status"`
	Note         string  `json:"note"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
}

type ApiResponsesMerchantDocument struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*MerchantDocumentResponse `json:"data"`
}

type ApiResponseMerchantDocument struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    *MerchantDocumentResponse `json:"data"`
}

type ApiResponseMerchantDocumentDeleteAt struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    *MerchantDocumentResponseDeleteAt `json:"data"`
}

type ApiResponseMerchantDocumentDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseMerchantDocumentAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationMerchantDocument struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*MerchantDocumentResponse `json:"data"`
	Pagination *PaginationMeta             `json:"pagination"`
}

type ApiResponsePaginationMerchantDocumentDeleteAt struct {
	Status     string                              `json:"status"`
	Message    string                              `json:"message"`
	Data       []*MerchantDocumentResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta                     `json:"pagination"`
}
