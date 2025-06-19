package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productHandleApi struct {
	client          pb.ProductServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.ProductResponseMapper
	upload_image    upload_image.ImageUploads
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerProduct(
	router *echo.Echo,
	client pb.ProductServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ProductResponseMapper,
	upload_image upload_image.ImageUploads,
) *productHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_handler_requests_total",
			Help: "Total number of product requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "product_handler_request_duration_seconds",
			Help:    "Duration of product requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	productHandler := &productHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		upload_image:    upload_image,
		trace:           otel.Tracer("product-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/product")

	routercategory.GET("", productHandler.FindAllProduct)
	routercategory.GET("/:id", productHandler.FindById)
	routercategory.GET("/merchant/:merchant_id", productHandler.FindByMerchant)
	routercategory.GET("/category/:category_name", productHandler.FindByCategory)

	routercategory.GET("/active", productHandler.FindByActive)
	routercategory.GET("/trashed", productHandler.FindByTrashed)

	routercategory.POST("/create", productHandler.Create)
	routercategory.POST("/update/:id", productHandler.Update)

	routercategory.POST("/trashed/:id", productHandler.TrashedProduct)
	routercategory.POST("/restore/:id", productHandler.RestoreProduct)
	routercategory.DELETE("/permanent/:id", productHandler.DeleteProductPermanent)

	routercategory.POST("/restore/all", productHandler.RestoreAllProduct)
	routercategory.POST("/permanent/all", productHandler.DeleteAllProductPermanent)

	return productHandler
}

