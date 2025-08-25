package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantsociallink_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantSocialLinkHandleApi struct {
	client          pb.MerchantSocialServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantSocialLinkMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchantSocialLink(
	router *echo.Echo,
	client pb.MerchantSocialServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantSocialLinkMapper,
) *merchantSocialLinkHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_social_link_handler_requests_total",
			Help: "Total number of banner requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_social_link_handler_request_duration_seconds",
			Help:    "Duration of banner requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantSocialLinkHandler := &merchantSocialLinkHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("merchant-social-link-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/merchant-social-link")

	routercategory.POST("/create", merchantSocialLinkHandler.CreateMerchantSocial)
	routercategory.POST("/update/:id", merchantSocialLinkHandler.UpdateMerchantSocial)

	return merchantSocialLinkHandler
}

// @Security Bearer
// @Summary Create a new merchant social link
// @Tags MerchantSocial
// @Description Create a new social media link for a merchant
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantSocialRequest true "Create merchant social link request"
// @Success 200 {object} response.ApiResponseMerchantSocial "Successfully created merchant social link"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create merchant social link"
// @Router /api/merchant-social-link/create [post]
func (h *merchantSocialLinkHandleApi) CreateMerchantSocial(c echo.Context) error {
	const method = "CreateMerchantSocial"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"
	defer func() { end(status) }()

	var body requests.CreateMerchantSocialRequest
	if err := c.Bind(&body); err != nil {
		status = "error"
		logError("Failed to bind request body", err, zap.Error(err))
		return merchantsociallink_errors.ErrApiBindCreateMerchantSocialLink(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"
		logError("Validation failed for request body", err, zap.Error(err))
		return merchantsociallink_errors.ErrApiValidateCreateMerchantSocialLink(c)
	}

	req := &pb.CreateMerchantSocialRequest{
		MerchantDetailId: int32(*body.MerchantDetailID),
		Platform:         body.Platform,
		Url:              body.Url,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		status = "error"
		logError("Failed to create merchant social link", err, zap.Error(err))
		return merchantsociallink_errors.ErrApiBindCreateMerchantSocialLink(c)
	}

	so := h.mapping.ToApiResponseMerchantSocialLink(res)

	logSuccess("Successfully created merchant social link", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update an existing merchant social link
// @Tags MerchantSocial
// @Description Update an existing merchant social media link
// @Accept json
// @Produce json
// @Param id path int true "Merchant Social Link ID"
// @Param request body requests.UpdateMerchantSocialRequest true "Update merchant social link request"
// @Success 200 {object} response.ApiResponseMerchantSocial "Successfully updated merchant social link"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update merchant social link"
// @Router /api/merchant-social/update/{id} [post]
func (h *merchantSocialLinkHandleApi) UpdateMerchantSocial(c echo.Context) error {
	const method = "UpdateMerchantSocial"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)
	status := "success"
	defer func() { end(status) }()

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		status = "error"
		logError("Invalid merchant social link ID", err, zap.Error(err))
		return merchantsociallink_errors.ErrApiBindCreateMerchantSocialLink(c)
	}

	var body requests.UpdateMerchantSocialRequest
	if err := c.Bind(&body); err != nil {
		status = "error"
		logError("Failed to bind request body", err, zap.Error(err))
		return merchantsociallink_errors.ErrApiBindUpdateMerchantSocialLink(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"
		logError("Validation failed for request body", err, zap.Error(err))
		return merchantsociallink_errors.ErrApiValidateUpdateMerchantSocialLink(c)
	}

	body.ID = idInt

	req := &pb.UpdateMerchantSocialRequest{
		Id:               int32(body.ID),
		MerchantDetailId: int32(*body.MerchantDetailID),
		Platform:         body.Platform,
		Url:              body.Url,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		status = "error"
		logError("Failed to update merchant social link", err, zap.Error(err))
		return merchantsociallink_errors.ErrApiFailedUpdateMerchantSocialLink(c)
	}

	so := h.mapping.ToApiResponseMerchantSocialLink(res)
	logSuccess("Successfully updated merchant social link", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantSocialLinkHandleApi) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (func(string), func(string, ...zap.Field), func(string, error, ...zap.Field)) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError := func(msg string, err error, fields ...zap.Field) {
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *merchantSocialLinkHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
