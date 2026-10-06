package reviewdetailapimapper

import (
	pbreview_detail "github.com/MamangRust/monolith-ecommerce-pb/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type reviewDetailQueryResponseMapper struct{}

func NewReviewDetailQueryResponseMapper() ReviewDetailQueryResponseMapper {
	return &reviewDetailQueryResponseMapper{}
}

func (m *reviewDetailQueryResponseMapper) ToResponseReviewDetail(reviewDetail *pbreview_detail.ReviewDetailsResponse) *response.ReviewDetailsResponse {
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

func (m *reviewDetailQueryResponseMapper) ToResponsesReviewDetail(ReviewDetails []*pbreview_detail.ReviewDetailsResponse) []*response.ReviewDetailsResponse {
	var mappedReviewDetails []*response.ReviewDetailsResponse
	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetail(ReviewDetail))
	}
	return mappedReviewDetails
}

func (m *reviewDetailQueryResponseMapper) ToApiResponseReviewDetail(pbResponse *pbreview_detail.ApiResponseReviewDetail) *response.ApiResponseReviewDetail {
	return &response.ApiResponseReviewDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponseReviewDetail(pbResponse.Data),
	}
}

func (m *reviewDetailQueryResponseMapper) ToApiResponsesReviewDetail(pbResponse *pbreview_detail.ApiResponsesReviewDetails) *response.ApiResponsesReviewDetails {
	return &response.ApiResponsesReviewDetails{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    m.ToResponsesReviewDetail(pbResponse.Data),
	}
}

func (m *reviewDetailQueryResponseMapper) ToApiResponsePaginationReviewDetail(pbResponse *pbreview_detail.ApiResponsePaginationReviewDetails) *response.ApiResponsePaginationReviewDetails {
	return &response.ApiResponsePaginationReviewDetails{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesReviewDetail(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *reviewDetailQueryResponseMapper) ToApiResponsePaginationReviewDetailDeleteAt(pbResponse *pbreview_detail.ApiResponsePaginationReviewDetailsDeleteAt) *response.ApiResponsePaginationReviewDetailsDeleteAt {
	return &response.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       m.ToResponsesReviewDetailDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (m *reviewDetailQueryResponseMapper) ToResponseReviewDetailDeleteAt(reviewDetail *pbreview_detail.ReviewDetailsResponseDeleteAt) *response.ReviewDetailsResponseDeleteAt {
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

func (m *reviewDetailQueryResponseMapper) ToResponsesReviewDetailDeleteAt(ReviewDetails []*pbreview_detail.ReviewDetailsResponseDeleteAt) []*response.ReviewDetailsResponseDeleteAt {
	var mappedReviewDetails []*response.ReviewDetailsResponseDeleteAt
	for _, ReviewDetail := range ReviewDetails {
		mappedReviewDetails = append(mappedReviewDetails, m.ToResponseReviewDetailDeleteAt(ReviewDetail))
	}
	return mappedReviewDetails
}
