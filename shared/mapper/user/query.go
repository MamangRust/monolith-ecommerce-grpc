package userapimapper

import (
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

// userQueryResponseMapper provides methods to map gRPC user responses to HTTP API responses.
type userQueryResponseMapper struct {
}

// NewUserResponseMapper returns a pointer to a userQueryResponseMapper.
//
// The userQueryResponseMapper is a helper struct for mapping user-related gRPC responses
// into HTTP/REST API response formats.
func NewUserQueryResponseMapper() *userQueryResponseMapper {
	return &userQueryResponseMapper{}
}

// ToApiResponseUser converts a single user response into an API response.
//
// Args:
//   - pbResponse: A pointer to a pb.ApiResponseUser containing the user data.
//
// Returns:
//   - A pointer to a response.ApiResponseUser containing the mapped user data, including
//     Status, Message, and Data.
func (u *userQueryResponseMapper) ToApiResponseUser(pbResponse *pb.ApiResponseUser) *response.ApiResponseUser {
	return &response.ApiResponseUser{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseUser(pbResponse.Data),
	}
}

// ToApiResponsePaginationUserDeleteAt maps a pagination meta, status, message, and a list of UserResponseDeleteAt
//
// Args:
//   - pbResponse: A pointer to a pb.ApiResponsePaginationUserDeleteAt containing the user data.
//
// Returns:
//   - A pointer to a response.ApiResponsePaginationUserDeleteAt containing the mapped data.
func (u *userQueryResponseMapper) ToApiResponsePaginationUserDeleteAt(pbResponse *pb.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt {
	return &response.ApiResponsePaginationUserDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       u.ToResponsesUserDeleteAt(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

// ToApiResponsePaginationUser maps a pagination meta, status, message, and a list of UserResponse
//
// Args:
//   - pbResponse: A pointer to a pb.ApiResponsePaginationUser containing the user data.
//
// Returns:
//   - A pointer to a response.ApiResponsePaginationUser containing the mapped data.
func (u *userQueryResponseMapper) ToApiResponsePaginationUser(pbResponse *pb.ApiResponsePaginationUser) *response.ApiResponsePaginationUser {
	return &response.ApiResponsePaginationUser{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       u.ToResponsesUser(pbResponse.Data),
		Pagination: *paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

// ToApiResponseUserDeleteAt maps a soft-deleted user response into an API response.
//
// Args:
//   - pbResponse: A pointer to a pb.ApiResponseUserDeleteAt containing the user data.
//
// Returns:
//   - A pointer to a response.ApiResponseUserDeleteAt containing the mapped user data, including
//     Status, Message, and Data.
func (u *userQueryResponseMapper) ToApiResponseUserDeleteAt(pbResponse *pb.ApiResponseUserDeleteAt) *response.ApiResponseUserDeleteAt {
	return &response.ApiResponseUserDeleteAt{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    u.ToResponseUserDeleteAt(pbResponse.Data),
	}
}

// ToResponseUser maps a protobuf UserResponse to a domain UserResponse.
//
// Args:
//   - user: A pointer to a pb.UserResponse containing the user data.
//
// Returns:
//   - A pointer to a response.UserResponse containing the mapped user data, including
//     ID, FirstName, LastName, Email, CreatedAt, and UpdatedAt.
func (u *userQueryResponseMapper) ToResponseUser(user *pb.UserResponse) *response.UserResponse {
	return &response.UserResponse{
		ID:        int(user.Id),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToResponsesUser maps a slice of pb.UserResponse to a slice of response.UserResponse.
//
// Args:
//   - users: A slice of pointers to pb.UserResponse to be mapped.
//
// Returns:
//   - A slice of pointers to response.UserResponse containing the mapped user data for each user, including
//     ID, FirstName, LastName, Email, CreatedAt, and UpdatedAt.
func (u *userQueryResponseMapper) ToResponsesUser(users []*pb.UserResponse) []*response.UserResponse {
	var mappedUsers []*response.UserResponse

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.ToResponseUser(user))
	}

	return mappedUsers
}

// ToResponseUserDeleteAt maps a protobuf UserResponseDeleteAt to a domain UserResponseDeleteAt.
//
// Args:
//   - user: A pointer to a pb.UserResponseDeleteAt containing the user data.
//
// Returns:
//   - A pointer to a response.UserResponseDeleteAt containing the mapped user data, including
//     ID, FirstName, LastName, Email, CreatedAt, UpdatedAt, and potentially DeletedAt, if available.
func (u *userQueryResponseMapper) ToResponseUserDeleteAt(user *pb.UserResponseDeleteAt) *response.UserResponseDeleteAt {
	var deletedAt string
	if user.DeletedAt != nil {
		deletedAt = user.DeletedAt.Value
	}

	return &response.UserResponseDeleteAt{
		ID:        int(user.Id),
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

// ToResponsesUserDeleteAt maps a slice of protobuf UserResponseDeleteAt to a slice of domain UserResponseDeleteAt.
//
// Args:
//   - users: A slice of pointers to pb.UserResponseDeleteAt to be mapped.
//
// Returns:
//   - A slice of pointers to response.UserResponseDeleteAt containing the mapped user data for each user, including
//     ID, FirstName, LastName, Email, CreatedAt, UpdatedAt, and potentially DeletedAt, if available.
func (u *userQueryResponseMapper) ToResponsesUserDeleteAt(users []*pb.UserResponseDeleteAt) []*response.UserResponseDeleteAt {
	var mappedUsers []*response.UserResponseDeleteAt

	for _, user := range users {
		mappedUsers = append(mappedUsers, u.ToResponseUserDeleteAt(user))
	}

	return mappedUsers
}
