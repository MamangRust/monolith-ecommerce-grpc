package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	merchantdocument_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_document_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type merchantDocumentCommandService struct {
	kafka                             kafka.Kafka
	ctx                               context.Context
	trace                             trace.Tracer
	merchantQueryRepository           repository.MerchantQueryRepository
	merchantDocumentCommandRepository repository.MerchantDocumentCommandRepository
	userRepository                    repository.UserQueryRepository
	logger                            logger.LoggerInterface
	mapping                           response_service.MerchantDocumentResponseMapper
	requestCounter                    *prometheus.CounterVec
	requestDuration                   *prometheus.HistogramVec
}

func NewMerchantDocumentCommandService(
	kafka kafka.Kafka,
	ctx context.Context,
	merchantDocumentCommandRepository repository.MerchantDocumentCommandRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	userRepository repository.UserQueryRepository,
	logger logger.LoggerInterface,
	mapping response_service.MerchantDocumentResponseMapper,
) *merchantDocumentCommandService {
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "merchant_document_command_request_count",
		Help: "Number of merchant document command requests MerchantDocumentCommandService",
	}, []string{"status"})

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "merchant_document_command_request_duration_seconds",
		Help:    "The duration of requests MerchantDocumentCommandService",
		Buckets: prometheus.DefBuckets,
	}, []string{"status"})

	prometheus.MustRegister(requestCounter, requestDuration)

	return &merchantDocumentCommandService{
		kafka:                             kafka,
		ctx:                               ctx,
		trace:                             otel.Tracer("merchant-document-command-service"),
		merchantQueryRepository:           merchantQueryRepository,
		merchantDocumentCommandRepository: merchantDocumentCommandRepository,
		userRepository:                    userRepository,
		logger:                            logger,
		mapping:                           mapping,
		requestCounter:                    requestCounter,
		requestDuration:                   requestDuration,
	}
}

