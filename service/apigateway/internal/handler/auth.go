package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	authapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/auth"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"github.com/go-playground/validator/v10"
	"fmt"
	auth_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/cache/auth"
)


type authHandleApi struct {
	client        pb.AuthServiceClient
	logger        logger.LoggerInterface
	queryMapper   authapimapper.AuthQueryResponseMapper
	commandMapper authapimapper.AuthCommandResponseMapper
	apiHandler    sharedErrors.ApiHandler
	cache         auth_cache.AuthMencache
}

type authHandleParams struct {
	client pb.AuthServiceClient
	router *echo.Echo
	cache  auth_cache.AuthMencache
	logger logger.LoggerInterface
	mapper authapimapper.AuthResponseMapper
	apiHandler sharedErrors.ApiHandler
}

func NewHandlerAuth(params *authHandleParams) *authHandleApi {
	authHandler := &authHandleApi{
		client:        params.client,
		logger:        params.logger,
		queryMapper:   params.mapper.QueryMapper(),
		commandMapper: params.mapper.CommandMapper(),
		apiHandler:    params.apiHandler,
		cache:         params.cache,
	}
	routerAuth := params.router.Group("/api/auth")

	routerAuth.GET("/hello", authHandler.HandleHello)
	routerAuth.POST("/register", params.apiHandler.Handle("register", authHandler.Register))
	routerAuth.POST("/login", params.apiHandler.Handle("login", authHandler.Login))
	routerAuth.POST("/refresh-token", params.apiHandler.Handle("refresh-token", authHandler.RefreshToken))
	routerAuth.GET("/me", params.apiHandler.Handle("GetMe", authHandler.GetMe))

	return authHandler
}

func (h *authHandleApi) HandleHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello")
}

// @Summary Register a new user
// @Tags Auth
// @Description Create a new user account
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "Registration details"
// @Success 201 {object} response.ApiResponseRegister "Successfully registered"
// @Failure 400 {object} errors.ErrorResponse "Validation error"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/auth/register [post]
func (h *authHandleApi) Register(c echo.Context) error {
	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		return sharedErrors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return sharedErrors.NewValidationError(validations)
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
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusCreated, h.commandMapper.ToResponseRegister(res))
}


// @Summary Login user
// @Tags Auth
// @Description Authenticate user and return tokens
// @Accept json
// @Produce json
// @Param request body requests.AuthRequest true "Login credentials"
// @Success 200 {object} response.ApiResponseLogin "Successfully logged in"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/auth/login [post]
func (h *authHandleApi) Login(c echo.Context) error {
	var body requests.AuthRequest

	if err := c.Bind(&body); err != nil {
		return sharedErrors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return sharedErrors.NewValidationError(validations)
	}

	ctx := c.Request().Context()
	cachedResponse, found := h.cache.GetCachedLogin(ctx, body.Email)
	if found {
		h.logger.Debug("Returning login response from cache", zap.String("email", body.Email))
		return c.JSON(http.StatusOK, cachedResponse)
	}

	res, err := h.client.LoginUser(ctx, &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}


	mappedResponse := h.commandMapper.ToResponseLogin(res)
	h.cache.SetCachedLogin(ctx, body.Email, mappedResponse)

	return c.JSON(http.StatusOK, mappedResponse)
}

// @Summary Refresh token
// @Tags Auth
// @Description refresh token
// @Accept json
// @Produce json
// @Param request body requests.RefreshTokenRequest true "Refresh token details"
// @Success 200 {object} response.ApiResponseRefreshToken "Successfully refreshed token"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/auth/refresh-token [post]
func (h *authHandleApi) RefreshToken(c echo.Context) error {
	var body requests.RefreshTokenRequest

	if err := c.Bind(&body); err != nil {
		return sharedErrors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return sharedErrors.NewValidationError(validations)
	}

	ctx := c.Request().Context()
	cachedResponse, found := h.cache.GetRefreshToken(ctx, body.RefreshToken)
	if found {
		h.logger.Debug("Returning refresh token response from cache")
		return c.JSON(http.StatusOK, cachedResponse)
	}

	res, err := h.client.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}


	mappedResponse := h.commandMapper.ToResponseRefreshToken(res)
	h.cache.SetRefreshToken(ctx, body.RefreshToken, mappedResponse)

	return c.JSON(http.StatusOK, mappedResponse)
}

// @Security Bearer
// @Summary Get current user info
// @Tags Auth
// @Description Retrieve current authenticated user details
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseGetMe "User info"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal server error"
// @Router /api/auth/me [get]
func (h *authHandleApi) GetMe(c echo.Context) error {
	userIdStr, ok := c.Get("userId").(string)
	if !ok {
		return sharedErrors.NewBadRequestError("user not authenticated")
	}

	uid, err := strconv.ParseInt(userIdStr, 10, 32)
	if err != nil {
		return sharedErrors.NewBadRequestError("invalid user ID format")
	}
	userID := int(uid)

	ctx := c.Request().Context()
	if cached, found := h.cache.GetCachedUserInfo(ctx, userIdStr); found {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.GetMe(ctx, &pb.GetMeRequest{UserId: int32(userID)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}


	response := h.queryMapper.ToResponseGetMe(res)
	h.cache.SetCachedUserInfo(ctx, userIdStr, response)

	return c.JSON(http.StatusOK, response)
}



func (h *authHandleApi) parseValidationErrors(err error) []sharedErrors.ValidationError {
	var validationErrs []sharedErrors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, sharedErrors.ValidationError{
				Field:   fe.Field(),
				Message: h.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []sharedErrors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (h *authHandleApi) getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Validation failed on '%s' tag", fe.Tag())
	}
}
