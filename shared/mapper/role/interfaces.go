package roleapimapper

import (
	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

// RoleBaseResponseMapper defines a set of methods to map gRPC Role API responses
type RoleBaseResponseMapper interface {
	// ToApiResponseRole maps a single gRPC role response
	// into an HTTP API response format.
	ToApiResponseRole(pbResponse *pbrole.ApiResponseRole) *response.ApiResponseRole
}

type RoleQueryResponseMapper interface {
	RoleBaseResponseMapper

	// ToApiResponsesRole maps a gRPC response containing multiple roles
	// into a list HTTP API response format.
	ToApiResponsesRole(pbResponse *pbrole.ApiResponsesRole) *response.ApiResponsesRole

	// ToApiResponsePaginationRole maps a paginated gRPC response of roles
	// into a paginated HTTP API response format.
	ToApiResponsePaginationRole(pbResponse *pbrole.ApiResponsePaginationRole) *response.ApiResponsePaginationRole

	// ToApiResponsePaginationRoleDeleteAt maps a paginated gRPC response
	// of soft-deleted roles into a paginated HTTP API response format.
	ToApiResponsePaginationRoleDeleteAt(pbResponse *pbrole.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt
}

type RoleCommandResponseMapper interface {
	RoleBaseResponseMapper

	// ToApiResponseRoleDelete maps a gRPC delete role response
	// into an HTTP API response format.
	ToApiResponseRoleDelete(pbResponse *pbrole.ApiResponseRoleDelete) *response.ApiResponseRoleDelete

	// ToApiResponseRoleAll maps a gRPC response containing all roles
	// into an HTTP API response format.
	ToApiResponseRoleAll(pbResponse *pbrole.ApiResponseRoleAll) *response.ApiResponseRoleAll
}

type RoleResponseMapper interface {
	QueryMapper() RoleQueryResponseMapper
	CommandMapper() RoleCommandResponseMapper
}
