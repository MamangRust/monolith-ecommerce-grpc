package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-merchant_detail/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"

	pbmerchant_social_link "github.com/MamangRust/monolith-ecommerce-pb/merchant_social_link"
)

type merchantSocialLinkCommandHandler struct {
	pbmerchant_social_link.UnimplementedMerchantSocialCommandServiceServer
	MerchantSocialLinkCommand service.MerchantSocialLinkCommandService
	logger                    logger.LoggerInterface
}

func NewMerchantSocialLinkCommandHandler(svc service.MerchantSocialLinkCommandService, logger logger.LoggerInterface) MerchantSocialLinkCommandHandler {
	return &merchantSocialLinkCommandHandler{
		MerchantSocialLinkCommand: svc,
		logger:                    logger,
	}
}

func (s *merchantSocialLinkCommandHandler) Create(ctx context.Context, request *pbmerchant_social_link.CreateMerchantSocialRequest) (*pbmerchant_social_link.ApiResponseMerchantSocial, error) {
	merchantDetailID := int(request.GetMerchantDetailId())
	req := &requests.CreateMerchantSocialRequest{
		MerchantDetailID: &merchantDetailID,
		Platform:         request.GetPlatform(),
		Url:              request.GetUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateCreateMerchantDetail
	}

	link, err := s.MerchantSocialLinkCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_social_link.ApiResponseMerchantSocial{
		Status:  "success",
		Message: "Successfully created merchant social link",
		Data:    mapToProtoMerchantSocialLinkResponse(link),
	}, nil
}

func (s *merchantSocialLinkCommandHandler) Update(ctx context.Context, request *pbmerchant_social_link.UpdateMerchantSocialRequest) (*pbmerchant_social_link.ApiResponseMerchantSocial, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	merchantDetailID := int(request.GetMerchantDetailId())
	req := &requests.UpdateMerchantSocialRequest{
		ID:               id,
		MerchantDetailID: &merchantDetailID,
		Platform:         request.GetPlatform(),
		Url:              request.GetUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateUpdateMerchantDetail
	}

	link, err := s.MerchantSocialLinkCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbmerchant_social_link.ApiResponseMerchantSocial{
		Status:  "success",
		Message: "Successfully updated merchant social link",
		Data:    mapToProtoMerchantSocialLinkResponse(link),
	}, nil
}
