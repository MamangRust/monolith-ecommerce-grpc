package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantPolicyCommandHandler struct {
	pb.UnimplementedMerchantPolicyCommandServiceServer
	merchantPolicyService service.MerchantPoliciesCommandService
	logger                logger.LoggerInterface
}

func NewMerchantPolicyCommandHandler(
	merchantPolicyService service.MerchantPoliciesCommandService,
	logger logger.LoggerInterface,
) pb.MerchantPolicyCommandServiceServer {
	return &merchantPolicyCommandHandler{
		merchantPolicyService: merchantPolicyService,
		logger:                logger,
	}
}

func (h *merchantPolicyCommandHandler) Create(ctx context.Context, req *pb.CreateMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	policy, err := h.merchantPolicyService.Create(ctx, &requests.CreateMerchantPolicyRequest{
		MerchantID:  int(req.GetMerchantId()),
		PolicyType:  req.GetPolicyType(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	})

	if err != nil {
		return nil, err
	}

	return mapToSingleResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) Update(ctx context.Context, req *pb.UpdateMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	id := int(req.GetMerchantPolicyId())
	policy, err := h.merchantPolicyService.Update(ctx, &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &id,
		PolicyType:       req.GetPolicyType(),
		Title:            req.GetTitle(),
		Description:      req.GetDescription(),
	})

	if err != nil {
		return nil, err
	}

	return mapToSingleResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) TrashedMerchantPolicies(ctx context.Context, req *pb.FindByIdMerchantPoliciesRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	policy, err := h.merchantPolicyService.Trash(ctx, int(req.GetId()))

	if err != nil {
		return nil, err
	}

	return mapToSingleDeleteAtResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) RestoreMerchantPolicies(ctx context.Context, req *pb.FindByIdMerchantPoliciesRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	policy, err := h.merchantPolicyService.Restore(ctx, int(req.GetId()))

	if err != nil {
		return nil, err
	}

	return mapToSingleDeleteAtResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) DeleteMerchantPoliciesPermanent(ctx context.Context, req *pb.FindByIdMerchantPoliciesRequest) (*pb.ApiResponseMerchantDelete, error) {
	_, err := h.merchantPolicyService.DeletePermanent(ctx, int(req.GetId()))

	if err != nil {
		return nil, err
	}

	return &pb.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant policy permanently",
	}, nil
}

func (h *merchantPolicyCommandHandler) RestoreAllMerchantPolicies(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := h.merchantPolicyService.RestoreAll(ctx)

	if err != nil {
		return nil, err
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all merchant policies",
	}, nil
}

func (h *merchantPolicyCommandHandler) DeleteAllMerchantPoliciesPermanent(ctx context.Context, req *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := h.merchantPolicyService.DeleteAll(ctx)

	if err != nil {
		return nil, err
	}

	return &pb.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchant policies permanently",
	}, nil
}
