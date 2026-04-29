package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-review-detail/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	reviewdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/review_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type reviewDetailCommandHandler struct {
	pb.UnimplementedReviewDetailCommandServiceServer
	service service.ReviewDetailCommandService
	logger  logger.LoggerInterface
}

func NewReviewDetailCommandHandler(service service.ReviewDetailCommandService, logger logger.LoggerInterface) pb.ReviewDetailCommandServiceServer {
	return &reviewDetailCommandHandler{
		service: service,
		logger:  logger,
	}
}

func (s *reviewDetailCommandHandler) Create(ctx context.Context, request *pb.CreateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	req := &requests.CreateReviewDetailRequest{
		ReviewID: int(request.GetReviewId()),
		Type:     request.GetType(),
		Url:      request.GetUrl(),
		Caption:  request.GetCaption(),
	}

	if err := req.Validate(); err != nil {
		return nil, reviewdetail_errors.ErrGrpcValidateCreateReviewDetail
	}

	reviewDetail, err := s.service.Create(ctx, req)
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
		Message: "Successfully created review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailCommandHandler) Update(ctx context.Context, request *pb.UpdateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error) {
	id := int(request.GetReviewDetailId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateReviewDetailRequest{
		ReviewDetailID: &id,
		Type:           request.GetType(),
		Url:            request.GetUrl(),
		Caption:        request.GetCaption(),
	}

	if err := req.Validate(); err != nil {
		return nil, reviewdetail_errors.ErrGrpcValidateUpdateReviewDetail
	}

	reviewDetail, err := s.service.Update(ctx, req)
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
		Message: "Successfully updated review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailCommandHandler) TrashedReviewDetail(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	reviewDetail, err := s.service.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if reviewDetail.DeletedAt.Valid {
		deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
	}

	protoReviewDetail := &pb.ReviewDetailsResponseDeleteAt{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseReviewDetailDeleteAt{
		Status:  "success",
		Message: "Successfully trashed review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailCommandHandler) RestoreReviewDetail(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetailDeleteAt, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	reviewDetail, err := s.service.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	var deletedAt string
	if reviewDetail.DeletedAt.Valid {
		deletedAt = reviewDetail.DeletedAt.Time.Format("2006-01-02")
	}

	protoReviewDetail := &pb.ReviewDetailsResponseDeleteAt{
		Id:        int32(reviewDetail.ReviewDetailID),
		ReviewId:  int32(reviewDetail.ReviewID),
		Type:      reviewDetail.Type,
		Url:       reviewDetail.Url,
		Caption:   *reviewDetail.Caption,
		CreatedAt: reviewDetail.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt: reviewDetail.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
	}

	return &pb.ApiResponseReviewDetailDeleteAt{
		Status:  "success",
		Message: "Successfully restored review detail",
		Data:    protoReviewDetail,
	}, nil
}

func (s *reviewDetailCommandHandler) DeleteReviewDetailPermanent(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDelete, error) {
	id := int(request.GetId())

	if id == 0 {
		return nil, reviewdetail_errors.ErrGrpcInvalidID
	}

	_, err := s.service.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewDelete{
		Status:  "success",
		Message: "Successfully deleted review detail permanently",
	}, nil
}

func (s *reviewDetailCommandHandler) RestoreAllReviewDetail(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.service.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully restored all review details",
	}, nil
}

func (s *reviewDetailCommandHandler) DeleteAllReviewDetailPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseReviewAll, error) {
	_, err := s.service.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pb.ApiResponseReviewAll{
		Status:  "success",
		Message: "Successfully deleted all review details permanently",
	}, nil
}
