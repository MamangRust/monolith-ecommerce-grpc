package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantPoliciesHandleApi struct {
	client          pb.MerchantPoliciesServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantPolicyResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
}

func NewHandlerMerchantPolicies(
	router *echo.Echo,
	client pb.MerchantPoliciesServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantPolicyResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
) *merchantPoliciesHandleApi {
	merchantPoliciesHandler := &merchantPoliciesHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
	}

	routercategory := router.Group("/api/merchant-policy")

	routercategory.GET("", merchantPoliciesHandler.FindAllMerchantPolicy)
	routercategory.GET("/:id", merchantPoliciesHandler.FindById)
	routercategory.GET("/active", merchantPoliciesHandler.FindByActive)
	routercategory.GET("/trashed", merchantPoliciesHandler.FindByTrashed)

	routercategory.POST("/create", merchantPoliciesHandler.Create)
	routercategory.POST("/update/:id", merchantPoliciesHandler.Update)

	routercategory.POST("/trashed/:id", merchantPoliciesHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantPoliciesHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantPoliciesHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantPoliciesHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantPoliciesHandler.DeleteAllMerchantPermanent)

	return merchantPoliciesHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags MerchantPolicy
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantPolicies "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-policy [get]
func (h *merchantPoliciesHandleApi) FindAllMerchantPolicy(c echo.Context) error {
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
		return merchantpolicy_errors.ErrApiFailedFindAllMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantPolicies(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantPolicy
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-policy/{id} [get]
func (h *merchantPoliciesHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch merchant details", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedFindByIdMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPolicies(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags MerchantPolicy
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantPoliciesDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-policy/active [get]
func (h *merchantPoliciesHandleApi) FindByActive(c echo.Context) error {
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
		return merchantpolicy_errors.ErrApiFailedFindByActiveMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantPoliciesDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantPolicy
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantPoliciesDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-policy/trashed [get]
func (h *merchantPoliciesHandleApi) FindByTrashed(c echo.Context) error {
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
		return merchantpolicy_errors.ErrApiFailedFindByTrashedMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantPoliciesDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new merchant policy
// @Tags MerchantPolicy
// @Description Create a new merchant policy with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantPolicyRequest true "Create merchant policy request"
// @Success 200 {object} response.ApiResponseMerchantPolicies "Successfully created merchant policy"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant policy"
// @Router /api/merchant-policy/create [post]
func (h *merchantPoliciesHandleApi) Create(c echo.Context) error {
	var body requests.CreateMerchantPolicyRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchantpolicy_errors.ErrApiBindCreateMerchantPolicy(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchantpolicy_errors.ErrValidateCreateMerchantPolicy(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateMerchantPoliciesRequest{
		MerchantId:  int32(body.MerchantID),
		PolicyType:  body.PolicyType,
		Title:       body.Title,
		Description: body.Description,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("Merchant policy creation failed", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedCreateMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPolicies(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update an existing merchant policy
// @Tags MerchantPolicy
// @Description Update an existing merchant policy with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Merchant Policy ID"
// @Param request body requests.UpdateMerchantPolicyRequest true "Update merchant policy request"
// @Success 200 {object} response.ApiResponseMerchantPolicies "Successfully updated merchant policy"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant policy"
// @Router /api/merchant-policy/update/{id} [post]
func (h *merchantPoliciesHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))
		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	var body requests.UpdateMerchantPolicyRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return merchantpolicy_errors.ErrApiBindUpdateMerchantPolicy(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return merchantpolicy_errors.ErrValidateUpdateMerchantPolicy(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateMerchantPoliciesRequest{
		MerchantPolicyId: int32(idInt),
		PolicyType:       body.PolicyType,
		Title:            body.Title,
		Description:      body.Description,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		h.logger.Error("Merchant policy update failed", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedUpdateMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPolicies(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags MerchantPolicy
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantPoliciesDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-policy/trashed/{id} [get]
func (h *merchantPoliciesHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantPolicies(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive merchant", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedTrashMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPoliciesDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantPolicy
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantPoliciesDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-policy/restore/{id} [post]
func (h *merchantPoliciesHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantPolicies(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore merchant", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedRestoreMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPoliciesDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantPolicy
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-policy/delete/{id} [delete]
func (h *merchantPoliciesHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantPoliciesPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete merchant", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedDeleteMerchantPolicy(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantPolicy
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-policy/restore/all [post]
func (h *merchantPoliciesHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllMerchantPolicies(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant restoration failed", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedRestoreAllMerchantPolicies(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully restored all merchant")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantPolicy
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-policy/delete/all [post]
func (h *merchantPoliciesHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllMerchantPoliciesPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant deletion failed", zap.Error(err))
		return merchantpolicy_errors.ErrApiFailedDeleteAllMerchantPolicies(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}
