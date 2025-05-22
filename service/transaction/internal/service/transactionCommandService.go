package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	traceunic "github.com/MamangRust/monolith-ecommerce-pkg/trace_unic"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	shippingaddress_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-ecommerce-shared/mapper/response/services"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	ctx                            context.Context
	trace                          trace.Tracer
	kafka                          kafka.Kafka
	userQueryRepository            repository.UserQueryRepository
	merchantQueryRepository        repository.MerchantQueryRepository
	transactionQueryRepository     repository.TransactionQueryRepository
	transactionCommandRepository   repository.TransactionCommandRepository
	orderQueryRepository           repository.OrderQueryRepository
	orderItemQueryRepository       repository.OrderItemRepository
	shippingAddressQueryRepository repository.ShippingAddressQueryRepository
	mapping                        response_service.TransactionResponseMapper
	logger                         logger.LoggerInterface
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewTransactionCommandService(
	ctx context.Context,
	kafka kafka.Kafka,
	userQueryRepository repository.UserQueryRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	transactionQueryRepository repository.TransactionQueryRepository,
	transactionCommandRepository repository.TransactionCommandRepository,
	orderQueryRepository repository.OrderQueryRepository,
	orderItemQueryRepository repository.OrderItemRepository,
	shippingAddressQueryRepository repository.ShippingAddressQueryRepository,
	mapping response_service.TransactionResponseMapper,
	logger logger.LoggerInterface,
) *transactionCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_command_service_request_total",
			Help: "Total number of requests to the TransactionCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_command_service_request_duration",
			Help:    "Histogram of request durations for the TransactionCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionCommandService{
		ctx:                            ctx,
		kafka:                          kafka,
		trace:                          otel.Tracer("transaction-command-service"),
		userQueryRepository:            userQueryRepository,
		merchantQueryRepository:        merchantQueryRepository,
		transactionQueryRepository:     transactionQueryRepository,
		transactionCommandRepository:   transactionCommandRepository,
		orderQueryRepository:           orderQueryRepository,
		orderItemQueryRepository:       orderItemQueryRepository,
		shippingAddressQueryRepository: shippingAddressQueryRepository,
		mapping:                        mapping,
		logger:                         logger,
		requestCounter:                 requestCounter,
		requestDuration:                requestDuration,
	}
}

