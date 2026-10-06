package handler

import (
	"math"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pbuser "github.com/MamangRust/monolith-ecommerce-pb/user"
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

func createPaginationMeta(page, pageSize, totalRecords int) *pbcommon.PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoUserResponse(m interface{}) *pbuser.UserResponse {
	switch v := m.(type) {
	case *db.User:
		return &pbuser.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.GetUsersRow:
		return &pbuser.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.GetUserByIDRow:
		return &pbuser.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.CreateUserRow:
		return &pbuser.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	default:
		return nil
	}
}

func mapToProtoUserResponseDeleteAt(m interface{}) *pbuser.UserResponseDeleteAt {
	switch v := m.(type) {
	case *db.User:
		return &pbuser.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.GetUsersActiveRow:
		return &pbuser.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.GetUserTrashedRow:
		return &pbuser.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.TrashUserRow:
		return &pbuser.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.RestoreUserRow:
		return &pbuser.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	default:
		return nil
	}
}
func mapToProtoUserResponseWithPassword(m interface{}) *pbuser.UserResponseWithPassword {
	switch v := m.(type) {
	case *db.User:
		return &pbuser.UserResponseWithPassword{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			Password:  v.Password,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.GetUserByEmailWithPasswordRow:
		return &pbuser.UserResponseWithPassword{
			Id:       int32(v.UserID),
			Email:    v.Email,
			Password: v.Password,
		}
	default:
		return nil
	}
}
