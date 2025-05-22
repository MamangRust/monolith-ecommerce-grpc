package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryHandleApi struct {
	client       pb.CategoryServiceClient
	logger       logger.LoggerInterface
	mapping      response_api.CategoryResponseMapper
	upload_image upload_image.ImageUploads
}

func NewHandlerCategory(
	router *echo.Echo,
	client pb.CategoryServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CategoryResponseMapper,
	upload_image upload_image.ImageUploads,
) *categoryHandleApi {
	categoryHandler := &categoryHandleApi{
		client:       client,
		logger:       logger,
		mapping:      mapping,
		upload_image: upload_image,
	}

	routercategory := router.Group("/api/category")

	routercategory.GET("", categoryHandler.FindAllCategory)
	routercategory.GET("/:id", categoryHandler.FindById)
	routercategory.GET("/active", categoryHandler.FindByActive)
	routercategory.GET("/trashed", categoryHandler.FindByTrashed)

	routercategory.GET("/monthly-total-pricing", categoryHandler.FindMonthTotalPrice)
	routercategory.GET("/yearly-total-pricing", categoryHandler.FindYearTotalPrice)
	routercategory.GET("/merchant/monthly-total-pricing", categoryHandler.FindMonthTotalPriceByMerchant)
	routercategory.GET("/merchant/yearly-total-pricing", categoryHandler.FindYearTotalPriceByMerchant)
	routercategory.GET("/mycategory/monthly-total-pricing", categoryHandler.FindMonthTotalPriceById)
	routercategory.GET("/mycategory/yearly-total-pricing", categoryHandler.FindYearTotalPriceById)

	routercategory.GET("/monthly-pricing", categoryHandler.FindMonthPrice)
	routercategory.GET("/yearly-pricing", categoryHandler.FindYearPrice)
	routercategory.GET("/merchant/monthly-pricing", categoryHandler.FindMonthPriceByMerchant)
	routercategory.GET("/merchant/yearly-pricing", categoryHandler.FindYearPriceByMerchant)
	routercategory.GET("/mycategory/monthly-pricing", categoryHandler.FindMonthPriceById)
	routercategory.GET("/mycategory/yearly-pricing", categoryHandler.FindYearPriceById)

	routercategory.POST("/create", categoryHandler.Create)
	routercategory.POST("/update/:id", categoryHandler.Update)

	routercategory.POST("/trashed/:id", categoryHandler.TrashedCategory)
	routercategory.POST("/restore/:id", categoryHandler.RestoreCategory)
	routercategory.DELETE("/permanent/:id", categoryHandler.DeleteCategoryPermanent)

	routercategory.POST("/restore/all", categoryHandler.RestoreAllCategory)
	routercategory.POST("/permanent/all", categoryHandler.DeleteAllCategoryPermanent)

	return categoryHandler
}

