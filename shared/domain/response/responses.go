package response

type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

type ApiResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type APIResponsePagination[T any] struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    T              `json:"data"`
	Meta    PaginationMeta `json:"pagination"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
