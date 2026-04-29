package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func getStringValue(v interface{}) *wrapperspb.StringValue {
	if ts, ok := v.(pgtype.Timestamptz); ok && ts.Valid {
		return wrapperspb.String(ts.Time.Format("2006-01-02 15:04:05"))
	}
	if ts, ok := v.(pgtype.Timestamp); ok && ts.Valid {
		return wrapperspb.String(ts.Time.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func formatTimestamp(v interface{}) string {
	if ts, ok := v.(pgtype.Timestamptz); ok && ts.Valid {
		return ts.Time.Format("2006-01-02 15:04:05")
	}
	if ts, ok := v.(pgtype.Timestamp); ok && ts.Valid {
		return ts.Time.Format("2006-01-02 15:04:05")
	}
	return ""
}

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
	totalPages := (totalRecords + pageSize - 1) / pageSize
	return &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecords),
	}
}

func mapToProtoTransactionResponse(m interface{}) *pb.TransactionResponse {
	switch v := m.(type) {
	case *db.Transaction:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionsRow:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionsActiveRow:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionByIDRow:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionByOrderIDRow:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.GetTransactionByMerchantRow:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateTransactionRow:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateTransactionRow:
		return &pb.TransactionResponse{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoTransactionResponseDeleteAt(m interface{}) *pb.TransactionResponseDeleteAt {
	switch v := m.(type) {
	case *db.Transaction:
		return &pb.TransactionResponseDeleteAt{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
			DeletedAt:     getStringValue(v.DeletedAt),
		}
	case *db.GetTransactionsActiveRow:
		return &pb.TransactionResponseDeleteAt{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
			DeletedAt:     getStringValue(v.DeletedAt),
		}
	case *db.GetTransactionsTrashedRow:
		return &pb.TransactionResponseDeleteAt{
			Id:            v.TransactionID,
			OrderId:       v.OrderID,
			MerchantId:    v.MerchantID,
			PaymentMethod: v.PaymentMethod,
			Amount:        v.Amount,
			PaymentStatus: v.PaymentStatus,
			CreatedAt:     formatTimestamp(v.CreatedAt),
			UpdatedAt:     formatTimestamp(v.UpdatedAt),
			DeletedAt:     getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}
