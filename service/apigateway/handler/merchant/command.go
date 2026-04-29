package merchanthandler

import (
	"net/http"
	"strconv"

	merchant_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/merchant"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant"
	"github.com/labstack/echo/v4"
"google.golang.org/protobuf/types/known/emptypb"
)

type merchantCommandHandlerApi struct {
	client       pb.MerchantCommandServiceClient
	logger       logger.LoggerInterface
	mapper       apimapper.MerchantCommandResponseMapper
	cache        merchant_cache.MerchantCommandCache
	upload_image upload_image.ImageUploads
	apiHandler   errors.ApiHandler
}

type merchantCommandHandleDeps struct {
	client       pb.MerchantCommandServiceClient
	router       *echo.Echo
	logger       logger.LoggerInterface
	mapper       apimapper.MerchantCommandResponseMapper
	cache        merchant_cache.MerchantCommandCache
	upload_image upload_image.ImageUploads
	apiHandler   errors.ApiHandler
}

func NewMerchantCommandHandleApi(params *merchantCommandHandleDeps) *merchantCommandHandlerApi {
	merchantCommandHandler := &merchantCommandHandlerApi{
		client:       params.client,
		logger:       params.logger,
		mapper:       params.mapper,
		cache:        params.cache,
		upload_image: params.upload_image,
		apiHandler:   params.apiHandler,
	}

	routerMerchant := params.router.Group("/api/merchant-command")

	routerMerchant.POST("/create", merchantCommandHandler.Create)
	routerMerchant.POST("/update/:id", merchantCommandHandler.Update)
	routerMerchant.POST("/trashed/:id", merchantCommandHandler.TrashedMerchant)
	routerMerchant.POST("/restore/:id", merchantCommandHandler.RestoreMerchant)
	routerMerchant.DELETE("/permanent/:id", merchantCommandHandler.DeleteMerchantPermanent)
	routerMerchant.POST("/restore/all", merchantCommandHandler.RestoreAllMerchant)
	routerMerchant.POST("/permanent/all", merchantCommandHandler.DeleteAllMerchantPermanent)

	return merchantCommandHandler
}

// @Security Bearer
// @Summary Create a new merchant
// @Tags Merchant Command
// @Description Create a new merchant for the current user
// @Accept json
// @Produce json
// @Param request body requests.CreateMerchantRequest true "Merchant details"
// @Success 201 {object} response.ApiResponseMerchant "Successfully created merchant"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 500 {object} errors.ErrorResponse "Failed to create merchant"
// @Router /api/merchant-command/create [post]
func (h *merchantCommandHandlerApi) Create(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID <= 0 { return errors.ErrUnauthorized }

	var req requests.CreateMerchantRequest
	if err := c.Bind(&req); err != nil { return errors.NewBadRequestError("invalid request").WithInternal(err) }
	if err := req.Validate(); err != nil { return errors.NewValidationError(nil) }

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateMerchantRequest{
		UserId:       int32(userID),
		Name:         req.Name,
		Description:  req.Description,
		Address:      req.Address,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	return c.JSON(http.StatusCreated, h.mapper.ToApiResponseMerchant(res))
}

// @Security Bearer
// @Summary Update an existing merchant
// @Tags Merchant Command
// @Description Update an existing merchant record
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Param request body requests.UpdateMerchantRequest true "Updated merchant details"
// @Success 200 {object} response.ApiResponseMerchant "Successfully updated merchant"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body"
// @Failure 500 {object} errors.ErrorResponse "Failed to update merchant"
// @Router /api/merchant-command/update/{id} [post]
func (h *merchantCommandHandlerApi) Update(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok || userID <= 0 { return errors.ErrUnauthorized }

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	var req requests.UpdateMerchantRequest
	if err := c.Bind(&req); err != nil { return errors.NewBadRequestError("invalid request").WithInternal(err) }
	if err := req.Validate(); err != nil { return errors.NewValidationError(nil) }

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateMerchantRequest{
		MerchantId:   int32(id),
		UserId:       int32(userID),
		Name:         req.Name,
		Description:  req.Description,
		Address:      req.Address,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedMerchant(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchant(res))
}

// @Security Bearer
// @Summary Move merchant to trash
// @Tags Merchant Command
// @Description Move a merchant record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDeleteAt "Successfully moved merchant to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move merchant to trash"
// @Router /api/merchant-command/trashed/{id} [post]
func (h *merchantCommandHandlerApi) TrashedMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.TrashedMerchant(ctx, &pb.FindByIdMerchantRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedMerchant(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a trashed merchant
// @Tags Merchant Command
// @Description Restore a trashed merchant record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchant "Successfully restored merchant"
// @Failure 400 {object} errors.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-command/restore/{id} [post]
func (h *merchantCommandHandlerApi) RestoreMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.RestoreMerchant(ctx, &pb.FindByIdMerchantRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedMerchant(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchant(res))
}

// @Security Bearer
// @Summary Permanently delete a merchant
// @Tags Merchant Command
// @Description Permanently delete a merchant record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant permanently"
// @Router /api/merchant-command/permanent/{id} [delete]
func (h *merchantCommandHandlerApi) DeleteMerchantPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.DeleteMerchantPermanent(ctx, &pb.FindByIdMerchantRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedMerchant(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed merchants
// @Tags Merchant Command
// @Description Restore all trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored all merchants"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchants"
// @Router /api/merchant-command/restore/all [post]
func (h *merchantCommandHandlerApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllMerchant(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	// Invalidate all for bulk?
	// Normaly we'd clear everything related to merchants.

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed merchants
// @Tags Merchant Command
// @Description Permanently delete all trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted all merchants permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchants permanently"
// @Router /api/merchant-command/permanent/all [post]
func (h *merchantCommandHandlerApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantAll(res))
}


