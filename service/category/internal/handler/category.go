package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-category/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/category_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryHandleGrpc struct {
	pb.UnimplementedCategoryServiceServer
	categoryQuery           service.CategoryQueryService
	categoryCommand         service.CategoryCommandService
	categoryStats           service.CategoryStatsService
	categoryStatsById       service.CategoryStatsByIdService
	categoryStatsByMerchant service.CategoryStatsByMerchantService
	mapping                 protomapper.CategoryProtoMapper
}

func NewCategoryHandleGrpc(
	service service.Service,
) *categoryHandleGrpc {
	return &categoryHandleGrpc{
		categoryQuery:           service.CategoryQuery,
		categoryCommand:         service.CategoryCommand,
		categoryStats:           service.CategoryStats,
		categoryStatsById:       service.CategoryStatsById,
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

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	category, totalRecords, err := s.categoryQuery.FindAll(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationCategory(paginationMeta, "success", "Successfully fetched categories", category)
	return so, nil
}

func (s *categoryHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategory, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := s.categoryQuery.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

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

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.categoryQuery.FindByActive(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationCategoryDeleteAt(paginationMeta, "success", "Successfully fetched active categories", users)

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

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationCategoryDeleteAt(paginationMeta, "success", "Successfully fetched trashed categories", categories)

	return so, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPrices(ctx context.Context, req *pb.FindYearMonthTotalPrices) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, category_errors.ErrGrpcCategoryInvalidMonth
	}

	reqService := requests.MonthTotalPrice{
		Year:  year,
		Month: month,
	}

	methods, err := s.categoryStats.FindMonthlyTotalPrice(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthlyTotalPrice("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPrices(ctx context.Context, req *pb.FindYearTotalPrices) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	methods, err := s.categoryStats.FindYearlyTotalPrice(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearlyTotalPrice("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesById(ctx context.Context, req *pb.FindYearMonthTotalPriceById) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month >= 12 {
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

	methods, err := s.categoryStatsById.FindMonthlyTotalPriceById(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthlyTotalPrice("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesById(ctx context.Context, req *pb.FindYearTotalPriceById) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
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

	methods, err := s.categoryStatsById.FindYearlyTotalPriceById(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearlyTotalPrice("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearMonthTotalPriceByMerchant) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	if month <= 0 || month >= 12 {
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

	methods, err := s.categoryStatsByMerchant.FindMonthlyTotalPriceByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthlyTotalPrice("success", "Monthly sales retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearTotalPriceByMerchant) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
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

	methods, err := s.categoryStatsByMerchant.FindYearlyTotalPriceByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearlyTotalPrice("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryMonthPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	methods, err := s.categoryStats.FindMonthPrice(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseCategoryMonthlyPrice("success", "Monthly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryYearPrice, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidYear
	}

	methods, err := s.categoryStats.FindYearPrice(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseCategoryYearlyPrice("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryMonthPrice, error) {
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

	methods, err := s.categoryStatsByMerchant.FindMonthPriceByMerchant(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseCategoryMonthlyPrice("success", "Merchant monthly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryYearPrice, error) {
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

	methods, err := s.categoryStatsByMerchant.FindYearPriceByMerchant(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseCategoryYearlyPrice("success", "Merchant yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindMonthPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryMonthPrice, error) {
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

	methods, err := s.categoryStatsById.FindMonthPriceById(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseCategoryMonthlyPrice("success", "Merchant monthly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) FindYearPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryYearPrice, error) {
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

	methods, err := s.categoryStatsById.FindYearPriceById(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseCategoryYearlyPrice("success", "Merchant yearly payment methods retrieved successfully", methods), nil
}

func (s *categoryHandleGrpc) Create(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.ApiResponseCategory, error) {
	req := &requests.CreateCategoryRequest{
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := s.categoryCommand.CreateCategory(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCategory("success", "Successfully created category", category)
	return so, nil
}

func (s *categoryHandleGrpc) Update(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ApiResponseCategory, error) {
	id := int(request.GetCategoryId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	req := &requests.UpdateCategoryRequest{
		CategoryID:    &id,
		Name:          request.GetName(),
		Description:   request.GetDescription(),
		ImageCategory: request.GetImageCategory(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := s.categoryCommand.UpdateCategory(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCategory("success", "Successfully updated category", category)
	return so, nil
}

func (s *categoryHandleGrpc) TrashedCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := s.categoryCommand.TrashedCategory(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCategoryDeleteAt("success", "Successfully trashed category", category)

	return so, nil
}

func (s *categoryHandleGrpc) RestoreCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	category, err := s.categoryCommand.RestoreCategory(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCategoryDeleteAt("success", "Successfully restored category", category)

	return so, nil
}

func (s *categoryHandleGrpc) DeleteCategoryPermanent(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, category_errors.ErrGrpcCategoryInvalidId
	}

	_, err := s.categoryCommand.DeleteCategoryPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCategoryDelete("success", "Successfully deleted category permanently")

	return so, nil
}

func (s *categoryHandleGrpc) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	_, err := s.categoryCommand.RestoreAllCategories()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCategoryAll("success", "Successfully restore all category")

	return so, nil
}

func (s *categoryHandleGrpc) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	_, err := s.categoryCommand.DeleteAllCategoriesPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseCategoryAll("success", "Successfully delete category permanen")

	return so, nil
}
