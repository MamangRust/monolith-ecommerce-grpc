package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionHandleGrpc struct {
	pb.UnimplementedTransactionServiceServer
	transactionQuery           service.TransactionQueryService
	transactionCommand         service.TransactionCommandService
	transactionStats           service.TransactionStatsService
	transactionStatsByMerchant service.TransactionStatsByMerchantService
	logger                     logger.LoggerInterface
	mapping                    protomapper.TransactionProtoMapper
}

func NewTransactionHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.TransactionServiceServer {
	return &transactionHandleGrpc{
		transactionQuery:           service.TransactionQuery,
		transactionCommand:         service.TransactionCommand,
		transactionStats:           service.TransactionStats,
		transactionStatsByMerchant: service.TransactionStatsByMerchant,
		logger:                     logger,
		mapping:                    protomapper.NewTransactionProtoMapper(),
	}
}

func (s *transactionHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all transactions",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.transactionQuery.FindAllTransactions(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all transactions",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched all transactions",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_transactions_count", len(transactions)),
	)

	so := s.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "success", "Successfully fetched transactions", transactions)
	return so, nil
}

func (s *transactionHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllTransactionMerchantRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchantID := int(request.GetMerchantId())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllTransactionByMerchant{
		MerchantID: merchantID,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transactions, totalRecords, err := s.transactionQuery.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch transactions by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int("fetched_transactions_count", len(transactions)),
	)

	so := s.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "success", "Successfully fetched transactions", transactions)
	return so, nil
}

func (s *transactionHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid transaction ID provided", zap.Int("transaction_id", id))
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Fetching transaction by ID", zap.Int("transaction_id", id))

	transaction, err := s.transactionQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch transaction by ID",
			zap.Int("transaction_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched transaction by ID",
		zap.Int("transaction_id", id),
	)

	so := s.mapping.ToProtoResponseTransaction("success", "Successfully fetched transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active transactions",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.transactionQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active transactions",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched active transactions",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_transactions_count", len(transactions)),
	)

	so := s.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetched active transactions", transactions)
	return so, nil
}

