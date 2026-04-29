package producthandler

import (
	"net/http"
	"strconv"

	product_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/product"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/product"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)



type productCommandHandlerApi struct {
	client       pb.ProductCommandServiceClient
	logger       logger.LoggerInterface
	mapper       apimapper.ProductCommandResponseMapper
	cache        product_cache.ProductCommandCache
	upload_image upload_image.ImageUploads
	errors       errors.ApiHandler
}
 
type productCommandHandleDeps struct {
	client       pb.ProductCommandServiceClient
	router       *echo.Echo
	logger       logger.LoggerInterface
	mapper       apimapper.ProductCommandResponseMapper
	cache        product_cache.ProductCommandCache
	upload_image upload_image.ImageUploads
	apiHandler   errors.ApiHandler
}



func NewProductCommandHandleApi(params *productCommandHandleDeps) *productCommandHandlerApi {
	handler := &productCommandHandlerApi{
		client:       params.client,
		logger:       params.logger,
		mapper:       params.mapper,
		cache:        params.cache,
		upload_image: params.upload_image,
		errors:       params.apiHandler,
	}



	routerProduct := params.router.Group("/api/product-command")
	routerProduct.POST("/create", handler.Create)
	routerProduct.POST("/update/:id", handler.Update)
	routerProduct.POST("/trashed/:id", handler.Trashed)
	routerProduct.POST("/restore/:id", handler.Restore)
	routerProduct.DELETE("/permanent/:id", handler.DeletePermanent)
	routerProduct.POST("/restore/all", handler.RestoreAll)
	routerProduct.POST("/permanent/all", handler.DeleteAllPermanent)

	return handler
}

// @Security Bearer
// @Summary Create a new product
// @Tags Product Command
// @Description Create a new product with the provided details
// @Accept mpfd
// @Produce json
// @Param merchant_id formData int true "Merchant ID"
// @Param category_id formData int true "Category ID"
// @Param name formData string true "Product name"
// @Param description formData string true "Product description"
// @Param price formData int true "Price"
// @Param count_in_stock formData int true "Stock count"
// @Param brand formData string true "Brand"
// @Param weight formData int true "Weight"
// @Param image formData file false "Product image"
// @Success 201 {object} response.ApiResponseProduct "Successfully created product"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create product"
// @Router /api/product-command/create [post]
func (h *productCommandHandlerApi) Create(c echo.Context) error {
	merchantID, _ := strconv.Atoi(c.FormValue("merchant_id"))
	categoryID, _ := strconv.Atoi(c.FormValue("category_id"))
	name := c.FormValue("name")
	description := c.FormValue("description")
	price, _ := strconv.Atoi(c.FormValue("price"))
	countInStock, _ := strconv.Atoi(c.FormValue("count_in_stock"))
	brand := c.FormValue("brand")
	weight, _ := strconv.Atoi(c.FormValue("weight"))
	rating, _ := strconv.Atoi(c.FormValue("rating"))
	slugProduct := c.FormValue("slug_product")

	imagePath := ""
	file, err := c.FormFile("image")
	if err == nil {
		path, err := h.upload_image.ProcessImageUpload(c, "uploads/products", file, false)
		if err == nil {
			imagePath = path
		}
	}

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateProductRequest{
		MerchantId:   int32(merchantID),
		CategoryId:   int32(categoryID),
		Name:         name,
		Description:  description,
		Price:        int32(price),
		CountInStock: int32(countInStock),
		Brand:        brand,
		Weight:       int32(weight),
		Rating:       int32(rating),
		SlugProduct:  slugProduct,
		ImageProduct: imagePath,
	})
	if err != nil {
		if imagePath != "" {
			h.upload_image.CleanupImageOnFailure(imagePath)
		}
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedProduct(ctx, 0)

	return c.JSON(http.StatusCreated, h.mapper.ToApiResponseProduct(res))
}

// @Security Bearer
// @Summary Update an existing product
// @Tags Product Command
// @Description Update an existing product record
// @Accept mpfd
// @Produce json
// @Param id path int true "Product ID"
// @Param merchant_id formData int false "Merchant ID"
// @Param category_id formData int false "Category ID"
// @Param name formData string false "Product name"
// @Param description formData string false "Product description"
// @Param price formData int false "Price"
// @Param count_in_stock formData int false "Stock count"
// @Param brand formData string false "Brand"
// @Param weight formData int false "Weight"
// @Param image formData file false "Product image"
// @Success 200 {object} response.ApiResponseProduct "Successfully updated product"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update product"
// @Router /api/product-command/update/{id} [post]
func (h *productCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	merchantID, _ := strconv.Atoi(c.FormValue("merchant_id"))
	categoryID, _ := strconv.Atoi(c.FormValue("category_id"))
	name := c.FormValue("name")
	description := c.FormValue("description")
	price, _ := strconv.Atoi(c.FormValue("price"))
	countInStock, _ := strconv.Atoi(c.FormValue("count_in_stock"))
	brand := c.FormValue("brand")
	weight, _ := strconv.Atoi(c.FormValue("weight"))
	rating, _ := strconv.Atoi(c.FormValue("rating"))
	slugProduct := c.FormValue("slug_product")

	imagePath := ""
	file, err := c.FormFile("image")
	if err == nil {
		path, err := h.upload_image.ProcessImageUpload(c, "uploads/products", file, false)
		if err == nil {
			imagePath = path
		} else {
             // Handle error properly or continue if image is optional but here we might want to fail if it's provided but invalid
		}
	}

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateProductRequest{
		ProductId:    int32(id),
		MerchantId:   int32(merchantID),
		CategoryId:   int32(categoryID),
		Name:         name,
		Description:  description,
		Price:        int32(price),
		CountInStock: int32(countInStock),
		Brand:        brand,
		Weight:       int32(weight),
		Rating:       int32(rating),
		SlugProduct:  slugProduct,
		ImageProduct: imagePath,
	})
	if err != nil {
		if imagePath != "" {
			h.upload_image.CleanupImageOnFailure(imagePath)
		}
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedProduct(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseProduct(res))
}

// @Security Bearer
// @Summary Move product to trash
// @Tags Product Command
// @Description Move a product record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDeleteAt "Successfully moved product to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid product ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move product to trash"
// @Router /api/product-command/trashed/{id} [post]
func (h *productCommandHandlerApi) Trashed(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.TrashedProduct(ctx, &pb.FindByIdProductRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedProduct(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponsesProductDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a trashed product
// @Tags Product Command
// @Description Restore a trashed product record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDeleteAt "Successfully restored product"
// @Failure 400 {object} errors.ErrorResponse "Invalid product ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore product"
// @Router /api/product-command/restore/{id} [post]
func (h *productCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.RestoreProduct(ctx, &pb.FindByIdProductRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedProduct(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponsesProductDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete a product
// @Tags Product Command
// @Description Permanently delete a product record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.ApiResponseProductDelete "Successfully deleted product record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid product ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete product permanently"
// @Router /api/product-command/permanent/{id} [delete]
func (h *productCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.DeleteProductPermanent(ctx, &pb.FindByIdProductRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedProduct(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseProductDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed products
// @Tags Product Command
// @Description Restore all trashed product records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseProductAll "Successfully restored all products"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore products"
// @Router /api/product-command/restore/all [post]
func (h *productCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllProduct(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedProduct(ctx, 0)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseProductAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed products
// @Tags Product Command
// @Description Permanently delete all trashed product records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseProductAll "Successfully deleted all products permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete products permanently"
// @Router /api/product-command/permanent/all [post]
func (h *productCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllProductPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedProduct(ctx, 0)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseProductAll(res))
}


