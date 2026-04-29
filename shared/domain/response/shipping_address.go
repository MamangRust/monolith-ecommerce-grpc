package response

type ShippingAddressResponse struct {
	ID             int    `json:"id"`
	OrderID        int    `json:"order_id"`
	Alamat         string `json:"alamat"`
	Provinsi       string `json:"provinsi"`
	Negara         string `json:"negara"`
	Kota           string `json:"kota"`
	ShippingMethod string `json:"shipping_method"`
	ShippingCost   int    `json:"shipping_cost"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type ShippingAddressResponseDeleteAt struct {
	ID             int     `json:"id"`
	OrderID        int     `json:"order_id"`
	Alamat         string  `json:"alamat"`
	Provinsi       string  `json:"provinsi"`
	Negara         string  `json:"negara"`
	Kota           string  `json:"kota"`
	ShippingMethod string  `json:"shipping_method"`
	ShippingCost   int     `json:"shipping_cost"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}

type ApiResponseShippingAddress struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    *ShippingAddressResponse `json:"data"`
}

type ApiResponseShippingAddressDeleteAt struct {
	Status  string                           `json:"status"`
	Message string                           `json:"message"`
	Data    *ShippingAddressResponseDeleteAt `json:"data"`
}

type ApiResponsesShippingAddress struct {
	Status  string                     `json:"status"`
	Message string                     `json:"message"`
	Data    []*ShippingAddressResponse `json:"data"`
}

type ApiResponseShippingAddressDelete struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponseShippingAddressAll struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ApiResponsePaginationShippingAddressDeleteAt struct {
	Status     string                             `json:"status"`
	Message    string                             `json:"message"`
	Data       []*ShippingAddressResponseDeleteAt `json:"data"`
	Pagination PaginationMeta                     `json:"pagination"`
}

type ApiResponsePaginationShippingAddress struct {
	Status     string                     `json:"status"`
	Message    string                     `json:"message"`
	Data       []*ShippingAddressResponse `json:"data"`
	Pagination PaginationMeta             `json:"pagination"`
}
