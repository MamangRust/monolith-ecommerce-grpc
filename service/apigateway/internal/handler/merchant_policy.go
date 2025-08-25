package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantPoliciesHandleApi struct {
	client          pb.MerchantPoliciesServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantPolicyResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchantPolicies(
	router *echo.Echo,
	client pb.MerchantPoliciesServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantPolicyResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
) *merchantPoliciesHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_policy_handler_requests_total",
			Help: "Total number of banner requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_policy_handler_request_duration_seconds",
			Help:    "Duration of banner requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantPoliciesHandler := &merchantPoliciesHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
		trace:           otel.Tracer("merchant-policy-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
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
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllMerchantPolicy"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedFindAllMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantPolicies(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("Success", true))

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
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to fetch merchant details", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to fetch merchant details", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedFindByIdMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPolicies(res)

	logSuccess("Successfully fetch merchant details", zap.Bool("Success", true))

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
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActive"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedFindByActiveMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantPoliciesDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("Success", true))

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
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashed"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	status := "success"

	defer func() { end(status) }()

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedFindByTrashedMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantPoliciesDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("Success", true))

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
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	var body requests.CreateMerchantPolicyRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind request body", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiBindCreateMerchantPolicy(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		return merchantpolicy_errors.ErrValidateCreateMerchantPolicy(c)
	}

	req := &pb.CreateMerchantPoliciesRequest{
		MerchantId:  int32(body.MerchantID),
		PolicyType:  body.PolicyType,
		Title:       body.Title,
		Description: body.Description,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to create merchant policy", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedCreateMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPolicies(res)

	logSuccess("Successfully create merchant policy", zap.Bool("Success", true))

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
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		status = "error"

		logError("Failed to parse merchant policy id", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	var body requests.UpdateMerchantPolicyRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("Failed to bind request body", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiBindUpdateMerchantPolicy(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("Failed to validate request body", err, zap.Error(err))

		return merchantpolicy_errors.ErrValidateUpdateMerchantPolicy(c)
	}

	req := &pb.UpdateMerchantPoliciesRequest{
		MerchantPolicyId: int32(idInt),
		PolicyType:       body.PolicyType,
		Title:            body.Title,
		Description:      body.Description,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to update merchant policy", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedUpdateMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPolicies(res)

	logSuccess("Successfully update merchant policy", zap.Bool("Success", true))

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
	const method = "TrashedMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to parse merchant policy id", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantPolicies(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to trashed merchant policy", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedTrashMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPoliciesDeleteAt(res)

	logSuccess("Successfully trashed merchant policy", zap.Bool("Success", true))

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
	const method = "RestoreMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to parse merchant policy id", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantPolicies(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore merchant policy", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedRestoreMerchantPolicy(c)
	}

	so := h.mapping.ToApiResponseMerchantPoliciesDeleteAt(res)

	logSuccess("Successfully restore merchant policy", zap.Bool("Success", true))

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
	const method = "DeleteMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Failed to parse merchant policy id", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantPoliciesRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantPoliciesPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to delete merchant policy", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedDeleteMerchantPolicy(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	logSuccess("Successfully deleted merchant policy", zap.Bool("Success", true))

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
	const method = "RestoreAllMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllMerchantPolicies(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all merchant", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedRestoreAllMerchantPolicies(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	logSuccess("Successfully restore all merchant", zap.Bool("Success", true))

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
	const method = "DeleteAllMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllMerchantPoliciesPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Bulk merchant deletion failed", err, zap.Error(err))

		return merchantpolicy_errors.ErrApiFailedDeleteAllMerchantPolicies(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	logSuccess("Bulk merchant deletion success", zap.Bool("Success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *merchantPoliciesHandleApi) startTracingAndLogging(
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

func (s *merchantPoliciesHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
