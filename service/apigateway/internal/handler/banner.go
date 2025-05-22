package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bannerHandleApi struct {
	client  pb.BannerServiceClient
	logger  logger.LoggerInterface
	mapping response_api.BannerResponseMapper
}

func NewHandleBanner(
	router *echo.Echo,
	client pb.BannerServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.BannerResponseMapper,
) *bannerHandleApi {
	bannerHandler := &bannerHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routercategory := router.Group("/api/banner")

	routercategory.GET("", bannerHandler.FindAllBanner)
	routercategory.GET("/:id", bannerHandler.FindById)
	routercategory.GET("/active", bannerHandler.FindByActive)
	routercategory.GET("/trashed", bannerHandler.FindByTrashed)

	routercategory.POST("/create", bannerHandler.Create)
	routercategory.POST("/update/:id", bannerHandler.Update)

	routercategory.POST("/trashed/:id", bannerHandler.TrashedBanner)
	routercategory.POST("/restore/:id", bannerHandler.RestoreBanner)
	routercategory.DELETE("/permanent/:id", bannerHandler.DeleteBannerPermanent)

	routercategory.POST("/restore/all", bannerHandler.RestoreAllBanner)
	routercategory.POST("/permanent/all", bannerHandler.DeleteAllBannerPermanent)

	return bannerHandler
}

// @Security Bearer
// @Summary Find all banners
// @Tags Banner
// @Description Retrieve a list of all banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBanner "List of banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner [get]
func (h *bannerHandleApi) FindAllBanner(c echo.Context) error {
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

	req := &pb.FindAllBannerRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch Banners", zap.Error(err))
		return banner_errors.ErrApiFailedFindAllBanner(c)
	}

	so := h.mapping.ToApiResponsePaginationBanner(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find banner by ID
// @Tags Banner
// @Description Retrieve a banner by ID
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBanner "Banner data"
// @Failure 400 {object} response.ErrorResponse "Invalid banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/{id} [get]
func (h *bannerHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid Banner ID format", zap.Error(err))
		return banner_errors.ErrApiBannerInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch Banner details", zap.Error(err))
		return banner_errors.ErrApiFailedFindById(c)
	}

	so := h.mapping.ToApiResponseBanner(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active banners
// @Tags Banner
// @Description Retrieve a list of active banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBannerDeleteAt "List of active banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/active [get]
func (h *bannerHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllBannerRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active Banners", zap.Error(err))
		return banner_errors.ErrApiFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationBannerDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed banners
// @Tags Banner
// @Description Retrieve a list of trashed banners
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationBannerDeleteAt "List of active banners"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve banner data"
// @Router /api/banner/trashed [get]
func (h *bannerHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllBannerRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch trashed Banners", zap.Error(err))
		return banner_errors.ErrApiFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationBannerDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new banner.
// @Summary Create a new banner
// @Tags Banner
// @Description Create a new banner with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateBannerRequest true "Create banner request"
// @Success 200 {object} response.ApiResponseBanner "Successfully created banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create banner"
// @Router /api/banner/create [post]
func (h *bannerHandleApi) Create(c echo.Context) error {
	var body requests.CreateBannerRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return banner_errors.ErrApiBindCreateBanner(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return banner_errors.ErrApiValidateCreateBanner(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateBannerRequest{
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("Banner creation failed", zap.Error(err))
		return banner_errors.ErrApiFailedCreateBanner(c)
	}

	so := h.mapping.ToApiResponseBanner(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing banner record.
// @Summary Update an existing banner
// @Tags Banner
// @Description Update an existing banner record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Param request body requests.UpdateBannerRequest true "Update banner request"
// @Success 200 {object} response.ApiResponseBanner "Successfully updated banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update banner"
// @Router /api/banner/update/{id} [post]
func (h *bannerHandleApi) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))
		return banner_errors.ErrApiBannerInvalidId(c)
	}

	var body requests.UpdateBannerRequest
	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return banner_errors.ErrApiBindUpdateBanner(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return banner_errors.ErrApiValidateCreateBanner(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateBannerRequest{
		BannerId:  int32(idInt),
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		h.logger.Error("Banner update failed", zap.Error(err))
		return banner_errors.ErrApiFailedUpdateBanner(c)
	}

	so := h.mapping.ToApiResponseBanner(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedBanner retrieves a trashed Banner record by its ID.
// @Summary Retrieve a trashed Banner
// @Tags Banner
// @Description Retrieve a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Successfully retrieved trashed Banner"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed Banner"
// @Router /api/banner/trashed/{id} [get]
func (h *bannerHandleApi) TrashedBanner(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid banner ID format", zap.Error(err))
		return banner_errors.ErrApiBannerInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedBanner(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive banner", zap.Error(err))
		return banner_errors.ErrApiFailedTrashedBanner(c)
	}

	so := h.mapping.ToApiResponseBannerDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreBanner restores a Banner record from the trash by its ID.
// @Summary Restore a trashed Banner
// @Tags Banner
// @Description Restore a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Successfully restored Banner"
// @Failure 400 {object} response.ErrorResponse "Invalid Banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore Banner"
// @Router /api/banner/restore/{id} [post]
func (h *bannerHandleApi) RestoreBanner(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid Banner ID format", zap.Error(err))
		return banner_errors.ErrApiBannerInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreBanner(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore Banner", zap.Error(err))
		return banner_errors.ErrApiFailedRestoreBanner(c)
	}

	so := h.mapping.ToApiResponseBannerDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteBannerPermanent permanently deletes a Banner record by its ID.
// @Summary Permanently delete a Banner
// @Tags Banner
// @Description Permanently delete a Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDelete "Successfully deleted Banner record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete Banner:"
// @Router /api/banner/delete/{id} [delete]
func (h *bannerHandleApi) DeleteBannerPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid Banner ID format", zap.Error(err))
		return banner_errors.ErrApiBannerInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdBannerRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteBannerPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete Banner", zap.Error(err))
		return banner_errors.ErrApiFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseBannerDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllBanner restores a Banner record from the trash by its ID.
// @Summary Restore a trashed Banner
// @Tags Banner
// @Description Restore a trashed Banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerAll "Successfully restored Banner all"
// @Failure 400 {object} response.ErrorResponse "Invalid Banner ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore Banner"
// @Router /api/banner/restore/all [post]
func (h *bannerHandleApi) RestoreAllBanner(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllBanner(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk Banner restoration failed", zap.Error(err))
		return banner_errors.ErrApiFailedRestoreAllBanner(c)
	}

	so := h.mapping.ToApiResponseBannerAll(res)

	h.logger.Debug("Successfully restored all Banner")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllBannerPermanent permanently deletes a banner record by its ID.
// @Summary Permanently delete a banner
// @Tags Banner
// @Description Permanently delete a banner record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "banner ID"
// @Success 200 {object} response.ApiResponseBannerAll "Successfully deleted banner record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete banner:"
// @Router /api/banner/delete/all [post]
func (h *bannerHandleApi) DeleteAllBannerPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllBannerPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk banner deletion failed", zap.Error(err))
		return banner_errors.ErrApiFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseBannerAll(res)

	h.logger.Debug("Successfully deleted all banner permanently")

	return c.JSON(http.StatusOK, so)
}
