package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/event"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"

	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	observability      observability.TraceLoggerObservability
	kafka              *kafka.Kafka
	outbox             repository.OutboxRepository
	pool               *pgxpool.Pool
	cache              cache.TransactionCommandCache
	transactionQuery   repository.TransactionQueryRepository
	transactionCommand repository.TransactionCommandRepository
	userQuery          repository.UserQueryRepository
	merchantQuery      repository.MerchantQueryRepository
	orderQuery         repository.OrderQueryRepository
	orderItem          repository.OrderItemRepository
	shippingAddress    repository.ShippingAddressQueryRepository
	logger             logger.LoggerInterface
}

type TransactionCommandServiceDeps struct {
	Observability      observability.TraceLoggerObservability
	Kafka              *kafka.Kafka
	Outbox             repository.OutboxRepository
	Pool               *pgxpool.Pool
	Cache              cache.TransactionCommandCache
	TransactionQuery   repository.TransactionQueryRepository
	TransactionCommand repository.TransactionCommandRepository
	UserQuery          repository.UserQueryRepository
	MerchantQuery      repository.MerchantQueryRepository
	OrderQuery         repository.OrderQueryRepository
	OrderItem          repository.OrderItemRepository
	ShippingAddress    repository.ShippingAddressQueryRepository
	Logger             logger.LoggerInterface
}

func NewTransactionCommandService(deps *TransactionCommandServiceDeps) TransactionCommandService {
	return &transactionCommandService{
		observability:      deps.Observability,
		kafka:              deps.Kafka,
		outbox:             deps.Outbox,
		pool:               deps.Pool,
		cache:              deps.Cache,
		transactionQuery:   deps.TransactionQuery,
		transactionCommand: deps.TransactionCommand,
		userQuery:          deps.UserQuery,
		merchantQuery:      deps.MerchantQuery,
		orderQuery:         deps.OrderQuery,
		orderItem:          deps.OrderItem,
		shippingAddress:    deps.ShippingAddress,
		logger:             deps.Logger,
	}
}

func (s *transactionCommandService) Create(ctx context.Context, req *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", req.UserID),
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("order_id", req.OrderID))

	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindByID(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	_, err = s.merchantQuery.FindByID(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	_, err = s.orderQuery.FindByID(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	orderItems, err := s.orderItem.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	shipping, err := s.shippingAddress.FindByID(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	var merchandiseTotal int
	for _, item := range orderItems {
		if item.Quantity <= 0 || item.Price < 0 {
			status = "error"
			return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, transaction_errors.ErrFailedOrderItemEmpty, method, span)
		}
		merchandiseTotal += int(item.Price) * int(item.Quantity)
	}

	// Shipping is charged once per order, not once per item.
	totalAmount := merchandiseTotal + int(shipping.ShippingCost)
	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	span.SetAttributes(attribute.Int("calculated_amount", totalAmountWithTax))

	// Validate the caller-supplied status (if any) before deriving the final
	// status. A create may only start from a canonical non-terminal status.
	if req.PaymentStatus != nil && *req.PaymentStatus != "" {
		if !IsValidPaymentStatus(*req.PaymentStatus) {
			status = "error"
			return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentStatusInvalid, method, span)
		}
	}

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = PaymentStatusSuccess
	} else {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentInsufficientBalance, method, span)
	}

	if req.PaymentStatus != nil && *req.PaymentStatus != "" && !CanTransitionPaymentStatus(*req.PaymentStatus, paymentStatus) {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transaction Successful",
		"Message": fmt.Sprintf("Your transaction of %d has been processed successfully.", req.Amount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/transaction/history",
	})

	payloadBytes, err := event.MarshalEmail("transaction.created", user.Email, "Transaction Successful - SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	// Transactional outbox: when a pool is available, the business insert and the
	// outbox event commit in a single database transaction, so a crash between
	// them cannot lose the event. The relay then publishes with retry and
	// dead-letter semantics, so a Kafka failure cannot roll back the payment.
	//
	// Without a pool this falls back to best-effort enqueue after commit, which
	// is NON-ATOMIC: a crash between the commit and the enqueue loses the event
	// silently. This path is intended for tests and local development only and
	// must not be relied on in production; production always supplies a pool.
	var transaction *db.CreateTransactionRow
	var tx pgx.Tx
	if s.pool != nil {
		tx, err = s.pool.Begin(ctx)
		if err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
		}
		defer func() { _ = tx.Rollback(ctx) }()
		transaction, err = s.transactionCommand.CreateInTx(ctx, tx, req)
	} else {
		if s.outbox != nil {
			s.logger.Warn("transactional outbox running in NON-ATOMIC fallback mode: no pgx pool configured; event loss is possible between commit and enqueue")
		}
		transaction, err = s.transactionCommand.Create(ctx, req)
	}
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	if s.outbox != nil {
		merchantPayload, marshalErr := json.Marshal(map[string]any{
			"merchantId":    transaction.MerchantID,
			"transactionId": transaction.TransactionID,
			"amount":        transaction.Amount,
			"status":        transaction.PaymentStatus,
			"timestamp":     time.Now().UnixMilli(),
		})
		if marshalErr != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, marshalErr, method, span)
		}

		enqueue := func(topic, key string, payload []byte) error {
			if tx != nil {
				_, enqueueErr := s.outbox.CreateInTx(ctx, tx, topic, key, payload)
				return enqueueErr
			}
			_, enqueueErr := s.outbox.Create(ctx, topic, key, payload)
			return enqueueErr
		}

		for _, event := range []struct {
			topic   string
			key     string
			payload []byte
		}{
			{topic: "email-service-topic-transaction-create", key: strconv.Itoa(int(transaction.TransactionID)), payload: payloadBytes},
			{topic: "merchant-service-topic-transaction-event", key: strconv.Itoa(int(transaction.MerchantID)), payload: merchantPayload},
		} {
			if enqueueErr := enqueue(event.topic, event.key, event.payload); enqueueErr != nil {
				s.logger.Error("failed to enqueue outbox event", zap.Error(enqueueErr), zap.String("topic", event.topic), zap.Int32("transaction_id", transaction.TransactionID))
				// In the transactional path the whole unit rolls back so business data
				// and both events stay consistent; in the fallback path the error is
				// surfaced through structured logs for manual reconciliation.
				if tx != nil {
					status = "error"
					return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, enqueueErr, method, span)
				}
			}
		}
	}

	if tx != nil {
		if err := tx.Commit(ctx); err != nil {
			status = "error"
			return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
		}
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully created transaction", zap.Int32("transaction_id", transaction.TransactionID))

	return transaction, nil
}

