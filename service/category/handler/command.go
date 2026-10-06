package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
)

type categoryCommandHandler struct {
	pbcategory.UnimplementedCategoryCommandServiceServer
	service service.CategoryCommandService
	logger  logger.LoggerInterface
}

func NewCategoryCommandHandler(service service.CategoryCommandService, logger logger.LoggerInterface) pbcategory.CategoryCommandServiceServer {
	return &categoryCommandHandler{
		service: service,
		logger:  logger,
	}
}

func (h *categoryCommandHandler) Create(ctx context.Context, request *pbcategory.CreateCategoryRequest) (*pbcategory.ApiResponseCategory, error) {
	slug := request.GetSlugCategory()
	req := &requests.CreateCategoryRequest{
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		SlugCategory:  &slug,
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := h.service.Create(ctx, req)
	if err != nil {
		return nil, category_errors.ErrGrpcCreateCategory
	}

	return &pbcategory.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pbcategory.CategoryResponse),
	}, nil
}

func (h *categoryCommandHandler) Update(ctx context.Context, request *pbcategory.UpdateCategoryRequest) (*pbcategory.ApiResponseCategory, error) {
	id := int(request.GetCategoryId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	slug := request.GetSlugCategory()
	req := &requests.UpdateCategoryRequest{
		CategoryID:    &id,
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		SlugCategory:  &slug,
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := h.service.Update(ctx, req)
	if err != nil {
		return nil, category_errors.ErrGrpcUpdateCategory
	}

	return &pbcategory.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pbcategory.CategoryResponse),
	}, nil
}

func (h *categoryCommandHandler) TrashedCategory(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.Trash(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pbcategory.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pbcategory.CategoryResponseDeleteAt),
	}, nil
}

func (h *categoryCommandHandler) RestoreCategory(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.Restore(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pbcategory.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pbcategory.CategoryResponseDeleteAt),
	}, nil
}

func (h *categoryCommandHandler) DeleteCategoryPermanent(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategoryDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	_, err := h.service.DeletePermanent(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcDeleteCategory
	}

	return &pbcategory.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (h *categoryCommandHandler) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pbcategory.ApiResponseCategoryAll, error) {
	_, err := h.service.RestoreAll(ctx)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pbcategory.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restored all categories",
	}, nil
}

func (h *categoryCommandHandler) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcategory.ApiResponseCategoryAll, error) {
	_, err := h.service.DeleteAll(ctx)
	if err != nil {
		return nil, category_errors.ErrGrpcDeleteCategory
	}

	return &pbcategory.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully deleted all categories permanently",
	}, nil
}
