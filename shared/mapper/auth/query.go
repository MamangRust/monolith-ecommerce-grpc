package authapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type authQueryResponseMapper struct{}

func NewAuthQueryResponseMapper() AuthQueryResponseMapper {
	return &authQueryResponseMapper{}
}

func (s *authQueryResponseMapper) ToResponseGetMe(res *pb.ApiResponseGetMe) *response.ApiResponseGetMe {
	if res == nil {
		return &response.ApiResponseGetMe{
			Status:  "error",
			Message: "response is nil",
			Data:    nil,
		}
	}

	var userResponse *response.UserResponse
	if res.Data != nil {
		userResponse = &response.UserResponse{
			ID:        int(res.Data.Id),
			FirstName: res.Data.Firstname,
			LastName:  res.Data.Lastname,
			Email:     res.Data.Email,
			CreatedAt: res.Data.CreatedAt,
			UpdatedAt: res.Data.UpdatedAt,
		}
	}

	return &response.ApiResponseGetMe{
		Status:  res.Status,
		Message: res.Message,
		Data:    userResponse,
	}
}
