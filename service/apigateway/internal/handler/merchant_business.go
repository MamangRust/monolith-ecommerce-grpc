package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantBusinessHandleApi struct {
	client          pb.MerchantBusinessServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantBusinessResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
}

func NewHandlerMerchantBusiness(
	router *echo.Echo,
	client pb.MerchantBusinessServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantBusinessResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
) *merchantBusinessHandleApi {
	merchantBusinessHandler := &merchantBusinessHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
	}

	routercategory := router.Group("/api/merchant-business")

	routercategory.GET("", merchantBusinessHandler.FindAllMerchantBusiness)
	routercategory.GET("/:id", merchantBusinessHandler.FindById)
	routercategory.GET("/active", merchantBusinessHandler.FindByActive)
	routercategory.GET("/trashed", merchantBusinessHandler.FindByTrashed)

	routercategory.POST("/create", merchantBusinessHandler.Create)
	routercategory.POST("/update/:id", merchantBusinessHandler.Update)

	routercategory.POST("/trashed/:id", merchantBusinessHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantBusinessHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantBusinessHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantBusinessHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantBusinessHandler.DeleteAllMerchantPermanent)

	return merchantBusinessHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags MerchantCertification
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusiness "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business [get]
func (h *merchantBusinessHandleApi) FindAllMerchantBusiness(c echo.Context) error {
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
		return merchantbusiness_errors.ErrApiFailedFindAllMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantBusiness(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantCertification
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusiness "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/{id} [get]
func (h *merchantBusinessHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch merchant details", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedFindByIdMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusiness(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags MerchantCertification
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/active [get]
func (h *merchantBusinessHandleApi) FindByActive(c echo.Context) error {
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
		return merchantbusiness_errors.ErrApiFailedFindByActiveMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantBusinessDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantBusinessDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-business/trashed [get]
func (h *merchantBusinessHandleApi) FindByTrashed(c echo.Context) error {
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
		return merchantbusiness_errors.ErrApiFailedFindByTrashedMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantBusinessDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new merchant business information
// @Tags MerchantBusiness
// @Description Create merchant business info (e.g., type, tax ID, website, etc.)
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantBusinessInformationRequest true "Create merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully created merchant business info"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant business info"
// @Router /api/merchant-business/create [post]
func (h *merchantBusinessHandleApi) Create(c echo.Context) error {
	var body requests.CreateMerchantBusinessInformationRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchantbusiness_errors.ErrApiBindCreateMerchantBusiness(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchantbusiness_errors.ErrApiValidateCreateMerchantBusiness(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateMerchantBusinessRequest{
		MerchantId:        int32(body.MerchantID),
		BusinessType:      body.BusinessType,
		TaxId:             body.TaxID,
		EstablishedYear:   int32(body.EstablishedYear),
		NumberOfEmployees: int32(body.NumberOfEmployees),
		WebsiteUrl:        body.WebsiteUrl,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("Merchant Business creation failed", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedCreateMerchantBusiness(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponseMerchantBusiness(res))
}

// @Security Bearer
// @Summary Update existing merchant business information
// @Tags MerchantBusiness
// @Description Update merchant business info by ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant Business Info ID"
// @Param request body requests.UpdateMerchantBusinessInformationRequest true "Update merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully updated merchant business info"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant business info"
// @Router /api/merchant-business/update/{id} [post]
func (h *merchantBusinessHandleApi) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	var body requests.UpdateMerchantBusinessInformationRequest
	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchantbusiness_errors.ErrApiBindUpdateMerchantBusiness(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchantbusiness_errors.ErrApiValidateUpdateMerchantBusiness(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateMerchantBusinessRequest{
		MerchantBusinessInfoId: int32(idInt),
		BusinessType:           body.BusinessType,
		TaxId:                  body.TaxID,
		EstablishedYear:        int32(body.EstablishedYear),
		NumberOfEmployees:      int32(body.NumberOfEmployees),
		WebsiteUrl:             body.WebsiteUrl,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		h.logger.Error("Merchant Business update failed", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedUpdateMerchantBusiness(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponseMerchantBusiness(res))
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-business/trashed/{id} [get]
func (h *merchantBusinessHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantBusiness(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive merchant", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedTrashMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusinessDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantCertification
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-business/restore/{id} [post]
func (h *merchantBusinessHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantBusiness(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore merchant", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedRestoreMerchantBusiness(c)
	}

	so := h.mapping.ToApiResponseMerchantBusinessDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantCertification
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-business/delete/{id} [delete]
func (h *merchantBusinessHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedInvalidIdMerchantBusiness(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantBusinessRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantBusinessPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete merchant", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedDeleteMerchantBusinessPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantCertification
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-business/restore/all [post]
func (h *merchantBusinessHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllMerchantBusiness(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant restoration failed", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedRestoreAllMerchantBusiness(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully restored all merchant")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantCertification
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-business/delete/all [post]
func (h *merchantBusinessHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllMerchantBusinessPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant deletion failed", zap.Error(err))
		return merchantbusiness_errors.ErrApiFailedDeleteAllMerchantBusinessPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}
