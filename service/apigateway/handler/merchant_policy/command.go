package merchantpolicyhandler

import (
	"net/http"
	"strconv"

	merchantpolicy_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_policies"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_policy"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantPolicyCommandHandlerApi struct {
	client         pb.MerchantPolicyCommandServiceClient
	logger         logger.LoggerInterface
	mapper         apimapper.MerchantPolicyCommandResponseMapper
	merchantMapper merchantapimapper.MerchantCommandResponseMapper
	cache          merchantpolicy_cache.MerchantPolicyCommandCache
	observability  observability.TraceLoggerObservability
}

type merchantPolicyCommandHandleDeps struct {
	client         pb.MerchantPolicyCommandServiceClient
	router         *echo.Echo
	logger         logger.LoggerInterface
	mapper         apimapper.MerchantPolicyCommandResponseMapper
	merchantMapper merchantapimapper.MerchantCommandResponseMapper
	cache          merchantpolicy_cache.MerchantPolicyCommandCache
	observability  observability.TraceLoggerObservability
}

func NewMerchantPolicyCommandHandleApi(params *merchantPolicyCommandHandleDeps) *merchantPolicyCommandHandlerApi {
	handler := &merchantPolicyCommandHandlerApi{
		client:         params.client,
		logger:         params.logger,
		mapper:         params.mapper,
		merchantMapper: params.merchantMapper,
		cache:          params.cache,
		observability:  params.observability,
	}

	routerPolicy := params.router.Group("/api/merchant-policy-command")
	routerPolicy.POST("/create", handler.Create)
	routerPolicy.POST("/update/:id", handler.Update)
	routerPolicy.POST("/trashed/:id", handler.Trash)
	routerPolicy.POST("/restore/:id", handler.Restore)
	routerPolicy.DELETE("/permanent/:id", handler.DeletePermanent)
	routerPolicy.POST("/restore/all", handler.RestoreAll)
	routerPolicy.POST("/permanent/all", handler.DeleteAllPermanent)

	return handler
}

// @Security Bearer
// @Summary Create merchant policy
// @Tags Merchant Policy Command
// @Description Create a new policy for a merchant
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantPolicyRequest true "Create merchant policy request"
// @Success 200 {object} response.ApiResponseMerchantPolicies "Successfully created merchant policy"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create merchant policy"
// @Router /api/merchant-policy-command/create [post]
func (h *merchantPolicyCommandHandlerApi) Create(c echo.Context) error {
	var body requests.CreateMerchantPolicyRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "CreateMerchantPolicy")
	defer func() {
		end(status)
	}()

	res, err := h.client.Create(ctx, &pb.CreateMerchantPoliciesRequest{
		MerchantId:  int32(body.MerchantID),
		PolicyType:  body.PolicyType,
		Title:       body.Title,
		Description: body.Description,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Create")
	}

	logSuccess("Successfully created merchant policy")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantPolicies(res))
}

// @Security Bearer
// @Summary Update merchant policy
// @Tags Merchant Policy Command
// @Description Update an existing policy for a merchant
// @Accept json
// @Produce json
// @Param id path int true "Policy ID"
// @Param body body requests.UpdateMerchantPolicyRequest true "Update merchant policy request"
// @Success 200 {object} response.ApiResponseMerchantPolicies "Successfully updated merchant policy"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update merchant policy"
// @Router /api/merchant-policy-command/update/{id} [post]
func (h *merchantPolicyCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	var body requests.UpdateMerchantPolicyRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	body.MerchantPolicyID = &id
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "UpdateMerchantPolicy")
	defer func() {
		end(status)
	}()

	res, err := h.client.Update(ctx, &pb.UpdateMerchantPoliciesRequest{
		MerchantPolicyId: int32(id),
		PolicyType:       body.PolicyType,
		Title:            body.Title,
		Description:      body.Description,
	})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Update")
	}

	h.cache.DeleteMerchantPolicyCache(ctx, 0) // Body has no MerchantID

	logSuccess("Successfully updated merchant policy")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantPolicies(res))
}

