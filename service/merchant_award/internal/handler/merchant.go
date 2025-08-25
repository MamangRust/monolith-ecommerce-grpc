package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantAwardHandleGrpc struct {
	pb.UnimplementedMerchantAwardServiceServer
	merchantAwardQueryService   service.MerchantAwardQueryService
	merchantAwardCommandService service.MerchantAwardCommandService
	mapping                     protomapper.MerchantAwardProtoMapper
	mappingMerchant             protomapper.MerchantProtoMapper
	logger                      logger.LoggerInterface
}

func NewMerchantAwardHandleGrpc(
	service *service.Service,
	mapping protomapper.MerchantAwardProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
	logger logger.LoggerInterface,
) pb.MerchantAwardServiceServer {
	return &merchantAwardHandleGrpc{
		merchantAwardQueryService:   service.MerchantAwardQuery,
		merchantAwardCommandService: service.MerchantAwardCommand,
		mapping:                     mapping,
		mappingMerchant:             mappingMerchant,
		logger:                      logger,
	}
}

func (s *merchantAwardHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAward, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching all merchant awards",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	awards, totalRecords, err := s.merchantAwardQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all merchant awards",
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

	s.logger.Info("Successfully fetched all merchant awards",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantAward(paginationMeta, "success", "Successfully fetched merchant awards", awards)
	return so, nil
}

func (s *merchantAwardHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid merchant award ID provided", zap.Int("award_id", id))
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	s.logger.Info("Fetching merchant award by ID", zap.Int("award_id", id))

	award, err := s.merchantAwardQueryService.FindById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to fetch merchant award by ID",
			zap.Int("award_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Successfully fetched merchant award by ID",
		zap.Int("award_id", id),
		zap.Int("merchant_id", int(award.MerchantID)),
		zap.String("award_name", award.MerchantName),
	)

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully fetched merchant award", award)
	return so, nil
}

func (s *merchantAwardHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAwardDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching active merchant awards",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	awards, totalRecords, err := s.merchantAwardQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active merchant awards",
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

	s.logger.Info("Successfully fetched active merchant awards",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantAwardDeleteAt(paginationMeta, "success", "Successfully fetched active merchant awards", awards)
	return so, nil
}

func (s *merchantAwardHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantAwardDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	s.logger.Info("Fetching trashed merchant awards",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	awards, totalRecords, err := s.merchantAwardQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed merchant awards",
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

	s.logger.Info("Successfully fetched trashed merchant awards",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
	)

	so := s.mapping.ToProtoResponsePaginationMerchantAwardDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant awards", awards)
	return so, nil
}

func (s *merchantAwardHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	s.logger.Info("Creating merchant award",
		zap.Int("merchant_id", int(request.GetMerchantId())),
		zap.String("title", request.GetTitle()),
		zap.String("issue_date", request.GetIssueDate()),
	)

	req := &requests.CreateMerchantCertificationOrAwardRequest{
		MerchantID:     int(request.GetMerchantId()),
		Title:          request.GetTitle(),
		Description:    request.GetDescription(),
		IssuedBy:       request.GetIssuedBy(),
		IssueDate:      request.GetIssueDate(),
		ExpiryDate:     request.GetExpiryDate(),
		CertificateUrl: request.GetCertificateUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant award creation",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("title", request.GetTitle()),
			zap.Error(err),
		)
		return nil, merchantaward_errors.ErrGrpcValidateCreateMerchantAward
	}

	award, err := s.merchantAwardCommandService.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create merchant award",
			zap.Int("merchant_id", int(request.GetMerchantId())),
			zap.String("title", request.GetTitle()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant award created successfully",
		zap.Int("award_id", int(award.ID)),
		zap.Int("merchant_id", int(award.MerchantID)),
		zap.String("title", award.Title),
	)

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully created merchant award", award)
	return so, nil
}

func (s *merchantAwardHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetMerchantCertificationId())

	if id == 0 {
		s.logger.Error("Invalid award ID provided for update", zap.Int("award_id", id))
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	s.logger.Info("Updating merchant award", zap.Int("award_id", id))

	req := &requests.UpdateMerchantCertificationOrAwardRequest{
		MerchantCertificationID: &id,
		Title:                   request.GetTitle(),
		Description:             request.GetDescription(),
		IssuedBy:                request.GetIssuedBy(),
		IssueDate:               request.GetIssueDate(),
		ExpiryDate:              request.GetExpiryDate(),
		CertificateUrl:          request.GetCertificateUrl(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on merchant award update",
			zap.Int("award_id", id),
			zap.String("title", request.GetTitle()),
			zap.Error(err),
		)
		return nil, merchantaward_errors.ErrGrpcValidateUpdateMerchantAward
	}

	award, err := s.merchantAwardCommandService.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update merchant award",
			zap.Int("award_id", id),
			zap.String("title", request.GetTitle()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant award updated successfully",
		zap.Int("award_id", id),
		zap.String("title", award.Title),
	)

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully updated merchant award", award)
	return so, nil
}

func (s *merchantAwardHandleGrpc) TrashedMerchantAward(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid award ID for trashing", zap.Int("award_id", id))
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	s.logger.Info("Moving merchant award to trash", zap.Int("award_id", id))

	award, err := s.merchantAwardCommandService.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash merchant award",
			zap.Int("award_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant award moved to trash successfully",
		zap.Int("award_id", id),
		zap.String("title", award.Title),
		zap.Int("merchant_id", int(award.MerchantID)),
	)

	so := s.mapping.ToProtoResponseMerchantAwardDeleteAt("success", "Successfully trashed merchant award", award)
	return so, nil
}

func (s *merchantAwardHandleGrpc) RestoreMerchantAward(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid award ID for restore", zap.Int("award_id", id))
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	s.logger.Info("Restoring merchant award from trash", zap.Int("award_id", id))

	award, err := s.merchantAwardCommandService.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore merchant award",
			zap.Int("award_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant award restored successfully",
		zap.Int("award_id", id),
		zap.String("title", award.Title),
	)

	so := s.mapping.ToProtoResponseMerchantAwardDeleteAt("success", "Successfully restored merchant award", award)
	return so, nil
}

func (s *merchantAwardHandleGrpc) DeleteMerchantAwardPermanent(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid award ID for permanent deletion", zap.Int("award_id", id))
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	s.logger.Info("Permanently deleting merchant award", zap.Int("award_id", id))

	_, err := s.merchantAwardCommandService.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete merchant award",
			zap.Int("award_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Merchant award permanently deleted", zap.Int("award_id", id))

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant award permanently")
	return so, nil
}

func (s *merchantAwardHandleGrpc) RestoreAllMerchantAward(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Restoring all trashed merchant awards")

	_, err := s.merchantAwardCommandService.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all merchant awards", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant awards restored successfully")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restored all merchant awards")
	return so, nil
}

func (s *merchantAwardHandleGrpc) DeleteAllMerchantAwardPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	s.logger.Info("Permanently deleting all trashed merchant awards")

	_, err := s.merchantAwardCommandService.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all merchant awards", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All merchant awards permanently deleted")

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully deleted all merchant awards permanently")
	return so, nil
}
