package handler

import (
	"math"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
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

func mapToProtoOrderResponse(m interface{}) *pborder.OrderResponse {
	switch v := m.(type) {
	case *db.Order:
		return &pborder.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	case *db.GetOrdersRow:
		return &pborder.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	case *db.GetOrderByIDRow:
		return &pborder.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateOrderRow:
		return &pborder.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateOrderRow:
		return &pborder.OrderResponse{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoOrderResponseDeleteAt(m interface{}) *pborder.OrderResponseDeleteAt {
	var res *pborder.OrderResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *db.Order:
		res = &pborder.OrderResponseDeleteAt{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetOrdersActiveRow:
		res = &pborder.OrderResponseDeleteAt{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetOrdersTrashedRow:
		res = &pborder.OrderResponseDeleteAt{
			Id:         v.OrderID,
			MerchantId: v.MerchantID,
			UserId:     v.UserID,
			TotalPrice: int32(v.TotalPrice),
			CreatedAt:  formatTimestamp(v.CreatedAt),
			UpdatedAt:  formatTimestamp(v.UpdatedAt),
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

func mapToProtoOrderMonthlyTotalRevenueResponse(m interface{}) *pborder.OrderMonthlyTotalRevenueResponse {
	switch v := m.(type) {
	case *db.GetMonthlyTotalRevenueRow:
		return &pborder.OrderMonthlyTotalRevenueResponse{
			Year:         v.Year,
			Month:        v.Month,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetMonthlyTotalRevenueByMerchantRow:
		return &pborder.OrderMonthlyTotalRevenueResponse{
			Year:         v.Year,
			Month:        v.Month,
			TotalRevenue: int32(v.TotalRevenue),
		}
	default:
		return nil
	}
}

func mapToProtoOrderYearlyTotalRevenueResponse(m interface{}) *pborder.OrderYearlyTotalRevenueResponse {
	switch v := m.(type) {
	case *db.GetYearlyTotalRevenueRow:
		return &pborder.OrderYearlyTotalRevenueResponse{
			Year:         v.Year,
			TotalRevenue: int32(v.TotalRevenue),
		}
	case *db.GetYearlyTotalRevenueByMerchantRow:
		return &pborder.OrderYearlyTotalRevenueResponse{
			Year:         v.Year,
			TotalRevenue: int32(v.TotalRevenue),
		}
	default:
		return nil
	}
}

func mapToProtoOrderMonthlyResponse(m interface{}) *pborder.OrderMonthlyResponse {
	switch v := m.(type) {
	case *db.GetMonthlyOrderRow:
		return &pborder.OrderMonthlyResponse{
			Month:          v.Month,
			OrderCount:     int32(v.OrderCount),
			TotalRevenue:   int32(v.TotalRevenue),
			TotalItemsSold: int32(v.TotalItemsSold),
		}
	case *db.GetMonthlyOrderByMerchantRow:
		return &pborder.OrderMonthlyResponse{
			Month:          v.Month,
			OrderCount:     int32(v.OrderCount),
			TotalRevenue:   int32(v.TotalRevenue),
			TotalItemsSold: int32(v.TotalItemsSold),
		}
	default:
		return nil
	}
}

func mapToProtoOrderYearlyResponse(m interface{}) *pborder.OrderYearlyResponse {
	switch v := m.(type) {
	case *db.GetYearlyOrderRow:
		return &pborder.OrderYearlyResponse{
			Year:               v.Year,
			OrderCount:         int32(v.OrderCount),
			TotalRevenue:       int32(v.TotalRevenue),
			TotalItemsSold:     int32(v.TotalItemsSold),
			UniqueProductsSold: int32(v.UniqueProductsSold),
		}
	case *db.GetYearlyOrderByMerchantRow:
		return &pborder.OrderYearlyResponse{
			Year:               v.Year,
			OrderCount:         int32(v.OrderCount),
			TotalRevenue:       int32(v.TotalRevenue),
			TotalItemsSold:     int32(v.TotalItemsSold),
			UniqueProductsSold: int32(v.UniqueProductsSold),
		}
	default:
		return nil
	}
}
