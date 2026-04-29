package requests

type CreateUserRoleRequest struct {
	UserId int `json:"user_id" validate:"required"`
	RoleId int `json:"role_id" validate:"required"`
}

type RemoveUserRoleRequest struct {
	UserId int `json:"user_id" validate:"required"`
	RoleId int `json:"role_id" validate:"required"`
}
