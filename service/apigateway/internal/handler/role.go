package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/role_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleHandleApi struct {
	role    pb.RoleServiceClient
	logger  logger.LoggerInterface
	mapping response_api.RoleResponseMapper
}

func NewHandlerRole(router *echo.Echo, role pb.RoleServiceClient, logger logger.LoggerInterface, mapping response_api.RoleResponseMapper) *roleHandleApi {
	roleHandler := &roleHandleApi{
		role:    role,
		logger:  logger,
		mapping: mapping,
	}

	routerRole := router.Group("/api/role")

	routerRole.GET("", roleHandler.FindAll)
	routerRole.GET("/:id", roleHandler.FindById)
	routerRole.GET("/active", roleHandler.FindByActive)
	routerRole.GET("/trashed", roleHandler.FindByTrashed)
	routerRole.GET("/user/:user_id", roleHandler.FindByUserId)
	routerRole.POST("", roleHandler.Create)
	routerRole.POST("/update/:id", roleHandler.Update)
	routerRole.POST("/trashed/:id", roleHandler.Trashed)
	routerRole.POST("/restore/:id", roleHandler.Restore)
	routerRole.DELETE("/permanent/:id", roleHandler.DeletePermanent)
	routerRole.POST("/restore/all", roleHandler.RestoreAll)
	routerRole.DELETE("/permanent-all", roleHandler.DeleteAllPermanent)

	return roleHandler
}

