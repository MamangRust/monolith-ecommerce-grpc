package handler

import (
	"net/http"
	"strings"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/auth_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandleApi struct {
	client  pb.AuthServiceClient
	logger  logger.LoggerInterface
	mapping response_api.AuthResponseMapper
}

func NewHandlerAuth(router *echo.Echo, client pb.AuthServiceClient, logger logger.LoggerInterface, mapper response_api.AuthResponseMapper) *authHandleApi {
	authHandler := &authHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapper,
	}
	routerAuth := router.Group("/api/auth")

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
// @Router /api/auth/hello [get]
func (h *authHandleApi) HandleHello(c echo.Context) error {
	return c.String(200, "Hello")
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
	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindRegister(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateRegister(c)
	}

	data := &pb.RegisterRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.RegisterUser(c.Request().Context(), data)

	if err != nil {
		h.logger.Error("Registration failed", zap.Error(err))
		return auth_errors.ErrApiRegister(c)
	}

	return c.JSON(http.StatusCreated, h.mapping.ToResponseRegister(res))
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
	var body requests.AuthRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindLogin(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateRegister(c)
	}

	res, err := h.client.LoginUser(c.Request().Context(), &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			h.logger.Debug("Invalid login attempt", zap.String("email", body.Email))
			return auth_errors.ErrInvalidLogin(c)
		}

		h.logger.Error("Login failed", zap.Error(err))

		if status.Code(err) == codes.Internal && strings.Contains(err.Error(), "empty token") {
			return auth_errors.ErrInvalidAccessToken(c)
		}

		return auth_errors.ErrApiLogin(c)
	}

	mappedResponse := h.mapping.ToResponseLogin(res)

	if mappedResponse.Data == nil || mappedResponse.Data.AccessToken == "" {
		h.logger.Error("Empty token in final response", zap.Any("response", mappedResponse))
		return auth_errors.ErrApiLogin(c)
	}

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
	var body requests.RefreshTokenRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return auth_errors.ErrBindRefreshToken(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return auth_errors.ErrValidateRefreshToken(c)
	}

	res, err := h.client.RefreshToken(c.Request().Context(), &pb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		h.logger.Error("Token refresh failed", zap.Error(err))
		return auth_errors.ErrApiRefreshToken(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToResponseRefreshToken(res))
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
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return auth_errors.ErrInvalidAccessToken(c)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	res, err := h.client.GetMe(c.Request().Context(), &pb.GetMeRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		h.logger.Error("Failed to get user information", zap.Error(err))
		return auth_errors.ErrApiGetMe(c)
	}

	so := h.mapping.ToResponseGetMe(res)

	return c.JSON(http.StatusOK, so)
}
