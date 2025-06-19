package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantBusinessHandleGrpc struct {
	pb.UnimplementedMerchantBusinessServiceServer
	merchantBusinessQueryService   service.MerchantBusinessQueryService
	merchantBusinessCommandService service.MerchantBusinessCommandService
	mapping                        protomapper.MerchantBusinessProtoMapper
	mappingMerchant                protomapper.MerchantProtoMapper
}

func NewMerchantBusinessHandleGrpc(
	service *service.Service,
	mapping protomapper.MerchantBusinessProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
) *merchantBusinessHandleGrpc {
	return &merchantBusinessHandleGrpc{
		merchantBusinessQueryService:   service.MerchantBusinessQuery,
		merchantBusinessCommandService: service.MerchantBusinessCommand,
		mapping:                        mapping,
		mappingMerchant:                mappingMerchant,
	}
}

func (s *merchantBusinessHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusiness, error) {
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

	merchant, totalRecords, err := s.merchantBusinessQueryService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantBusiness(paginationMeta, "success", "Successfully fetched merchant", merchant)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessQueryService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantBusiness("success", "Successfully fetched merchant", merchant)

	return so, nil

}

func (s *merchantBusinessHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
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

	merchant, totalRecords, err := s.merchantBusinessQueryService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantBusinessDeleteAt(paginationMeta, "success", "Successfully fetched active merchant", merchant)

	return so, nil
}

func (s *merchantBusinessHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantBusinessDeleteAt, error) {
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

	users, totalRecords, err := s.merchantBusinessQueryService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantBusinessDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant", users)

	return so, nil
}

func (s *merchantBusinessHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	req := &requests.CreateMerchantBusinessInformationRequest{
		MerchantID:        int(request.GetMerchantId()),
		BusinessType:      request.GetBusinessType(),
		TaxID:             request.GetTaxId(),
		EstablishedYear:   int(request.GetEstablishedYear()),
		NumberOfEmployees: int(request.GetNumberOfEmployees()),
		WebsiteUrl:        request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantbusiness_errors.ErrGrpcValidateCreateMerchantBusiness
	}

	merchant, err := s.merchantBusinessCommandService.CreateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantBusiness("success", "Successfully created merchant business information", merchant)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantBusinessRequest) (*pb.ApiResponseMerchantBusiness, error) {
	id := int(request.GetMerchantBusinessInfoId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	req := &requests.UpdateMerchantBusinessInformationRequest{
		MerchantBusinessInfoID: &id,
		BusinessType:           request.GetBusinessType(),
		TaxID:                  request.GetTaxId(),
		EstablishedYear:        int(request.GetEstablishedYear()),
		NumberOfEmployees:      int(request.GetNumberOfEmployees()),
		WebsiteUrl:             request.GetWebsiteUrl(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchantbusiness_errors.ErrGrpcValidateUpdateMerchantBusiness
	}

	merchant, err := s.merchantBusinessCommandService.UpdateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantBusiness("success", "Successfully updated merchant business information", merchant)
	return so, nil
}

func (s *merchantBusinessHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessCommandService.TrashedMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantBusinessDeleteAt("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantBusinessHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantBusinessDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	merchant, err := s.merchantBusinessCommandService.RestoreMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantBusinessDeleteAt("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantBusinessHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantbusiness_errors.ErrGrpcInvalidMerchantBusinessId
	}

	_, err := s.merchantBusinessCommandService.DeleteMerchantPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")

	return so, nil
}

func (s *merchantBusinessHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantBusinessCommandService.RestoreAllMerchant()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantBusinessHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantBusinessCommandService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully delete merchant permanen")

	return so, nil
}
