package handler

import (
	"encoding/json"
	"log"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (h *Handler) mapToReviewDetailResponse(data interface{}) interface{} {
	switch v := data.(type) {
	case *db.GetReviewDetailRow:
		return &pb.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetReviewDetailsRow:
		return &pb.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.GetReviewDetailsActiveRow:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pb.ReviewDetailsResponseDeleteAt{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	case *db.GetReviewDetailsTrashedRow:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pb.ReviewDetailsResponseDeleteAt{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	case *db.CreateReviewDetailRow:
		return &pb.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.UpdateReviewDetailRow:
		return &pb.ReviewDetailsResponse{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
		}
	case *db.ReviewDetail:
		var deletedAt string
		if v.DeletedAt.Valid {
			deletedAt = v.DeletedAt.Time.Format("2006-01-02")
		}
		return &pb.ReviewDetailsResponseDeleteAt{
			Id:        int32(v.ReviewDetailID),
			ReviewId:  int32(v.ReviewID),
			Type:      v.Type,
			Url:       v.Url,
			Caption:   *v.Caption,
			CreatedAt: v.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt: v.UpdatedAt.Time.Format("2006-01-02"),
			DeletedAt: &wrapperspb.StringValue{Value: deletedAt},
		}
	default:
		log.Printf("Unknown type for mapping: %T", v)
		return nil
	}
}

func (h *Handler) mapToPayload(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
