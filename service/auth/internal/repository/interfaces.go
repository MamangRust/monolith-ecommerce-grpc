package repository

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*record.UserRecord, error)
	FindById(ctx context.Context, id int) (*record.UserRecord, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*record.UserRecord, error)
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*record.UserRecord, error)
	UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*record.UserRecord, error)
	UpdateUserPassword(ctx context.Context, user_id int, password string) (*record.UserRecord, error)
	FindByVerificationCode(ctx context.Context, verification_code string) (*record.UserRecord, error)
}

type ResetTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*record.ResetTokenRecord, error)
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*record.ResetTokenRecord, error)
	DeleteResetToken(ctx context.Context, user_id int) error
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*record.RefreshTokenRecord, error)
	FindByUserId(ctx context.Context, user_id int) (*record.RefreshTokenRecord, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*record.RefreshTokenRecord, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*record.RefreshTokenRecord, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

type UserRoleRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

type RoleRepository interface {
	FindById(ctx context.Context, role_id int) (*record.RoleRecord, error)
	FindByName(ctx context.Context, name string) (*record.RoleRecord, error)
}
