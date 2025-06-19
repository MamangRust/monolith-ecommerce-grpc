package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userHandleApi struct {
	client          pb.UserServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.UserResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerUser(router *echo.Echo, client pb.UserServiceClient, logger logger.LoggerInterface, mapping response_api.UserResponseMapper) *userHandleApi {
	userHandler := &userHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}
	routerUser := router.Group("/api/user")

	routerUser.GET("", userHandler.FindAllUser)
	routerUser.GET("/:id", userHandler.FindById)
	routerUser.GET("/active", userHandler.FindByActive)
	routerUser.GET("/trashed", userHandler.FindByTrashed)

	routerUser.POST("/create", userHandler.Create)
	routerUser.POST("/update/:id", userHandler.Update)

	routerUser.POST("/trashed/:id", userHandler.TrashedUser)
	routerUser.POST("/restore/:id", userHandler.RestoreUser)
	routerUser.DELETE("/permanent/:id", userHandler.DeleteUserPermanent)

	routerUser.POST("/restore/all", userHandler.RestoreAllUser)
	routerUser.POST("/permanent/all", userHandler.DeleteAllUserPermanent)

	return userHandler
}

// @Security Bearer
// @Summary Find all users
// @Tags User
// @Description Retrieve a list of all users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUser "List of users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user [get]
func (h *userHandleApi) FindAllUser(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllUser"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationUser(res)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find user by ID
// @Tags User
// @Description Retrieve a user by ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUser "User data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user/{id} [get]
func (h *userHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid user ID", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiUserNotFound(c)
	}

	so := h.mapping.ToApiResponseUser(user)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active users
// @Tags User
// @Description Retrieve a list of active users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUserDeleteAt "List of active users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user/active [get]
func (h *userHandleApi) FindByActive(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActive"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiFailedFindActive(c)
	}

	so := h.mapping.ToApiResponsePaginationUserDeleteAt(res)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed user records.
// @Summary Retrieve trashed users
// @Tags User
// @Description Retrieve a list of trashed user records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUserDeleteAt "List of trashed user data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user/trashed [get]
func (h *userHandleApi) FindByTrashed(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashed"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve user data", err, zap.Error(err))

		return user_errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationUserDeleteAt(res)

	logSuccess("Successfully retrieve user data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new user.
// @Summary Create a new user
// @Tags User
// @Description Create a new user with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateUserRequest true "Create user request"
// @Success 200 {object} response.ApiResponseUser "Successfully created user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create user"
// @Router /api/user/create [post]
func (h *userHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	var body requests.CreateUserRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to create user", err, zap.Error(err))

		return user_errors.ErrApiBindCreateUser(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to create user", err, zap.Error(err))

		return user_errors.ErrApiValidateCreateUser(c)
	}

	req := &pb.CreateUserRequest{
		Firstname:       body.FirstName,
		Lastname:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to create user", err, zap.Error(err))

		return user_errors.ErrApiFailedCreateUser(c)
	}

	so := h.mapping.ToApiResponseUser(res)

	logSuccess("Successfully create user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing user record.
// @Summary Update an existing user
// @Tags User
// @Description Update an existing user record with the provided details
// @Accept json
// @Produce json
// @Param UpdateUserRequest body requests.UpdateUserRequest true "Update user request"
// @Success 200 {object} response.ApiResponseUser "Successfully updated user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update user"
// @Router /api/user/update/{id} [post]
func (h *userHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"

		logError("Failed to update user", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	var body requests.UpdateUserRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to update user", err, zap.Error(err))

		return user_errors.ErrApiBindUpdateUser(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to update user", err, zap.Error(err))

		return user_errors.ErrApiValidateUpdateUser(c)
	}

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
		status = "error"

		logError("Failed to update user", err, zap.Error(err))

		return user_errors.ErrApiFailedUpdateUser(c)
	}

	so := h.mapping.ToApiResponseUser(res)

	logSuccess("Successfully update user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedUser retrieves a trashed user record by its ID.
// @Summary Retrieve a trashed user
// @Tags User
// @Description Retrieve a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDeleteAt "Successfully retrieved trashed user"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed user"
// @Router /api/user/trashed/{id} [get]
func (h *userHandleApi) TrashedUser(c echo.Context) error {
	const method = "TrashedUser"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to retrieve trashed user", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.TrashedUser(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve trashed user", err, zap.Error(err))

		return user_errors.ErrApiFailedTrashedUser(c)
	}

	so := h.mapping.ToApiResponseUserDeleteAt(user)

	logSuccess("Successfully retrieve trashed user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreUser restores a user record from the trash by its ID.
// @Summary Restore a trashed user
// @Tags User
// @Description Restore a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDeleteAt "Successfully restored user"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore user"
// @Router /api/user/restore/{id} [post]
func (h *userHandleApi) RestoreUser(c echo.Context) error {
	const method = "RestoreUser"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to restore user", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.RestoreUser(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore user", err, zap.Error(err))

		return user_errors.ErrApiFailedRestoreUser(c)
	}

	so := h.mapping.ToApiResponseUserDeleteAt(user)

	logSuccess("Successfully restore user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteUserPermanent permanently deletes a user record by its ID.
// @Summary Permanently delete a user
// @Tags User
// @Description Permanently delete a user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete user:"
// @Router /api/user/delete/{id} [delete]
func (h *userHandleApi) DeleteUserPermanent(c echo.Context) error {
	const method = "DeleteUserPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to delete user", err, zap.Error(err))

		return user_errors.ErrApiUserInvalidId(c)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.DeleteUserPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to delete user", err, zap.Error(err))

		return user_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseUserDelete(user)

	logSuccess("Successfully delete user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreUser restores a user record from the trash by its ID.
// @Summary Restore a trashed user
// @Tags User
// @Description Restore a trashed user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserAll "Successfully restored user all"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore user"
// @Router /api/user/restore/all [post]
func (h *userHandleApi) RestoreAllUser(c echo.Context) error {
	const method = "RestoreAllUser"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllUser(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all user", err, zap.Error(err))

		return user_errors.ErrApiFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseUserAll(res)

	logSuccess("Successfully restore all user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteUserPermanent permanently deletes a user record by its ID.
// @Summary Permanently delete a user
// @Tags User
// @Description Permanently delete a user record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUserDelete "Successfully deleted user record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete user:"
// @Router /api/user/delete/all [post]
func (h *userHandleApi) DeleteAllUserPermanent(c echo.Context) error {
	const method = "DeleteAllUserPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllUserPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to permanently delete all user", err, zap.Error(err))

		return user_errors.ErrApiFailedDeleteAll(c)
	}

	so := h.mapping.ToApiResponseUserAll(res)

	logSuccess("Successfully permanently delete all user", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *userHandleApi) startTracingAndLogging(
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

func (s *userHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
