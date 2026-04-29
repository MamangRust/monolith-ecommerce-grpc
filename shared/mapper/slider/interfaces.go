package sliderapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type SliderBaseResponseMapper interface {
	ToResponseSlider(pbResponse *pb.SliderResponse) *response.SliderResponse
	ToResponsesSlider(pbResponses []*pb.SliderResponse) []*response.SliderResponse
}

type SliderQueryResponseMapper interface {
	SliderBaseResponseMapper
	ToApiResponseSlider(pbResponse *pb.ApiResponseSlider) *response.ApiResponseSlider
	ToApiResponsesSlider(pbResponse *pb.ApiResponsesSlider) *response.ApiResponsesSlider
	ToApiResponsePaginationSlider(pbResponse *pb.ApiResponsePaginationSlider) *response.ApiResponsePaginationSlider
	ToApiResponsePaginationSliderDeleteAt(pbResponse *pb.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt
}

type SliderCommandResponseMapper interface {
	SliderBaseResponseMapper
	ToResponseSliderDeleteAt(pbResponse *pb.SliderResponseDeleteAt) *response.SliderResponseDeleteAt
	ToResponsesSliderDeleteAt(pbResponses []*pb.SliderResponseDeleteAt) []*response.SliderResponseDeleteAt
	ToApiResponseSliderDeleteAt(pbResponse *pb.ApiResponseSliderDeleteAt) *response.ApiResponseSliderDeleteAt
	ToApiResponseSliderDelete(pbResponse *pb.ApiResponseSliderDelete) *response.ApiResponseSliderDelete
	ToApiResponseSliderAll(pbResponse *pb.ApiResponseSliderAll) *response.ApiResponseSliderAll
	ToApiResponsePaginationSliderDeleteAt(pbResponse *pb.ApiResponsePaginationSliderDeleteAt) *response.ApiResponsePaginationSliderDeleteAt
}
