package reviewdetailhandler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	reviewdetail_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/review_detail"
	reviewapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/review"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewDetailCommandHandleApi struct {
	client        pb.ReviewDetailCommandServiceClient
	logger        logger.LoggerInterface
	mapper        apimapper.ReviewDetailCommandResponseMapper
	queryMapper   apimapper.ReviewDetailQueryResponseMapper
	reviewMapper  reviewapimapper.ReviewCommandResponseMapper
	cache         reviewdetail_cache.ReviewDetailCommandCache
	upload        upload_image.ImageUploads
	observability observability.TraceLoggerObservability
}

type reviewDetailCommandHandleDeps struct {
	client        pb.ReviewDetailCommandServiceClient
	router        *echo.Echo
	logger        logger.LoggerInterface
	mapper        apimapper.ReviewDetailCommandResponseMapper
	queryMapper   apimapper.ReviewDetailQueryResponseMapper
	reviewMapper  reviewapimapper.ReviewCommandResponseMapper
	cache         reviewdetail_cache.ReviewDetailCommandCache
	upload        upload_image.ImageUploads
	observability observability.TraceLoggerObservability
}

func NewReviewDetailCommandHandleApi(deps *reviewDetailCommandHandleDeps) {
	handler := &reviewDetailCommandHandleApi{
		client:        deps.client,
		logger:        deps.logger,
		mapper:        deps.mapper,
		queryMapper:   deps.queryMapper,
		reviewMapper:  deps.reviewMapper,
		cache:         deps.cache,
		upload:        deps.upload,
		observability: deps.observability,
	}

	router := deps.router.Group("/api/review-detail-command")
	router.POST("/create", handler.Create)
	router.POST("/update/:id", handler.Update)
	router.POST("/trashed/:id", handler.TrashedReviewDetail)
	router.POST("/restore/:id", handler.RestoreReviewDetail)
	router.DELETE("/permanent/:id", handler.DeleteReviewDetailPermanent)
	router.POST("/restore/all", handler.RestoreAllReviewDetail)
	router.POST("/permanent/all", handler.DeleteAllReviewDetailPermanent)
}

// @Security Bearer
// @Summary Create a new review detail
// @Tags Review Detail Command
// @Description Create a new review detail (e.g., photo/video attachment)
// @Accept mpfd
// @Produce json
// @Param review_id formData int true "Review ID"
// @Param type formData string true "Attachment type (e.g., image, video)"
// @Param caption formData string true "Attachment caption"
// @Param url formData file true "Attachment file"
// @Success 201 {object} response.ApiResponseReviewDetail "Successfully created review detail"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create review detail"
// @Router /api/review-detail-command/create [post]
func (h *reviewDetailCommandHandleApi) Create(c echo.Context) error {
	formData, err := h.parseReviewDetailForm(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	grpcReq := &pb.CreateReviewDetailRequest{
		ReviewId: int32(formData.ReviewID),
		Type:     formData.Type,
		Url:      formData.Url,
		Caption:  formData.Caption,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "CreateReviewDetail")
	defer func() {
		end(status)
	}()

	res, err := h.client.Create(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Create")
	}

	logSuccess("Successfully created review detail")
	return c.JSON(http.StatusCreated, h.queryMapper.ToApiResponseReviewDetail(res))
}

// @Security Bearer
// @Summary Update review detail
// @Tags Review Detail Command
// @Description Update an existing review detail
// @Accept mpfd
// @Produce json
// @Param id path int true "Review Detail ID"
// @Param type formData string true "Attachment type"
// @Param caption formData string true "Attachment caption"
// @Param url formData file false "New attachment file"
// @Success 200 {object} response.ApiResponseReviewDetail "Successfully updated review detail"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update review detail"
// @Router /api/review-detail-command/update/{id} [post]
func (h *reviewDetailCommandHandleApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	formData, err := h.parseReviewDetailForm(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	grpcReq := &pb.UpdateReviewDetailRequest{
		ReviewDetailId: int32(id),
		Type:           formData.Type,
		Url:            formData.Url,
		Caption:        formData.Caption,
	}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "UpdateReviewDetail")
	defer func() {
		end(status)
	}()

	res, err := h.client.Update(ctx, grpcReq)
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Update")
	}

	h.cache.DeleteReviewDetailCache(ctx, id)

	logSuccess("Successfully updated review detail")
	return c.JSON(http.StatusOK, h.queryMapper.ToApiResponseReviewDetail(res))
}

