package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-merchant_detail/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	"google.golang.org/protobuf/types/known/emptypb"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbmerchant_detail "github.com/MamangRust/monolith-ecommerce-pb/merchant_detail"
)

type merchantDetailCommandHandler struct {
	pbmerchant_detail.UnimplementedMerchantDetailCommandServiceServer
	MerchantDetailCommand service.MerchantDetailCommandService
	logger                logger.LoggerInterface
}

func NewMerchantDetailCommandHandler(svc service.MerchantDetailCommandService, logger logger.LoggerInterface) MerchantDetailCommandHandler {
	return &merchantDetailCommandHandler{
		MerchantDetailCommand: svc,
		logger:                logger,
	}
}

func (s *merchantDetailCommandHandler) Create(ctx context.Context, request *pbmerchant_detail.CreateMerchantDetailRequest) (*pbmerchant_detail.ApiResponseMerchantDetail, error) {
	req := &requests.CreateMerchantDetailRequest{
		MerchantID:       int(request.GetMerchantId()),
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateCreateMerchantDetail
	}

	detail, err := s.MerchantDetailCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_detail.ApiResponseMerchantDetail{
		Status:  "success",
		Message: "Successfully created merchant detail",
		Data:    mapToProtoMerchantDetailResponse(detail),
	}, nil
}

func (s *merchantDetailCommandHandler) Update(ctx context.Context, request *pbmerchant_detail.UpdateMerchantDetailRequest) (*pbmerchant_detail.ApiResponseMerchantDetail, error) {
	id := int(request.GetMerchantDetailId())
	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	req := &requests.UpdateMerchantDetailRequest{
		MerchantDetailID: &id,
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateUpdateMerchantDetail
	}

	detail, err := s.MerchantDetailCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_detail.ApiResponseMerchantDetail{
		Status:  "success",
		Message: "Successfully updated merchant detail",
		Data:    mapToProtoMerchantDetailResponse(detail),
	}, nil
}

func (s *merchantDetailCommandHandler) TrashedMerchantDetail(ctx context.Context, request *pbmerchant_detail.FindByIdMerchantDetailRequest) (*pbmerchant_detail.ApiResponseMerchantDetailDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	detail, err := s.MerchantDetailCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_detail.ApiResponseMerchantDetailDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant detail",
		Data:    mapToProtoMerchantDetailResponseDeleteAt(detail),
	}, nil
}

func (s *merchantDetailCommandHandler) RestoreMerchantDetail(ctx context.Context, request *pbmerchant_detail.FindByIdMerchantDetailRequest) (*pbmerchant_detail.ApiResponseMerchantDetailDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	detail, err := s.MerchantDetailCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_detail.ApiResponseMerchantDetailDeleteAt{
		Status:  "success",
		Message: "Successfully restored merchant detail",
		Data:    mapToProtoMerchantDetailResponseDeleteAt(detail),
	}, nil
}

func (s *merchantDetailCommandHandler) DeleteMerchantDetailPermanent(ctx context.Context, request *pbmerchant_detail.FindByIdMerchantDetailRequest) (*pbmerchant.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	_, err := s.MerchantDetailCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant detail permanently",
	}, nil
}

func (s *merchantDetailCommandHandler) RestoreAllMerchantDetail(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := s.MerchantDetailCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all merchant details",
	}, nil
}

func (s *merchantDetailCommandHandler) DeleteAllMerchantDetailPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := s.MerchantDetailCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchant details permanently",
	}, nil
}
