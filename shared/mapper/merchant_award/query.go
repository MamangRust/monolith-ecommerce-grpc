package merchantawardapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type merchantAwardQueryResponseMapper struct{}

func NewMerchantAwardQueryResponseMapper() MerchantAwardQueryResponseMapper {
	return &merchantAwardQueryResponseMapper{}
}

func (m *merchantAwardQueryResponseMapper) ToResponseMerchantAward(MerchantAward *pb.MerchantAwardResponse) *response.MerchantAwardResponse {
	if MerchantAward == nil { return nil }
	return &response.MerchantAwardResponse{
		ID:             int(MerchantAward.Id),
		MerchantID:     int(MerchantAward.MerchantId),
		Title:          MerchantAward.Title,
		Description:    MerchantAward.Description,
		IssuedBy:       MerchantAward.IssuedBy,
		IssueDate:      MerchantAward.IssueDate,
		ExpiryDate:     MerchantAward.ExpiryDate,
		CertificateUrl: MerchantAward.CertificateUrl,
		CreatedAt:      MerchantAward.CreatedAt,
		UpdatedAt:      MerchantAward.UpdatedAt,
	}
}

func (m *merchantAwardQueryResponseMapper) ToResponsesMerchantAward(MerchantAwards []*pb.MerchantAwardResponse) []*response.MerchantAwardResponse {
	var mappedMerchantAwards []*response.MerchantAwardResponse
	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.ToResponseMerchantAward(MerchantAward))
	}
	return mappedMerchantAwards
}

func (m *merchantAwardQueryResponseMapper) ToApiResponseMerchantAward(pbResponse *pb.ApiResponseMerchantAward) *response.ApiResponseMerchantAward {
	return &response.ApiResponseMerchantAward{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantAward(pbResponse.Data),
	}
}

func (m *merchantAwardQueryResponseMapper) ToApiResponsesMerchantAward(pbResponse *pb.ApiResponsesMerchantAward) *response.ApiResponsesMerchantAward {
	return &response.ApiResponsesMerchantAward{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesMerchantAward(pbResponse.Data),
	}
}

func (m *merchantAwardQueryResponseMapper) ToApiResponsePaginationMerchantAward(pbResponse *pb.ApiResponsePaginationMerchantAward) *response.ApiResponsePaginationMerchantAward {
	return &response.ApiResponsePaginationMerchantAward{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesMerchantAward(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *merchantAwardQueryResponseMapper) ToApiResponsePaginationMerchantAwardDeleteAt(pbResponse *pb.ApiResponsePaginationMerchantAwardDeleteAt) *response.ApiResponsePaginationMerchantAwardDeleteAt {
	var data []*response.MerchantAwardResponseDeleteAt
	for _, b := range pbResponse.Data {
		var deletedAt string
		if b.DeletedAt != nil { deletedAt = b.DeletedAt.Value }
		data = append(data, &response.MerchantAwardResponseDeleteAt{
			ID:             int(b.Id),
			MerchantID:     int(b.MerchantId),
			Title:          b.Title,
			Description:    b.Description,
			IssuedBy:       b.IssuedBy,
			IssueDate:      b.IssueDate,
			ExpiryDate:     b.ExpiryDate,
			CertificateUrl: b.CertificateUrl,
			CreatedAt:      b.CreatedAt,
			UpdatedAt:      b.UpdatedAt,
			DeletedAt:      &deletedAt,
		})
	}
	return &response.ApiResponsePaginationMerchantAwardDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       data,
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
