package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
)

type categoryStatsHandler struct {
	pbcategory.UnimplementedCategoryStatsServiceServer
	categoryStats service.CategoryStatsService
	logger        logger.LoggerInterface
}

func NewCategoryStatsHandler(categoryStats service.CategoryStatsService, logger logger.LoggerInterface) CategoryStatsHandler {
	return &categoryStatsHandler{
		categoryStats: categoryStats,
		logger:        logger,
	}
}

func (h *categoryStatsHandler) FindMonthlyTotalPrices(ctx context.Context, req *pbcategory.FindYearMonthTotalPrices) (*pbcategory.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	reqService := requests.MonthTotalPrice{
		Year:  year,
		Month: month,
	}

	serviceResults, err := h.categoryStats.FindMonthlyTotalPrice(ctx, &reqService)
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

func (h *categoryStatsHandler) FindYearlyTotalPrices(ctx context.Context, req *pbcategory.FindYearTotalPrices) (*pbcategory.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := h.categoryStats.FindYearlyTotalPrice(ctx, year)
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

func (h *categoryStatsHandler) FindMonthPrice(ctx context.Context, req *pbcategory.FindYearCategory) (*pbcategory.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := h.categoryStats.FindMonthPrice(ctx, year)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoryMonthPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoryMonthPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsHandler) FindYearPrice(ctx context.Context, req *pbcategory.FindYearCategory) (*pbcategory.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := h.categoryStats.FindYearPrice(ctx, year)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoryYearPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoryYearPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func mapToCategoryResponse(data interface{}) interface{} {
	return (&Handler{}).mapToCategoryResponse(data)
}
