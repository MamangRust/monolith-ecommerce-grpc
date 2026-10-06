package reviewapimapper

import (
	pbreview "github.com/MamangRust/monolith-ecommerce-pb/review"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewBaseResponseMapper interface {
	ToResponseReview(pbResponse *pbreview.ReviewResponse) *response.ReviewResponse
	ToResponsesReview(pbResponses []*pbreview.ReviewResponse) []*response.ReviewResponse
	ToResponseReviewsDetail(pbResponse *pbreview.ReviewsDetailResponse) *response.ReviewsDetailResponse
	ToResponsesReviewsDetail(pbResponses []*pbreview.ReviewsDetailResponse) []*response.ReviewsDetailResponse
}

type ReviewQueryResponseMapper interface {
	ReviewBaseResponseMapper
	ToApiResponseReview(pbResponse *pbreview.ApiResponseReview) *response.ApiResponseReview
	ToApiResponsesReview(pbResponse *pbreview.ApiResponsesReview) *response.ApiResponsesReview
	ToApiResponsePaginationReview(pbResponse *pbreview.ApiResponsePaginationReview) *response.ApiResponsePaginationReview
	ToApiResponsePaginationReviewsDetail(pbResponse *pbreview.ApiResponsePaginationReviewDetail) *response.ApiResponsePaginationReviewsDetail
	ToApiResponsePaginationReviewDeleteAt(pbResponse *pbreview.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt
}

type ReviewCommandResponseMapper interface {
	ReviewBaseResponseMapper
	ToResponseReviewDeleteAt(pbResponse *pbreview.ReviewResponseDeleteAt) *response.ReviewResponseDeleteAt
	ToResponsesReviewDeleteAt(pbResponses []*pbreview.ReviewResponseDeleteAt) []*response.ReviewResponseDeleteAt
	ToApiResponseReviewDeleteAt(pbResponse *pbreview.ApiResponseReviewDeleteAt) *response.ApiResponseReviewDeleteAt
	ToApiResponseReviewDelete(pbResponse *pbreview.ApiResponseReviewDelete) *response.ApiResponseReviewDelete
	ToApiResponseReviewAll(pbResponse *pbreview.ApiResponseReviewAll) *response.ApiResponseReviewAll
	ToApiResponsePaginationReviewDeleteAt(pbResponse *pbreview.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt
}
