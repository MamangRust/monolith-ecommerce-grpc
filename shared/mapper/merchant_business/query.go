package merchantbusinessapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type merchantBusinessQueryResponseMapper struct{}

func NewMerchantBusinessQueryResponseMapper() MerchantBusinessQueryResponseMapper {
	return &merchantBusinessQueryResponseMapper{}
}

func (m *merchantBusinessQueryResponseMapper) ToResponseMerchantBusiness(merchant *pb.MerchantBusinessResponse) *response.MerchantBusinessResponse {
	if merchant == nil { return nil }
	return &response.MerchantBusinessResponse{
		ID:                int(merchant.Id),
		MerchantID:        int(merchant.MerchantId),
		BusinessType:      merchant.BusinessType,
		TaxID:             merchant.TaxId,
		EstablishedYear:   int(merchant.EstablishedYear),
		NumberOfEmployees: int(merchant.NumberOfEmployees),
		WebsiteUrl:        merchant.WebsiteUrl,
		MerchantName:      &merchant.MerchantName,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
	}
}

func (m *merchantBusinessQueryResponseMapper) ToResponsesMerchantBusiness(merchants []*pb.MerchantBusinessResponse) []*response.MerchantBusinessResponse {
	var mappedMerchants []*response.MerchantBusinessResponse
	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantBusiness(merchant))
	}
	return mappedMerchants
}

func (m *merchantBusinessQueryResponseMapper) ToApiResponseMerchantBusiness(pbResponse *pb.ApiResponseMerchantBusiness) *response.ApiResponseMerchantBusiness {
	return &response.ApiResponseMerchantBusiness{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantBusiness(pbResponse.Data),
	}
}

func (m *merchantBusinessQueryResponseMapper) ToApiResponsesMerchantBusiness(pbResponse *pb.ApiResponsesMerchantBusiness) *response.ApiResponsesMerchantBusiness {
	return &response.ApiResponsesMerchantBusiness{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantBusiness(pbResponse.Data),
	}
}

func (m *merchantBusinessQueryResponseMapper) ToApiResponsePaginationMerchantBusiness(pbResponse *pb.ApiResponsePaginationMerchantBusiness) *response.ApiResponsePaginationMerchantBusiness {
	return &response.ApiResponsePaginationMerchantBusiness{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantBusiness(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantBusinessQueryResponseMapper) ToApiResponsePaginationMerchantBusinessDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantBusinessDeleteAt) *response.ApiResponsePaginationMerchantBusinessDeleteAt {
	var data []*response.MerchantBusinessResponseDeleteAt
	for _, b := range pbResponse.Data {
		var deletedAt string
		if b.DeletedAt != nil { deletedAt = b.DeletedAt.Value }
		data = append(data, &response.MerchantBusinessResponseDeleteAt{
			ID:                int(b.Id),
			MerchantID:        int(b.MerchantId),
			BusinessType:      b.BusinessType,
			TaxID:             b.TaxId,
			EstablishedYear:   int(b.EstablishedYear),
			NumberOfEmployees: int(b.NumberOfEmployees),
			WebsiteUrl:        b.WebsiteUrl,
			MerchantName:      b.MerchantName,
			CreatedAt:         b.CreatedAt,
			UpdatedAt:         b.UpdatedAt,
			DeletedAt:         &deletedAt,
		})
	}
	return &response.ApiResponsePaginationMerchantBusinessDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       data,
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
