package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant/internal/services"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchant_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantHandleGrpc struct {
	pb.UnimplementedMerchantServiceServer
	merchantQuery   services.MerchantQueryService
	merchantCommand services.MerchantCommandService

	mapping protomapper.MerchantProtoMapper
}

func NewMerchantHandleGrpc(service *services.Service, mapping protomapper.MerchantProtoMapper) *merchantHandleGrpc {
	return &merchantHandleGrpc{
		merchantQuery:   service.MerchantQuery,
		merchantCommand: service.MerchantCommand,
		mapping:         mapping,
	}
}

func (s *merchantHandleGrpc) FindAllMerchant(ctx context.Context, req *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchant, error) {
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

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

	merchants, totalRecords, err := s.merchantQuery.FindAll(ctx, &reqService)

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
	so := s.mapping.ToProtoResponsePaginationMerchant(paginationMeta, "success", "Successfully fetched merchant record", merchants)

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

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByActive(ctx, &reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched merchant record", res)
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

	reqService := requests.FindAllMerchant{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByTrashed(ctx, &reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDeleteAt(paginationMeta, "success", "Successfully fetched merchant record", res)

	return so, nil
}

func (s *merchantHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantQuery.FindById(ctx, id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully fetched merchant", merchant)

	return so, nil

}

func (s *merchantHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantRequest) (*pb.ApiResponseMerchant, error) {
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
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantCommand.CreateMerchant(ctx, req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully created merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

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
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.UpdateMerchant(ctx, req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully updated merchant", merchant)
	return so, nil
}

func (s *merchantHandleGrpc) UpdateStatus(ctx context.Context, req *pb.UpdateMerchantStatusRequest) (*pb.ApiResponseMerchant, error) {
	id := int(req.GetMerchantId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	request := requests.UpdateMerchantStatusRequest{
		MerchantID: &id,
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchantStatus
	}

	merchant, err := s.merchantCommand.UpdateMerchantStatus(ctx, &request)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully updated merchant status", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantCommand.TrashedMerchant(ctx, id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDeleteAt("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchant, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	merchant, err := s.merchantCommand.RestoreMerchant(ctx, id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchant("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchant_errors.ErrGrpcInvalidMerchantId
	}

	_, err := s.merchantCommand.DeleteMerchantPermanent(ctx, id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")

	return so, nil
}

func (s *merchantHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantCommand.RestoreAllMerchant(ctx)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantCommand.DeleteAllMerchantPermanent(ctx)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantAll("success", "Successfully delete merchant permanen")

	return so, nil
}
