package merchantdetailapimapper

import (
	pbmerchant_detail "github.com/MamangRust/monolith-ecommerce-pb/merchant_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type merchantDetailQueryResponseMapper struct {
	MerchantDetailCommandResponseMapper
}

func NewMerchantDetailQueryResponseMapper() MerchantDetailQueryResponseMapper {
	return &merchantDetailQueryResponseMapper{
		MerchantDetailCommandResponseMapper: NewMerchantDetailCommandResponseMapper(),
	}
}

func (m *merchantDetailQueryResponseMapper) ToResponseMerchantDetail(merchant *pbmerchant_detail.MerchantDetailResponse) *response.MerchantDetailResponse {
	return m.MerchantDetailCommandResponseMapper.ToResponseMerchantDetail(merchant)
}

func (m *merchantDetailQueryResponseMapper) ToResponseMerchantDetailRelation(merchant *pbmerchant_detail.MerchantDetailResponse) *response.MerchantDetailResponse {
	return m.MerchantDetailCommandResponseMapper.ToResponseMerchantDetailRelation(merchant)
}

func (m *merchantDetailQueryResponseMapper) ToResponsesMerchantDetail(merchants []*pbmerchant_detail.MerchantDetailResponse) []*response.MerchantDetailResponse {
	return m.MerchantDetailCommandResponseMapper.ToResponsesMerchantDetail(merchants)
}

func (m *merchantDetailQueryResponseMapper) ToApiResponseMerchantDetail(pbResponse *pbmerchant_detail.ApiResponseMerchantDetail) *response.ApiResponseMerchantDetail {
	return m.MerchantDetailCommandResponseMapper.ToApiResponseMerchantDetail(pbResponse)
}

func (m *merchantDetailQueryResponseMapper) ToApiResponseMerchantDetailRelation(pbResponse *pbmerchant_detail.ApiResponseMerchantDetail) *response.ApiResponseMerchantDetailRelation {
	return &response.ApiResponseMerchantDetailRelation{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantDetailRelation(pbResponse.Data),
	}
}

func (m *merchantDetailQueryResponseMapper) ToApiResponsesMerchantDetail(pbResponse *pbmerchant_detail.ApiResponsesMerchantDetail) *response.ApiResponsesMerchantDetail {
	return &response.ApiResponsesMerchantDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantDetail(pbResponse.Data),
	}
}

func (m *merchantDetailQueryResponseMapper) ToApiResponsePaginationMerchantDetail(pbResponse *pbmerchant_detail.ApiResponsePaginationMerchantDetail) *response.ApiResponsePaginationMerchantDetail {
	return &response.ApiResponsePaginationMerchantDetail{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantDetail(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantDetailQueryResponseMapper) ToApiResponsePaginationMerchantDetailDeleteAt(pbResponse *pbmerchant_detail.ApiResponsePaginationMerchantDetailDeleteAt) *response.ApiResponsePaginationMerchantDetailDeleteAt {
	return m.MerchantDetailCommandResponseMapper.ToApiResponsePaginationMerchantDetailDeleteAt(pbResponse)
}
