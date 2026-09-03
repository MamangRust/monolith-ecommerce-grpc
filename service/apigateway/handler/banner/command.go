package bannerhandler

import (
	"fmt"
	"net/http"
	"strconv"

	banner_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/banner"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/banner"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bannerCommandHandlerApi struct {
	client        pb.BannerCommandServiceClient
	logger        logger.LoggerInterface
	mapper        apimapper.BannerCommandResponseMapper
	cache         banner_cache.BannerCommandCache
	observability observability.TraceLoggerObservability
}

type bannerCommandHandleDeps struct {
	client        pb.BannerCommandServiceClient
	router        *echo.Echo
	logger        logger.LoggerInterface
	mapper        apimapper.BannerCommandResponseMapper
	cache         banner_cache.BannerCommandCache
	observability observability.TraceLoggerObservability
}

func NewBannerCommandHandleApi(params *bannerCommandHandleDeps) *bannerCommandHandlerApi {
	handler := &bannerCommandHandlerApi{
		client:        params.client,
		logger:        params.logger,
		mapper:        params.mapper,
		cache:         params.cache,
		observability: params.observability,
	}

	routerBanner := params.router.Group("/api/banner-command")
	routerBanner.POST("/create", handler.Create)
	routerBanner.POST("/update/:id", handler.Update)
	routerBanner.POST("/trashed/:id", handler.Trash)
	routerBanner.POST("/restore/:id", handler.Restore)
	routerBanner.DELETE("/permanent/:id", handler.DeletePermanent)
	routerBanner.POST("/restore/all", handler.RestoreAll)
	routerBanner.POST("/permanent/all", handler.DeleteAllPermanent)

	return handler
}

// @Security Bearer
// @Summary Create a new banner
// @Tags Banner Command
// @Description Create a new banner
// @Accept json
// @Produce json
// @Param request body requests.CreateBannerRequest true "Banner details"
// @Success 200 {object} response.ApiResponseBanner "Banner created"
// @Failure 400 {object} errors.ErrorResponse "Invalid request"
// @Failure 500 {object} errors.ErrorResponse "Failed to create banner"
// @Router /api/banner-command/create [post]
func (h *bannerCommandHandlerApi) Create(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"CreateBanner",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	var body requests.CreateBannerRequest
	if err := c.Bind(&body); err != nil {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := body.Validate(); err != nil {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := h.client.Create(ctx, &pb.CreateBannerRequest{
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Create")
	}

	logSuccess("Banner created successfully")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseBanner(res))
}

// @Security Bearer
// @Summary Update a banner
// @Tags Banner Command
// @Description Update an existing banner
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Param request body requests.UpdateBannerRequest true "Banner details"
// @Success 200 {object} response.ApiResponseBanner "Banner updated"
// @Failure 400 {object} errors.ErrorResponse "Invalid request"
// @Failure 500 {object} errors.ErrorResponse "Failed to update banner"
// @Router /api/banner-command/update/{id} [post]
func (h *bannerCommandHandlerApi) Update(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"UpdateBanner",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	var body requests.UpdateBannerRequest
	if err := c.Bind(&body); err != nil {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	if err := body.Validate(); err != nil {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := h.client.Update(ctx, &pb.UpdateBannerRequest{
		BannerId:  int32(id),
		Name:      body.Name,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		EndTime:   body.EndTime,
		IsActive:  body.IsActive,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Update")
	}

	h.cache.DeleteBannerCache(ctx, id)

	logSuccess("Banner updated successfully")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseBanner(res))
}

// @Security Bearer
// @Summary Trash a banner
// @Tags Banner Command
// @Description Move a banner to trash
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Banner trashed"
// @Failure 400 {object} errors.ErrorResponse "Invalid ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to trash banner"
// @Router /api/banner-command/trashed/{id} [post]
func (h *bannerCommandHandlerApi) Trash(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"TrashBanner",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	res, err := h.client.Trash(ctx, &pb.FindByIdBannerRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Trash")
	}

	h.cache.DeleteBannerCache(ctx, id)

	logSuccess("Banner moved to trash")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseBannerDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a banner
// @Tags Banner Command
// @Description Restore a trashed banner
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDeleteAt "Banner restored"
// @Failure 400 {object} errors.ErrorResponse "Invalid ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore banner"
// @Router /api/banner-command/restore/{id} [post]
func (h *bannerCommandHandlerApi) Restore(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"RestoreBanner",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	res, err := h.client.Restore(ctx, &pb.FindByIdBannerRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Restore")
	}

	h.cache.DeleteBannerCache(ctx, id)

	logSuccess("Banner restored successfully")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseBannerDeleteAt(res))
}

// @Security Bearer
// @Summary Delete a banner permanently
// @Tags Banner Command
// @Description Permanently delete a banner
// @Accept json
// @Produce json
// @Param id path int true "Banner ID"
// @Success 200 {object} response.ApiResponseBannerDelete "Banner deleted"
// @Failure 400 {object} errors.ErrorResponse "Invalid ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete banner"
// @Router /api/banner-command/permanent/{id} [delete]
func (h *bannerCommandHandlerApi) DeletePermanent(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"DeletePermanent",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		status = "error"
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	res, err := h.client.DeletePermanent(ctx, &pb.FindByIdBannerRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Delete")
	}

	h.cache.DeleteBannerCache(ctx, id)

	logSuccess("Banner deleted permanently")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseBannerDelete(res))
}

// @Security Bearer
// @Summary Restore all banners
// @Tags Banner Command
// @Description Restore all trashed banners
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseBannerAll "All banners restored"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore banners"
// @Router /api/banner-command/restore/all [post]
func (h *bannerCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"RestoreAll",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	res, err := h.client.RestoreAll(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "RestoreAll")
	}

	logSuccess("All banners restored")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseBannerAll(res))
}

// @Security Bearer
// @Summary Delete all banners permanently
// @Tags Banner Command
// @Description Permanently delete all banners
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseBannerAll "All banners deleted"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete banners"
// @Router /api/banner-command/permanent/all [post]
func (h *bannerCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(
		c.Request().Context(),
		"DeleteAllPermanent",
		attribute.String("path", c.Request().URL.Path),
		attribute.String("method", c.Request().Method),
	)
	defer end(status)
	c.SetRequest(c.Request().WithContext(ctx))

	res, err := h.client.DeleteAll(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "DeleteAll")
	}

	logSuccess("All banners deleted permanently")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseBannerAll(res))
}

func (h *bannerCommandHandlerApi) handleError(c echo.Context, err error, span trace.Span, method string) error {
	appErr := errors.ParseGrpcError(err)
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		fmt.Sprintf("Banner command error in %s", method),
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	return errors.HandleApiError(c, appErr, traceID)
}
