package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-merchant_policy/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"

	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	pbmerchant_policy "github.com/MamangRust/monolith-ecommerce-pb/merchant_policy"
)

type merchantPolicyQueryHandler struct {
	pbmerchant_policy.UnimplementedMerchantPolicyQueryServiceServer
	merchantPolicyService service.MerchantPoliciesQueryService
	logger                logger.LoggerInterface
}

func NewMerchantPolicyQueryHandler(
	merchantPolicyService service.MerchantPoliciesQueryService,
	logger logger.LoggerInterface,
) pbmerchant_policy.MerchantPolicyQueryServiceServer {
	return &merchantPolicyQueryHandler{
		merchantPolicyService: merchantPolicyService,
		logger:                logger,
	}
}

func (h *merchantPolicyQueryHandler) FindAll(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant_policy.ApiResponsePaginationMerchantPolicies, error) {
	merchants, total, err := h.merchantPolicyService.FindAll(ctx, &requests.FindAllMerchant{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	})

	if err != nil {
		return nil, err
	}

	return mapToPaginationResponse(merchants, total), nil
}

func (h *merchantPolicyQueryHandler) FindById(ctx context.Context, req *pbmerchant_policy.FindByIdMerchantPoliciesRequest) (*pbmerchant_policy.ApiResponseMerchantPolicies, error) {
	merchant, err := h.merchantPolicyService.FindByID(ctx, int(req.GetId()))

	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return mapToSingleResponse(merchant), nil
}

func (h *merchantPolicyQueryHandler) FindByActive(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant_policy.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
	merchants, total, err := h.merchantPolicyService.FindActive(ctx, &requests.FindAllMerchant{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	})

	if err != nil {
		return nil, err
	}

	return mapToPaginationDeleteAtResponse(merchants, total), nil
}

func (h *merchantPolicyQueryHandler) FindByTrashed(ctx context.Context, req *pbmerchant.FindAllMerchantRequest) (*pbmerchant_policy.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
	merchants, total, err := h.merchantPolicyService.FindTrashed(ctx, &requests.FindAllMerchant{
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
		Search:   req.GetSearch(),
	})

	if err != nil {
		return nil, err
	}

	return mapToPaginationDeleteAtResponse(merchants, total), nil
}
