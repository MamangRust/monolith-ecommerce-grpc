package reviewdetailapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type ReviewDetailBaseResponseMapper interface {
	ToResponseReviewDetail(reviewDetail *pb.ReviewDetailsResponse) *response.ReviewDetailsResponse
	ToResponsesReviewDetail(ReviewDetails []*pb.ReviewDetailsResponse) []*response.ReviewDetailsResponse
}

type ReviewDetailQueryResponseMapper interface {
	ReviewDetailBaseResponseMapper
	ToApiResponseReviewDetail(pbResponse *pb.ApiResponseReviewDetail) *response.ApiResponseReviewDetail
	ToApiResponsesReviewDetail(pbResponse *pb.ApiResponsesReviewDetails) *response.ApiResponsesReviewDetails
	ToApiResponsePaginationReviewDetail(pbResponse *pb.ApiResponsePaginationReviewDetails) *response.ApiResponsePaginationReviewDetails
	ToApiResponsePaginationReviewDetailDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDetailsDeleteAt) *response.ApiResponsePaginationReviewDetailsDeleteAt
}

type ReviewDetailCommandResponseMapper interface {
	ReviewDetailBaseResponseMapper
	ToResponseReviewDetailDeleteAt(reviewDetail *pb.ReviewDetailsResponseDeleteAt) *response.ReviewDetailsResponseDeleteAt
	ToResponsesReviewDetailDeleteAt(ReviewDetails []*pb.ReviewDetailsResponseDeleteAt) []*response.ReviewDetailsResponseDeleteAt
	ToApiResponseReviewDetailDeleteAt(pbResponse *pb.ApiResponseReviewDetailDeleteAt) *response.ApiResponseReviewDetailDeleteAt
	ToApiResponsePaginationReviewDetailDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDetailsDeleteAt) *response.ApiResponsePaginationReviewDetailsDeleteAt
}
