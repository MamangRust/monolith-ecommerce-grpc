package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbmerchant_policy "github.com/MamangRust/monolith-ecommerce-pb/merchant_policy"
)

func mapToSingleResponse(data interface{}) *pbmerchant_policy.ApiResponseMerchantPolicies {
	return &pbmerchant_policy.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policy",
		Data:    mapToMerchantPolicyResponse(data).(*pbmerchant_policy.MerchantPoliciesResponse),
	}
}

func mapToPaginationResponse(data []*db.GetMerchantPoliciesRow, total *int) *pbmerchant_policy.ApiResponsePaginationMerchantPolicies {
	var policies []*pbmerchant_policy.MerchantPoliciesResponse
	for _, v := range data {
		policies = append(policies, mapToMerchantPolicyResponse(v).(*pbmerchant_policy.MerchantPoliciesResponse))
	}

	return &pbmerchant_policy.ApiResponsePaginationMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pbcommon.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToPaginationDeleteAtResponse(data interface{}, total *int) *pbmerchant_policy.ApiResponsePaginationMerchantPoliciesDeleteAt {
	var policies []*pbmerchant_policy.MerchantPoliciesResponseDeleteAt

	switch v := data.(type) {
	case []*db.GetMerchantPoliciesActiveRow:
		for _, item := range v {
			policies = append(policies, mapToMerchantPolicyResponse(item).(*pbmerchant_policy.MerchantPoliciesResponseDeleteAt))
		}
	case []*db.GetMerchantPoliciesTrashedRow:
		for _, item := range v {
			policies = append(policies, mapToMerchantPolicyResponse(item).(*pbmerchant_policy.MerchantPoliciesResponseDeleteAt))
		}
	}

	return &pbmerchant_policy.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pbcommon.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToSingleDeleteAtResponse(data *db.MerchantPolicy) *pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt {
	return &pbmerchant_policy.ApiResponseMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully processed merchant policy",
		Data:    mapToMerchantPolicyResponse(data).(*pbmerchant_policy.MerchantPoliciesResponseDeleteAt),
	}
}

func mapToMerchantPolicyResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *db.GetMerchantPolicyRow:
		return &pbmerchant_policy.MerchantPoliciesResponse{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Time.String(),
			UpdatedAt:   v.UpdatedAt.Time.String(),
		}
	case *db.GetMerchantPoliciesRow:
		return &pbmerchant_policy.MerchantPoliciesResponse{
			Id:           int32(v.MerchantPolicyID),
			MerchantId:   int32(v.MerchantID),
			PolicyType:   v.PolicyType,
			Title:        v.Title,
			Description:  v.Description,
			CreatedAt:    v.CreatedAt.Time.String(),
			UpdatedAt:    v.UpdatedAt.Time.String(),
			MerchantName: v.MerchantName,
		}
	case *db.GetMerchantPoliciesActiveRow:
		return &pbmerchant_policy.MerchantPoliciesResponseDeleteAt{
			Id:           int32(v.MerchantPolicyID),
			MerchantId:   int32(v.MerchantID),
			PolicyType:   v.PolicyType,
			Title:        v.Title,
			Description:  v.Description,
			CreatedAt:    v.CreatedAt.Time.String(),
			UpdatedAt:    v.UpdatedAt.Time.String(),
			MerchantName: v.MerchantName,
			DeletedAt:    &wrapperspb.StringValue{Value: v.DeletedAt.Time.String()},
		}
	case *db.GetMerchantPoliciesTrashedRow:
		return &pbmerchant_policy.MerchantPoliciesResponseDeleteAt{
			Id:           int32(v.MerchantPolicyID),
			MerchantId:   int32(v.MerchantID),
			PolicyType:   v.PolicyType,
			Title:        v.Title,
			Description:  v.Description,
			CreatedAt:    v.CreatedAt.Time.String(),
			UpdatedAt:    v.UpdatedAt.Time.String(),
			DeletedAt:    &wrapperspb.StringValue{Value: v.DeletedAt.Time.String()},
			MerchantName: v.MerchantName,
		}
	case *db.CreateMerchantPolicyRow:
		return &pbmerchant_policy.MerchantPoliciesResponse{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Time.String(),
			UpdatedAt:   v.UpdatedAt.Time.String(),
		}
	case *db.UpdateMerchantPolicyRow:
		return &pbmerchant_policy.MerchantPoliciesResponse{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Time.String(),
			UpdatedAt:   v.UpdatedAt.Time.String(),
		}
	case *db.MerchantPolicy:
		return &pbmerchant_policy.MerchantPoliciesResponseDeleteAt{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Time.String(),
			UpdatedAt:   v.UpdatedAt.Time.String(),
			DeletedAt:   &wrapperspb.StringValue{Value: v.DeletedAt.Time.String()},
		}
	default:
		return nil
	}
}
