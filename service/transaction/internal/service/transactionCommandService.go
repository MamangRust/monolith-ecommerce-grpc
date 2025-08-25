package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/email"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
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
	mencache                       mencache.TransactionCommandCache
	errorhandler                   errorhandler.TransactionCommandError
	trace                          trace.Tracer
	userQueryRepository            repository.UserQueryRepository
	merchantQueryRepository        repository.MerchantQueryRepository
	transactionQueryRepository     repository.TransactionQueryRepository
	transactionCommandRepository   repository.TransactionCommandRepository
	orderQueryRepository           repository.OrderQueryRepository
	orderItemQueryRepository       repository.OrderItemRepository
	shippingAddressQueryRepository repository.ShippingAddressQueryRepository
	mapping                        response_service.TransactionResponseMapper
	logger                         logger.LoggerInterface
	kafka                          *kafka.Kafka
	requestCounter                 *prometheus.CounterVec
	requestDuration                *prometheus.HistogramVec
}

func NewTransactionCommandService(
	mencache mencache.TransactionCommandCache,
	errorhandler errorhandler.TransactionCommandError,
	kafka *kafka.Kafka,
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
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionCommandService{
		mencache:                       mencache,
		errorhandler:                   errorhandler,
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

func (s *transactionCommandService) CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("user.id", req.UserID), attribute.Int("merchant.id", req.MerchantID), attribute.Int("order.id", req.OrderID))

	defer func() {
		end(status)
	}()

	user, err := s.userQueryRepository.FindById(ctx, req.UserID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_USER_BY_ID", span, &status, user_errors.ErrUserNotFoundRes, zap.Error(err))
	}

	_, err = s.merchantQueryRepository.FindById(ctx, req.MerchantID)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	_, err = s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ORDER_BY_ID", span, &status, order_errors.ErrFailedFindOrderById, zap.Error(err))
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	shipping, err := s.shippingAddressQueryRepository.FindByOrder(ctx, req.OrderID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SHIPPING_ADDRESS_BY_ORDER", span, &status, shippingaddress_errors.ErrFailedFindShippingAddressByOrder, zap.Error(err))
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			return s.errorhandler.HandleInvalidOrderItem(err, method, "FAILED_INVALID_ORDER_ITEM", span, &status, zap.Error(err))
		}
		totalAmount += item.Price*item.Quantity + shipping.ShippingCost
	}

	ppn := totalAmount * 11 / 100 // ppn
	totalAmountWithTax := totalAmount + ppn + shipping.ShippingCost

	span.SetAttributes(
		attribute.Int("calculated.amount", totalAmountWithTax),
	)

	s.logger.Debug("Calculated amount", zap.Int("amount", totalAmountWithTax))

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		paymentStatus = "failed"
		return s.errorhandler.HandleInsufficientBalance(err, method, "FAILED_INSUFFICIENT_BALANCE", span, &status, zap.Error(err))
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.CreateTransaction(ctx, req)
	if err != nil {
		return s.errorhandler.HandleCreateTransactionError(err, method, "FAILED_CREATE_TRANSACTION", span, &status, zap.Error(err))
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
		return errorhandler.HandleErrorJSONMarshal[*response.TransactionResponse](s.logger, err, method, "FAILED_CREATE_TRANSACTION", span, &status, transaction_errors.ErrFailedSendEmail, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-transaction-create", strconv.Itoa(transaction.ID), payloadBytes)
	if err != nil {
		return errorhandler.HandleErrorKafkaSend[*response.TransactionResponse](s.logger, err, method, "FAILED_CREATE_TRANSACTION", span, &status, transaction_errors.ErrFailedSendEmail, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	logSuccess("Successfully created transaction", zap.Int("transaction.id", transaction.ID), zap.Bool("success", true))

	return so, nil
}

func (s *transactionCommandService) UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("transaction.id", *req.TransactionID), attribute.Int("merchant.id", req.MerchantID), attribute.Int("order.id", req.OrderID))

	defer func() {
		end(status)
	}()

	existingTx, err := s.transactionQueryRepository.FindById(ctx, *req.TransactionID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_TRANSACTION_BY_ID", span, &status, transaction_errors.ErrFailedFindTransactionById, zap.Error(err))
	}

	if existingTx.PaymentStatus == "success" || existingTx.PaymentStatus == "refunded" {
		return s.errorhandler.HandleCannotModifiedStatus(err, method, "FAILED_PAYMENT_STATUS_CANNOT_BE_MODIFIED", span, &status, zap.Error(err))
	}

	_, err = s.merchantQueryRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	_, err = s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ORDER_BY_ID", span, &status, order_errors.ErrFailedFindOrderById, zap.Error(err))
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	shipping, err := s.shippingAddressQueryRepository.FindByOrder(ctx, req.OrderID)
	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_SHIPPING_ADDRESS_BY_ORDER", span, &status, shippingaddress_errors.ErrFailedFindShippingAddressByOrder, zap.Error(err))
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			return s.errorhandler.HandleInsufficientBalance(err, method, "FAILED_UPDATE_TRANSACTION", span, &status, zap.Error(err))
		}
		totalAmount += item.Price*item.Quantity + shipping.ShippingCost
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn + shipping.ShippingCost

	span.SetAttributes(attribute.Int("calculated.amount", totalAmountWithTax))
	s.logger.Debug("Calculated amount", zap.Int("amount", totalAmountWithTax))

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		paymentStatus = "failed"
		return s.errorhandler.HandleInvalidOrderItem(err, method, "FAILED_UPDATE_TRANSACTION", span, &status, zap.Error(err))
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.UpdateTransaction(ctx, req)
	if err != nil {
		return s.errorhandler.HandleUpdateTransactionError(err, method, "FAILED_UPDATE_TRANSACTION", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponse(transaction)

	s.mencache.DeleteTransactionCache(ctx, *req.TransactionID)

	logSuccess("Successfully updated transaction", zap.Int("transaction.id", transaction.ID), zap.Bool("success", true))

	return so, nil
}

func (s *transactionCommandService) TrashedTransaction(ctx context.Context, transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.TrashTransaction(ctx, transactionID)

	if err != nil {
		return s.errorhandler.HandleTrashedTransactionError(err, method, "FAILED_TRASH_TRANSACTION", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponseDeleteAt(res)

	s.mencache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully trashed transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", true))

	return so, nil
}

func (s *transactionCommandService) RestoreTransaction(ctx context.Context, transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.RestoreTransaction(ctx, transactionID)

	if err != nil {
		return s.errorhandler.HandleRestoreTransactionError(err, method, "FAILED_RESTORE_TRANSACTION", span, &status, zap.Error(err))
	}

	so := s.mapping.ToTransactionResponseDeleteAt(res)

	logSuccess("Successfully restored transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", true))

	return so, nil
}

func (s *transactionCommandService) DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, *response.ErrorResponse) {
	const method = "DeleteTransactionPermanently"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.DeleteTransactionPermanently(ctx, transactionID)

	if err != nil {
		return s.errorhandler.HandleDeleteTransactionPermanentError(err, method, "FAILED_DELETE_TRANSACTION_PERMANENTLY", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", success))

	return success, nil
}

func (s *transactionCommandService) RestoreAllTransactions(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllTransactions"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.RestoreAllTransactions(ctx)

	if err != nil {
		return s.errorhandler.HandleRestoreAllTransactionError(
			err, method, "FAILED_RESTORE_ALL_TRANSACTIONS", span, &status, zap.Error(err),
		)
	}

	logSuccess("All trashed transactions restored successfully", zap.Bool("success", success))

	return success, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.DeleteAllTransactionPermanent(ctx)

	if err != nil {
		return s.errorhandler.HandleDeleteAllTransactionPermanentError(err, method, "FAILED_DELETE_ALL_TRANSACTION_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully permanently deleted all trashed transactions", zap.Bool("success", success))

	return success, nil
}

func (s *transactionCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *transactionCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
