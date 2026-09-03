package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/event"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-pkg/outbox"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantDocumentCommandService struct {
	observability observability.TraceLoggerObservability
	kafka         *kafka.Kafka
	cache         cache.MerchantDocumentCommandCache
	repository    repository.MerchantDocumentCommandRepository
	merchantQuery repository.MerchantQueryRepository
	userQuery     repository.UserQueryRepository
	pool          *pgxpool.Pool
	outbox        *outbox.OutboxService
	logger        logger.LoggerInterface
}

type MerchantDocumentCommandServiceDeps struct {
	Observability observability.TraceLoggerObservability
	Kafka         *kafka.Kafka
	Cache         cache.MerchantDocumentCommandCache
	Repository    repository.MerchantDocumentCommandRepository
	MerchantQuery repository.MerchantQueryRepository
	UserQuery     repository.UserQueryRepository
	Pool          *pgxpool.Pool
	Outbox        *outbox.OutboxService
	Logger        logger.LoggerInterface
}

func NewMerchantDocumentCommandService(deps *MerchantDocumentCommandServiceDeps) MerchantDocumentCommandService {
	return &merchantDocumentCommandService{
		observability: deps.Observability,
		kafka:         deps.Kafka,
		cache:         deps.Cache,
		repository:    deps.Repository,
		merchantQuery: deps.MerchantQuery,
		userQuery:     deps.UserQuery,
		pool:          deps.Pool,
		outbox:        deps.Outbox,
		logger:        deps.Logger,
	}
}

// publishDocumentEmail resolves the merchant owner's email and publishes the
// given topic with an {email, subject, body} payload keyed by document id.
// Email delivery must never fail the document operation, so errors are logged
// and swallowed, mirroring the merchant-account email flow.
func (s *merchantDocumentCommandService) publishDocumentEmail(ctx context.Context, tx pgx.Tx, topic, eventType string, documentID, merchantID int32, subject, message string) {
	if s.merchantQuery == nil || s.userQuery == nil {
		s.logger.Warn("merchant document email skipped: query dependencies missing")
		return
	}

	merchant, err := s.merchantQuery.FindByID(ctx, int(merchantID))
	if err != nil {
		s.logger.Error("failed to resolve merchant for document email", zap.Error(err), zap.Int32("merchant_id", merchantID))
		return
	}

	user, err := s.userQuery.FindByID(ctx, int(merchant.UserID))
	if err != nil {
		s.logger.Error("failed to resolve user for document email", zap.Error(err), zap.Int32("user_id", merchant.UserID))
		return
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   subject,
		"Message": message,
		"Button":  "Go to Portal",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", merchantID),
	})

	payloadBytes, err := event.MarshalEmail(eventType, user.Email, subject, htmlBody)
	if err != nil {
		s.logger.Error("failed to marshal merchant document email payload", zap.Error(err), zap.Int32("document_id", documentID))
		return
	}

	if tx != nil && s.outbox != nil {
		if err := s.outbox.EnqueueInTx(ctx, tx, topic, strconv.Itoa(int(documentID)), payloadBytes); err != nil {
			s.logger.Error("failed to enqueue merchant document email to outbox", zap.Error(err), zap.String("topic", topic), zap.Int32("document_id", documentID))
		}
		return
	}
	if s.kafka != nil {
		if err := s.kafka.SendMessage(topic, strconv.Itoa(int(documentID)), payloadBytes); err != nil {
			s.logger.Error("failed to publish merchant document email", zap.Error(err), zap.String("topic", topic), zap.Int32("document_id", documentID))
		}
	}
}

