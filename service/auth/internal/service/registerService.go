package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-auth/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/randomstring"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb/errors/role_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb/errors/user_errors"
	userrole_errors "github.com/MamangRust/monolith-point-of-sale-shared/pb/errors/user_role_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/pb/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type registerService struct {
	ctx             context.Context
	trace           trace.Tracer
	user            repository.UserRepository
	role            repository.RoleRepository
	userRole        repository.UserRoleRepository
	hash            hash.HashPassword
	kafka           kafka.Kafka
	logger          logger.LoggerInterface
	mapping         response_service.UserResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewRegisterService(ctx context.Context, user repository.UserRepository, role repository.RoleRepository, userRole repository.UserRoleRepository, hash hash.HashPassword, kafka kafka.Kafka, logger logger.LoggerInterface, mapping response_service.UserResponseMapper) *registerService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "register_service_requests_total",
			Help: "Total number of requests to the RegisterService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "register_service_request_duration_seconds",
			Help:    "Histogram of request durations for the RegisterService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &registerService{
		ctx:             ctx,
		trace:           otel.Tracer("register-service"),
		user:            user,
		role:            role,
		userRole:        userRole,
		hash:            hash,
		kafka:           kafka,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *registerService) Register(request *requests.RegisterRequest) (*response.UserResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("Register", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RegisterUser")
	defer span.End()

	span.SetAttributes(
		attribute.String("email", request.Email),
		attribute.String("first_name", request.FirstName),
		attribute.String("last_name", request.LastName),
	)

	s.logger.Debug("Starting user registration",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	existingUser, err := s.user.FindByEmail(request.Email)
	if err == nil && existingUser != nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Email already exists", zap.String("trace_id", traceID))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Email already exists")

		status = "email_already_exists"
		return nil, user_errors.ErrUserEmailAlready
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Failed to hash password", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to hash password")

		status = "hash_password_failed"
		return nil, user_errors.ErrUserPassword
	}
	request.Password = passwordHash

	const defaultRoleName = "Cashier"

	role, err := s.role.FindByName(defaultRoleName)

	if err != nil || role == nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Failed to find default role", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find default role")

		status = "role_not_found"
		return nil, role_errors.ErrRoleNotFoundRes
	}

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Failed to generate verification code", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to generate verification code")
		status = "verification_code_failed"
		return nil, user_errors.ErrFailedCreateUser
	}

	request.VerifiedCode = random
	request.IsVerified = false

	newUser, err := s.user.CreateUser(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Failed to create user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create user")

		status = "create_user_failed"
		return nil, user_errors.ErrFailedCreateUser
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge",
		"Message": "Your account has been successfully created.",
		"Button":  "Login Now",
		"Link":    "https://sanedge.example.com/login?verify_code=" + request.VerifiedCode,
	})

	emailPayload := map[string]any{
		"email":   request.Email,
		"subject": "Welcome to SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Failed to marshal email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal email payload")

		status = "marshal_email_failed"
		return nil, user_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-auth-register", strconv.Itoa(newUser.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Failed to send email", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send email")

		status = "kafka_send_failed"
		return nil, user_errors.ErrFailedSendEmail
	}

	_, err = s.userRole.AssignRoleToUser(&requests.CreateUserRoleRequest{
		UserId: newUser.ID,
		RoleId: role.ID,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("REGISTER_ERR")

		s.logger.Error("Failed to assign role to user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to assign role to user")

		status = "assign_role_failed"
		return nil, userrole_errors.ErrFailedAssignRoleToUser
	}

	userResponse := s.mapping.ToUserResponse(newUser)

	s.logger.Debug("User registered successfully", zap.Int("user_id", newUser.ID), zap.String("email", request.Email))
	span.SetStatus(codes.Ok, "User registered successfully")

	return userResponse, nil
}

func (s *registerService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
