package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type categoryQueryHandler struct {
	pb.UnimplementedCategoryQueryServiceServer
	service service.CategoryQueryService
	logger  logger.LoggerInterface
}

func NewCategoryQueryHandler(service service.CategoryQueryService, logger logger.LoggerInterface) pb.CategoryQueryServiceServer {
	return &categoryQueryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *categoryQueryHandler) FindAll(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategory, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := h.service.FindAll(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcFindAllCategory
	}

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(math.Ceil(float64(*totalRecords) / float64(pageSize))),
		TotalRecords: int32(*totalRecords),
	}

	results := make([]*pb.CategoryResponse, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToCategoryResponse(v).(*pb.CategoryResponse)
	}

	return &pb.ApiResponsePaginationCategory{
		Status:     "success",
		Message:    "Successfully fetched categories",
		Data:       results,
		Pagination: paginationMeta,
	}, nil
}

func (h *categoryQueryHandler) FindById(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategory, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.FindByID(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pb.CategoryResponse),
	}, nil
}

func (h *categoryQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := h.service.FindActive(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcFindAllCategory
	}

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(math.Ceil(float64(*totalRecords) / float64(pageSize))),
		TotalRecords: int32(*totalRecords),
	}

	results := make([]*pb.CategoryResponseDeleteAt, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToCategoryResponse(v).(*pb.CategoryResponseDeleteAt)
	}

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active categories",
		Data:       results,
		Pagination: paginationMeta,
	}, nil
}

func (h *categoryQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := h.service.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcFindAllCategory
	}

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(math.Ceil(float64(*totalRecords) / float64(pageSize))),
		TotalRecords: int32(*totalRecords),
	}

	results := make([]*pb.CategoryResponseDeleteAt, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToCategoryResponse(v).(*pb.CategoryResponseDeleteAt)
	}

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed categories",
		Data:       results,
		Pagination: paginationMeta,
	}, nil
}
