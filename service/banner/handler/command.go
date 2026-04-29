package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-banner/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/banner_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bannerCommandHandler struct {
	pb.UnimplementedBannerCommandServiceServer
	BannerCommand service.BannerCommandService
	logger        logger.LoggerInterface
}

func NewBannerCommandHandler(svc service.BannerCommandService, logger logger.LoggerInterface) BannerCommandHandler {
	return &bannerCommandHandler{
		BannerCommand: svc,
		logger:        logger,
	}
}

func (s *bannerCommandHandler) Create(ctx context.Context, request *pb.CreateBannerRequest) (*pb.ApiResponseBanner, error) {
	req := &requests.CreateBannerRequest{
		Name:      request.GetName(),
		StartDate: request.GetStartDate(),
		EndDate:   request.GetEndDate(),
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
		IsActive:  request.GetIsActive(),
	}

	if err := req.Validate(); err != nil {
		return nil, banner_errors.ErrGrpcValidateCreateBanner
	}

	banner, err := s.BannerCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBanner{
		Status:  "success",
		Message: "Successfully created banner",
		Data:    mapToProtoBannerResponse(banner),
	}, nil
}

func (s *bannerCommandHandler) Update(ctx context.Context, request *pb.UpdateBannerRequest) (*pb.ApiResponseBanner, error) {
	id := int(request.GetBannerId())
	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	req := &requests.UpdateBannerRequest{
		BannerID:  &id,
		Name:      request.GetName(),
		StartDate: request.GetStartDate(),
		EndDate:   request.GetEndDate(),
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
		IsActive:  request.GetIsActive(),
	}

	if err := req.Validate(); err != nil {
		return nil, banner_errors.ErrGrpcValidateUpdateBanner
	}

	banner, err := s.BannerCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBanner{
		Status:  "success",
		Message: "Successfully updated banner",
		Data:    mapToProtoBannerResponse(banner),
	}, nil
}

func (s *bannerCommandHandler) Trash(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	banner, err := s.BannerCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerDeleteAt{
		Status:  "success",
		Message: "Successfully trashed banner",
		Data:    mapToProtoBannerResponseDeleteAt(banner),
	}, nil
}

func (s *bannerCommandHandler) Restore(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	banner, err := s.BannerCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerDeleteAt{
		Status:  "success",
		Message: "Successfully restored banner",
		Data:    mapToProtoBannerResponseDeleteAt(banner),
	}, nil
}

func (s *bannerCommandHandler) DeletePermanent(ctx context.Context, request *pb.FindByIdBannerRequest) (*pb.ApiResponseBannerDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, banner_errors.ErrGrpcBannerInvalidId
	}

	_, err := s.BannerCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerDelete{
		Status:  "success",
		Message: "Successfully deleted banner permanently",
	}, nil
}

func (s *bannerCommandHandler) RestoreAll(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	_, err := s.BannerCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerAll{
		Status:  "success",
		Message: "Successfully restored all banners",
	}, nil
}

func (s *bannerCommandHandler) DeleteAll(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseBannerAll, error) {
	_, err := s.BannerCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseBannerAll{
		Status:  "success",
		Message: "Successfully deleted all banners permanently",
	}, nil
}
