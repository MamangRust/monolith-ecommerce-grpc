package response

type SliderResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SliderResponseDeleteAt struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Image     string  `json:"image"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
}

type ApiResponseSlider struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    *SliderResponse `json:"data"`
}

type ApiResponseSliderDeleteAt struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Data    *SliderResponseDeleteAt `json:"data"`
}

type ApiResponsesSlider struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []*SliderResponse `json:"data"`
}

type ApiResponseSliderDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseSliderAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationSliderDeleteAt struct {
	Status     string                    `json:"status"`
	Message    string                    `json:"message"`
	Data       []*SliderResponseDeleteAt `json:"data"`
	Pagination PaginationMeta            `json:"pagination"`
}

type ApiResponsePaginationSlider struct {
	Status     string            `json:"status"`
	Message    string            `json:"message"`
	Data       []*SliderResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}
