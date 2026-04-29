package reviewapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type reviewQueryResponseMapper struct{}

func NewReviewQueryResponseMapper() ReviewQueryResponseMapper {
	return &reviewQueryResponseMapper{}
}

func (r *reviewQueryResponseMapper) ToResponseReview(pbResponse *pb.ReviewResponse) *response.ReviewResponse {
	return &response.ReviewResponse{
		ID:        int(pbResponse.Id),
		UserID:    int(pbResponse.UserId),
		ProductID: int(pbResponse.ProductId),
		Rating:    int(pbResponse.Rating),
		Comment:   pbResponse.Comment,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
	}
}

func (r *reviewQueryResponseMapper) ToResponsesReview(pbResponses []*pb.ReviewResponse) []*response.ReviewResponse {
	var reviews []*response.ReviewResponse
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReview(review))
	}
	return reviews
}

func (r *reviewQueryResponseMapper) ToResponseReviewsDetail(pbResponse *pb.ReviewsDetailResponse) *response.ReviewsDetailResponse {
	if pbResponse == nil {
		return nil
	}

	var deletedAt *string
	if pbResponse.DeletedAt != "" {
		deletedAt = &pbResponse.DeletedAt
	}

	var reviewDetail *response.ReviewDetailResponse
	if pbResponse.ReviewDetail != nil {
		reviewDetail = &response.ReviewDetailResponse{
			ID:        int(pbResponse.ReviewDetail.Id),
			Type:      pbResponse.ReviewDetail.Type,
			Url:       pbResponse.ReviewDetail.Url,
			Caption:   pbResponse.ReviewDetail.Caption,
			CreatedAt: pbResponse.ReviewDetail.CreatedAt,
		}
	}

	return &response.ReviewsDetailResponse{
		ID:           int(pbResponse.Id),
		UserID:       int(pbResponse.UserId),
		ProductID:    int(pbResponse.ProductId),
		Name:         pbResponse.Name,
		Comment:      pbResponse.Comment,
		Rating:       int(pbResponse.Rating),
		ReviewDetail: reviewDetail,
		CreatedAt:    pbResponse.CreatedAt,
		UpdatedAt:    pbResponse.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

func (r *reviewQueryResponseMapper) ToResponsesReviewsDetail(pbResponses []*pb.ReviewsDetailResponse) []*response.ReviewsDetailResponse {
	var reviews []*response.ReviewsDetailResponse
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReviewsDetail(review))
	}
	return reviews
}

func (r *reviewQueryResponseMapper) ToApiResponseReview(pbResponse *pb.ApiResponseReview) *response.ApiResponseReview {
	return &response.ApiResponseReview{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponseReview(pbResponse.Data),
	}
}

func (r *reviewQueryResponseMapper) ToApiResponsesReview(pbResponse *pb.ApiResponsesReview) *response.ApiResponsesReview {
	return &response.ApiResponsesReview{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponsesReview(pbResponse.Data),
	}
}

func (r *reviewQueryResponseMapper) ToApiResponsePaginationReview(pbResponse *pb.ApiResponsePaginationReview) *response.ApiResponsePaginationReview {
	return &response.ApiResponsePaginationReview{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       r.ToResponsesReview(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (r *reviewQueryResponseMapper) ToApiResponsePaginationReviewsDetail(pbResponse *pb.ApiResponsePaginationReviewDetail) *response.ApiResponsePaginationReviewsDetail {
	return &response.ApiResponsePaginationReviewsDetail{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponsesReviewsDetail(pbResponse.Data),
	}
}

func (r *reviewQueryResponseMapper) ToApiResponsePaginationReviewDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt {
	return &response.ApiResponsePaginationReviewDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       r.ToResponsesReviewDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (r *reviewQueryResponseMapper) ToResponseReviewDeleteAt(pbResponse *pb.ReviewResponseDeleteAt) *response.ReviewResponseDeleteAt {
	var deletedAt *string
	if pbResponse.DeletedAt != nil {
		val := pbResponse.DeletedAt.Value
		deletedAt = &val
	}

	return &response.ReviewResponseDeleteAt{
		ID:        int(pbResponse.Id),
		UserID:    int(pbResponse.UserId),
		ProductID: int(pbResponse.ProductId),
		Rating:    int(pbResponse.Rating),
		Comment:   pbResponse.Comment,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func (r *reviewQueryResponseMapper) ToResponsesReviewDeleteAt(pbResponses []*pb.ReviewResponseDeleteAt) []*response.ReviewResponseDeleteAt {
	var reviews []*response.ReviewResponseDeleteAt
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReviewDeleteAt(review))
	}
	return reviews
}
