package reviewapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewBaseResponseMapper interface {
	ToResponseReview(pbResponse *pb.ReviewResponse) *response.ReviewResponse
	ToResponsesReview(pbResponses []*pb.ReviewResponse) []*response.ReviewResponse
	ToResponseReviewsDetail(pbResponse *pb.ReviewsDetailResponse) *response.ReviewsDetailResponse
	ToResponsesReviewsDetail(pbResponses []*pb.ReviewsDetailResponse) []*response.ReviewsDetailResponse
}

type ReviewQueryResponseMapper interface {
	ReviewBaseResponseMapper
	ToApiResponseReview(pbResponse *pb.ApiResponseReview) *response.ApiResponseReview
	ToApiResponsesReview(pbResponse *pb.ApiResponsesReview) *response.ApiResponsesReview
	ToApiResponsePaginationReview(pbResponse *pb.ApiResponsePaginationReview) *response.ApiResponsePaginationReview
	ToApiResponsePaginationReviewsDetail(pbResponse *pb.ApiResponsePaginationReviewDetail) *response.ApiResponsePaginationReviewsDetail
	ToApiResponsePaginationReviewDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt
}

type ReviewCommandResponseMapper interface {
	ReviewBaseResponseMapper
	ToResponseReviewDeleteAt(pbResponse *pb.ReviewResponseDeleteAt) *response.ReviewResponseDeleteAt
	ToResponsesReviewDeleteAt(pbResponses []*pb.ReviewResponseDeleteAt) []*response.ReviewResponseDeleteAt
	ToApiResponseReviewDeleteAt(pbResponse *pb.ApiResponseReviewDeleteAt) *response.ApiResponseReviewDeleteAt
	ToApiResponseReviewDelete(pbResponse *pb.ApiResponseReviewDelete) *response.ApiResponseReviewDelete
	ToApiResponseReviewAll(pbResponse *pb.ApiResponseReviewAll) *response.ApiResponseReviewAll
	ToApiResponsePaginationReviewDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt
}
