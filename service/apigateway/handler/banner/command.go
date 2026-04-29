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
		return h.handleError(c, errors.NewBadRequestError("Invalid request"), span, "Create")
	}
	if err := body.Validate(); err != nil {
		status = "error"
		return h.handleError(c, errors.NewBadRequestError(err.Error()), span, "Create")
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
		return h.handleError(c, errors.NewBadRequestError("Invalid ID"), span, "Update")
	}

	var body requests.UpdateBannerRequest
	if err := c.Bind(&body); err != nil {
		status = "error"
		return h.handleError(c, errors.NewBadRequestError("Invalid request"), span, "Update")
	}
	if err := body.Validate(); err != nil {
		status = "error"
		return h.handleError(c, errors.NewBadRequestError(err.Error()), span, "Update")
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
		return h.handleError(c, errors.NewBadRequestError("Invalid ID"), span, "Trash")
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
		return h.handleError(c, errors.NewBadRequestError("Invalid ID"), span, "Restore")
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
		return h.handleError(c, errors.NewBadRequestError("Invalid ID"), span, "Delete")
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