func (s *merchantDocumentCommandService) Create(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*db.CreateMerchantDocumentRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", request.MerchantID))

	defer func() {
		end(status)
	}()

	var res *db.CreateMerchantDocumentRow
	var tx pgx.Tx
	var err error
	if s.pool != nil {
		var beginErr error
		tx, beginErr = s.pool.Begin(ctx)
		if beginErr != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateMerchantDocumentRow](s.logger, beginErr, method, span, zap.Int("merchantID", request.MerchantID))
		}
		defer func() { _ = tx.Rollback(ctx) }()
		res, err = s.repository.CreateInTx(ctx, tx, request)
	} else {
		res, err = s.repository.Create(ctx, request)
	}
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantDocumentRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchantID", request.MerchantID),
		)
	}

	s.publishDocumentEmail(
		ctx,
		tx,
		"email-service-topic-merchant-document-create",
		"merchant_document.created",
		res.DocumentID,
		res.MerchantID,
		"Document Uploaded - SanEdge",
		"Your document has been uploaded successfully. Our team will review it and notify you once its status changes.",
	)

	if tx != nil {
		if err := tx.Commit(ctx); err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateMerchantDocumentRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("merchantID", request.MerchantID),
			)
		}
	}

	logSuccess("Successfully created merchant document", zap.Int("merchantID", request.MerchantID))

	return res, nil
}

func (s *merchantDocumentCommandService) Update(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*db.UpdateMerchantDocumentRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("document_id", *request.DocumentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Update(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantDocumentRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", *request.DocumentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(res.DocumentID))

	logSuccess("Successfully updated merchant document", zap.Int("document_id", *request.DocumentID))

	return res, nil
}

func (s *merchantDocumentCommandService) UpdateStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*db.UpdateMerchantDocumentStatusRow, error) {
	const method = "UpdateStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", *request.DocumentID))

	defer func() {
		end(status)
	}()

	var res *db.UpdateMerchantDocumentStatusRow
	var tx pgx.Tx
	var err error
	if s.pool != nil {
		var beginErr error
		tx, beginErr = s.pool.Begin(ctx)
		if beginErr != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateMerchantDocumentStatusRow](s.logger, beginErr, method, span, zap.Int("document_id", *request.DocumentID))
		}
		defer func() { _ = tx.Rollback(ctx) }()
		res, err = s.repository.UpdateStatusInTx(ctx, tx, request)
	} else {
		res, err = s.repository.UpdateStatus(ctx, request)
	}
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantDocumentStatusRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", *request.DocumentID),
		)
	}

	subject, message := "Document Status Updated - SanEdge", "Your document status has been updated."
	switch res.Status {
	case "approved":
		subject = "Document Approved - SanEdge"
		message = "Congratulations! Your document has been <b>approved</b>. Your merchant account is one step closer to activation."
	case "rejected":
		subject = "Document Rejected - SanEdge"
		message = "We're sorry, but your document has been <b>rejected</b>. Please review the note and re-upload the corrected document."
	}

	s.publishDocumentEmail(
		ctx,
		tx,
		"email-service-topic-merchant-document-update-status",
		"merchant_document.status_updated",
		res.DocumentID,
		res.MerchantID,
		subject,
		message,
	)

	if tx != nil {
		if err := tx.Commit(ctx); err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateMerchantDocumentStatusRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("document_id", *request.DocumentID),
			)
		}
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(res.DocumentID))

	logSuccess("Successfully updated merchant document status", zap.Int("document_id", *request.DocumentID))

	return res, nil
}

func (s *merchantDocumentCommandService) Trash(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", documentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Trash(ctx, documentID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", documentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(documentID))

	logSuccess("Successfully trashed merchant document", zap.Int("document_id", documentID))

	return res, nil
}

func (s *merchantDocumentCommandService) Restore(ctx context.Context, documentID int) (*db.MerchantDocument, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", documentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.Restore(ctx, documentID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.MerchantDocument](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", documentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(documentID))

	logSuccess("Successfully restored merchant document", zap.Int("document_id", documentID))

	return res, nil
}

func (s *merchantDocumentCommandService) DeletePermanent(ctx context.Context, documentID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchantDocument.id", documentID))

	defer func() {
		end(status)
	}()

	res, err := s.repository.DeletePermanent(ctx, documentID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("document_id", documentID),
		)
	}

	s.cache.DeleteCachedMerchantDocuments(ctx, int(documentID))

	logSuccess("Successfully permanently deleted merchant document", zap.Int("document_id", documentID))

	return res, nil
}

func (s *merchantDocumentCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.repository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all merchant documents")

	return res, nil
}

func (s *merchantDocumentCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.repository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all merchant documents")

	return res, nil
}
