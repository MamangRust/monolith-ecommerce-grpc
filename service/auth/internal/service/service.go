package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
)

type Service struct {
	Login         LoginService
	Register      RegistrationService
	PasswordReset PasswordResetService
	Identify      IdentifyService
}

type Deps struct {
	Context      context.Context
	Repositories *repository.Repositories
	Token        auth.TokenManager
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
	Kafka        kafka.Kafka
	Mapper       response_service.UserResponseMapper
}

func NewService(deps Deps) *Service {
	tokenService := NewTokenService(deps.Repositories.RefreshToken, deps.Token, deps.Logger)

	return &Service{
		Login:         NewLoginService(deps.Context, deps.Logger, deps.Hash, deps.Repositories.User, deps.Repositories.RefreshToken, deps.Token, *tokenService),
		Register:      NewRegisterService(deps.Context, deps.Repositories.User, deps.Repositories.Role, deps.Repositories.UserRole, deps.Hash, deps.Kafka, deps.Logger, deps.Mapper),
		PasswordReset: NewPasswordResetService(deps.Context, deps.Kafka, deps.Logger, deps.Repositories.User, deps.Repositories.ResetToken),
		Identify:      NewIdentityService(deps.Context, deps.Token, deps.Repositories.RefreshToken, deps.Repositories.User, deps.Logger, deps.Mapper, *tokenService),
	}
}
