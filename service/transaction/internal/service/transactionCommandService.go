package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"

	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	observability   observability.TraceLoggerObservability
	kafka           *kafka.Kafka
	cache           cache.TransactionCommandCache
	transactionQuery   repository.TransactionQueryRepository
	transactionCommand repository.TransactionCommandRepository
	userQuery       repository.UserQueryRepository
	merchantQuery   repository.MerchantQueryRepository
	orderQuery      repository.OrderQueryRepository
	orderItem       repository.OrderItemRepository
	shippingAddress repository.ShippingAddressQueryRepository
	logger          logger.LoggerInterface
}

type TransactionCommandServiceDeps struct {
	Observability      observability.TraceLoggerObservability
	Kafka              *kafka.Kafka
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

func NewTransactionCommandService(deps *TransactionCommandServiceDeps) *transactionCommandService {
	return &transactionCommandService{
		observability:      deps.Observability,
		kafka:              deps.Kafka,
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

func (s *transactionCommandService) CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*db.CreateTransactionRow, error) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user_id", req.UserID),
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("order_id", req.OrderID))

	defer func() {
		end(status)
	}()

	user, err := s.userQuery.FindById(ctx, req.UserID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	_, err = s.merchantQuery.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	_, err = s.orderQuery.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	orderItems, err := s.orderItem.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	shipping, err := s.shippingAddress.FindByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			status = "error"
			return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, transaction_errors.ErrFailedOrderItemEmpty, method, span)
		}
		totalAmount += int(item.Price)*int(item.Quantity) + int(shipping.ShippingCost)
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn + int(shipping.ShippingCost)

	span.SetAttributes(attribute.Int("calculated_amount", totalAmountWithTax))

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentInsufficientBalance, method, span)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommand.CreateTransaction(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateTransactionRow](s.logger, err, method, span)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Transaction Successful",
		"Message": fmt.Sprintf("Your transaction of %d has been processed successfully.", req.Amount),
		"Button":  "View History",
		"Link":    "https://sanedge.example.com/transaction/history",
	})

	emailPayload := map[string]any{
		"email":   user.Email,
		"subject": "Transaction Successful - SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, _ := json.Marshal(emailPayload)
	err = s.kafka.SendMessage("email-service-topic-transaction-create", strconv.Itoa(int(transaction.TransactionID)), payloadBytes)
	if err != nil {
		s.logger.Error("failed to send kafka message", zap.Error(err))
	}

	logSuccess("Successfully created transaction", zap.Int32("transaction_id", transaction.TransactionID))

	return transaction, nil
}

func (s *transactionCommandService) UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*db.UpdateTransactionRow, error) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", *req.TransactionID))

	defer func() {
		end(status)
	}()

	existingTx, err := s.transactionQuery.FindById(ctx, *req.TransactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	if existingTx.PaymentStatus == "success" || existingTx.PaymentStatus == "refunded" {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
	}

	_, err = s.merchantQuery.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	_, err = s.orderQuery.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	orderItems, err := s.orderItem.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	shipping, err := s.shippingAddress.FindByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += int(item.Price)*int(item.Quantity) + int(shipping.ShippingCost)
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn + int(shipping.ShippingCost)

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentInsufficientBalance, method, span)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommand.UpdateTransaction(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, *req.TransactionID)

	logSuccess("Successfully updated transaction", zap.Int32("transaction_id", transaction.TransactionID))

	return transaction, nil
}

func (s *transactionCommandService) TrashedTransaction(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommand.TrashTransaction(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](s.logger, err, method, span)
	}

	s.cache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully trashed transaction", zap.Int("transaction_id", transactionID))

	return res, nil
}

func (s *transactionCommandService) RestoreTransaction(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommand.RestoreTransaction(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.Transaction](s.logger, err, method, span)
	}

	logSuccess("Successfully restored transaction", zap.Int("transaction_id", transactionID))

	return res, nil
}

func (s *transactionCommandService) DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, error) {
	const method = "DeleteTransactionPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction_id", transactionID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeleteTransactionPermanently(ctx, transactionID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction_id", transactionID))

	return success, nil
}

func (s *transactionCommandService) RestoreAllTransactions(ctx context.Context) (bool, error) {
	const method = "RestoreAllTransactions"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.RestoreAllTransactions(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	logSuccess("Successfully restored all transactions")

	return success, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommand.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](s.logger, err, method, span)
	}

	logSuccess("Successfully permanently deleted all transactions")

	return success, nil
}
