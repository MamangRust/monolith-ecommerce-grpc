package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
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

func NewTransactionCommandService(deps *TransactionCommandServiceDeps) TransactionCommandService {
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

	transaction, err := s.transactionCommand.Create(ctx, req)
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

	if (existingTx.PaymentStatus == "success" || existingTx.PaymentStatus == "refunded") && (req.PaymentStatus != nil && *req.PaymentStatus != existingTx.PaymentStatus) {
		status = "error"
		return errorhandler.HandleError[*db.UpdateTransactionRow](s.logger, transaction_errors.ErrFailedPaymentStatusCannotBeModified, method, span)
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

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += int(item.Price)*int(item.Quantity) + int(shipping.ShippingCost)
	}

	if req.Amount == 0 {
		req.Amount = int(existingTx.Amount)
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = existingTx.PaymentMethod
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

	logSuccess("Successfully permanently deleted all transactions")

	return success, nil
}
