package rolehandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-apigateway/middlewares"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache"
	role_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/role"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/role"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type roleCommandHandlerApi struct {
	kafka      *kafka.Kafka
	role       pb.RoleCommandServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.RoleCommandResponseMapper
	cache      role_cache.RoleCommandCache
	apiHandler errors.ApiHandler
}

type roleCommandHandleDeps struct {
	client     pb.RoleCommandServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.RoleCommandResponseMapper
	kafka      *kafka.Kafka
	cache_role mencache.RoleCache
	cache      role_cache.RoleCommandCache
	apiHandler errors.ApiHandler
}

func NewRoleCommandHandleApi(params *roleCommandHandleDeps) *roleCommandHandlerApi {
	handler := &roleCommandHandlerApi{
		role:       params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
		kafka:      params.kafka,
	}

	roleMiddleware := middlewares.NewRoleValidator(params.kafka, "request-role", "response-role", 5*time.Second, params.logger, params.cache_role)
	routerRole := params.router.Group("/api/role-command")
	roleMiddlewareChain := roleMiddleware.Middleware()
	requireAdmin := middlewares.RequireRoles("Admin_Admin_14")

	routerRole.POST("/create", roleMiddlewareChain(requireAdmin(handler.Create)))
	routerRole.POST("/update/:id", roleMiddlewareChain(requireAdmin(handler.Update)))
	routerRole.POST("/trashed/:id", roleMiddlewareChain(requireAdmin(handler.Trash)))
	routerRole.POST("/restore/:id", roleMiddlewareChain(requireAdmin(handler.Restore)))
	routerRole.DELETE("/permanent/:id", roleMiddlewareChain(requireAdmin(handler.DeletePermanent)))
	routerRole.POST("/restore/all", roleMiddlewareChain(requireAdmin(handler.RestoreAll)))
	routerRole.POST("/permanent/all", roleMiddlewareChain(requireAdmin(handler.DeleteAllPermanent)))

	return handler
}

// @Security Bearer
// @Summary Create a new role
// @Tags Role Command
// @Description Create a new role with the provided details
// @Accept json
// @Produce json
// @Param body body requests.CreateRoleRequest true "Create role request"
// @Success 200 {object} response.ApiResponseRole "Successfully created role"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create role"
// @Router /api/role-command/create [post]
func (h *roleCommandHandlerApi) Create(c echo.Context) error {
	var body requests.CreateRoleRequest
	if err := c.Bind(&body); err != nil { return errors.NewBadRequestError("Invalid request").WithInternal(err) }

	ctx := c.Request().Context()
	res, err := h.role.CreateRole(ctx, &pb.CreateRoleRequest{Name: body.Name})
	if err != nil { return errors.ParseGrpcError(err) }


	return c.JSON(http.StatusOK, h.mapper.ToApiResponseRole(res))
}

// @Security Bearer
// @Summary Update role details
// @Tags Role Command
// @Description Update the name of an existing role
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param body body requests.UpdateRoleRequest true "Update role request"
// @Success 200 {object} response.ApiResponseRole "Successfully updated role"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update role"
// @Router /api/role-command/update/{id} [post]
func (h *roleCommandHandlerApi) Update(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 { return errors.NewBadRequestError("id is required") }

	var body requests.UpdateRoleRequest
	if err := c.Bind(&body); err != nil { return errors.NewBadRequestError("Invalid request").WithInternal(err) }

	ctx := c.Request().Context()
	res, err := h.role.UpdateRole(ctx, &pb.UpdateRoleRequest{Id: int32(roleID), Name: body.Name})
	if err != nil { return errors.ParseGrpcError(err) }


	h.cache.DeleteCachedRole(ctx, roleID)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseRole(res))
}

// @Security Bearer
// @Summary Move role to trash
// @Tags Role Command
// @Description Move a role record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Successfully moved role to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid role ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move role to trash"
// @Router /api/role-command/trashed/{id} [post]
func (h *roleCommandHandlerApi) Trash(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.role.TrashedRole(ctx, &pb.FindByIdRoleRequest{RoleId: int32(roleID)})
	if err != nil { return errors.ParseGrpcError(err) }

	h.cache.DeleteCachedRole(ctx, roleID)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseRole(res))
}

// @Security Bearer
// @Summary Restore a trashed role
// @Tags Role Command
// @Description Restore a trashed role record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Successfully restored role"
// @Failure 400 {object} errors.ErrorResponse "Invalid role ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore role"
// @Router /api/role-command/restore/{id} [put]
func (h *roleCommandHandlerApi) Restore(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.role.RestoreRole(ctx, &pb.FindByIdRoleRequest{RoleId: int32(roleID)})
	if err != nil { return errors.ParseGrpcError(err) }


	h.cache.DeleteCachedRole(ctx, roleID)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseRole(res))
}

// @Security Bearer
// @Summary Permanently delete a role
// @Tags Role Command
// @Description Permanently delete a role record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRoleDelete "Successfully deleted role record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid role ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete role permanently"
// @Router /api/role-command/permanent/{id} [delete]
func (h *roleCommandHandlerApi) DeletePermanent(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.role.DeleteRolePermanent(ctx, &pb.FindByIdRoleRequest{RoleId: int32(roleID)})
	if err != nil { return errors.ParseGrpcError(err) }


	h.cache.DeleteCachedRole(ctx, roleID)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseRoleDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed roles
// @Tags Role Command
// @Description Restore all trashed role records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Successfully restored all roles"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore roles"
// @Router /api/role-command/restore/all [post]
func (h *roleCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.role.RestoreAllRole(ctx, &emptypb.Empty{})
	if err != nil { return errors.ParseGrpcError(err) }


	return c.JSON(http.StatusOK, h.mapper.ToApiResponseRoleAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed roles
// @Tags Role Command
// @Description Permanently delete all trashed role records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseRoleAll "Successfully deleted all roles permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete roles permanently"
// @Router /api/role-command/permanent/all [delete]
func (h *roleCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.role.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	if err != nil { return errors.ParseGrpcError(err) }


	return c.JSON(http.StatusOK, h.mapper.ToApiResponseRoleAll(res))
}


