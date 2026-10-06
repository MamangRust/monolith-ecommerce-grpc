package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
)

type categoryStatsByIdHandler struct {
	pbcategory.UnimplementedCategoryStatsByIdServiceServer
	categoryStatsById service.CategoryStatsByIdService
	logger            logger.LoggerInterface
}

func NewCategoryStatsByIdHandler(categoryStatsById service.CategoryStatsByIdService, logger logger.LoggerInterface) CategoryStatsByIdHandler {
	return &categoryStatsByIdHandler{
		categoryStatsById: categoryStatsById,
		logger:            logger,
	}
}

func (h *categoryStatsByIdHandler) FindMonthlyTotalPricesById(ctx context.Context, req *pbcategory.FindYearMonthTotalPriceById) (*pbcategory.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.MonthTotalPriceCategory{
		Year:       year,
		Month:      month,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindMonthlyTotalPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoriesMonthlyTotalPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoriesMonthlyTotalPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByIdHandler) FindYearlyTotalPricesById(ctx context.Context, req *pbcategory.FindYearTotalPriceById) (*pbcategory.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.YearTotalPriceCategory{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindYearlyTotalPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoriesYearlyTotalPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoriesYearlyTotalPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByIdHandler) FindMonthPriceById(ctx context.Context, req *pbcategory.FindYearCategoryById) (*pbcategory.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.MonthPriceId{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindMonthPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoryMonthPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoryMonthPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Category monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByIdHandler) FindYearPriceById(ctx context.Context, req *pbcategory.FindYearCategoryById) (*pbcategory.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	reqService := requests.YearPriceId{
		Year:       year,
		CategoryID: id,
	}

	serviceResults, err := h.categoryStatsById.FindYearPriceById(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoryYearPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoryYearPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Category yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}
