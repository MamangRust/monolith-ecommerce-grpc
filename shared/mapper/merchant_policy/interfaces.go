package merchantpolicyapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantPolicyBaseResponseMapper interface {
	ToResponseMerchantPolicy(merchant *pb.MerchantPoliciesResponse) *response.MerchantPoliciesResponse
	ToResponsesMerchantPolicy(merchants []*pb.MerchantPoliciesResponse) []*response.MerchantPoliciesResponse
	ToApiResponseMerchantPolicies(pbResponse *pb.ApiResponseMerchantPolicies) *response.ApiResponseMerchantPolicies
}

type MerchantPolicyQueryResponseMapper interface {
	MerchantPolicyBaseResponseMapper
	ToApiResponsesMerchantPolicies(pbResponse *pb.ApiResponsesMerchantPolicies) *response.ApiResponsesMerchantPolicies
	ToApiResponsePaginationMerchantPolicies(pbResponse *pb.ApiResponsePaginationMerchantPolicies) *response.ApiResponsePaginationMerchantPolicies
	ToApiResponsePaginationMerchantPoliciesDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantPoliciesDeleteAt) *response.ApiResponsePaginationMerchantPoliciesDeleteAt
}

type MerchantPolicyCommandResponseMapper interface {
	MerchantPolicyBaseResponseMapper
	ToResponseMerchantPolicyDeleteAt(merchant *pb.MerchantPoliciesResponseDeleteAt) *response.MerchantPoliciesResponseDeleteAt
	ToResponsesMerchantPolicyDeleteAt(merchants []*pb.MerchantPoliciesResponseDeleteAt) []*response.MerchantPoliciesResponseDeleteAt
	ToApiResponseMerchantPoliciesDeleteAt(pbResponse *pb.ApiResponseMerchantPoliciesDeleteAt) *response.ApiResponseMerchantPoliciesDeleteAt
}