// @Security Bearer
// @Summary Move review detail to trash
// @Tags Review Detail Command
// @Description Move a review detail record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Review Detail ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully moved review detail to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid review detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move review detail to trash"
// @Router /api/review-detail-command/trashed/{id} [post]
func (h *reviewDetailCommandHandleApi) TrashedReviewDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "TrashedReviewDetail")
	defer func() {
		end(status)
	}()

	res, err := h.client.TrashedReviewDetail(ctx, &pb.FindByIdReviewDetailRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Trash")
	}

	h.cache.DeleteReviewDetailCache(ctx, id)

	logSuccess("Successfully moved review detail to trash")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseReviewDetailDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a trashed review detail
// @Tags Review Detail Command
// @Description Restore a trashed review detail record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Review Detail ID"
// @Success 200 {object} response.ApiResponseReviewDetailDeleteAt "Successfully restored review detail"
// @Failure 400 {object} errors.ErrorResponse "Invalid review detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore review detail"
// @Router /api/review-detail-command/restore/{id} [post]
func (h *reviewDetailCommandHandleApi) RestoreReviewDetail(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "RestoreReviewDetail")
	defer func() {
		end(status)
	}()

	res, err := h.client.RestoreReviewDetail(ctx, &pb.FindByIdReviewDetailRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Restore")
	}

	h.cache.DeleteReviewDetailCache(ctx, id)

	logSuccess("Successfully restored review detail")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseReviewDetailDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete a review detail
// @Tags Review Detail Command
// @Description Permanently delete a review detail record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Review Detail ID"
// @Success 200 {object} response.ApiResponseReviewDelete "Successfully deleted review detail record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid review detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete review detail permanently"
// @Router /api/review-detail-command/permanent/{id} [delete]
func (h *reviewDetailCommandHandleApi) DeleteReviewDetailPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "DeleteReviewDetailPermanent")
	defer func() {
		end(status)
	}()

	res, err := h.client.DeleteReviewDetailPermanent(ctx, &pb.FindByIdReviewDetailRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Delete")
	}

	h.cache.DeleteReviewDetailCache(ctx, id)

	logSuccess("Successfully deleted review detail permanently")
	return c.JSON(http.StatusOK, h.reviewMapper.ToApiResponseReviewDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed review details
// @Tags Review Detail Command
// @Description Restore all trashed review detail records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseReviewAll "Successfully restored all review details"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore review details"
// @Router /api/review-detail-command/restore/all [post]
func (h *reviewDetailCommandHandleApi) RestoreAllReviewDetail(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "RestoreAllReviewDetail")
	defer func() {
		end(status)
	}()

	res, err := h.client.RestoreAllReviewDetail(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "RestoreAll")
	}

	logSuccess("Successfully restored all review details")
	return c.JSON(http.StatusOK, h.reviewMapper.ToApiResponseReviewAll(res))
}

func (h *reviewDetailCommandHandleApi) DeleteAllReviewDetailPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "DeleteAllReviewDetailPermanent")
	defer func() {
		end(status)
	}()

	res, err := h.client.DeleteAllReviewDetailPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "DeleteAll")
	}

	logSuccess("Successfully deleted all review details permanently")
	return c.JSON(http.StatusOK, h.reviewMapper.ToApiResponseReviewAll(res))
}

func (h *reviewDetailCommandHandleApi) parseReviewDetailForm(c echo.Context) (requests.ReviewDetailFormData, error) {
	var formData requests.ReviewDetailFormData
	var err error

	reviewIDStr := c.FormValue("review_id")
	if reviewIDStr != "" {
		formData.ReviewID, err = strconv.Atoi(reviewIDStr)
		if err != nil || formData.ReviewID <= 0 {
			return formData, echo.NewHTTPError(http.StatusBadRequest, "Invalid Review ID")
		}
	}

	formData.Type = strings.TrimSpace(c.FormValue("type"))
	if formData.Type == "" {
		return formData, echo.NewHTTPError(http.StatusBadRequest, "Type is required")
	}

	formData.Caption = strings.TrimSpace(c.FormValue("caption"))
	if formData.Caption == "" {
		return formData, echo.NewHTTPError(http.StatusBadRequest, "Caption is required")
	}

	file, err := c.FormFile("url")
	if err == nil {
		uploadPath, err := h.upload.ProcessImageUpload(c, "uploads/review_detail", file, false)
		if err != nil {
			return formData, err
		}
		formData.Url = uploadPath
	} else if c.FormValue("url") != "" {
		formData.Url = c.FormValue("url")
	}

	return formData, nil
}

func (h *reviewDetailCommandHandleApi) handleError(c echo.Context, err error, span trace.Span, method string) error {
	appErr := errors.ParseGrpcError(err)
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		"Review detail command error in "+method,
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	return errors.HandleApiError(c, appErr, traceID)
}
