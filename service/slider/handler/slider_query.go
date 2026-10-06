package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-slider/service"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbslider "github.com/MamangRust/monolith-ecommerce-pb/slider"
)

type sliderQueryHandler struct {
	pbslider.UnimplementedSliderQueryServiceServer
	sliderQuery service.SliderQueryService
	logger      logger.LoggerInterface
}

func NewSliderQueryHandler(sliderQuery service.SliderQueryService, logger logger.LoggerInterface) *sliderQueryHandler {
	return &sliderQueryHandler{
		sliderQuery: sliderQuery,
		logger:      logger,
	}
}

func (s *sliderQueryHandler) FindAll(ctx context.Context, request *pbslider.FindAllSliderRequest) (*pbslider.ApiResponsePaginationSlider, error) {
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

	sliders, totalRecords, err := s.sliderQuery.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pbslider.SliderResponse, len(sliders))
	for i, slider := range sliders {
		protoSliders[i] = MapToSliderResponseGetSlidersRow(slider)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pbslider.ApiResponsePaginationSlider{
		Status:     "success",
		Message:    "Successfully fetched slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
}

func (s *sliderQueryHandler) FindByActive(ctx context.Context, request *pbslider.FindAllSliderRequest) (*pbslider.ApiResponsePaginationSliderDeleteAt, error) {
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

	sliders, totalRecords, err := s.sliderQuery.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pbslider.SliderResponseDeleteAt, len(sliders))
	for i, slider := range sliders {
		protoSliders[i] = MapToSliderResponseDeleteAtGetSlidersActiveRow(slider)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pbslider.ApiResponsePaginationSliderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
}

func (s *sliderQueryHandler) FindByTrashed(ctx context.Context, request *pbslider.FindAllSliderRequest) (*pbslider.ApiResponsePaginationSliderDeleteAt, error) {
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

	sliders, totalRecords, err := s.sliderQuery.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoSliders := make([]*pbslider.SliderResponseDeleteAt, len(sliders))
	for i, slider := range sliders {
		protoSliders[i] = MapToSliderResponseDeleteAtGetSlidersTrashedRow(slider)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pbslider.ApiResponsePaginationSliderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed slider records",
		Data:       protoSliders,
		Pagination: paginationMeta,
	}, nil
}

func (s *sliderQueryHandler) FindById(ctx context.Context, request *pbslider.FindByIdSliderRequest) (*pbslider.ApiResponseSlider, error) {
	id := int(request.GetId())

	slider, err := s.sliderQuery.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbslider.ApiResponseSlider{
		Status:  "success",
		Message: "Successfully fetched slider by ID",
		Data:    MapToSliderResponseGetSliderByIDRow(slider),
	}, nil
}
