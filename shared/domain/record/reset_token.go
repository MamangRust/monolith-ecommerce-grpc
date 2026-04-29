package record

type ResetTokenRecord struct {
	ID        int64  `json:"id"`
	Token     string `json:"token"`
	UserID    int64  `json:"user_id"`
	ExpiredAt string `json:"expired_at"`
}
