package producthandler

import (
	"net/http"
	"strconv"

	product_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/product"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/product"
	"github.com/labstack/echo/v4"
)


type productQueryHandlerApi struct {
	client pb.ProductQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.ProductQueryResponseMapper
	cache  product_cache.ProductQueryCache
	errors errors.ApiHandler
}

type productQueryHandleDeps struct {
	client     pb.ProductQueryServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.ProductQueryResponseMapper
	cache      product_cache.ProductQueryCache
	apiHandler errors.ApiHandler
}

func NewProductQueryHandleApi(params *productQueryHandleDeps) *productQueryHandlerApi {
	handler := &productQueryHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
		cache:  params.cache,
		errors: params.apiHandler,
	}



	routerProduct := params.router.Group("/api/product-query")
	routerProduct.GET("", handler.FindAll)
	routerProduct.GET("/:id", handler.FindById)
	routerProduct.GET("/merchant/:merchant_id", handler.FindByMerchant)
	routerProduct.GET("/category/:category_name", handler.FindByCategory)
	routerProduct.GET("/active", handler.FindByActive)
	routerProduct.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all products
// @Tags Product Query
// @Description Retrieve a list of all products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProduct "List of products"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve product data"
// @Router /api/product-query [get]
func (h *productQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllProduct{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedProducts(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindAll(ctx, &pb.FindAllProductRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationProduct(res)
	h.cache.SetCachedProducts(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find product by ID
// @Tags Product Query
// @Description Retrieve a product by ID
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProduct "Product data"
// @Failure 400 {object} errors.ErrorResponse "Invalid product ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve product data"
// @Router /api/product-query/{id} [get]
func (h *productQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid Product ID") }

	ctx := c.Request().Context()
	if cachedData, found := h.cache.GetCachedProduct(ctx, id); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindById(ctx, &pb.FindByIdProductRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponseProduct(res)
	h.cache.SetCachedProduct(ctx, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find products by merchant ID
// @Tags Product Query
// @Description Retrieve a list of products belonging to a specific merchant
// @Accept json
// @Produce json
// @Param merchant_id path int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProduct "List of products by merchant"
// @Failure 400 {object} errors.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve product data"
// @Router /api/product-query/merchant/{merchant_id} [get]
func (h *productQueryHandlerApi) FindByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid Merchant ID") }

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllProductByMerchant{MerchantID: merchantID, Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedProductsByMerchant(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByMerchant(ctx, &pb.FindAllProductMerchantRequest{
		MerchantId: int32(merchantID), Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationProduct(res)
	h.cache.SetCachedProductsByMerchant(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find products by category name
// @Tags Product Query
// @Description Retrieve a list of products belonging to a specific category
// @Accept json
// @Produce json
// @Param category_name path string true "Category Name"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProduct "List of products by category"
// @Failure 400 {object} errors.ErrorResponse "Invalid category name"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve product data"
// @Router /api/product-query/category/{category_name} [get]
func (h *productQueryHandlerApi) FindByCategory(c echo.Context) error {
	categoryName := c.Param("category_name")
	if categoryName == "" { return echo.NewHTTPError(http.StatusBadRequest, "Category Name is required") }

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllProductByCategory{CategoryName: categoryName, Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedProductsByCategory(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByCategory(ctx, &pb.FindAllProductCategoryRequest{
		CategoryName: categoryName, Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationProduct(res)
	h.cache.SetCachedProductsByCategory(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active products
// @Tags Product Query
// @Description Retrieve a list of active products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProductDeleteAt "List of active products"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve product data"
// @Router /api/product-query/active [get]
func (h *productQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllProduct{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedProductActive(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByActive(ctx, &pb.FindAllProductRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationProductDeleteAt(res)
	h.cache.SetCachedProductActive(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed products
// @Tags Product Query
// @Description Retrieve a list of trashed product records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProductDeleteAt "List of trashed product data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve product data"
// @Router /api/product-query/trashed [get]
func (h *productQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllProduct{Page: page, PageSize: pageSize, Search: search}

	if cachedData, found := h.cache.GetCachedProductTrashed(ctx, req); found {
		return c.JSON(http.StatusOK, cachedData)
	}

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllProductRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	apiResponse := h.mapper.ToApiResponsePaginationProductDeleteAt(res)
	h.cache.SetCachedProductTrashed(ctx, req, apiResponse)

	return c.JSON(http.StatusOK, apiResponse)
}