func (s *merchantDocumentCommandService) CreateMerchantDocument(request *requests.CreateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateMerchantDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateMerchantDocument")
	defer span.End()

	span.SetAttributes(
		attribute.String("document-type", request.DocumentType),
	)

	s.logger.Debug("Creating new merchant document", zap.String("document-type", request.DocumentType))

	merchant, err := s.merchantQueryRepository.FindById(request.MerchantID)

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

	merchantDocument, err := s.merchantDocumentCommandRepository.CreateMerchantDocument(request)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_MERCHANT_DOCUMENT")

		s.logger.Error("Failed to create merchant document", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create merchant document")
		status = "failed_to_create_merchant_document"

		return nil, merchantdocument_errors.ErrFailedCreateMerchantDocument
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Thank you for registering your merchant account. Your account is currently <b>inactive</b> and under initial review. To proceed, please upload all required documents for verification. Once your documents are submitted, our team will review them and activate your account accordingly.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.ID),
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": "Merchant Verification Pending - Action Required",
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

	err = s.kafka.SendMessage("email-service-topic-merchant-created", strconv.Itoa(merchantDocument.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MERCHANT_CREATE_EMAIL_ERR")
		s.logger.Error("Failed to send merchant creation email message", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send merchant creation email")
		return nil, merchant_errors.ErrFailedSendEmail
	}

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.logger.Debug("Successfully created merchant document", zap.String("document-type", request.DocumentType))

	return so, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocument(request *requests.UpdateMerchantDocumentRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchantDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchantDocument")
	defer span.End()

	span.SetAttributes(
		attribute.Int("document_id", *request.DocumentID),
	)

	s.logger.Debug("Updating merchant document", zap.Int("document_id", *request.DocumentID))

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocument(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_MERCHANT_DOCUMENT")

		s.logger.Error("Failed to update merchant document", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update merchant document")
		status = "failed_to_update_merchant_document"

		return nil, merchantdocument_errors.ErrFailedUpdateMerchantDocument
	}

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.logger.Debug("Successfully updated merchant document", zap.Int("document_id", *request.DocumentID))

	return so, nil
}

func (s *merchantDocumentCommandService) UpdateMerchantDocumentStatus(request *requests.UpdateMerchantDocumentStatusRequest) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateMerchantDocumentStatus", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateMerchantDocumentStatus")
	defer span.End()

	span.SetAttributes(
		attribute.Int("document_id", *request.DocumentID),
	)

	s.logger.Debug("Updating merchant document status", zap.Int("document_id", *request.DocumentID))

	merchant, err := s.merchantQueryRepository.FindById(request.MerchantID)

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

	merchantDocument, err := s.merchantDocumentCommandRepository.UpdateMerchantDocumentStatus(request)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_MERCHANT_DOCUMENT_STATUS")

		s.logger.Error("Failed to update merchant document status", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update merchant document status")
		status = "failed_to_update_merchant_document_status"

		return nil, merchantdocument_errors.ErrFailedUpdateMerchantDocument
	}

	statusReq := request.Status
	note := request.Note
	subject := ""
	message := ""
	buttonLabel := ""
	link := fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", request.MerchantID)

	switch statusReq {
	case "pending":
		subject = "Merchant Document Status: Pending Review"
		message = "Your merchant documents have been submitted and are currently pending review."
		buttonLabel = "View Documents"
	case "approved":
		subject = "Merchant Document Status: Approved"
		message = "Congratulations! Your merchant documents have been approved. Your account is now active and fully functional."
		buttonLabel = "Go to Dashboard"
		link = fmt.Sprintf("https://sanedge.example.com/merchant/%d/dashboard", request.MerchantID)
	case "rejected":
		subject = "Merchant Document Status: Rejected"
		message = "Unfortunately, your merchant documents were rejected. Please review the feedback below and re-upload the necessary documents."
		buttonLabel = "Re-upload Documents"
	default:
		return nil, nil
	}

	if note != "" {
		message += fmt.Sprintf(`<br><br><b>Reviewer Note:</b><br><i>%s</i>`, note)
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
		traceID := traceunic.GenerateTraceID("MERCHANT_DOC_STATUS_EMAIL_ERR")
		s.logger.Error("Failed to marshal merchant document status email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal merchant document status email payload")
		return nil, merchant_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-document-update-status", strconv.Itoa(request.MerchantID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("MERCHANT_DOC_STATUS_EMAIL_ERR")
		s.logger.Error("Failed to send merchant document status email", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send merchant document status email")
		return nil, merchant_errors.ErrFailedSendEmail
	}

	so := s.mapping.ToMerchantDocumentResponse(merchantDocument)

	s.logger.Debug("Successfully updated merchant document status", zap.Int("document_id", *request.DocumentID))

	return so, nil
}

func (s *merchantDocumentCommandService) TrashedMerchantDocument(documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedDocument")
	defer span.End()

	span.SetAttributes(attribute.Int("document_id", documentID))

	s.logger.Debug("Trashing merchant document", zap.Int("document_id", documentID))

	res, err := s.merchantDocumentCommandRepository.TrashedMerchantDocument(documentID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_DOCUMENT")

		s.logger.Error("Failed to trash document", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash document")
		status = "failed_to_trash_document"

		return nil, merchantdocument_errors.ErrFailedTrashMerchantDocument
	}

	s.logger.Debug("Successfully trashed document", zap.Int("document_id", documentID))

	return s.mapping.ToMerchantDocumentResponse(res), nil
}

func (s *merchantDocumentCommandService) RestoreMerchantDocument(documentID int) (*response.MerchantDocumentResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreDocument", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreDocument")
	defer span.End()

	span.SetAttributes(attribute.Int("document_id", documentID))

	s.logger.Debug("Restoring merchant document", zap.Int("document_id", documentID))

	res, err := s.merchantDocumentCommandRepository.RestoreMerchantDocument(documentID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_DOCUMENT")

		s.logger.Error("Failed to restore document", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore document")
		status = "failed_to_restore_document"

		return nil, merchantdocument_errors.ErrFailedRestoreMerchantDocument
	}

	s.logger.Debug("Successfully restored document", zap.Int("document_id", documentID))

	return s.mapping.ToMerchantDocumentResponse(res), nil
}

func (s *merchantDocumentCommandService) DeleteMerchantDocumentPermanent(documentID int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteDocumentPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteDocumentPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting merchant document", zap.Int("document_id", documentID))

	_, err := s.merchantDocumentCommandRepository.DeleteMerchantDocumentPermanent(documentID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_DOCUMENT")

		s.logger.Error("Failed to permanently delete document", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete document permanently")
		status = "failed_to_delete_document"

		return false, merchantdocument_errors.ErrFailedDeleteMerchantDocument
	}

	s.logger.Debug("Successfully deleted document permanently", zap.Int("document_id", documentID))

	return true, nil
}

func (s *merchantDocumentCommandService) RestoreAllMerchantDocument() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllDocuments", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllDocuments")
	defer span.End()

	s.logger.Debug("Restoring all merchant documents")

	_, err := s.merchantDocumentCommandRepository.RestoreAllMerchantDocument()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_DOCUMENTS")

		s.logger.Error("Failed to restore all documents", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all documents")
		status = "failed_to_restore_all_documents"

		return false, merchantdocument_errors.ErrFailedRestoreAllMerchantDocuments
	}

	s.logger.Debug("Successfully restored all merchant documents")

	return true, nil
}

func (s *merchantDocumentCommandService) DeleteAllMerchantDocumentPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllDocumentsPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllDocumentsPermanent")
	defer span.End()

	s.logger.Debug("Deleting all merchant documents permanently")

	_, err := s.merchantDocumentCommandRepository.DeleteAllMerchantDocumentPermanent()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_DOCUMENTS")

		s.logger.Error("Failed to delete all documents permanently", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all documents permanently")
		status = "failed_to_delete_all_documents"

		return false, merchantdocument_errors.ErrFailedDeleteAllMerchantDocuments
	}

	s.logger.Debug("Successfully deleted all merchant documents permanently")

	return true, nil
}

func (s *merchantDocumentCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
