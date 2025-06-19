package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/service"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sliderHandleGrpc struct {
	pb.UnimplementedSliderServiceServer
	sliderQuery   service.SliderQueryService
	sliderCommand service.SliderCommandService
	mapping       protomapper.SliderProtoMapper
}

func NewSliderHandleGrpc(service *service.Service) *sliderHandleGrpc {
	return &sliderHandleGrpc{
		sliderQuery:   service.SliderQuery,
		sliderCommand: service.SliderCommand,
		mapping:       protomapper.NewSliderProtoMapper(),
	}
}

func (s *sliderHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllSliderRequest) (*pb.ApiResponsePaginationSlider, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	category, totalRecords, err := s.sliderQuery.FindAll(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationSlider(paginationMeta, "success", "Successfully fetched slider", category)
	return so, nil
}

func (s *sliderHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllSliderRequest) (*pb.ApiResponsePaginationSliderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.sliderQuery.FindByActive(&reqService)

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
	so := s.mapping.ToProtoResponsePaginationSliderDeleteAt(paginationMeta, "success", "Successfully fetched active slider", users)

	return so, nil
}

func (s *sliderHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllSliderRequest) (*pb.ApiResponsePaginationSliderDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.sliderQuery.FindByTrashed(&reqService)

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

	so := s.mapping.ToProtoResponsePaginationSliderDeleteAt(paginationMeta, "success", "Successfully fetched trashed slider", users)

	return so, nil
}

func (s *sliderHandleGrpc) Create(ctx context.Context, request *pb.CreateSliderRequest) (*pb.ApiResponseSlider, error) {
	req := &requests.CreateSliderRequest{
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, slider_errors.ErrGrpcValidateCreateSlider
	}

	slider, err := s.sliderCommand.CreateSlider(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseSlider("success", "Successfully created slider", slider), nil
}

func (s *sliderHandleGrpc) Update(ctx context.Context, request *pb.UpdateSliderRequest) (*pb.ApiResponseSlider, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateSliderRequest{
		ID:       &id,
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		return nil, slider_errors.ErrGrpcValidateUpdateSlider
	}

	slider, err := s.sliderCommand.UpdateSlider(req)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	return s.mapping.ToProtoResponseSlider("success", "Successfully updated slider", slider), nil
}

func (s *sliderHandleGrpc) TrashedSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	slider, err := s.sliderCommand.TrashedSlider(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSliderDeleteAt("success", "Successfully trashed slider", slider)

	return so, nil
}

func (s *sliderHandleGrpc) RestoreSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	slider, err := s.sliderCommand.RestoreSlider(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSliderDeleteAt("success", "Successfully restored slider", slider)

	return so, nil
}

func (s *sliderHandleGrpc) DeleteSliderPermanent(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, slider_errors.ErrGrpcInvalidID
	}

	_, err := s.sliderCommand.DeleteSliderPermanent(id)

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSliderDelete("success", "Successfully deleted slider permanently")

	return so, nil
}

func (s *sliderHandleGrpc) RestoreAllSlider(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	_, err := s.sliderCommand.RestoreAllSliders()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSliderAll("success", "Successfully restored all sliders")

	return so, nil
}

func (s *sliderHandleGrpc) DeleteAllSliderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	_, err := s.sliderCommand.DeleteAllSlidersPermanent()

	if err != nil {
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	so := s.mapping.ToProtoResponseSliderAll("success", "Successfully deleted all sliders permanently")

	return so, nil
}