// FindAll godoc.
// @Summary Get all roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRole "List of roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch roles"
// @Router /api/role [get]
func (h *roleHandleApi) FindAll(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindAllRole(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to fetch role records", zap.Error(err))
		return role_errors.ErrApiFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationRole(res)

	return c.JSON(http.StatusOK, so)
}

// FindById godoc.
// @Summary Get a role by ID
// @Tags Role
// @Security Bearer
// @Description Retrieve a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch role"
// @Router /api/role/{id} [get]
func (h *roleHandleApi) FindById(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return role_errors.ErrApiRoleInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.FindByIdRole(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to fetch role", zap.Error(err))
		return role_errors.ErrApiRoleNotFound(c)
	}

	so := h.mapping.ToApiResponseRole(res)

	return c.JSON(http.StatusOK, so)
}

// FindByActive godoc.
// @Summary Get active roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of active roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of active roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch active roles"
// @Router /api/role/active [get]
func (h *roleHandleApi) FindByActive(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindByActive(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to fetch active roles", zap.Error(err))
		return role_errors.ErrApiFailedFindActive(c)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// FindByTrashed godoc.
// @Summary Get trashed roles
// @Tags Role
// @Security Bearer
// @Description Retrieve a paginated list of trashed roles with optional search and pagination parameters.
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10)"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of trashed roles"
// @Failure 400 {object} response.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch trashed roles"
// @Router /api/role/trashed [get]
func (h *roleHandleApi) FindByTrashed(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllRoleRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.role.FindByTrashed(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to fetch trashed roles", zap.Error(err))
		return role_errors.ErrApiFailedFindTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationRoleDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// FindByUserId godoc.
// @Summary Get role by user ID
// @Tags Role
// @Security Bearer
// @Description Retrieve a role by the associated user ID.
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to fetch role by user ID"
// @Router /api/role/user/{user_id} [get]
func (h *roleHandleApi) FindByUserId(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 {
		return role_errors.ErrApiRoleInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdUserRoleRequest{
		UserId: int32(userID),
	}

	res, err := h.role.FindByUserId(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to fetch role by user ID", zap.Error(err))
		return role_errors.ErrApiRoleNotFound(c)
	}

	so := h.mapping.ToApiResponsesRole(res)

	return c.JSON(http.StatusOK, so)
}

// Create godoc.
// @Summary Create a new role
// @Tags Role
// @Security Bearer
// @Description Create a new role with the provided details.
// @Accept json
// @Produce json
// @Param request body requests.CreateRoleRequest true "Role data"
// @Success 200 {object} response.ApiResponseRole "Created role data"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to create role"
// @Router /api/role/create [post]
func (h *roleHandleApi) Create(c echo.Context) error {
	var req requests.CreateRoleRequest

	if err := c.Bind(&req); err != nil {
		return role_errors.ErrApiBindCreateRole(c)
	}

	if err := req.Validate(); err != nil {
		return role_errors.ErrApiValidateCreateRole(c)
	}

	reqPb := &pb.CreateRoleRequest{
		Name: req.Name,
	}

	ctx := c.Request().Context()

	res, err := h.role.CreateRole(ctx, reqPb)
	if err != nil {
		h.logger.Debug("Failed to create role", zap.Error(err))
		return role_errors.ErrApiFailedCreateRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

	return c.JSON(http.StatusOK, so)
}

// Update godoc.
// @Summary Update a role
// @Tags Role
// @Security Bearer
// @Description Update an existing role with the provided details.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body requests.UpdateRoleRequest true "Role data"
// @Success 200 {object} response.ApiResponseRole "Updated role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID or request body"
// @Failure 500 {object} response.ErrorResponse "Failed to update role"
// @Router /api/role/update/{id} [post]
func (h *roleHandleApi) Update(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return role_errors.ErrInvalidRoleId(c)
	}

	var req requests.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return role_errors.ErrApiBindUpdateRole(c)
	}

	if err := req.Validate(); err != nil {
		return role_errors.ErrApiValidateUpdateRole(c)
	}

	reqPb := &pb.UpdateRoleRequest{
		Id:   int32(roleID),
		Name: req.Name,
	}

	ctx := c.Request().Context()

	res, err := h.role.UpdateRole(ctx, reqPb)
	if err != nil {
		h.logger.Debug("Failed to update role", zap.Error(err))
		return role_errors.ErrApiFailedUpdateRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

	return c.JSON(http.StatusOK, so)
}

// Trashed godoc.
// @Summary Soft-delete a role
// @Tags Role
// @Security Bearer
// @Description Soft-delete a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Soft-deleted role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to soft-delete role"
// @Router /api/role/trashed/{id} [post]
func (h *roleHandleApi) Trashed(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return role_errors.ErrInvalidRoleId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.TrashedRole(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to trash role", zap.Error(err))
		return role_errors.ErrApiFailedTrashedRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

	return c.JSON(http.StatusOK, so)
}

// Restore godoc.
// @Summary Restore a soft-deleted role
// @Tags Role
// @Security Bearer
// @Description Restore a soft-deleted role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Restored role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore role"
// @Router /api/role/restore/{id} [post]
func (h *roleHandleApi) Restore(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return role_errors.ErrInvalidRoleId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.RestoreRole(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to restore role", zap.Error(err))
		return role_errors.ErrApiFailedRestoreRole(c)
	}

	so := h.mapping.ToApiResponseRole(res)

	return c.JSON(http.StatusOK, so)
}

// DeletePermanent godoc.
// @Summary Permanently delete a role
// @Tags Role
// @Security Bearer
// @Description Permanently delete a role by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Permanently deleted role data"
// @Failure 400 {object} response.ErrorResponse "Invalid role ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete role permanently"
// @Router /api/role/permanent/{id} [delete]
func (h *roleHandleApi) DeletePermanent(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 {
		return role_errors.ErrInvalidRoleId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdRoleRequest{
		RoleId: int32(roleID),
	}

	res, err := h.role.DeleteRolePermanent(ctx, req)
	if err != nil {
		h.logger.Debug("Failed to delete role permanently", zap.Error(err))
		return role_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseRoleDelete(res)

	return c.JSON(http.StatusOK, so)
}

// RestoreAll godoc.
// @Summary Restore all soft-deleted roles
// @Tags Role
// @Security Bearer
// @Description Restore all soft-deleted roles.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Restored roles data"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all roles"
// @Router /api/role/restore/all [post]
func (h *roleHandleApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.role.RestoreAllRole(ctx, &emptypb.Empty{})
	if err != nil {
		h.logger.Debug("Failed to restore all roles", zap.Error(err))
		return role_errors.ErrApiFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseRoleAll(res)

	return c.JSON(http.StatusOK, so)
}

// DeleteAllPermanent godoc.
// @Summary Permanently delete all roles
// @Tags Role
// @Security Bearer
// @Description Permanently delete all roles.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Permanently deleted roles data"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all roles permanently"
// @Router /api/role/permanent/all [delete]
func (h *roleHandleApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.role.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	if err != nil {
		h.logger.Debug("Failed to delete all roles permanently", zap.Error(err))
		return role_errors.ErrApiFailedDeleteAll(c)
	}

	so := h.mapping.ToApiResponseRoleAll(res)

	return c.JSON(http.StatusOK, so)
}
