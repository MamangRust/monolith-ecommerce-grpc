package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-grpc-merchant/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/email"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type merchantCommandService struct {
	kafka                     kafka.Kafka
	ctx                       context.Context
	trace                     trace.Tracer
	userRepository            repository.UserQueryRepository
	merchantQueryRepository   repository.MerchantQueryRepository
	merchantCommandRepository repository.MerchantCommandRepository
	logger                    logger.LoggerInterface
	mapping                   response_service.MerchantResponseMapper
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewMerchantCommandService(kafka kafka.Kafka, ctx context.Context,
	userRepository repository.UserQueryRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	merchantCommandRepository repository.MerchantCommandRepository, logger logger.LoggerInterface, mapping response_service.MerchantResponseMapper) *merchantCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "merchant_command_service_requests_total",
			Help: "Total number of requests to the MerchantCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "merchant_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the MerchantCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantCommandService{
		kafka:                     kafka,
		ctx:                       ctx,
		trace:                     otel.Tracer("merchant-command-service"),
		merchantCommandRepository: merchantCommandRepository,
		userRepository:            userRepository,
		merchantQueryRepository:   merchantQueryRepository,
		logger:                    logger,
		mapping:                   mapping,
		requestCounter:            requestCounter,
		requestDuration:           requestDuration,
	}
}