func (s *transactionCommandService) Update(ctx context.Context, req *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", *req.TransactionID))

	defer func() {
		end(status)
	}()

	existingTx, err := s.transactionQuery.FindByID(ctx, *req.TransactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	// Payment status changes must follow the canonical state machine. The
	// requested status (if provided) must be canonical and a legal transition
	// from the current status; terminal statuses cannot be reopened.
	if req.PaymentStatus != nil && *req.PaymentStatus != "" {
		if !IsValidPaymentStatus(*req.PaymentStatus) {
			status = "error"
			return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentStatusInvalid, method, span)
		}
		if *req.PaymentStatus != existingTx.PaymentStatus && !CanTransitionPaymentStatus(existingTx.PaymentStatus, *req.PaymentStatus) {
			status = "error"
			return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
		}
	}

	if req.MerchantID == 0 {
		req.MerchantID = int(existingTx.MerchantID)
	}
	_, err = s.merchantQuery.FindByID(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	if req.OrderID == 0 {
		req.OrderID = int(existingTx.OrderID)
	}
	_, err = s.orderQuery.FindByID(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	orderItems, err := s.orderItem.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	shipping, err := s.shippingAddress.FindByID(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	var merchandiseTotal int
	for _, item := range orderItems {
		if item.Quantity <= 0 || item.Price < 0 {
			status = "error"
			return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedOrderItemEmpty, method, span)
		}
		merchandiseTotal += int(item.Price) * int(item.Quantity)
	}

	totalAmount := merchandiseTotal + int(shipping.ShippingCost)

	if req.Amount == 0 {
		req.Amount = int(existingTx.Amount)
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = existingTx.PaymentMethod
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	// Derive the final status from the verified amount. The derived status must
	// itself be a legal transition from the current status (e.g. a failed
	// transaction cannot be reopened as success by a later update).
	paymentStatus := PaymentStatusSuccess
	if req.Amount < totalAmountWithTax {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentInsufficientBalance, method, span)
	}
	if req.PaymentStatus != nil && *req.PaymentStatus != "" && CanTransitionPaymentStatus(existingTx.PaymentStatus, *req.PaymentStatus) && *req.PaymentStatus != PaymentStatusSuccess {
		paymentStatus = *req.PaymentStatus
	}
	if !CanTransitionPaymentStatus(existingTx.PaymentStatus, paymentStatus) {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommand.Update(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, *req.TransactionID)

	logSuccess("Successfully updated transaction", zap.Int32("transaction_id", transaction.TransactionID))

	return transaction, nil
}

func (s *transactionCommandService) Trash(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommand.Trash(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully trashed transaction", zap.Int("transaction_id", transactionID))

	return res, nil
}

func (s *transactionCommandService) Restore(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommand.Restore(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully restored transaction", zap.Int("transaction_id", transactionID))

	return res, nil
}

func (s *transactionCommandService) DeletePermanent(ctx context.Context, transactionID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeletePermanent(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction_id", transactionID))

	return success, nil
}

func (s *transactionCommandService) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	const method = "DeleteByOrderIDPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeleteByOrderIDPermanent(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.InvalidateTransactionCache(ctx)

	logSuccess("Successfully permanently deleted transactions by order", zap.Int("order_id", orderID))

	return success, nil
}

func (s *transactionCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.RestoreAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.InvalidateTransactionCache(ctx)
	logSuccess("Successfully restored all transactions")

	return success, nil
}

func (s *transactionCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeleteAll(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	s.cache.InvalidateTransactionCache(ctx)
	logSuccess("Successfully permanently deleted all transactions")

	return success, nil
}
