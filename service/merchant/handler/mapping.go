package handler

import (
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

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

func mapToProtoMerchantResponse(m interface{}) *pb.MerchantResponse {
	switch v := m.(type) {
	case *db.Merchant:
		return &pb.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantsRow:
		return &pb.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantByIDRow:
		return &pb.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantRow:
		return &pb.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantRow:
		return &pb.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantStatusRow:
		return &pb.MerchantResponse{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantResponseDeleteAt(m interface{}) *pb.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *db.Merchant:
		return &pb.MerchantResponseDeleteAt{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	case *db.GetMerchantsActiveRow:
		return &pb.MerchantResponseDeleteAt{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantResponseTrashed(m interface{}) *pb.MerchantResponseDeleteAt {
	switch v := m.(type) {
	case *db.GetMerchantsTrashedRow:
		return &pb.MerchantResponseDeleteAt{
			Id:           int32(v.MerchantID),
			UserId:       int32(v.UserID),
			Name:         v.Name,
			Description:  getString(v.Description),
			Address:      getString(v.Address),
			ContactEmail: getString(v.ContactEmail),
			ContactPhone: getString(v.ContactPhone),
			Status:       v.Status,
			CreatedAt:    formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponse(m interface{}) *pb.MerchantDocument {
	switch v := m.(type) {
	case *db.MerchantDocument:
		return &pb.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantDocumentsRow:
		return &pb.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetActiveMerchantDocumentsRow:
		return &pb.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantDocumentRow:
		return &pb.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantDocumentRow:
		return &pb.MerchantDocument{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantDocumentResponseAt(m interface{}) *pb.MerchantDocumentDeleteAt {
	switch v := m.(type) {
	case *db.GetTrashedMerchantDocumentsRow:
		return &pb.MerchantDocumentDeleteAt{
			DocumentId:   int32(v.DocumentID),
			MerchantId:   int32(v.MerchantID),
			DocumentType: v.DocumentType,
			DocumentUrl:  v.DocumentUrl,
			Status:       v.Status,
			Note:         getString(v.Note),
			UploadedAt:   formatTimestamp(v.CreatedAt),
			UpdatedAt:    formatTimestamp(v.UpdatedAt),
			DeletedAt:    getStringValue(v.DeletedAt),
		}
	default:
		return nil
	}
}
