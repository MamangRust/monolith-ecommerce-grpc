package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
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

type merchantDocumentHandleApi struct {
	client          pb.MerchantDocumentServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.MerchantDocumentResponseMapper
	upload_image    upload_image.ImageUploads
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerMerchantDocument(
	router *echo.Echo,
	client pb.MerchantDocumentServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantDocumentResponseMapper,
	upload_image upload_image.ImageUploads,
) *merchantDocumentHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_document_handler_requests_total",
			Help: "Total number of banner requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_document_handler_request_duration_seconds",
			Help:    "Duration of banner requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	merchantDocumentHandler := &merchantDocumentHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		upload_image:    upload_image,
		trace:           otel.Tracer("merchant-document-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/merchant-document")

	routercategory.GET("", merchantDocumentHandler.FindAllMerchantDocument)
	routercategory.GET("/:id", merchantDocumentHandler.FindById)
	routercategory.GET("/active", merchantDocumentHandler.FindByActive)
	routercategory.GET("/trashed", merchantDocumentHandler.FindByTrashed)

	routercategory.POST("/create", merchantDocumentHandler.Create)
	routercategory.POST("/update/:id", merchantDocumentHandler.Update)

	routercategory.POST("/trashed/:id", merchantDocumentHandler.TrashedDocument)
	routercategory.POST("/restore/:id", merchantDocumentHandler.RestoreDocument)
	routercategory.DELETE("/permanent/:id", merchantDocumentHandler.Delete)

	routercategory.POST("/restore/all", merchantDocumentHandler.RestoreAllDocuments)
	routercategory.POST("/permanent/all", merchantDocumentHandler.DeleteAllDocumentsPermanent)

	return merchantDocumentHandler
}

