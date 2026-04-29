package merchantbusinessapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantBusinessCommandResponseMapper struct{}

func NewMerchantBusinessCommandResponseMapper() MerchantBusinessCommandResponseMapper {
	return &merchantBusinessCommandResponseMapper{}
}

func (m *merchantBusinessCommandResponseMapper) ToResponseMerchantBusiness(merchant *pb.MerchantBusinessResponse) *response.MerchantBusinessResponse {
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

func (m *merchantBusinessCommandResponseMapper) ToResponsesMerchantBusiness(merchants []*pb.MerchantBusinessResponse) []*response.MerchantBusinessResponse {
	var mappedMerchants []*response.MerchantBusinessResponse
	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantBusiness(merchant))
	}
	return mappedMerchants
}

func (m *merchantBusinessCommandResponseMapper) ToApiResponseMerchantBusiness(pbResponse *pb.ApiResponseMerchantBusiness) *response.ApiResponseMerchantBusiness {
	return &response.ApiResponseMerchantBusiness{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantBusiness(pbResponse.Data),
	}
}

func (m *merchantBusinessCommandResponseMapper) ToResponseMerchantBusinessDeleteAt(merchant *pb.MerchantBusinessResponseDeleteAt) *response.MerchantBusinessResponseDeleteAt {
	if merchant == nil { return nil }
	var deletedAt string
	if merchant.DeletedAt != nil {
		deletedAt = merchant.DeletedAt.Value
	}

	return &response.MerchantBusinessResponseDeleteAt{
		ID:                int(merchant.Id),
		MerchantID:        int(merchant.MerchantId),
		BusinessType:      merchant.BusinessType,
		TaxID:             merchant.TaxId,
		EstablishedYear:   int(merchant.EstablishedYear),
		NumberOfEmployees: int(merchant.NumberOfEmployees),
		WebsiteUrl:        merchant.WebsiteUrl,
		MerchantName:      merchant.MerchantName,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
		DeletedAt:         &deletedAt,
	}
}

func (m *merchantBusinessCommandResponseMapper) ToResponsesMerchantBusinessDeleteAt(merchants []*pb.MerchantBusinessResponseDeleteAt) []*response.MerchantBusinessResponseDeleteAt {
	var mappedMerchants []*response.MerchantBusinessResponseDeleteAt
	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, m.ToResponseMerchantBusinessDeleteAt(merchant))
	}
	return mappedMerchants
}

func (m *merchantBusinessCommandResponseMapper) ToApiResponseMerchantBusinessDeleteAt(pbResponse *pb.ApiResponseMerchantBusinessDeleteAt) *response.ApiResponseMerchantBusinessDeleteAt {
	return &response.ApiResponseMerchantBusinessDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantBusinessDeleteAt(pbResponse.Data),
	}
}
