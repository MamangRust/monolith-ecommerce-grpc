package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	registerService      service.RegistrationService
	loginService         service.LoginService
	passwordResetService service.PasswordResetService
	identifyService      service.IdentifyService
	mapping              protomapper.AuthProtoMapper
}

func NewAuthHandleGrpc(authService service.Service) pb.AuthServiceServer {
	return &authHandleGrpc{
		registerService:      authService.Register,
		loginService:         authService.Login,
		passwordResetService: authService.PasswordReset,
		identifyService:      authService.Identify,
		mapping:              protomapper.NewAuthProtoMapper(),
	}
}

func (s *authHandleGrpc) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.ApiResponseVerifyCode, error) {
	_, err := s.passwordResetService.VerifyCode(ctx, req.Code)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseVerifyCode("success", "Verify code successful"), nil

}

func (s *authHandleGrpc) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ApiResponseForgotPassword, error) {
	_, err := s.passwordResetService.ForgotPassword(ctx, req.Email)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseForgotPassword("success", "Forgot password successful"), nil
}

func (s *authHandleGrpc) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseResetPassword, error) {
	_, err := s.passwordResetService.ResetPassword(ctx, &requests.CreateResetPasswordRequest{
		ResetToken:      req.ResetToken,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseResetPassword("success", "Reset password successful"), nil
}

func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	request := &requests.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.loginService.Login(ctx, request)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseLogin("success", "Login successful", res), nil
}

func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	res, err := s.identifyService.RefreshToken(ctx, req.RefreshToken)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseRefreshToken("success", "Refresh token successful", res), nil
}

func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	res, err := s.identifyService.GetMe(ctx, req.AccessToken)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseGetMe("success", "Refresh token successful", res), nil
}

func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error) {
	request := &requests.RegisterRequest{
		FirstName:       req.Firstname,
		LastName:        req.Lastname,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := s.registerService.Register(ctx, request)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseRegister("success", "Registration successful", res), nil
}