// @Security Bearer
// @Summary Find all products
// @Tags Product
// @Description Retrieve a list of all products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProduct "List of products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product [get]
func (h *productHandleApi) FindAllProduct(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllProduct"
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

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve product data", err, zap.Error(err))

		return product_errors.ErrApiProductFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationProduct(res)

	logSuccess("Successfully retrieved product data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find products by merchant
// @Tags Product
// @Description Retrieve a list of products filtered by merchant
// @Accept json
// @Produce json
// @Param merchant_id query string true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProduct "List of products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/merchant [get]
func (h *productHandleApi) FindByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		h.logger.Debug("Invalid merchant ID", zap.Error(err))
		return product_errors.ErrApiProductInvalidMerchantId(c)
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByMerchant"
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

	req := &pb.FindAllProductMerchantRequest{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to retrieve product data", err, zap.Error(err))

		return product_errors.ErrApiProductFailedFindByMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationProduct(res)

	logSuccess("Successfully retrieved product data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find products by category
// @Tags Product
// @Description Retrieve a list of products filtered by category
// @Accept json
// @Produce json
// @Param category_name query string true "Category Name"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProduct "List of products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/category [get]
func (h *productHandleApi) FindByCategory(c echo.Context) error {
	categoryName := c.Param("category_name")
	if categoryName == "" {
		return product_errors.ErrApiProductInvalidCategoryName(c)
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByCategory"
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

	req := &pb.FindAllProductCategoryRequest{
		CategoryName: categoryName,
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Search:       search,
	}

	res, err := h.client.FindByCategory(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve product data", err, zap.Error(err))

		return product_errors.ErrApiProductFailedFindByCategory(c)
	}

	so := h.mapping.ToApiResponsePaginationProduct(res)

	logSuccess("Successfully retrieved product data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find product by ID
// @Tags Product
// @Description Retrieve a product by ID
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProduct "Product data"
// @Failure 400 {object} response.ErrorResponse "Invalid product ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/{id} [get]
func (h *productHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to retrieve product data", err, zap.Error(err))

		return product_errors.ErrApiProductInvalidId(c)
	}

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve product data", err, zap.Error(err))

		return product_errors.ErrApiProductFailedFindById(c)
	}

	so := h.mapping.ToApiResponseProduct(res)

	logSuccess("Successfully retrieved product data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active products
// @Tags Product
// @Description Retrieve a list of active products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProductDeleteAt "List of active products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/active [get]
func (h *productHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve product data", err, zap.Error(err))

		return product_errors.ErrApiProductFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationProductDeleteAt(res)

	logSuccess("Successfully retrieved product data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed products
// @Tags Product
// @Description Retrieve a list of trashed products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationProductDeleteAt "List of trashed products"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve product data"
// @Router /api/product/trashed [get]
func (h *productHandleApi) FindByTrashed(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllMerchantPolicy"
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

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve product data", err, zap.Error(err))

		return product_errors.ErrApiProductFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationProductDeleteAt(res)

	logSuccess("Successfully retrieved product data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// CreateProduct handles the creation of a new product with image upload.
// @Summary Create a new product
// @Tags Product
// @Description Create a new product with the provided details and an image file
// @Accept multipart/form-data
// @Produce json
// @Param merchant_id formData string true "Merchant ID"
// @Param category_id formData string true "Category ID"
// @Param name formData string true "Product name"
// @Param description formData string true "Product description"
// @Param price formData number true "Product price"
// @Param count_in_stock formData integer true "Product count in stock"
// @Param brand formData string true "Product brand"
// @Param weight formData number true "Product weight"
// @Param rating formData number true "Product rating"
// @Param slug_product formData string true "Product slug"
// @Param image formData file true "Product image file"
// @Success 200 {object} response.ApiResponseProduct "Successfully created product"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create product"
// @Router /api/product/create [post]
func (h *productHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	formData, err := h.parseProductForm(c, true)
	if err != nil {
		status = "error"

		logError("Failed to create product", err, zap.Error(err))

		return product_errors.ErrApiInvalidBody(c)
	}

	req := &pb.CreateProductRequest{
		MerchantId:   int32(formData.MerchantID),
		CategoryId:   int32(formData.CategoryID),
		Name:         formData.Name,
		Description:  formData.Description,
		Price:        int32(formData.Price),
		CountInStock: int32(formData.CountInStock),
		Brand:        formData.Brand,
		Weight:       int32(formData.Weight),
		ImageProduct: formData.ImagePath,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		if formData.ImagePath != "" {
			status = "error"

			logError("Failed to create product", err, zap.Error(err))

			h.upload_image.CleanupImageOnFailure(formData.ImagePath)

			return product_errors.ErrApiProductFailedCreate(c)
		}

		status = "error"

		logError("Failed to create product", err, zap.Error(err))

		return product_errors.ErrApiProductFailedCreate(c)
	}

	so := h.mapping.ToApiResponseProduct(res)

	logSuccess("Successfully created product", zap.Bool("success", true))

	return c.JSON(http.StatusCreated, so)
}

// @Security Bearer
// UpdateProduct handles the update of an existing product with optional image upload.
// @Summary Update an existing product
// @Tags Product
// @Description Update an existing product record with the provided details and an optional image file
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param merchant_id formData string true "Merchant ID"
// @Param category_id formData string true "Category ID"
// @Param name formData string true "Product name"
// @Param description formData string true "Product description"
// @Param price formData number true "Product price"
// @Param count_in_stock formData integer true "Product count in stock"
// @Param brand formData string true "Product brand"
// @Param weight formData number true "Product weight"
// @Param rating formData number true "Product rating"
// @Param slug_product formData string true "Product slug"
// @Param image formData file true "Product image file"
// @Success 200 {object} response.ApiResponseProduct "Successfully updated product"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update product"
// @Router /api/product/update/{id} [post]
func (h *productHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		status = "error"

		logError("Failed to update product", err, zap.Error(err))

		return product_errors.ErrApiProductInvalidId(c)
	}

	formData, err := h.parseProductForm(c, false)
	if err != nil {
		status = "error"

		logError("Failed to update product", err, zap.Error(err))

		return product_errors.ErrApiInvalidBody(c)
	}

	req := &pb.UpdateProductRequest{
		ProductId:    int32(productID),
		MerchantId:   int32(formData.MerchantID),
		CategoryId:   int32(formData.CategoryID),
		Name:         formData.Name,
		Description:  formData.Description,
		Price:        int32(formData.Price),
		CountInStock: int32(formData.CountInStock),
		Brand:        formData.Brand,
		Weight:       int32(formData.Weight),
		ImageProduct: formData.ImagePath,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		if formData.ImagePath != "" {
			h.upload_image.CleanupImageOnFailure(formData.ImagePath)
		}

		status = "error"

		logError("Failed to update product", err, zap.Error(err))

		return product_errors.ErrApiProductFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseProduct(res)

	logSuccess("Successfully updated product", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedProduct retrieves a trashed product record by its ID.
// @Summary Retrieve a trashed product
// @Tags Product
// @Description Retrieve a trashed product record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDeleteAt "Successfully retrieved trashed product"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed product"
// @Router /api/product/trashed/{id} [get]
func (h *productHandleApi) TrashedProduct(c echo.Context) error {
	const method = "TrashedProduct"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to archive product", err, zap.Error(err))

		return product_errors.ErrApiProductInvalidId(c)
	}

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedProduct(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to archive product", err, zap.Error(err))

		return product_errors.ErrApiProductFailedTrashed(c)
	}

	so := h.mapping.ToApiResponsesProductDeleteAt(res)

	logSuccess("Successfully retrieved trashed product", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreProduct restores a product record from the trash by its ID.
// @Summary Restore a trashed product
// @Tags Product
// @Description Restore a trashed product record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDeleteAt "Successfully restored product"
// @Failure 400 {object} response.ErrorResponse "Invalid product ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore product"
// @Router /api/product/restore/{id} [post]
func (h *productHandleApi) RestoreProduct(c echo.Context) error {
	const method = "RestoreProduct"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to restore product", err, zap.Error(err))

		return product_errors.ErrApiProductInvalidId(c)
	}

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreProduct(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore product", err, zap.Error(err))

		return product_errors.ErrApiProductFailedRestore(c)
	}

	so := h.mapping.ToApiResponsesProductDeleteAt(res)

	logSuccess("Successfully restored product", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteProductPermanent permanently deletes a product record by its ID.
// @Summary Permanently delete a product
// @Tags Product
// @Description Permanently delete a product record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDelete "Successfully deleted product record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete product:"
// @Router /api/product/delete/{id} [delete]
func (h *productHandleApi) DeleteProductPermanent(c echo.Context) error {
	const method = "DeleteProductPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to permanently delete product", err, zap.Error(err))

		return product_errors.ErrApiProductInvalidId(c)
	}

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteProductPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to permanently delete product", err, zap.Error(err))

		return product_errors.ErrApiProductFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseProductDelete(res)

	logSuccess("Successfully deleted product record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllProduct restores all trashed product records.
// @Summary Restore all trashed products
// @Tags Product
// @Description Restore all trashed product records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseProductAll "Successfully restored all products"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all products"
// @Router /api/product/restore/all [post]
func (h *productHandleApi) RestoreAllProduct(c echo.Context) error {
	const method = "RestoreAllProduct"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllProduct(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all products", err, zap.Error(err))

		return product_errors.ErrApiProductFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseProductAll(res)

	logSuccess("Successfully restored all products", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllProductPermanent permanently deletes all product records.
// @Summary Permanently delete all products
// @Tags Product
// @Description Permanently delete all product records.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseProductAll "Successfully deleted all product records permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete all products"
// @Router /api/product/delete/all [post]
func (h *productHandleApi) DeleteAllProductPermanent(c echo.Context) error {
	const method = "DeleteAllProductPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllProductPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to permanently delete all products", err, zap.Error(err))

		return product_errors.ErrApiProductFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseProductAll(res)

	logSuccess("Successfully deleted all product records permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (h *productHandleApi) parseProductForm(c echo.Context, requireImage bool) (requests.ProductFormData, error) {
	var formData requests.ProductFormData
	var err error

	formData.MerchantID, err = strconv.Atoi(c.FormValue("merchant_id"))
	if err != nil || formData.MerchantID <= 0 {
		return formData, product_errors.ErrApiProductInvalidMerchantId(c)
	}

	formData.CategoryID, err = strconv.Atoi(c.FormValue("category_id"))
	if err != nil || formData.CategoryID <= 0 {
		return formData, product_errors.ErrApiProductInvalidCategoryId(c)
	}

	formData.Name = strings.TrimSpace(c.FormValue("name"))
	if formData.Name == "" {
		return formData, product_errors.ErrApiProductNameRequired(c)
	}

	formData.Description = strings.TrimSpace(c.FormValue("description"))
	formData.Brand = strings.TrimSpace(c.FormValue("brand"))

	formData.Price, err = strconv.Atoi(c.FormValue("price"))
	if err != nil || formData.Price <= 0 {
		return formData, product_errors.ErrApiProductInvalidPrice(c)
	}

	formData.CountInStock, err = strconv.Atoi(c.FormValue("count_in_stock"))
	if err != nil || formData.CountInStock < 0 {
		return formData, product_errors.ErrApiProductInvalidStock(c)
	}

	formData.Weight, err = strconv.Atoi(c.FormValue("weight"))
	if err != nil || formData.Weight <= 0 {
		return formData, product_errors.ErrApiProductInvalidWeight(c)
	}

	file, err := c.FormFile("image_product")
	if err != nil {
		if requireImage {
			h.logger.Debug("Image upload error", zap.Error(err))
			return formData, product_errors.ErrApiProductImageRequired(c)
		}

		return formData, nil
	}

	imagePath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}

	formData.ImagePath = imagePath
	return formData, nil
}

func (s *productHandleApi) startTracingAndLogging(
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

func (s *productHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
