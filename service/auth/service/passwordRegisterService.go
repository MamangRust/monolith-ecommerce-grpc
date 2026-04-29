package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	mencache "github.com/MamangRust/monolith-ecommerce-auth/cache"
	"github.com/MamangRust/monolith-ecommerce-auth/repository"

	emails "github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	randomstring "github.com/MamangRust/monolith-ecommerce-pkg/randomstring"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// PasswordResetServiceDeps defines dependencies required by PasswordResetService.
type PasswordResetServiceDeps struct {
	Cache         mencache.PasswordResetCache
	Kafka         *kafka.Kafka
	Logger        logger.LoggerInterface
	User          repository.UserRepository
	ResetToken    repository.ResetTokenRepository
	Observability observability.TraceLoggerObservability
}

// passwordResetService implements PasswordResetService.
type passwordResetService struct {
	mencache      mencache.PasswordResetCache
	kafka         *kafka.Kafka
	logger        logger.LoggerInterface
	user          repository.UserRepository
	resetToken    repository.ResetTokenRepository
	observability observability.TraceLoggerObservability
}

func NewPasswordResetService(params *PasswordResetServiceDeps) *passwordResetService {

	return &passwordResetService{
		mencache:      params.Cache,
		kafka:         params.Kafka,
		logger:        params.Logger,
		user:          params.User,
		resetToken:    params.ResetToken,
		observability: params.Observability,
	}
}

func (s *passwordResetService) ForgotPassword(ctx context.Context, email string) (bool, error) {
	const method = "ForgotPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", email))

	defer func() {
		end(status)
	}()

	res, err := s.user.FindByEmail(ctx, email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("email", email))
	}

	span.SetAttributes(attribute.Int("user.id", int(res.UserID)))

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("email", email))
	}

	_, err = s.resetToken.CreateResetToken(ctx, &requests.CreateResetTokenRequest{
		UserID:     int(res.UserID),
		ResetToken: random,
		ExpiredAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
	}

	s.mencache.SetResetTokenCache(ctx, random, int(res.UserID), 5*time.Minute)

	htmlBody := emails.GenerateEmailHTML(map[string]string{
		"Title":   "Reset Your Password",
		"Message": "Click the button below to reset your password.",
		"Button":  "Reset Password",
		"Link":    "https://sanedge.example.com/reset-password?token=" + random,
	})

	emailPayload := map[string]any{
		"email":   res.Email,
		"subject": "Password Reset Request",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("email", email))
	}

	if s.kafka != nil {
		err = s.kafka.SendMessage("email-service-topic-auth-forgot-password", strconv.Itoa(int(res.UserID)), payloadBytes)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
		}
	}

	logSuccess("Successfully sent password reset email", zap.String("email", email))

	return true, nil
}

func (s *passwordResetService) ResetPassword(ctx context.Context, req *requests.CreateResetPasswordRequest) (bool, error) {
	const method = "ResetPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("reset_token", req.ResetToken))

	defer func() {
		end(status)
	}()

	var userID int
	var found bool

	userID, found = s.mencache.GetResetTokenCache(ctx, req.ResetToken)
	if !found {
		res, err := s.resetToken.FindByToken(ctx, req.ResetToken)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("reset_token", req.ResetToken))
		}
		userID = int(res.UserID)

		s.mencache.SetResetTokenCache(ctx, req.ResetToken, userID, 5*time.Minute)
	}

	if req.Password != req.ConfirmPassword {
		status = "error"
		err := errors.New("password and confirm password do not match")
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("reset_token", req.ResetToken))
	}

	_, err := s.user.UpdateUserPassword(ctx, userID, req.Password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Int("user.id", userID))
	}

	_ = s.resetToken.DeleteResetToken(ctx, userID)
	s.mencache.DeleteResetTokenCache(ctx, req.ResetToken)

	logSuccess("Successfully reset password", zap.String("reset_token", req.ResetToken))

	return true, nil
}

func (s *passwordResetService) VerifyCode(ctx context.Context, code string) (bool, error) {
	const method = "VerifyCode"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("code", code))

	defer func() {
		end(status)
	}()

	res, err := s.user.FindByVerificationCode(ctx, code)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("code", code))
	}

	_, err = s.user.UpdateUserIsVerified(ctx, int(res.UserID), true)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
	}

	s.mencache.DeleteVerificationCodeCache(ctx, res.Email)

	htmlBody := emails.GenerateEmailHTML(map[string]string{
		"Title":   "Verification Success",
		"Message": "Your account has been successfully verified. Click the button below to view or manage your card.",
		"Button":  "Go to Dashboard",
		"Link":    "https://sanedge.example.com/card/create",
	})

	emailPayload := map[string]any{
		"email":   res.Email,
		"subject": "Verification Success",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.String("code", code))
	}

	if s.kafka != nil {
		err = s.kafka.SendMessage("email-service-topic-auth-verify-code-success", strconv.Itoa(int(res.UserID)), payloadBytes)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[bool](s.logger, err, method, span, zap.Int("user.id", int(res.UserID)))
		}
	}

	logSuccess("Successfully verify code", zap.String("code", code))

	return true, nil
}
