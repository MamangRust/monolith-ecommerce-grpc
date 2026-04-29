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
)

type roleQueryHandlerApi struct {
	kafka      *kafka.Kafka
	role       pb.RoleQueryServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.RoleQueryResponseMapper
	cache      role_cache.RoleQueryCache
	apiHandler errors.ApiHandler
}

type roleQueryHandleDeps struct {
	client     pb.RoleQueryServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.RoleQueryResponseMapper
	kafka      *kafka.Kafka
	cache_role mencache.RoleCache
	cache      role_cache.RoleQueryCache
	apiHandler errors.ApiHandler
}

func NewRoleQueryHandleApi(params *roleQueryHandleDeps) *roleQueryHandlerApi {
	handler := &roleQueryHandlerApi{
		role:       params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		kafka:      params.kafka,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	roleMiddleware := middlewares.NewRoleValidator(params.kafka, "request-role", "response-role", 5*time.Second, params.logger, params.cache_role)
	routerRole := params.router.Group("/api/role-query")
	roleMiddlewareChain := roleMiddleware.Middleware()
	requireAdmin := middlewares.RequireRoles("Admin_Role_10")

	routerRole.GET("", roleMiddlewareChain(requireAdmin(handler.FindAll)))
	routerRole.GET("/:id", roleMiddlewareChain(requireAdmin(handler.FindById)))
	routerRole.GET("/active", roleMiddlewareChain(requireAdmin(handler.FindByActive)))
	routerRole.GET("/trashed", roleMiddlewareChain(requireAdmin(handler.FindByTrashed)))
	routerRole.GET("/user/:user_id", roleMiddlewareChain(requireAdmin(handler.FindByUserId)))

	return handler
}

// @Security Bearer
// @Summary Find all roles
// @Tags Role Query
// @Description Retrieve a list of all roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRole "List of roles"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query [get]
func (h *roleQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoles(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindAllRole(ctx, &pb.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return errors.ParseGrpcError(err) }


	apiResponse := h.mapper.ToApiResponsePaginationRole(res)
	h.cache.SetCachedRoles(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find role by ID
// @Tags Role Query
// @Description Retrieve a role by ID
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} response.ApiResponseRole "Role data"
// @Failure 400 {object} errors.ErrorResponse "Invalid role ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/{id} [get]
func (h *roleQueryHandlerApi) FindById(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil || roleID <= 0 { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedRoleById(ctx, roleID); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindByIdRole(ctx, &pb.FindByIdRoleRequest{RoleId: int32(roleID)})
	if err != nil { return errors.ParseGrpcError(err) }


	apiResponse := h.mapper.ToApiResponseRole(res)
	h.cache.SetCachedRoleById(ctx, roleID, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active roles
// @Tags Role Query
// @Description Retrieve a list of active roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of active roles"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/active [get]
func (h *roleQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoleActive(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindByActive(ctx, &pb.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return errors.ParseGrpcError(err) }


	apiResponse := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)
	h.cache.SetCachedRoleActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed roles
// @Tags Role Query
// @Description Retrieve a list of trashed role records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationRoleDeleteAt "List of trashed role data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/trashed [get]
func (h *roleQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllRole{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedRoleTrashed(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindByTrashed(ctx, &pb.FindAllRoleRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return errors.ParseGrpcError(err) }


	apiResponse := h.mapper.ToApiResponsePaginationRoleDeleteAt(res)
	h.cache.SetCachedRoleTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find roles by user ID
// @Tags Role Query
// @Description Retrieve a list of roles assigned to a specific user
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} response.ApiResponsesRole "List of user roles"
// @Failure 400 {object} errors.ErrorResponse "Invalid user ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve role data"
// @Router /api/role-query/user/{user_id} [get]
func (h *roleQueryHandlerApi) FindByUserId(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil || userID <= 0 { return errors.NewBadRequestError("user_id is required") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedRoleByUserId(ctx, userID); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.role.FindByUserId(ctx, &pb.FindByIdUserRoleRequest{UserId: int32(userID)})
	if err != nil { return errors.ParseGrpcError(err) }


	apiResponse := h.mapper.ToApiResponsesRole(res)
	h.cache.SetCachedRoleByUserId(ctx, userID, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}


