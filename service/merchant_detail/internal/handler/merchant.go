package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantDetailHandleGrpc struct {
	pb.UnimplementedMerchantDetailServiceServer
	merchantDetailQueryService   service.MerchantDetailQueryService
	merchantDetailCommandService service.MerchantDetailCommandService
	mapping                      protomapper.MerchantDetailProtoMapper
	mappingMerchant              protomapper.MerchantProtoMapper
}

func NewMerchantDetailHandleGrpc(
	service *service.Service,
	mapping protomapper.MerchantDetailProtoMapper,
	mappingMerchant protomapper.MerchantProtoMapper,
) *merchantDetailHandleGrpc {
	return &merchantDetailHandleGrpc{
		merchantDetailQueryService:   service.MerchantDetailQuery,
		merchantDetailCommandService: service.MerchantDetailCommand,
		mapping:                      mapping,
		mappingMerchant:              mappingMerchant,
	}
}

func (s *merchantDetailHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetail, error) {
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

	merchant, totalRecords, err := s.merchantDetailQueryService.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDetail(paginationMeta, "success", "Successfully fetched merchant", merchant)
	return so, nil
}

func (s *merchantDetailHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	merchant, err := s.merchantDetailQueryService.FindById(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDetailRelation("success", "Successfully fetched merchant", merchant)

	return so, nil

}

func (s *merchantDetailHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
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

	merchant, totalRecords, err := s.merchantDetailQueryService.FindByActive(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDetailDeleteAt(paginationMeta, "success", "Successfully fetched active merchant", merchant)

	return so, nil
}

func (s *merchantDetailHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllMerchantRequest) (*pb.ApiResponsePaginationMerchantDetailDeleteAt, error) {
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

	users, totalRecords, err := s.merchantDetailQueryService.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationMerchantDetailDeleteAt(paginationMeta, "success", "Successfully fetched trashed merchant", users)

	return so, nil
}

func (s *merchantDetailHandleGrpc) Create(ctx context.Context, request *pb.CreateMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	socialLinks := make([]*requests.CreateMerchantSocialRequest, 0)
	for _, link := range request.GetSocialLinks() {
		socialLinks = append(socialLinks, &requests.CreateMerchantSocialRequest{
			Platform: link.GetPlatform(),
			Url:      link.GetUrl(),
		})
	}

	req := &requests.CreateMerchantDetailRequest{
		MerchantID:       int(request.GetMerchantId()),
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
		SocialLink:       socialLinks,
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateCreateMerchantDetail
	}

	merchant, err := s.merchantDetailCommandService.CreateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDetail("success", "Successfully created merchant Detail", merchant)
	return so, nil
}

func (s *merchantDetailHandleGrpc) Update(ctx context.Context, request *pb.UpdateMerchantDetailRequest) (*pb.ApiResponseMerchantDetail, error) {
	id := int(request.GetMerchantDetailId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	socialLinks := make([]*requests.UpdateMerchantSocialRequest, 0)
	for _, link := range request.GetSocialLinks() {
		socialLinks = append(socialLinks, &requests.UpdateMerchantSocialRequest{
			ID:               int(link.GetId()),
			Platform:         link.GetPlatform(),
			Url:              link.GetUrl(),
			MerchantDetailID: &id,
		})
	}

	req := &requests.UpdateMerchantDetailRequest{
		MerchantDetailID: &id,
		DisplayName:      request.GetDisplayName(),
		CoverImageUrl:    request.GetCoverImageUrl(),
		LogoUrl:          request.GetLogoUrl(),
		ShortDescription: request.GetShortDescription(),
		WebsiteUrl:       request.GetWebsiteUrl(),
		SocialLink:       socialLinks,
	}

	if err := req.Validate(); err != nil {
		return nil, merchantdetail_errors.ErrGrpcValidateUpdateMerchantDetail
	}

	merchant, err := s.merchantDetailCommandService.UpdateMerchant(req)
	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDetail("success", "Successfully updated merchant Detail", merchant)
	return so, nil
}

func (s *merchantDetailHandleGrpc) TrashedMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	merchant, err := s.merchantDetailCommandService.TrashedMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDetailDeleteAt("success", "Successfully trashed merchant", merchant)

	return so, nil
}

func (s *merchantDetailHandleGrpc) RestoreMerchant(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	merchant, err := s.merchantDetailCommandService.RestoreMerchant(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseMerchantDetailDeleteAt("success", "Successfully restored merchant", merchant)

	return so, nil
}

func (s *merchantDetailHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pb.FindByIdMerchantRequest) (*pb.ApiResponseMerchantDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, merchantdetail_errors.ErrGrpcInvalidMerchantDetailId
	}

	_, err := s.merchantDetailCommandService.DeleteMerchantPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantDelete("success", "Successfully deleted merchant permanently")

	return so, nil
}

func (s *merchantDetailHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantDetailCommandService.RestoreAllMerchant()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully restore all merchant")

	return so, nil
}

func (s *merchantDetailHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseMerchantAll, error) {
	_, err := s.merchantDetailCommandService.DeleteAllMerchantPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mappingMerchant.ToProtoResponseMerchantAll("success", "Successfully delete merchant permanen")

	return so, nil
}
