package handler

import (
	"math"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
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

func mapToProtoCartResponse(m interface{}) *pb.CartResponse {
	switch v := m.(type) {
	case *db.Cart:
		return &pb.CartResponse{
			Id:        v.CartID,
			UserId:    v.UserID,
			ProductId: v.ProductID,
			Name:      v.Name,
			Price:     v.Price,
			Image:     v.Image,
			Quantity:  v.Quantity,
			Weight:    v.Weight,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetCartsRow:
		return &pb.CartResponse{
			Id:        v.CartID,
			UserId:    v.UserID,
			ProductId: v.ProductID,
			Name:      v.Name,
			Price:     v.Price,
			Image:     v.Image,
			Quantity:  v.Quantity,
			Weight:    v.Weight,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	default:
		return nil
	}
}
