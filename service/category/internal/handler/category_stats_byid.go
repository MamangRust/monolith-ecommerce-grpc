package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	category_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type categoryStatsByIdHandler struct {
	pb.UnimplementedCategoryStatsByIdServiceServer
	categoryStatsById service.CategoryStatsByIdService
	logger            logger.LoggerInterface
}

func NewCategoryStatsByIdHandler(service *service.Service, logger logger.LoggerInterface) *categoryStatsByIdHandler {
	return &categoryStatsByIdHandler{
		categoryStatsById: service.CategoryStatsById,
		logger:            logger,
	}
}

func (h *categoryStatsByIdHandler) FindMonthlyTotalPricesById(ctx context.Context, req *pb.FindYearMonthTotalPriceById) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
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

func (h *categoryStatsByIdHandler) FindYearlyTotalPricesById(ctx context.Context, req *pb.FindYearTotalPriceById) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
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

func (h *categoryStatsByIdHandler) FindMonthPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryMonthPrice, error) {
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

	data := make([]*pb.CategoryMonthPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = (&Handler{}).mapToCategoryResponse(result).(*pb.CategoryMonthPriceResponse)
	}

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Category monthly payment methods retrieved successfully",
		Data:    data,
	}, nil
}

func (h *categoryStatsByIdHandler) FindYearPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryYearPrice, error) {
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

	data := make([]*pb.CategoryYearPriceResponse, len(serviceResults))
	for i, result := range serviceResults {
		data[i] = (&Handler{}).mapToCategoryResponse(result).(*pb.CategoryYearPriceResponse)
	}

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Category yearly payment methods retrieved successfully",
		Data:    data,
	}, nil
}