func (s *merchantCommandService) CreateMerchant(request *requests.CreateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.String("name", request.Name),
	)

	s.logger.Debug("Creating new merchant", zap.String("merchant_name", request.Name))

	user, err := s.userRepository.FindById(request.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_USER")

		s.logger.Error("Failed to find user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find user")
		status = "failed_to_find_user"

		return nil, user_errors.ErrUserNotFoundRes
	}

	res, err := s.merchantCommandRepository.CreateMerchant(request)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_MERCHANT")

		s.logger.Error("Failed to create merchant", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create merchant")
		status = "failed_to_create_merchant"

		return nil, merchant_errors.ErrFailedCreateMerchant
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Your merchant account has been created successfully. To continue, please upload the required documents for verification. Once completed, our team will review and activate your account.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.ID),
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": "Initial Verification - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MERCHANT_CREATE_EMAIL_ERR")
		s.logger.Error("Failed to marshal merchant creation email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal merchant creation email payload")
		return nil, merchant_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-created", strconv.Itoa(res.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MERCHANT_CREATE_EMAIL_ERR")
		s.logger.Error("Failed to send merchant creation email message", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send merchant creation email")
		return nil, merchant_errors.ErrFailedSendEmail
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully created merchant", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantCommandService) UpdateMerchant(request *requests.UpdateMerchantRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", *request.MerchantID),
	)

	s.logger.Debug("Updating merchant", zap.Int("merchant_id", *request.MerchantID))

	res, err := s.merchantCommandRepository.UpdateMerchant(request)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_MERCHANT")

		s.logger.Error("Failed to update merchant", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update merchant")
		status = "failed_to_update_merchant"

		return nil, merchant_errors.ErrFailedUpdateMerchant
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully updated merchant", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantCommandService) UpdateMerchantStatus(request *requests.UpdateMerchantStatusRequest) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchantStatus", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchantStatus")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", *request.MerchantID),
	)

	s.logger.Debug("Updating merchant status", zap.Int("merchant_id", *request.MerchantID))

	merchant, err := s.merchantQueryRepository.FindById(*request.MerchantID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT")

		s.logger.Error("Failed to find merchant", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find merchant")
		status = "failed_to_find_merchant"

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	user, err := s.userRepository.FindById(merchant.UserID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_USER")

		s.logger.Error("Failed to find user", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find user")
		status = "failed_to_find_user"

		return nil, user_errors.ErrUserNotFoundRes
	}

	res, err := s.merchantCommandRepository.UpdateMerchantStatus(request)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_MERCHANT_STATUS")

		s.logger.Error("Failed to update merchant status", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update merchant status")
		status = "failed_to_update_merchant_status"

		return nil, merchant_errors.ErrFailedUpdateMerchant
	}

	statusReq := request.Status
	subject := ""
	message := ""
	buttonLabel := "Go to Portal"
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", *request.MerchantID)

	switch statusReq {
	case "active":
		subject = "Your Merchant Account is Now Active"
		message = "Congratulations! Your merchant account has been verified and is now <b>active</b>. You can now fully access all features in the SanEdge Merchant Portal."
	case "inactive":
		subject = "Merchant Account Set to Inactive"
		message = "Your merchant account status has been set to <b>inactive</b>. Please contact support if you believe this is a mistake."
	case "rejected":
		subject = "Merchant Account Rejected"
		message = "We're sorry to inform you that your merchant account has been <b>rejected</b>. Please contact support or review your submissions."
	default:
		return nil, nil
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   subject,
		"Message": message,
		"Button":  buttonLabel,
		"Link":    link,
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": subject,
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MERCHANT_STATUS_UPDATE_EMAIL_ERR")
		s.logger.Error("Failed to marshal merchant status update email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal merchant status update email payload")
		return nil, merchant_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-update-status", strconv.Itoa(*request.MerchantID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MERCHANT_STATUS_UPDATE_EMAIL_ERR")
		s.logger.Error("Failed to send merchant status update email", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send merchant status update email")
		return nil, merchant_errors.ErrFailedSendEmail
	}

	so := s.mapping.ToMerchantResponse(res)

	s.logger.Debug("Successfully updated merchant status", zap.Int("merchant_id", res.ID))

	return so, nil
}

func (s *merchantCommandService) TrashedMerchant(merchant_id int) (*response.MerchantResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Trashing merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantCommandRepository.TrashedMerchant(merchant_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_MERCHANT")

		s.logger.Error("Failed to trash merchant", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash merchant")
		status = "failed_to_trash_merchant"

		return nil, merchant_errors.ErrFailedTrashedMerchant
	}

	s.logger.Debug("Successfully trashed merchant", zap.Int("merchant_id", merchant_id))

	so := s.mapping.ToMerchantResponseDeleteAt(res)

	return so, nil
}

func (s *merchantCommandService) RestoreMerchant(merchant_id int) (*response.MerchantResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("Restoring merchant", zap.Int("merchant_id", merchant_id))

	res, err := s.merchantCommandRepository.RestoreMerchant(merchant_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_MERCHANT")

		s.logger.Error("Failed to restore merchant", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore merchant")
		status = "failed_to_restore_merchant"

		return nil, merchant_errors.ErrFailedRestoreMerchant
	}
	s.logger.Debug("Successfully restored merchant", zap.Int("merchant_id", merchant_id))

	so := s.mapping.ToMerchantResponse(res)

	return so, nil
}

func (s *merchantCommandService) DeleteMerchantPermanent(merchant_id int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteMerchantPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteMerchantPermanent")
	defer span.End()

	s.logger.Debug("Deleting merchant permanently", zap.Int("merchant_id", merchant_id))

	_, err := s.merchantCommandRepository.DeleteMerchantPermanent(merchant_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_MERCHANT")

		s.logger.Error("Failed to delete merchant permanently", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete merchant permanently")
		status = "failed_to_delete_merchant"

		return false, merchant_errors.ErrFailedDeleteMerchantPermanent
	}

	s.logger.Debug("Successfully deleted merchant permanently", zap.Int("merchant_id", merchant_id))

	return true, nil
}

func (s *merchantCommandService) RestoreAllMerchant() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllMerchant")
	defer span.End()

	s.logger.Debug("Restoring all merchants")

	_, err := s.merchantCommandRepository.RestoreAllMerchant()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_MERCHANTS")

		s.logger.Error("Failed to restore all merchants", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all merchants")
		status = "failed_to_restore_all_merchants"
		return false, merchant_errors.ErrFailedRestoreAllMerchants
	}

	s.logger.Debug("Successfully restored all merchants")
	return true, nil
}

func (s *merchantCommandService) DeleteAllMerchantPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllMerchantPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllMerchantPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all merchants")

	_, err := s.merchantCommandRepository.DeleteAllMerchantPermanent()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_MERCHANTS")

		s.logger.Error("Failed to permanently delete all merchants", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all merchants")
		status = "failed_delete_all_merchants"

		return false, merchant_errors.ErrFailedDeleteAllMerchantsPermanent
	}

	s.logger.Debug("Successfully deleted all merchants permanently")
	return true, nil
}

func (s *merchantCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
