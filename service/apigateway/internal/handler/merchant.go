package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantHandleApi struct {
	client  pb.MerchantServiceClient
	logger  logger.LoggerInterface
	mapping response_api.MerchantResponseMapper
}

func NewHandlerMerchant(
	router *echo.Echo,
	client pb.MerchantServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantResponseMapper,
) *merchantHandleApi {
	merchantHandler := &merchantHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routercategory := router.Group("/api/merchant")

	routercategory.GET("", merchantHandler.FindAllMerchant)
	routercategory.GET("/:id", merchantHandler.FindById)
	routercategory.GET("/active", merchantHandler.FindByActive)
	routercategory.GET("/trashed", merchantHandler.FindByTrashed)

	routercategory.POST("/create", merchantHandler.Create)
	routercategory.POST("/update/:id", merchantHandler.Update)

	routercategory.POST("/trashed/:id", merchantHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantHandler.DeleteAllMerchantPermanent)

	return merchantHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags Merchant
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchant "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant [get]
func (h *merchantHandleApi) FindAllMerchant(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch merchants", zap.Error(err))
		return merchant_errors.ErrApiFailedFindAllMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchant(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags Merchant
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant/{id} [get]
func (h *merchantHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchant_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch merchant details", zap.Error(err))
		return merchant_errors.ErrApiFailedFindMerchantById(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags Merchant
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant/active [get]
func (h *merchantHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active merchants", zap.Error(err))
		return merchant_errors.ErrApiFailedFindActiveMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags Merchant
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant/trashed [get]
func (h *merchantHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived merchants", zap.Error(err))
		return merchant_errors.ErrApiFailedFindTrashedMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new merchant.
// @Summary Create a new merchant
// @Tags Merchant
// @Description Create a new merchant with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantRequest true "Create merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Successfully created merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant"
// @Router /api/merchant/create [post]
func (h *merchantHandleApi) Create(c echo.Context) error {
	var body requests.CreateMerchantRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchant_errors.ErrApiBindCreateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchant_errors.ErrApiValidateCreateMerchant(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateMerchantRequest{
		UserId:       int32(body.UserID),
		Name:         body.Name,
		Description:  body.Description,
		Address:      body.Address,
		ContactEmail: body.ContactEmail,
		ContactPhone: body.ContactPhone,
		Status:       body.Status,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		h.logger.Error("Merchant creation failed", zap.Error(err))
		return merchant_errors.ErrApiFailedCreateMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing merchant record.
// @Summary Update an existing merchant
// @Tags Merchant
// @Description Update an existing merchant record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Param request body requests.UpdateMerchantRequest true "Update merchant request"
// @Success 200 {object} response.ApiResponseMerchant "Successfully updated merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant"
// @Router /api/merchant/update [post]
func (h *merchantHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))

		return merchant_errors.ErrApiInvalidId(c)
	}

	var body requests.UpdateMerchantRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchant_errors.ErrApiBindUpdateMerchant(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchant_errors.ErrApiValidateUpdateMerchant(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateMerchantRequest{
		MerchantId:   int32(idInt),
		UserId:       int32(body.UserID),
		Name:         body.Name,
		Description:  body.Description,
		Address:      body.Address,
		ContactEmail: body.ContactEmail,
		ContactPhone: body.ContactPhone,
		Status:       body.Status,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		h.logger.Error("Merchant update failed", zap.Error(err))
		return merchant_errors.ErrApiFailedUpdateMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags Merchant
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant/trashed/{id} [get]
func (h *merchantHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchant_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchant(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive merchant", zap.Error(err))
		return merchant_errors.ErrApiFailedTrashMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchantDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags Merchant
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant/restore/{id} [post]
func (h *merchantHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchant_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchant(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore merchant", zap.Error(err))
		return merchant_errors.ErrApiFailedRestoreMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags Merchant
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant/delete/{id} [delete]
func (h *merchantHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchant_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete merchant", zap.Error(err))
		return merchant_errors.ErrApiFailedDeleteMerchantPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags Merchant
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant/restore/all [post]
func (h *merchantHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllMerchant(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant restoration failed", zap.Error(err))
		return merchant_errors.ErrApiFailedRestoreAllMerchant(c)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully restored all merchant")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags Merchant
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant/delete/all [post]
func (h *merchantHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant deletion failed", zap.Error(err))
		return merchant_errors.ErrApiFailedDeleteAllMerchantPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}
