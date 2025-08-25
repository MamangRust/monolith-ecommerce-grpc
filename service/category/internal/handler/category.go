package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryHandleGrpc struct {
	pb.UnimplementedCategoryServiceServer
	categoryQuery           service.CategoryQueryService
	categoryCommand         service.CategoryCommandService
	categoryStats           service.CategoryStatsService
	categoryStatsById       service.CategoryStatsByIdService
	categoryStatsByMerchant service.CategoryStatsByMerchantService
	logger                  logger.LoggerInterface
	mapping                 protomapper.CategoryProtoMapper
}

func NewCategoryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.CategoryServiceServer {
	return &categoryHandleGrpc{
		categoryQuery:           service.CategoryQuery,
		categoryCommand:         service.CategoryCommand,
		categoryStats:           service.CategoryStats,
		categoryStatsById:       service.CategoryStatsById,
		logger:                  logger,
		categoryStatsByMerchant: service.CategoryStatsByMerchant,
		mapping:                 protomapper.NewCategoryProtoMapper(),
	}
}

func (s *categoryHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategory, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all categories",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all categories",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched all categories",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationCategory(paginationMeta, "success", "Successfully fetched categories", categories)
	return so, nil
}

func (s *categoryHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategory, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid category ID provided", zap.Int("category_id", id))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Fetching category by ID", zap.Int("category_id", id))

	category, err := s.categoryQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch category by ID",
			zap.Int("category_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched category by ID", zap.Int("category_id", id))

	so := s.mapping.ToProtoResponseCategory("success", "Successfully fetched category", category)
	return so, nil
}

func (s *categoryHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active categories",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active categories",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active categories",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
	)

	so := s.mapping.ToProtoResponsePaginationCategoryDeleteAt(paginationMeta, "success", "Successfully fetched active categories", categories)
	return so, nil
}