func (s *transactionHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed transactions",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transactions, totalRecords, err := s.transactionQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed transactions",
			zap.Int("page", page),
			zap.Int("page_size", pageSize),
			zap.String("search", search),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("Successfully fetched trashed transactions",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_transactions_count", len(transactions)),
	)

	so := s.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetched trashed transactions", transactions)
	return so, nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccess(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		s.logger.Error("Invalid month", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	s.logger.Info("Fetching monthly success transactions",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly success transactions", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly success transactions",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	return s.mapping.ToProtoResponseMonthAmountSuccess("success", "Monthly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccess(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	year := int(request.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	s.logger.Info("Fetching yearly success transactions",
		zap.Int("year", year),
	)

	res, err := s.transactionStats.FindYearlyAmountSuccess(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly success transactions", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly success transactions",
		zap.Int("year", year),
	)

	return s.mapping.ToProtoResponseYearAmountSuccess("success", "Yearly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailed(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		s.logger.Error("Invalid month", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	s.logger.Info("Fetching monthly failed transactions",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	res, err := s.transactionStats.FindMonthlyAmountFailed(ctx, &reqService)

	if err != nil {
		s.logger.Error("Failed to fetch monthly failed transactions", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly failed transactions",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	return s.mapping.ToProtoResponseMonthAmountFailed("success", "Monthly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusFailed(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	year := int(request.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	s.logger.Info("Fetching yearly failed transactions",
		zap.Int("year", year),
	)

	res, err := s.transactionStats.FindYearlyAmountFailed(ctx, year)

	if err != nil {
		s.logger.Error("Failed to fetch yearly failed transactions", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly failed transactions",
		zap.Int("year", year),
	)

	return s.mapping.ToProtoResponseYearAmountFailed("success", "Yearly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccessByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())
	merchantID := int(request.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant monthly success", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for merchant monthly success", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly success", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Fetching monthly successful transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantID,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly successful transactions by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly successful transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	return s.mapping.ToProtoResponseMonthAmountSuccess("success", "Merchant monthly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccessByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	year := int(request.GetYear())
	merchantID := int(request.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant yearly success", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly success", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Fetching yearly successful transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly successful transactions by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly successful transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	return s.mapping.ToProtoResponseYearAmountSuccess("success", "Merchant yearly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailedByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())
	merchantID := int(request.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant monthly failed transactions", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for merchant monthly failed transactions", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly failed transactions", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Fetching monthly failed transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantID,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly failed transactions by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly failed transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	return s.mapping.ToProtoResponseMonthAmountFailed("success", "Merchant monthly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusFailedByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	year := int(request.GetYear())
	merchantID := int(request.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant yearly failed transactions", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly failed transactions", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Fetching yearly failed transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly failed transactions by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly failed transactions by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	return s.mapping.ToProtoResponseYearAmountFailed("success", "Merchant yearly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthMethodSuccess(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly payment method", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for monthly payment method", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	s.logger.Info("Fetching monthly successful payment methods",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	methods, err := s.transactionStats.FindMonthlyMethodSuccess(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("Failed to fetch monthly payment methods",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly payment methods",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseMonthMethod("success", "Monthly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodSuccess(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly payment method", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	s.logger.Info("Fetching yearly successful payment methods", zap.Int("year", year))

	methods, err := s.transactionStats.FindYearlyMethodSuccess(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly payment methods",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly payment methods",
		zap.Int("year", year),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseYearMethod("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantSuccess(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant monthly payment method", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly payment method", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for merchant monthly payment method", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	s.logger.Info("Fetching monthly payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantID,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly payment methods by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseMonthMethod("success", "Merchant monthly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodByMerchantSuccess(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant yearly payment method", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly payment method", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Fetching yearly payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly payment methods by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseYearMethod("success", "Merchant yearly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindMonthMethodFailed(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year provided for monthly failed payment methods", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for monthly failed payment methods", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	s.logger.Info("Fetching monthly failed payment methods",
		zap.Int("year", year),
		zap.Int("month", month),
	)

	methods, err := s.transactionStats.FindMonthlyMethodFailed(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("Failed to fetch monthly failed payment methods",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly failed payment methods",
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseMonthMethod("success", "Monthly failed payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodFailed(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		s.logger.Error("Invalid year provided for yearly failed payment methods", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	s.logger.Info("Fetching yearly failed payment methods", zap.Int("year", year))

	methods, err := s.transactionStats.FindYearlyMethodFailed(ctx, year)
	if err != nil {
		s.logger.Error("Failed to fetch yearly failed payment methods",
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly failed payment methods",
		zap.Int("year", year),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseYearMethod("success", "Yearly failed payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantFailed(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant monthly failed payment methods", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for monthly failed payment methods", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	if month <= 0 || month > 12 {
		s.logger.Error("Invalid month provided for merchant monthly failed payment methods", zap.Int("month", month))
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	s.logger.Info("Fetching monthly failed payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
	)

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantID,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch monthly failed payment methods by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched monthly failed payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("month", month),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseMonthMethod("success", "Merchant monthly failed payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodByMerchantFailed(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())
	merchantID := int(req.GetMerchantId())

	if year <= 0 {
		s.logger.Error("Invalid year provided for merchant yearly failed payment methods", zap.Int("year", year))
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if merchantID <= 0 {
		s.logger.Error("Invalid merchant ID provided for yearly failed payment methods", zap.Int("merchant_id", merchantID))
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Fetching yearly failed payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
	)

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: merchantID,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch yearly failed payment methods by merchant",
			zap.Int("merchant_id", merchantID),
			zap.Int("year", year),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched yearly failed payment methods by merchant",
		zap.Int("merchant_id", merchantID),
		zap.Int("year", year),
		zap.Int("methods_count", len(methods)),
	)

	return s.mapping.ToProtoResponseYearMethod("success", "Merchant yearly failed payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) Create(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	s.logger.Info("Creating new transaction",
		zap.Int("user_id", int(request.GetUserId())),
		zap.Int("order_id", int(request.GetOrderId())),
		zap.Int("merchant_id", int(request.GetMerchantId())),
		zap.String("payment_method", request.GetPaymentMethod()),
		zap.Int("amount", int(request.GetAmount())),
	)

	req := &requests.CreateTransactionRequest{
		UserID:        int(request.GetUserId()),
		OrderID:       int(request.GetOrderId()),
		MerchantID:    int(request.GetMerchantId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on transaction creation",
			zap.Int("user_id", int(request.GetUserId())),
			zap.Int("order_id", int(request.GetOrderId())),
			zap.Error(err),
		)
		return nil, transaction_errors.ErrGrpcValidateCreateTransaction
	}

	transaction, err := s.transactionCommand.CreateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create transaction",
			zap.Int("user_id", int(request.GetUserId())),
			zap.Int("order_id", int(request.GetOrderId())),
			zap.String("payment_method", request.GetPaymentMethod()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Transaction created successfully",
		zap.Int("transaction_id", int(transaction.ID)),
		zap.Int("amount", transaction.Amount),
	)

	so := s.mapping.ToProtoResponseTransaction("success", "Successfully created transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) Update(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		s.logger.Error("Invalid transaction ID provided for update", zap.Int("transaction_id", id))
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Updating transaction", zap.Int("transaction_id", id))

	req := &requests.UpdateTransactionRequest{
		TransactionID: &id,
		OrderID:       int(request.GetOrderId()),
		MerchantID:    int(request.GetMerchantId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on transaction update",
			zap.Int("transaction_id", id),
			zap.String("payment_method", request.GetPaymentMethod()),
			zap.Error(err),
		)
		return nil, transaction_errors.ErrGrpcValidateUpdateTransaction
	}

	transaction, err := s.transactionCommand.UpdateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update transaction",
			zap.Int("transaction_id", id),
			zap.String("payment_method", request.GetPaymentMethod()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Transaction updated successfully",
		zap.Int("transaction_id", id),
		zap.Int("amount", transaction.Amount),
	)

	so := s.mapping.ToProtoResponseTransaction("success", "Successfully updated transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid transaction ID for trashing", zap.Int("transaction_id", id))
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Moving transaction to trash", zap.Int("transaction_id", id))

	transaction, err := s.transactionCommand.TrashedTransaction(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash transaction",
			zap.Int("transaction_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Transaction moved to trash successfully",
		zap.Int("transaction_id", id),
		zap.Int("amount", transaction.Amount),
	)

	so := s.mapping.ToProtoResponseTransactionDeleteAt("success", "Successfully trashed transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid transaction ID for restore", zap.Int("transaction_id", id))
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Restoring transaction from trash", zap.Int("transaction_id", id))

	transaction, err := s.transactionCommand.RestoreTransaction(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore transaction",
			zap.Int("transaction_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Transaction restored successfully",
		zap.Int("transaction_id", id),
	)

	so := s.mapping.ToProtoResponseTransactionDeleteAt("success", "Successfully restored transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) DeleteTransactionPermanent(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid transaction ID for permanent deletion", zap.Int("transaction_id", id))
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Permanently deleting transaction", zap.Int("transaction_id", id))

	_, err := s.transactionCommand.DeleteTransactionPermanently(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete transaction",
			zap.Int("transaction_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Transaction permanently deleted", zap.Int("transaction_id", id))

	so := s.mapping.ToProtoResponseTransactionDelete("success", "Successfully deleted transaction permanently")
	return so, nil
}

func (s *transactionHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	s.logger.Info("Restoring all trashed transactions")

	_, err := s.transactionCommand.RestoreAllTransactions(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all transactions", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All transactions restored successfully")

	so := s.mapping.ToProtoResponseTransactionAll("success", "Successfully restored all transactions")
	return so, nil
}

func (s *transactionHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	s.logger.Info("Permanently deleting all trashed transactions")

	_, err := s.transactionCommand.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all transactions", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All transactions permanently deleted")

	so := s.mapping.ToProtoResponseTransactionAll("success", "Successfully deleted all transactions permanently")
	return so, nil
}
