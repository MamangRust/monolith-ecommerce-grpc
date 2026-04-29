package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"

	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	sharederrorhandler "github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	user_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// LoginServiceDeps defines all dependencies required by LoginService.
type LoginServiceDeps struct {
	Cache  mencache.LoginCache
	Logger logger.LoggerInterface
	Hash   hash.HashPassword

	UserRepository repository.UserRepository
	RefreshToken   repository.RefreshTokenRepository

	Token        auth.TokenManager
	TokenService *tokenService

	Observability observability.TraceLoggerObservability
}

type loginService struct {
	mencache mencache.LoginCache
	logger   logger.LoggerInterface
	hash     hash.HashPassword

	user         repository.UserRepository
	refreshToken repository.RefreshTokenRepository

	token        auth.TokenManager
	tokenService *tokenService

	observability observability.TraceLoggerObservability
}

func NewLoginService(params *LoginServiceDeps) *loginService {
	return &loginService{
		mencache:      params.Cache,
		logger:        params.Logger,
		hash:          params.Hash,
		user:          params.UserRepository,
		refreshToken:  params.RefreshToken,
		token:         params.Token,
		tokenService:  params.TokenService,
		observability: params.Observability,
	}
}

func (s *loginService) Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, error) {
	const method = "Login"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	// Check if account is locked
	locked, err := s.mencache.IsAccountLocked(ctx, request.Email)
	if err != nil {
		s.logger.Error("Failed to check account lock status", zap.Error(err), zap.String("email", request.Email))
	}
	if locked {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, user_errors.ErrAccountLocked, method, span, zap.String("email", request.Email))
	}

	res, err := s.user.FindByEmailAndVerify(ctx, request.Email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.String("email", request.Email))
	}

	err = s.hash.ComparePassword(res.Password, request.Password)
	if err != nil {
		status = "error"

		// Increment failed login attempts
		_, incErr := s.mencache.IncrementFailedLogin(ctx, request.Email)
		if incErr != nil {
			s.logger.Error("Failed to increment failed login counter", zap.Error(incErr), zap.String("email", request.Email))
		}

		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, user_errors.ErrFailedPasswordNoMatch, method, span, zap.String("email", request.Email))
	}

	token, err := s.tokenService.createAccessToken(ctx, int(res.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, int(res.UserID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*response.TokenResponse](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
	}

	// Reset failed login attempts on successful login
	if err := s.mencache.ResetFailedLogin(ctx, request.Email); err != nil {
		s.logger.Error("Failed to reset failed login data", zap.Error(err), zap.String("email", request.Email))
	}

	tokenResp := &response.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	logSuccess("Successfully logged in", zap.String("email", request.Email))

	return tokenResp, nil
}