func (s *categoryHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed categories",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed categories",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed categories",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
	)

	so := s.mapping.ToProtoResponsePaginationCategoryDeleteAt(paginationMeta, "success", "Successfully fetched trashed categories", categories)
	return so, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPrices(ctx context.Context, req *pb.FindYearMonthTotalPrices) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly total prices", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for monthly total prices", zap.Int("month", month))
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	s.logger.Info("Fetching monthly total prices for all categories",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthTotalPrice{
		Year:  year,
		Month: month,
	}

	methods, err := s.categoryStats.FindMonthlyTotalPrice(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly total prices for all categories",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly total prices for all categories",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("category_count", len(methods)),
	)

	return s.mapping.ToProtoResponseMonthlyTotalPrice("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPrices(ctx context.Context, req *pb.FindYearTotalPrices) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly total prices", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	s.logger.Info("Fetching yearly total prices for all categories", zap.Int("year", year))

	methods, err := s.categoryStats.FindYearlyTotalPrice(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly total prices for all categories",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly total prices for all categories",
		zap.Int("year", year),
		zap.Int("category_count", len(methods)),
	)

	return s.mapping.ToProtoResponseYearlyTotalPrice("success", "Yearly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesById(ctx context.Context, req *pb.FindYearMonthTotalPriceById) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly category stats", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for monthly category stats", zap.Int("month", month))
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	if id <= 0 {
		s.logger.Error("Invalid category ID provided for monthly stats", zap.Int("category_id", id))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Fetching monthly total price for category",
		zap.Int("category_id", id),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthTotalPriceCategory{
		Year:       year,
		Month:      month,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindMonthlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly total price for category",
			zap.Int("category_id", id),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly total price for category",
		zap.Int("category_id", id),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	return s.mapping.ToProtoResponseMonthlyTotalPrice("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesById(ctx context.Context, req *pb.FindYearTotalPriceById) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly category stats", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if id <= 0 {
		s.logger.Error("Invalid category ID provided for yearly stats", zap.Int("category_id", id))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Fetching yearly total price for category",
		zap.Int("category_id", id),
		zap.Int("year", year),
	)

	reqService := requests.YearTotalPriceCategory{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindYearlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly total price for category",
			zap.Int("category_id", id),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly total price for category",
		zap.Int("category_id", id),
		zap.Int("year", year),
	)

	return s.mapping.ToProtoResponseYearlyTotalPrice("success", "Yearly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearMonthTotalPriceByMerchant) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly merchant stats", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for monthly merchant stats", zap.Int("month", month))
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly stats", zap.Int("merchant_id", merchantID))
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	s.logger.Info("Fetching monthly total prices by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthTotalPriceMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantID,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly total prices by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly total prices by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("category_count", len(methods)),
	)

	return s.mapping.ToProtoResponseMonthlyTotalPrice("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearTotalPriceByMerchant) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly merchant stats", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly stats", zap.Int("merchant_id", merchantID))
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	s.logger.Info("Fetching yearly total prices by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearTotalPriceMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	methods, err := s.categoryStatsByMerchant.FindYearlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly total prices by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly total prices by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("category_count", len(methods)),
	)

	return s.mapping.ToProtoResponseYearlyTotalPrice("success", "Yearly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly category trend", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	s.logger.Info("Fetching monthly price trend for all categories", zap.Int("year", year))

	methods, err := s.categoryStats.FindMonthPrice(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch monthly price trend for categories",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly price trend",
		zap.Int("year", year),
		zap.Int("data_points", len(methods)),
	)

	return s.mapping.ToProtoResponseCategoryMonthlyPrice("success", "Monthly sales trend retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly category trend", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	s.logger.Info("Fetching yearly price trend for all categories", zap.Int("year", year))

	methods, err := s.categoryStats.FindYearPrice(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly price trend for categories",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly price trend",
		zap.Int("year", year),
		zap.Int("data_points", len(methods)),
	)

	return s.mapping.ToProtoResponseCategoryYearlyPrice("success", "Yearly sales trend retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly merchant category trend", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly merchant trend", zap.Int("merchant_id", merchantID))
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	s.logger.Info("Fetching monthly price trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.MonthPriceMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly price trend by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly price trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("data_points", len(methods)),
	)

	return s.mapping.ToProtoResponseCategoryMonthlyPrice("success", "Merchant monthly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly merchant category trend", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly merchant trend", zap.Int("merchant_id", merchantID))
		return nil, category_errors.ErrGrpcCategoryInvalidMerchantId
	}

	s.logger.Info("Fetching yearly price trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearPriceMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	methods, err := s.categoryStatsByMerchant.FindYearPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly price trend by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly price trend by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("data_points", len(methods)),
	)

	return s.mapping.ToProtoResponseCategoryYearlyPrice("success", "Merchant yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())
	categoryID := int(req.GetCategoryId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly category trend by ID", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if categoryID <= 0 {
		s.logger.Error("Invalid category ID provided for monthly trend", zap.Int("category_id", categoryID))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Fetching monthly price trend by category ID",
		zap.Int("category_id", categoryID),
		zap.Int("year", year),
	)

	reqService := requests.MonthPriceId{
		Year:       year,
		CategoryID: categoryID,
	}

	methods, err := s.categoryStatsById.FindMonthPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly price trend by category ID",
			zap.Int("category_id", categoryID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly price trend by category ID",
		zap.Int("category_id", categoryID),
		zap.Int("year", year),
		zap.Int("data_points", len(methods)),
	)

	return s.mapping.ToProtoResponseCategoryMonthlyPrice("success", "Monthly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())
	categoryID := int(req.GetCategoryId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly category trend by ID", zap.Int("year", year))
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if categoryID <= 0 {
		s.logger.Error("Invalid category ID provided for yearly trend", zap.Int("category_id", categoryID))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Fetching yearly price trend by category ID",
		zap.Int("category_id", categoryID),
		zap.Int("year", year),
	)

	reqService := requests.YearPriceId{
		Year:       year,
		CategoryID: categoryID,
	}

	methods, err := s.categoryStatsById.FindYearPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly price trend by category ID",
			zap.Int("category_id", categoryID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly price trend by category ID",
		zap.Int("category_id", categoryID),
		zap.Int("year", year),
		zap.Int("data_points", len(methods)),
	)

	return s.mapping.ToProtoResponseCategoryYearlyPrice("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) Create(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.ApiResponseCategory, error) {
	s.logger.Info("Creating new category",
		zap.String("name", request.GetName()),
		zap.String("slug", request.GetSlugCategory()),
	)

	req := &requests.CreateCategoryRequest{
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		ImageCategory: request.GetImageCategory(),
		SlugCategory:  &request.SlugCategory,
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on category creation",
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := s.categoryCommand.CreateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create category",
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Category created successfully",
		zap.Int("category_id", int(category.ID)),
		zap.String("name", category.Name),
	)

	so := s.mapping.ToProtoResponseCategory("success", "Successfully created category", category)
	return so, nil
}

func (s *categoryHandleGrpc) Update(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ApiResponseCategory, error) {
	id := int(request.GetCategoryId())

	if id == 0 {
		s.logger.Error("Invalid category ID provided for update", zap.Int("category_id", id))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Updating category", zap.Int("category_id", id))

	req := &requests.UpdateCategoryRequest{
		CategoryID:    &id,
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		ImageCategory: request.GetImageCategory(),
		SlugCategory:  &request.SlugCategory,
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on category update",
			zap.Int("category_id", id),
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := s.categoryCommand.UpdateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update category",
			zap.Int("category_id", id),
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Category updated successfully",
		zap.Int("category_id", id),
		zap.String("name", category.Name),
	)

	so := s.mapping.ToProtoResponseCategory("success", "Successfully updated category", category)
	return so, nil
}

func (s *categoryHandleGrpc) TrashedCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid category ID for trashing", zap.Int("category_id", id))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Moving category to trash", zap.Int("category_id", id))

	category, err := s.categoryCommand.TrashedCategory(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash category",
			zap.Int("category_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Category moved to trash successfully",
		zap.Int("category_id", id),
		zap.String("name", category.Name),
	)

	so := s.mapping.ToProtoResponseCategoryDeleteAt("success", "Successfully trashed category", category)
	return so, nil
}

func (s *categoryHandleGrpc) RestoreCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid category ID for restore", zap.Int("category_id", id))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Restoring category from trash", zap.Int("category_id", id))

	category, err := s.categoryCommand.RestoreCategory(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore category",
			zap.Int("category_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Category restored successfully",
		zap.Int("category_id", id),
		zap.String("name", category.Name),
	)

	so := s.mapping.ToProtoResponseCategoryDeleteAt("success", "Successfully restored category", category)
	return so, nil
}

func (s *categoryHandleGrpc) DeleteCategoryPermanent(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid category ID for permanent deletion", zap.Int("category_id", id))
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	s.logger.Info("Permanently deleting category", zap.Int("category_id", id))

	_, err := s.categoryCommand.DeleteCategoryPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete category",
			zap.Int("category_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Category permanently deleted", zap.Int("category_id", id))

	so := s.mapping.ToProtoResponseCategoryDelete("success", "Successfully deleted category permanently")
	return so, nil
}

func (s *categoryHandleGrpc) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	s.logger.Info("Restoring all trashed categories")

	_, err := s.categoryCommand.RestoreAllCategories(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all categories", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All categories restored successfully")

	so := s.mapping.ToProtoResponseCategoryAll("success", "Successfully restored all categories")
	return so, nil
}

func (s *categoryHandleGrpc) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	s.logger.Info("Permanently deleting all trashed categories")

	_, err := s.categoryCommand.DeleteAllCategoriesPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all categories", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All categories permanently deleted")

	so := s.mapping.ToProtoResponseCategoryAll("success", "Successfully deleted all categories permanently")
	return so, nil
}