func (s *transactionCommandService) CreateTransaction(req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "CreateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order.id", req.OrderID),
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("request.amount", req.Amount),
	)

	s.logger.Debug("Creating new transaction", zap.Int("orderID", req.OrderID))

	user, err := s.userQueryRepository.FindById(req.UserID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_USER")

		s.logger.Error("Failed to fetch user", zap.Int("userID", req.UserID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch user")

		status = "failed_fetch_user"
		return nil, user_errors.ErrUserNotFoundRes
	}

	_, err = s.merchantQueryRepository.FindById(req.MerchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT")
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found")

		status = "failed_find_merchant"
		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	_, err = s.orderQueryRepository.FindById(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER")
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Order not found")

		status = "failed_find_order"
		return nil, order_errors.ErrFailedFindOrderById
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEM_BY_ORDER")
		s.logger.Error("Failed to fetch order items", zap.Int("orderID", req.OrderID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch order items")

		status = "failed_fetch_order_items"
		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}
	if len(orderItems) == 0 {
		status = "empty_order_items"
		span.SetStatus(codes.Error, "No order items found")
		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	shipping, err := s.shippingAddressQueryRepository.FindByOrder(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_SHIPPING_ADDRESS")

		s.logger.Error("Failed to fetch shipping address", zap.Int("orderID", req.OrderID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch shipping address")

		status = "failed_fetch_shipping"
		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByOrder
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			status = "invalid_item_quantity"
			span.SetStatus(codes.Error, "Invalid item quantity")
			return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
		}
		totalAmount += item.Price*item.Quantity + shipping.ShippingCost
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn + shipping.ShippingCost

	span.SetAttributes(
		attribute.Int("calculated.amount", totalAmountWithTax),
	)

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "insufficient_balance"
		span.SetStatus(codes.Error, "Insufficient balance for transaction")
		return nil, transaction_errors.ErrFailedPaymentInsufficientBalance
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.CreateTransaction(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_TRANSACTION")
		s.logger.Error("Failed to create transaction record", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create transaction")

		status = "failed_create_transaction"
		return nil, transaction_errors.ErrFailedCreateTransaction
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

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TransactionErr")
		s.logger.Error("Failed to marshal transaction email payload", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal transaction email payload")
		return nil, transaction_errors.ErrFailedSendEmail
	}

	err = s.kafka.SendMessage("email-service-topic-transaction-create", strconv.Itoa(transaction.ID), payloadBytes)
	if err != nil {
		traceID := traceunic.GenerateTraceID("TransactionErr")
		s.logger.Error("Failed to send transaction email message", zap.String("trace_id", traceID), zap.Error(err))
		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send transaction email")
		return nil, transaction_errors.ErrFailedSendEmail
	}

	s.logger.Debug("Successfully created transaction", zap.Any("transaction", transaction))

	return s.mapping.ToTransactionResponse(transaction), nil
}
func (s *transactionCommandService) UpdateTransaction(req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateTransaction", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transaction.id", *req.TransactionID),
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("order.id", req.OrderID),
		attribute.Int("request.amount", req.Amount),
	)

	s.logger.Debug("Updating transaction", zap.Int("transactionID", *req.TransactionID))

	existingTx, err := s.transactionQueryRepository.FindById(*req.TransactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTION")
		s.logger.Error("Transaction not found", zap.Int("transactionID", *req.TransactionID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Transaction not found")

		status = "failed_find_transaction"
		return nil, transaction_errors.ErrFailedFindTransactionById
	}

	if existingTx.PaymentStatus == "success" || existingTx.PaymentStatus == "refunded" {
		span.SetStatus(codes.Error, "Payment status cannot be modified")
		status = "invalid_payment_status"
		return nil, transaction_errors.ErrFailedPaymentStatusCannotBeModified
	}

	_, err = s.merchantQueryRepository.FindById(req.MerchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT")
		s.logger.Error("Merchant not found", zap.Int("merchantId", req.MerchantID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found")

		status = "failed_find_merchant"
		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	_, err = s.orderQueryRepository.FindById(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER")
		s.logger.Error("Order not found", zap.Int("orderID", req.OrderID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Order not found")

		status = "failed_find_order"
		return nil, order_errors.ErrFailedFindOrderById
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FETCH_ORDER_ITEMS")
		s.logger.Error("Failed to retrieve order items", zap.Int("orderID", req.OrderID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve order items")

		status = "failed_fetch_order_items"
		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	shipping, err := s.shippingAddressQueryRepository.FindByOrder(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FETCH_SHIPPING")
		s.logger.Error("Failed to retrieve shipping address", zap.Int("orderID", req.OrderID), zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve shipping address")

		status = "failed_fetch_shipping"
		return nil, shippingaddress_errors.ErrFailedFindShippingAddressByOrder
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			status = "invalid_item_quantity"
			span.SetStatus(codes.Error, "Invalid order item quantity")
			return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
		}
		totalAmount += item.Price*item.Quantity + shipping.ShippingCost
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn + shipping.ShippingCost

	span.SetAttributes(attribute.Int("calculated.amount", totalAmountWithTax))

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "insufficient_balance"
		span.SetStatus(codes.Error, "Insufficient balance")
		return nil, transaction_errors.ErrFailedPaymentInsufficientBalance
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.UpdateTransaction(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION")
		s.logger.Error("Failed to update transaction", zap.String("trace_id", traceID), zap.Error(err))

		span.SetAttributes(attribute.String("trace.id", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transaction")

		status = "failed_update_transaction"
		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	s.logger.Debug("Successfully updated transaction", zap.Any("transaction", transaction))
	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionCommandService) TrashedTransaction(transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("TrashedTransaction", status, start) }()

	_, span := s.trace.Start(s.ctx, "TrashedTransaction")
	defer span.End()

	span.SetAttributes(attribute.Int("transaction.id", transactionID))
	s.logger.Debug("Trashing transaction", zap.Int("transaction_id", transactionID))

	res, err := s.transactionCommandRepository.TrashTransaction(transactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_TRANSACTION")
		s.logger.Error("Failed to move transaction to trash",
			zap.Int("transaction_id", transactionID), zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash transaction")
		status = "failed_trash"
		return nil, transaction_errors.ErrFailedTrashedTransaction
	}

	span.SetStatus(codes.Ok, "Transaction trashed")
	s.logger.Debug("Successfully trashed transaction", zap.Int("transaction_id", transactionID))
	return s.mapping.ToTransactionResponseDeleteAt(res), nil
}

func (s *transactionCommandService) RestoreTransaction(transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("RestoreTransaction", status, start) }()

	_, span := s.trace.Start(s.ctx, "RestoreTransaction")
	defer span.End()

	span.SetAttributes(attribute.Int("transaction.id", transactionID))
	s.logger.Debug("Restoring transaction", zap.Int("transaction_id", transactionID))

	res, err := s.transactionCommandRepository.RestoreTransaction(transactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_TRANSACTION")
		s.logger.Error("Failed to restore transaction",
			zap.Int("transaction_id", transactionID), zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore transaction")
		status = "failed_restore"
		return nil, transaction_errors.ErrFailedRestoreTransaction
	}

	span.SetStatus(codes.Ok, "Transaction restored")
	s.logger.Debug("Successfully restored transaction", zap.Int("transaction_id", transactionID))
	return s.mapping.ToTransactionResponseDeleteAt(res), nil
}

func (s *transactionCommandService) DeleteTransactionPermanently(transactionID int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("DeleteTransactionPermanently", status, start) }()

	_, span := s.trace.Start(s.ctx, "DeleteTransactionPermanently")
	defer span.End()

	span.SetAttributes(attribute.Int("transaction.id", transactionID))
	s.logger.Debug("Permanently deleting transaction", zap.Int("transactionID", transactionID))

	success, err := s.transactionCommandRepository.DeleteTransactionPermanently(transactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_TRANSACTION")
		s.logger.Error("Failed to permanently delete transaction",
			zap.Int("transaction_id", transactionID), zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete transaction")
		status = "failed_delete"
		return false, transaction_errors.ErrFailedDeleteTransactionPermanently
	}

	span.SetStatus(codes.Ok, "Transaction permanently deleted")
	return success, nil
}

func (s *transactionCommandService) RestoreAllTransactions() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("RestoreAllTransactions", status, start) }()

	_, span := s.trace.Start(s.ctx, "RestoreAllTransactions")
	defer span.End()

	s.logger.Debug("Restoring all trashed transactions")

	success, err := s.transactionCommandRepository.RestoreAllTransactions()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL")
		s.logger.Error("Failed to restore all trashed transactions", zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all transactions")
		status = "failed_restore_all"
		return false, transaction_errors.ErrFailedRestoreAllTransactions
	}

	span.SetStatus(codes.Ok, "All transactions restored")
	return success, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("DeleteAllTransactionPermanent", status, start) }()

	_, span := s.trace.Start(s.ctx, "DeleteAllTransactionPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all transactions")

	success, err := s.transactionCommandRepository.DeleteAllTransactionPermanent()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL")
		s.logger.Error("Failed to permanently delete all transactions", zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all transactions")
		status = "failed_delete_all"
		return false, transaction_errors.ErrFailedDeleteAllTransactionPermanent
	}

	span.SetStatus(codes.Ok, "All transactions permanently deleted")
	return success, nil
}

func (s *transactionCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
