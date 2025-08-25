package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-auth/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/auth"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loginService struct {
	errorPassword   errorhandler.PasswordErrorHandler
	errorToken      errorhandler.TokenErrorHandler
	errorHandler    errorhandler.LoginErrorHandler
	mencache        mencache.LoginCache
	logger          logger.LoggerInterface
	hash            hash.HashPassword
	user            repository.UserRepository
	refreshToken    repository.RefreshTokenRepository
	token           auth.TokenManager
	trace           trace.Tracer
	tokenService    tokenService
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewLoginService(
	errorPassword errorhandler.PasswordErrorHandler,
	errorToken errorhandler.TokenErrorHandler,
	errorHandler errorhandler.LoginErrorHandler,
	mencache mencache.LoginCache,
	logger logger.LoggerInterface,
	hash hash.HashPassword,
	userRepository repository.UserRepository,
	refreshToken repository.RefreshTokenRepository,
	token auth.TokenManager,
	tokenService tokenService,
) *loginService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_service_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "login_service_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &loginService{
		errorPassword:   errorPassword,
		errorToken:      errorToken,
		errorHandler:    errorHandler,
		mencache:        mencache,
		logger:          logger,
		hash:            hash,
		user:            userRepository,
		refreshToken:    refreshToken,
		token:           token,
		trace:           otel.Tracer("login-service"),
		tokenService:    tokenService,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *loginService) Login(ctx context.Context, request *requests.AuthRequest) (*response.TokenResponse, *response.ErrorResponse) {
	const method = "Login"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	if cachedToken, found := s.mencache.GetCachedLogin(ctx, request.Email); found {
		logSuccess("Successfully logged in", zap.String("email", request.Email))
		return cachedToken, nil
	}

	res, err := s.user.FindByEmailAndVerify(ctx, request.Email)
	if err != nil {
		return s.errorHandler.HandleFindEmailError(err, method, "LOGIN_ERR", span, &status, zap.Error(err))
	}

	err = s.hash.ComparePassword(res.Password, request.Password)
	if err != nil {
		return s.errorPassword.HandleComparePasswordError(err, method, "COMPARE_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	token, err := s.tokenService.createAccessToken(ctx, res.ID)
	if err != nil {
		return s.errorToken.HandleCreateAccessTokenError(err, method, "CREATE_ACCESS_TOKEN_ERR", span, &status, zap.Error(err))
	}

	refreshToken, err := s.tokenService.createRefreshToken(ctx, res.ID)
	if err != nil {
		return s.errorToken.HandleCreateRefreshTokenError(err, method, "CREATE_REFRESH_TOKEN_ERR", span, &status, zap.Error(err))
	}

	tokenResp := &response.TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}

	s.mencache.SetCachedLogin(ctx, request.Email, tokenResp, time.Minute)

	logSuccess("Successfully logged in", zap.String("email", request.Email))

	return tokenResp, nil
}

func (s *loginService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *loginService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
