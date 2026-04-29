package record

type ShippingAddressRecord struct {
	ID             int     `json:"id"`
	OrderID        int     `json:"order_id"`
	Alamat         string  `json:"alamat"`
	Provinsi       string  `json:"provinsi"`
	Negara         string  `json:"negara"`
	Kota           string  `json:"kota"`
	Courier        string  `json:"courier"`
	ShippingMethod string  `json:"shipping_method"`
	ShippingCost   int     `json:"shipping_cost"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}
