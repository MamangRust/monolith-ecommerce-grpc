package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryCommandHandler struct {
	pb.UnimplementedCategoryCommandServiceServer
	service service.CategoryCommandService
	logger  logger.LoggerInterface
}

func NewCategoryCommandHandler(service service.CategoryCommandService, logger logger.LoggerInterface) pb.CategoryCommandServiceServer {
	return &categoryCommandHandler{
		service: service,
		logger:  logger,
	}
}

func (h *categoryCommandHandler) Create(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.ApiResponseCategory, error) {
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

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pb.CategoryResponse),
	}, nil
}

func (h *categoryCommandHandler) Update(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ApiResponseCategory, error) {
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

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pb.CategoryResponse),
	}, nil
}

func (h *categoryCommandHandler) TrashedCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.Trash(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pb.CategoryResponseDeleteAt),
	}, nil
}

func (h *categoryCommandHandler) RestoreCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.Restore(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pb.CategoryResponseDeleteAt),
	}, nil
}

func (h *categoryCommandHandler) DeleteCategoryPermanent(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	_, err := h.service.DeletePermanent(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcDeleteCategory
	}

	return &pb.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (h *categoryCommandHandler) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	_, err := h.service.RestoreAll(ctx)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restored all categories",
	}, nil
}

func (h *categoryCommandHandler) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	_, err := h.service.DeleteAll(ctx)
	if err != nil {
		return nil, category_errors.ErrGrpcDeleteCategory
	}

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully deleted all categories permanently",
	}, nil
}
