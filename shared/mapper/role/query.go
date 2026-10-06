package roleapimapper

import (
	pbrole "github.com/MamangRust/monolith-ecommerce-pb/role"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	paginationapimapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/pagination"
)

type roleQueryResponseMapper struct{}

func NewRoleQueryResponseMapper() RoleQueryResponseMapper {
	return &roleQueryResponseMapper{}
}

func (s *roleQueryResponseMapper) ToApiResponseRole(pbResponse *pbrole.ApiResponseRole) *response.ApiResponseRole {
	return &response.ApiResponseRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponseRole(pbResponse.Data),
	}
}

func (s *roleQueryResponseMapper) ToApiResponsesRole(pbResponse *pbrole.ApiResponsesRole) *response.ApiResponsesRole {
	return &response.ApiResponsesRole{
		Status:  pbResponse.Status,
		Message: pbResponse.Message,
		Data:    s.mapResponsesRole(pbResponse.Data),
	}
}

func (s *roleQueryResponseMapper) ToApiResponsePaginationRole(pbResponse *pbrole.ApiResponsePaginationRole) *response.ApiResponsePaginationRole {
	return &response.ApiResponsePaginationRole{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRole(pbResponse.Data),
		Pagination: paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *roleQueryResponseMapper) ToApiResponsePaginationRoleDeleteAt(pbResponse *pbrole.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt {
	return &response.ApiResponsePaginationRoleDeleteAt{
		Status:     pbResponse.Status,
		Message:    pbResponse.Message,
		Data:       s.mapResponsesRoleDeleteAt(pbResponse.Data),
		Pagination: paginationapimapper.MapPaginationMeta(pbResponse.Pagination),
	}
}

func (s *roleQueryResponseMapper) mapResponseRole(role *pbrole.RoleResponse) *response.RoleResponse {
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

func (s *roleQueryResponseMapper) mapResponsesRole(roles []*pbrole.RoleResponse) []*response.RoleResponse {
	var responseRoles []*response.RoleResponse
	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRole(role))
	}
	return responseRoles
}

func (s *roleQueryResponseMapper) mapResponseRoleDeleteAt(role *pbrole.RoleResponseDeleteAt) *response.RoleResponseDeleteAt {
	if role == nil {
		return nil
	}
	var deletedAt string
	if role.DeletedAt != nil {
		deletedAt = role.DeletedAt.Value
	}
	return &response.RoleResponseDeleteAt{
		ID:        int(role.Id),
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		DeletedAt: &deletedAt,
	}
}

func (s *roleQueryResponseMapper) mapResponsesRoleDeleteAt(roles []*pbrole.RoleResponseDeleteAt) []*response.RoleResponseDeleteAt {
	var responseRoles []*response.RoleResponseDeleteAt
	for _, role := range roles {
		responseRoles = append(responseRoles, s.mapResponseRoleDeleteAt(role))
	}
	return responseRoles
}
