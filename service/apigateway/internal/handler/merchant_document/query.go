package merchantdocumenthandler

import (
	"net/http"
	"strconv"

	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/merchant_documents"
	"github.com/labstack/echo/v4"
)

type merchantDocumentQueryHandlerApi struct {
	client pb.MerchantDocumentQueryServiceClient
	logger logger.LoggerInterface
	mapper apimapper.MerchantDocumentQueryResponseMapper
}

type merchantDocumentQueryHandleDeps struct {
	client pb.MerchantDocumentQueryServiceClient
	router *echo.Echo
	logger logger.LoggerInterface
	mapper apimapper.MerchantDocumentQueryResponseMapper
}

func NewMerchantDocumentQueryHandleApi(params *merchantDocumentQueryHandleDeps) *merchantDocumentQueryHandlerApi {
	handler := &merchantDocumentQueryHandlerApi{
		client: params.client,
		logger: params.logger,
		mapper: params.mapper,
	}

	routerDoc := params.router.Group("/api/merchant-document-query")
	routerDoc.GET("", handler.FindAll)
	routerDoc.GET("/:id", handler.FindById)
	routerDoc.GET("/active", handler.FindByActive)
	routerDoc.GET("/trashed", handler.FindByTrashed)

	return handler
}

// @Security Bearer
// @Summary Find all merchant documents
// @Tags Merchant Document Query
// @Description Retrieve a list of all merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocument "List of merchant documents"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant document data"
// @Router /api/merchant-document-query [get]
func (h *merchantDocumentQueryHandlerApi) FindAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	res, err := h.client.FindAll(ctx, &pb.FindAllMerchantDocumentsRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return h.handleGrpcError(err, "FindAll") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponsePaginationMerchantDocument(res))
}

// @Security Bearer
// @Summary Find merchant document by ID
// @Tags Merchant Document Query
// @Description Retrieve a merchant document by ID
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 200 {object} response.ApiResponseMerchantDocument "Merchant document data"
// @Failure 400 {object} errors.ErrorResponse "Invalid document ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant document data"
// @Router /api/merchant-document-query/{id} [get]
func (h *merchantDocumentQueryHandlerApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 { return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID") }

	ctx := c.Request().Context()
	res, err := h.client.FindById(ctx, &pb.FindMerchantDocumentByIdRequest{DocumentId: int32(id)})
	if err != nil { return h.handleGrpcError(err, "FindById") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseMerchantDocument(res))
}

// @Security Bearer
// @Summary Retrieve active merchant documents
// @Tags Merchant Document Query
// @Description Retrieve a list of active merchant documents
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocument "List of active merchant documents"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant document data"
// @Router /api/merchant-document-query/active [get]
func (h *merchantDocumentQueryHandlerApi) FindByActive(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	res, err := h.client.FindAllActive(ctx, &pb.FindAllMerchantDocumentsRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return h.handleGrpcError(err, "FindByActive") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponsePaginationMerchantDocument(res))
}

// @Security Bearer
// @Summary Retrieve trashed merchant documents
// @Tags Merchant Document Query
// @Description Retrieve a list of trashed merchant document records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationMerchantDocumentDeleteAt "List of trashed merchant document data"
// @Failure 500 {object} errors.ErrorResponse "Failed to retrieve merchant document data"
// @Router /api/merchant-document-query/trashed [get]
func (h *merchantDocumentQueryHandlerApi) FindByTrashed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 { page = 1 }
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 { pageSize = 10 }
	search := c.QueryParam("search")

	ctx := c.Request().Context()
	res, err := h.client.FindAllTrashed(ctx, &pb.FindAllMerchantDocumentsRequest{
		Page: int32(page), PageSize: int32(pageSize), Search: search,
	})
	if err != nil { return h.handleGrpcError(err, "FindByTrashed") }

	return c.JSON(http.StatusOK, h.mapper.ToApiResponsePaginationMerchantDocumentDeleteAt(res))
}

func (h *merchantDocumentQueryHandlerApi) handleGrpcError(err error, operation string) error {
	h.logger.Error("Failed to " + operation)
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to "+operation)
}
