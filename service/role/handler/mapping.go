package handler

import (
	"math"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func createPaginationMeta(page, pageSize, totalRecords int) *pb.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func formatTimestamp(v interface{}) string {
	switch t := v.(type) {
	case pgtype.Timestamptz:
		if t.Valid {
			return t.Time.Format("2006-01-02 15:04:05.000")
		}
	case pgtype.Timestamp:
		if t.Valid {
			return t.Time.Format("2006-01-02 15:04:05.000")
		}
	}
	return ""
}

func mapToProtoRoleResponse(m interface{}) *pb.RoleResponse {
	switch v := m.(type) {
	case *db.Role:
		return &pb.RoleResponse{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *db.GetRolesRow:
		return &pb.RoleResponse{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoRoleResponseDeleteAt(m interface{}) *pb.RoleResponseDeleteAt {
	var res *pb.RoleResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *db.Role:
		res = &pb.RoleResponseDeleteAt{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetActiveRolesRow:
		res = &pb.RoleResponseDeleteAt{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetTrashedRolesRow:
		res = &pb.RoleResponseDeleteAt{
			Id:        v.RoleID,
			Name:      v.RoleName,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	default:
		return nil
	}

	if val := formatTimestamp(deletedAt); val != "" {
		res.DeletedAt = &wrapperspb.StringValue{Value: val}
	}

	return res
}
func mapToProtoUserRoleResponse(v *db.UserRole) *pb.UserRoleResponse {
	if v == nil {
		return nil
	}
	return &pb.UserRoleResponse{
		UserRoleId: v.UserRoleID,
		UserId:     v.UserID,
		RoleId:     v.RoleID,
		CreatedAt:  formatTimestamp(v.CreatedAt),
		UpdatedAt:  formatTimestamp(v.UpdatedAt),
	}
}
