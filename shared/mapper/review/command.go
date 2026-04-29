package reviewapimapper

import (
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type reviewCommandResponseMapper struct{}

func NewReviewCommandResponseMapper() ReviewCommandResponseMapper {
	return &reviewCommandResponseMapper{}
}

func (r *reviewCommandResponseMapper) ToResponseReview(pbResponse *pb.ReviewResponse) *response.ReviewResponse {
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

func (r *reviewCommandResponseMapper) ToResponsesReview(pbResponses []*pb.ReviewResponse) []*response.ReviewResponse {
	var reviews []*response.ReviewResponse
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReview(review))
	}
	return reviews
}

func (r *reviewCommandResponseMapper) ToResponseReviewsDetail(pbResponse *pb.ReviewsDetailResponse) *response.ReviewsDetailResponse {
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

func (r *reviewCommandResponseMapper) ToResponsesReviewsDetail(pbResponses []*pb.ReviewsDetailResponse) []*response.ReviewsDetailResponse {
	var reviews []*response.ReviewsDetailResponse
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReviewsDetail(review))
	}
	return reviews
}

func (r *reviewCommandResponseMapper) ToResponseReviewDeleteAt(pbResponse *pb.ReviewResponseDeleteAt) *response.ReviewResponseDeleteAt {
	var deletedAt string
	if pbResponse.DeletedAt != nil {
		deletedAt = pbResponse.DeletedAt.Value
	}

	return &response.ReviewResponseDeleteAt{
		ID:        int(pbResponse.Id),
		UserID:    int(pbResponse.UserId),
		ProductID: int(pbResponse.ProductId),
		Rating:    int(pbResponse.Rating),
		Comment:   pbResponse.Comment,
		CreatedAt: pbResponse.CreatedAt,
		UpdatedAt: pbResponse.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (r *reviewCommandResponseMapper) ToResponsesReviewDeleteAt(pbResponses []*pb.ReviewResponseDeleteAt) []*response.ReviewResponseDeleteAt {
	var reviews []*response.ReviewResponseDeleteAt
	for _, review := range pbResponses {
		reviews = append(reviews, r.ToResponseReviewDeleteAt(review))
	}
	return reviews
}

func (r *reviewCommandResponseMapper) ToApiResponseReviewDeleteAt(pbResponse *pb.ApiResponseReviewDeleteAt) *response.ApiResponseReviewDeleteAt {
	return &response.ApiResponseReviewDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    r.ToResponseReviewDeleteAt(pbResponse.Data),
	}
}

func (r *reviewCommandResponseMapper) ToApiResponseReviewDelete(pbResponse *pb.ApiResponseReviewDelete) *response.ApiResponseReviewDelete {
	return &response.ApiResponseReviewDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (r *reviewCommandResponseMapper) ToApiResponseReviewAll(pbResponse *pb.ApiResponseReviewAll) *response.ApiResponseReviewAll {
	return &response.ApiResponseReviewAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (r *reviewCommandResponseMapper) ToApiResponsePaginationReviewDeleteAt(pbResponse *pb.ApiResponsePaginationReviewDeleteAt) *response.ApiResponsePaginationReviewDeleteAt {
	return &response.ApiResponsePaginationReviewDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       r.ToResponsesReviewDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}
