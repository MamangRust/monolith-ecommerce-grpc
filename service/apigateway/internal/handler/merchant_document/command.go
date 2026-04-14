package merchantdocumenthandler

import (
	"net/http"
	"strconv"
	"strings"

	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_documents"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDocumentCommandHandlerApi struct {
	client       pb.MerchantDocumentCommandServiceClient
	logger       logger.LoggerInterface
	mapper       apimapper.MerchantDocumentCommandResponseMapper
	upload_image upload_image.ImageUploads
}

type merchantDocumentCommandHandleDeps struct {
	client       pb.MerchantDocumentCommandServiceClient
	router       *echo.Echo
	logger       logger.LoggerInterface
	mapper       apimapper.MerchantDocumentCommandResponseMapper
	upload_image upload_image.ImageUploads
}

func NewMerchantDocumentCommandHandleApi(params *merchantDocumentCommandHandleDeps) *merchantDocumentCommandHandlerApi {
	handler := &merchantDocumentCommandHandlerApi{
		client:       params.client,
		logger:       params.logger,
		mapper:       params.mapper,
		upload_image: params.upload_image,
	}

	routerDoc := params.router.Group("/api/merchant-document-command")
	routerDoc.POST("/create", handler.Create)
	routerDoc.POST("/update/:id", handler.Update)
	routerDoc.POST("/update-status/:id", handler.UpdateStatus)
	routerDoc.POST("/trashed/:id", handler.Trash)
	routerDoc.POST("/restore/:id", handler.Restore)
	routerDoc.DELETE("/permanent/:id", handler.DeletePermanent)
	routerDoc.POST("/restore/all", handler.RestoreAll)
	routerDoc.POST("/permanent/all", handler.DeleteAllPermanent)

	return handler
}

// @Security Bearer
// @Summary Create a new merchant document
// @Tags Merchant Document Command
// @Description Upload and create a new merchant document record
// @Accept mpfd
// @Produce json
// @Param merchant_id formData int true "Merchant ID"
// @Param document_type formData string true "Document type (e.g., NIB, SIUP)"
// @Param document_url formData file true "Document file"
// @Success 200 {object} response.ApiResponseMerchantDocument "Successfully created merchant document"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to create merchant document"
// @Router /api/merchant-document-command/create [post]
func (h *merchantDocumentCommandHandlerApi) Create(c echo.Context) error {
	formData, err := h.parseMerchantDocumentCreate(c)
	if err != nil { return err }

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateMerchantDocumentRequest{
		MerchantId:   int32(formData.MerchantID),
		DocumentType: formData.DocumentType,
		DocumentUrl:  formData.DocumentUrl,
	})
	if err != nil { return h.handleGrpcError(err, "Create") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocument(res))
}

// @Security Bearer
// @Summary Update merchant document
// @Tags Merchant Document Command
// @Description Update an existing merchant document record
// @Accept mpfd
// @Produce json
// @Param id path int true "Document ID"
// @Param merchant_id formData int true "Merchant ID"
// @Param document_type formData string true "Document type"
// @Param document_url formData file false "New document file"
// @Param document_url_old formData string false "Old document URL if no new file is uploaded"
// @Param status formData string true "Document status"
// @Param note formData string false "Document note"
// @Success 200 {object} response.ApiResponseMerchantDocument "Successfully updated merchant document"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update merchant document"
// @Router /api/merchant-document-command/update/{id} [post]
func (h *merchantDocumentCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	formData, err := h.parseMerchantDocumentUpdate(c)
	if err != nil { return err }

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateMerchantDocumentRequest{
		DocumentId:   int32(id),
		MerchantId:   int32(formData.MerchantID),
		DocumentType: formData.DocumentType,
		DocumentUrl:  formData.DocumentUrl,
		Status:       formData.Status,
		Note:       formData.Note,
	})
	if err != nil { return h.handleGrpcError(err, "Update") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocument(res))
}

// @Security Bearer
// @Summary Update merchant document status
// @Tags Merchant Document Command
// @Description Update the status and note of a merchant document
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Param body body requests.UpdateMerchantDocumentStatusRequest true "Update status request"
// @Success 200 {object} response.ApiResponseMerchantDocument "Successfully updated document status"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 400 {object} errors.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} errors.ErrorResponse "Failed to update document status"
// @Router /api/merchant-document-command/update-status/{id} [post]
func (h *merchantDocumentCommandHandlerApi) UpdateStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	var body requests.UpdateMerchantDocumentStatusRequest
	if err := c.Bind(&body); err != nil { return echo.NewHTTPError(http.StatusBadRequest, "Invalid request") }
	if err := body.Validate(); err != nil { return echo.NewHTTPError(http.StatusBadRequest, err.Error()) }

	ctx := c.Request().Context()
	res, err := h.client.UpdateStatus(ctx, &pb.UpdateMerchantDocumentStatusRequest{
		DocumentId: int32(id),
		MerchantId: int32(body.MerchantID),
		Status:     body.Status,
		Note:       body.Note,
	})
	if err != nil { return h.handleGrpcError(err, "UpdateStatus") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocument(res))
}

