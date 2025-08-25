package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
)

type authHandleGrpc struct {
	pb.UnimplementedAuthServiceServer
	logger               logger.LoggerInterface
	registerService      service.RegistrationService
	loginService         service.LoginService
	passwordResetService service.PasswordResetService
	identifyService      service.IdentifyService
	mapping              protomapper.AuthProtoMapper
}

func NewAuthHandleGrpc(authService *service.Service, logger logger.LoggerInterface) pb.AuthServiceServer {
	return &authHandleGrpc{
		registerService:      authService.Register,
		loginService:         authService.Login,
		passwordResetService: authService.PasswordReset,
		identifyService:      authService.Identify,
		logger:               logger,
		mapping:              protomapper.NewAuthProtoMapper(),
	}
}

func (s *authHandleGrpc) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.ApiResponseVerifyCode, error) {
	s.logger.Info("Verify code", zap.String("code", req.Code))

	_, err := s.passwordResetService.VerifyCode(ctx, req.Code)
	if err != nil {
		s.logger.Error("Verify code failed", zap.String("code", req.Code), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Verify code successful", zap.String("code", req.Code))
	return s.mapping.ToProtoResponseVerifyCode("success", "Verify code successful"), nil
}

func (s *authHandleGrpc) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ApiResponseForgotPassword, error) {
	s.logger.Info("Forgot password request", zap.String("email", req.Email))

	_, err := s.passwordResetService.ForgotPassword(ctx, req.Email)
	if err != nil {
		s.logger.Error("Forgot password failed", zap.String("email", req.Email), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Forgot password successful", zap.String("email", req.Email))
	return s.mapping.ToProtoResponseForgotPassword("success", "Forgot password successful"), nil
}

func (s *authHandleGrpc) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ApiResponseResetPassword, error) {
	s.logger.Info("Reset password request", zap.String("reset_token", req.ResetToken))

	_, err := s.passwordResetService.ResetPassword(ctx, &requests.CreateResetPasswordRequest{
		ResetToken:      req.ResetToken,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		s.logger.Error("Reset password failed", zap.String("reset_token", req.ResetToken), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Reset password successful", zap.String("reset_token", req.ResetToken))
	return s.mapping.ToProtoResponseResetPassword("success", "Reset password successful"), nil
}

func (s *authHandleGrpc) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.ApiResponseLogin, error) {
	s.logger.Info("Login attempt", zap.String("email", req.Email))

	request := &requests.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.loginService.Login(ctx, request)
	if err != nil {
		s.logger.Error("Login failed", zap.String("email", req.Email), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Login successful", zap.String("email", req.Email))
	return s.mapping.ToProtoResponseLogin("success", "Login successful", res), nil
}

func (s *authHandleGrpc) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.ApiResponseRefreshToken, error) {
	s.logger.Info("Refresh token request", zap.String("refresh_token", req.RefreshToken))

	res, err := s.identifyService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		s.logger.Error("Refresh token failed", zap.String("refresh_token", req.RefreshToken), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Refresh token successful", zap.String("refresh_token", req.RefreshToken))
	return s.mapping.ToProtoResponseRefreshToken("success", "Refresh token successful", res), nil
}

func (s *authHandleGrpc) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.ApiResponseGetMe, error) {
	s.logger.Info("Get user profile request", zap.String("access_token", req.AccessToken))

	res, err := s.identifyService.GetMe(ctx, req.AccessToken)
	if err != nil {
		s.logger.Error("GetMe failed", zap.String("access_token", req.AccessToken), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("GetMe successful", zap.Any("user_id", res.ID))
	return s.mapping.ToProtoResponseGetMe("success", "Get user profile successful", res), nil
}

func (s *authHandleGrpc) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.ApiResponseRegister, error) {
	s.logger.Info("User registration request", zap.String("email", req.Email))

	request := &requests.RegisterRequest{
		FirstName:       req.Firstname,
		LastName:        req.Lastname,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	res, err := s.registerService.Register(ctx, request)
	if err != nil {
		s.logger.Error("Registration failed", zap.String("email", req.Email), zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Registration successful", zap.String("email", req.Email), zap.Any("user_id", res.ID))
	return s.mapping.ToProtoResponseRegister("success", "Registration successful", res), nil
}
