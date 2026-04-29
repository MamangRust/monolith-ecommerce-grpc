package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
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
	logger             logger.LoggerInterface
}

type MerchantCommandServiceDeps struct {
	Kafka              *kafka.Kafka
	Observability      observability.TraceLoggerObservability
	Cache              cache.MerchantCommandCache
	MerchantRepository repository.MerchantCommandRepository
	MerchantQuery      repository.MerchantQueryRepository
	UserRepository     repository.UserQueryRepository
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

	res, err := s.merchantRepository.Create(ctx, request)
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

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": "Initial Verification - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant.id", int(res.MerchantID)),
		)
	}

	err = s.kafka.SendMessage("email-service-topic-merchant-created", strconv.Itoa(int(res.MerchantID)), payloadBytes)
	if err != nil {
		s.logger.Error("Failed to send email to Kafka", zap.Error(err))
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

	res, err := s.merchantRepository.UpdateStatus(ctx, request)
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

	if subject != "" {
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

		payloadBytes, _ := json.Marshal(emailPayload)
		_ = s.kafka.SendMessage("email-service-topic-merchant-status-updated", strconv.Itoa(int(res.MerchantID)), payloadBytes)
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
