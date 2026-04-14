package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type categoryStatsHandler struct {
	pb.UnimplementedCategoryStatsServiceServer
	categoryStats service.CategoryStatsService
	logger        logger.LoggerInterface
}

func NewCategoryStatsHandler(service *service.Service, logger logger.LoggerInterface) *categoryStatsHandler {
	return &categoryStatsHandler{
		categoryStats: service.CategoryStats,
		logger:        logger,
	}
}

func (h *categoryStatsHandler) FindMonthlyTotalPrices(ctx context.Context, req *pb.FindYearMonthTotalPrices) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
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

	data := make([]*pb.CategoriesMonthlyTotalPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = (&Handler{}).mapToCategoryResponse(result).(*pb.CategoriesMonthlyTotalPriceResponse)
	}

	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsHandler) FindYearlyTotalPrices(ctx context.Context, req *pb.FindYearTotalPrices) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := h.categoryStats.FindYearlyTotalPrice(ctx, year)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pb.CategoriesYearlyTotalPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = (&Handler{}).mapToCategoryResponse(result).(*pb.CategoriesYearlyTotalPriceResponse)
	}

	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsHandler) FindMonthPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := h.categoryStats.FindMonthPrice(ctx, year)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pb.CategoryMonthPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = (&Handler{}).mapToCategoryResponse(result).(*pb.CategoryMonthPriceResponse)
	}

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsHandler) FindYearPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	serviceResults, err := h.categoryStats.FindYearPrice(ctx, year)
	if err != nil {
		return nil, category_errors.ErrGrpcCategoryStats
	}

	data := make([]*pb.CategoryYearPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = (&Handler{}).mapToCategoryResponse(result).(*pb.CategoryYearPriceResponse)
	}

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}