// @Security Bearer
// @Summary Move merchant policy to trash
// @Tags Merchant Policy Command
// @Description Move a merchant policy record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Policy ID"
// @Success 200 {object} response.ApiResponseMerchantPoliciesDeleteAt "Successfully moved merchant policy to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid policy ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move merchant policy to trash"
// @Router /api/merchant-policy-command/trashed/{id} [post]
func (h *merchantPolicyCommandHandlerApi) Trash(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "TrashedMerchantPolicy")
	defer func() {
		end(status)
	}()

	res, err := h.client.TrashedMerchantPolicies(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Trash")
	}

	h.cache.DeleteMerchantPolicyCache(ctx, 0)

	logSuccess("Successfully moved merchant policy to trash")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantPoliciesDeleteAt(res))
}

// @Security Bearer
// @Summary Restore trashed merchant policy
// @Tags Merchant Policy Command
// @Description Restore a trashed merchant policy record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Policy ID"
// @Success 200 {object} response.ApiResponseMerchantPoliciesDeleteAt "Successfully restored merchant policy"
// @Failure 400 {object} errors.ErrorResponse "Invalid policy ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant policy"
// @Router /api/merchant-policy-command/restore/{id} [post]
func (h *merchantPolicyCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "RestoreMerchantPolicy")
	defer func() {
		end(status)
	}()

	res, err := h.client.RestoreMerchantPolicies(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "Restore")
	}

	h.cache.DeleteMerchantPolicyCache(ctx, 0)

	logSuccess("Successfully restored merchant policy")
	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantPoliciesDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete merchant policy
// @Tags Merchant Policy Command
// @Description Permanently delete a merchant policy record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Policy ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant policy record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid policy ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant policy permanently"
// @Router /api/merchant-policy-command/permanent/{id} [delete]
func (h *merchantPolicyCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "DeleteMerchantPolicyPermanent")
	defer func() {
		end(status)
	}()

	res, err := h.client.DeleteMerchantPoliciesPermanent(ctx, &pb.FindByIdMerchantPoliciesRequest{Id: int32(id)})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "DeletePermanent")
	}

	h.cache.DeleteMerchantPolicyCache(ctx, 0)

	logSuccess("Successfully deleted merchant policy permanently")
	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed merchant policies
// @Tags Merchant Policy Command
// @Description Restore all trashed merchant policy records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchant policies"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant policies"
// @Router /api/merchant-policy-command/restore/all [post]
func (h *merchantPolicyCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "RestoreAllMerchantPolicies")
	defer func() {
		end(status)
	}()

	res, err := h.client.RestoreAllMerchantPolicies(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "RestoreAll")
	}

	logSuccess("Successfully restored all merchant policies")
	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed merchant policies
// @Tags Merchant Policy Command
// @Description Permanently delete all trashed merchant policy records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchant policies permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant policies permanently"
// @Router /api/merchant-policy-command/permanent/all [post]
func (h *merchantPolicyCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	ctx, span, end, status, logSuccess := h.observability.StartTracingAndLogging(ctx, "DeleteAllMerchantPoliciesPermanent")
	defer func() {
		end(status)
	}()

	res, err := h.client.DeleteAllMerchantPoliciesPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		status = "error"
		return h.handleError(c, err, span, "DeleteAllPermanent")
	}

	logSuccess("Successfully deleted all merchant policies permanently")
	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantAll(res))
}

func (h *merchantPolicyCommandHandlerApi) handleError(c echo.Context, err error, span trace.Span, method string) error {
	appErr := errors.ParseGrpcError(err)
	traceID := span.SpanContext().TraceID().String()

	h.logger.Error(
		"Merchant policy command error in "+method,
		zap.Error(err),
		zap.String("trace.id", traceID),
	)

	return errors.HandleApiError(c, appErr, traceID)
}
