package merchantdetailhandler

import (
	"net/http"
	"strconv"

	merchant_detail_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant_detail"
	merchantapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDetailCommandHandlerApi struct {
	client         pb.MerchantDetailCommandServiceClient
	logger         logger.LoggerInterface
	mapper         apimapper.MerchantDetailCommandResponseMapper
	merchantMapper merchantapimapper.MerchantCommandResponseMapper
	cache          merchant_detail_cache.MerchantDetailCommandCache
	upload_image   upload_image.ImageUploads
	apiHandler     sharedErrors.ApiHandler
}

type merchantDetailCommandHandleDeps struct {
	client         pb.MerchantDetailCommandServiceClient
	router         *echo.Echo
	logger         logger.LoggerInterface
	mapper         apimapper.MerchantDetailCommandResponseMapper
	merchantMapper merchantapimapper.MerchantCommandResponseMapper
	cache          merchant_detail_cache.MerchantDetailCommandCache
	upload_image   upload_image.ImageUploads
	apiHandler     sharedErrors.ApiHandler
}

func NewMerchantDetailCommandHandleApi(params *merchantDetailCommandHandleDeps) *merchantDetailCommandHandlerApi {
	handler := &merchantDetailCommandHandlerApi{
		client:         params.client,
		logger:         params.logger,
		mapper:         params.mapper,
		merchantMapper: params.merchantMapper,
		cache:          params.cache,
		upload_image:   params.upload_image,
		apiHandler:     params.apiHandler,
	}

	router := params.router.Group("/api/merchant-detail-command")
	router.POST("/create", handler.Create)
	router.POST("/update/:id", handler.Update)
	router.POST("/trashed/:id", handler.Trashed)
	router.POST("/restore/:id", handler.Restore)
	router.DELETE("/permanent/:id", handler.DeletePermanent)
	router.POST("/restore/all", handler.RestoreAll)
	router.POST("/permanent/all", handler.DeleteAllPermanent)

	return handler
}

// @Security Bearer
// @Summary Create merchant details
// @Tags Merchant Detail Command
// @Description Create detailed information for a merchant
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantDetailRequest true "Create merchant detail request"
// @Success 201 {object} response.ApiResponseMerchantDetail "Successfully created merchant detail"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create merchant detail"
// @Router /api/merchant-detail-command/create [post]
func (h *merchantDetailCommandHandlerApi) Create(c echo.Context) error {
	var req requests.CreateMerchantDetailRequest
	if err := c.Bind(&req); err != nil { return sharedErrors.NewBadRequestError("invalid request").WithInternal(err) }

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateMerchantDetailRequest{
		MerchantId:       int32(req.MerchantID),
		DisplayName:      req.DisplayName,
		CoverImageUrl:    req.CoverImageUrl,
		LogoUrl:          req.LogoUrl,
		ShortDescription: req.ShortDescription,
		WebsiteUrl:       req.WebsiteUrl,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusCreated, h.mapper.ToApiResponseMerchantDetail(res))
}

// @Security Bearer
// @Summary Update merchant details
// @Tags Merchant Detail Command
// @Description Update existing detailed information for a merchant
// @Accept json
// @Produce json
// @Param id path int true "Detail ID"
// @Param body body requests.UpdateMerchantDetailRequest true "Update merchant detail request"
// @Success 200 {object} response.ApiResponseMerchantDetail "Successfully updated merchant detail"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update merchant detail"
// @Router /api/merchant-detail-command/update/{id} [post]
func (h *merchantDetailCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return sharedErrors.NewBadRequestError("id is required") }

	var req requests.UpdateMerchantDetailRequest
	if err := c.Bind(&req); err != nil { return sharedErrors.NewBadRequestError("invalid request").WithInternal(err) }
	req.MerchantDetailID = &id
	if err := req.Validate(); err != nil { return sharedErrors.NewBadRequestError(err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateMerchantDetailRequest{
		MerchantDetailId: int32(id),
		DisplayName:      req.DisplayName,
		CoverImageUrl:    req.CoverImageUrl,
		LogoUrl:          req.LogoUrl,
		ShortDescription: req.ShortDescription,
		WebsiteUrl:       req.WebsiteUrl,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantDetailCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDetail(res))
}

// @Security Bearer
// @Summary Move merchant detail to trash
// @Tags Merchant Detail Command
// @Description Move a merchant detail record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Detail ID"
// @Success 200 {object} response.ApiResponseMerchantDetailDeleteAt "Successfully moved merchant detail to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move merchant detail to trash"
// @Router /api/merchant-detail-command/trashed/{id} [post]
func (h *merchantDetailCommandHandlerApi) Trashed(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return sharedErrors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.TrashedMerchantDetail(ctx, &pb.FindByIdMerchantDetailRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantDetailCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDetailDeleteAt(res))
}

// @Security Bearer
// @Summary Restore trashed merchant detail
// @Tags Merchant Detail Command
// @Description Restore a trashed merchant detail record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Detail ID"
// @Success 200 {object} response.ApiResponseMerchantDetailDeleteAt "Successfully restored merchant detail"
// @Failure 400 {object} errors.ErrorResponse "Invalid detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant detail"
// @Router /api/merchant-detail-command/restore/{id} [post]
func (h *merchantDetailCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return sharedErrors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.RestoreMerchantDetail(ctx, &pb.FindByIdMerchantDetailRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantDetailCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDetailDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete merchant detail
// @Tags Merchant Detail Command
// @Description Permanently delete a merchant detail record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Detail ID"
// @Success 204 "Successfully deleted merchant detail record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid detail ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant detail permanently"
// @Router /api/merchant-detail-command/permanent/{id} [delete]
func (h *merchantDetailCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return sharedErrors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.DeleteMerchantDetailPermanent(ctx, &pb.FindByIdMerchantDetailRequest{Id: int32(id)})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	h.cache.DeleteMerchantDetailCache(ctx, id)

	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed merchant details
// @Tags Merchant Detail Command
// @Description Restore all trashed merchant detail records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchant details"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant details"
// @Router /api/merchant-detail-command/restore/all [post]
func (h *merchantDetailCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllMerchantDetail(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed merchant details
// @Tags Merchant Detail Command
// @Description Permanently delete all trashed merchant detail records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchant details permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant details permanently"
// @Router /api/merchant-detail-command/permanent/all [post]
func (h *merchantDetailCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllMerchantDetailPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.merchantMapper.ToApiResponseMerchantAll(res))
}

