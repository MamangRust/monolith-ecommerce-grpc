package sliderhandler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	slider_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/slider"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/slider"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sliderCommandHandleApi struct {
	client      pb.SliderCommandServiceClient
	logger      logger.LoggerInterface
	mapper      apimapper.SliderCommandResponseMapper
	queryMapper apimapper.SliderQueryResponseMapper
	cache       slider_cache.SliderCommandCache
	upload      upload_image.ImageUploads
}

type sliderCommandHandleDeps struct {
	client      pb.SliderCommandServiceClient
	router      *echo.Echo
	logger      logger.LoggerInterface
	mapper      apimapper.SliderCommandResponseMapper
	queryMapper apimapper.SliderQueryResponseMapper
	cache       slider_cache.SliderCommandCache
	upload      upload_image.ImageUploads
}

func NewSliderCommandHandleApi(deps *sliderCommandHandleDeps) {
	handler := &sliderCommandHandleApi{
		client:      deps.client,
		logger:      deps.logger,
		mapper:      deps.mapper,
		queryMapper: deps.queryMapper,
		cache:       deps.cache,
		upload:      deps.upload,
	}

	router := deps.router.Group("/api/slider-command")
	router.POST("/create", handler.Create)
	router.POST("/update/:id", handler.Update)
	router.POST("/trashed/:id", handler.TrashedSlider)
	router.POST("/restore/:id", handler.RestoreSlider)
	router.DELETE("/permanent/:id", handler.DeleteSliderPermanent)
	router.POST("/restore/all", handler.RestoreAllSlider)
	router.POST("/permanent/all", handler.DeleteAllSliderPermanent)
}

// @Security Bearer
// @Summary Create slider
// @Tags Slider Command
// @Description Create a new slider with an image upload
// @Accept mpfd
// @Produce json
// @Param name formData string true "Slider name"
// @Param image_slider formData file true "Slider image"
// @Success 201 {object} response.ApiResponseSlider "Successfully created slider"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create slider"
// @Router /api/slider-command/create [post]
func (h *sliderCommandHandleApi) Create(c echo.Context) error {
	formData, err := h.parseSliderForm(c, true)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	grpcReq := &pb.CreateSliderRequest{
		Name:  formData.Nama,
		Image: formData.FilePath,
	}

	res, err := h.client.Create(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "Create")
	}

	return c.JSON(http.StatusCreated, h.queryMapper.ToApiResponseSlider(res))
}

// @Security Bearer
// @Summary Update slider
// @Tags Slider Command
// @Description Update an existing slider's name or image
// @Accept mpfd
// @Produce json
// @Param id path int true "Slider ID"
// @Param name formData string false "Slider name"
// @Param image_slider formData file false "Slider image"
// @Success 200 {object} response.ApiResponseSlider "Successfully updated slider"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update slider"
// @Router /api/slider-command/update/{id} [post]
func (h *sliderCommandHandleApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	formData, err := h.parseSliderForm(c, false)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	grpcReq := &pb.UpdateSliderRequest{
		Id:    int32(id),
		Name:  formData.Nama,
		Image: formData.FilePath,
	}

	res, err := h.client.Update(ctx, grpcReq)
	if err != nil {
		return h.handleGrpcError(err, "Update")
	}

	h.cache.DeleteSliderCache(ctx, id)

	return c.JSON(http.StatusOK, h.queryMapper.ToApiResponseSlider(res))
}

// @Security Bearer
// @Summary Move slider to trash
// @Tags Slider Command
// @Description Move a slider record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Slider ID"
// @Success 200 {object} response.ApiResponseSliderDeleteAt "Successfully moved slider to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move slider to trash"
// @Router /api/slider-command/trashed/{id} [post]
func (h *sliderCommandHandleApi) TrashedSlider(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	res, err := h.client.TrashedSlider(ctx, &pb.FindByIdSliderRequest{Id: int32(id)})
	if err != nil {
		return h.handleGrpcError(err, "Trash")
	}

	h.cache.DeleteSliderCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseSliderDeleteAt(res))
}

// @Security Bearer
// @Summary Restore trashed slider
// @Tags Slider Command
// @Description Restore a trashed slider record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Slider ID"
// @Success 200 {object} response.ApiResponseSliderDeleteAt "Successfully restored slider"
// @Failure 400 {object} errors.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore slider"
// @Router /api/slider-command/restore/{id} [post]
func (h *sliderCommandHandleApi) RestoreSlider(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	res, err := h.client.RestoreSlider(ctx, &pb.FindByIdSliderRequest{Id: int32(id)})
	if err != nil {
		return h.handleGrpcError(err, "Restore")
	}

	h.cache.DeleteSliderCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseSliderDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete slider
// @Tags Slider Command
// @Description Permanently delete a slider record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Slider ID"
// @Success 200 {object} response.ApiResponseSliderDelete "Successfully deleted slider record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid slider ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete slider permanently"
// @Router /api/slider-command/permanent/{id} [delete]
func (h *sliderCommandHandleApi) DeleteSliderPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	ctx := c.Request().Context()
	res, err := h.client.DeleteSliderPermanent(ctx, &pb.FindByIdSliderRequest{Id: int32(id)})
	if err != nil {
		return h.handleGrpcError(err, "Delete")
	}

	h.cache.DeleteSliderCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseSliderDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed sliders
// @Tags Slider Command
// @Description Restore all trashed slider records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseSliderAll "Successfully restored all sliders"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore sliders"
// @Router /api/slider-command/restore/all [post]
func (h *sliderCommandHandleApi) RestoreAllSlider(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllSlider(ctx, &emptypb.Empty{})
	if err != nil {
		return h.handleGrpcError(err, "RestoreAll")
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseSliderAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed sliders
// @Tags Slider Command
// @Description Permanently delete all trashed slider records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseSliderAll "Successfully deleted all sliders permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete sliders permanently"
// @Router /api/slider-command/permanent/all [post]
func (h *sliderCommandHandleApi) DeleteAllSliderPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllSliderPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return h.handleGrpcError(err, "DeleteAll")
	}

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseSliderAll(res))
}

func (h *sliderCommandHandleApi) parseSliderForm(c echo.Context, requireImage bool) (requests.SliderFormData, error) {
	var formData requests.SliderFormData

	formData.Nama = strings.TrimSpace(c.FormValue("name"))
	if formData.Nama == "" {
		return formData, echo.NewHTTPError(http.StatusBadRequest, "Name is required")
	}

	file, err := c.FormFile("image_slider")
	if err == nil {
		imagePath, err := h.upload.ProcessImageUpload(c, "uploads/slider", file, false)
		if err != nil {
			return formData, err
		}
		formData.FilePath = imagePath
	} else if requireImage {
		return formData, echo.NewHTTPError(http.StatusBadRequest, "Image is required")
	}

	return formData, nil
}

func (h *sliderCommandHandleApi) handleGrpcError(err error, operation string) error {
	h.logger.Error("Failed to "+operation, zap.Error(err))
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to "+operation)
}
