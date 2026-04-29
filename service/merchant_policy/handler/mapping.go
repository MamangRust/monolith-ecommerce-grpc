package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func mapToSingleResponse(data interface{}) *pb.ApiResponseMerchantPolicies {
	return &pb.ApiResponseMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policy",
		Data:    mapToMerchantPolicyResponse(data).(*pb.MerchantPoliciesResponse),
	}
}

func mapToPaginationResponse(data []*db.GetMerchantPoliciesRow, total *int) *pb.ApiResponsePaginationMerchantPolicies {
	var policies []*pb.MerchantPoliciesResponse
	for _, v := range data {
		policies = append(policies, mapToMerchantPolicyResponse(v).(*pb.MerchantPoliciesResponse))
	}

	return &pb.ApiResponsePaginationMerchantPolicies{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pb.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToPaginationDeleteAtResponse(data interface{}, total *int) *pb.ApiResponsePaginationMerchantPoliciesDeleteAt {
	var policies []*pb.MerchantPoliciesResponseDeleteAt
	
	switch v := data.(type) {
	case []*db.GetMerchantPoliciesActiveRow:
		for _, item := range v {
			policies = append(policies, mapToMerchantPolicyResponse(item).(*pb.MerchantPoliciesResponseDeleteAt))
		}
	case []*db.GetMerchantPoliciesTrashedRow:
		for _, item := range v {
			policies = append(policies, mapToMerchantPolicyResponse(item).(*pb.MerchantPoliciesResponseDeleteAt))
		}
	}

	return &pb.ApiResponsePaginationMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully fetched merchant policies",
		Data:    policies,
		Pagination: &pb.PaginationMeta{
			TotalRecords: int32(*total),
		},
	}
}

func mapToSingleDeleteAtResponse(data *db.MerchantPolicy) *pb.ApiResponseMerchantPoliciesDeleteAt {
	return &pb.ApiResponseMerchantPoliciesDeleteAt{
		Status:  "success",
		Message: "Successfully processed merchant policy",
		Data:    mapToMerchantPolicyResponse(data).(*pb.MerchantPoliciesResponseDeleteAt),
	}
}

func mapToMerchantPolicyResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *db.GetMerchantPolicyRow:
		return &pb.MerchantPoliciesResponse{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Time.String(),
			UpdatedAt:   v.UpdatedAt.Time.String(),
		}
	case *db.GetMerchantPoliciesRow:
		return &pb.MerchantPoliciesResponse{
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
		return &pb.MerchantPoliciesResponseDeleteAt{
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
		return &pb.MerchantPoliciesResponseDeleteAt{
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
		return &pb.MerchantPoliciesResponse{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Time.String(),
			UpdatedAt:   v.UpdatedAt.Time.String(),
		}
	case *db.UpdateMerchantPolicyRow:
		return &pb.MerchantPoliciesResponse{
			Id:          int32(v.MerchantPolicyID),
			MerchantId:  int32(v.MerchantID),
			PolicyType:  v.PolicyType,
			Title:       v.Title,
			Description: v.Description,
			CreatedAt:   v.CreatedAt.Time.String(),
			UpdatedAt:   v.UpdatedAt.Time.String(),
		}
	case *db.MerchantPolicy:
		return &pb.MerchantPoliciesResponseDeleteAt{
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
