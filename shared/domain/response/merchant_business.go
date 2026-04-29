package response

type MerchantBusinessResponse struct {
	ID                int     `json:"id"`
	MerchantID        int     `json:"merchant_id"`
	BusinessType      string  `json:"business_type"`
	TaxID             string  `json:"tax_id"`
	EstablishedYear   int     `json:"established_year"`
	NumberOfEmployees int     `json:"number_of_employees"`
	WebsiteUrl        string  `json:"website_url"`
	MerchantName      *string `json:"merchant_name"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type MerchantBusinessResponseDeleteAt struct {
	ID                int     `json:"id"`
	MerchantID        int     `json:"merchant_id"`
	BusinessType      string  `json:"business_type"`
	TaxID             string  `json:"tax_id"`
	EstablishedYear   int     `json:"established_year"`
	NumberOfEmployees int     `json:"number_of_employees"`
	WebsiteUrl        string  `json:"website_url"`
	MerchantName      string  `json:"merchant_name"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	DeletedAt         *string `json:"deleted_at"`
}

type ApiResponseMerchantBusiness struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    *MerchantBusinessResponse `json:"data"`
}

type ApiResponseMerchantBusinessDeleteAt struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    *MerchantBusinessResponseDeleteAt `json:"data"`
}

type ApiResponsesMerchantBusiness struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*MerchantBusinessResponse `json:"data"`
}

type ApiResponsePaginationMerchantBusinessDeleteAt struct {
	Status     string                              `json:"status"`
	Message    string                              `json:"message"`
	Data       []*MerchantBusinessResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                      `json:"pagination"`
}

type ApiResponsePaginationMerchantBusiness struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*MerchantBusinessResponse `json:"data"`
	Pagination PaginationMeta              `json:"pagination"`
}
