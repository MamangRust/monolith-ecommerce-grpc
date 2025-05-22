package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	review_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review"
	response_api "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/api"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type reviewHandleApi struct {
	client  pb.ReviewServiceClient
	logger  logger.LoggerInterface
	mapping response_api.ReviewResponseMapper
}

func NewHandlerReview(
	router *echo.Echo,
	client pb.ReviewServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.ReviewResponseMapper,
) *reviewHandleApi {
	reviewHandler := &reviewHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routerreview := router.Group("/api/review")

	routerreview.GET("", reviewHandler.FindAll)
	routerreview.GET("/product/:id", reviewHandler.FindByProduct)
	routerreview.GET("/active", reviewHandler.FindByActive)
	routerreview.GET("/trashed", reviewHandler.FindByTrashed)

	routerreview.POST("/create", reviewHandler.Create)
	routerreview.POST("/update/:id", reviewHandler.Update)

	routerreview.POST("/trashed/:id", reviewHandler.TrashedReview)
	routerreview.POST("/restore/:id", reviewHandler.RestoreReview)
	routerreview.DELETE("/permanent/:id", reviewHandler.DeleteReviewPermanent)

	routerreview.POST("/restore/all", reviewHandler.RestoreAllReview)
	routerreview.POST("/permanent/all", reviewHandler.DeleteAllReviewPermanent)

	return reviewHandler
}

