package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantCommandService struct {
	kafka              *kafka.Kafka
	observability      observability.TraceLoggerObservability
	cache              cache.MerchantCommandCache
	merchantRepository repository.MerchantCommandRepository
	merchantQuery      repository.MerchantQueryRepository
	userRepository     repository.UserQueryRepository
	pool               *pgxpool.Pool
	outbox             *outbox.OutboxService
	logger             logger.LoggerInterface
}

type MerchantCommandServiceDeps struct {
	Kafka              *kafka.Kafka
	Observability      observability.TraceLoggerObservability
	Cache              cache.MerchantCommandCache
	MerchantRepository repository.MerchantCommandRepository
	MerchantQuery      repository.MerchantQueryRepository
	UserRepository     repository.UserQueryRepository
	Pool               *pgxpool.Pool
	Outbox             *outbox.OutboxService
	Logger             logger.LoggerInterface
}

func NewMerchantCommandService(deps *MerchantCommandServiceDeps) MerchantCommandService {
	return &merchantCommandService{
		kafka:              deps.Kafka,
		observability:      deps.Observability,
		cache:              deps.Cache,
		merchantRepository: deps.MerchantRepository,
		merchantQuery:      deps.MerchantQuery,
		userRepository:     deps.UserRepository,
		pool:               deps.Pool,
		outbox:             deps.Outbox,
		logger:             deps.Logger,
	}
}

func (s *merchantCommandService) Create(ctx context.Context, request *requests.CreateMerchantRequest) (*db.CreateMerchantRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("user.id", request.UserID))

	defer func() {
		end(status)
	}()

	user, err := s.userRepository.FindByID(ctx, request.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", request.UserID),
		)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge Merchant Portal",
		"Message": "Your merchant account has been created successfully. To continue, please upload the required documents for verification. Once completed, our team will review and activate your account.",
		"Button":  "Upload Documents",
		"Link":    fmt.Sprintf("https://sanedge.example.com/merchant/%d/documents", user.UserID),
	})

	payloadBytes, err := event.MarshalEmail("merchant.created", user.Email, "Initial Verification - SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", request.UserID),
		)
	}

	var res *db.CreateMerchantRow
	if s.pool != nil {
		tx, beginErr := s.pool.Begin(ctx)
		if beginErr != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateMerchantRow](s.logger, beginErr, method, span, zap.Int("user.id", request.UserID))
		}
		defer func() { _ = tx.Rollback(ctx) }()

		res, err = s.merchantRepository.CreateInTx(ctx, tx, request)
		if err == nil && s.outbox != nil {
			err = s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-merchant-create", strconv.Itoa(int(res.MerchantID)), payloadBytes)
		}
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateMerchantRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("user.id", request.UserID),
			)
		}
		if err := tx.Commit(ctx); err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateMerchantRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("user.id", request.UserID),
			)
		}
	} else {
		res, err = s.merchantRepository.Create(ctx, request)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateMerchantRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("user.id", request.UserID),
			)
		}
		if s.kafka != nil {
			if sendErr := s.kafka.SendMessage("email-service-topic-merchant-create", strconv.Itoa(int(res.MerchantID)), payloadBytes); sendErr != nil {
				s.logger.Error("Failed to send email to Kafka", zap.Error(sendErr))
			}
		}
	}

	logSuccess("Successfully created merchant", zap.Int("merchant.id", int(res.MerchantID)))

	return res, nil
}

func (s *merchantCommandService) Update(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.UpdateMerchantRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", *request.MerchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.Update(ctx, request)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, *request.MerchantID)

	logSuccess("Successfully updated merchant", zap.Int("merchant.id", *request.MerchantID))

	return res, nil
}

