package merchantbusinesshandler

import (
	"net/http"
	"strconv"

	merchantbusiness_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_business"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_business"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantBusinessCommandHandlerApi struct {
	client         pb.MerchantBusinessCommandServiceClient
	logger         logger.LoggerInterface
	mapper         apimapper.MerchantBusinessCommandResponseMapper
	merchantMapper merchantapimapper.MerchantCommandResponseMapper
	cache          merchantbusiness_cache.MerchantBusinessCommandCache
}

type merchantBusinessCommandHandleDeps struct {
	client         pb.MerchantBusinessCommandServiceClient
	router         *echo.Echo
	logger         logger.LoggerInterface
	mapper         apimapper.MerchantBusinessCommandResponseMapper
	merchantMapper merchantapimapper.MerchantCommandResponseMapper
	cache          merchantbusiness_cache.MerchantBusinessCommandCache
}

func NewMerchantBusinessCommandHandleApi(params *merchantBusinessCommandHandleDeps) *merchantBusinessCommandHandlerApi {
	handler := &merchantBusinessCommandHandlerApi{
		client:         params.client,
		logger:         params.logger,
		mapper:         params.mapper,
		merchantMapper: params.merchantMapper,
		cache:          params.cache,
	}

	routerBusiness := params.router.Group("/api/merchant-business-command")
	routerBusiness.POST("/create", handler.Create)
	routerBusiness.POST("/update/:id", handler.Update)
	routerBusiness.POST("/trashed/:id", handler.Trash)
	routerBusiness.POST("/restore/:id", handler.Restore)
	routerBusiness.DELETE("/permanent/:id", handler.DeletePermanent)
	routerBusiness.POST("/restore/all", handler.RestoreAll)
	routerBusiness.POST("/permanent/all", handler.DeleteAllPermanent)

	return handler
}

// @Security Bearer
// @Summary Create merchant business information
// @Tags Merchant Business Command
// @Description Create business information for a merchant
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantBusinessInformationRequest true "Create merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully created merchant business info"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create merchant business info"
// @Router /api/merchant-business-command/create [post]
func (h *merchantBusinessCommandHandlerApi) Create(c echo.Context) error {
	var body requests.CreateMerchantBusinessInformationRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateMerchantBusinessRequest{
		MerchantId:        int32(body.MerchantID),
		BusinessType:      body.BusinessType,
		TaxId:             body.TaxID,
		EstablishedYear:   int32(body.EstablishedYear),
		NumberOfEmployees: int32(body.NumberOfEmployees),
		WebsiteUrl:        body.WebsiteUrl,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantBusiness(res))
}

// @Security Bearer
// @Summary Update merchant business information
// @Tags Merchant Business Command
// @Description Update existing business information for a merchant
// @Accept json
// @Produce json
// @Param id path int true "Business ID"
// @Param body body requests.UpdateMerchantBusinessInformationRequest true "Update merchant business request"
// @Success 200 {object} response.ApiResponseMerchantBusiness "Successfully updated merchant business info"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update merchant business info"
// @Router /api/merchant-business-command/update/{id} [post]
func (h *merchantBusinessCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	var body requests.UpdateMerchantBusinessInformationRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	body.MerchantBusinessInfoID = &id
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateMerchantBusinessRequest{
		MerchantBusinessInfoId: int32(id),
		BusinessType:           body.BusinessType,
		TaxId:                  body.TaxID,
		EstablishedYear:        int32(body.EstablishedYear),
		NumberOfEmployees:      int32(body.NumberOfEmployees),
		WebsiteUrl:             body.WebsiteUrl,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantBusinessCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantBusiness(res))
}

// @Security Bearer
// @Summary Move merchant business info to trash
// @Tags Merchant Business Command
// @Description Move a merchant business info record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Business ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully moved merchant business info to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid business ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move merchant business info to trash"
// @Router /api/merchant-business-command/trashed/{id} [post]
func (h *merchantBusinessCommandHandlerApi) Trash(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.TrashedMerchantBusiness(ctx, &pb.FindByIdMerchantBusinessRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantBusinessCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantBusinessDeleteAt(res))
}

// @Security Bearer
// @Summary Restore trashed merchant business info
// @Tags Merchant Business Command
// @Description Restore a trashed merchant business info record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Business ID"
// @Success 200 {object} response.ApiResponseMerchantBusinessDeleteAt "Successfully restored merchant business info"
// @Failure 400 {object} errors.ErrorResponse "Invalid business ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant business info"
// @Router /api/merchant-business-command/restore/{id} [post]
func (h *merchantBusinessCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.RestoreMerchantBusiness(ctx, &pb.FindByIdMerchantBusinessRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantBusinessCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantBusinessDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete merchant business info
// @Tags Merchant Business Command
// @Description Permanently delete a merchant business info record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Business ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant business info record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid business ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant business info permanently"
// @Router /api/merchant-business-command/permanent/{id} [delete]
func (h *merchantBusinessCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.DeleteMerchantBusinessPermanent(ctx, &pb.FindByIdMerchantBusinessRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantBusinessCache(ctx, id)

	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed merchant business info
// @Tags Merchant Business Command
// @Description Restore all trashed merchant business info records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchant business info"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant business info"
// @Router /api/merchant-business-command/restore/all [post]
func (h *merchantBusinessCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllMerchantBusiness(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed merchant business info
// @Tags Merchant Business Command
// @Description Permanently delete all trashed merchant business info records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchant business info permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant business info permanently"
// @Router /api/merchant-business-command/permanent/all [post]
func (h *merchantBusinessCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllMerchantBusinessPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantAll(res))
}

