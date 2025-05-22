package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productHandleApi struct {
	client       pb.ProductServiceClient
	logger       logger.LoggerInterface
	mapping      response_api.ProductResponseMapper
	upload_image upload_image.ImageUploads
}

func NewHandlerProduct(
	router *echo.Echo,
	client pb.ProductServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ProductResponseMapper,
	upload_image upload_image.ImageUploads,
) *productHandleApi {
	productHandler := &productHandleApi{
		client:       client,
		logger:       logger,
		mapping:      mapping,
		upload_image: upload_image,
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

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve product data", zap.Error(err))
		return product_errors.ErrApiProductFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationProduct(res)

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

	req := &pb.FindAllProductMerchantRequest{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, req)
	if err != nil {
		h.logger.Error("Failed to retrieve product data by merchant",
			zap.Error(err),
			zap.Int("merchant_id", merchantID),
			zap.Any("request", req),
		)

		return product_errors.ErrApiProductFailedFindByMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationProduct(res)
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

	req := &pb.FindAllProductCategoryRequest{
		CategoryName: categoryName,
		Page:         int32(page),
		PageSize:     int32(pageSize),
		Search:       search,
	}

	res, err := h.client.FindByCategory(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve product data by category", zap.Error(err))
		return product_errors.ErrApiProductFailedFindByCategory(c)
	}

	so := h.mapping.ToApiResponsePaginationProduct(res)
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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid product ID", zap.Error(err))
		return product_errors.ErrApiProductInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve product data", zap.Error(err))
		return product_errors.ErrApiProductFailedFindById(c)
	}

	so := h.mapping.ToApiResponseProduct(res)

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

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve product data", zap.Error(err))
		return product_errors.ErrApiProductFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationProductDeleteAt(res)

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

	req := &pb.FindAllProductRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Debug("Failed to retrieve product data", zap.Error(err))
		return product_errors.ErrApiProductFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationProductDeleteAt(res)

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
	formData, err := h.parseProductForm(c, true)
	if err != nil {
		return product_errors.ErrApiInvalidBody(c)
	}

	ctx := c.Request().Context()
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
			h.upload_image.CleanupImageOnFailure(formData.ImagePath)

			return product_errors.ErrApiProductFailedCreate(c)
		}

		h.logger.Error("Product creation failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return product_errors.ErrApiProductFailedCreate(c)
	}

	so := h.mapping.ToApiResponseProduct(res)
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
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return product_errors.ErrApiProductInvalidId(c)
	}

	formData, err := h.parseProductForm(c, false)
	if err != nil {
		return product_errors.ErrApiInvalidBody(c)
	}

	ctx := c.Request().Context()
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

		h.logger.Error("Product update failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return product_errors.ErrApiProductFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseProduct(res)
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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid product ID format", zap.Error(err))
		return product_errors.ErrApiProductInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedProduct(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive product", zap.Error(err))
		return product_errors.ErrApiProductFailedTrashed(c)
	}

	so := h.mapping.ToApiResponsesProductDeleteAt(res)

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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid product ID format", zap.Error(err))
		return product_errors.ErrApiProductInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreProduct(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore product", zap.Error(err))
		return product_errors.ErrApiProductFailedRestore(c)
	}

	so := h.mapping.ToApiResponsesProductDeleteAt(res)

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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid product ID format", zap.Error(err))
		return product_errors.ErrApiProductInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdProductRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteProductPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to permanently delete product", zap.Error(err))
		return product_errors.ErrApiProductFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseProductDelete(res)

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
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllProduct(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk products restoration failed", zap.Error(err))
		return product_errors.ErrApiProductFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseProductAll(res)

	h.logger.Debug("Successfully restored all products")

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
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllProductPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk products deletion failed", zap.Error(err))
		return product_errors.ErrApiProductFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseProductAll(res)

	h.logger.Debug("Successfully deleted all products permanently")

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
