package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productHandleGrpc struct {
	pb.UnimplementedProductServiceServer
	logger                logger.LoggerInterface
	productQueryService   service.ProductQueryService
	productCommandService service.ProductCommandService
	mapping               protomapper.ProductProtoMapper
}

func NewProductHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.ProductServiceServer {
	return &productHandleGrpc{
		logger:                logger,
		productQueryService:   service.ProductQuery,
		productCommandService: service.ProductCommand,
		mapping:               protomapper.NewProductProtoMapper(),
	}
}

func (s *productHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProduct, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all products",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all products",
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

	s.logger.Info("Successfully fetched all products",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_products_count", len(products)),
	)

	so := s.mapping.ToProtoResponsePaginationProduct(paginationMeta, "success", "Successfully fetched products", products)
	return so, nil
}

func (s *productHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllProductMerchantRequest) (*pb.ApiResponsePaginationProduct, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchantID := int(request.GetMerchantId())
	minPrice := int(request.GetMinPrice())
	maxPrice := int(request.GetMaxPrice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if minPrice <= 0 {
		minPrice = 0
	}
	if maxPrice <= 0 {
		maxPrice = 0
	}

	s.logger.Info("Fetching products by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
		zap.Int("min_price", minPrice),
		zap.Int("max_price", maxPrice),
	)

	reqService := requests.FindAllProductByMerchant{
		MerchantID: merchantID,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
	}

	products, totalRecords, err := s.productQueryService.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch products by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
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

	s.logger.Info("Successfully fetched products by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int("fetched_products_count", len(products)),
	)

	so := s.mapping.ToProtoResponsePaginationProduct(paginationMeta, "success", "Successfully fetched products", products)
	return so, nil
}

func (s *productHandleGrpc) FindByCategory(ctx context.Context, request *pb.FindAllProductCategoryRequest) (*pb.ApiResponsePaginationProduct, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	categoryName := request.GetCategoryName()
	minPrice := int(request.GetMinPrice())
	maxPrice := int(request.GetMaxPrice())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if minPrice <= 0 {
		minPrice = 0
	}
	if maxPrice <= 0 {
		maxPrice = 0
	}

	s.logger.Info("Fetching products by category",
		zap.String("category_name", categoryName),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
		zap.Int("min_price", minPrice),
		zap.Int("max_price", maxPrice),
	)

	reqService := requests.FindAllProductByCategory{
		Page:         page,
		PageSize:     pageSize,
		Search:       search,
		CategoryName: categoryName,
		MinPrice:     &minPrice,
		MaxPrice:     &maxPrice,
	}

	products, totalRecords, err := s.productQueryService.FindByCategory(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch products by category",
			zap.String("category_name", categoryName),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
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

	s.logger.Info("Successfully fetched products by category",
		zap.String("category_name", categoryName),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int("fetched_products_count", len(products)),
	)

	so := s.mapping.ToProtoResponsePaginationProduct(paginationMeta, "success", "Successfully fetched products", products)
	return so, nil
}

func (s *productHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProduct, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid product ID provided", zap.Int("product_id", id))
		return nil, product_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Fetching product by ID", zap.Int("product_id", id))

	product, err := s.productQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch product by ID",
			zap.Int("product_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched product by ID",
		zap.Int("product_id", id),
		zap.String("product_name", product.Name),
		zap.Int("merchant_id", int(product.MerchantID)),
	)

	so := s.mapping.ToProtoResponseProduct("success", "Successfully fetched product", product)
	return so, nil
}

func (s *productHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active products",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active products",
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

	s.logger.Info("Successfully fetched active products",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int("fetched_products_count", len(products)),
	)

	so := s.mapping.ToProtoResponsePaginationProductDeleteAt(paginationMeta, "success", "Successfully fetched active products", products)
	return so, nil
}

func (s *productHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllProductRequest) (*pb.ApiResponsePaginationProductDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed products",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllProduct{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	products, totalRecords, err := s.productQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed products",
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

	s.logger.Info("Successfully fetched trashed products",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int("fetched_products_count", len(products)),
	)

	so := s.mapping.ToProtoResponsePaginationProductDeleteAt(paginationMeta, "success", "Successfully fetched trashed products", products)
	return so, nil
}

func (s *productHandleGrpc) Create(ctx context.Context, request *pb.CreateProductRequest) (*pb.ApiResponseProduct, error) {
	s.logger.Info("Creating new product",
		zap.Int("merchant_id", int(request.GetMerchantId())),
		zap.Int("category_id", int(request.GetCategoryId())),
		zap.String("name", request.GetName()),
		zap.Int("price", int(request.GetPrice())),
		zap.Int("count_in_stock", int(request.GetCountInStock())),
	)

	rating := int(request.GetRating())

	req := &requests.CreateProductRequest{
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		SlugProduct:  &request.SlugProduct,
		Rating:       &rating,
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on product creation",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, product_errors.ErrGrpcValidateCreateProduct
	}

	product, err := s.productCommandService.CreateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create product",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Product created successfully",
		zap.Int("product_id", int(product.ID)),
		zap.String("name", product.Name),
		zap.Int("price", product.Price),
	)

	so := s.mapping.ToProtoResponseProduct("success", "Successfully created product", product)
	return so, nil
}

func (s *productHandleGrpc) Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.ApiResponseProduct, error) {
	id := int(request.GetProductId())

	if id == 0 {
		s.logger.Error("Invalid product ID provided for update", zap.Int("product_id", id))
		return nil, product_errors.ErrGrpcInvalidID
	}

	rating := int(request.GetRating())

	s.logger.Info("Updating product", zap.Int("product_id", id))

	req := &requests.UpdateProductRequest{
		ProductID:    &id,
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		SlugProduct:  &request.SlugProduct,
		Rating:       &rating,
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on product update",
			zap.Int("product_id", id),
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, product_errors.ErrGrpcValidateUpdateProduct
	}

	product, err := s.productCommandService.UpdateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update product",
			zap.Int("product_id", id),
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Product updated successfully",
		zap.Int("product_id", id),
		zap.String("name", product.Name),
		zap.Int("price", product.Price),
	)

	so := s.mapping.ToProtoResponseProduct("success", "Successfully updated product", product)
	return so, nil
}

