package record

type MerchantSocialLinkRecord struct {
	ID               int     `json:"id"`
	MerchantDetailID int     `json:"merchant_detail_id"`
	Platform         string  `json:"platform"`
	URL              string  `json:"url"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	DeletedAt        *string `json:"deleted_at"`
}
