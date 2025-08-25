package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/services"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantHandleGrpc struct {
	pb.UnimplementedMerchantServiceServer
	merchantQuery   services.MerchantQueryService
	merchantCommand services.MerchantCommandService
	logger          logger.LoggerInterface
	mapping         protomapper.MerchantProtoMapper
}

func NewMerchantHandleGrpc(service *services.Service, mapping protomapper.MerchantProtoMapper, logger logger.LoggerInterface) pb.MerchantServiceServer {
	return &merchantHandleGrpc{
		merchantQuery:   service.MerchantQuery,
		merchantCommand: service.MerchantCommand,
		logger:          logger,
		mapping:         mapping,
	}
}

func (s *merchantHandleGrpc) FindAll(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchant, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all merchants",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all merchants",
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

	s.logger.Info("Successfully fetched all merchants",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchant(paginationMeta, "success", "Successfully fetched merchant records", merchants)
	return so, nil
}

func (s *merchantHandleGrpc) FindByActive(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active merchants",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active merchants",
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

	s.logger.Info("Successfully fetched active merchants",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched active merchant records", merchants)
	return so, nil
}

func (s *merchantHandleGrpc) FindByTrashed(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDeleteAt, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed merchants",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed merchants",
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

	s.logger.Info("Successfully fetched trashed merchants",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant records", merchants)
	return so, nil
}

func (s *merchantHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant ID provided for lookup", zap.Int("merchant_id", id))
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Fetching merchant by ID", zap.Int("merchant_id", id))

	merchant, err := s.merchantQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch merchant by ID",
			zap.Int("merchant_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched merchant by ID", zap.Int("merchant_id", id))

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully fetched merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	s.logger.Info("Creating new merchant",
		zap.String("name", request.GetName()),
		zap.String("contact_email", request.GetContactEmail()),
	)

	req := &requests.CreateMerchantRequest{
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant creation",
			zap.String("name", request.GetName()),
			zap.String("contact_email", request.GetContactEmail()),
			zap.Error(err),
		)
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantCommand.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create merchant",
			zap.String("name", request.GetName()),
			zap.String("contact_email", request.GetContactEmail()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant created successfully",
		zap.Int("merchant_id", int(merchant.ID)),
		zap.String("name", merchant.Name),
		zap.Int("user_id", merchant.UserID),
	)

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully created merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetMerchantId())

	if id == 0 {
		s.logger.Error("Invalid merchant ID provided for update", zap.Int("merchant_id", id))
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Updating merchant", zap.Int("merchant_id", id))

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &id,
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant update",
			zap.Int("merchant_id", id),
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update merchant",
			zap.Int("merchant_id", id),
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant updated successfully",
		zap.Int("merchant_id", id),
		zap.String("name", merchant.Name),
	)

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully updated merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) UpdateStatus(ctx context.Context, req *pb.UpdateMerchantStatusRequest) (*pb.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		s.logger.Error("Invalid merchant ID provided for status update", zap.Int("merchant_id", id))
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Updating merchant status",
		zap.Int("merchant_id", id),
		zap.String("status", req.GetStatus()),
	)

	request := requests.UpdateMerchantStatusRequest{
		MerchantID: &id,
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant status update",
			zap.Int("merchant_id", id),
			zap.String("status", req.GetStatus()),
			zap.Error(err),
		)
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchantStatus
	}

	merchant, err := s.merchantCommand.UpdateMerchantStatus(ctx, &request)
	if err != nil {
		s.logger.Error("Failed to update merchant status",
			zap.Int("merchant_id", id),
			zap.String("status", req.GetStatus()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant status updated successfully",
		zap.Int("merchant_id", id),
		zap.String("status", merchant.Status),
	)

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully updated merchant status", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant ID for trashing", zap.Int("merchant_id", id))
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Moving merchant to trash", zap.Int("merchant_id", id))

	merchant, err := s.merchantCommand.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash merchant",
			zap.Int("merchant_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant moved to trash successfully",
		zap.Int("merchant_id", id),
		zap.String("name", merchant.Name),
	)

	so := s.mapping.ToProtoResponseMerchantDeleteAt("success", "Successfully trashed merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant ID for restore", zap.Int("merchant_id", id))
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Restoring merchant from trash", zap.Int("merchant_id", id))

	merchant, err := s.merchantCommand.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore merchant",
			zap.Int("merchant_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant restored successfully",
		zap.Int("merchant_id", id),
		zap.String("name", merchant.Name),
	)

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully restored merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant ID for permanent deletion", zap.Int("merchant_id", id))
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	s.logger.Info("Permanently deleting merchant", zap.Int("merchant_id", id))

	_, err := s.merchantCommand.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete merchant",
			zap.Int("merchant_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant permanently deleted", zap.Int("merchant_id", id))

	so := s.mapping.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")
	return so, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Restoring all trashed merchants")

	_, err := s.merchantCommand.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all merchants", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchants restored successfully")

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully restored all merchants")
	return so, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Permanently deleting all trashed merchants")

	_, err := s.merchantCommand.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all merchants", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchants permanently deleted")

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully deleted all merchants permanently")
	return so, nil
}
