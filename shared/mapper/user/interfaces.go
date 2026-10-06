package userapimapper

import (
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type UserBaseResponseMapper interface {
	// Converts a single user response into an API response.
	ToApiResponseUser(pbResponse *pbuser.ApiResponseUser) *response.ApiResponseUser
}

type UserQueryResponseMapper interface {
	UserBaseResponseMapper

	// Converts paginated user records into an API response.
	ToApiResponsePaginationUser(pbResponse *pbuser.ApiResponsePaginationUser) *response.ApiResponsePaginationUser

	// Converts paginated soft-deleted users into an API response.
	ToApiResponsePaginationUserDeleteAt(pbResponse *pbuser.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt
}

type UserCommandResponseMapper interface {
	UserBaseResponseMapper

	// Converts a soft-deleted user response into an API response.
	ToApiResponseUserDeleteAt(pbResponse *pbuser.ApiResponseUserDeleteAt) *response.ApiResponseUserDeleteAt

	// Converts a permanently deleted user response into an API response.
	ToApiResponseUserDelete(pbResponse *pbuser.ApiResponseUserDelete) *response.ApiResponseUserDelete

	// Converts all user records into an API response.
	ToApiResponseUserAll(pbResponse *pbuser.ApiResponseUserAll) *response.ApiResponseUserAll
}
