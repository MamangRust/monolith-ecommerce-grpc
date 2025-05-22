package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	emails "github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/randomstring"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type passwordRegisterService struct {
	ctx             context.Context
	trace           trace.Tracer
	kafka           kafka.Kafka
	logger          logger.LoggerInterface
	user            repository.UserRepository
	resetToken      repository.ResetTokenRepository
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewPasswordResetService(ctx context.Context, kafka kafka.Kafka, logger logger.LoggerInterface, user repository.UserRepository, resetToken repository.ResetTokenRepository) *passwordRegisterService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "password_reset_service_requests_total",
			Help: "Total number of requests to the PasswordResetService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "password_reset_service_request_duration_seconds",
			Help:    "Histogram of request durations for the PasswordResetService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &passwordRegisterService{
		ctx:             ctx,
		kafka:           kafka,
		trace:           otel.Tracer("password-reset-service"),
		user:            user,
		logger:          logger,
		resetToken:      resetToken,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *passwordRegisterService) ForgotPassword(email string) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("ForgotPassword", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "PasswordService.ForgotPassword")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.email", email),
	)

	s.logger.Debug("Starting forgot password process",
		zap.String("email", email),
	)

	res, err := s.user.FindByEmail(email)
	if err != nil {
		traceID := traceunic.GenerateTraceID("USER_NOT_FOUND")
		s.logger.Error("Failed to find user by email",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")
		status = "user_not_found"
		return false, user_errors.ErrUserNotFoundRes
	}

	span.SetAttributes(
		attribute.Int("user.id", res.ID),
	)

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		traceID := traceunic.GenerateTraceID("RANDOM_STR_FAILED")
		s.logger.Error("Failed to generate random string",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to generate reset token")
		status = "token_generation_failed"
		return false, user_errors.ErrInternalServerError
	}

	_, err = s.resetToken.CreateResetToken(&requests.CreateResetTokenRequest{
		UserID:     res.ID,
		ResetToken: random,
		ExpiredAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("CREATE_TOKEN_FAILED")
		s.logger.Error("Failed to create reset token",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to store reset token")
		status = "token_storage_failed"
		return false, user_errors.ErrInternalServerError
	}

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
		traceID := traceunic.GenerateTraceID("MARSHAL_EMAIL_FAILED")
		s.logger.Error("Failed to marshal email payload",
			zap.String("trace_id", traceID),
			zap.Int("user_id", res.ID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to prepare email")
		status = "email_preparation_failed"
		return false, user_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-auth-forgot-password", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("KAFKA_SEND_FAILED")
		s.logger.Error("Failed to send email via Kafka",
			zap.String("trace_id", traceID),
			zap.Int("user_id", res.ID),
			zap.String("email", res.Email),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send email")
		status = "email_send_failed"
		return false, user_errors.ErrFailedSendEmail
	}

	s.logger.Debug("Password reset email sent successfully",
		zap.Int("user_id", res.ID),
		zap.String("email", res.Email),
	)
	span.SetStatus(codes.Ok, "Password reset initiated")
	return true, nil
}

func (s *passwordRegisterService) ResetPassword(req *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("ResetPassword", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "PasswordService.ResetPassword")
	defer span.End()

	span.SetAttributes(
		attribute.String("reset_token", req.ResetToken),
	)

	s.logger.Debug("Starting password reset process",
		zap.String("reset_token", req.ResetToken),
	)

	res, err := s.resetToken.FindByToken(req.ResetToken)
	if err != nil {
		traceID := traceunic.GenerateTraceID("INVALID_RESET_TOKEN")
		s.logger.Error("Failed to find user by reset token",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid reset token")
		status = "invalid_token"
		return false, user_errors.ErrUserNotFoundRes
	}

	if req.Password != req.ConfirmPassword {
		traceID := traceunic.GenerateTraceID("PASSWORD_MISMATCH")
		s.logger.Error("Password and confirmation do not match",
			zap.String("trace_id", traceID),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Password mismatch")
		status = "password_mismatch"
		return false, user_errors.ErrFailedPasswordNoMatch
	}

	_, err = s.user.UpdateUserPassword(int(res.UserID), req.Password)
	if err != nil {
		traceID := traceunic.GenerateTraceID("PASSWORD_UPDATE_FAILED")
		s.logger.Error("Failed to update user password",
			zap.String("trace_id", traceID),
			zap.Int("user_id", int(res.UserID)),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update password")
		status = "password_update_failed"
		return false, user_errors.ErrInternalServerError
	}

	if err := s.resetToken.DeleteResetToken(int(res.UserID)); err != nil {
		traceID := traceunic.GenerateTraceID("DELETE_RESET_TOKEN_FAILED")
		s.logger.Error("Failed to delete reset token",
			zap.String("trace_id", traceID),
			zap.Int("user_id", int(res.UserID)),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete reset token")
		status = "delete_reset_token_failed"
		return false, user_errors.ErrInternalServerError

	}

	s.logger.Debug("Password reset successfully",
		zap.Int("user_id", int(res.UserID)),
	)
	span.SetStatus(codes.Ok, "Password reset successful")
	return true, nil
}

func (s *passwordRegisterService) VerifyCode(code string) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("VerifyCode", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "PasswordService.VerifyCode")
	defer span.End()

	span.SetAttributes(
		attribute.String("verification_code", code),
	)

	s.logger.Debug("Starting verification code process",
		zap.String("code", code),
	)

	res, err := s.user.FindByVerificationCode(code)
	if err != nil {
		traceID := traceunic.GenerateTraceID("INVALID_VERIFICATION_CODE")
		s.logger.Error("Failed to find user by verification code",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid verification code")
		status = "invalid_code"
		return false, user_errors.ErrUserNotFoundRes
	}

	_, err = s.user.UpdateUserIsVerified(res.ID, true)
	if err != nil {
		traceID := traceunic.GenerateTraceID("VERIFICATION_UPDATE_FAILED")
		s.logger.Error("Failed to update user verification status",
			zap.String("trace_id", traceID),
			zap.Int("user_id", res.ID),
			zap.Error(err),
		)

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to verify user")
		status = "verification_failed"
		return false, user_errors.ErrInternalServerError
	}

	s.logger.Debug("User verified successfully",
		zap.Int("user_id", res.ID),
	)
	span.SetStatus(codes.Ok, "User verified successfully")
	return true, nil
}

func (s *passwordRegisterService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
