package authapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type AuthBaseResponseMapper interface {
}

type AuthQueryResponseMapper interface {
	AuthBaseResponseMapper
	ToResponseGetMe(res *pb.ApiResponseGetMe) *response.ApiResponseGetMe
}

type AuthCommandResponseMapper interface {
	AuthBaseResponseMapper
	ToResponseVerifyCode(res *pb.ApiResponseVerifyCode) *response.ApiResponseVerifyCode
	ToResponseForgotPassword(res *pb.ApiResponseForgotPassword) *response.ApiResponseForgotPassword
	ToResponseResetPassword(res *pb.ApiResponseResetPassword) *response.ApiResponseResetPassword
	ToResponseLogin(res *pb.ApiResponseLogin) *response.ApiResponseLogin
	ToResponseRegister(res *pb.ApiResponseRegister) *response.ApiResponseRegister
	ToResponseRefreshToken(res *pb.ApiResponseRefreshToken) *response.ApiResponseRefreshToken
}
