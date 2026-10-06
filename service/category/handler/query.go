package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
)

type categoryQueryHandler struct {
	pbcategory.UnimplementedCategoryQueryServiceServer
	service service.CategoryQueryService
	logger  logger.LoggerInterface
}

func NewCategoryQueryHandler(service service.CategoryQueryService, logger logger.LoggerInterface) pbcategory.CategoryQueryServiceServer {
	return &categoryQueryHandler{
		service: service,
		logger:  logger,
	}
}

func (h *categoryQueryHandler) FindAll(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategory, error) {
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

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(math.Ceil(float64(*totalRecords) / float64(pageSize))),
		TotalRecords: int32(*totalRecords),
	}

	results := make([]*pbcategory.CategoryResponse, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToCategoryResponse(v).(*pbcategory.CategoryResponse)
	}

	return &pbcategory.ApiResponsePaginationCategory{
		Status:     "success",
		Message:    "Successfully fetched categories",
		Data:       results,
		Pagination: paginationMeta,
	}, nil
}

func (h *categoryQueryHandler) FindById(ctx context.Context, request *pbcategory.FindByIdCategoryRequest) (*pbcategory.ApiResponseCategory, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := h.service.FindByID(ctx, id)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryNotFound
	}

	return &pbcategory.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    (&Handler{}).mapToCategoryResponse(category).(*pbcategory.CategoryResponse),
	}, nil
}

func (h *categoryQueryHandler) FindByActive(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error) {
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

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(math.Ceil(float64(*totalRecords) / float64(pageSize))),
		TotalRecords: int32(*totalRecords),
	}

	results := make([]*pbcategory.CategoryResponseDeleteAt, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToCategoryResponse(v).(*pbcategory.CategoryResponseDeleteAt)
	}

	return &pbcategory.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active categories",
		Data:       results,
		Pagination: paginationMeta,
	}, nil
}

func (h *categoryQueryHandler) FindByTrashed(ctx context.Context, request *pbcategory.FindAllCategoryRequest) (*pbcategory.ApiResponsePaginationCategoryDeleteAt, error) {
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

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(math.Ceil(float64(*totalRecords) / float64(pageSize))),
		TotalRecords: int32(*totalRecords),
	}

	results := make([]*pbcategory.CategoryResponseDeleteAt, len(categories))
	for i, v := range categories {
		results[i] = (&Handler{}).mapToCategoryResponse(v).(*pbcategory.CategoryResponseDeleteAt)
	}

	return &pbcategory.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed categories",
		Data:       results,
		Pagination: paginationMeta,
	}, nil
}
