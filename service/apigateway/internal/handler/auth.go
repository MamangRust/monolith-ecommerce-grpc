package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/auth_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type authHandleApi struct {
	client          pb.AuthServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.AuthResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerAuth(router *echo.Echo, client pb.AuthServiceClient, logger logger.LoggerInterface, mapper response_api.AuthResponseMapper) *authHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_handler_requests_total",
			Help: "Total number of auth requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_handler_request_duration_seconds",
			Help:    "Duration of auth requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	authHandler := &authHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapper,
		trace:           otel.Tracer("auth-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerAuth := router.Group("/api/auth")

	routerAuth.GET("/verify-code", authHandler.VerifyCode)
	routerAuth.POST("/forgot-password", authHandler.ForgotPassword)
	routerAuth.POST("/reset-password", authHandler.ResetPassword)
	routerAuth.GET("/hello", authHandler.HandleHello)
	routerAuth.POST("/register", authHandler.Register)
	routerAuth.POST("/login", authHandler.Login)
	routerAuth.POST("/refresh-token", authHandler.RefreshToken)
	routerAuth.GET("/me", authHandler.GetMe)

	return authHandler
}

// HandleHello godoc
// @Summary Returns a "Hello" message
// @Tags Auth
// @Description Returns a simple "Hello" message for testing purposes.
// @Produce json
// @Success 200 {string} string "Hello"
// @Router /auth/hello [get]
func (h *authHandleApi) HandleHello(c echo.Context) error {
	return c.String(200, "Hello")
}

// VerifyCode godoc
// @Summary Verifies the user using a verification code
// @Tags Auth
// @Description Verifies the user's email using the verification code provided in the query parameter.
// @Produce json
// @Param verify_code query string true "Verification Code"
// @Success 200 {object} response.ApiResponseVerifyCode
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/verify-code [get]
func (h *authHandleApi) VerifyCode(c echo.Context) error {
	const method = "VerifyCode"

	verifyCode := c.QueryParam("verify_code")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.String("verify_code", verifyCode),
	)
	status := "success"

	defer func() { end(status) }()

	res, err := h.client.VerifyCode(ctx, &pb.VerifyCodeRequest{
		Code: verifyCode,
	})

	if err != nil {
		status = "error"
		logError("Failed to verify code", err, zap.Error(err))

		return auth_errors.ErrApiVerifyCode(c)
	}

	resp := h.mapping.ToResponseVerifyCode(res)

	logSuccess("Verification success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, resp)
}

// ForgotPassword godoc
// @Summary Sends a reset token to the user's email
// @Tags Auth
// @Description Initiates password reset by sending a reset token to the provided email.
// @Accept json
// @Produce json
// @Param request body requests.ForgotPasswordRequest true "Forgot Password Request"
// @Success 200 {object} response.ApiResponseForgotPassword
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *authHandleApi) ForgotPassword(c echo.Context) error {
	const method = "ForgotPassword"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() { end(status) }()

	var body requests.ForgotPasswordRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind forgot password request", err, zap.Error(err))

		return auth_errors.ErrBindForgotPassword(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate forgot password request", err, zap.Error(err))

		return auth_errors.ErrValidateForgotPassword(c)
	}

	logSuccess("Forgot password request received", zap.String("email", body.Email))

	res, err := h.client.ForgotPassword(ctx, &pb.ForgotPasswordRequest{
		Email: body.Email,
	})
	if err != nil {
		status = "error"

		logError("Failed to forgot password", err, zap.Error(err))

		return auth_errors.ErrApiForgotPassword(c)
	}

	resp := h.mapping.ToResponseForgotPassword(res)

	logSuccess("Forgot password success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, resp)
}

// ResetPassword godoc
// @Summary Resets the user's password using a reset token
// @Tags Auth
// @Description Allows user to reset their password using a valid reset token.
// @Accept json
// @Produce json
// @Param request body requests.CreateResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} response.ApiResponseResetPassword
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/reset-password [post]
func (h *authHandleApi) ResetPassword(c echo.Context) error {
	const method = "ResetPassword"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() { end(status) }()

	var body requests.CreateResetPasswordRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind reset password request", err, zap.Error(err))

		return auth_errors.ErrBindResetPassword(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate reset password request", err, zap.Error(err))

		return auth_errors.ErrValidateResetPassword(c)
	}

	res, err := h.client.ResetPassword(ctx, &pb.ResetPasswordRequest{
		ResetToken:      body.ResetToken,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	})

	if err != nil {
		status = "error"

		logError("Failed to reset password", err, zap.Error(err))

		return auth_errors.ErrApiResetPassword(c)
	}

	so := h.mapping.ToResponseResetPassword(res)

	logSuccess("Reset password success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Description Registers a new user with the provided details.
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "User registration data"
// @Success 200 {object} response.ApiResponseRegister "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/register [post]
func (h *authHandleApi) Register(c echo.Context) error {
	const method = "Register"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() { end(status) }()

	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind register request", err, zap.Error(err))

		return auth_errors.ErrBindRegister(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate register request", err, zap.Error(err))

		return auth_errors.ErrValidateRegister(c)
	}

	data := &pb.RegisterRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.RegisterUser(ctx, data)

	if err != nil {
		status = "error"

		logError("Failed to register user", err, zap.Error(err))

		return auth_errors.ErrApiRegister(c)
	}

	so := h.mapping.ToResponseRegister(res)

	logSuccess("Register success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Login godoc
// @Summary Authenticate a user
// @Tags Auth
// @Description Authenticates a user using the provided email and password.
// @Accept json
// @Produce json
// @Param request body requests.AuthRequest true "User login credentials"
// @Success 200 {object} response.ApiResponseLogin "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/login [post]
func (h *authHandleApi) Login(c echo.Context) error {
	const method = "Login"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() { end(status) }()

	var body requests.AuthRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind login request", err, zap.Error(err))

		return auth_errors.ErrBindLogin(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate login request", err, zap.Error(err))

		return auth_errors.ErrValidateRegister(c)
	}

	data := &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	}

	res, err := h.client.LoginUser(ctx, data)

	if err != nil {
		status = "error"

		logError("Failed to login user", err, zap.Error(err))

		return auth_errors.ErrApiLogin(c)
	}

	mappedResponse := h.mapping.ToResponseLogin(res)

	logSuccess("Login success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, mappedResponse)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags Auth
// @Security Bearer
// @Description Refreshes the access token using a valid refresh token.
// @Accept json
// @Produce json
// @Param request body requests.RefreshTokenRequest true "Refresh token data"
// @Success 200 {object} response.ApiResponseRefreshToken "Success"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/refresh-token [post]
func (h *authHandleApi) RefreshToken(c echo.Context) error {
	const method = "RefreshToken"
	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() { end(status) }()

	var body requests.RefreshTokenRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind refresh token request", err, zap.Error(err))

		return auth_errors.ErrBindRefreshToken(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate refresh token request", err, zap.Error(err))

		return auth_errors.ErrValidateRefreshToken(c)
	}

	res, err := h.client.RefreshToken(c.Request().Context(), &pb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})

	if err != nil {
		status = "error"

		logError("Failed to refresh token", err, zap.Error(err))

		return auth_errors.ErrApiRefreshToken(c)
	}

	so := h.mapping.ToResponseRefreshToken(res)

	logSuccess("Refresh token success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// GetMe godoc
// @Summary Get current user information
// @Tags Auth
// @Security Bearer
// @Description Retrieves the current user's information using a valid access token from the Authorization header.
// @Produce json
// @Security BearerToken
// @Success 200 {object} response.ApiResponseGetMe "Success"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/auth/me [get]
func (h *authHandleApi) GetMe(c echo.Context) error {
	const method = "GetMe"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"

	defer func() { end(status) }()

	authHeader := c.Request().Header.Get("Authorization")

	h.logger.Debug("Authorization header: ", zap.String("authHeader", authHeader))

	if !strings.HasPrefix(authHeader, "Bearer ") {
		status = "error"

		err := errors.New("invalid authorization header")

		logError("Invalid authorization header", err, zap.String("authHeader", authHeader))

		return auth_errors.ErrInvalidAccessToken(c)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	res, err := h.client.GetMe(c.Request().Context(), &pb.GetMeRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		status = "error"

		logError("Failed to get me", err, zap.Error(err))

		return auth_errors.ErrApiGetMe(c)
	}

	so := h.mapping.ToResponseGetMe(res)

	logSuccess("Get me success", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *authHandleApi) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (func(string), func(string, ...zap.Field), func(string, error, ...zap.Field)) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError := func(msg string, err error, fields ...zap.Field) {
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *authHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