func (s *merchantCommandService) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.UpdateMerchantStatusRow, error) {
	const method = "UpdateMerchantStatus"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", *request.MerchantID))

	defer func() {
		end(status)
	}()

	merchant, err := s.merchantQuery.FindByID(ctx, *request.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantStatusRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", *request.MerchantID),
		)
	}

	user, err := s.userRepository.FindByID(ctx, int(merchant.UserID))
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateMerchantStatusRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("user.id", int(merchant.UserID)),
		)
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
	}

	var res *db.UpdateMerchantStatusRow
	if s.pool != nil {
		tx, beginErr := s.pool.Begin(ctx)
		if beginErr != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateMerchantStatusRow](s.logger, beginErr, method, span, zap.Int("merchant.id", *request.MerchantID))
		}
		defer func() { _ = tx.Rollback(ctx) }()

		res, err = s.merchantRepository.UpdateStatusInTx(ctx, tx, request)
		if err == nil && subject != "" && s.outbox != nil {
			htmlBody := email.GenerateEmailHTML(map[string]string{
				"Title":   subject,
				"Message": message,
				"Button":  buttonLabel,
				"Link":    link,
			})
			payloadBytes, marshalErr := event.MarshalEmail("merchant.status_updated", user.Email, subject, htmlBody)
			if marshalErr != nil {
				s.logger.Error("failed to marshal merchant status email", zap.Error(marshalErr), zap.Int32("merchant_id", res.MerchantID))
			} else {
				err = s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-merchant-update-status", strconv.Itoa(int(res.MerchantID)), payloadBytes)
			}
		}
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateMerchantStatusRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("merchant.id", *request.MerchantID),
			)
		}
		if err := tx.Commit(ctx); err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateMerchantStatusRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("merchant.id", *request.MerchantID),
			)
		}
	} else {
		res, err = s.merchantRepository.UpdateStatus(ctx, request)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.UpdateMerchantStatusRow](
				s.logger,
				err,
				method,
				span,
				zap.Int("merchant.id", *request.MerchantID),
			)
		}
		if subject != "" {
			htmlBody := email.GenerateEmailHTML(map[string]string{
				"Title":   subject,
				"Message": message,
				"Button":  buttonLabel,
				"Link":    link,
			})
			payloadBytes, marshalErr := event.MarshalEmail("merchant.status_updated", user.Email, subject, htmlBody)
			if marshalErr != nil {
				s.logger.Error("failed to marshal merchant status email", zap.Error(marshalErr), zap.Int32("merchant_id", res.MerchantID))
			} else if s.kafka != nil {
				if sendErr := s.kafka.SendMessage("email-service-topic-merchant-update-status", strconv.Itoa(int(res.MerchantID)), payloadBytes); sendErr != nil {
					s.logger.Error("failed to publish merchant status email", zap.Error(sendErr), zap.Int32("merchant_id", res.MerchantID))
				}
			}
		}
	}

	// The merchant status event for the transaction service remains a direct
	// Kafka publish (it is not an email and has no outbox contract).
	if s.kafka != nil && subject != "" {
		if statusEvent, marshalErr := json.Marshal(map[string]any{
			"merchantId": res.MerchantID,
			"status":     request.Status,
			"timestamp":  time.Now().UnixMilli(),
		}); marshalErr != nil {
			s.logger.Error("failed to marshal merchant status event", zap.Error(marshalErr), zap.Int32("merchant_id", res.MerchantID))
		} else if sendErr := s.kafka.SendMessage("transaction-service-topic-merchant-status-event", strconv.Itoa(int(res.MerchantID)), statusEvent); sendErr != nil {
			s.logger.Error("failed to publish merchant status event", zap.Error(sendErr), zap.Int32("merchant_id", res.MerchantID))
		}
	}

	s.cache.DeleteCachedMerchant(ctx, *request.MerchantID)

	logSuccess("Successfully updated merchant status", zap.Int("merchant.id", *request.MerchantID))

	return res, nil
}

func (s *merchantCommandService) Trash(ctx context.Context, merchantID int) (*db.Merchant, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.Trash(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", merchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, merchantID)

	logSuccess("Successfully trashed merchant", zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *merchantCommandService) Restore(ctx context.Context, merchantID int) (*db.Merchant, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.Restore(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Merchant](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", merchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, merchantID)

	logSuccess("Successfully restored merchant", zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *merchantCommandService) DeletePermanent(ctx context.Context, merchantID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.DeletePermanent(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", merchantID),
		)
	}

	s.cache.DeleteCachedMerchant(ctx, merchantID)

	logSuccess("Successfully permanently deleted merchant", zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *merchantCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully restored all merchants")

	return res, nil
}

func (s *merchantCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.merchantRepository.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			err,
			method,
			span,
		)
	}

	logSuccess("Successfully permanently deleted all merchants")

	return res, nil
}
