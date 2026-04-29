package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type shippingCommandHandler struct {
	pb.UnimplementedShippingCommandServiceServer
	shippingCommand service.ShippingAddressCommandService
	logger          logger.LoggerInterface
}

func NewShippingCommandHandler(svc service.ShippingAddressCommandService, logger logger.LoggerInterface) pb.ShippingCommandServiceServer {
	return &shippingCommandHandler{
		shippingCommand: svc,
		logger:          logger,
	}
}

func (s *shippingCommandHandler) CreateShipping(ctx context.Context, request *pb.CreateShippingAddressRequest) (*pb.ApiResponseShipping, error) {
	orderID := int(request.OrderId)
	req := &requests.CreateShippingAddressRequest{
		OrderID:        &orderID,
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   int(request.ShippingCost),
	}

	shipping, err := s.shippingCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShipping{
		Status:  "success",
		Message: "Successfully created shipping address",
		Data:    mapToProtoShippingResponse(shipping),
	}, nil
}

func (s *shippingCommandHandler) UpdateShipping(ctx context.Context, request *pb.UpdateShippingAddressRequest) (*pb.ApiResponseShipping, error) {
	shippingID := int(request.ShippingId)
	orderID := int(request.OrderId)
	req := &requests.UpdateShippingAddressRequest{
		ShippingID:     &shippingID,
		OrderID:        &orderID,
		Alamat:         request.Alamat,
		Provinsi:       request.Provinsi,
		Kota:           request.Kota,
		Negara:         request.Negara,
		Courier:        request.Courier,
		ShippingMethod: request.ShippingMethod,
		ShippingCost:   int(request.ShippingCost),
	}

	shipping, err := s.shippingCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShipping{
		Status:  "success",
		Message: "Successfully updated shipping address",
		Data:    mapToProtoShippingResponse(shipping),
	}, nil
}

func (s *shippingCommandHandler) TrashedShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingDeleteAt{
		Status:  "success",
		Message: "Successfully trashed shipping address",
		Data:    mapToProtoShippingResponseDeleteAt(shipping),
	}, nil
}

func (s *shippingCommandHandler) RestoreShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingDeleteAt{
		Status:  "success",
		Message: "Successfully restored shipping address",
		Data:    mapToProtoShippingResponseDeleteAt(shipping),
	}, nil
}

func (s *shippingCommandHandler) DeleteShippingPermanent(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	_, err := s.shippingCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingDelete{
		Status:  "success",
		Message: "Successfully deleted shipping address permanently",
	}, nil
}

func (s *shippingCommandHandler) RestoreAllShipping(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingAll{
		Status:  "success",
		Message: "Successfully restored all shipping addresses",
	}, nil
}

func (s *shippingCommandHandler) DeleteAllShippingPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingAll{
		Status:  "success",
		Message: "Successfully deleted all shipping addresses permanently",
	}, nil
}

func (s *shippingCommandHandler) DeleteShippingByOrderPermanent(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	_, err := s.shippingCommand.DeleteShippingAddressByOrderPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingDelete{
		Status:  "success",
		Message: "Successfully deleted shipping addresses by order permanently",
	}, nil
}
