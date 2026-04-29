package categoryhandler

import (
	"net/http"
	"strconv"

	category_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/category"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/category"
	"github.com/labstack/echo/v4"
)

type categoryQueryHandlerApi struct {
	client     pb.CategoryQueryServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.CategoryQueryResponseMapper
	cache      category_cache.CategoryMencache
	errors     errors.ApiHandler
}



type categoryQueryHandleDeps struct {
	client     pb.CategoryQueryServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.CategoryQueryResponseMapper
	cache      category_cache.CategoryMencache
	apiHandler errors.ApiHandler
}

func NewCategoryQueryHandleApi(params *categoryQueryHandleDeps) *categoryQueryHandlerApi {
	categoryQueryHandler := &categoryQueryHandlerApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		cache:      params.cache,
		errors:     params.apiHandler,
	}


	routerCategory := params.router.Group("/api/category-query")

	routerCategory.GET("", categoryQueryHandler.FindAll)
	routerCategory.GET("/:id", categoryQueryHandler.FindById)
	routerCategory.GET("/active", categoryQueryHandler.FindByActive)
	routerCategory.GET("/trashed", categoryQueryHandler.FindByTrashed)

	return categoryQueryHandler
}

// @Security Bearer
// @Summary Find all categories
// @Tags Category Query
// @Description Retrieve a list of all categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategory "List of categories"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve category data"
// @Router /api/category-query [get]
func (h *categoryQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllCategory{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedCategoriesCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllCategoryRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}



	apiResponse := h.mapper.ToApiResponsePaginationCategory(res)
	h.cache.SetCachedCategoriesCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find category by ID
// @Tags Category Query
// @Description Retrieve a category by ID
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategory "Category data"
// @Failure 400 {object} errors.ErrorResponse "Invalid category ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve category data"
// @Router /api/category-query/{id} [get]
func (h *categoryQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedCategoryCache(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdCategoryRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponseCategory(res)
	h.cache.SetCachedCategoryCache(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active categories
// @Tags Category Query
// @Description Retrieve a list of active categories
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of active categories"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve category data"
// @Router /api/category-query/active [get]
func (h *categoryQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllCategory{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedCategoryActiveCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllCategoryRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}



	apiResponse := h.mapper.ToApiResponsePaginationCategoryDeleteAt(res)
	h.cache.SetCachedCategoryActiveCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed categories
// @Tags Category Query
// @Description Retrieve a list of trashed category records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of trashed category data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve category data"
// @Router /api/category-query/trashed [get]
func (h *categoryQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllCategory{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedCategoryTrashedCache(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllCategoryRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}



	apiResponse := h.mapper.ToApiResponsePaginationCategoryDeleteAt(res)
	h.cache.SetCachedCategoryTrashedCache(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}


