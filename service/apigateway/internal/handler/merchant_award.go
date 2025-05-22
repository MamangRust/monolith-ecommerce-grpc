package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantAwardHandleApi struct {
	client          pb.MerchantAwardServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantAwardResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
}

func NewHandlerMerchantAward(
	router *echo.Echo,
	client pb.MerchantAwardServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantAwardResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
) *merchantAwardHandleApi {
	merchantAwardHandler := &merchantAwardHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
	}

	routercategory := router.Group("/api/merchant-certification")

	routercategory.GET("", merchantAwardHandler.FindAllMerchantAward)
	routercategory.GET("/:id", merchantAwardHandler.FindById)
	routercategory.GET("/active", merchantAwardHandler.FindByActive)
	routercategory.GET("/trashed", merchantAwardHandler.FindByTrashed)

	routercategory.POST("/create", merchantAwardHandler.Create)
	routercategory.POST("/update/:id", merchantAwardHandler.Update)

	routercategory.POST("/trashed/:id", merchantAwardHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantAwardHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantAwardHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantAwardHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantAwardHandler.DeleteAllMerchantPermanent)

	return merchantAwardHandler
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
// @Success 200 {object} response.ApiResponsePaginationMerchantAward "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification [get]
func (h *merchantAwardHandleApi) FindAllMerchantAward(c echo.Context) error {
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
		return merchantaward_errors.ErrApiFailedFindAllMerchantAward(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantAward(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantCertification
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAward "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification/{id} [get]
func (h *merchantAwardHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch merchant details", zap.Error(err))
		return merchantaward_errors.ErrApiFailedFindMerchantAwardById(c)
	}

	so := h.mapping.ToApiResponseMerchantAward(res)

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
// @Success 200 {object} response.ApiResponsePaginationMerchantAwardDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification/active [get]
func (h *merchantAwardHandleApi) FindByActive(c echo.Context) error {
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
		return merchantaward_errors.ErrApiFailedFindActiveMerchantAward(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantAwardDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantAwardDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-certification/trashed [get]
func (h *merchantAwardHandleApi) FindByTrashed(c echo.Context) error {
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
		return merchantaward_errors.ErrApiFailedFindTrashedMerchantAward(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantAwardDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new merchant certification or award.
// @Summary Create a new merchant certification or award
// @Tags MerchantCertificationCertification
// @Description Create a new merchant certification or award with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantCertificationOrAwardRequest true "Create merchant certification or award request"
// @Success 200 {object} response.ApiResponseMerchantAward "Successfully created merchant certification or award"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant certification or award"
// @Router /api/merchant-certification/create [post]
func (h *merchantAwardHandleApi) Create(c echo.Context) error {
	var body requests.CreateMerchantCertificationOrAwardRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchantaward_errors.ErrApiBindCreateMerchantAward(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchantaward_errors.ErrApiValidateCreateMerchantAward(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateMerchantAwardRequest{
		MerchantId:     int32(body.MerchantID),
		Title:          body.Title,
		Description:    body.Description,
		IssuedBy:       body.IssuedBy,
		IssueDate:      body.IssueDate,
		ExpiryDate:     body.ExpiryDate,
		CertificateUrl: body.CertificateUrl,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("Merchant award creation failed", zap.Error(err))
		return merchantaward_errors.ErrApiFailedCreateMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAward(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Update handles the update of an existing merchant certification or award.
// @Summary Update an existing merchant certification or award
// @Tags MerchantCertification
// @Description Update an existing merchant certification or award with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Merchant Certification ID"
// @Param request body requests.UpdateMerchantCertificationOrAwardRequest true "Update merchant certification or award request"
// @Success 200 {object} response.ApiResponseMerchantAward "Successfully updated merchant certification or award"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant certification or award"
// @Router /api/merchant-certification/update/{id} [post]
func (h *merchantAwardHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))
		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	var body requests.UpdateMerchantCertificationOrAwardRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchantaward_errors.ErrApiBindUpdateMerchantAward(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchantaward_errors.ErrApiValidateUpdateMerchantAward(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateMerchantAwardRequest{
		MerchantCertificationId: int32(idInt),
		Title:                   body.Title,
		Description:             body.Description,
		IssuedBy:                body.IssuedBy,
		IssueDate:               body.IssueDate,
		ExpiryDate:              body.ExpiryDate,
		CertificateUrl:          body.CertificateUrl,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		h.logger.Error("Merchant award update failed", zap.Error(err))
		return merchantaward_errors.ErrApiFailedUpdateMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAward(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags MerchantCertification
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAwardDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-certification/trashed/{id} [get]
func (h *merchantAwardHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantAward(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive merchant", zap.Error(err))
		return merchantaward_errors.ErrApiFailedTrashedMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAwardDeleteAt(res)

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
// @Success 200 {object} response.ApiResponseMerchantAwardDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-certification/restore/{id} [post]
func (h *merchantAwardHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantAward(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore merchant", zap.Error(err))
		return merchantaward_errors.ErrApiFailedRestoreMerchantAward(c)
	}

	so := h.mapping.ToApiResponseMerchantAwardDeleteAt(res)

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
// @Router /api/merchant-certification/delete/{id} [delete]
func (h *merchantAwardHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantaward_errors.ErrApiFailedInvalidMerchantAwardId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantAwardRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantAwardPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete merchant", zap.Error(err))
		return merchantaward_errors.ErrApiFailedDeleteMerchantAwardPermanent(c)
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
// @Router /api/merchant-certification/restore/all [post]
func (h *merchantAwardHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllMerchantAward(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant restoration failed", zap.Error(err))
		return merchantaward_errors.ErrApiFailedRestoreAllMerchantAward(c)
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
// @Router /api/merchant-certification/delete/all [post]
func (h *merchantAwardHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllMerchantAwardPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant deletion failed", zap.Error(err))
		return merchantaward_errors.ErrApiFailedDeleteAllMerchantAwardPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}
