package handler

import (
	"context"
	"log"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/transaction_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type transactionHandleGrpc struct {
	pb.UnimplementedTransactionServiceServer
	transactionQuery           service.TransactionQueryService
	transactionCommand         service.TransactionCommandService
	transactionStats           service.TransactionStatsService
	transactionStatsByMerchant service.TransactionStatsByMerchantService
	mapping                    protomapper.TransactionProtoMapper
}

func NewTransactionHandleGrpc(
	service service.Service,
) *transactionHandleGrpc {
	return &transactionHandleGrpc{
		transactionQuery:           service.TransactionQuery,
		transactionCommand:         service.TransactionCommand,
		transactionStats:           service.TransactionStats,
		transactionStatsByMerchant: service.TransactionStatsByMerchant,
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

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindAllTransactions(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "success", "Successfully fetched transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllTransactionMerchantRequest) (*pb.ApiResponsePaginationTransaction, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransactionByMerchant{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByMerchant(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationTransaction(paginationMeta, "success", "Successfully fetched transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionQuery.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

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

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByActive(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}
	so := s.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetched active transaction", transaction)

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

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByTrashed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	so := s.mapping.ToProtoResponsePaginationTransactionDeleteAt(paginationMeta, "success", "Successfully fetched trashed transaction", transaction)

	return so, nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccess(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountSuccess(&reqService)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthAmountSuccess("success", "Monthly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccess(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountSuccess(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearAmountSuccess("success", "Yearly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailed(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountFailed(&reqService)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthAmountFailed("success", "Monthly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusFailed(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	year := int(request.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountFailed(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearAmountFailed("success", "Yearly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccessByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountSuccessByMerchant(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthAmountSuccess("success", "Merchant monthly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccessByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountSuccessByMerchant(
		&reqService,
	)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearAmountSuccess("success", "Merchant yearly success data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailedByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountFailedByMerchant(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthAmountFailed("success", "Merchant monthly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindYearStatusFailedByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountFailedByMerchant(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearAmountFailed("success", "Merchant yearly failed data retrieved successfully", res), nil
}

func (s *transactionHandleGrpc) FindMonthMethodSuccess(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodSuccess(&requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthMethod("success", "Monthly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodSuccess(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodSuccess(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearMethod("success", "Yearly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantSuccess(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantSuccess(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthMethod("success", "Merchant monthly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodByMerchantSuccess(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantSuccess(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearMethod("success", "Merchant yearly payment methods retrieved successfully", methods), nil
}

func (s *transactionHandleGrpc) FindMonthMethodFailed(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodFailed(&requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthMethod("Failed", "Monthly payment methods retrieved Failedfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodFailed(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodFailed(year)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearMethod("Failed", "Yearly payment methods retrieved Failedfully", methods), nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantFailed(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantFailed(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseMonthMethod("Failed", "Merchant monthly payment methods retrieved Failedfully", methods), nil
}

func (s *transactionHandleGrpc) FindYearMethodByMerchantFailed(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantFailed(
		&reqService,
	)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseYearMethod("Failed", "Merchant yearly payment methods retrieved Failedfully", methods), nil
}

func (s *transactionHandleGrpc) Create(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	req := &requests.CreateTransactionRequest{
		UserID:        int(request.GetUserId()),
		OrderID:       int(request.GetOrderId()),
		MerchantID:    int(request.GetMerchantId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		log.Fatal(err)
		return nil, transaction_errors.ErrGrpcValidateCreateTransaction
	}

	transaction, err := s.transactionCommand.CreateTransaction(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransaction("success", "Successfully created transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) Update(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	id := int(request.GetTransactionId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateTransactionRequest{
		TransactionID: &id,
		OrderID:       int(request.GetOrderId()),
		MerchantID:    int(request.GetMerchantId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateUpdateTransaction
	}

	transaction, err := s.transactionCommand.UpdateTransaction(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransaction("success", "Successfully updated transaction", transaction)
	return so, nil
}

func (s *transactionHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.TrashedTransaction(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionDeleteAt("success", "Successfully trashed transaction", transaction)

	return so, nil
}

func (s *transactionHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.RestoreTransaction(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionDeleteAt("success", "Successfully restored transaction", transaction)

	return so, nil
}

func (s *transactionHandleGrpc) DeleteTransactionPermanent(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	_, err := s.transactionCommand.DeleteTransactionPermanently(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionDelete("success", "Successfully deleted Transaction permanently")

	return so, nil
}

func (s *transactionHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := s.transactionCommand.RestoreAllTransactions()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionAll("success", "Successfully restore all Transaction")

	return so, nil
}

func (s *transactionHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	_, err := s.transactionCommand.DeleteAllTransactionPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseTransactionAll("success", "Successfully delete Transaction permanen")

	return so, nil
}
