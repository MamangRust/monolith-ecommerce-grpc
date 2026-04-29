package response

type BannerResponse struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BannerResponseDeleteAt struct {
	ID        int32   `json:"id"`
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

type ApiResponseBanner struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    *BannerResponse `json:"data"`
}

type ApiResponseBannerDeleteAt struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    *BannerResponseDeleteAt `json:"data"`
}

type ApiResponsesBanner struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []*BannerResponse `json:"data"`
}

type ApiResponseBannerDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseBannerAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationBannerDeleteAt struct {
	Status     string                    `json:"status"`
	Message    string                    `json:"message"`
	Data       []*BannerResponseDeleteAt `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}

type ApiResponsePaginationBanner struct {
	Status     string            `json:"status"`
	Message    string            `json:"message"`
	Data       []*BannerResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}
