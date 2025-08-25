package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingAddressHandleGrpc struct {
	pb.UnimplementedShippingServiceServer
	shippingAddressQueryService   service.ShippingAddressQueryService
	shippingAddressCommandService service.ShippingAddressCommandService
	logger                        logger.LoggerInterface
	mapping                       protomapper.ShippingAddresProtoMapper
}

func NewShippingAddressHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.ShippingServiceServer {
	return &shippingAddressHandleGrpc{
		shippingAddressQueryService:   service.ShippingAddressQuery,
		shippingAddressCommandService: service.ShippingAddressCommand,
		logger:                        logger,
		mapping:                       protomapper.NewShippingAddressProtoMapper(),
	}
}

func (s *shippingAddressHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShipping, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all shipping addresses",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	addresses, totalRecords, err := s.shippingAddressQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all shipping addresses",
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

	s.logger.Info("Successfully fetched all shipping addresses",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_addresses_count", len(addresses)),
	)

	so := s.mapping.ToProtoResponsePaginationShippingAddress(paginationMeta, "success", "Successfully fetched shipping addresses", addresses)
	return so, nil
}

func (s *shippingAddressHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShipping, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid shipping address ID provided", zap.Int("address_id", id))
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Fetching shipping address by ID", zap.Int("address_id", id))

	address, err := s.shippingAddressQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch shipping address by ID",
			zap.Int("address_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched shipping address by ID",
		zap.Int("address_id", id),
		zap.String("city", address.Kota),
	)

	so := s.mapping.ToProtoResponseShippingAddress("success", "Successfully fetched shipping address", address)
	return so, nil
}

func (s *shippingAddressHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active shipping addresses",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	addresses, totalRecords, err := s.shippingAddressQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active shipping addresses",
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

	s.logger.Info("Successfully fetched active shipping addresses",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_addresses_count", len(addresses)),
	)

	so := s.mapping.ToProtoResponsePaginationShippingAddressDeleteAt(paginationMeta, "success", "Successfully fetched active shipping addresses", addresses)
	return so, nil
}

func (s *shippingAddressHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllShippingRequest) (*pb.ApiResponsePaginationShippingDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed shipping addresses",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllShippingAddress{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	addresses, totalRecords, err := s.shippingAddressQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed shipping addresses",
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

	s.logger.Info("Successfully fetched trashed shipping addresses",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_addresses_count", len(addresses)),
	)

	so := s.mapping.ToProtoResponsePaginationShippingAddressDeleteAt(paginationMeta, "success", "Successfully fetched trashed shipping addresses", addresses)
	return so, nil
}

func (s *shippingAddressHandleGrpc) TrashedShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid shipping address ID for trashing", zap.Int("address_id", id))
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Moving shipping address to trash", zap.Int("address_id", id))

	address, err := s.shippingAddressCommandService.TrashShippingAddress(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash shipping address",
			zap.Int("address_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Shipping address moved to trash successfully",
		zap.Int("address_id", id),
		zap.String("city", address.Kota),
	)

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully trashed shipping address", address)
	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid shipping address ID for restore", zap.Int("address_id", id))
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Restoring shipping address from trash", zap.Int("address_id", id))

	address, err := s.shippingAddressCommandService.RestoreShippingAddress(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore shipping address",
			zap.Int("address_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Shipping address restored successfully",
		zap.Int("address_id", id),
	)

	so := s.mapping.ToProtoResponseShippingAddressDeleteAt("success", "Successfully restored shipping address", address)
	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingPermanent(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid shipping address ID for permanent deletion", zap.Int("address_id", id))
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Permanently deleting shipping address", zap.Int("address_id", id))

	_, err := s.shippingAddressCommandService.DeleteShippingAddressPermanently(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete shipping address",
			zap.Int("address_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Shipping address permanently deleted", zap.Int("address_id", id))

	so := s.mapping.ToProtoResponseShippingAddressDelete("success", "Successfully deleted shipping address permanently")
	return so, nil
}

func (s *shippingAddressHandleGrpc) RestoreAllShipping(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	s.logger.Info("Restoring all trashed shipping addresses")

	_, err := s.shippingAddressCommandService.RestoreAllShippingAddress(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all shipping addresses", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All shipping addresses restored successfully")

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully restored all shipping addresses")
	return so, nil
}

func (s *shippingAddressHandleGrpc) DeleteShippingAddressPermanently(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	s.logger.Info("Permanently deleting all trashed shipping addresses")

	_, err := s.shippingAddressCommandService.DeleteAllPermanentShippingAddress(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all shipping addresses", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All shipping addresses permanently deleted")

	so := s.mapping.ToProtoResponseShippingAddressAll("success", "Successfully deleted all shipping addresses permanently")
	return so, nil
}
