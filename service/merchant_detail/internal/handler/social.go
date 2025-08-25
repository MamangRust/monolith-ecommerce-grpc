package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantsocial_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_social_link_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
)

type merchanSocialLinkHandleGrpc struct {
	pb.UnimplementedMerchantSocialServiceServer
	merchantSocialLink service.MerchantSocialLinkService
	mapping            protomapper.MerchantSocialLinkProtoMapper
	logger             logger.LoggerInterface
}

func NewMerchantSocialLinkHandleGrpc(
	service *service.Service,
	mapping protomapper.MerchantSocialLinkProtoMapper,
	logger logger.LoggerInterface,
) pb.MerchantSocialServiceServer {
	return &merchanSocialLinkHandleGrpc{
		merchantSocialLink: service.MerchantSocialLink,
		logger:             logger,
		mapping:            mapping,
	}
}

func (s *merchanSocialLinkHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantSocialRequest) (*pb.ApiResponseMerchantSocial, error) {
	s.logger.Info("Creating merchant social link",
		zap.Int("merchant_detail_id", int(request.GetMerchantDetailId())),
		zap.String("platform", request.GetPlatform()),
		zap.String("url", request.GetUrl()),
	)

	id := int(request.GetMerchantDetailId())

	req := &requests.CreateMerchantSocialRequest{
		MerchantDetailID: &id,
		Platform:         request.GetPlatform(),
		Url:              request.GetUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant social link creation",
			zap.Int("merchant_detail_id", int(request.GetMerchantDetailId())),
			zap.String("platform", request.GetPlatform()),
			zap.Error(err),
		)
		return nil, merchantsocial_errors.ErrGrpcValidateCreateMerchantSocialLink
	}

	social, err := s.merchantSocialLink.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create merchant social link",
			zap.Int("merchant_detail_id", int(request.GetMerchantDetailId())),
			zap.String("platform", request.GetPlatform()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant social link created successfully",
		zap.Int("social_id", int(social.ID)),
		zap.Int("merchant_detail_id", int(social.MerchantDetailID)),
		zap.String("platform", social.Platform),
		zap.String("url", social.URL),
	)

	so := s.mapping.ToProtoResponseMerchantSocialLink("success", "Successfully created merchant social link", social)
	return so, nil
}

func (s *merchanSocialLinkHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantSocialRequest) (*pb.ApiResponseMerchantSocial, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant social link ID provided for update", zap.Int("social_id", id))
		return nil, merchantsocial_errors.ErrGrpcMerchantSocialLinkInvalidID
	}

	s.logger.Info("Updating merchant social link", zap.Int("social_id", id))

	id_merchant := int(request.GetMerchantDetailId())

	req := &requests.UpdateMerchantSocialRequest{
		ID:               id,
		MerchantDetailID: &id_merchant,
		Platform:         request.GetPlatform(),
		Url:              request.GetUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant social link update",
			zap.Int("social_id", id),
			zap.String("platform", request.GetPlatform()),
			zap.Error(err),
		)
		return nil, merchantsocial_errors.ErrGrpcFailedUpdateMerchantSocialLink
	}

	social, err := s.merchantSocialLink.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update merchant social link",
			zap.Int("social_id", id),
			zap.String("platform", request.GetPlatform()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant social link updated successfully",
		zap.Int("social_id", id),
		zap.Int("merchant_detail_id", int(social.MerchantDetailID)),
		zap.String("platform", social.Platform),
		zap.String("url", social.URL),
	)

	so := s.mapping.ToProtoResponseMerchantSocialLink("success", "Successfully updated merchant social link", social)
	return so, nil
}
