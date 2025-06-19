package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
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

type merchantDetailHandleApi struct {
	client          pb.MerchantDetailServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantDetailResponseMapper
	mappingMerchant response_api.MerchantResponseMapper
	upload_image    upload_image.ImageUploads
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchantDetail(
	router *echo.Echo,
	client pb.MerchantDetailServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantDetailResponseMapper,
	mappingMerchant response_api.MerchantResponseMapper,
	upload_image upload_image.ImageUploads,
) *merchantDetailHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_detail_handler_requests_total",
			Help: "Total number of banner requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_detail_handler_request_duration_seconds",
			Help:    "Duration of banner requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	merchantDetailHandler := &merchantDetailHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		mappingMerchant: mappingMerchant,
		upload_image:    upload_image,
		trace:           otel.Tracer("merchant-detail-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/merchant-detail")

	routercategory.GET("", merchantDetailHandler.FindAllMerchantDetail)
	routercategory.GET("/:id", merchantDetailHandler.FindById)
	routercategory.GET("/active", merchantDetailHandler.FindByActive)
	routercategory.GET("/trashed", merchantDetailHandler.FindByTrashed)

	routercategory.POST("/create", merchantDetailHandler.Create)
	routercategory.POST("/update/:id", merchantDetailHandler.Update)

	routercategory.POST("/trashed/:id", merchantDetailHandler.TrashedMerchant)
	routercategory.POST("/restore/:id", merchantDetailHandler.RestoreMerchant)
	routercategory.DELETE("/permanent/:id", merchantDetailHandler.DeleteMerchantPermanent)

	routercategory.POST("/restore/all", merchantDetailHandler.RestoreAllMerchant)
	routercategory.POST("/permanent/all", merchantDetailHandler.DeleteAllMerchantPermanent)

	return merchantDetailHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags MerchantDetail
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDetail "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail [get]
func (h *merchantDetailHandleApi) FindAllMerchantDetail(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllMerchantDetail"
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

		return merchantdetail_errors.ErrApiFailedFindAllMerchantDetail(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDetail(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantDetail
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDetail "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail/{id} [get]
func (h *merchantDetailHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID format", err, zap.Error(err))

		return merchantdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedFindMerchantDetailById(c)
	}

	so := h.mapping.ToApiResponseMerchantDetail(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags MerchantDetail
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDetailDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail/active [get]
func (h *merchantDetailHandleApi) FindByActive(c echo.Context) error {
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

		return merchantdetail_errors.ErrApiFailedFindActiveMerchantDetail(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDetailDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantDetail
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantDetailDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-detail/trashed [get]
func (h *merchantDetailHandleApi) FindByTrashed(c echo.Context) error {
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
		return merchantdetail_errors.ErrApiFailedFindActiveMerchantDetail(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDetailDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new merchant detail
// @Tags MerchantDetail
// @Description Create a new merchant detail with display name, cover image, logo, etc.
// @Accept multipart/form-data
// @Produce json
// @Param merchant_id formData int true "Merchant ID"
// @Param display_name formData string true "Display name"
// @Param short_description formData string true "Short description"
// @Param website_url formData string false "Website URL"
// @Param cover_image_url formData file true "Cover image file"
// @Param logo_url formData file true "Logo file"
// @Param social_links formData string true "Social links in JSON format (e.g., [{\"platform\": \"instagram\", \"url\": \"https://insta...\", \"merchant_detail_id\": 1}])"
// @Success 200 {object} response.ApiResponseMerchantDetail "Successfully created merchant detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/merchant-detail/create [post]
func (h *merchantDetailHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	formData, err := h.parseMerchantDetailCreate(c)
	if err != nil {
		status = "error"

		logError("Invalid request or validation error", err, zap.Error(err))

		return merchantdetail_errors.ErrApiInvalidBody(c)
	}

	var pbSocialLinks []*pb.CreateMerchantSocialRequest
	for _, link := range formData.SocialLinks {
		pbSocialLinks = append(pbSocialLinks, &pb.CreateMerchantSocialRequest{
			MerchantDetailId: int32(*link.MerchantDetailID),
			Platform:         link.Platform,
			Url:              link.Url,
		})
	}

	req := &pb.CreateMerchantDetailRequest{
		MerchantId:       int32(formData.MerchantID),
		DisplayName:      formData.DisplayName,
		CoverImageUrl:    formData.CoverImageUrl,
		LogoUrl:          formData.LogoUrl,
		ShortDescription: formData.ShortDescription,
		WebsiteUrl:       formData.WebsiteUrl,
		SocialLinks:      pbSocialLinks,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to create merchant detail", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedCreateMerchantDetail(c)
	}

	so := h.mapping.ToApiResponseMerchantDetail(res)

	logSuccess("Successfully create merchant detail", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update existing merchant detail
// @Tags MerchantDetail
// @Description Update an existing merchant detail by ID
// @Accept json
// @Produce json
// @Param id path int true "Merchant Detail ID"
// @Param merchant_id formData int true "Merchant ID"
// @Param display_name formData string true "Display name"
// @Param short_description formData string true "Short description"
// @Param website_url formData string false "Website URL"
// @Param cover_image_url formData file true "Cover image file"
// @Param logo_url formData file true "Logo file"
// @Param social_links formData string true "Social links in JSON format (e.g., [{\"platform\": \"instagram\", \"url\": \"https://insta...\", \"merchant_detail_id\": 1}])"
// @Success 200 {object} response.ApiResponseMerchantDetail "Successfully updated merchant detail"
// @Failure 400 {object} response.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/merchant-detail/update/{id} [post]
func (h *merchantDetailHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		status = "error"

		logError("Invalid merchant detail ID format", err, zap.Error(err))

		return merchantdetail_errors.ErrApiInvalidId(c)
	}

	formData, err := h.parseMerchantDetailUpdate(c)
	if err != nil {
		status = "error"

		logError("Invalid request or validation error", err, zap.Error(err))

		return merchantdetail_errors.ErrApiInvalidBody(c)
	}

	var pbSocialLinks []*pb.UpdateMerchantSocialRequest
	for _, link := range formData.SocialLinks {
		pbSocialLinks = append(pbSocialLinks, &pb.UpdateMerchantSocialRequest{
			Id:               int32(link.ID),
			MerchantDetailId: int32(*link.MerchantDetailID),
			Platform:         link.Platform,
			Url:              link.Url,
		})
	}

	req := &pb.UpdateMerchantDetailRequest{
		MerchantDetailId: int32(idInt),
		DisplayName:      formData.DisplayName,
		CoverImageUrl:    formData.CoverImageUrl,
		LogoUrl:          formData.LogoUrl,
		ShortDescription: formData.ShortDescription,
		WebsiteUrl:       formData.WebsiteUrl,
		SocialLinks:      pbSocialLinks,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		status = "error"

		logError("Failed to update merchant detail", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedUpdateMerchantDetail(c)
	}

	so := h.mapping.ToApiResponseMerchantDetail(res)

	logSuccess("Successfully update merchant detail", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedMerchant retrieves a trashed merchant record by its ID.
// @Summary Retrieve a trashed merchant
// @Tags MerchantDetail
// @Description Retrieve a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDetailDeleteAt "Successfully retrieved trashed merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed merchant"
// @Router /api/merchant-detail/trashed/{id} [get]
func (h *merchantDetailHandleApi) TrashedMerchant(c echo.Context) error {
	const method = "TrashedMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID format", err, zap.Error(err))

		return merchantdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchantDetail(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve trashed merchant", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedTrashedMerchantDetail(c)
	}

	so := h.mapping.ToApiResponseMerchantDetailDeleteAt(res)

	logSuccess("Successfully retrieve trashed merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDeleteAt "Successfully restored merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-detail/restore/{id} [post]
func (h *merchantDetailHandleApi) RestoreMerchant(c echo.Context) error {
	const method = "RestoreMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID format", err, zap.Error(err))

		return merchantdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchantDetail(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to restore merchant", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedRestoreMerchantDetail(c)
	}

	so := h.mapping.ToApiResponseMerchantDetailDeleteAt(res)

	logSuccess("Successfully restore merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDelete "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-detail/delete/{id} [delete]
func (h *merchantDetailHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	const method = "DeleteMerchantPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID format", err, zap.Error(err))

		return merchantdetail_errors.ErrApiInvalidId(c)
	}

	req := &pb.FindByIdMerchantDetailRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantDetailPermanent(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to delete merchant", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedDeleteMerchantDetailPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantDelete(res)

	logSuccess("Successfully delete merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllMerchant restores a merchant record from the trash by its ID.
// @Summary Restore a trashed merchant
// @Tags MerchantDetail
// @Description Restore a trashed merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully restored merchant all"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore merchant"
// @Router /api/merchant-detail/restore/all [post]
func (h *merchantDetailHandleApi) RestoreAllMerchant(c echo.Context) error {
	const method = "RestoreAllMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.RestoreAllMerchantDetail(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to restore all merchant", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedRestoreAllMerchantDetail(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	logSuccess("Successfully restore all merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllMerchantPermanent permanently deletes a merchant record by its ID.
// @Summary Permanently delete a merchant
// @Tags MerchantDetail
// @Description Permanently delete a merchant record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantAll "Successfully deleted merchant record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete merchant:"
// @Router /api/merchant-detail/delete/all [post]
func (h *merchantDetailHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	res, err := h.client.DeleteAllMerchantDetailPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("Failed to delete all merchant permanently", err, zap.Error(err))

		return merchantdetail_errors.ErrApiFailedDeleteAllMerchantDetailPermanent(c)
	}

	so := h.mappingMerchant.ToApiResponseMerchantAll(res)

	logSuccess("Successfully delete all merchant permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (h *merchantDetailHandleApi) parseMerchantDetailCreate(c echo.Context) (requests.CreateMerchantDetailFormData, error) {
	var formData requests.CreateMerchantDetailFormData
	var err error

	formData.MerchantID, err = strconv.Atoi(c.FormValue("merchant_id"))
	if err != nil || formData.MerchantID <= 0 {
		return formData, merchantdetail_errors.ErrApiInvalidMerchantId(c)
	}

	formData.DisplayName = strings.TrimSpace(c.FormValue("display_name"))
	if formData.DisplayName == "" {
		return formData, merchantdetail_errors.ErrApiDisplayNameRequired(c)
	}

	formData.ShortDescription = strings.TrimSpace(c.FormValue("short_description"))
	if formData.ShortDescription == "" {
		return formData, merchantdetail_errors.ErrApiShortDescriptionRequired(c)
	}

	formData.WebsiteUrl = strings.TrimSpace(c.FormValue("website_url"))

	coverFile, err := c.FormFile("cover_image_url")
	if err != nil {
		return formData, merchantdetail_errors.ErrApiCoverImageRequired(c)
	}

	coverPath, err := h.upload_image.ProcessImageUpload(c, coverFile)
	if err != nil {
		return formData, err
	}
	formData.CoverImageUrl = coverPath

	logoFile, err := c.FormFile("logo_url")
	if err != nil {
		return formData, merchantdetail_errors.ErrApiLogoRequired(c)
	}

	logoPath, err := h.upload_image.ProcessImageUpload(c, logoFile)
	if err != nil {
		return formData, err
	}
	formData.LogoUrl = logoPath

	socialLinksJson := c.FormValue("social_links")
	if socialLinksJson == "" {
		return formData, merchantdetail_errors.ErrApiSocialLinksRequired(c)
	}

	var parsedSocialLinks []requests.CreateMerchantSocialFormData
	if err := json.Unmarshal([]byte(socialLinksJson), &parsedSocialLinks); err != nil {
		return formData, merchantdetail_errors.ErrApiInvalidSocialLinks(c)
	}
	formData.SocialLinks = parsedSocialLinks

	return formData, nil
}

func (h *merchantDetailHandleApi) parseMerchantDetailUpdate(c echo.Context) (requests.UpdateMerchantDetailFormData, error) {
	var formData requests.UpdateMerchantDetailFormData
	var err error

	formData.DisplayName = strings.TrimSpace(c.FormValue("display_name"))
	formData.ShortDescription = strings.TrimSpace(c.FormValue("short_description"))
	formData.WebsiteUrl = strings.TrimSpace(c.FormValue("website_url"))

	coverFile, err := c.FormFile("cover_image_url")
	if err == nil {
		coverPath, err := h.upload_image.ProcessImageUpload(c, coverFile)
		if err != nil {
			return formData, err
		}
		formData.CoverImageUrl = coverPath
	}

	logoFile, err := c.FormFile("logo_url")
	if err == nil {
		logoPath, err := h.upload_image.ProcessImageUpload(c, logoFile)
		if err != nil {
			return formData, err
		}
		formData.LogoUrl = logoPath
	}

	socialLinksRaw := c.FormValue("social_links")
	if socialLinksRaw != "" {
		var links []requests.UpdateMerchantSocialFormData
		err = json.Unmarshal([]byte(socialLinksRaw), &links)
		if err != nil {
			return formData, merchantdetail_errors.ErrApiInvalidSocialLinks(c)
		}
		formData.SocialLinks = links
	}

	return formData, nil
}

func (s *merchantDetailHandleApi) startTracingAndLogging(
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

func (s *merchantDetailHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