// @Security Bearer
// @Summary Move merchant document to trash
// @Tags Merchant Document Command
// @Description Move a merchant document record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Successfully moved merchant document to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid document ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move merchant document to trash"
// @Router /api/merchant-document-command/trashed/{id} [post]
func (h *merchantDocumentCommandHandlerApi) Trash(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.Trashed(ctx, &pb.TrashedMerchantDocumentRequest{DocumentId: int32(id)})
	if err != nil { return h.handleGrpcError(err, "Trash") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocument(res))
}

// @Security Bearer
// @Summary Restore trashed merchant document
// @Tags Merchant Document Command
// @Description Restore a trashed merchant document record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Successfully restored merchant document"
// @Failure 400 {object} errors.ErrorResponse "Invalid document ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant document"
// @Router /api/merchant-document-command/restore/{id} [post]
func (h *merchantDocumentCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.Restore(ctx, &pb.RestoreMerchantDocumentRequest{DocumentId: int32(id)})
	if err != nil { return h.handleGrpcError(err, "Restore") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocument(res))
}

// @Security Bearer
// @Summary Permanently delete merchant document
// @Tags Merchant Document Command
// @Description Permanently delete a merchant document record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocumentDeleteAt "Successfully deleted merchant document record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid document ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant document permanently"
// @Router /api/merchant-document-command/permanent/{id} [delete]
func (h *merchantDocumentCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.DeletePermanent(ctx, &pb.DeleteMerchantDocumentPermanentRequest{DocumentId: int32(id)})
	if err != nil { return h.handleGrpcError(err, "Delete") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocumentDeleteAt(res))
}

// @Security Bearer
// @Summary Restore all trashed merchant documents
// @Tags Merchant Document Command
// @Description Restore all trashed merchant document records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully restored all merchant documents"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore merchant documents"
// @Router /api/merchant-document-command/restore/all [post]
func (h *merchantDocumentCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAll(ctx, &emptypb.Empty{})
	if err != nil { return h.handleGrpcError(err, "RestoreAll") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocumentAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed merchant documents
// @Tags Merchant Document Command
// @Description Permanently delete all trashed merchant document records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseMerchantDocumentAll "Successfully deleted all merchant documents permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete merchant documents permanently"
// @Router /api/merchant-document-command/permanent/all [post]
func (h *merchantDocumentCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllPermanent(ctx, &emptypb.Empty{})
	if err != nil { return h.handleGrpcError(err, "DeleteAll") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocumentAll(res))
}

func (h *merchantDocumentCommandHandlerApi) parseMerchantDocumentCreate(c echo.Context) (requests.CreateMerchantDocumentRequest, error) {
	var formData requests.CreateMerchantDocumentRequest
	var err error

	formData.MerchantID, err = strconv.Atoi(c.FormValue("merchant_id"))
	if err != nil || formData.MerchantID <= 0 { return formData, echo.NewHTTPError(http.StatusBadRequest, "Invalid merchant ID") }

	formData.DocumentType = strings.TrimSpace(c.FormValue("document_type"))
	if formData.DocumentType == "" { return formData, echo.NewHTTPError(http.StatusBadRequest, "Document type is required") }

	file, err := c.FormFile("document_url")
	if err != nil { return formData, echo.NewHTTPError(http.StatusBadRequest, "Document file is required") }

	uploadedPath, err := h.upload_image.ProcessImageUpload(c, "uploads/merchant_document", file, true)
	if err != nil { return formData, err }
	formData.DocumentUrl = uploadedPath

	return formData, nil
}

func (h *merchantDocumentCommandHandlerApi) parseMerchantDocumentUpdate(c echo.Context) (requests.UpdateMerchantDocumentRequest, error) {
	var formData requests.UpdateMerchantDocumentRequest
	var err error

	formData.MerchantID, err = strconv.Atoi(c.FormValue("merchant_id"))
	if err != nil || formData.MerchantID <= 0 { return formData, echo.NewHTTPError(http.StatusBadRequest, "Invalid merchant ID") }

	formData.DocumentType = strings.TrimSpace(c.FormValue("document_type"))
	if formData.DocumentType == "" { return formData, echo.NewHTTPError(http.StatusBadRequest, "Document type is required") }

	file, err := c.FormFile("document_url")
	if err == nil {
		uploadedPath, err := h.upload_image.ProcessImageUpload(c, "uploads/merchant_document/files", file, true)
		if err != nil { return formData, err }
		formData.DocumentUrl = uploadedPath
	} else {
		formData.DocumentUrl = c.FormValue("document_url_old")
	}

	formData.Status = c.FormValue("status")
	formData.Note = c.FormValue("note")

	return formData, nil
}

func (h *merchantDocumentCommandHandlerApi) handleGrpcError(err error, operation string) error {
	h.logger.Error("Failed to " + operation)
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to "+operation)
}
