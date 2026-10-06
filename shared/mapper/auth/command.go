package authapimapper

import (
	pbauth "github.com/MamangRust/monolith-ecommerce-pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type authCommandResponseMapper struct{}

func NewAuthCommandResponseMapper() AuthCommandResponseMapper {
	return &authCommandResponseMapper{}
}

func (s *authCommandResponseMapper) ToResponseVerifyCode(res *pbauth.ApiResponseVerifyCode) *response.ApiResponseVerifyCode {
	return &response.ApiResponseVerifyCode{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authCommandResponseMapper) ToResponseForgotPassword(res *pbauth.ApiResponseForgotPassword) *response.ApiResponseForgotPassword {
	return &response.ApiResponseForgotPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authCommandResponseMapper) ToResponseResetPassword(res *pbauth.ApiResponseResetPassword) *response.ApiResponseResetPassword {
	return &response.ApiResponseResetPassword{
		Status:  res.Status,
		Message: res.Message,
	}
}

func (s *authCommandResponseMapper) ToResponseLogin(res *pbauth.ApiResponseLogin) *response.ApiResponseLogin {
	if res == nil {
		return &response.ApiResponseLogin{
			Status:  "error",
			Message: "response is nil",
			Data:    nil,
		}
	}

	var tokenResponse *response.TokenResponse
	if res.Data != nil {
		tokenResponse = &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		}
	}

	return &response.ApiResponseLogin{
		Status:  res.Status,
		Message: res.Message,
		Data:    tokenResponse,
	}
}

func (s *authCommandResponseMapper) ToResponseRegister(res *pbauth.ApiResponseRegister) *response.ApiResponseRegister {
	if res == nil {
		return &response.ApiResponseRegister{
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

	return &response.ApiResponseRegister{
		Status:  res.Status,
		Message: res.Message,
		Data:    userResponse,
	}
}

func (s *authCommandResponseMapper) ToResponseRefreshToken(res *pbauth.ApiResponseRefreshToken) *response.ApiResponseRefreshToken {
	if res == nil {
		return &response.ApiResponseRefreshToken{
			Status:  "error",
			Message: "response is nil",
			Data:    nil,
		}
	}

	var tokenResponse *response.TokenResponse
	if res.Data != nil {
		tokenResponse = &response.TokenResponse{
			AccessToken:  res.Data.AccessToken,
			RefreshToken: res.Data.RefreshToken,
		}
	}

	return &response.ApiResponseRefreshToken{
		Status:  res.Status,
		Message: res.Message,
		Data:    tokenResponse,
	}
}
