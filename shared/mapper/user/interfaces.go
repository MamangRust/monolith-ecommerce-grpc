package userapimapper

import (
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type UserBaseResponseMapper interface {
	// Converts a single user response into an API response.
	ToApiResponseUser(pbResponse *pb.ApiResponseUser) *response.ApiResponseUser
}

type UserQueryResponseMapper interface {
	UserBaseResponseMapper

	// Converts paginated user records into an API response.
	ToApiResponsePaginationUser(pbResponse *pb.ApiResponsePaginationUser) *response.ApiResponsePaginationUser

	// Converts paginated soft-deleted users into an API response.
	ToApiResponsePaginationUserDeleteAt(pbResponse *pb.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt
}

type UserCommandResponseMapper interface {
	UserBaseResponseMapper

	// Converts a soft-deleted user response into an API response.
	ToApiResponseUserDeleteAt(pbResponse *pb.ApiResponseUserDeleteAt) *response.ApiResponseUserDeleteAt

	// Converts a permanently deleted user response into an API response.
	ToApiResponseUserDelete(pbResponse *pb.ApiResponseUserDelete) *response.ApiResponseUserDelete

	// Converts all user records into an API response.
	ToApiResponseUserAll(pbResponse *pb.ApiResponseUserAll) *response.ApiResponseUserAll
}
