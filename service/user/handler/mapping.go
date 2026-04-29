package handler

import (
	"math"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
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

func mapToProtoUserResponse(m interface{}) *pb.UserResponse {
	switch v := m.(type) {
	case *db.User:
		return &pb.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.GetUsersRow:
		return &pb.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.GetUserByIDRow:
		return &pb.UserResponse{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.CreateUserRow:
		return &pb.UserResponse{
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

func mapToProtoUserResponseDeleteAt(m interface{}) *pb.UserResponseDeleteAt {
	switch v := m.(type) {
	case *db.User:
		return &pb.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.GetUsersActiveRow:
		return &pb.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.GetUserTrashedRow:
		return &pb.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.TrashUserRow:
		return &pb.UserResponseDeleteAt{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
			DeletedAt: &wrapperspb.StringValue{Value: v.DeletedAt.Time.Format(time.RFC3339)},
		}
	case *db.RestoreUserRow:
		return &pb.UserResponseDeleteAt{
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
func mapToProtoUserResponseWithPassword(m interface{}) *pb.UserResponseWithPassword {
	switch v := m.(type) {
	case *db.User:
		return &pb.UserResponseWithPassword{
			Id:        int32(v.UserID),
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Email:     v.Email,
			Password:  v.Password,
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		}
	case *db.GetUserByEmailWithPasswordRow:
		return &pb.UserResponseWithPassword{
			Id:        int32(v.UserID),
			Email:     v.Email,
			Password:  v.Password,
		}
	default:
		return nil
	}
}
