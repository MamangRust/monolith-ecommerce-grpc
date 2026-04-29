package response

// TODO: Lanjutkan lagi di proto, database, response(record, response(service,api), dan proto)
type MerchantSocialLinkResponse struct {
	ID               int    `json:"id"`
	MerchantDetailID int    `json:"merchant_detail_id"`
	Platform         string `json:"platform"`
	URL              string `json:"url"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type MerchantSocialLinkResponseDeleteAt struct {
	ID               int     `json:"id"`
	MerchantDetailID int     `json:"merchant_detail_id"`
	Platform         string  `json:"platform"`
	URL              string  `json:"url"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	DeletedAt        *string `json:"deleted_at"`
}

type ApiResponsesMerchantSocialLink struct {
	Status  string                        `json:"status"`
	Message string                        `json:"message"`
	Data    []*MerchantSocialLinkResponse `json:"data"`
}

type ApiResponseMerchantSocialLink struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    *MerchantSocialLinkResponse `json:"data"`
}

type ApiResponseMerchantSocialLinkDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseMerchantSocialLinkAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationMerchantSocialLink struct {
	Status     string                        `json:"status"`
	Message    string                        `json:"message"`
	Data       []*MerchantSocialLinkResponse `json:"data"`
	Pagination *PaginationMeta               `json:"pagination"`
}

type ApiResponsePaginationMerchantSocialLinkDeleteAt struct {
	Status     string                                `json:"status"`
	Message    string                                `json:"message"`
	Data       []*MerchantSocialLinkResponseDeleteAt `json:"data"`
	Pagination *PaginationMeta                       `json:"pagination"`
}
