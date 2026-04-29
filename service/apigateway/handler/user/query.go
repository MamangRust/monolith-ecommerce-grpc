package userhandler

import (
	"net/http"
	"strconv"

	user_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/user"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/user"
	"github.com/labstack/echo/v4"
)


type userQueryHandleApi struct {
	client pb.UserQueryServiceClient

	logger logger.LoggerInterface

	mapper apimapper.UserQueryResponseMapper

	cache user_cache.UserMencache

	apiHandler errors.ApiHandler
}

type userQueryHandleDeps struct {
	client pb.UserQueryServiceClient

	router *echo.Echo

	logger logger.LoggerInterface

	mapper apimapper.UserQueryResponseMapper

	cache user_cache.UserMencache

	apiHandler errors.ApiHandler
}

func NewUserQueryHandleApi(params *userQueryHandleDeps) *userQueryHandleApi {

	userQueryHandleApi := &userQueryHandleApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		apiHandler: params.apiHandler,
	}

	routerUser := params.router.Group("/api/user-query")

	routerUser.GET("", params.apiHandler.Handle("find-all-users", userQueryHandleApi.FindAllUser))
	routerUser.GET("/:id", params.apiHandler.Handle("find-user-by-id", userQueryHandleApi.FindById))
	routerUser.GET("/active", params.apiHandler.Handle("find-active-users", userQueryHandleApi.FindByActive))
	routerUser.GET("/trashed", params.apiHandler.Handle("find-trashed-users", userQueryHandleApi.FindByTrashed))

	return userQueryHandleApi
}

// @Security Bearer
// @Summary Find all users
// @Tags User Query
// @Description Retrieve a list of all users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUser "List of users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query [get]
func (h *userQueryHandleApi) FindAllUser(c echo.Context) error {
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

	req := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedUsersCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, grpcReq)
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationUser(res)

	h.cache.SetCachedUsersCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find user by ID
// @Tags User Query
// @Description Retrieve a user by ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.ApiResponseUser "User data"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query/{id} [get]
func (h *userQueryHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("id is required")
	}

	ctx := c.Request().Context()

	cachedData, found := h.cache.GetCachedUserCache(ctx, id)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	req := &pb.FindByIdUserRequest{
		Id: int32(id),
	}

	user, err := h.client.FindById(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponseUser(user)

	h.cache.SetCachedUserCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active users
// @Tags User Query
// @Description Retrieve a list of active users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUserDeleteAt "List of active users"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query/active [get]
func (h *userQueryHandleApi) FindByActive(c echo.Context) error {
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

	req := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedUserActiveCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, grpcReq)
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationUserDeleteAt(res)

	h.cache.SetCachedUserActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed users
// @Tags User Query
// @Description Retrieve a list of trashed user records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationUserDeleteAt "List of trashed user data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve user data"
// @Router /api/user-query/trashed [get]
func (h *userQueryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	cachedData, found := h.cache.GetCachedUserTrashedCache(ctx, req)
	if found {
		return c.JSON(http.StatusOK, cachedData)
	}

	grpcReq := &pb.FindAllUserRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, grpcReq)
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationUserDeleteAt(res)

	h.cache.SetCachedUserTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}


