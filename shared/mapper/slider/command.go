package sliderapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type sliderCommandResponseMapper struct{}

func NewSliderCommandResponseMapper() SliderCommandResponseMapper {
	return &sliderCommandResponseMapper{}
}

func (s *sliderCommandResponseMapper) ToResponseSlider(pbResponse *pb.SliderResponse) *response.SliderResponse {
	return &response.SliderResponse{
		ID:        int(pbResponse.Id),
		Name:      pbResponse.Name,
		Image:     pbResponse.Image,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
	}
}

func (s *sliderCommandResponseMapper) ToResponsesSlider(pbResponses []*pb.SliderResponse) []*response.SliderResponse {
	var sliders []*response.SliderResponse
	for _, slider := range pbResponses {
		sliders = append(sliders, s.ToResponseSlider(slider))
	}
	return sliders
}

func (s *sliderCommandResponseMapper) ToResponseSliderDeleteAt(pbResponse *pb.SliderResponseDeleteAt) *response.SliderResponseDeleteAt {
	var deletedAt string
	if pbResponse.DeletedAt != nil {
		deletedAt = pbResponse.DeletedAt.Value
	}

	return &response.SliderResponseDeleteAt{
		ID:        int(pbResponse.Id),
		Name:      pbResponse.Name,
		Image:     pbResponse.Image,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (s *sliderCommandResponseMapper) ToResponsesSliderDeleteAt(pbResponses []*pb.SliderResponseDeleteAt) []*response.SliderResponseDeleteAt {
	var sliders []*response.SliderResponseDeleteAt
	for _, slider := range pbResponses {
		sliders = append(sliders, s.ToResponseSliderDeleteAt(slider))
	}
	return sliders
}

func (s *sliderCommandResponseMapper) ToApiResponseSliderDeleteAt(pbResponse *pb.ApiResponseSliderDeleteAt) *response.ApiResponseSliderDeleteAt {
	return &response.ApiResponseSliderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.ToResponseSliderDeleteAt(pbResponse.Data),
	}
}

func (s *sliderCommandResponseMapper) ToApiResponseSliderDelete(pbResponse *pb.ApiResponseSliderDelete) *response.ApiResponseSliderDelete {
	return &response.ApiResponseSliderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *sliderCommandResponseMapper) ToApiResponseSliderAll(pbResponse *pb.ApiResponseSliderAll) *response.ApiResponseSliderAll {
	return &response.ApiResponseSliderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *sliderCommandResponseMapper) ToApiResponsePaginationSliderDeleteAt(pbResponse *pb.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt {
	return &response.ApiResponsePaginationSliderDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.ToResponsesSliderDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
