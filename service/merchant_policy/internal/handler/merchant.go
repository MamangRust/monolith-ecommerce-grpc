package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_policy/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantpolicy_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_policy_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantPolicyHandleGrpc struct {
	pb.UnimplementedMerchantPoliciesServiceServer
	merchantPolicyQueryService   service.MerchantPoliciesQueryService
	merchantPolicyCommandService service.MerchantPoliciesCommandService
	mapping                      protomapper.MerchantPolicyProtoMapper
	mappingMerchant              protomapper.MerchantProtoMapper
	logger                       logger.LoggerInterface
}

func NewMerchantPolicyHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.MerchantPoliciesServiceServer {
	return &merchantPolicyHandleGrpc{
		merchantPolicyQueryService:   service.MerchantPolicyQuery,
		merchantPolicyCommandService: service.MerchantPolicyCmd,
		mapping:                      protomapper.NewMerchantPolicyProtoMapper(),
		mappingMerchant:              protomapper.NewMerchantProtoMaper(),
		logger:                       logger,
	}
}

func (s *merchantPolicyHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantPolicies, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all merchant policies",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	policies, totalRecords, err := s.merchantPolicyQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all merchant policies",
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

	s.logger.Info("Successfully fetched all merchant policies",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantPolicy(paginationMeta, "success", "Successfully fetched merchant policies", policies)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant policy ID provided", zap.Int("policy_id", id))
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	s.logger.Info("Fetching merchant policy by ID", zap.Int("policy_id", id))

	policy, err := s.merchantPolicyQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch merchant policy by ID",
			zap.Int("policy_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched merchant policy by ID",
		zap.Int("policy_id", id),
		zap.Int("merchant_id", int(policy.MerchantID)),
		zap.String("policy_type", policy.PolicyType),
	)

	so := s.mapping.ToProtoResponseMerchantPolicy("success", "Successfully fetched merchant policy", policy)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active merchant policies",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	policies, totalRecords, err := s.merchantPolicyQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active merchant policies",
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

	s.logger.Info("Successfully fetched active merchant policies",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantPolicyDeleteAt(paginationMeta, "success", "Successfully fetched active merchant policies", policies)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantPoliciesDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed merchant policies",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	policies, totalRecords, err := s.merchantPolicyQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed merchant policies",
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

	s.logger.Info("Successfully fetched trashed merchant policies",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantPolicyDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant policies", policies)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	s.logger.Info("Creating merchant policy",
		zap.Int("merchant_id", int(request.GetMerchantId())),
		zap.String("policy_type", request.GetPolicyType()),
		zap.String("title", request.GetTitle()),
	)

	req := &requests.CreateMerchantPolicyRequest{
		MerchantID:  int(request.GetMerchantId()),
		PolicyType:  request.GetPolicyType(),
		Title:       request.GetTitle(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant policy creation",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("policy_type", request.GetPolicyType()),
			zap.String("title", request.GetTitle()),
			zap.Error(err),
		)
		return nil, merchantpolicy_errors.ErrGrpcValidateCreateMerchantPolicy
	}

	policy, err := s.merchantPolicyCommandService.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create merchant policy",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("policy_type", request.GetPolicyType()),
			zap.String("title", request.GetTitle()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant policy created successfully",
		zap.Int("policy_id", int(policy.ID)),
		zap.Int("merchant_id", int(policy.MerchantID)),
		zap.String("policy_type", policy.PolicyType),
		zap.String("title", policy.Title),
	)

	so := s.mapping.ToProtoResponseMerchantPolicy("success", "Successfully created merchant policy", policy)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantPoliciesRequest) (*pb.ApiResponseMerchantPolicies, error) {
	id := int(request.GetMerchantPolicyId())

	if id == 0 {
		s.logger.Error("Invalid merchant policy ID provided for update", zap.Int("policy_id", id))
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	s.logger.Info("Updating merchant policy", zap.Int("policy_id", id))

	req := &requests.UpdateMerchantPolicyRequest{
		MerchantPolicyID: &id,
		PolicyType:       request.GetPolicyType(),
		Title:            request.GetTitle(),
		Description:      request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant policy update",
			zap.Int("policy_id", id),
			zap.String("policy_type", request.GetPolicyType()),
			zap.String("title", request.GetTitle()),
			zap.Error(err),
		)
		return nil, merchantpolicy_errors.ErrGrpcValidateUpdateMerchantPolicy
	}

	policy, err := s.merchantPolicyCommandService.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update merchant policy",
			zap.Int("policy_id", id),
			zap.String("policy_type", request.GetPolicyType()),
			zap.String("title", request.GetTitle()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant policy updated successfully",
		zap.Int("policy_id", id),
		zap.String("title", policy.Title),
		zap.String("policy_type", policy.PolicyType),
	)

	so := s.mapping.ToProtoResponseMerchantPolicy("success", "Successfully updated merchant policy", policy)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant policy ID for trashing", zap.Int("policy_id", id))
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	s.logger.Info("Moving merchant policy to trash", zap.Int("policy_id", id))

	policy, err := s.merchantPolicyCommandService.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash merchant policy",
			zap.Int("policy_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant policy moved to trash successfully",
		zap.Int("policy_id", id),
		zap.Int("merchant_id", int(policy.MerchantID)),
		zap.String("title", policy.Title),
	)

	so := s.mapping.ToProtoResponseMerchantPolicyDeleteAt("success", "Successfully trashed merchant policy", policy)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantPoliciesDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant policy ID for restore", zap.Int("policy_id", id))
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	s.logger.Info("Restoring merchant policy from trash", zap.Int("policy_id", id))

	policy, err := s.merchantPolicyCommandService.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore merchant policy",
			zap.Int("policy_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant policy restored successfully",
		zap.Int("policy_id", id),
		zap.String("title", policy.Title),
	)

	so := s.mapping.ToProtoResponseMerchantPolicyDeleteAt("success", "Successfully restored merchant policy", policy)
	return so, nil
}

func (s *merchantPolicyHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant policy ID for permanent deletion", zap.Int("policy_id", id))
		return nil, merchantpolicy_errors.ErrGrpcInvalidMerchantPolicyID
	}

	s.logger.Info("Permanently deleting merchant policy", zap.Int("policy_id", id))

	_, err := s.merchantPolicyCommandService.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete merchant policy",
			zap.Int("policy_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant policy permanently deleted", zap.Int("policy_id", id))

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant policy permanently")
	return so, nil
}

func (s *merchantPolicyHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Restoring all trashed merchant policies")

	_, err := s.merchantPolicyCommandService.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all merchant policies", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant policies restored successfully")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restored all merchant policies")
	return so, nil
}

func (s *merchantPolicyHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Permanently deleting all trashed merchant policies")

	_, err := s.merchantPolicyCommandService.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all merchant policies", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant policies permanently deleted")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully deleted all merchant policies permanently")
	return so, nil
}
