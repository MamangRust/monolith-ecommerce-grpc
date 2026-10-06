package handler

import (
	"math"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"

	pbcommon "github.com/MamangRust/monolith-ecommerce-pb/common"
	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
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

func mapToProtoOrderItemResponse(item interface{}) *pborder_item.OrderItemResponse {
	switch v := item.(type) {
	case *db.OrderItem:
		return &pborder_item.OrderItemResponse{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateOrderItemRow:
		return &pborder_item.OrderItemResponse{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateOrderItemRow:
		return &pborder_item.OrderItemResponse{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *db.GetOrderItemsRow:
		return &pborder_item.OrderItemResponse{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	case *db.GetOrderItemsByOrderRow:
		return &pborder_item.OrderItemResponse{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoOrderItemResponseDeleteAt(item interface{}) *pborder_item.OrderItemResponseDeleteAt {
	var res *pborder_item.OrderItemResponseDeleteAt
	var deletedAt interface{}

	switch v := item.(type) {
	case *db.OrderItem:
		res = &pborder_item.OrderItemResponseDeleteAt{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetOrderItemsActiveRow:
		res = &pborder_item.OrderItemResponseDeleteAt{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
			CreatedAt: formatTimestamp(v.CreatedAt),
			UpdatedAt: formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetOrderItemsTrashedRow:
		res = &pborder_item.OrderItemResponseDeleteAt{
			Id:        v.OrderItemID,
			OrderId:   v.OrderID,
			ProductId: v.ProductID,
			Quantity:  v.Quantity,
			Price:     v.Price,
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
