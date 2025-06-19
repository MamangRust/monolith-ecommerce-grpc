package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_award/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantaward_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_award"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantAwardHandleGrpc struct {
	pb.UnimplementedMerchantAwardServiceServer
	merchantAwardQueryService   service.MerchantAwardQueryService
	merchantAwardCommandService service.MerchantAwardCommandService
	mapping                     protomapper.MerchantAwardProtoMapper
	mappingMerchant             protomapper.MerchantProtoMapper
}

func NewMerchantAwardHandleGrpc(
	service *service.Service,
	mapping protomapper.MerchantAwardProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
) *merchantAwardHandleGrpc {
	return &merchantAwardHandleGrpc{
		merchantAwardQueryService:   service.MerchantAwardQuery,
		merchantAwardCommandService: service.MerchantAwardCommand,
		mapping:                     mapping,
		mappingMerchant:             mappingMerchant,
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

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.merchantAwardQueryService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantAward(paginationMeta, "success", "Successfully fetched merchant", merchant)
	return so, nil
}

func (s *merchantAwardHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardQueryService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully fetched merchant", merchant)

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

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.merchantAwardQueryService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantAwardDeleteAt(paginationMeta, "success", "Successfully fetched active merchant", merchant)

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

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.merchantAwardQueryService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantAwardDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant", users)

	return so, nil
}

func (s *merchantAwardHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
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
		return nil, merchantaward_errors.ErrGrpcValidateCreateMerchantAward
	}

	merchant, err := s.merchantAwardCommandService.CreateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully created merchant award", merchant)
	return so, nil
}

func (s *merchantAwardHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantAwardRequest) (*pb.ApiResponseMerchantAward, error) {
	id := int(request.GetMerchantCertificationId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

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
		return nil, merchantaward_errors.ErrGrpcValidateUpdateMerchantAward
	}

	merchant, err := s.merchantAwardCommandService.UpdateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAward("success", "Successfully updated merchant award", merchant)
	return so, nil
}

func (s *merchantAwardHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardCommandService.TrashedMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAwardDeleteAt("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantAwardHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantAwardDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	merchant, err := s.merchantAwardCommandService.RestoreMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAwardDeleteAt("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantAwardHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantaward_errors.ErrGrpcMerchantInvalidId
	}

	_, err := s.merchantAwardCommandService.DeleteMerchantPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")

	return so, nil
}

func (s *merchantAwardHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardCommandService.RestoreAllMerchant()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantAwardHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantAwardCommandService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully delete merchant permanen")

	return so, nil
}
