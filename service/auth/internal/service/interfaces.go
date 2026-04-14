package service

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go
type RegistrationService interface {
	Register(ctx context.Context, request *requests.RegisterRequest) (*db.CreateUserRow, error)
}

type LoginService interface {
	Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, error)
}

type PasswordResetService interface {
	ForgotPassword(ctx context.Context, email string) (bool, error)

	ResetPassword(ctx context.Context, request *requests.CreateResetPasswordRequest) (bool, error)

	VerifyCode(ctx context.Context, code string) (bool, error)
}

type IdentifyService interface {
	RefreshToken(ctx context.Context, token string) (*response.TokenResponse, error)

	GetMe(ctx context.Context, userId int) (*db.GetUserByIDRow, error)
}
