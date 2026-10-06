package roleapimapper

import (
	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
)

type roleCommandResponseMapper struct {
}

func NewRoleCommandResponseMapper() RoleCommandResponseMapper {
	return &roleCommandResponseMapper{}
}

func (s *roleCommandResponseMapper) ToApiResponseRole(pbResponse *pbrole.ApiResponseRole) *response.ApiResponseRole {
	return &response.ApiResponseRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponseRole(pbResponse.Data),
	}
}

func (s *roleCommandResponseMapper) ToApiResponseRoleDelete(pbResponse *pbrole.ApiResponseRoleDelete) *response.ApiResponseRoleDelete {
	return &response.ApiResponseRoleDelete{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleCommandResponseMapper) ToApiResponseRoleAll(pbResponse *pbrole.ApiResponseRoleAll) *response.ApiResponseRoleAll {
	return &response.ApiResponseRoleAll{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
	}
}

func (s *roleCommandResponseMapper) mapResponseRole(role *pbrole.RoleResponse) *response.RoleResponse {
	if role == nil {
		return nil
	}
	return &response.RoleResponse{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}
