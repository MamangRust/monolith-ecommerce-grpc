package record

import "time"

type UserRoleRecord struct {
	UserRoleID int32      `json:"user_role_id"`
	UserID     int32      `json:"user_id"`
	RoleID     int32      `json:"role_id"`
	RoleName   string     `json:"role_name,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
