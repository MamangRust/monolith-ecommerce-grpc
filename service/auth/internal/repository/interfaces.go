package repository

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/pb/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb/domain/requests"
)

type UserRepository interface {
	FindByEmail(email string) (*record.UserRecord, error)
	FindById(id int) (*record.UserRecord, error)
	CreateUser(request *requests.RegisterRequest) (*record.UserRecord, error)
	UpdateUserIsVerified(user_id int, is_verified bool) (*record.UserRecord, error)
	UpdateUserPassword(user_id int, password string) (*record.UserRecord, error)
	FindByVerificationCode(verification_code string) (*record.UserRecord, error)
}

type ResetTokenRepository interface {
	FindByToken(token string) (*record.ResetTokenRecord, error)
	CreateResetToken(req *requests.CreateResetTokenRequest) (*record.ResetTokenRecord, error)
	DeleteResetToken(user_id int) error
}

type RefreshTokenRepository interface {
	FindByToken(token string) (*record.RefreshTokenRecord, error)
	FindByUserId(user_id int) (*record.RefreshTokenRecord, error)
	CreateRefreshToken(req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error)
	UpdateRefreshToken(req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByUserId(user_id int) error
}

type UserRoleRepository interface {
	AssignRoleToUser(req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error)
	RemoveRoleFromUser(req *requests.RemoveUserRoleRequest) error
}

type RoleRepository interface {
	FindById(role_id int) (*record.RoleRecord, error)
	FindByName(name string) (*record.RoleRecord, error)
}
