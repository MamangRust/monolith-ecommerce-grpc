package userhandler

import (
	"fmt"
	"net/http"
	"strconv"

	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/user"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/user"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)


type userCommandHandleApi struct {
	client pb.UserCommandServiceClient

	logger logger.LoggerInterface

	mapper apimapper.UserCommandResponseMapper

	cache user_cache.UserMencache

	apiHandler errors.ApiHandler
}

type userCommandHandleDeps struct {
	client pb.UserCommandServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.UserCommandResponseMapper

	cache user_cache.UserMencache

	apiHandler errors.ApiHandler
}

func NewUserCommandHandleApi(params *userCommandHandleDeps) *userCommandHandleApi {

	userCommandHandleApi := &userCommandHandleApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	routerUser := params.router.Group("/api/user-command")

	routerUser.POST("/create", params.apiHandler.Handle("create-user", userCommandHandleApi.Create))
	routerUser.POST("/update/:id", params.apiHandler.Handle("update-user", userCommandHandleApi.Update))

	routerUser.POST("/trashed/:id", params.apiHandler.Handle("trash-user", userCommandHandleApi.TrashedUser))
	routerUser.POST("/restore/:id", params.apiHandler.Handle("restore-user", userCommandHandleApi.RestoreUser))
	routerUser.DELETE("/permanent/:id", params.apiHandler.Handle("delete-user-permanent", userCommandHandleApi.DeleteUserPermanent))

	routerUser.POST("/restore/all", params.apiHandler.Handle("restore-all-users", userCommandHandleApi.RestoreAllUser))
	routerUser.POST("/permanent/all", params.apiHandler.Handle("delete-all-users-permanent", userCommandHandleApi.DeleteAllUserPermanent))

	return userCommandHandleApi
}

// @Security Bearer
// Create handles the creation of a new user.
// @Summary Create a new user
// @Tags User Command
// @Description Create a new user with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "Create user request"
// @Success 200 {object} response.ApiResponseUser "Successfully created user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create user"
// @Router /api/user-command/create [post]
func (h *userCommandHandleApi) Create(c echo.Context) error {
	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.CreateUserRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	so := h.mapper.ToApiResponseUser(res)

	h.cache.SetCachedUserCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing user record.
// @Summary Update an existing user
// @Tags User Command
// @Description Update an existing user record with the provided details
// @Accept json
// @Produce json
// @Param UpdateUserRequest body requests.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.ApiResponseUser "Successfully updated user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update user"
// @Router /api/user-command/update/{id} [post]
func (h *userCommandHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	var body requests.UpdateUserRequest

	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("Invalid request format").WithInternal(err)
	}

	if err := body.Validate(); err != nil {
		validations := h.parseValidationErrors(err)
		return errors.NewValidationError(validations)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateUserRequest{
		Id:              int32(idInt),
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	so := h.mapper.ToApiResponseUser(res)

	h.cache.DeleteUserCache(ctx, idInt)
	h.cache.SetCachedUserCache(ctx, so)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve a trashed user
// @Tags User Command
// @Description Retrieve a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDeleteAt "Successfully retrieved trashed user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed user"
// @Router /api/user-command/trashed/{id} [get]
func (h *userCommandHandleApi) TrashedUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.TrashedUser(ctx, req)

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	so := h.mapper.ToApiResponseUserDeleteAt(user)

	h.cache.DeleteUserCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore a trashed user
// @Tags User Command
// @Description Restore a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDeleteAt "Successfully restored user"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore user"
// @Router /api/user-command/restore/{id} [post]
func (h *userCommandHandleApi) RestoreUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.RestoreUser(ctx, req)

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	so := h.mapper.ToApiResponseUserDeleteAt(user)

	h.cache.DeleteUserCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Permanently delete a user
// @Tags User Command
// @Description Permanently delete a user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete user"
// @Router /api/user-command/permanent/{id} [delete]
func (h *userCommandHandleApi) DeleteUserPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.DeleteUserPermanent(ctx, req)

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	so := h.mapper.ToApiResponseUserDelete(user)

	h.cache.DeleteUserCache(ctx, id)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Restore all trashed users
// @Tags User Command
// @Description Restore all trashed user records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseUserAll "Successfully restored all users"
// @Failure 500 {object} response.ErrorResponse "Failed to restore users"
// @Router /api/user-command/restore/all [post]
func (h *userCommandHandleApi) RestoreAllUser(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllUser(ctx, &emptypb.Empty{})

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	so := h.mapper.ToApiResponseUserAll(res)

	h.logger.Debug("Successfully restored all user")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Permanently delete all trashed users
// @Tags User Command
// @Description Permanently delete all trashed user records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseUserAll "Successfully deleted all users permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete users"
// @Router /api/user-command/permanent/all [post]
func (h *userCommandHandleApi) DeleteAllUserPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllUserPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	so := h.mapper.ToApiResponseUserAll(res)

	h.logger.Debug("Successfully deleted all user permanently")

	return c.JSON(http.StatusOK, so)
}



func (h *userCommandHandleApi) parseValidationErrors(err error) []errors.ValidationError {
	var validationErrs []errors.ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			validationErrs = append(validationErrs, errors.ValidationError{
				Field:   fe.Field(),
				Message: h.getValidationMessage(fe),
			})
		}
		return validationErrs
	}

	return []errors.ValidationError{
		{
			Field:   "general",
			Message: err.Error(),
		},
	}
}

func (h *userCommandHandleApi) getValidationMessage(fe validator.FieldError) string {
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
