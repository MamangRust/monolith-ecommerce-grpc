package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-merchant_policy/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"google.golang.org/protobuf/types/known/emptypb"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbmerchant_policy "github.com/MamangRust/monolith-ecommerce-pb/merchant_policy"
)

type merchantPolicyCommandHandler struct {
	pbmerchant_policy.UnimplementedMerchantPolicyCommandServiceServer
	merchantPolicyService service.MerchantPoliciesCommandService
	logger                logger.LoggerInterface
}

func NewMerchantPolicyCommandHandler(
	merchantPolicyService service.MerchantPoliciesCommandService,
	logger logger.LoggerInterface,
) pbmerchant_policy.MerchantPolicyCommandServiceServer {
	return &merchantPolicyCommandHandler{
		merchantPolicyService: merchantPolicyService,
		logger:                logger,
	}
}

func (h *merchantPolicyCommandHandler) Create(ctx context.Context, req *pbmerchant_policy.CreateMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPolicies, error) {
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

func (h *merchantPolicyCommandHandler) Update(ctx context.Context, req *pbmerchant_policy.UpdateMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPolicies, error) {
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

func (h *merchantPolicyCommandHandler) TrashedMerchantPolicies(ctx context.Context, req *pbmerchant_policy.FindByIdMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt, error) {
	policy, err := h.merchantPolicyService.Trash(ctx, int(req.GetId()))

	if err != nil {
		return nil, err
	}

	return mapToSingleDeleteAtResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) RestoreMerchantPolicies(ctx context.Context, req *pbmerchant_policy.FindByIdMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt, error) {
	policy, err := h.merchantPolicyService.Restore(ctx, int(req.GetId()))

	if err != nil {
		return nil, err
	}

	return mapToSingleDeleteAtResponse(policy), nil
}

func (h *merchantPolicyCommandHandler) DeleteMerchantPoliciesPermanent(ctx context.Context, req *pbmerchant_policy.FindByIdMerchantPoliciesRequest) (*pbmerchant.ApiResponseMerchantDelete, error) {
	_, err := h.merchantPolicyService.DeletePermanent(ctx, int(req.GetId()))

	if err != nil {
		return nil, err
	}

	return &pbmerchant.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant policy permanently",
	}, nil
}

func (h *merchantPolicyCommandHandler) RestoreAllMerchantPolicies(ctx context.Context, req *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := h.merchantPolicyService.RestoreAll(ctx)

	if err != nil {
		return nil, err
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restored all merchant policies",
	}, nil
}

func (h *merchantPolicyCommandHandler) DeleteAllMerchantPoliciesPermanent(ctx context.Context, req *emptypb.Empty) (*pbmerchant.ApiResponseMerchantAll, error) {
	_, err := h.merchantPolicyService.DeleteAll(ctx)

	if err != nil {
		return nil, err
	}

	return &pbmerchant.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully deleted all merchant policies permanently",
	}, nil
}
