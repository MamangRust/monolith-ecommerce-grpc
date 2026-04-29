package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type productQueryHandler struct {
	pb.UnimplementedProductQueryServiceServer
	productService service.ProductQueryService
	logger         logger.LoggerInterface
}

func NewProductQueryHandler(productService service.ProductQueryService, logger logger.LoggerInterface) *productQueryHandler {
	return &productQueryHandler{
		productService: productService,
		logger:         logger,
	}
}

func (h *productQueryHandler) FindAll(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProduct, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := h.productService.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbProducts := make([]*pb.ProductResponse, len(products))
	for i, product := range products {
		pbProducts[i] = mapToProtoProductResponse(product)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched products",
		Data:       pbProducts,
		Pagination: paginationMeta,
	}, nil
}

func (h *productQueryHandler) FindByMerchant(ctx context.Context, request *pb.FindAllProductMerchantRequest) (*pb.ApiResponsePaginationProduct, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()
	merchantId := int(request.GetMerchantId())
	minPrice := int(request.GetMinPrice())
	maxPrice := int(request.GetMaxPrice())

	reqService := requests.FindAllProductByMerchant{
		MerchantID: merchantId,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
	}

	products, totalRecords, err := h.productService.FindByMerchant(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbProducts := make([]*pb.ProductResponse, len(products))
	for i, product := range products {
		pbProducts[i] = mapToProtoProductResponse(product)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched merchant products",
		Data:       pbProducts,
		Pagination: paginationMeta,
	}, nil
}

func (h *productQueryHandler) FindByCategory(ctx context.Context, request *pb.FindAllProductCategoryRequest) (*pb.ApiResponsePaginationProduct, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()
	categoryName := request.GetCategoryName()
	minPrice := int(request.GetMinPrice())
	maxPrice := int(request.GetMaxPrice())

	reqService := requests.FindAllProductByCategory{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		CategoryName: categoryName,
		MinPrice:     &minPrice,
		MaxPrice:     &maxPrice,
	}

	products, totalRecords, err := h.productService.FindByCategory(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbProducts := make([]*pb.ProductResponse, len(products))
	for i, product := range products {
		pbProducts[i] = mapToProtoProductResponse(product)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationProduct{
		Status:     "success",
		Message:    "Successfully fetched category products",
		Data:       pbProducts,
		Pagination: paginationMeta,
	}, nil
}

func (h *productQueryHandler) FindById(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProduct, error) {
	id := int(request.GetId())

	product, err := h.productService.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully fetched product",
		Data:    mapToProtoProductResponse(product),
	}, nil
}

func (h *productQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := h.productService.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbProducts := make([]*pb.ProductResponseDeleteAt, len(products))
	for i, product := range products {
		pbProducts[i] = mapToProtoProductResponseDeleteAt(product)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active products",
		Data:       pbProducts,
		Pagination: paginationMeta,
	}, nil
}

func (h *productQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	page, pageSize := normalizePage(int(request.GetPage()), int(request.GetPageSize()))
	search := request.GetSearch()

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := h.productService.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	pbProducts := make([]*pb.ProductResponseDeleteAt, len(products))
	for i, product := range products {
		pbProducts[i] = mapToProtoProductResponseDeleteAt(product)
	}

	paginationMeta := createPaginationMeta(page, pageSize, *totalRecords)

	return &pb.ApiResponsePaginationProductDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed products",
		Data:       pbProducts,
		Pagination: paginationMeta,
	}, nil
}
