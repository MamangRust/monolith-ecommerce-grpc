package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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

func (s *shippingCommandHandler) TrashedShipping(ctx context.Context, request *pb.FindByIdShippingRequest) (*pb.ApiResponseShippingDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, shippingaddress_errors.ErrGrpcInvalidID
	}

	shipping, err := s.shippingCommand.TrashShippingAddress(ctx, id)
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

	shipping, err := s.shippingCommand.RestoreShippingAddress(ctx, id)
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

	_, err := s.shippingCommand.DeleteShippingAddressPermanently(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingDelete{
		Status:  "success",
		Message: "Successfully deleted shipping address permanently",
	}, nil
}

func (s *shippingCommandHandler) RestoreAllShipping(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingCommand.RestoreAllShippingAddress(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingAll{
		Status:  "success",
		Message: "Successfully restored all shipping addresses",
	}, nil
}

func (s *shippingCommandHandler) DeleteAllShippingPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseShippingAll, error) {
	_, err := s.shippingCommand.DeleteAllPermanentShippingAddress(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseShippingAll{
		Status:  "success",
		Message: "Successfully deleted all shipping addresses permanently",
	}, nil
}
