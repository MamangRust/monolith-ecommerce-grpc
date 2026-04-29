package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/pb"
)

type ReviewDetailQueryHandler interface {
	pb.ReviewDetailQueryServiceServer
}

type ReviewDetailCommandHandler interface {
	pb.ReviewDetailCommandServiceServer
}

type ReviewDetailHandleGrpc interface {
	FindAll(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetails, error)
	FindById(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetail, error)
	FindByActive(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pb.FindAllReviewRequest) (*pb.ApiResponsePaginationReviewDetailsDeleteAt, error)
	Create(ctx context.Context, request *pb.CreateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error)
	Update(ctx context.Context, request *pb.UpdateReviewDetailRequest) (*pb.ApiResponseReviewDetail, error)
	TrashedReviewDetail(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetailDeleteAt, error)
	RestoreReviewDetail(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDetailDeleteAt, error)
	DeleteReviewDetailPermanent(ctx context.Context, request *pb.FindByIdReviewDetailRequest) (*pb.ApiResponseReviewDelete, error)
}
