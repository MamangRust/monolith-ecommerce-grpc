package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-auth/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	emails "github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/randomstring"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type passwordResetService struct {
	errorhandler      errorhandler.PasswordResetErrorHandler
	errorRandomString errorhandler.RandomStringErrorHandler
	errorMarshal      errorhandler.MarshalErrorHandler
	errorPassword     errorhandler.PasswordErrorHandler
	errorKafka        errorhandler.KafkaErrorHandler
	mencache          mencache.PasswordResetCache
	trace             trace.Tracer
	kafka             *kafka.Kafka
	logger            logger.LoggerInterface
	user              repository.UserRepository
	resetToken        repository.ResetTokenRepository
	requestCounter    *prometheus.CounterVec
	requestDuration   *prometheus.HistogramVec
}

func NewPasswordResetService(
	errorhandler errorhandler.PasswordResetErrorHandler,
	errorRandomString errorhandler.RandomStringErrorHandler,
	errorMarshal errorhandler.MarshalErrorHandler,
	errorPassword errorhandler.PasswordErrorHandler,
	errorKafka errorhandler.KafkaErrorHandler,
	mencache mencache.PasswordResetCache,
	kafka *kafka.Kafka, logger logger.LoggerInterface, user repository.UserRepository, resetToken repository.ResetTokenRepository) *passwordResetService {
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

	return &passwordResetService{
		errorhandler:      errorhandler,
		errorRandomString: errorRandomString,
		errorPassword:     errorPassword,
		errorMarshal:      errorMarshal,
		errorKafka:        errorKafka,
		mencache:          mencache,
		kafka:             kafka,
		trace:             otel.Tracer("password-reset-service"),
		user:              user,
		logger:            logger,
		resetToken:        resetToken,
		requestCounter:    requestCounter,
		requestDuration:   requestDuration,
	}
}

func (s *passwordResetService) ForgotPassword(ctx context.Context, email string) (bool, *response.ErrorResponse) {
	const method = "ForgotPassword"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("email", email))

	defer func() {
		end(status)
	}()

	res, err := s.user.FindByEmail(ctx, email)
	if err != nil {
		return s.errorhandler.HandleFindEmailError(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email))
	}

	span.SetAttributes(attribute.Int("user.id", res.ID))

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		return s.errorRandomString.HandleRandomStringErrorForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email), zap.Error(err))
	}

	_, err = s.resetToken.CreateResetToken(ctx, &requests.CreateResetTokenRequest{
		UserID:     res.ID,
		ResetToken: random,
		ExpiredAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return s.errorhandler.HandleCreateResetTokenError(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email), zap.Error(err))
	}

	s.mencache.SetResetTokenCache(ctx, random, res.ID, 5*time.Minute)

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
		return s.errorMarshal.HandleMarsalForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-forgot-password", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return s.errorKafka.HandleSendEmailForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	logSuccess("Successfully sent password reset email", zap.String("email", email))

	return true, nil
}

func (s *passwordResetService) ResetPassword(ctx context.Context, req *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse) {
	const method = "ResetPassword"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("reset_token", req.ResetToken))

	defer func() {
		end(status)
	}()

	var userID int
	var found bool

	userID, found = s.mencache.GetResetTokenCache(ctx, req.ResetToken)
	if !found {
		res, err := s.resetToken.FindByToken(ctx, req.ResetToken)
		if err != nil {
			return s.errorhandler.HandleFindTokenError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
		}
		userID = int(res.UserID)

		s.mencache.SetResetTokenCache(ctx, req.ResetToken, userID, 5*time.Minute)
	}

	if req.Password != req.ConfirmPassword {
		err := errors.New("password and confirm password do not match")
		return s.errorPassword.HandlePasswordNotMatchError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
	}

	_, err := s.user.UpdateUserPassword(ctx, userID, req.Password)
	if err != nil {
		return s.errorhandler.HandleUpdatePasswordError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", req.ResetToken))
	}

	_ = s.resetToken.DeleteResetToken(ctx, userID)
	s.mencache.DeleteResetTokenCache(ctx, req.ResetToken)

	logSuccess("Successfully reset password", zap.String("reset_token", req.ResetToken))

	return true, nil
}

func (s *passwordResetService) VerifyCode(ctx context.Context, code string) (bool, *response.ErrorResponse) {
	const method = "VerifyCode"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.String("code", code))

	defer func() {
		end(status)
	}()

	res, err := s.user.FindByVerificationCode(ctx, code)
	if err != nil {
		return s.errorhandler.HandleVerifyCodeError(err, method, "VERIFY_CODE_ERR", span, &status, zap.String("code", code))
	}

	_, err = s.user.UpdateUserIsVerified(ctx, res.ID, true)
	if err != nil {
		return s.errorhandler.HandleUpdateVerifiedError(err, method, "VERIFY_CODE_ERR", span, &status, zap.Int("user.id", res.ID))
	}

	s.mencache.DeleteVerificationCodeCache(ctx, res.Email)

	htmlBody := emails.GenerateEmailHTML(map[string]string{
		"Title":   "Verification Success",
		"Message": "Your account has been successfully verified. Click the button below to view or manage your card.",
		"Button":  "Go to Dashboard",
		"Link":    "https://sanedge.example.com/card/create",
	})

	payloadBytes, err := json.Marshal(htmlBody)
	if err != nil {
		return s.errorMarshal.HandleMarshalVerifyCode(err, method, "SEND_EMAIL_VERIFY_CODE_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-verify-code-success", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		return s.errorKafka.HandleSendEmailVerifyCode(err, method, "SEND_EMAIL_VERIFY_CODE_ERR", span, &status, zap.Error(err))
	}

	logSuccess("Successfully verify code", zap.String("code", code))

	return true, nil
}

func (s *passwordResetService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Info("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *passwordResetService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
