package merchantpolicyhandler

import (
	"net/http"
	"strconv"

	merchantpolicy_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_policies"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_policy"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type merchantPolicyQueryHandlerApi struct {
	client        pb.MerchantPolicyQueryServiceClient
	logger        logger.LoggerInterface
	mapper        apimapper.MerchantPolicyQueryResponseMapper
	cache         merchantpolicy_cache.MerchantPolicyQueryCache
	observability observability.TraceLoggerObservability
}

type merchantPolicyQueryHandleDeps struct {
	client        pb.MerchantPolicyQueryServiceClient
	router        *echo.Echo
	logger        logger.LoggerInterface
	mapper        apimapper.MerchantPolicyQueryResponseMapper
	cache         merchantpolicy_cache.MerchantPolicyQueryCache
	observability observability.TraceLoggerObservability
}

func NewMerchantPolicyQueryHandleApi(params *merchantPolicyQueryHandleDeps) *merchantPolicyQueryHandlerApi {
	handler := &merchantPolicyQueryHandlerApi{
		client:        params.client,
		logger:        params.logger,
		mapper:        params.mapper,
		cache:         params.cache,
		observability: params.observability,
	}

	routerPolicy := params.router.Group("/api/merchant-policy-query")
	routerPolicy.GET("", handler.FindAll)
	routerPolicy.GET("/:id", handler.FindById)
	routerPolicy.GET("/active", handler.FindByActive)
	routerPolicy.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all merchant policies
// @Tags Merchant Policy Query
// @Description Retrieve a list of all merchant policies
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantPolicies "List of merchant policies"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant policy data"
// @Router /api/merchant-policy-query [get]
func (h *merchantPolicyQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindAllMerchantPolicies")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindAll(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindAll")
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantPolicies(res)
	h.cache.SetCachedMerchantPolicyAll(ctx, req, apiResponse)

	logSuccess("Successfully fetched all merchant policies")
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Find merchant policy by ID
// @Tags Merchant Policy Query
// @Description Retrieve a merchant policy by ID
// @Accept json
// @Produce json
// @Param id path int true "Policy ID"
// @Success 200 {object} response.ApiResponseMerchantPolicies "Merchant policy data"
// @Failure 400 {object} errors.ErrorResponse "Invalid policy ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant policy data"
// @Router /api/merchant-policy-query/{id} [get]
func (h *merchantPolicyQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByIdMerchantPolicy")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindById(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindById")
	}

	apiResponse := h.mapper.ToApiResponseMerchantPolicies(res)
	h.cache.SetCachedMerchantPolicy(ctx, apiResponse)

	logSuccess("Successfully fetched merchant policy by ID")
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve active merchant policies
// @Tags Merchant Policy Query
// @Description Retrieve a list of active merchant policies
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantPoliciesDeleteAt "List of active merchant policies"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant policy data"
// @Router /api/merchant-policy-query/active [get]
func (h *merchantPolicyQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByActiveMerchantPolicies")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByActive(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByActive")
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantPoliciesDeleteAt(res)
	h.cache.SetCachedMerchantPolicyActive(ctx, req, apiResponse)

	logSuccess("Successfully fetched active merchant policies")
	return c.JSON(http.StatusOK, apiResponse)
}

// @Security Bearer
// @Summary Retrieve trashed merchant policies
// @Tags Merchant Policy Query
// @Description Retrieve a list of trashed merchant policy records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantPoliciesDeleteAt "List of trashed merchant policy data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant policy data"
// @Router /api/merchant-policy-query/trashed [get]
func (h *merchantPolicyQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	req := &requests.FindAllMerchant{Page: page, PageSize: pageSize, Search: search}

	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "FindByTrashedMerchantPolicies")
	defer func() {
		end(status)
	}()

	res, err := h.client.FindByTrashed(ctx, &pb.FindAllMerchantRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "FindByTrashed")
	}

	apiResponse := h.mapper.ToApiResponsePaginationMerchantPoliciesDeleteAt(res)
	h.cache.SetCachedMerchantPolicyTrashed(ctx, req, apiResponse)

	logSuccess("Successfully fetched trashed merchant policies")
	return c.JSON(http.StatusOK, apiResponse)
}

func (h *merchantPolicyQueryHandlerApi) handleError(c echo.Context, err error, span trace.Span, method string) error {
	appErr := errors.ParseGrpcError(err)
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		"Merchant policy query error in "+method,
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	return errors.HandleApiError(c, appErr, traceID)
}
