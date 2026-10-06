package merchantapimapper

import (
	pbmerchant "github.com/MamangRust/monolith-ecommerce-pb/merchant"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type MerchantBaseResponseMapper interface {
	ToResponseMerchant(merchant *pbmerchant.MerchantResponse) *response.MerchantResponse
	ToResponsesMerchant(merchants []*pbmerchant.MerchantResponse) []*response.MerchantResponse
}

type MerchantQueryResponseMapper interface {
	MerchantBaseResponseMapper
	ToApiResponseMerchant(pbResponse *pbmerchant.ApiResponseMerchant) *response.ApiResponseMerchant
	ToApiResponsesMerchant(pbResponse *pbmerchant.ApiResponsesMerchant) *response.ApiResponsesMerchant
	ToApiResponsePaginationMerchant(pbResponse *pbmerchant.ApiResponsePaginationMerchant) *response.ApiResponsePaginationMerchant
	ToApiResponsePaginationMerchantDeleteAt(pbResponse *pbmerchant.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
}

type MerchantCommandResponseMapper interface {
	MerchantBaseResponseMapper
	ToApiResponseMerchant(pbResponse *pbmerchant.ApiResponseMerchant) *response.ApiResponseMerchant
	ToResponseMerchantDeleteAt(merchant *pbmerchant.MerchantResponseDeleteAt) *response.MerchantResponseDeleteAt
	ToResponsesMerchantDeleteAt(merchants []*pbmerchant.MerchantResponseDeleteAt) []*response.MerchantResponseDeleteAt
	ToApiResponseMerchantDeleteAt(pbResponse *pbmerchant.ApiResponseMerchantDeleteAt) *response.ApiResponseMerchantDeleteAt
	ToApiResponseMerchantDelete(pbResponse *pbmerchant.ApiResponseMerchantDelete) *response.ApiResponseMerchantDelete
	ToApiResponseMerchantAll(pbResponse *pbmerchant.ApiResponseMerchantAll) *response.ApiResponseMerchantAll
	ToApiResponsePaginationMerchantDeleteAt(pbResponse *pbmerchant.ApiResponsePaginationMerchantDeleteAt) *response.ApiResponsePaginationMerchantDeleteAt
}