// @Security Bearer
// @Summary Find all merchant
// @Tags MerchantDocument
// @Description Retrieve a list of all merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocument "List of merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-document [get]
func (h *merchantDocumentHandleApi) FindAllMerchantDocument(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllMerchantDocument"
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

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		status = "error"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find merchant by ID
// @Tags MerchantDocument
// @Description Retrieve a merchant by ID
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "merchant data"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-document/{id} [get]
func (h *merchantDocumentHandleApi) FindById(c echo.Context) error {
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	status := "success"

	defer func() { end(status) }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("Invalid merchant ID format", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	req := &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		status = "error"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindByIdMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active merchant
// @Tags MerchantDocument
// @Description Retrieve a list of active merchant
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocumentDeleteAt "List of active merchant"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-document/active [get]
func (h *merchantDocumentHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllActive(ctx, req)

	if err != nil {
		status = "error"
		logError("Failed to retrieve merchant data", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedFindAllActiveMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed merchant records.
// @Summary Retrieve trashed merchant
// @Tags MerchantDocument
// @Description Retrieve a list of trashed merchant records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponsePaginationMerchantDocumentDeleteAt "List of trashed merchant data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve merchant data"
// @Router /api/merchant-document/trashed [get]
func (h *merchantDocumentHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAllTrashed(ctx, req)

	if err != nil {
		status = "error"

		logError("Failed to retrieve merchant data", err, zap.Error(err))
		return merchantdocument_errors.ErrApiFailedFindAllTrashedMerchantDocuments(c)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocumentDeleteAt(res)

	logSuccess("Successfully retrieve merchant data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Create godoc
// @Summary Create a new merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Create a new document for a merchant
// @Accept json
// @Produce json
// @Param body body requests.CreateMerchantDocumentRequest true "Create merchant document request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Created document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create document"
// @Router /api/merchant-document/create [post]
func (h *merchantDocumentHandleApi) Create(c echo.Context) error {
	const method = "Create"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	formData, err := h.parseMerchantDocumentCreate(c)

	if err != nil {
		status = "error"

		logError("invalid request or validation erro", err, zap.Error(err))

		return merchantdocument_errors.ErrApiValidateCreateMerchantDocument(c)
	}

	req := &pb.CreateMerchantDocumentRequest{
		MerchantId:   int32(formData.MerchantID),
		DocumentType: formData.DocumentType,
		DocumentUrl:  formData.DocumentUrl,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		logError("failed create", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedCreateMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success create", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// Update godoc
// @Summary Update a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Update a merchant document with the given ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Param body body requests.UpdateMerchantDocumentRequest true "Update merchant document request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Updated document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update document"
// @Router /api/merchant-document/update/{id} [post]
func (h *merchantDocumentHandleApi) Update(c echo.Context) error {
	const method = "Update"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"
		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	formData, err := h.parseMerchantDocumentUpdate(c)

	if err != nil {
		status = "error"

		logError("invalid request or validation erro", err, zap.Error(err))

		return merchantdocument_errors.ErrApiValidateUpdateMerchantDocument(c)
	}

	req := &pb.UpdateMerchantDocumentRequest{
		DocumentId:   int32(id),
		MerchantId:   int32(formData.MerchantID),
		DocumentType: formData.DocumentType,
		DocumentUrl:  formData.DocumentUrl,
		Status:       formData.Status,
		Note:         formData.Note,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		status = "error"

		logError("failed update", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedUpdateMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success update", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// UpdateStatus godoc
// @Summary Update merchant document status
// @Tags Merchant Document Command
// @Security Bearer
// @Description Update the status of a merchant document
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Param body body requests.UpdateMerchantDocumentStatusRequest true "Update status request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Updated document"
// @Failure 400 {object} response.ErrorResponse "Bad request or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update document status"
// @Router /api/merchants-document/update-status/{id} [post]
func (h *merchantDocumentHandleApi) UpdateStatus(c echo.Context) error {
	const method = "UpdateStatus"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"
		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	var body requests.UpdateMerchantDocumentStatusRequest

	if err := c.Bind(&body); err != nil {
		status = "error"

		logError("failed bind update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	if err := body.Validate(); err != nil {
		status = "error"

		logError("failed validate update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	req := &pb.UpdateMerchantDocumentStatusRequest{
		DocumentId: int32(id),
		MerchantId: int32(body.MerchantID),
		Status:     body.Status,
		Note:       body.Note,
	}

	res, err := h.client.UpdateStatus(ctx, req)

	if err != nil {
		status = "error"

		logError("failed update status", err, zap.Error(err))

		return merchantdocument_errors.ErrApiBindUpdateMerchantDocumentStatus(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success update status", zap.Error(err))

	return c.JSON(http.StatusOK, so)
}

// TrashedDocument godoc
// @Summary Trashed a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Trashed a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Trashed document"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to trashed document"
// @Router /api/merchant-document/trashed/{id} [post]
func (h *merchantDocumentHandleApi) TrashedDocument(c echo.Context) error {
	const method = "TrashedDocument"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"

		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.client.Trashed(ctx, &pb.TrashedMerchantDocumentRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		logError("failed trashed", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedTrashMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("success trashed", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreDocument godoc
// @Summary Restore a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Restore a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Restored document"
// @Failure 400 {object} response.ErrorResponse "Bad request or invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore document"
// @Router /api/merchant-document/restore/{id} [post]
func (h *merchantDocumentHandleApi) RestoreDocument(c echo.Context) error {
	const method = "RestoreDocument"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		status = "error"

		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.client.Restore(ctx, &pb.RestoreMerchantDocumentRequest{
		DocumentId: int32(idInt),
	})

	if err != nil {
		status = "error"

		logError("failed restore", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedRestoreMerchantDocument(c)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	logSuccess("Success restore merchant document", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// Delete godoc
// @Summary Delete a merchant document
// @Tags Merchant Document Command
// @Security Bearer
// @Description Delete a merchant document by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocumentDelete "Deleted document"
// @Failure 400 {object} response.ErrorResponse "Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete document"
// @Router /api/merchant-document/permanent/{id} [delete]
func (h *merchantDocumentHandleApi) Delete(c echo.Context) error {
	const method = "Delete"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		status = "error"

		logError("failed parse id", err, zap.Error(err))

		return merchantdocument_errors.ErrApiInvalidMerchantDocumentID(c)
	}

	res, err := h.client.DeletePermanent(ctx, &pb.DeleteMerchantDocumentPermanentRequest{
		DocumentId: int32(id),
	})

	if err != nil {
		status = "error"

		logError("failed delete", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedDeleteMerchantDocumentPermanent(c)
	}

	so := h.mapping.ToApiResponseMerchantDocumentDeleteAt(res)

	logSuccess("Successfully deleted merchant document", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// RestoreAllDocuments godoc
// @Summary Restore all merchant documents
// @Tags Merchant Document Command
// @Security Bearer
// @Description Restore all merchant documents that were previously deleted
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully restored all documents"
// @Failure 500 {object} response.ErrorResponse "Failed to restore all documents"
// @Router /api/merchant-document/restore/all [post]
func (h *merchantDocumentHandleApi) RestoreAllDocuments(c echo.Context) error {
	const method = "RestoreAllDocuments"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := h.client.RestoreAll(ctx, &emptypb.Empty{})

	if err != nil {
		status = "error"

		logError("failed restore all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedRestoreAllMerchantDocuments(c)
	}

	response := h.mapping.ToApiResponseMerchantDocumentAll(res)

	logSuccess("Successfully restored all merchant documents")

	return c.JSON(http.StatusOK, response)
}

// DeleteAllDocumentsPermanent godoc
// @Summary Permanently delete all merchant documents
// @Tags Merchant Document Command
// @Security Bearer
// @Description Permanently delete all merchant documents from the database
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully deleted all documents permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to permanently delete all documents"
// @Router /api/merchant-document/permanent/all [post]
func (h *merchantDocumentHandleApi) DeleteAllDocumentsPermanent(c echo.Context) error {
	const method = "DeleteAllDocumentsPermanent"
	status := "success"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := h.client.DeleteAllPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		status = "status"

		logError("failed delete all", err, zap.Error(err))

		return merchantdocument_errors.ErrApiFailedDeleteAllMerchantDocumentsPermanent(c)
	}

	response := h.mapping.ToApiResponseMerchantDocumentAll(res)

	logSuccess("Successfully deleted all merchant documents permanently")

	return c.JSON(http.StatusOK, response)
}

func (h *merchantDocumentHandleApi) parseMerchantDocumentCreate(c echo.Context) (requests.CreateMerchantDocumentRequest, error) {
	var formData requests.CreateMerchantDocumentRequest
	var err error

	formData.MerchantID, err = strconv.Atoi(c.FormValue("merchant_id"))
	if err != nil || formData.MerchantID <= 0 {
		return formData, merchantdocument_errors.ErrApiInvalidMerchantId(c)
	}

	formData.DocumentType = strings.TrimSpace(c.FormValue("document_type"))
	if formData.DocumentType == "" {
		return formData, merchantdocument_errors.ErrApiDocumentTypeRequired(c)
	}

	file, err := c.FormFile("document_url")
	if err != nil {
		return formData, merchantdocument_errors.ErrApiDocumentFileRequired(c)
	}

	uploadedPath, err := h.upload_image.ProcessImageUpload(c, "uploads/merchant_document", file, true)
	if err != nil {
		return formData, err
	}
	formData.DocumentUrl = uploadedPath

	return formData, nil
}

func (h *merchantDocumentHandleApi) parseMerchantDocumentUpdate(c echo.Context) (requests.UpdateMerchantDocumentRequest, error) {
	var formData requests.UpdateMerchantDocumentRequest
	var err error

	formData.MerchantID, err = strconv.Atoi(c.FormValue("merchant_id"))
	if err != nil || formData.MerchantID <= 0 {
		return formData, merchantdocument_errors.ErrApiInvalidMerchantId(c)
	}

	formData.DocumentType = strings.TrimSpace(c.FormValue("document_type"))
	if formData.DocumentType == "" {
		return formData, merchantdocument_errors.ErrApiDocumentTypeRequired(c)
	}

	file, err := c.FormFile("document_url")
	if err == nil {
		uploadedPath, err := h.upload_image.ProcessImageUpload(c, "uploads/merchant_document/files", file, true)
		if err != nil {
			return formData, err
		}
		formData.DocumentUrl = uploadedPath
	}

	formData.Status = strings.TrimSpace(c.FormValue("status"))
	if formData.Status == "" {
		return formData, merchantdocument_errors.ErrApiStatusRequired(c)
	}

	formData.Note = strings.TrimSpace(c.FormValue("note"))
	if formData.Note == "" {
		return formData, merchantdocument_errors.ErrApiNoteRequired(c)
	}

	return formData, nil
}

func (s *merchantDocumentHandleApi) startTracingAndLogging(
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

func (s *merchantDocumentHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
