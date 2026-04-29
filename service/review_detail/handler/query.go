package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type reviewDetailQueryHandler struct {
	pb.UnimplementedReviewDetailQueryServiceServer
	service service.ReviewDetailQueryService
	logger  logger.LoggerInterface
}

func NewReviewDetailQueryHandler(service service.ReviewDetailQueryService, logger logger.LoggerInterface) pb.ReviewDetailQueryServiceServer {
	return &reviewDetailQueryHandler{
		service: service,
		logger:  logger,
	}
}

func (s *reviewDetailQueryHandler) FindAll(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetails, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviewDetails, totalRecords, err := s.service.FindAll(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponse, len(reviewDetails))
	for i, reviewDetail := range reviewDetails {
		protoReviewDetails[i] = &pb.ReviewDetailsResponse{
			Id:        int32(reviewDetail.ReviewDetailID),
			ReviewId:  int32(reviewDetail.ReviewID),
			Type:      reviewDetail.Type,
			Url:       reviewDetail.Url,
			Caption:   *reviewDetail.Caption,
			CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDetails{
		Status:     "success",
		Message:    "Successfully fetched review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewDetailQueryHandler) FindById(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	reviewDetail, err := s.service.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetail := &pb.ReviewDetailsResponse{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
	}

	return &pb.ApiResponseReviewDetail{
		Status:  "success",
		Message: "Successfully fetched review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailQueryHandler) FindByActive(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviewDetails, totalRecords, err := s.service.FindActive(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponseDeleteAt, len(reviewDetails))
	for i, reviewDetail := range reviewDetails {
		var deletedAt string
		if reviewDetail.DeletedAt.Valid {
			deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
		}

		protoReviewDetails[i] = &pb.ReviewDetailsResponseDeleteAt{
			Id:        int32(reviewDetail.ReviewDetailID),
			ReviewId:  int32(reviewDetail.ReviewID),
			Type:      reviewDetail.Type,
			Url:       reviewDetail.Url,
			Caption:   *reviewDetail.Caption,
			CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}

func (s *reviewDetailQueryHandler) FindByTrashed(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error) {
	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllReview{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	reviewDetails, totalRecords, err := s.service.FindTrashed(ctx, &reqService)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	protoReviewDetails := make([]*pb.ReviewDetailsResponseDeleteAt, len(reviewDetails))
	for i, reviewDetail := range reviewDetails {
		var deletedAt string
		if reviewDetail.DeletedAt.Valid {
			deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
		}

		protoReviewDetails[i] = &pb.ReviewDetailsResponseDeleteAt{
			Id:        int32(reviewDetail.ReviewDetailID),
			ReviewId:  int32(reviewDetail.ReviewID),
			Type:      reviewDetail.Type,
			Url:       reviewDetail.Url,
			Caption:   *reviewDetail.Caption,
			CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	return &pb.ApiResponsePaginationReviewDetailsDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed review details",
		Data:       protoReviewDetails,
		Pagination: paginationMeta,
	}, nil
}
