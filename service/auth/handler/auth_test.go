package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-auth/service"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	sharederrors "github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type passwordResetServiceStub struct {
	forgotPassword func(context.Context, string) (bool, error)
	resetPassword  func(context.Context, *requests.CreateResetPasswordRequest) (bool, error)
	verifyCode     func(context.Context, string) (bool, error)
}

func (s passwordResetServiceStub) ForgotPassword(ctx context.Context, email string) (bool, error) {
	if s.forgotPassword == nil {
		return false, fmt.Errorf("unexpected ForgotPassword call")
	}
	return s.forgotPassword(ctx, email)
}

func (s passwordResetServiceStub) ResetPassword(ctx context.Context, req *requests.CreateResetPasswordRequest) (bool, error) {
	if s.resetPassword == nil {
		return false, fmt.Errorf("unexpected ResetPassword call")
	}
	return s.resetPassword(ctx, req)
}

func (s passwordResetServiceStub) VerifyCode(ctx context.Context, code string) (bool, error) {
	if s.verifyCode == nil {
		return false, fmt.Errorf("unexpected VerifyCode call")
	}
	return s.verifyCode(ctx, code)
}

type registrationServiceStub struct{}

func (registrationServiceStub) Register(context.Context, *requests.RegisterRequest) (*db.CreateUserRow, error) {
	return nil, nil
}

type loginServiceStub struct{}

func (loginServiceStub) Login(context.Context, *requests.AuthRequest) (*response.TokenResponse, error) {
	return nil, nil
}

type identifyServiceStub struct{}

func (identifyServiceStub) RefreshToken(context.Context, string) (*response.TokenResponse, error) {
	return nil, nil
}

func (identifyServiceStub) GetMe(context.Context, int) (*db.GetUserByIDRow, error) {
	return nil, nil
}

func newAuthHandlerForTest(passwordReset service.PasswordResetService) pb.AuthServiceServer {
	return NewAuthHandleGrpc(&service.Service{
		Register:      registrationServiceStub{},
		Login:         loginServiceStub{},
		PasswordReset: passwordReset,
		Identify:      identifyServiceStub{},
	}, &logger.Logger{Log: zap.NewNop()})
}

func TestAuthHandlerVerifyCode(t *testing.T) {
	var gotCode string
	handler := newAuthHandlerForTest(passwordResetServiceStub{
		verifyCode: func(_ context.Context, code string) (bool, error) {
			gotCode = code
			return true, nil
		},
	})

	got, err := handler.VerifyCode(context.Background(), &pb.VerifyCodeRequest{Code: "verify-123"})
	if err != nil {
		t.Fatalf("VerifyCode() error = %v", err)
	}
	if gotCode != "verify-123" {
		t.Fatalf("VerifyCode() passed code %q, want %q", gotCode, "verify-123")
	}
	if got.GetStatus() != "success" || got.GetMessage() != "Verification successfully" {
		t.Fatalf("VerifyCode() response = %#v, want success response", got)
	}
}

func TestAuthHandlerForgotPassword(t *testing.T) {
	var gotEmail string
	handler := newAuthHandlerForTest(passwordResetServiceStub{
		forgotPassword: func(_ context.Context, email string) (bool, error) {
			gotEmail = email
			return true, nil
		},
	})

	got, err := handler.ForgotPassword(context.Background(), &pb.ForgotPasswordRequest{Email: "user@example.com"})
	if err != nil {
		t.Fatalf("ForgotPassword() error = %v", err)
	}
	if gotEmail != "user@example.com" {
		t.Fatalf("ForgotPassword() passed email %q, want %q", gotEmail, "user@example.com")
	}
	if got.GetStatus() != "success" || got.GetMessage() != "ForgotPassword successful" {
		t.Fatalf("ForgotPassword() response = %#v, want success response", got)
	}
}

func TestAuthHandlerResetPassword(t *testing.T) {
	var gotRequest *requests.CreateResetPasswordRequest
	handler := newAuthHandlerForTest(passwordResetServiceStub{
		resetPassword: func(_ context.Context, req *requests.CreateResetPasswordRequest) (bool, error) {
			gotRequest = req
			return true, nil
		},
	})

	got, err := handler.ResetPassword(context.Background(), &pb.ResetPasswordRequest{
		ResetToken:      "reset-123",
		Password:        "new-password",
		ConfirmPassword: "new-password",
	})
	if err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if gotRequest == nil {
		t.Fatal("ResetPassword() did not pass a request to the service")
	}
	if gotRequest.ResetToken != "reset-123" || gotRequest.Password != "new-password" || gotRequest.ConfirmPassword != "new-password" {
		t.Fatalf("ResetPassword() request = %#v, want all request fields mapped", gotRequest)
	}
	if got.GetStatus() != "success" || got.GetMessage() != "Reset password successful" {
		t.Fatalf("ResetPassword() response = %#v, want success response", got)
	}
}

func TestAuthHandlerPasswordResetErrorsBecomeGRPCErrors(t *testing.T) {
	tests := []struct {
		name     string
		wantCode codes.Code
		call     func(pb.AuthServiceServer) error
	}{
		{
			name:     "verify code",
			wantCode: codes.NotFound,
			call: func(handler pb.AuthServiceServer) error {
				_, err := handler.VerifyCode(context.Background(), &pb.VerifyCodeRequest{Code: "bad-code"})
				return err
			},
		},
		{
			name:     "forgot password",
			wantCode: codes.NotFound,
			call: func(handler pb.AuthServiceServer) error {
				_, err := handler.ForgotPassword(context.Background(), &pb.ForgotPasswordRequest{Email: "missing@example.com"})
				return err
			},
		},
		{
			name:     "reset password",
			wantCode: codes.InvalidArgument,
			call: func(handler pb.AuthServiceServer) error {
				_, err := handler.ResetPassword(context.Background(), &pb.ResetPasswordRequest{ResetToken: "bad-token"})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newAuthHandlerForTest(passwordResetServiceStub{
				verifyCode: func(context.Context, string) (bool, error) {
					return false, sharederrors.ErrNotFound
				},
				forgotPassword: func(context.Context, string) (bool, error) {
					return false, sharederrors.ErrNotFound
				},
				resetPassword: func(context.Context, *requests.CreateResetPasswordRequest) (bool, error) {
					return false, sharederrors.ErrBadRequest
				},
			})

			err := tt.call(handler)
			if err == nil {
				t.Fatal("handler call error = nil, want error")
			}
			if gotCode := status.Code(err); gotCode != tt.wantCode {
				t.Fatalf("handler call status code = %s, want %s: %v", gotCode, tt.wantCode, err)
			}
		})
	}
}
