package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-category/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"

	pbcategory "github.com/MamangRust/monolith-ecommerce-pb/category"
)

type categoryStatsByMerchantHandler struct {
	pbcategory.UnimplementedCategoryStatsByMerchantServiceServer
	categoryStatsByMerchant service.CategoryStatsByMerchantService
	logger                  logger.LoggerInterface
}

func NewCategoryStatsByMerchantHandler(categoryStatsByMerchant service.CategoryStatsByMerchantService, logger logger.LoggerInterface) CategoryStatsByMerchantHandler {
	return &categoryStatsByMerchantHandler{
		categoryStatsByMerchant: categoryStatsByMerchant,
		logger:                  logger,
	}
}

func (h *categoryStatsByMerchantHandler) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pbcategory.FindYearMonthTotalPriceByMerchant) (*pbcategory.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.MonthTotalPriceMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindMonthlyTotalPriceByMerchant(ctx, &reqService)
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

func (h *categoryStatsByMerchantHandler) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pbcategory.FindYearTotalPriceByMerchant) (*pbcategory.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.YearTotalPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindYearlyTotalPriceByMerchant(ctx, &reqService)
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

func (h *categoryStatsByMerchantHandler) FindMonthPriceByMerchant(ctx context.Context, req *pbcategory.FindYearCategoryByMerchant) (*pbcategory.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.MonthPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoryMonthPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoryMonthPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByMerchantHandler) FindYearPriceByMerchant(ctx context.Context, req *pbcategory.FindYearCategoryByMerchant) (*pbcategory.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	reqService := requests.YearPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	serviceResults, err := h.categoryStatsByMerchant.FindYearPriceByMerchant(ctx, &reqService)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pbcategory.CategoryYearPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = mapToCategoryResponse(result).(*pbcategory.CategoryYearPriceResponse)
	}

	return &pbcategory.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}
