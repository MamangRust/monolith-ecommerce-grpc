package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandleGrpc interface {
	pb.UserServiceServer

	FindAll(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUser, error)
	FindById(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUser, error)
	FindByActive(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error)
	FindByTrashed(ctx context.Context, request *pb.FindAllUserRequest) (*pb.ApiResponsePaginationUserDeleteAt, error)
	Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.ApiResponseUser, error)
	Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.ApiResponseUser, error)
	TrashedUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error)
	RestoreUser(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDeleteAt, error)
	DeleteUserPermanent(ctx context.Context, request *pb.FindByIdUserRequest) (*pb.ApiResponseUserDelete, error)

	RestoreAllUser(context.Context, *emptypb.Empty) (*pb.ApiResponseUserAll, error)
	DeleteAllUserPermanent(context.Context, *emptypb.Empty) (*pb.ApiResponseUserAll, error)
}
