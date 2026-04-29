package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/product_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productCommandHandler struct {
	pb.UnimplementedProductCommandServiceServer
	productService service.ProductCommandService
	logger         logger.LoggerInterface
}

func NewProductCommandHandler(productService service.ProductCommandService, logger logger.LoggerInterface) *productCommandHandler {
	return &productCommandHandler{
		productService: productService,
		logger:         logger,
	}
}

func (h *productCommandHandler) Create(ctx context.Context, request *pb.CreateProductRequest) (*pb.ApiResponseProduct, error) {
	rating := int(request.GetRating())
	slug := request.GetSlugProduct()

	req := &requests.CreateProductRequest{
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		Rating:       &rating,
		SlugProduct:  &slug,
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateCreateProduct
	}

	product, err := h.productService.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully created product",
		Data:    mapToProtoProductResponse(product),
	}, nil
}

func (h *productCommandHandler) Update(ctx context.Context, request *pb.UpdateProductRequest) (*pb.ApiResponseProduct, error) {
	id := int(request.GetProductId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	rating := int(request.GetRating())
	slug := request.GetSlugProduct()

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
		Rating:       &rating,
		SlugProduct:  &slug,
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateUpdateProduct
	}

	product, err := h.productService.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully updated product",
		Data:    mapToProtoProductResponse(product),
	}, nil
}

func (h *productCommandHandler) TrashedProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := h.productService.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully trashed product",
		Data:    mapToProtoProductResponseDeleteAt(product),
	}, nil
}

func (h *productCommandHandler) RestoreProduct(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := h.productService.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully restored product",
		Data:    mapToProtoProductResponseDeleteAt(product),
	}, nil
}

func (h *productCommandHandler) DeleteProductPermanent(ctx context.Context, request *pb.FindByIdProductRequest) (*pb.ApiResponseProductDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	_, err := h.productService.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductDelete{
		Status:  "success",
		Message: "Successfully deleted product permanently",
	}, nil
}

func (h *productCommandHandler) RestoreAllProduct(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	_, err := h.productService.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully restored all products",
	}, nil
}

func (h *productCommandHandler) DeleteAllProductPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseProductAll, error) {
	_, err := h.productService.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully deleted all products permanently",
	}, nil
}

func (h *productCommandHandler) UpdateProductCountStock(ctx context.Context, request *pb.UpdateProductCountStockRequest) (*pb.ApiResponseProduct, error) {
	product, err := h.productService.UpdateProductCountStock(ctx, int(request.ProductId), int(request.Stock))
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully updated product count stock",
		Data:    mapToProtoProductResponse(product),
	}, nil
}
