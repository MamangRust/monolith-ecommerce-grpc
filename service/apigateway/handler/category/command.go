package categoryhandler

import (
	"net/http"
	"strconv"

	category_cache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/cache/category"
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/upload_image"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	apimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/category"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)


type categoryCommandHandlerApi struct {
	client       pb.CategoryCommandServiceClient
	logger       logger.LoggerInterface
	mapper       apimapper.CategoryCommandResponseMapper
	cache        category_cache.CategoryMencache
	upload_image upload_image.ImageUploads
	errors       errors.ApiHandler
}



type categoryCommandHandleDeps struct {
	client       pb.CategoryCommandServiceClient
	router       *echo.Echo
	logger       logger.LoggerInterface
	mapper       apimapper.CategoryCommandResponseMapper
	cache        category_cache.CategoryMencache
	upload_image upload_image.ImageUploads
	apiHandler   errors.ApiHandler
}

func NewCategoryCommandHandleApi(params *categoryCommandHandleDeps) *categoryCommandHandlerApi {
	categoryCommandHandler := &categoryCommandHandlerApi{
		client:       params.client,
		logger:       params.logger,
		mapper:       params.mapper,
		cache:        params.cache,
		upload_image: params.upload_image,
		errors:       params.apiHandler,
	}

	routerCategory := params.router.Group("/api/category-command")

	routerCategory.POST("/create", categoryCommandHandler.Create)
	routerCategory.POST("/update/:id", categoryCommandHandler.Update)
	routerCategory.POST("/trashed/:id", categoryCommandHandler.Trashed)
	routerCategory.POST("/restore/:id", categoryCommandHandler.Restore)
	routerCategory.DELETE("/permanent/:id", categoryCommandHandler.DeletePermanent)
	routerCategory.POST("/restore/all", categoryCommandHandler.RestoreAll)
	routerCategory.DELETE("/permanent/all", categoryCommandHandler.DeleteAllPermanent)

	return categoryCommandHandler
}

// @Security Bearer
// @Summary Create a new category
// @Tags Category Command
// @Description Create a new category with the provided details
// @Accept mpfd
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param slug_category formData string false "Category slug"
// @Param image formData file false "Category image"
// @Success 201 {object} response.ApiResponseCategory "Successfully created category"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} errors.ErrorResponse "Failed to create category"
// @Router /api/category-command/create [post]
func (h *categoryCommandHandlerApi) Create(c echo.Context) error {
	var req requests.CreateCategoryRequest
	if err := c.Bind(&req); err != nil { return errors.NewBadRequestError("invalid request").WithInternal(err) }

	if err := req.Validate(); err != nil { return errors.NewValidationError(nil) } // Simplified validation error

	file, err := c.FormFile("image")
	var imageURL string
	if err == nil {
		imageURL, err = h.upload_image.ProcessImageUpload(c, "uploads/category", file, false)
	}

	slugCategory := ""
	if req.SlugCategory != nil { slugCategory = *req.SlugCategory }

	ctx := c.Request().Context()
	res, err := h.client.Create(ctx, &pb.CreateCategoryRequest{
		Name: req.Name, Description: req.Description, SlugCategory: slugCategory, ImageCategory: imageURL,
	})

	if err != nil {
		return errors.ParseGrpcError(err)
	}



	return c.JSON(http.StatusCreated, h.mapper.ToApiResponseCategory(res))
}

// @Security Bearer
// @Summary Update an existing category
// @Tags Category Command
// @Description Update an existing category record with the provided details
// @Accept mpfd
// @Produce json
// @Param id path int true "Category ID"
// @Param name formData string false "Category name"
// @Param description formData string false "Category description"
// @Param slug_category formData string false "Category slug"
// @Param image formData file false "Category image"
// @Success 200 {object} response.ApiResponseCategory "Successfully updated category"
// @Failure 400 {object} errors.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} errors.ErrorResponse "Failed to update category"
// @Router /api/category-command/update/{id} [post]
func (h *categoryCommandHandlerApi) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	var req requests.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil { return errors.NewBadRequestError("invalid request").WithInternal(err) }

	file, err := c.FormFile("image")
	var imageURL string
	if err == nil {
		imageURL, err = h.upload_image.ProcessImageUpload(c, "uploads/category", file, false)
	}

	slugCategory := ""
	if req.SlugCategory != nil { slugCategory = *req.SlugCategory }

	ctx := c.Request().Context()
	res, err := h.client.Update(ctx, &pb.UpdateCategoryRequest{
		CategoryId: int32(id), Name: req.Name, Description: req.Description, SlugCategory: slugCategory, ImageCategory: imageURL,
	})

	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedCategoryCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCategory(res))
}

// @Security Bearer
// @Summary Move category to trash
// @Tags Category Command
// @Description Move a category record to trash by its ID
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully moved category to trash"
// @Failure 400 {object} errors.ErrorResponse "Invalid category ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to move category to trash"
// @Router /api/category-command/trashed/{id} [post]
func (h *categoryCommandHandlerApi) Trashed(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.TrashedCategory(ctx, &pb.FindByIdCategoryRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedCategoryCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCategoryDeleteAt(res))
}

// @Security Bearer
// @Summary Restore a trashed category
// @Tags Category Command
// @Description Restore a trashed category record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully restored category"
// @Failure 400 {object} errors.ErrorResponse "Invalid category ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore category"
// @Router /api/category-command/restore/{id} [post]
func (h *categoryCommandHandlerApi) Restore(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.RestoreCategory(ctx, &pb.FindByIdCategoryRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedCategoryCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCategoryDeleteAt(res))
}

// @Security Bearer
// @Summary Permanently delete a category
// @Tags Category Command
// @Description Permanently delete a category record by its ID
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDelete "Successfully deleted category record permanently"
// @Failure 400 {object} errors.ErrorResponse "Invalid category ID"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete category permanently"
// @Router /api/category-command/permanent/{id} [delete]
func (h *categoryCommandHandlerApi) DeletePermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { return errors.NewBadRequestError("id is required") }

	ctx := c.Request().Context()
	res, err := h.client.DeleteCategoryPermanent(ctx, &pb.FindByIdCategoryRequest{Id: int32(id)})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	h.cache.DeleteCachedCategoryCache(ctx, id)

	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCategoryDelete(res))
}

// @Security Bearer
// @Summary Restore all trashed categories
// @Tags Category Command
// @Description Restore all trashed category records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully restored all categories"
// @Failure 500 {object} errors.ErrorResponse "Failed to restore categories"
// @Router /api/category-command/restore/all [post]
func (h *categoryCommandHandlerApi) RestoreAll(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllCategory(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCategoryAll(res))
}

// @Security Bearer
// @Summary Permanently delete all trashed categories
// @Tags Category Command
// @Description Permanently delete all trashed category records
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully deleted all categories permanently"
// @Failure 500 {object} errors.ErrorResponse "Failed to delete categories permanently"
// @Router /api/category-command/permanent/all [delete]
func (h *categoryCommandHandlerApi) DeleteAllPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}


	return c.JSON(http.StatusOK, h.mapper.ToApiResponseCategoryAll(res))
}


