package response

type MerchantPoliciesResponse struct {
	ID           int    `json:"id"`
	MerchantID   int    `json:"merchant_id"`
	PolicyType   string `json:"policy_type"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	MerchantName string `json:"merchant_name"`
}

type MerchantPoliciesResponseDeleteAt struct {
	ID           int     `json:"id"`
	MerchantID   int     `json:"merchant_id"`
	PolicyType   string  `json:"policy_type"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
	MerchantName string  `json:"merchant_name"`
}

type ApiResponseMerchantPolicies struct {
	Status  string                    `json:"status"`
	Message string                    `json:"message"`
	Data    *MerchantPoliciesResponse `json:"data"`
}

type ApiResponseMerchantPoliciesDeleteAt struct {
	Status  string                            `json:"status"`
	Message string                            `json:"message"`
	Data    *MerchantPoliciesResponseDeleteAt `json:"data"`
}

type ApiResponsesMerchantPolicies struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    []*MerchantPoliciesResponse `json:"data"`
}

type ApiResponsePaginationMerchantPoliciesDeleteAt struct {
	Status     string                              `json:"status"`
	Message    string                              `json:"message"`
	Data       []*MerchantPoliciesResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                      `json:"pagination"`
}

type ApiResponsePaginationMerchantPolicies struct {
	Status     string                      `json:"status"`
	Message    string                      `json:"message"`
	Data       []*MerchantPoliciesResponse `json:"data"`
	Pagination PaginationMeta              `json:"pagination"`
}
