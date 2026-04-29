package record

type MerchantRecord struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Address      string  `json:"address"`
	ContactEmail string  `json:"contact_email"`
	ContactPhone string  `json:"contact_phone"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DeletedAt    *string `json:"deleted_at"`
}

type MerchantDetailRecord struct {
	ID               int                              `json:"id"`
	MerchantID       int                              `json:"merchant_id"`
	DisplayName      string                           `json:"display_name"`
	CoverImageUrl    string                           `json:"cover_image_url"`
	LogoUrl          string                           `json:"logo_url"`
	ShortDescription string                           `json:"short_description"`
	WebsiteUrl       string                           `json:"website_url"`
	SocialMediaLinks []*MerchantSocialMediaLinkRecord `json:"social_media_links"`
	CreatedAt        string                           `json:"created_at"`
	UpdatedAt        string                           `json:"updated_at"`
	DeletedAt        *string                          `json:"deleted_at"`
}

type MerchantSocialMediaLinkRecord struct {
	ID       int    `json:"id"`
	Platform string `json:"platform"`
	Url      string `json:"url"`
}

type MerchantBusinessRecord struct {
	ID                int     `json:"id"`
	MerchantID        int     `json:"merchant_id"`
	BusinessType      string  `json:"business_type"`
	TaxID             string  `json:"tax_id"`
	EstablishedYear   int     `json:"established_year"`
	NumberOfEmployees int     `json:"number_of_employees"`
	WebsiteUrl        string  `json:"website_url"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	DeletedAt         *string `json:"deleted_at"`
}

type MerchantPoliciesRecord struct {
	ID          int     `json:"id"`
	MerchantID  int     `json:"merchant_id"`
	PolicyType  string  `json:"policy_type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at"`
}

type MerchantAwardRecord struct {
	ID             int     `json:"id"`
	MerchantID     int     `json:"merchant_id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	IssuedBy       string  `json:"issued_by"`
	IssueDate      string  `json:"issue_date"`
	ExpiryDate     string  `json:"expiry_date"`
	CertificateUrl string  `json:"certificate_url"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	DeletedAt      *string `json:"deleted_at"`
}
