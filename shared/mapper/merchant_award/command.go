package merchantawardapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type merchantAwardCommandResponseMapper struct{}

func NewMerchantAwardCommandResponseMapper() MerchantAwardCommandResponseMapper {
	return &merchantAwardCommandResponseMapper{}
}

func (m *merchantAwardCommandResponseMapper) ToResponseMerchantAward(MerchantAward *pb.MerchantAwardResponse) *response.MerchantAwardResponse {
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

func (m *merchantAwardCommandResponseMapper) ToResponsesMerchantAward(MerchantAwards []*pb.MerchantAwardResponse) []*response.MerchantAwardResponse {
	var mappedMerchantAwards []*response.MerchantAwardResponse
	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.ToResponseMerchantAward(MerchantAward))
	}
	return mappedMerchantAwards
}

func (m *merchantAwardCommandResponseMapper) ToApiResponseMerchantAward(pbResponse *pb.ApiResponseMerchantAward) *response.ApiResponseMerchantAward {
	return &response.ApiResponseMerchantAward{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantAward(pbResponse.Data),
	}
}

func (m *merchantAwardCommandResponseMapper) ToResponseMerchantAwardDeleteAt(MerchantAward *pb.MerchantAwardResponseDeleteAt) *response.MerchantAwardResponseDeleteAt {
	if MerchantAward == nil { return nil }
	var deletedAt string
	if MerchantAward.DeletedAt != nil {
		deletedAt = MerchantAward.DeletedAt.Value
	}

	return &response.MerchantAwardResponseDeleteAt{
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
		DeletedAt:      &deletedAt,
	}
}

func (m *merchantAwardCommandResponseMapper) ToResponsesMerchantAwardDeleteAt(MerchantAwards []*pb.MerchantAwardResponseDeleteAt) []*response.MerchantAwardResponseDeleteAt {
	var mappedMerchantAwards []*response.MerchantAwardResponseDeleteAt
	for _, MerchantAward := range MerchantAwards {
		mappedMerchantAwards = append(mappedMerchantAwards, m.ToResponseMerchantAwardDeleteAt(MerchantAward))
	}
	return mappedMerchantAwards
}

func (m *merchantAwardCommandResponseMapper) ToApiResponseMerchantAwardDeleteAt(pbResponse *pb.ApiResponseMerchantAwardDeleteAt) *response.ApiResponseMerchantAwardDeleteAt {
	return &response.ApiResponseMerchantAwardDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseMerchantAwardDeleteAt(pbResponse.Data),
	}
}
