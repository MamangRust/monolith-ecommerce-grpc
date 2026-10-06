package sliderapimapper

import (
	pbslider "github.com/MamangRust/monolith-ecommerce-pb/slider"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderBaseResponseMapper interface {
	ToResponseSlider(pbResponse *pbslider.SliderResponse) *response.SliderResponse
	ToResponsesSlider(pbResponses []*pbslider.SliderResponse) []*response.SliderResponse
}

type SliderQueryResponseMapper interface {
	SliderBaseResponseMapper
	ToApiResponseSlider(pbResponse *pbslider.ApiResponseSlider) *response.ApiResponseSlider
	ToApiResponsesSlider(pbResponse *pbslider.ApiResponsesSlider) *response.ApiResponsesSlider
	ToApiResponsePaginationSlider(pbResponse *pbslider.ApiResponsePaginationSlider) *response.ApiResponsePaginationSlider
	ToApiResponsePaginationSliderDeleteAt(pbResponse *pbslider.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt
}

type SliderCommandResponseMapper interface {
	SliderBaseResponseMapper
	ToResponseSliderDeleteAt(pbResponse *pbslider.SliderResponseDeleteAt) *response.SliderResponseDeleteAt
	ToResponsesSliderDeleteAt(pbResponses []*pbslider.SliderResponseDeleteAt) []*response.SliderResponseDeleteAt
	ToApiResponseSliderDeleteAt(pbResponse *pbslider.ApiResponseSliderDeleteAt) *response.ApiResponseSliderDeleteAt
	ToApiResponseSliderDelete(pbResponse *pbslider.ApiResponseSliderDelete) *response.ApiResponseSliderDelete
	ToApiResponseSliderAll(pbResponse *pbslider.ApiResponseSliderAll) *response.ApiResponseSliderAll
	ToApiResponsePaginationSliderDeleteAt(pbResponse *pbslider.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt
}
