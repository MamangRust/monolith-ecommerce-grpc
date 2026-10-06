package service

import (
	mencache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
)

// Service aggregates authentication and identity-related services.
type Service struct {
	Login         LoginService
	Register      RegistrationService
	PasswordReset PasswordResetService
	Identify      IdentifyService
}

// Deps defines dependencies required to initialize Service.
type Deps struct {
	Mencache      *mencache.Mencache
	Repositories  *repository.Repositories
	Token         auth.TokenManager
	Hash          hash.HashPassword
	Logger        logger.LoggerInterface
	Kafka         *kafka.Kafka
	Observability observability.TraceLoggerObservability
}

// NewService initializes and returns the core authentication service bundle.
func NewService(deps *Deps) *Service {

	tokenService := NewTokenService(&tokenServiceDeps{
		Token:         deps.Token,
		RefreshToken:  deps.Repositories.RefreshToken,
		Logger:        deps.Logger,
		Observability: deps.Observability,
	})

	return &Service{
		Login:         newLogin(deps, tokenService, deps.Observability, deps.Mencache.LoginCache),
		Register:      newRegister(deps, deps.Observability, deps.Mencache.RegisterCache),
		PasswordReset: newPasswordReset(deps, deps.Observability, deps.Mencache.PasswordResetCache),
		Identify:      newIdentity(deps, tokenService, deps.Observability, deps.Mencache.IdentityCache),
	}
}

// newLogin initializes and returns the LoginService.
func newLogin(deps *Deps, tokenService *tokenService, observability observability.TraceLoggerObservability, cache mencache.LoginCache) LoginService {
	return NewLoginService(&LoginServiceDeps{
		Cache:          cache,
		Logger:         deps.Logger,
		Hash:           deps.Hash,
		UserRepository: deps.Repositories.User,
		RefreshToken:   deps.Repositories.RefreshToken,
		Token:          deps.Token,
		TokenService:   tokenService,
		Observability:  observability,
	})
}

// newRegister initializes and returns the RegistrationService.
func newRegister(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.RegisterCache) RegistrationService {
	return NewRegisterService(&RegisterServiceDeps{
		Cache:         cache,
		User:          deps.Repositories.User,
		Role:          deps.Repositories.Role,
		UserRole:      deps.Repositories.UserRole,
		Hash:          deps.Hash,
		Kafka:         deps.Kafka,
		Logger:        deps.Logger,
		Observability: observability,
	})
}

// newPasswordReset initializes the reset forgot password, reset password and verify code services.
func newPasswordReset(deps *Deps, observability observability.TraceLoggerObservability, cache mencache.PasswordResetCache) PasswordResetService {
	return NewPasswordResetService(&PasswordResetServiceDeps{
		Cache:         cache,
		Kafka:         deps.Kafka,
		Logger:        deps.Logger,
		User:          deps.Repositories.User,
		ResetToken:    deps.Repositories.ResetToken,
		Observability: observability,
	})
}

// newIdentity initializes the IdentifyService for identity verification and token refresh.
func newIdentity(deps *Deps, tokenService *tokenService, observability observability.TraceLoggerObservability, cache mencache.IdentityCache) IdentifyService {
	return NewIdentityService(&IdentityServiceDeps{
		Cache:         cache,
		Token:         deps.Token,
		RefreshToken:  deps.Repositories.RefreshToken,
		User:          deps.Repositories.User,
		Logger:        deps.Logger,
		TokenService:  tokenService,
		Observability: observability,
	})
}
