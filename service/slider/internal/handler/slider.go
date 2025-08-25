package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-slider/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/slider_errors"
	protomapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/proto"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type sliderHandleGrpc struct {
	pb.UnimplementedSliderServiceServer
	sliderQuery   service.SliderQueryService
	sliderCommand service.SliderCommandService
	logger        logger.LoggerInterface
	mapping       protomapper.SliderProtoMapper
}

func NewSliderHandleGrpc(service *service.Service, logger logger.LoggerInterface) pb.SliderServiceServer {
	return &sliderHandleGrpc{
		sliderQuery:   service.SliderQuery,
		sliderCommand: service.SliderCommand,
		logger:        logger,
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

	s.logger.Info("Fetching all sliders",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch all sliders",
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

	s.logger.Info("Successfully fetched all sliders",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_sliders_count", len(sliders)),
	)

	so := s.mapping.ToProtoResponsePaginationSlider(paginationMeta, "success", "Successfully fetched sliders", sliders)
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

	s.logger.Info("Fetching active sliders",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch active sliders",
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

	s.logger.Info("Successfully fetched active sliders",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_sliders_count", len(sliders)),
	)

	so := s.mapping.ToProtoResponsePaginationSliderDeleteAt(paginationMeta, "success", "Successfully fetched active sliders", sliders)
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

	s.logger.Info("Fetching trashed sliders",
		zap.Int("page", page),
		zap.Int("page_size", pageSize),
		zap.String("search", search),
	)

	reqService := requests.FindAllSlider{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	sliders, totalRecords, err := s.sliderQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("Failed to fetch trashed sliders",
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

	s.logger.Info("Successfully fetched trashed sliders",
		zap.Int("page", page),
		zap.Int32("total_records", int32(*totalRecords)),
		zap.Int32("total_pages", int32(totalPages)),
		zap.Int("fetched_sliders_count", len(sliders)),
	)

	so := s.mapping.ToProtoResponsePaginationSliderDeleteAt(paginationMeta, "success", "Successfully fetched trashed sliders", sliders)
	return so, nil
}

func (s *sliderHandleGrpc) Create(ctx context.Context, request *pb.CreateSliderRequest) (*pb.ApiResponseSlider, error) {
	s.logger.Info("Creating new slider",
		zap.String("name", request.GetName()),
	)

	req := &requests.CreateSliderRequest{
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on slider creation",
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, slider_errors.ErrGrpcValidateCreateSlider
	}

	slider, err := s.sliderCommand.CreateSlider(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create slider",
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Slider created successfully",
		zap.Int("slider_id", int(slider.ID)),
		zap.String("name", slider.Name),
	)

	return s.mapping.ToProtoResponseSlider("success", "Successfully created slider", slider), nil
}

func (s *sliderHandleGrpc) Update(ctx context.Context, request *pb.UpdateSliderRequest) (*pb.ApiResponseSlider, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid slider ID provided for update", zap.Int("slider_id", id))
		return nil, slider_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Updating slider", zap.Int("slider_id", id))

	req := &requests.UpdateSliderRequest{
		ID:       &id,
		Nama:     request.GetName(),
		FilePath: request.GetImage(),
	}

	if err := req.Validate(); err != nil {
		s.logger.Error("Validation failed on slider update",
			zap.Int("slider_id", id),
			zap.String("name", request.GetName()),
			zap.Error(err),
		)
		return nil, slider_errors.ErrGrpcValidateUpdateSlider
	}

	slider, err := s.sliderCommand.UpdateSlider(ctx, req)
	if err != nil {
		s.logger.Error("Failed to update slider",
			zap.Int("slider_id", id),
			zap.String("name", request.GetName()),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Slider updated successfully",
		zap.Int("slider_id", id),
		zap.String("name", slider.Name),
	)

	return s.mapping.ToProtoResponseSlider("success", "Successfully updated slider", slider), nil
}

func (s *sliderHandleGrpc) TrashedSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid slider ID for trashing", zap.Int("slider_id", id))
		return nil, slider_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Moving slider to trash", zap.Int("slider_id", id))

	slider, err := s.sliderCommand.TrashedSlider(ctx, id)
	if err != nil {
		s.logger.Error("Failed to trash slider",
			zap.Int("slider_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Slider moved to trash successfully",
		zap.Int("slider_id", id),
		zap.String("name", slider.Name),
	)

	so := s.mapping.ToProtoResponseSliderDeleteAt("success", "Successfully trashed slider", slider)
	return so, nil
}

func (s *sliderHandleGrpc) RestoreSlider(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid slider ID for restore", zap.Int("slider_id", id))
		return nil, slider_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Restoring slider from trash", zap.Int("slider_id", id))

	slider, err := s.sliderCommand.RestoreSlider(ctx, id)
	if err != nil {
		s.logger.Error("Failed to restore slider",
			zap.Int("slider_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Slider restored successfully",
		zap.Int("slider_id", id),
		zap.String("name", slider.Name),
	)

	so := s.mapping.ToProtoResponseSliderDeleteAt("success", "Successfully restored slider", slider)
	return so, nil
}

func (s *sliderHandleGrpc) DeleteSliderPermanent(ctx context.Context, request *pb.FindByIdSliderRequest) (*pb.ApiResponseSliderDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		s.logger.Error("Invalid slider ID for permanent deletion", zap.Int("slider_id", id))
		return nil, slider_errors.ErrGrpcInvalidID
	}

	s.logger.Info("Permanently deleting slider", zap.Int("slider_id", id))

	_, err := s.sliderCommand.DeleteSliderPermanent(ctx, id)
	if err != nil {
		s.logger.Error("Failed to permanently delete slider",
			zap.Int("slider_id", id),
			zap.Any("error", err),
		)
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("Slider permanently deleted", zap.Int("slider_id", id))

	so := s.mapping.ToProtoResponseSliderDelete("success", "Successfully deleted slider permanently")
	return so, nil
}

func (s *sliderHandleGrpc) RestoreAllSlider(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	s.logger.Info("Restoring all trashed sliders")

	_, err := s.sliderCommand.RestoreAllSliders(ctx)
	if err != nil {
		s.logger.Error("Failed to restore all sliders", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All sliders restored successfully")

	so := s.mapping.ToProtoResponseSliderAll("success", "Successfully restored all sliders")
	return so, nil
}

func (s *sliderHandleGrpc) DeleteAllSliderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseSliderAll, error) {
	s.logger.Info("Permanently deleting all trashed sliders")

	_, err := s.sliderCommand.DeleteAllSlidersPermanent(ctx)
	if err != nil {
		s.logger.Error("Failed to permanently delete all sliders", zap.Any("error", err))
		return nil, response.ToGrpcErrorFromErrorResponse(err)
	}

	s.logger.Info("All sliders permanently deleted")

	so := s.mapping.ToProtoResponseSliderAll("success", "Successfully deleted all sliders permanently")
	return so, nil
}