// @Security Bearer
// @Summary Find all category
// @Tags Category
// @Description Retrieve a list of all category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategory "List of category"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category [get]
func (h *categoryHandleApi) FindAllCategory(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch categories", zap.Error(err))
		return category_errors.ErrApiFailedFindAllCategory(c)
	}

	so := h.mapping.ToApiResponsePaginationCategory(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find category by ID
// @Tags Category
// @Description Retrieve a category by ID
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategory "Category data"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/{id} [get]
func (h *categoryHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch category details", zap.Error(err))
		return category_errors.ErrApiFailedFindByIdCategory(c)
	}

	so := h.mapping.ToApiResponseCategory(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active category
// @Tags Category
// @Description Retrieve a list of active category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of active category"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/active [get]
func (h *categoryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active categories", zap.Error(err))
		return category_errors.ErrApiFailedFindByActiveCategory(c)
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed category records.
// @Summary Retrieve trashed category
// @Tags Category
// @Description Retrieve a list of trashed category records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of trashed category data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/trashed [get]
func (h *categoryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived categories", zap.Error(err))
		return category_errors.ErrApiFailedFindByTrashedCategory(c)
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthTotalPrice retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalPrices(ctx, &pb.FindYearMonthTotalPrices{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))

		return category_errors.ErrApiFailedFindMonthTotalPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalPrice retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-total-pricing [get]
func (h *categoryHandleApi) FindYearTotalPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalPrices(ctx, &pb.FindYearTotalPrices{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))

		return category_errors.ErrApiFailedFindYearTotalPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthTotalPriceById retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return category_errors.ErrApiFailedFindMonthTotalPriceById(c)
	}

	categoryStr := c.QueryParam("category_id")

	category, err := strconv.Atoi(categoryStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))
		return category_errors.ErrApiFailedFindMonthTotalPriceById(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalPricesById(ctx, &pb.FindYearMonthTotalPriceById{
		Year:       int32(year),
		Month:      int32(month),
		CategoryId: int32(category),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))

		return category_errors.ErrApiFailedFindMonthTotalPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalPriceById retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-total-pricing/{id} [get]
func (h *categoryHandleApi) FindYearTotalPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	categoryStr := c.QueryParam("category_id")

	category, err := strconv.Atoi(categoryStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalPricesById(ctx, &pb.FindYearTotalPriceById{
		Year:       int32(year),
		CategoryId: int32(category),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))

		return category_errors.ErrApiFailedFindYearTotalPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthTotalPriceByMerchant retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalPricesByMerchant(ctx, &pb.FindYearMonthTotalPriceByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))

		return category_errors.ErrApiFailedFindMonthTotalPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalPriceByMerchant retrieves yearly category total pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/merchant/yearly-total-pricing [get]
func (h *categoryHandleApi) FindYearTotalPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalPricesByMerchant(ctx, &pb.FindYearTotalPriceByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))
		return category_errors.ErrApiFailedFindYearTotalPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthPrice retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))
		return category_errors.ErrApiFailedFindMonthPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearPrice retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPrice(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))
		return category_errors.ErrApiFailedFindYearPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthPriceByMerchant retrieves monthly category pricing by merchant
// @Summary Get monthly category pricing by merchant
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for categories by specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))
		return category_errors.ErrApiFailedFindMonthPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearPriceByMerchant retrieves yearly category pricing by merchant
// @Summary Get yearly category pricing by merchant
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for categories by specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPriceByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	merchantIdStr := c.QueryParam("merchant_id")

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))
		return category_errors.ErrApiFailedFindYearPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthPriceById retrieves monthly pricing for specific category
// @Summary Get monthly pricing by category ID
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for specific category
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly pricing by category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	categoryIdStr := c.QueryParam("category_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	category_id, err := strconv.Atoi(categoryIdStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))
		return category_errors.ErrApiFailedFindMonthPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearPriceById retrieves yearly pricing for specific category
// @Summary Get yearly pricing by category ID
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for specific category
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly pricing by category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPriceById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	categoryIdStr := c.QueryParam("category_id")

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	category_id, err := strconv.Atoi(categoryIdStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))
		return category_errors.ErrApiFailedFindYearPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new category with image upload.