// @Security Bearer
// @Summary Find all review
// @Tags Review
// @Description Retrieve a list of all review
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of review"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review [get]
func (h *reviewHandleApi) FindAll(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch reviews", zap.Error(err))
		return review_errors.ErrApiFailedFindAllReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReview(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find reviews by product ID
// @Tags Review
// @Description Retrieve a list of reviews for a specific product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of reviews for the product"
// @Failure 400 {object} response.ErrorResponse "Invalid product ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/product/{id} [get]
func (h *reviewHandleApi) FindByProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return review_errors.ErrApiReviewInvalidProductId(c)
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")
	ctx := c.Request().Context()

	req := &pb.FindAllReviewProductRequest{
		ProductId: int32(id),
		Page:      int32(page),
		PageSize:  int32(pageSize),
		Search:    search,
	}

	res, err := h.client.FindByProduct(ctx, req)
	if err != nil {
		h.logger.Error("Failed to fetch reviews product", zap.Error(err))
		return review_errors.ErrApiFailedFindProductReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewsDetail(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find reviews by merchant ID
// @Tags Review
// @Description Retrieve a list of reviews for a specific merchant
// @Accept json
// @Produce json
// @Param id path int true "merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReview "List of reviews for the merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/merchant/{id} [get]
func (h *reviewHandleApi) FindByMerchant(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return review_errors.ErrApiReviewInvalidMerchantId(c)
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")
	ctx := c.Request().Context()

	req := &pb.FindAllReviewMerchantRequest{
		MerchantId: int32(id),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, req)
	if err != nil {
		h.logger.Error("Failed to fetch reviews merchant", zap.Error(err))
		return review_errors.ErrApiFailedFindMerchantReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewsDetail(res)
	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active review
// @Tags Review
// @Description Retrieve a list of active review
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of active review"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/active [get]
func (h *reviewHandleApi) FindByActive(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch review details", zap.Error(err))

		return review_errors.ErrApiFailedFindActiveReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed review records.
// @Summary Retrieve trashed review
// @Tags Review
// @Description Retrieve a list of trashed review records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationReviewDeleteAt "List of trashed review data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve review data"
// @Router /api/review/trashed [get]
func (h *reviewHandleApi) FindByTrashed(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllReviewRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived reviews", zap.Error(err))
		return review_errors.ErrApiFailedFindTrashedReviews(c)
	}

	so := h.mapping.ToApiResponsePaginationReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new review without image upload.
// @Summary Create a new review
// @Tags Review
// @Description Create a new review with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateReviewRequest true "review details"
// @Success 201 {object} response.ApiResponseReview "Successfully created review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create review"
// @Router /api/review/create [post]
func (h *reviewHandleApi) Create(c echo.Context) error {
	var req requests.CreateReviewRequest

	if err := c.Bind(&req); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return review_errors.ErrApiBindCreateReview(c)
	}

	if err := req.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return review_errors.ErrApiValidateCreateReview(c)
	}

	ctx := c.Request().Context()

	grpcReq := &pb.CreateReviewRequest{
		UserId:    int32(req.UserID),
		ProductId: int32(req.ProductID),
		Comment:   req.Comment,
		Rating:    int32(req.Rating),
	}

	res, err := h.client.Create(ctx, grpcReq)

	if err != nil {
		h.logger.Error("review creation failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return review_errors.ErrApiFailedCreateReview(c)
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// Update handles the update of an existing review.
// @Summary Update an existing review
// @Tags Review
// @Description Update an existing review record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Param request body requests.UpdateReviewRequest true "review update details"
// @Success 200 {object} response.ApiResponseReview "Successfully updated review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update review"
// @Router /api/review/update/{id} [post]
func (h *reviewHandleApi) Update(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	var req requests.UpdateReviewRequest

	if err := c.Bind(&req); err != nil {
		return review_errors.ErrApiBindUpdateReview(c)
	}

	ctx := c.Request().Context()

	grpcReq := &pb.UpdateReviewRequest{
		ReviewId: int32(idInt),
		Name:     req.Name,
		Comment:  req.Comment,
		Rating:   int32(req.Rating),
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		h.logger.Error("review update failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return review_errors.ErrApiFailedUpdateReview(c)
	}

	return c.JSON(http.StatusOK, res)
}

// @Security Bearer
// TrashedReview retrieves a trashed review record by its ID.
// @Summary Retrieve a trashed review
// @Tags Review
// @Description Retrieve a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewDeleteAt "Successfully retrieved trashed review"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed review"
// @Router /api/review/trashed/{id} [get]
func (h *reviewHandleApi) TrashedReview(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid review ID format", zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedReview(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive review", zap.Error(err))
		return review_errors.ErrApiFailedTrashedReview(c)
	}

	so := h.mapping.ToApiResponseReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreReview restores a review record from the trash by its ID.
// @Summary Restore a trashed review
// @Tags Review
// @Description Restore a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} response.ApiResponseReviewDeleteAt "Successfully restored review"
// @Failure 400 {object} response.ErrorResponse "Invalid review ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore review"
// @Router /api/review/restore/{id} [post]
func (h *reviewHandleApi) RestoreReview(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid review ID format", zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreReview(ctx, req)

	if err != nil {
		h.logger.Error("Failed to restore review", zap.Error(err))
		return review_errors.ErrApiFailedRestoreReview(c)
	}

	so := h.mapping.ToApiResponseReviewDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteReviewPermanent permanently deletes a review record by its ID.
// @Summary Permanently delete a review
// @Tags review
// @Description Permanently delete a review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewDelete "Successfully deleted review record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete review:"
// @Router /api/review/delete/{id} [delete]
func (h *reviewHandleApi) DeleteReviewPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid review ID format", zap.Error(err))
		return review_errors.ErrApiReviewInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdReviewRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteReviewPermanent(ctx, req)

	if err != nil {
		h.logger.Error("Failed to permanently delete review", zap.Error(err))
		return review_errors.ErrApiFailedDeleteReviewPermanent(c)
	}

	so := h.mapping.ToApiResponseReviewDelete(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllReview restores a review record from the trash by its ID.
// @Summary Restore a trashed review
// @Tags review
// @Description Restore a trashed review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewAll "Successfully restored review all"
// @Failure 400 {object} response.ErrorResponse "Invalid review ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore review"
// @Router /api/review/restore/all [post]
func (h *reviewHandleApi) RestoreAllReview(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllReview(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk review restoration failed", zap.Error(err))
		return review_errors.ErrApiFailedRestoreAllReviews(c)
	}

	so := h.mapping.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully restored all review")

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllReviewPermanent permanently deletes a review record by its ID.
// @Summary Permanently delete a review
// @Tags Review
// @Description Permanently delete a review record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "review ID"
// @Success 200 {object} response.ApiResponseReviewAll "Successfully deleted review record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete review:"
// @Router /api/review/delete/all [post]
func (h *reviewHandleApi) DeleteAllReviewPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllReviewPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		h.logger.Error("Bulk review deletion failed", zap.Error(err))
		return review_errors.ErrApiFailedDeleteAllReviewsPermanent(c)
	}

	so := h.mapping.ToApiResponseReviewAll(res)

	h.logger.Debug("Successfully deleted all review permanently")

	return c.JSON(http.StatusOK, so)
}
