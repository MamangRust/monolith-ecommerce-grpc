package merchantsociallinkhandler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_social_link"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
)

type merchantSocialLinkCommandHandlerApi struct {
	client     pb.MerchantSocialCommandServiceClient
	logger     logger.LoggerInterface
	mapper     apimapper.MerchantSocialLinkCommandResponseMapper
	apiHandler sharedErrors.ApiHandler
}

type merchantSocialLinkCommandHandleDeps struct {
	client     pb.MerchantSocialCommandServiceClient
	router     *echo.Echo
	logger     logger.LoggerInterface
	mapper     apimapper.MerchantSocialLinkCommandResponseMapper
	apiHandler sharedErrors.ApiHandler
}

func NewMerchantSocialLinkCommandHandleApi(params *merchantSocialLinkCommandHandleDeps) *merchantSocialLinkCommandHandlerApi {
	handler := &merchantSocialLinkCommandHandlerApi{
		client:     params.client,
		logger:     params.logger,
		mapper:     params.mapper,
		apiHandler: params.apiHandler,
	}

	routerSocial := params.router.Group("/api/merchant-social-link")
	routerSocial.POST("/create", handler.Create)
	routerSocial.POST("/update/:id", handler.Update)

	return handler
}

// @Security Bearer
// @Summary Create merchant social link
// @Tags Merchant Social Link Command
// @Description Create a new social media link for a merchant
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantSocialRequest true "Create merchant social link request"
// @Success 200 {object} response.ApiResponseMerchantSocialLink "Successfully created merchant social link"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create merchant social link"
// @Router /api/merchant-social-link/create [post]
func (h *merchantSocialLinkCommandHandlerApi) Create(c echo.Context) error {
	var body requests.CreateMerchantSocialRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateMerchantSocialRequest{
		MerchantDetailId: int32(*body.MerchantDetailID),
		Platform:         body.Platform,
		Url:              body.Url,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantSocialLink(res))
}

// @Security Bearer
// @Summary Update merchant social link
// @Tags Merchant Social Link Command
// @Description Update an existing social media link for a merchant
// @Accept json
// @Produce json
// @Param id path int true "Link ID"
// @Param body body requests.UpdateMerchantSocialRequest true "Update merchant social link request"
// @Success 200 {object} response.ApiResponseMerchantSocialLink "Successfully updated merchant social link"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update merchant social link"
// @Router /api/merchant-social-link/update/{id} [post]
func (h *merchantSocialLinkCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	var body requests.UpdateMerchantSocialRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateMerchantSocialRequest{
		Id:               int32(id),
		MerchantDetailId: int32(*body.MerchantDetailID),
		Platform:         body.Platform,
		Url:              body.Url,
	})
	if err != nil {
		return sharedErrors.ParseGrpcError(err)
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantSocialLink(res))
}

