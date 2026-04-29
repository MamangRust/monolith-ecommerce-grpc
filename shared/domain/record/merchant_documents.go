package record

type MerchantDocumentRecord struct {
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
