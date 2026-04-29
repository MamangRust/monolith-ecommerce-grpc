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

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
	case pgtype.Date:
		if t.Valid {
			return t.Time.Format("2006-01-02")
		}
	}
	return ""
}

func mapToProtoMerchantAwardResponse(m interface{}) *pb.MerchantAwardResponse {
	switch v := m.(type) {
	case *db.MerchantCertificationsAndAward:
		return &pb.MerchantAwardResponse{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantCertificationsAndAwardsRow:
		return &pb.MerchantAwardResponse{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
	case *db.CreateMerchantCertificationOrAwardRow:
		return &pb.MerchantAwardResponse{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
	case *db.UpdateMerchantCertificationOrAwardRow:
		return &pb.MerchantAwardResponse{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
	case *db.GetMerchantCertificationOrAwardRow:
		return &pb.MerchantAwardResponse{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
	default:
		return nil
	}
}

func mapToProtoMerchantAwardResponseDeleteAt(m interface{}) *pb.MerchantAwardResponseDeleteAt {
	var res *pb.MerchantAwardResponseDeleteAt
	var deletedAt interface{}

	switch v := m.(type) {
	case *db.MerchantCertificationsAndAward:
		res = &pb.MerchantAwardResponseDeleteAt{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetMerchantCertificationsAndAwardsActiveRow:
		res = &pb.MerchantAwardResponseDeleteAt{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
		}
		deletedAt = v.DeletedAt
	case *db.GetMerchantCertificationsAndAwardsTrashedRow:
		res = &pb.MerchantAwardResponseDeleteAt{
			Id:             v.MerchantCertificationID,
			MerchantId:     v.MerchantID,
			Title:          v.Title,
			Description:    getString(v.Description),
			IssuedBy:       getString(v.IssuedBy),
			CertificateUrl: getString(v.CertificateUrl),
			IssueDate:      formatTimestamp(v.IssueDate),
			ExpiryDate:     formatTimestamp(v.ExpiryDate),
			CreatedAt:      formatTimestamp(v.CreatedAt),
			UpdatedAt:      formatTimestamp(v.UpdatedAt),
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
