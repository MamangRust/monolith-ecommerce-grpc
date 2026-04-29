package response

type ReviewDetailsResponse struct {
	ID        int    `json:"id"`
	ReviewID  int    `json:"review_id"`
	Type      string `json:"type"`
	Url       string `json:"url"`
	Caption   string `json:"caption"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ReviewDetailsResponseDeleteAt struct {
	ID        int     `json:"id"`
	ReviewID  int     `json:"review_id"`
	Type      string  `json:"type"`
	Url       string  `json:"url"`
	Caption   string  `json:"caption"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type ApiResponseReviewDetail struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    *ReviewDetailsResponse `json:"data"`
}

type ApiResponseReviewDetailDeleteAt struct {
	Status  string                         `json:"status"`
	Message string                         `json:"message"`
	Data    *ReviewDetailsResponseDeleteAt `json:"data"`
}

type ApiResponsesReviewDetails struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []*ReviewDetailsResponse `json:"data"`
}

type ApiResponsePaginationReviewDetailsDeleteAt struct {
	Status     string                           `json:"status"`
	Message    string                           `json:"message"`
	Data       []*ReviewDetailsResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                   `json:"pagination"`
}

type ApiResponsePaginationReviewDetails struct {
	Status     string                   `json:"status"`
	Message    string                   `json:"message"`
	Data       []*ReviewDetailsResponse `json:"data"`
	Pagination PaginationMeta           `json:"pagination"`
}
