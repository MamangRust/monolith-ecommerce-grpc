package roleapimapper

import (
	pb "github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

// RoleBaseResponseMapper defines a set of methods to map gRPC Role API responses
type RoleBaseResponseMapper interface {
	// ToApiResponseRole maps a single gRPC role response
	// into an HTTP API response format.
	ToApiResponseRole(pbResponse *pb.ApiResponseRole) *response.ApiResponseRole
}

type RoleQueryResponseMapper interface {
	RoleBaseResponseMapper

	// ToApiResponsesRole maps a gRPC response containing multiple roles
	// into a list HTTP API response format.
	ToApiResponsesRole(pbResponse *pb.ApiResponsesRole) *response.ApiResponsesRole

	// ToApiResponsePaginationRole maps a paginated gRPC response of roles
	// into a paginated HTTP API response format.
	ToApiResponsePaginationRole(pbResponse *pb.ApiResponsePaginationRole) *response.ApiResponsePaginationRole

	// ToApiResponsePaginationRoleDeleteAt maps a paginated gRPC response
	// of soft-deleted roles into a paginated HTTP API response format.
	ToApiResponsePaginationRoleDeleteAt(pbResponse *pb.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt
}

type RoleCommandResponseMapper interface {
	RoleBaseResponseMapper

	// ToApiResponseRoleDelete maps a gRPC delete role response
	// into an HTTP API response format.
	ToApiResponseRoleDelete(pbResponse *pb.ApiResponseRoleDelete) *response.ApiResponseRoleDelete

	// ToApiResponseRoleAll maps a gRPC response containing all roles
	// into an HTTP API response format.
	ToApiResponseRoleAll(pbResponse *pb.ApiResponseRoleAll) *response.ApiResponseRoleAll
}

type RoleResponseMapper interface {
	QueryMapper() RoleQueryResponseMapper
	CommandMapper() RoleCommandResponseMapper
}
