package sliderapimapper

import (
	pbslider "github.com/MamangRust/monolith-ecommerce-pb/slider"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type sliderCommandResponseMapper struct{}

func NewSliderCommandResponseMapper() SliderCommandResponseMapper {
	return &sliderCommandResponseMapper{}
}

func (s *sliderCommandResponseMapper) ToResponseSlider(pbResponse *pbslider.SliderResponse) *response.SliderResponse {
	return &response.SliderResponse{
		ID:        int(pbResponse.Id),
		Name:      pbResponse.Name,
		Image:     pbResponse.Image,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
	}
}

func (s *sliderCommandResponseMapper) ToResponsesSlider(pbResponses []*pbslider.SliderResponse) []*response.SliderResponse {
	var sliders []*response.SliderResponse
	for _, slider := range pbResponses {
		sliders = append(sliders, s.ToResponseSlider(slider))
	}
	return sliders
}

func (s *sliderCommandResponseMapper) ToResponseSliderDeleteAt(pbResponse *pbslider.SliderResponseDeleteAt) *response.SliderResponseDeleteAt {
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

func (s *sliderCommandResponseMapper) ToResponsesSliderDeleteAt(pbResponses []*pbslider.SliderResponseDeleteAt) []*response.SliderResponseDeleteAt {
	var sliders []*response.SliderResponseDeleteAt
	for _, slider := range pbResponses {
		sliders = append(sliders, s.ToResponseSliderDeleteAt(slider))
	}
	return sliders
}

func (s *sliderCommandResponseMapper) ToApiResponseSliderDeleteAt(pbResponse *pbslider.ApiResponseSliderDeleteAt) *response.ApiResponseSliderDeleteAt {
	return &response.ApiResponseSliderDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.ToResponseSliderDeleteAt(pbResponse.Data),
	}
}

func (s *sliderCommandResponseMapper) ToApiResponseSliderDelete(pbResponse *pbslider.ApiResponseSliderDelete) *response.ApiResponseSliderDelete {
	return &response.ApiResponseSliderDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *sliderCommandResponseMapper) ToApiResponseSliderAll(pbResponse *pbslider.ApiResponseSliderAll) *response.ApiResponseSliderAll {
	return &response.ApiResponseSliderAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *sliderCommandResponseMapper) ToApiResponsePaginationSliderDeleteAt(pbResponse *pbslider.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt {
	return &response.ApiResponsePaginationSliderDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.ToResponsesSliderDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