func (s *productHandleGrpc) TrashedProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid product ID for trashing", zap.Int("product_id", id))
		return nil, product_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Moving product to trash", zap.Int("product_id", id))

	product, err := s.productCommandService.TrashProduct(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash product",
			zap.Int("product_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Product moved to trash successfully",
		zap.Int("product_id", id),
		zap.String("name", product.Name),
		zap.Int("merchant_id", int(product.MerchantID)),
	)

	so := s.mapping.ToProtoResponseProductDeleteAt("success", "Successfully trashed product", product)
	return so, nil
}

func (s *productHandleGrpc) RestoreProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid product ID for restore", zap.Int("product_id", id))
		return nil, product_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Restoring product from trash", zap.Int("product_id", id))

	product, err := s.productCommandService.RestoreProduct(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore product",
			zap.Int("product_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Product restored successfully",
		zap.Int("product_id", id),
		zap.String("name", product.Name),
	)

	so := s.mapping.ToProtoResponseProductDeleteAt("success", "Successfully restored product", product)
	return so, nil
}

func (s *productHandleGrpc) DeleteProductPermanent(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid product ID for permanent deletion", zap.Int("product_id", id))
		return nil, product_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Permanently deleting product", zap.Int("product_id", id))

	_, err := s.productCommandService.DeleteProductPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete product",
			zap.Int("product_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Product permanently deleted", zap.Int("product_id", id))

	so := s.mapping.ToProtoResponseProductDelete("success", "Successfully deleted product permanently")
	return so, nil
}

func (s *productHandleGrpc) RestoreAllProduct(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	s.logger.Info("Restoring all trashed products")

	_, err := s.productCommandService.RestoreAllProducts(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all products", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All products restored successfully")

	so := s.mapping.ToProtoResponseProductAll("success", "Successfully restored all products")
	return so, nil
}

func (s *productHandleGrpc) DeleteAllProductPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	s.logger.Info("Permanently deleting all trashed products")

	_, err := s.productCommandService.DeleteAllProductsPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all products", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All products permanently deleted")

	so := s.mapping.ToProtoResponseProductAll("success", "Successfully deleted all products permanently")
	return so, nil
}
