package merchantawardapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantAwardBaseResponseMapper interface {
	ToResponseMerchantAward(MerchantAward *pb.MerchantAwardResponse) *response.MerchantAwardResponse
	ToResponsesMerchantAward(MerchantAwards []*pb.MerchantAwardResponse) []*response.MerchantAwardResponse
	ToApiResponseMerchantAward(pbResponse *pb.ApiResponseMerchantAward) *response.ApiResponseMerchantAward
}

type MerchantAwardQueryResponseMapper interface {
	MerchantAwardBaseResponseMapper
	ToApiResponsesMerchantAward(pbResponse *pb.ApiResponsesMerchantAward) *response.ApiResponsesMerchantAward
	ToApiResponsePaginationMerchantAward(pbResponse *pb.ApiResponsePaginationMerchantAward) *response.ApiResponsePaginationMerchantAward
	ToApiResponsePaginationMerchantAwardDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantAwardDeleteAt) *response.ApiResponsePaginationMerchantAwardDeleteAt
}

type MerchantAwardCommandResponseMapper interface {
	MerchantAwardBaseResponseMapper
	ToResponseMerchantAwardDeleteAt(MerchantAward *pb.MerchantAwardResponseDeleteAt) *response.MerchantAwardResponseDeleteAt
	ToResponsesMerchantAwardDeleteAt(MerchantAwards []*pb.MerchantAwardResponseDeleteAt) []*response.MerchantAwardResponseDeleteAt
	ToApiResponseMerchantAwardDeleteAt(pbResponse *pb.ApiResponseMerchantAwardDeleteAt) *response.ApiResponseMerchantAwardDeleteAt
}
