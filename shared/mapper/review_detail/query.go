package reviewdetailapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type reviewDetailQueryResponseMapper struct{}

func NewReviewDetailQueryResponseMapper() ReviewDetailQueryResponseMapper {
	return &reviewDetailQueryResponseMapper{}
}

func (m *reviewDetailQueryResponseMapper) ToResponseReviewDetail(reviewDetail *pb.ReviewDetailsResponse) *response.ReviewDetailsResponse {
	return &response.ReviewDetailsResponse{
		ID:        int(reviewDetail.Id),
		ReviewID:  int(reviewDetail.ReviewId),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
	}
}

func (m *reviewDetailQueryResponseMapper) ToResponsesReviewDetail(ReviewDetails []*pb.ReviewDetailsResponse) []*response.ReviewDetailsResponse {
	var mappedReviewDetails []*response.ReviewDetailsResponse
	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetail(ReviewDetail))
	}
	return mappedReviewDetails
}

func (m *reviewDetailQueryResponseMapper) ToApiResponseReviewDetail(pbResponse *pb.ApiResponseReviewDetail) *response.ApiResponseReviewDetail {
	return &response.ApiResponseReviewDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseReviewDetail(pbResponse.Data),
	}
}

func (m *reviewDetailQueryResponseMapper) ToApiResponsesReviewDetail(pbResponse *pb.ApiResponsesReviewDetails) *response.ApiResponsesReviewDetails {
	return &response.ApiResponsesReviewDetails{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesReviewDetail(pbResponse.Data),
	}
}

func (m *reviewDetailQueryResponseMapper) ToApiResponsePaginationReviewDetail(pbResponse *pb.ApiResponsePaginationReviewDetails) *response.ApiResponsePaginationReviewDetails {
	return &response.ApiResponsePaginationReviewDetails{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesReviewDetail(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *reviewDetailQueryResponseMapper) ToApiResponsePaginationReviewDetailDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDetailsDeleteAt) *response.ApiResponsePaginationReviewDetailsDeleteAt {
	return &response.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesReviewDetailDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *reviewDetailQueryResponseMapper) ToResponseReviewDetailDeleteAt(reviewDetail *pb.ReviewDetailsResponseDeleteAt) *response.ReviewDetailsResponseDeleteAt {
	var deletedAt *string
	if reviewDetail.DeletedAt != nil {
		val := reviewDetail.DeletedAt.Value
		deletedAt = &val
	}

	return &response.ReviewDetailsResponseDeleteAt{
		ID:        int(reviewDetail.Id),
		ReviewID:  int(reviewDetail.ReviewId),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt,
		UpdatedAt: reviewDetail.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (m *reviewDetailQueryResponseMapper) ToResponsesReviewDetailDeleteAt(ReviewDetails []*pb.ReviewDetailsResponseDeleteAt) []*response.ReviewDetailsResponseDeleteAt {
	var mappedReviewDetails []*response.ReviewDetailsResponseDeleteAt
	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetailDeleteAt(ReviewDetail))
	}
	return mappedReviewDetails
}
