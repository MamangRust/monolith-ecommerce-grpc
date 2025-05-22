package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewDetailHandleApi struct {
	client        pb.ReviewDetailServiceClient
	logger        logger.LoggerInterface
	mapping       response_api.ReviewDetailResponseMapper
	mappingReview response_api.ReviewResponseMapper
	upload_image  upload_image.ImageUploads
}

func NewHandlerReviewDetail(
	router *echo.Echo,
	client pb.ReviewDetailServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ReviewDetailResponseMapper,
	mappingReview response_api.ReviewResponseMapper,
	upload_image upload_image.ImageUploads,
) *reviewDetailHandleApi {
	reviewDetailHandler := &reviewDetailHandleApi{
		client:        client,
		logger:        logger,
		mapping:       mapping,
		mappingReview: mappingReview,
		upload_image:  upload_image,
	}

	routercategory := router.Group("/api/review-detail")

	routercategory.GET("", reviewDetailHandler.FindAllReviewDetail)
	routercategory.GET("/:id", reviewDetailHandler.FindById)
	routercategory.GET("/active", reviewDetailHandler.FindByActive)
	routercategory.GET("/trashed", reviewDetailHandler.FindByTrashed)

	routercategory.POST("/create", reviewDetailHandler.Create)
	routercategory.POST("/update/:id", reviewDetailHandler.Update)

	routercategory.POST("/trashed/:id", reviewDetailHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", reviewDetailHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", reviewDetailHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", reviewDetailHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", reviewDetailHandler.DeleteAllMerchantPermanent)

	return reviewDetailHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags ReviewDetail
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetails "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail [get]
func (h *reviewDetailHandleApi) FindAllReviewDetail(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch merchants", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedFindAllReviewDetails(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDetail(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags ReviewDetail
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetail "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/{id} [get]
func (h *reviewDetailHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch merchant details", zap.Error(err))
		return reviewdetail_errors.ErrApiReviewDetailNotFound(c)
	}

	so := h.mapping.ToApiResponseReviewDetail(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags ReviewDetail
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/active [get]
func (h *reviewDetailHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active merchants", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedFindActiveReviewDetails(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags ReviewDetail
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationReviewDetailsDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/review-detail/trashed [get]
func (h *reviewDetailHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived merchants", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedFindTrashedReviewDetails(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new merchant review detail.
// @Summary Create a new merchant review detail
// @Tags ReviewDetail
// @Description Create a new merchant review detail with the provided details
// @Accept multipart/form-data
// @Produce json
// @Param type formData string true "Type"
// @Param url formData file true "url"
// @Param caption formData string true "Product name"
// @Success 200 {object} response.ApiResponseReviewDetail "Successfully created review detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create review detail"
// @Router /api/review-detail/create [post]
func (h *reviewDetailHandleApi) Create(c echo.Context) error {
	formData, err := h.parseReviewDetailForm(c)
	if err != nil {
		return reviewdetail_errors.ErrApiInvalidBody(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateReviewDetailRequest{
		ReviewId: int32(formData.ReviewID),
		Type:     formData.Type,
		Url:      formData.Url,
		Caption:  formData.Caption,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("Review detail creation failed", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedCreateReviewDetail(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponseReviewDetail(res))
}

// @Security Bearer
// Update handles the update of an existing merchant review detail.
// @Summary Update an existing merchant review detail
// @Tags ReviewDetail
// @Description Update an existing merchant review detail with the provided details
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Review Detail ID"
// @Param type formData string true "Type"
// @Param url formData file true "url"
// @Param caption formData string true "Product name"
// @Param request body requests.UpdateReviewDetailRequest true "Update review detail request"
// @Success 200 {object} response.ApiResponseReviewDetail "Successfully updated review detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update review detail"
// @Router /api/review-detail/update/{id} [post]
func (h *reviewDetailHandleApi) Update(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))
		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	formData, err := h.parseReviewDetailForm(c)
	if err != nil {
		return reviewdetail_errors.ErrApiInvalidBody(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateReviewDetailRequest{
		ReviewDetailId: int32(idInt),
		Type:           formData.Type,
		Url:            formData.Url,
		Caption:        formData.Caption,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		h.logger.Error("Review detail update failed", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedUpdateReviewDetail(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponseReviewDetail(res))
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags ReviewDetail
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/review-detail/trashed/{id} [get]
func (h *reviewDetailHandleApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedReviewDetail(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive merchant", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedTrashedReviewDetail(c)
	}

	so := h.mapping.ToApiResponseReviewDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags ReviewDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/review-detail/restore/{id} [post]
func (h *reviewDetailHandleApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreReviewDetail(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore merchant", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedRestoreReviewDetail(c)
	}

	so := h.mapping.ToApiResponseReviewDetailDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags ReviewDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/review-detail/delete/{id} [delete]
func (h *reviewDetailHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return reviewdetail_errors.ErrApiInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteReviewDetailPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to delete merchant", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedDeleteReviewDetailPermanent(c)
	}

	so := h.mappingReview.ToApiResponseReviewDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags ReviewDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/review-detail/restore/all [post]
func (h *reviewDetailHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllReviewDetail(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant restoration failed", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedRestoreAllReviewDetail(c)
	}

	so := h.mappingReview.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully restored all merchant")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags ReviewDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/review-detail/delete/all [post]
func (h *reviewDetailHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllReviewDetailPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk merchant deletion failed", zap.Error(err))
		return reviewdetail_errors.ErrApiFailedDeleteAllReviewDetailPermanent(c)
	}

	so := h.mappingReview.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully deleted all merchant permanently")

	return c.JSON(http.StatusOK, so)
}

func (h *reviewDetailHandleApi) parseReviewDetailForm(c echo.Context) (requests.ReviewDetailFormData, error) {
	var formData requests.ReviewDetailFormData
	var err error

	formData.ReviewID, err = strconv.Atoi(c.FormValue("review_id"))
	if err != nil || formData.ReviewID <= 0 {
		return formData, reviewdetail_errors.ErrApiInvalidReviewId(c)
	}

	formData.Type = strings.TrimSpace(c.FormValue("type"))
	if formData.Type == "" {
		return formData, reviewdetail_errors.ErrApiReviewDetailTypeRequired(c)
	}

	formData.Caption = strings.TrimSpace(c.FormValue("caption"))
	if formData.Caption == "" {
		return formData, reviewdetail_errors.ErrApiReviewDetailCaptionRequired(c)
	}

	file, err := c.FormFile("url")
	if err != nil {
		return formData, reviewdetail_errors.ErrApiReviewDetailFileRequired(c)
	}

	uploadPath, err := h.upload_image.ProcessImageUpload(c, file)
	if err != nil {
		return formData, err
	}
	formData.Url = uploadPath

	return formData, nil
}