// @Summary Create a new category
// @Tags Category
// @Description Create a new category with the provided details and an image file
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param slug_category formData string true "Category slug"
// @Param image_category formData file true "Category image file"
// @Success 200 {object} response.ApiResponseCategory "Successfully created category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create category"
// @Router /api/category/create [post]
func (h *categoryHandleApi) Create(c echo.Context) error {
	formData, err := h.parseCategoryForm(c, true)
	if err != nil {
		return category_errors.ErrApiInvalidBody(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateCategoryRequest{
		Name:          formData.Name,
		Description:   formData.Description,
		SlugCategory:  *formData.SlugCategory,
		ImageCategory: formData.ImageCategory,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		h.logger.Error("Category creation failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		if formData.ImageCategory != "" {
			h.upload_image.CleanupImageOnFailure(formData.ImageCategory)
		}

		return category_errors.ErrApiFailedCreateCategory(c)
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// Update handles the update of an existing category with image upload.
// @Summary Update an existing category
// @Tags Category
// @Description Update an existing category record with the provided details and an optional image file
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Order ID"
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param slug_category formData string true "Category slug"
// @Param image_category formData file false "New category image file"
// @Success 200 {object} response.ApiResponseCategory "Successfully updated category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update category"
// @Router /api/category/update [post]
func (h *categoryHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidId(c)
	}

	formData, err := h.parseCategoryForm(c, false)
	if err != nil {
		return category_errors.ErrApiInvalidBody(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateCategoryRequest{
		CategoryId:    int32(idInt),
		Name:          formData.Name,
		Description:   formData.Description,
		SlugCategory:  *formData.SlugCategory,
		ImageCategory: formData.ImageCategory,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		if formData.ImageCategory != "" {
			h.upload_image.CleanupImageOnFailure(formData.ImageCategory)
		}

		h.logger.Error("Category update failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return category_errors.ErrApiFailedUpdateCategory(c)
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// TrashedCategory retrieves a trashed category record by its ID.
// @Summary Retrieve a trashed category
// @Tags Category
// @Description Retrieve a trashed category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully retrieved trashed category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed category"
// @Router /api/category/trashed/{id} [get]
func (h *categoryHandleApi) TrashedCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID format", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedCategory(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive category", zap.Error(err))
		return category_errors.ErrApiFailedTrashedCategory(c)
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreCategory restores a category record from the trash by its ID.
// @Summary Restore a trashed category
// @Tags Category
// @Description Restore a trashed category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully restored category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore category"
// @Router /api/category/restore/{id} [post]
func (h *categoryHandleApi) RestoreCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID format", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreCategory(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore category", zap.Error(err))
		return category_errors.ErrApiFailedRestoreCategory(c)
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteCategoryPermanent permanently deletes a category record by its ID.
// @Summary Permanently delete a category
// @Tags Category
// @Description Permanently delete a category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "category ID"
// @Success 200 {object} response.ApiResponseCategoryDelete "Successfully deleted category record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete category:"
// @Router /api/category/delete/{id} [delete]
func (h *categoryHandleApi) DeleteCategoryPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID format", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteCategoryPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to permanently delete category", zap.Error(err))
		return category_errors.ErrApiFailedDeleteCategoryPermanent(c)
	}

	so := h.mapping.ToApiResponseCategoryDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllCategory restores a category record from the trash by its ID.
// @Summary Restore a trashed category
// @Tags Category
// @Description Restore a trashed category record by its ID.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully restored category all"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore category"
// @Router /api/category/restore/all [post]
func (h *categoryHandleApi) RestoreAllCategory(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllCategory(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk category restoration failed", zap.Error(err))
		return category_errors.ErrApiFailedRestoreAllCategories(c)
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	h.logger.Debug("Successfully restored all category")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllCategoryPermanent permanently deletes a category record by its ID.
// @Summary Permanently delete a category
// @Tags Category
// @Description Permanently delete a category record by its ID.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully deleted category record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete category:"
// @Router /api/category/delete/all [post]
func (h *categoryHandleApi) DeleteAllCategoryPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk category deletion failed", zap.Error(err))
		return category_errors.ErrApiFailedDeleteAllCategoriesPermanent(c)
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	h.logger.Debug("Successfully deleted all category permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *categoryHandleApi) parseCategoryForm(c echo.Context, requireImage bool) (requests.CategoryFormData, error) {
	var formData requests.CategoryFormData

	formData.Name = strings.TrimSpace(c.FormValue("name"))
	if formData.Name == "" {
		return formData, category_errors.ErrApiCategoryNameRequired(c)
	}

	formData.Description = strings.TrimSpace(c.FormValue("description"))
	if formData.Description == "" {
		return formData, category_errors.ErrApiCategoryDescriptionRequired(c)
	}

	slug := strings.TrimSpace(c.FormValue("slug_category"))
	if slug == "" {
		return formData, category_errors.ErrApiCategorySlugRequired(c)
	}
	formData.SlugCategory = &slug

	file, err := c.FormFile("image_category")
	if err != nil {
		if requireImage {
			h.logger.Debug("Image upload error", zap.Error(err))
			return formData, category_errors.ErrApiCategoryImageRequired(c)
		}
		return formData, nil
	}

	imagePath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}

	formData.ImageCategory = imagePath
	return formData, nil
}
